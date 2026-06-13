package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"pos-system/backend/internal/app/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SaleRepo struct{ db *pgxpool.Pool }

const saleRefundStatusSQL = `
	CASE
		WHEN COALESCE((
			SELECT SUM(im.quantity)
			FROM inventory_movements im
			WHERE im.reference_id=s.id::text AND im.movement_type='RETURN' AND im.deleted_at IS NULL
		), 0) = 0 THEN 'NONE'
		WHEN COALESCE((
			SELECT SUM(im.quantity)
			FROM inventory_movements im
			WHERE im.reference_id=s.id::text AND im.movement_type='RETURN' AND im.deleted_at IS NULL
		), 0) >= COALESCE((
			SELECT SUM(si.quantity)
			FROM sale_items si
			WHERE si.sale_id=s.id AND si.deleted_at IS NULL
		), 0) THEN 'REFUNDED'
		ELSE 'PARTIAL_REFUND'
	END`

func (r *SaleRepo) CreateSale(ctx context.Context, actor domain.User, input domain.CreateSaleInput) (*domain.Sale, error) {
	if isBranchStaff(actor) && (actor.BranchID == nil || *actor.BranchID != input.BranchID) {
		return nil, errors.New("staff cannot sell outside assigned branch")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	subtotal := int64(0)
	for _, item := range input.Items {
		var stock int64
		var originalPrice int64
		err = tx.QueryRow(ctx, `
			SELECT i.quantity, p.sell_price
			FROM inventories i
			JOIN products p ON p.id=i.product_id
			WHERE i.branch_id=$1 AND i.product_id=$2 AND i.deleted_at IS NULL
			FOR UPDATE`, input.BranchID, item.ProductID).Scan(&stock, &originalPrice)
		if err != nil {
			return nil, err
		}
		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s", item.ProductID)
		}
		subtotal += item.FinalPrice * item.Quantity
	}

	total := subtotal - input.Discount + input.Tax
	if total < 0 {
		return nil, errors.New("sale total cannot be negative")
	}
	paid := int64(0)
	for _, p := range input.Payments {
		paid += p.Amount
	}
	if paid < total {
		return nil, errors.New("payment amount is less than sale total")
	}

	receipt := fmt.Sprintf("RCPT-%s-%d", time.Now().Format("20060102"), time.Now().UnixNano())
	var sale domain.Sale
	err = tx.QueryRow(ctx, `
		INSERT INTO sales (receipt_number, branch_id, employee_id, subtotal, discount, tax, total, payment_status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,'PAID')
		RETURNING id, receipt_number, branch_id, employee_id, subtotal, discount, tax, total, payment_status, created_at`,
		receipt, input.BranchID, actor.ID, subtotal, input.Discount, input.Tax, total).
		Scan(&sale.ID, &sale.ReceiptNumber, &sale.BranchID, &sale.EmployeeID, &sale.Subtotal, &sale.Discount, &sale.Tax, &sale.Total, &sale.PaymentStatus, &sale.CreatedAt)
	if err != nil {
		return nil, err
	}
	sale.RefundStatus = "NONE"

	for _, item := range input.Items {
		var originalPrice int64
		if err := tx.QueryRow(ctx, `SELECT sell_price FROM products WHERE id=$1`, item.ProductID).Scan(&originalPrice); err != nil {
			return nil, err
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO sale_items (sale_id, product_id, quantity, original_price, final_price, discount_amount, discount_reason, employee_id)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
			sale.ID, item.ProductID, item.Quantity, originalPrice, item.FinalPrice, item.DiscountAmount, item.DiscountReason, actor.ID)
		if err != nil {
			return nil, err
		}
		_, err = tx.Exec(ctx, `UPDATE inventories SET quantity=quantity-$3, updated_at=now() WHERE branch_id=$1 AND product_id=$2`,
			input.BranchID, item.ProductID, item.Quantity)
		if err != nil {
			return nil, err
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO inventory_movements (branch_id, product_id, movement_type, quantity, reference_id, created_by)
			VALUES ($1,$2,$3,$4,$5,$6)`, input.BranchID, item.ProductID, domain.MovementSale, -item.Quantity, sale.ID.String(), actor.ID)
		if err != nil {
			return nil, err
		}
	}

	for _, payment := range input.Payments {
		_, err = tx.Exec(ctx, `INSERT INTO payments (sale_id, payment_method, amount) VALUES ($1,$2,$3)`,
			sale.ID, payment.Method, payment.Amount)
		if err != nil {
			return nil, err
		}
	}

	return &sale, tx.Commit(ctx)
}

func (r *SaleRepo) List(ctx context.Context, actor domain.User, branchID *uuid.UUID, dateFrom, dateTo *time.Time) ([]domain.Sale, error) {
	filterBranch := branchID
	var employeeID *uuid.UUID
	if isBranchStaff(actor) {
		filterBranch = actor.BranchID
		employeeID = &actor.ID
	}
	rows, err := r.db.Query(ctx, `
		SELECT s.id, s.receipt_number, s.branch_id, s.employee_id, s.subtotal, s.discount, s.tax, s.total, s.payment_status, `+saleRefundStatusSQL+`, s.created_at
		FROM sales s
		WHERE s.deleted_at IS NULL
			AND ($1::uuid IS NULL OR s.branch_id=$1)
			AND ($2::uuid IS NULL OR s.employee_id=$2)
			AND ($3::timestamptz IS NULL OR s.created_at >= $3)
			AND ($4::timestamptz IS NULL OR s.created_at < $4)
		ORDER BY s.created_at DESC LIMIT 100`, filterBranch, employeeID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sales := make([]domain.Sale, 0)
	for rows.Next() {
		var s domain.Sale
		if err := rows.Scan(&s.ID, &s.ReceiptNumber, &s.BranchID, &s.EmployeeID, &s.Subtotal, &s.Discount, &s.Tax, &s.Total, &s.PaymentStatus, &s.RefundStatus, &s.CreatedAt); err != nil {
			return nil, err
		}
		sales = append(sales, s)
	}
	return sales, rows.Err()
}

func (r *SaleRepo) ListBranch(ctx context.Context, actor domain.User, branchID *uuid.UUID) ([]domain.Sale, error) {
	if actor.Role == domain.RoleEmployee {
		return nil, errors.New("employee cannot view branch sales")
	}
	rows, err := r.db.Query(ctx, `
		SELECT s.id, s.receipt_number, s.branch_id, s.employee_id, s.subtotal, s.discount, s.tax, s.total, s.payment_status, `+saleRefundStatusSQL+`, s.created_at
		FROM sales s
		WHERE s.deleted_at IS NULL
			AND ($1::uuid IS NULL OR s.branch_id=$1)
			AND (
				$2 = 'OWNER'
				OR EXISTS (
					SELECT 1 FROM user_branches ub
					WHERE ub.user_id=$3 AND ub.branch_id=s.branch_id
				)
			)
		ORDER BY s.created_at DESC LIMIT 150`, branchID, actor.Role, actor.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sales := make([]domain.Sale, 0)
	for rows.Next() {
		var s domain.Sale
		if err := rows.Scan(&s.ID, &s.ReceiptNumber, &s.BranchID, &s.EmployeeID, &s.Subtotal, &s.Discount, &s.Tax, &s.Total, &s.PaymentStatus, &s.RefundStatus, &s.CreatedAt); err != nil {
			return nil, err
		}
		sales = append(sales, s)
	}
	return sales, rows.Err()
}

func (r *SaleRepo) FindDetail(ctx context.Context, actor domain.User, saleID uuid.UUID) (*domain.SaleDetail, error) {
	var detail domain.SaleDetail
	err := r.db.QueryRow(ctx, `
		SELECT s.id, s.receipt_number, s.branch_id, s.employee_id, s.subtotal, s.discount, s.tax, s.total,
			s.payment_status, `+saleRefundStatusSQL+`, s.created_at, b.code, b.name, u.name
		FROM sales s
		JOIN branches b ON b.id=s.branch_id
		JOIN users u ON u.id=s.employee_id
		WHERE s.id=$1 AND s.deleted_at IS NULL
			AND (
				$2 = 'OWNER'
				OR ($2 = 'EMPLOYEE' AND s.branch_id=$3 AND s.employee_id=$4)
				OR ($2 = 'MANAGER' AND EXISTS (
					SELECT 1 FROM user_branches ub
					WHERE ub.user_id=$4 AND ub.branch_id=s.branch_id
				))
			)`,
		saleID, actor.Role, actor.BranchID, actor.ID).
		Scan(&detail.ID, &detail.ReceiptNumber, &detail.BranchID, &detail.EmployeeID, &detail.Subtotal, &detail.Discount,
			&detail.Tax, &detail.Total, &detail.PaymentStatus, &detail.RefundStatus, &detail.CreatedAt, &detail.BranchCode, &detail.BranchName, &detail.EmployeeName)
	if err != nil {
		return nil, err
	}

	items, err := r.saleItems(ctx, saleID)
	if err != nil {
		return nil, err
	}
	payments, err := r.salePayments(ctx, saleID)
	if err != nil {
		return nil, err
	}
	detail.Items = items
	detail.Payments = payments
	return &detail, nil
}

func isBranchStaff(actor domain.User) bool {
	return actor.Role == domain.RoleEmployee || actor.Role == domain.RoleManager
}

func (r *SaleRepo) saleItems(ctx context.Context, saleID uuid.UUID) ([]domain.SaleItemDetail, error) {
	rows, err := r.db.Query(ctx, `
		SELECT si.id, si.product_id, p.sku, p.barcode, p.name, si.quantity,
			COALESCE((
				SELECT SUM(im.quantity)
				FROM inventory_movements im
				WHERE im.reference_id=$3 AND im.product_id=si.product_id AND im.movement_type=$2 AND im.deleted_at IS NULL
			), 0),
			si.original_price, si.final_price, si.discount_amount, si.discount_reason, si.quantity * si.final_price
		FROM sale_items si
		JOIN products p ON p.id=si.product_id
		WHERE si.sale_id=$1 AND si.deleted_at IS NULL
		ORDER BY si.created_at, si.id`, saleID, domain.MovementReturn, saleID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.SaleItemDetail
	for rows.Next() {
		var item domain.SaleItemDetail
		if err := rows.Scan(&item.ID, &item.ProductID, &item.SKU, &item.Barcode, &item.ProductName, &item.Quantity, &item.ReturnedQty,
			&item.OriginalPrice, &item.FinalPrice, &item.DiscountAmount, &item.DiscountReason, &item.LineTotal); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *SaleRepo) salePayments(ctx context.Context, saleID uuid.UUID) ([]domain.PaymentDetail, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, payment_method, amount
		FROM payments
		WHERE sale_id=$1 AND deleted_at IS NULL
		ORDER BY created_at, id`, saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []domain.PaymentDetail
	for rows.Next() {
		var payment domain.PaymentDetail
		if err := rows.Scan(&payment.ID, &payment.Method, &payment.Amount); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, rows.Err()
}

func (r *SaleRepo) Refund(ctx context.Context, actor domain.User, saleID uuid.UUID, items []domain.CartItemInput) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var branchID uuid.UUID
	if err := tx.QueryRow(ctx, `SELECT branch_id FROM sales WHERE id=$1 AND deleted_at IS NULL`, saleID).Scan(&branchID); err != nil {
		return err
	}
	if actor.Role == domain.RoleEmployee && (actor.BranchID == nil || *actor.BranchID != branchID) {
		return errors.New("employee cannot refund outside assigned branch")
	}
	if actor.Role == domain.RoleManager {
		var canAccess bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM user_branches WHERE user_id=$1 AND branch_id=$2)`, actor.ID, branchID).Scan(&canAccess); err != nil {
			return err
		}
		if !canAccess {
			return errors.New("manager cannot refund outside managed branches")
		}
	}

	for _, item := range items {
		if item.Quantity <= 0 {
			return errors.New("refund quantity must be greater than zero")
		}
		var soldQty, returnedQty int64
		if err := tx.QueryRow(ctx, `
			SELECT COALESCE(SUM(quantity), 0)
			FROM sale_items
			WHERE sale_id=$1 AND product_id=$2 AND deleted_at IS NULL`, saleID, item.ProductID).Scan(&soldQty); err != nil {
			return err
		}
		if soldQty == 0 {
			return errors.New("refund item is not part of this sale")
		}
		if err := tx.QueryRow(ctx, `
			SELECT COALESCE(SUM(quantity), 0)
			FROM inventory_movements
			WHERE reference_id=$1 AND product_id=$2 AND movement_type=$3 AND deleted_at IS NULL`,
			saleID.String(), item.ProductID, domain.MovementReturn).Scan(&returnedQty); err != nil {
			return err
		}
		if item.Quantity > soldQty-returnedQty {
			return errors.New("refund quantity exceeds remaining refundable quantity")
		}
		_, err = tx.Exec(ctx, `UPDATE inventories SET quantity=quantity+$3, updated_at=now() WHERE branch_id=$1 AND product_id=$2`,
			branchID, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO inventory_movements (branch_id, product_id, movement_type, quantity, reference_id, created_by)
			VALUES ($1,$2,$3,$4,$5,$6)`, branchID, item.ProductID, domain.MovementReturn, item.Quantity, saleID.String(), actor.ID)
		if err != nil {
			return err
		}
	}
	_, err = tx.Exec(ctx, `INSERT INTO audit_logs (user_id, action, entity_type, entity_id, new_data) VALUES ($1,'REFUND','sale',$2,$3)`,
		actor.ID, saleID, fmt.Sprintf(`{"items":%d}`, len(items)))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
