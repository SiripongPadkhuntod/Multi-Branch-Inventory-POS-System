package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"pos-system/backend/internal/app/domain"
	"pos-system/backend/internal/app/port"
	"pos-system/backend/internal/infrastructure/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Services struct {
	Auth      *AuthService
	Products  *ProductService
	Inventory *InventoryService
	Sales     *SaleService
	Users     *UserService
	Dashboard *DashboardService
	System    *SystemService
	Audit     *AuditService
}

func NewServices(repos port.Repositories, cfg config.Config) *Services {
	return &Services{
		Auth:      &AuthService{repo: repos.Auth(), cfg: cfg},
		Products:  &ProductService{repo: repos.Products()},
		Inventory: &InventoryService{repo: repos.Inventories()},
		Sales:     &SaleService{repo: repos.Sales()},
		Users:     &UserService{repo: repos.Users()},
		Dashboard: &DashboardService{repo: repos.Dashboard()},
		System:    &SystemService{repo: repos.System()},
		Audit:     &AuditService{repo: repos.Audit()},
	}
}

type Claims struct {
	UserID   uuid.UUID   `json:"user_id"`
	Role     domain.Role `json:"role"`
	BranchID *uuid.UUID  `json:"branch_id"`
	TokenUse string      `json:"token_use"`
	jwt.RegisteredClaims
}

type AuthService struct {
	repo port.AuthRepository
	cfg  config.Config
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, *domain.User, error) {
	user, err := s.repo.FindUserByEmail(ctx, email)
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
	if err := s.repo.StoreRefreshToken(ctx, user.ID, tokenHash(refresh), time.Now().Add(s.cfg.RefreshTTL)); err != nil {
		return "", "", nil, err
	}
	return access, refresh, user, nil
}

