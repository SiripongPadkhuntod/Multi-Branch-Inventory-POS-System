package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleOwner    Role = "OWNER"
	RoleManager  Role = "MANAGER"
	RoleEmployee Role = "EMPLOYEE"
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

type PaymentMethod string

const (
	PaymentCash         PaymentMethod = "CASH"
	PaymentPromptPay    PaymentMethod = "PROMPTPAY"
	PaymentBankTransfer PaymentMethod = "BANK_TRANSFER"
	PaymentCreditCard   PaymentMethod = "CREDIT_CARD"
)

type User struct {
	ID           uuid.UUID   `json:"id"`
	BranchID     *uuid.UUID  `json:"branch_id"`
	BranchIDs    []uuid.UUID `json:"branch_ids,omitempty"`
	Role         Role        `json:"role"`
	Name         string      `json:"name"`
	Email        string      `json:"email"`
	PasswordHash string      `json:"-"`
	Status       string      `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
}

type EmployeeSalesSummary struct {
	UserID     uuid.UUID `json:"user_id"`
	BranchID   uuid.UUID `json:"branch_id"`
	BranchCode string    `json:"branch_code"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Role       Role      `json:"role"`
	Status     string    `json:"status"`
	SalesCount int64     `json:"sales_count"`
	Revenue    int64     `json:"revenue"`
}

type Branch struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type AuditLog struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	UserName   string    `json:"user_name"`
	UserEmail  string    `json:"user_email"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   uuid.UUID `json:"entity_id"`
	OldData    string    `json:"old_data"`
	NewData    string    `json:"new_data"`
	IPAddress  string    `json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`
}

type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

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

type DashboardSummary struct {
	DailySales       int64                `json:"daily_sales"`
	MonthlySales     int64                `json:"monthly_sales"`
	Revenue          int64                `json:"revenue"`
	Profit           int64                `json:"profit"`
	LowStock         []LowStockItem       `json:"low_stock"`
	TopProducts      []TopProduct         `json:"top_products"`
	BranchComparison []BranchSalesSummary `json:"branch_comparison"`
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	SKU         string    `json:"sku"`
	Barcode     string    `json:"barcode"`
	QRCode      string    `json:"qr_code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CategoryID  uuid.UUID `json:"category_id"`
	ImageURL    string    `json:"image_url"`
	CostPrice   int64     `json:"cost_price"`
	SellPrice   int64     `json:"sell_price"`
	Status      string    `json:"status"`
}

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

type CartItemInput struct {
	ProductID      uuid.UUID `json:"product_id" binding:"required"`
	Quantity       int64     `json:"quantity" binding:"required,min=1"`
	FinalPrice     int64     `json:"final_price" binding:"required,min=0"`
	DiscountAmount int64     `json:"discount_amount" binding:"min=0"`
	DiscountReason string    `json:"discount_reason"`
}

type PaymentInput struct {
	Method PaymentMethod `json:"payment_method" binding:"required"`
	Amount int64         `json:"amount" binding:"required,min=1"`
}

type CreateSaleInput struct {
	BranchID uuid.UUID       `json:"branch_id" binding:"required"`
	Items    []CartItemInput `json:"items" binding:"required,min=1"`
	Payments []PaymentInput  `json:"payments" binding:"required,min=1"`
	Discount int64           `json:"discount" binding:"min=0"`
	Tax      int64           `json:"tax" binding:"min=0"`
}

type Sale struct {
	ID            uuid.UUID `json:"id"`
	ReceiptNumber string    `json:"receipt_number"`
	BranchID      uuid.UUID `json:"branch_id"`
	EmployeeID    uuid.UUID `json:"employee_id"`
	Subtotal      int64     `json:"subtotal"`
	Discount      int64     `json:"discount"`
	Tax           int64     `json:"tax"`
	Total         int64     `json:"total"`
	PaymentStatus string    `json:"payment_status"`
	RefundStatus  string    `json:"refund_status"`
	CreatedAt     time.Time `json:"created_at"`
}

type SaleItemDetail struct {
	ID             uuid.UUID `json:"id"`
	ProductID      uuid.UUID `json:"product_id"`
	SKU            string    `json:"sku"`
	Barcode        string    `json:"barcode"`
	ProductName    string    `json:"product_name"`
	Quantity       int64     `json:"quantity"`
	ReturnedQty    int64     `json:"returned_quantity"`
	OriginalPrice  int64     `json:"original_price"`
	FinalPrice     int64     `json:"final_price"`
	DiscountAmount int64     `json:"discount_amount"`
	DiscountReason string    `json:"discount_reason"`
	LineTotal      int64     `json:"line_total"`
}

type PaymentDetail struct {
	ID     uuid.UUID     `json:"id"`
	Method PaymentMethod `json:"payment_method"`
	Amount int64         `json:"amount"`
}

type SaleDetail struct {
	Sale
	BranchCode   string           `json:"branch_code"`
	BranchName   string           `json:"branch_name"`
	EmployeeName string           `json:"employee_name"`
	Items        []SaleItemDetail `json:"items"`
	Payments     []PaymentDetail  `json:"payments"`
}
