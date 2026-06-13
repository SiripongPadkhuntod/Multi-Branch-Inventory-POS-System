package domain

import appdomain "pos-system/backend/internal/app/domain"

type MovementType = appdomain.MovementType

const (
	MovementReceive     = appdomain.MovementReceive
	MovementSale        = appdomain.MovementSale
	MovementReturn      = appdomain.MovementReturn
	MovementAdjustment  = appdomain.MovementAdjustment
	MovementTransferIn  = appdomain.MovementTransferIn
	MovementTransferOut = appdomain.MovementTransferOut
)

type Inventory = appdomain.Inventory
type InventoryMovementDetail = appdomain.InventoryMovementDetail
type BranchStockDetail = appdomain.BranchStockDetail
type ProductStockSummary = appdomain.ProductStockSummary
type Transfer = appdomain.Transfer
