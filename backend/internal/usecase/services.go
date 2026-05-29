package usecase

import (
	"context"
	"errors"
	"time"

	"pos-system/backend/internal/domain"
	"pos-system/backend/internal/infrastructure/config"
	"pos-system/backend/internal/repository"

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
}

func NewServices(repos repository.Repositories, cfg config.Config) *Services {
	return &Services{
		Auth:      &AuthService{repo: repos.Auth(), cfg: cfg},
		Products:  &ProductService{repo: repos.Products()},
		Inventory: &InventoryService{repo: repos.Inventories()},
		Sales:     &SaleService{repo: repos.Sales()},
		Users:     &UserService{repo: repos.Users()},
		Dashboard: &DashboardService{repo: repos.Dashboard()},
		System:    &SystemService{repo: repos.System()},
	}
}

type Claims struct {
	UserID   uuid.UUID   `json:"user_id"`
	Role     domain.Role `json:"role"`
	BranchID *uuid.UUID  `json:"branch_id"`
	jwt.RegisteredClaims
}

type AuthService struct {
	repo repository.AuthRepository
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
	access, err := s.sign(user, s.cfg.AccessTTL)
	if err != nil {
		return "", "", nil, err
	}
	refresh, err := s.sign(user, s.cfg.RefreshTTL)
	return access, refresh, user, err
}

func (s *AuthService) FindUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.FindUserByID(ctx, id)
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
	return claims, nil
}

func (s *AuthService) sign(user *domain.User, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   user.ID,
		Role:     user.Role,
		BranchID: user.BranchID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWTSecret))
}

type ProductService struct{ repo repository.ProductRepository }

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
	repo repository.InventoryRepository
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
func (s *InventoryService) Transfer(ctx context.Context, fromBranchID, toBranchID, productID, actorID uuid.UUID, quantity int64) error {
	return s.repo.Transfer(ctx, fromBranchID, toBranchID, productID, actorID, quantity)
}

type SaleService struct{ repo repository.SaleRepository }

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

type UserService struct{ repo repository.UserRepository }

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
	repo repository.DashboardRepository
}

func (s *DashboardService) AccessibleBranches(ctx context.Context, actor domain.User) ([]domain.Branch, error) {
	return s.repo.AccessibleBranches(ctx, actor)
}

func (s *DashboardService) Summary(ctx context.Context, actor domain.User, branchID *uuid.UUID) (*domain.DashboardSummary, error) {
	return s.repo.Summary(ctx, actor, branchID)
}

type SystemService struct {
	repo repository.SystemRepository
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
