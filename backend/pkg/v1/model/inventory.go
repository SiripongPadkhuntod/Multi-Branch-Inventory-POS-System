package model

import "time"

type Inventory struct {
	ID               string    `json:"id"`
	BranchID         string    `json:"branch_id"`
	ProductID        string    `json:"product_id"`
	Quantity         int64     `json:"quantity"`
	ReservedQuantity int64     `json:"reserved_quantity"`
	ReorderThreshold int64     `json:"reorder_threshold"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type InventoryMovement struct {
	ID            string    `json:"id"`
	BranchID      string    `json:"branch_id"`
	BranchCode    string    `json:"branch_code"`
	BranchName    string    `json:"branch_name"`
	ProductID     string    `json:"product_id"`
	SKU           string    `json:"sku"`
	Barcode       string    `json:"barcode"`
	ProductName   string    `json:"product_name"`
	MovementType  string    `json:"movement_type"`
	Quantity      int64     `json:"quantity"`
	ReferenceID   string    `json:"reference_id"`
	CreatedBy     string    `json:"created_by"`
	CreatedByName string    `json:"created_by_name"`
	CreatedAt     time.Time `json:"created_at"`
}

type BranchStockDetail struct {
	BranchID         string    `json:"branch_id"`
	BranchCode       string    `json:"branch_code"`
	BranchName       string    `json:"branch_name"`
	Quantity         int64     `json:"quantity"`
	ReservedQuantity int64     `json:"reserved_quantity"`
	ReorderThreshold int64     `json:"reorder_threshold"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ProductStockSummary struct {
	ProductID     string              `json:"product_id"`
	SKU           string              `json:"sku"`
	Barcode       string              `json:"barcode"`
	ImageURL      string              `json:"image_url"`
	ProductName   string              `json:"product_name"`
	TotalQuantity int64               `json:"total_quantity"`
	TotalReserved int64               `json:"total_reserved"`
	Branches      []BranchStockDetail `json:"branches"`
}