func (s *AuthService) FindUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.FindUserByID(ctx, id)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, string, *domain.User, error) {
	claims, err := s.ParseToken(refreshToken)
	if err != nil {
		return "", "", nil, errors.New("invalid refresh token")
	}
	if claims.TokenUse != "refresh" {
		return "", "", nil, errors.New("invalid refresh token")
	}
	refreshHash := tokenHash(refreshToken)
	active, err := s.repo.IsRefreshTokenActive(ctx, claims.UserID, refreshHash)
	if err != nil {
		return "", "", nil, err
	}
	if !active {
		return "", "", nil, errors.New("invalid refresh token")
	}
	user, err := s.FindUser(ctx, claims.UserID)
	if err != nil {
		return "", "", nil, errors.New("user not found")
	}
	if user.Status != "ACTIVE" {
		return "", "", nil, errors.New("user is inactive")
	}
	if err := s.repo.RevokeRefreshToken(ctx, refreshHash); err != nil {
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
	if err := s.repo.StoreRefreshToken(ctx, user.ID, tokenHash(refresh), time.Now().Add(s.cfg.RefreshTTL)); err != nil {
		return "", "", nil, err
	}
	return access, refresh, user, nil
}

func (s *AuthService) RevokeRefresh(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	claims, err := s.ParseToken(refreshToken)
	if err != nil || claims.TokenUse != "refresh" {
		return nil
	}
	return s.repo.RevokeRefreshToken(ctx, tokenHash(refreshToken))
}

func (s *AuthService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.TokenUse != "" && claims.TokenUse != "access" && claims.TokenUse != "refresh" {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *AuthService) sign(user *domain.User, ttl time.Duration, tokenUse string) (string, error) {
	now := time.Now()
	claims := Claims{
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

type ProductService struct{ repo port.ProductRepository }

func (s *ProductService) List(ctx context.Context, query string, limit, offset int) ([]domain.Product, error) {
	return s.repo.List(ctx, query, limit, offset)
}
func (s *ProductService) ByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.repo.FindByID(ctx, id)
}
func (s *ProductService) ByBarcode(ctx context.Context, barcode string) (*domain.Product, error) {
	return s.repo.FindByBarcode(ctx, barcode)
}
func (s *ProductService) Create(ctx context.Context, p domain.Product) (*domain.Product, error) {
	return s.repo.Create(ctx, p)
}
func (s *ProductService) Update(ctx context.Context, id uuid.UUID, p domain.Product) (*domain.Product, error) {
	return s.repo.Update(ctx, id, p)
}
func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.SoftDelete(ctx, id)
}

type InventoryService struct {
	repo port.InventoryRepository
}

func (s *InventoryService) List(ctx context.Context, branchID *uuid.UUID, categoryID *uuid.UUID, query string) ([]domain.Inventory, error) {
	return s.repo.List(ctx, branchID, categoryID, query)
}
func (s *InventoryService) Movements(ctx context.Context, branchID *uuid.UUID, query string, limit int) ([]domain.InventoryMovementDetail, error) {
	return s.repo.ListMovements(ctx, branchID, query, limit)
}
func (s *InventoryService) AllStock(ctx context.Context, query string) ([]domain.ProductStockSummary, error) {
	return s.repo.AllStock(ctx, query)
}
func (s *InventoryService) Adjust(ctx context.Context, branchID, productID, actorID uuid.UUID, delta int64, reason string) error {
	return s.repo.Adjust(ctx, branchID, productID, actorID, delta, reason)
}
func (s *InventoryService) SetReorderThreshold(ctx context.Context, branchID, productID uuid.UUID, threshold int64) error {
	return s.repo.SetReorderThreshold(ctx, branchID, productID, threshold)
}
func (s *InventoryService) CreateTransfer(ctx context.Context, fromBranchID, toBranchID, productID, actorID uuid.UUID, quantity int64) (*domain.Transfer, error) {
	return s.repo.CreateTransfer(ctx, fromBranchID, toBranchID, productID, actorID, quantity)
}
func (s *InventoryService) Transfers(ctx context.Context, status string, limit int) ([]domain.Transfer, error) {
	return s.repo.ListTransfers(ctx, status, limit)
}
func (s *InventoryService) ApproveTransfer(ctx context.Context, transferID, actorID uuid.UUID) (*domain.Transfer, error) {
	return s.repo.ApproveTransfer(ctx, transferID, actorID)
}
func (s *InventoryService) RejectTransfer(ctx context.Context, transferID, actorID uuid.UUID) (*domain.Transfer, error) {
	return s.repo.RejectTransfer(ctx, transferID, actorID)
}
func (s *InventoryService) CompleteTransfer(ctx context.Context, transferID, actorID uuid.UUID) (*domain.Transfer, error) {
	return s.repo.CompleteTransfer(ctx, transferID, actorID)
}

type SaleService struct{ repo port.SaleRepository }

func (s *SaleService) Create(ctx context.Context, actor domain.User, input domain.CreateSaleInput) (*domain.Sale, error) {
	return s.repo.CreateSale(ctx, actor, input)
}
func (s *SaleService) List(ctx context.Context, actor domain.User, branchID *uuid.UUID, dateFrom, dateTo *time.Time) ([]domain.Sale, error) {
	return s.repo.List(ctx, actor, branchID, dateFrom, dateTo)
}
func (s *SaleService) BranchList(ctx context.Context, actor domain.User, branchID *uuid.UUID) ([]domain.Sale, error) {
	return s.repo.ListBranch(ctx, actor, branchID)
}
func (s *SaleService) Detail(ctx context.Context, actor domain.User, saleID uuid.UUID) (*domain.SaleDetail, error) {
	return s.repo.FindDetail(ctx, actor, saleID)
}
func (s *SaleService) Refund(ctx context.Context, actor domain.User, saleID uuid.UUID, items []domain.CartItemInput) error {
	return s.repo.Refund(ctx, actor, saleID, items)
}

type UserService struct{ repo port.UserRepository }

func (s *UserService) List(ctx context.Context, actor domain.User) ([]domain.User, error) {
	return s.repo.List(ctx, actor)
}
func (s *UserService) Create(ctx context.Context, actor domain.User, user domain.User, password string) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(hash)
	if user.Status == "" {
		user.Status = "ACTIVE"
	}
	return s.repo.Create(ctx, actor, user)
}
func (s *UserService) Update(ctx context.Context, actor domain.User, id uuid.UUID, user domain.User) (*domain.User, error) {
	return s.repo.Update(ctx, actor, id, user)
}
func (s *UserService) SalesSummary(ctx context.Context, actor domain.User) ([]domain.EmployeeSalesSummary, error) {
	return s.repo.SalesSummary(ctx, actor)
}

type DashboardService struct {
	repo port.DashboardRepository
}

func (s *DashboardService) AccessibleBranches(ctx context.Context, actor domain.User) ([]domain.Branch, error) {
	return s.repo.AccessibleBranches(ctx, actor)
}

func (s *DashboardService) Summary(ctx context.Context, actor domain.User, branchID *uuid.UUID) (*domain.DashboardSummary, error) {
	return s.repo.Summary(ctx, actor, branchID)
}

type SystemService struct {
	repo port.SystemRepository
}

func (s *SystemService) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return s.repo.ListCategories(ctx)
}
func (s *SystemService) CreateCategory(ctx context.Context, category domain.Category) (*domain.Category, error) {
	return s.repo.CreateCategory(ctx, category)
}
func (s *SystemService) UpdateCategory(ctx context.Context, id uuid.UUID, category domain.Category) (*domain.Category, error) {
	return s.repo.UpdateCategory(ctx, id, category)
}
func (s *SystemService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteCategory(ctx, id)
}
func (s *SystemService) CreateBranch(ctx context.Context, branch domain.Branch) (*domain.Branch, error) {
	return s.repo.CreateBranch(ctx, branch)
}
func (s *SystemService) UpdateBranch(ctx context.Context, id uuid.UUID, branch domain.Branch) (*domain.Branch, error) {
	return s.repo.UpdateBranch(ctx, id, branch)
}

type AuditService struct {
	repo port.AuditRepository
}

func (s *AuditService) List(ctx context.Context, action, entityType, query string, limit int) ([]domain.AuditLog, error) {
	return s.repo.List(ctx, action, entityType, query, limit)
}
