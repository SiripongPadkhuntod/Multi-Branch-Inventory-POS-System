package repository

import (
	"context"
	"time"

	"pos-system/backend/internal/domain"

	"github.com/google/uuid"
)

type Repositories interface {
	Auth() AuthRepository
	Products() ProductRepository
	Inventories() InventoryRepository
	Sales() SaleRepository
	Users() UserRepository
	Dashboard() DashboardRepository
	System() SystemRepository
	Audit() AuditRepository
}

type AuthRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type ProductRepository interface {
	List(ctx context.Context, query string, limit, offset int) ([]domain.Product, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	FindByBarcode(ctx context.Context, barcode string) (*domain.Product, error)
	Create(ctx context.Context, product domain.Product) (*domain.Product, error)
	Update(ctx context.Context, id uuid.UUID, product domain.Product) (*domain.Product, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type InventoryRepository interface {
	List(ctx context.Context, branchID *uuid.UUID, categoryID *uuid.UUID, query string) ([]domain.Inventory, error)
	ListMovements(ctx context.Context, branchID *uuid.UUID, query string, limit int) ([]domain.InventoryMovementDetail, error)
	AllStock(ctx context.Context, query string) ([]domain.ProductStockSummary, error)
	Adjust(ctx context.Context, branchID, productID, actorID uuid.UUID, quantityDelta int64, reason string) error
	SetReorderThreshold(ctx context.Context, branchID, productID uuid.UUID, threshold int64) error
	Transfer(ctx context.Context, fromBranchID, toBranchID, productID, actorID uuid.UUID, quantity int64) error
}

type SaleRepository interface {
	CreateSale(ctx context.Context, actor domain.User, input domain.CreateSaleInput) (*domain.Sale, error)
	List(ctx context.Context, actor domain.User, branchID *uuid.UUID, dateFrom, dateTo *time.Time) ([]domain.Sale, error)
	ListBranch(ctx context.Context, actor domain.User, branchID *uuid.UUID) ([]domain.Sale, error)
	FindDetail(ctx context.Context, actor domain.User, saleID uuid.UUID) (*domain.SaleDetail, error)
	Refund(ctx context.Context, actor domain.User, saleID uuid.UUID, items []domain.CartItemInput) error
}

type UserRepository interface {
	List(ctx context.Context, actor domain.User) ([]domain.User, error)
	Create(ctx context.Context, actor domain.User, user domain.User) (*domain.User, error)
	Update(ctx context.Context, actor domain.User, id uuid.UUID, user domain.User) (*domain.User, error)
	SalesSummary(ctx context.Context, actor domain.User) ([]domain.EmployeeSalesSummary, error)
}

type DashboardRepository interface {
	AccessibleBranches(ctx context.Context, actor domain.User) ([]domain.Branch, error)
	Summary(ctx context.Context, actor domain.User, branchID *uuid.UUID) (*domain.DashboardSummary, error)
}

type SystemRepository interface {
	ListCategories(ctx context.Context) ([]domain.Category, error)
	CreateCategory(ctx context.Context, category domain.Category) (*domain.Category, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, category domain.Category) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	CreateBranch(ctx context.Context, branch domain.Branch) (*domain.Branch, error)
	UpdateBranch(ctx context.Context, id uuid.UUID, branch domain.Branch) (*domain.Branch, error)
}

type AuditRepository interface {
	List(ctx context.Context, action, entityType, query string, limit int) ([]domain.AuditLog, error)
}
