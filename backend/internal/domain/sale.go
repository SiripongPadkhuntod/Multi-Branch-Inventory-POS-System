package domain

import appdomain "pos-system/backend/internal/app/domain"

type PaymentMethod = appdomain.PaymentMethod

const (
	PaymentCash         = appdomain.PaymentCash
	PaymentPromptPay    = appdomain.PaymentPromptPay
	PaymentBankTransfer = appdomain.PaymentBankTransfer
	PaymentCreditCard   = appdomain.PaymentCreditCard
)

type CartItemInput = appdomain.CartItemInput
type PaymentInput = appdomain.PaymentInput
type CreateSaleInput = appdomain.CreateSaleInput
type Sale = appdomain.Sale
type SaleItemDetail = appdomain.SaleItemDetail
type PaymentDetail = appdomain.PaymentDetail
type SaleDetail = appdomain.SaleDetail
