package dto

import "pos-system/backend/pkg/v1/model"

type InventoryListRequest struct {
	BranchID   string `form:"branch_id"`
	CategoryID string `form:"category_id"`
	Query      string `form:"q"`
}

type InventoryMovementsRequest struct {
	BranchID string `form:"branch_id"`
	Query    string `form:"q"`
	Limit    int    `form:"limit"`
}

type InventoryAllStockRequest struct {
	Query string `form:"q"`
}

type InventoryAdjustRequest struct {
	BranchID      string `json:"branch_id"`
	ProductID     string `json:"product_id"`
	QuantityDelta int64  `json:"quantity_delta"`
	Reason        string `json:"reason"`
}

type InventoryListResponse struct {
	Code        string            `json:"code"`
	Description string            `json:"description"`
	Data        []model.Inventory `json:"data"`
}

type InventoryMovementsResponse struct {
	Code        string                    `json:"code"`
	Description string                    `json:"description"`
	Data        []model.InventoryMovement `json:"data"`
}

type InventoryAllStockResponse struct {
	Code        string                      `json:"code"`
	Description string                      `json:"description"`
	Data        []model.ProductStockSummary `json:"data"`
}
