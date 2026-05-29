package postgres

import (
	"context"
	"errors"

	"pos-system/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepo struct{ db *pgxpool.Pool }

func (r *InventoryRepo) List(ctx context.Context, branchID *uuid.UUID, categoryID *uuid.UUID, query string) ([]domain.Inventory, error) {
	rows, err := r.db.Query(ctx, `
		SELECT i.id, i.branch_id, i.product_id, i.quantity, i.reserved_quantity, i.updated_at
		FROM inventories i
		JOIN products p ON p.id=i.product_id
		WHERE i.deleted_at IS NULL
		  AND ($1::uuid IS NULL OR i.branch_id=$1)
		  AND ($2::uuid IS NULL OR p.category_id=$2)
		  AND ($3='' OR p.name ILIKE '%'||$3||'%' OR p.sku ILIKE '%'||$3||'%' OR p.barcode ILIKE '%'||$3||'%')
		ORDER BY i.updated_at DESC`, branchID, categoryID, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventories := make([]domain.Inventory, 0)
	for rows.Next() {
		var i domain.Inventory
		if err := rows.Scan(&i.ID, &i.BranchID, &i.ProductID, &i.Quantity, &i.ReservedQuantity, &i.UpdatedAt); err != nil {
			return nil, err
		}
		inventories = append(inventories, i)
	}
	return inventories, rows.Err()
}

func (r *InventoryRepo) ListMovements(ctx context.Context, branchID *uuid.UUID, query string, limit int) ([]domain.InventoryMovementDetail, error) {
	if limit <= 0 || limit > 300 {
		limit = 150
	}
	rows, err := r.db.Query(ctx, `
		SELECT im.id, im.branch_id, b.code, b.name, im.product_id, p.sku, p.barcode, p.name,
			im.movement_type, im.quantity, im.reference_id, im.created_by, u.name, im.created_at
		FROM inventory_movements im
		JOIN branches b ON b.id=im.branch_id
		JOIN products p ON p.id=im.product_id
		JOIN users u ON u.id=im.created_by
		WHERE im.deleted_at IS NULL
			AND ($1::uuid IS NULL OR im.branch_id=$1)
			AND ($2='' OR p.name ILIKE '%'||$2||'%' OR p.sku ILIKE '%'||$2||'%' OR p.barcode ILIKE '%'||$2||'%' OR b.code ILIKE '%'||$2||'%')
		ORDER BY im.created_at DESC
		LIMIT $3`, branchID, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movements := make([]domain.InventoryMovementDetail, 0)
	for rows.Next() {
		var movement domain.InventoryMovementDetail
		if err := rows.Scan(&movement.ID, &movement.BranchID, &movement.BranchCode, &movement.BranchName, &movement.ProductID,
			&movement.SKU, &movement.Barcode, &movement.ProductName, &movement.MovementType, &movement.Quantity,
			&movement.ReferenceID, &movement.CreatedBy, &movement.CreatedByName, &movement.CreatedAt); err != nil {
			return nil, err
		}
		movements = append(movements, movement)
	}
	return movements, rows.Err()
}

func (r *InventoryRepo) AllStock(ctx context.Context, query string) ([]domain.ProductStockSummary, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.sku, p.barcode, p.image_url, p.name, i.branch_id, b.code, b.name,
			COALESCE(i.quantity, 0), COALESCE(i.reserved_quantity, 0), COALESCE(i.updated_at, p.updated_at)
		FROM products p
		JOIN branches b ON b.deleted_at IS NULL
		LEFT JOIN inventories i ON i.product_id=p.id AND i.branch_id=b.id AND i.deleted_at IS NULL
		WHERE p.deleted_at IS NULL
			AND ($1='' OR p.name ILIKE '%'||$1||'%' OR p.sku ILIKE '%'||$1||'%' OR p.barcode ILIKE '%'||$1||'%' OR b.code ILIKE '%'||$1||'%')
		ORDER BY p.name, b.code`, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summaries := make([]domain.ProductStockSummary, 0)
	indexByProduct := map[uuid.UUID]int{}
	for rows.Next() {
		var productID uuid.UUID
		var item domain.ProductStockSummary
		var branch domain.BranchStockDetail
		if err := rows.Scan(&productID, &item.SKU, &item.Barcode, &item.ImageURL, &item.ProductName, &branch.BranchID, &branch.BranchCode,
			&branch.BranchName, &branch.Quantity, &branch.ReservedQuantity, &branch.UpdatedAt); err != nil {
			return nil, err
		}
		idx, exists := indexByProduct[productID]
		if !exists {
			item.ProductID = productID
			item.Branches = make([]domain.BranchStockDetail, 0)
			summaries = append(summaries, item)
			idx = len(summaries) - 1
			indexByProduct[productID] = idx
		}
		summaries[idx].TotalQuantity += branch.Quantity
		summaries[idx].TotalReserved += branch.ReservedQuantity
		summaries[idx].Branches = append(summaries[idx].Branches, branch)
	}
	return summaries, rows.Err()
}

func (r *InventoryRepo) Adjust(ctx context.Context, branchID, productID, actorID uuid.UUID, quantityDelta int64, reason string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var current int64
	err = tx.QueryRow(ctx, `
		SELECT quantity FROM inventories
		WHERE branch_id=$1 AND product_id=$2 AND deleted_at IS NULL
		FOR UPDATE`, branchID, productID).Scan(&current)
	if errors.Is(err, pgx.ErrNoRows) && quantityDelta > 0 {
		_, err = tx.Exec(ctx, `
			INSERT INTO inventories (branch_id, product_id, quantity)
			VALUES ($1,$2,0)
			ON CONFLICT (branch_id, product_id) DO UPDATE SET deleted_at=NULL, updated_at=now()`,
			branchID, productID)
		if err != nil {
			return err
		}
		err = tx.QueryRow(ctx, `
			SELECT quantity FROM inventories
			WHERE branch_id=$1 AND product_id=$2 AND deleted_at IS NULL
			FOR UPDATE`, branchID, productID).Scan(&current)
	}
	if err != nil {
		return err
	}
	if current+quantityDelta < 0 {
		return errors.New("stock cannot become negative")
	}

	_, err = tx.Exec(ctx, `
		UPDATE inventories SET quantity=quantity+$3, updated_at=now()
		WHERE branch_id=$1 AND product_id=$2`, branchID, productID, quantityDelta)
	if err != nil {
		return err
	}

	movementType := domain.MovementAdjustment
	if quantityDelta > 0 {
		movementType = domain.MovementReceive
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO inventory_movements (branch_id, product_id, movement_type, quantity, reference_id, created_by)
		VALUES ($1,$2,$3,$4,$5,$6)`, branchID, productID, movementType, quantityDelta, reason, actorID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *InventoryRepo) Transfer(ctx context.Context, fromBranchID, toBranchID, productID, actorID uuid.UUID, quantity int64) error {
	if fromBranchID == toBranchID {
		return errors.New("source and destination branches must be different")
	}
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var sourceQty int64
	err = tx.QueryRow(ctx, `
		SELECT quantity FROM inventories
		WHERE branch_id=$1 AND product_id=$2 AND deleted_at IS NULL
		FOR UPDATE`, fromBranchID, productID).Scan(&sourceQty)
	if err != nil {
		return err
	}
	if sourceQty < quantity {
		return errors.New("source branch has insufficient stock")
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO inventories (branch_id, product_id, quantity)
		VALUES ($1,$2,0)
		ON CONFLICT (branch_id, product_id) DO UPDATE SET deleted_at=NULL, updated_at=now()`,
		toBranchID, productID)
	if err != nil {
		return err
	}
	var destinationQty int64
	if err := tx.QueryRow(ctx, `
		SELECT quantity FROM inventories
		WHERE branch_id=$1 AND product_id=$2 AND deleted_at IS NULL
		FOR UPDATE`, toBranchID, productID).Scan(&destinationQty); err != nil {
		return err
	}

	var transferID uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO transfers (from_branch_id, to_branch_id, status, requested_by, approved_by)
		VALUES ($1,$2,'COMPLETED',$3,$3)
		RETURNING id`, fromBranchID, toBranchID, actorID).Scan(&transferID)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO transfer_items (transfer_id, product_id, quantity)
		VALUES ($1,$2,$3)`, transferID, productID, quantity); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
		UPDATE inventories SET quantity=quantity-$3, updated_at=now()
		WHERE branch_id=$1 AND product_id=$2`, fromBranchID, productID, quantity); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		UPDATE inventories SET quantity=quantity+$3, updated_at=now()
		WHERE branch_id=$1 AND product_id=$2`, toBranchID, productID, quantity); err != nil {
		return err
	}

	referenceID := transferID.String()
	if _, err := tx.Exec(ctx, `
		INSERT INTO inventory_movements (branch_id, product_id, movement_type, quantity, reference_id, created_by)
		VALUES ($1,$2,$3,$4,$5,$6)`, fromBranchID, productID, domain.MovementTransferOut, -quantity, referenceID, actorID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO inventory_movements (branch_id, product_id, movement_type, quantity, reference_id, created_by)
		VALUES ($1,$2,$3,$4,$5,$6)`, toBranchID, productID, domain.MovementTransferIn, quantity, referenceID, actorID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
