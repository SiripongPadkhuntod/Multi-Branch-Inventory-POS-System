package domain

import (
	"time"

	"github.com/google/uuid"
)

type MovementType string

const (
	MovementReceive     MovementType = "RECEIVE"
	MovementSale        MovementType = "SALE"
	MovementReturn      MovementType = "RETURN"
	MovementAdjustment  MovementType = "ADJUSTMENT"
	MovementTransferIn  MovementType = "TRANSFER_IN"
	MovementTransferOut MovementType = "TRANSFER_OUT"
)

type Inventory struct {
	ID               uuid.UUID `json:"id"`
	BranchID         uuid.UUID `json:"branch_id"`
	ProductID        uuid.UUID `json:"product_id"`
	Quantity         int64     `json:"quantity"`
	ReservedQuantity int64     `json:"reserved_quantity"`
	ReorderThreshold int64     `json:"reorder_threshold"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type InventoryMovementDetail struct {
	ID            uuid.UUID    `json:"id"`
	BranchID      uuid.UUID    `json:"branch_id"`
	BranchCode    string       `json:"branch_code"`
	BranchName    string       `json:"branch_name"`
	ProductID     uuid.UUID    `json:"product_id"`
	SKU           string       `json:"sku"`
	Barcode       string       `json:"barcode"`
	ProductName   string       `json:"product_name"`
	MovementType  MovementType `json:"movement_type"`
	Quantity      int64        `json:"quantity"`
	ReferenceID   string       `json:"reference_id"`
	CreatedBy     uuid.UUID    `json:"created_by"`
	CreatedByName string       `json:"created_by_name"`
	CreatedAt     time.Time    `json:"created_at"`
}

type BranchStockDetail struct {
	BranchID         uuid.UUID `json:"branch_id"`
	BranchCode       string    `json:"branch_code"`
	BranchName       string    `json:"branch_name"`
	Quantity         int64     `json:"quantity"`
	ReservedQuantity int64     `json:"reserved_quantity"`
	ReorderThreshold int64     `json:"reorder_threshold"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ProductStockSummary struct {
	ProductID     uuid.UUID           `json:"product_id"`
	SKU           string              `json:"sku"`
	Barcode       string              `json:"barcode"`
	ImageURL      string              `json:"image_url"`
	ProductName   string              `json:"product_name"`
	TotalQuantity int64               `json:"total_quantity"`
	TotalReserved int64               `json:"total_reserved"`
	Branches      []BranchStockDetail `json:"branches"`
}

type Transfer struct {
	ID              uuid.UUID  `json:"id"`
	FromBranchID    uuid.UUID  `json:"from_branch_id"`
	FromBranchCode  string     `json:"from_branch_code"`
	FromBranchName  string     `json:"from_branch_name"`
	ToBranchID      uuid.UUID  `json:"to_branch_id"`
	ToBranchCode    string     `json:"to_branch_code"`
	ToBranchName    string     `json:"to_branch_name"`
	Status          string     `json:"status"`
	RequestedBy     uuid.UUID  `json:"requested_by"`
	RequestedByName string     `json:"requested_by_name"`
	ApprovedBy      *uuid.UUID `json:"approved_by"`
	ApprovedByName  string     `json:"approved_by_name"`
	ProductID       uuid.UUID  `json:"product_id"`
	ProductSKU      string     `json:"product_sku"`
	ProductBarcode  string     `json:"product_barcode"`
	ProductName     string     `json:"product_name"`
	Quantity        int64      `json:"quantity"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
