package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string

const (
	PaymentCash         PaymentMethod = "CASH"
	PaymentPromptPay    PaymentMethod = "PROMPTPAY"
	PaymentBankTransfer PaymentMethod = "BANK_TRANSFER"
	PaymentCreditCard   PaymentMethod = "CREDIT_CARD"
)

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
