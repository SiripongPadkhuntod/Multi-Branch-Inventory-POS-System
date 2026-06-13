package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"pos-system/backend/internal/app/domain"
	"pos-system/backend/internal/app/port"
	"pos-system/backend/internal/property"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo port.Repository
	cfg  property.AppConfig
}

func New(repo port.Repository) port.Service {
	cfg := property.Load()
	return &service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *service) Health(context.Context) error {
	return nil
}

func (s *service) Authenticate(ctx context.Context, authorization string) (*domain.User, error) {
	if strings.TrimSpace(authorization) == "" {
		return nil, errors.New("missing authorization header")
	}
	tokenString, ok := strings.CutPrefix(authorization, "Bearer ")
	if !ok || strings.TrimSpace(tokenString) == "" {
		return nil, errors.New("invalid authorization header")
	}
	claims, err := s.parseToken(strings.TrimSpace(tokenString))
	if err != nil || (claims.TokenUse != "" && claims.TokenUse != "access") {
		return nil, errors.New("invalid access token")
	}
	user, err := s.repo.SQL.Auth().FindUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Status != "ACTIVE" {
		return nil, errors.New("user is inactive")
	}
	return user, nil
}

func (s *service) Login(ctx context.Context, email, password string) (string, string, *domain.User, error) {
	user, err := s.repo.SQL.Auth().FindUserByEmail(ctx, email)
	if err != nil {
		return "", "", nil, errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", nil, errors.New("invalid email or password")
	}
	access, err := s.sign(user, s.cfg.AccessTTL, "access")
	if err != nil {
		return "", "", nil, err
	}
	refresh, err := s.sign(user, s.cfg.RefreshTTL, "refresh")
	if err != nil {
		return "", "", nil, err
	}
	if err := s.repo.SQL.Auth().StoreRefreshToken(ctx, user.ID, tokenHash(refresh), time.Now().Add(s.cfg.RefreshTTL)); err != nil {
		return "", "", nil, err
	}
	return access, refresh, user, nil
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (string, string, *domain.User, error) {
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return "", "", nil, errors.New("invalid refresh token")
	}
	if claims.TokenUse != "refresh" {
		return "", "", nil, errors.New("invalid refresh token")
	}
	refreshHash := tokenHash(refreshToken)
	active, err := s.repo.SQL.Auth().IsRefreshTokenActive(ctx, claims.UserID, refreshHash)
	if err != nil {
		return "", "", nil, err
	}
	if !active {
		return "", "", nil, errors.New("invalid refresh token")
	}
	user, err := s.repo.SQL.Auth().FindUserByID(ctx, claims.UserID)
	if err != nil {
		return "", "", nil, errors.New("user not found")
	}
	if user.Status != "ACTIVE" {
		return "", "", nil, errors.New("user is inactive")
	}
	if err := s.repo.SQL.Auth().RevokeRefreshToken(ctx, refreshHash); err != nil {
		return "", "", nil, err
	}
	access, err := s.sign(user, s.cfg.AccessTTL, "access")
	if err != nil {
		return "", "", nil, err
	}
	refresh, err := s.sign(user, s.cfg.RefreshTTL, "refresh")
	if err != nil {
		return "", "", nil, err
	}
	if err := s.repo.SQL.Auth().StoreRefreshToken(ctx, user.ID, tokenHash(refresh), time.Now().Add(s.cfg.RefreshTTL)); err != nil {
		return "", "", nil, err
	}
	return access, refresh, user, nil
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	claims, err := s.parseToken(refreshToken)
	if err != nil || claims.TokenUse != "refresh" {
		return nil
	}
	return s.repo.SQL.Auth().RevokeRefreshToken(ctx, tokenHash(refreshToken))
}

func (s *service) ListProducts(ctx context.Context, query string, limit, offset int) ([]domain.Product, error) {
	return s.repo.SQL.Products().List(ctx, query, limit, offset)
}

func (s *service) ProductByBarcode(ctx context.Context, barcode string) (*domain.Product, error) {
	return s.repo.SQL.Products().FindByBarcode(ctx, barcode)
}

func (s *service) ListInventories(ctx context.Context, branchID *uuid.UUID, categoryID *uuid.UUID, query string) ([]domain.Inventory, error) {
	return s.repo.SQL.Inventories().List(ctx, branchID, categoryID, query)
}

func (s *service) ListInventoryMovements(ctx context.Context, branchID *uuid.UUID, query string, limit int) ([]domain.InventoryMovementDetail, error) {
	return s.repo.SQL.Inventories().ListMovements(ctx, branchID, query, limit)
}

func (s *service) AllStock(ctx context.Context, query string) ([]domain.ProductStockSummary, error) {
	return s.repo.SQL.Inventories().AllStock(ctx, query)
}

func (s *service) AdjustInventory(ctx context.Context, branchID, productID, actorID uuid.UUID, delta int64, reason string) error {
	return s.repo.SQL.Inventories().Adjust(ctx, branchID, productID, actorID, delta, reason)
}

type claims struct {
	UserID   uuid.UUID   `json:"user_id"`
	Role     domain.Role `json:"role"`
	BranchID *uuid.UUID  `json:"branch_id"`
	TokenUse string      `json:"token_use"`
	jwt.RegisteredClaims
}

func (s *service) parseToken(tokenString string) (*claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	parsedClaims, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if parsedClaims.TokenUse != "" && parsedClaims.TokenUse != "access" && parsedClaims.TokenUse != "refresh" {
		return nil, errors.New("invalid token")
	}
	return parsedClaims, nil
}

func (s *service) sign(user *domain.User, ttl time.Duration, tokenUse string) (string, error) {
	now := time.Now()
	claims := claims{
		UserID:   user.ID,
		Role:     user.Role,
		BranchID: user.BranchID,
		TokenUse: tokenUse,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWTSecret))
}

func tokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
