package port

import (
	"context"
	"time"

	"pos-system/backend/internal/app/domain"

	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, string, *domain.User, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, *domain.User, error)
	RevokeRefresh(ctx context.Context, refreshToken string) error
	FindUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type ProductService interface {
	List(ctx context.Context, query string, limit, offset int) ([]domain.Product, error)
	ByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	ByBarcode(ctx context.Context, barcode string) (*domain.Product, error)
	Create(ctx context.Context, product domain.Product) (*domain.Product, error)
	Update(ctx context.Context, id uuid.UUID, product domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type InventoryService interface {
	List(ctx context.Context, branchID *uuid.UUID, categoryID *uuid.UUID, query string) ([]domain.Inventory, error)
	Movements(ctx context.Context, branchID *uuid.UUID, query string, limit int) ([]domain.InventoryMovementDetail, error)
	AllStock(ctx context.Context, query string) ([]domain.ProductStockSummary, error)
	Adjust(ctx context.Context, branchID, productID, actorID uuid.UUID, delta int64, reason string) error
	SetReorderThreshold(ctx context.Context, branchID, productID uuid.UUID, threshold int64) error
	CreateTransfer(ctx context.Context, fromBranchID, toBranchID, productID, actorID uuid.UUID, quantity int64) (*domain.Transfer, error)
	Transfers(ctx context.Context, status string, limit int) ([]domain.Transfer, error)
	ApproveTransfer(ctx context.Context, transferID, actorID uuid.UUID) (*domain.Transfer, error)
	RejectTransfer(ctx context.Context, transferID, actorID uuid.UUID) (*domain.Transfer, error)
	CompleteTransfer(ctx context.Context, transferID, actorID uuid.UUID) (*domain.Transfer, error)
}

type SaleService interface {
	Create(ctx context.Context, actor domain.User, input domain.CreateSaleInput) (*domain.Sale, error)
	List(ctx context.Context, actor domain.User, branchID *uuid.UUID, dateFrom, dateTo *time.Time) ([]domain.Sale, error)
	BranchList(ctx context.Context, actor domain.User, branchID *uuid.UUID) ([]domain.Sale, error)
	Detail(ctx context.Context, actor domain.User, saleID uuid.UUID) (*domain.SaleDetail, error)
	Refund(ctx context.Context, actor domain.User, saleID uuid.UUID, items []domain.CartItemInput) error
}

type Service interface {
	Health(ctx context.Context) error
	Authenticate(ctx context.Context, authorization string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, string, *domain.User, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, *domain.User, error)
	Logout(ctx context.Context, refreshToken string) error
	ListProducts(ctx context.Context, query string, limit, offset int) ([]domain.Product, error)
	ProductByBarcode(ctx context.Context, barcode string) (*domain.Product, error)
	ListInventories(ctx context.Context, branchID *uuid.UUID, categoryID *uuid.UUID, query string) ([]domain.Inventory, error)
	ListInventoryMovements(ctx context.Context, branchID *uuid.UUID, query string, limit int) ([]domain.InventoryMovementDetail, error)
	AllStock(ctx context.Context, query string) ([]domain.ProductStockSummary, error)
	AdjustInventory(ctx context.Context, branchID, productID, actorID uuid.UUID, delta int64, reason string) error
}
