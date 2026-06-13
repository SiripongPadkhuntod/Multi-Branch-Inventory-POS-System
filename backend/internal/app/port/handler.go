package port

import (
	"context"

	"pos-system/backend/pkg/v1/dto"
)

type Handler interface {
	HealthHandler(ctx context.Context, req dto.EmptyStruct) (*dto.HealthResponse, error)
	LoginHandler(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshHandler(ctx context.Context, req dto.RefreshRequest) (*dto.LoginResponse, error)
	LogoutHandler(ctx context.Context, req dto.LogoutRequest) (*dto.SuccessResponse, error)
	ProductListHandler(ctx context.Context, req dto.ProductListRequest) (*dto.ProductListResponse, error)
	ProductBarcodeHandler(ctx context.Context, req dto.ProductBarcodeRequest) (*dto.ProductResponse, error)
	InventoryListHandler(ctx context.Context, req dto.InventoryListRequest) (*dto.InventoryListResponse, error)
	InventoryMovementsHandler(ctx context.Context, req dto.InventoryMovementsRequest) (*dto.InventoryMovementsResponse, error)
	InventoryAllStockHandler(ctx context.Context, req dto.InventoryAllStockRequest) (*dto.InventoryAllStockResponse, error)
	InventoryAdjustHandler(ctx context.Context, req dto.InventoryAdjustRequest) (*dto.SuccessResponse, error)
	InventoryReceiveHandler(ctx context.Context, req dto.InventoryAdjustRequest) (*dto.SuccessResponse, error)
}
