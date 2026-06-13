package domain

import "github.com/google/uuid"

type TopProduct struct {
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Quantity  int64     `json:"quantity"`
	Revenue   int64     `json:"revenue"`
}

type LowStockItem struct {
	BranchID         uuid.UUID `json:"branch_id"`
	BranchCode       string    `json:"branch_code"`
	ProductID        uuid.UUID `json:"product_id"`
	ProductName      string    `json:"product_name"`
	Quantity         int64     `json:"quantity"`
	ReorderThreshold int64     `json:"reorder_threshold"`
}

type BranchSalesSummary struct {
	BranchID   uuid.UUID `json:"branch_id"`
	BranchCode string    `json:"branch_code"`
	BranchName string    `json:"branch_name"`
	Revenue    int64     `json:"revenue"`
	SalesCount int64     `json:"sales_count"`
}

type DashboardSummary struct {
	DailySales       int64                `json:"daily_sales"`
	MonthlySales     int64                `json:"monthly_sales"`
	Revenue          int64                `json:"revenue"`
	Profit           int64                `json:"profit"`
	LowStock         []LowStockItem       `json:"low_stock"`
	TopProducts      []TopProduct         `json:"top_products"`
	BranchComparison []BranchSalesSummary `json:"branch_comparison"`
}
