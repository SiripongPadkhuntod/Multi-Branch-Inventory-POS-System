package handler

import (
	"context"
	"errors"
	"strings"

	"pos-system/backend/internal/app/domain"
	"pos-system/backend/pkg/v1/dto"
	"pos-system/backend/pkg/v1/model"

	"github.com/google/uuid"
)

func (h *handler) InventoryListHandler(ctx context.Context, req dto.InventoryListRequest) (*dto.InventoryListResponse, error) {
	branchID, err := optionalUUID(req.BranchID)
	if err != nil {
		return nil, err
	}
	categoryID, err := optionalUUID(req.CategoryID)
	if err != nil {
		return nil, err
	}
	inventories, err := h.svc.ListInventories(ctx, branchID, categoryID, strings.TrimSpace(req.Query))
	if err != nil {
		return nil, err
	}
	data := make([]model.Inventory, 0, len(inventories))
	for _, item := range inventories {
		data = append(data, inventoryModel(item))
	}
	return &dto.InventoryListResponse{Code: "SUCCESS", Description: "success", Data: data}, nil
}

func (h *handler) InventoryMovementsHandler(ctx context.Context, req dto.InventoryMovementsRequest) (*dto.InventoryMovementsResponse, error) {
	branchID, err := optionalUUID(req.BranchID)
	if err != nil {
		return nil, err
	}
	if req.Limit <= 0 {
		req.Limit = 150
	}
	movements, err := h.svc.ListInventoryMovements(ctx, branchID, strings.TrimSpace(req.Query), req.Limit)
	if err != nil {
		return nil, err
	}
	data := make([]model.InventoryMovement, 0, len(movements))
	for _, movement := range movements {
		data = append(data, inventoryMovementModel(movement))
	}
	return &dto.InventoryMovementsResponse{Code: "SUCCESS", Description: "success", Data: data}, nil
}

func (h *handler) InventoryAllStockHandler(ctx context.Context, req dto.InventoryAllStockRequest) (*dto.InventoryAllStockResponse, error) {
	summaries, err := h.svc.AllStock(ctx, strings.TrimSpace(req.Query))
	if err != nil {
		return nil, err
	}
	data := make([]model.ProductStockSummary, 0, len(summaries))
	for _, summary := range summaries {
		data = append(data, productStockSummaryModel(summary))
	}
	return &dto.InventoryAllStockResponse{Code: "SUCCESS", Description: "success", Data: data}, nil
}

func (h *handler) InventoryAdjustHandler(ctx context.Context, req dto.InventoryAdjustRequest) (*dto.SuccessResponse, error) {
	if err := h.adjustInventory(ctx, req, false); err != nil {
		return nil, err
	}
	return &dto.SuccessResponse{Code: "SUCCESS", Description: "stock adjusted"}, nil
}

func (h *handler) InventoryReceiveHandler(ctx context.Context, req dto.InventoryAdjustRequest) (*dto.SuccessResponse, error) {
	if err := h.adjustInventory(ctx, req, true); err != nil {
		return nil, err
	}
	return &dto.SuccessResponse{Code: "SUCCESS", Description: "stock received"}, nil
}

func (h *handler) adjustInventory(ctx context.Context, req dto.InventoryAdjustRequest, receiveOnly bool) error {
	branchID, err := requiredUUID(req.BranchID, "branch_id is required")
	if err != nil {
		return err
	}
	productID, err := requiredUUID(req.ProductID, "product_id is required")
	if err != nil {
		return err
	}
	user, ok := domain.UserFromContext(ctx)
	if !ok {
		return errors.New("authenticated user is required")
	}
	if req.QuantityDelta == 0 {
		return errors.New("quantity_delta is required")
	}
	if receiveOnly && req.QuantityDelta < 1 {
		return errors.New("receive quantity must be greater than zero")
	}
	return h.svc.AdjustInventory(ctx, branchID, productID, user.ID, req.QuantityDelta, strings.TrimSpace(req.Reason))
}

func optionalUUID(value string) (*uuid.UUID, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	id, err := uuid.Parse(value)
	if err != nil {
		return nil, errors.New("invalid uuid")
	}
	return &id, nil
}

func requiredUUID(value, message string) (uuid.UUID, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return uuid.Nil, errors.New(message)
	}
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, errors.New("invalid uuid")
	}
	return id, nil
}

func inventoryModel(item domain.Inventory) model.Inventory {
	return model.Inventory{
		ID:               item.ID.String(),
		BranchID:         item.BranchID.String(),
		ProductID:        item.ProductID.String(),
		Quantity:         item.Quantity,
		ReservedQuantity: item.ReservedQuantity,
		ReorderThreshold: item.ReorderThreshold,
		UpdatedAt:        item.UpdatedAt,
	}
}

func inventoryMovementModel(movement domain.InventoryMovementDetail) model.InventoryMovement {
	return model.InventoryMovement{
		ID:            movement.ID.String(),
		BranchID:      movement.BranchID.String(),
		BranchCode:    movement.BranchCode,
		BranchName:    movement.BranchName,
		ProductID:     movement.ProductID.String(),
		SKU:           movement.SKU,
		Barcode:       movement.Barcode,
		ProductName:   movement.ProductName,
		MovementType:  string(movement.MovementType),
		Quantity:      movement.Quantity,
		ReferenceID:   movement.ReferenceID,
		CreatedBy:     movement.CreatedBy.String(),
		CreatedByName: movement.CreatedByName,
		CreatedAt:     movement.CreatedAt,
	}
}

func productStockSummaryModel(summary domain.ProductStockSummary) model.ProductStockSummary {
	branches := make([]model.BranchStockDetail, 0, len(summary.Branches))
	for _, branch := range summary.Branches {
		branches = append(branches, model.BranchStockDetail{
			BranchID:         branch.BranchID.String(),
			BranchCode:       branch.BranchCode,
			BranchName:       branch.BranchName,
			Quantity:         branch.Quantity,
			ReservedQuantity: branch.ReservedQuantity,
			ReorderThreshold: branch.ReorderThreshold,
			UpdatedAt:        branch.UpdatedAt,
		})
	}
	return model.ProductStockSummary{
		ProductID:     summary.ProductID.String(),
		SKU:           summary.SKU,
		Barcode:       summary.Barcode,
		ImageURL:      summary.ImageURL,
		ProductName:   summary.ProductName,
		TotalQuantity: summary.TotalQuantity,
		TotalReserved: summary.TotalReserved,
		Branches:      branches,
	}
}
