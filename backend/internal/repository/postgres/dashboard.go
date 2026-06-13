package postgres

import (
	"context"
	"errors"

	"pos-system/backend/internal/app/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepo struct{ db *pgxpool.Pool }

func (r *DashboardRepo) AccessibleBranches(ctx context.Context, actor domain.User) ([]domain.Branch, error) {
	var rows pgxRows
	var err error
	if actor.Role == domain.RoleOwner {
		rows, err = r.db.Query(ctx, `
			SELECT id, code, name, address, phone, status, created_at
			FROM branches
			WHERE deleted_at IS NULL
			ORDER BY code`)
	} else if actor.Role == domain.RoleEmployee {
		rows, err = r.db.Query(ctx, `
			SELECT id, code, name, address, phone, status, created_at
			FROM branches
			WHERE id=$1 AND deleted_at IS NULL
			ORDER BY code`, actor.BranchID)
	} else {
		rows, err = r.db.Query(ctx, `
			SELECT b.id, b.code, b.name, b.address, b.phone, b.status, b.created_at
			FROM branches b
			JOIN user_branches ub ON ub.branch_id=b.id
			WHERE ub.user_id=$1 AND b.deleted_at IS NULL
			UNION
			SELECT id, code, name, address, phone, status, created_at
			FROM branches
			WHERE id=$2 AND $2::uuid IS NOT NULL AND deleted_at IS NULL
			ORDER BY code`, actor.ID, actor.BranchID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	branches := make([]domain.Branch, 0)
	for rows.Next() {
		var branch domain.Branch
		if err := rows.Scan(&branch.ID, &branch.Code, &branch.Name, &branch.Address, &branch.Phone, &branch.Status, &branch.CreatedAt); err != nil {
			return nil, err
		}
		branches = append(branches, branch)
	}
	return branches, rows.Err()
}

func (r *DashboardRepo) Summary(ctx context.Context, actor domain.User, branchID *uuid.UUID) (*domain.DashboardSummary, error) {
	branches, err := r.AccessibleBranches(ctx, actor)
	if err != nil {
		return nil, err
	}
	allowedIDs := make([]uuid.UUID, 0, len(branches))
	for _, branch := range branches {
		allowedIDs = append(allowedIDs, branch.ID)
	}
	if actor.Role != domain.RoleOwner && len(allowedIDs) == 0 {
		return nil, errors.New("no branch assigned")
	}
	if branchID != nil && !containsBranch(allowedIDs, *branchID) && actor.Role != domain.RoleOwner {
		return nil, errors.New("branch access denied")
	}

	summary := &domain.DashboardSummary{
		LowStock:         make([]domain.LowStockItem, 0),
		TopProducts:      make([]domain.TopProduct, 0),
		BranchComparison: make([]domain.BranchSalesSummary, 0),
	}
	filterBranch := branchID
	allowedParam := allowedIDs
	if actor.Role == domain.RoleOwner {
		allowedParam = nil
	}

	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FILTER (WHERE created_at::date = CURRENT_DATE),
			COUNT(*) FILTER (WHERE date_trunc('month', created_at) = date_trunc('month', now())),
			COALESCE(SUM(total), 0)
		FROM sales
		WHERE deleted_at IS NULL
			AND ($1::uuid IS NULL OR branch_id=$1)
			AND ($2::uuid[] IS NULL OR branch_id=ANY($2))`, filterBranch, allowedParam).
		Scan(&summary.DailySales, &summary.MonthlySales, &summary.Revenue); err != nil {
		return nil, err
	}

	if err := r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM((si.final_price - p.cost_price) * si.quantity), 0)
		FROM sale_items si
		JOIN sales s ON s.id=si.sale_id
		JOIN products p ON p.id=si.product_id
		WHERE si.deleted_at IS NULL AND s.deleted_at IS NULL
			AND ($1::uuid IS NULL OR s.branch_id=$1)
			AND ($2::uuid[] IS NULL OR s.branch_id=ANY($2))`, filterBranch, allowedParam).
		Scan(&summary.Profit); err != nil {
		return nil, err
	}

	lowStock, err := r.lowStock(ctx, filterBranch, allowedParam)
	if err != nil {
		return nil, err
	}
	topProducts, err := r.topProducts(ctx, filterBranch, allowedParam)
	if err != nil {
		return nil, err
	}
	branchComparison, err := r.branchComparison(ctx, filterBranch, allowedParam)
	if err != nil {
		return nil, err
	}
	summary.LowStock = lowStock
	summary.TopProducts = topProducts
	summary.BranchComparison = branchComparison
	return summary, nil
}

type pgxRows interface {
	Close()
	Next() bool
	Scan(dest ...any) error
	Err() error
}

func containsBranch(branches []uuid.UUID, branchID uuid.UUID) bool {
	for _, id := range branches {
		if id == branchID {
			return true
		}
	}
	return false
}

func (r *DashboardRepo) lowStock(ctx context.Context, branchID *uuid.UUID, allowedIDs []uuid.UUID) ([]domain.LowStockItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT i.branch_id, b.code, i.product_id, p.name, i.quantity, i.reorder_threshold
		FROM inventories i
		JOIN branches b ON b.id=i.branch_id
		JOIN products p ON p.id=i.product_id
		WHERE i.deleted_at IS NULL AND i.quantity <= i.reorder_threshold
			AND ($1::uuid IS NULL OR i.branch_id=$1)
			AND ($2::uuid[] IS NULL OR i.branch_id=ANY($2))
		ORDER BY (i.reorder_threshold - i.quantity) DESC, i.quantity, p.name
		LIMIT 20`, branchID, allowedIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.LowStockItem, 0)
	for rows.Next() {
		var item domain.LowStockItem
		if err := rows.Scan(&item.BranchID, &item.BranchCode, &item.ProductID, &item.ProductName, &item.Quantity, &item.ReorderThreshold); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *DashboardRepo) topProducts(ctx context.Context, branchID *uuid.UUID, allowedIDs []uuid.UUID) ([]domain.TopProduct, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.name, COALESCE(SUM(si.quantity), 0), COALESCE(SUM(si.quantity * si.final_price), 0)
		FROM sale_items si
		JOIN sales s ON s.id=si.sale_id
		JOIN products p ON p.id=si.product_id
		WHERE si.deleted_at IS NULL AND s.deleted_at IS NULL
			AND ($1::uuid IS NULL OR s.branch_id=$1)
			AND ($2::uuid[] IS NULL OR s.branch_id=ANY($2))
		GROUP BY p.id, p.name
		ORDER BY SUM(si.quantity) DESC, p.name
		LIMIT 5`, branchID, allowedIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := make([]domain.TopProduct, 0)
	for rows.Next() {
		var product domain.TopProduct
		if err := rows.Scan(&product.ProductID, &product.Name, &product.Quantity, &product.Revenue); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (r *DashboardRepo) branchComparison(ctx context.Context, branchID *uuid.UUID, allowedIDs []uuid.UUID) ([]domain.BranchSalesSummary, error) {
	rows, err := r.db.Query(ctx, `
		SELECT b.id, b.code, b.name, COALESCE(SUM(s.total), 0), COUNT(s.id)
		FROM branches b
		LEFT JOIN sales s ON s.branch_id=b.id AND s.deleted_at IS NULL
		WHERE b.deleted_at IS NULL
			AND ($1::uuid IS NULL OR b.id=$1)
			AND ($2::uuid[] IS NULL OR b.id=ANY($2))
		GROUP BY b.id, b.code, b.name
		ORDER BY COALESCE(SUM(s.total), 0) DESC, b.code`, branchID, allowedIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	branches := make([]domain.BranchSalesSummary, 0)
	for rows.Next() {
		var branch domain.BranchSalesSummary
		if err := rows.Scan(&branch.BranchID, &branch.BranchCode, &branch.BranchName, &branch.Revenue, &branch.SalesCount); err != nil {
			return nil, err
		}
		branches = append(branches, branch)
	}
	return branches, rows.Err()
}
