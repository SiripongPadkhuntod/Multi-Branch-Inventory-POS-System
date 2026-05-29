package postgres

import (
	"context"
	"errors"

	"pos-system/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct{ db *pgxpool.Pool }

func (r *UserRepo) List(ctx context.Context, actor domain.User) ([]domain.User, error) {
	if actor.Role == domain.RoleOwner {
		return r.listOwner(ctx)
	}
	return r.listManaged(ctx, actor)
}

func (r *UserRepo) listOwner(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, branch_id, role, name, email, password_hash, status, created_at
		FROM users WHERE deleted_at IS NULL ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []domain.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, rows.Err()
}

func (r *UserRepo) listManaged(ctx context.Context, actor domain.User) ([]domain.User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT u.id, u.branch_id, u.role, u.name, u.email, u.password_hash, u.status, u.created_at
		FROM users u
		JOIN user_branches ub ON ub.branch_id=u.branch_id
		WHERE ub.user_id=$1 AND u.role='EMPLOYEE' AND u.deleted_at IS NULL
		ORDER BY u.name`, actor.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]domain.User, 0)
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, rows.Err()
}

func (r *UserRepo) Create(ctx context.Context, actor domain.User, user domain.User) (*domain.User, error) {
	if err := r.canManageUser(ctx, actor, user); err != nil {
		return nil, err
	}
	return scanUser(r.db.QueryRow(ctx, `
		INSERT INTO users (branch_id, role, name, email, password_hash, status)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, branch_id, role, name, email, password_hash, status, created_at`,
		user.BranchID, user.Role, user.Name, user.Email, user.PasswordHash, user.Status))
}

func (r *UserRepo) Update(ctx context.Context, actor domain.User, id uuid.UUID, user domain.User) (*domain.User, error) {
	if err := r.canManageUser(ctx, actor, user); err != nil {
		return nil, err
	}
	if actor.Role == domain.RoleManager && !r.canAccessUser(ctx, actor, id) {
		return nil, errors.New("employee access denied")
	}
	return scanUser(r.db.QueryRow(ctx, `
		UPDATE users SET branch_id=$2, role=$3, name=$4, email=$5, status=$6, updated_at=now()
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING id, branch_id, role, name, email, password_hash, status, created_at`,
		id, user.BranchID, user.Role, user.Name, user.Email, user.Status))
}

func (r *UserRepo) SalesSummary(ctx context.Context, actor domain.User) ([]domain.EmployeeSalesSummary, error) {
	var rows pgxRows
	var err error
	if actor.Role == domain.RoleOwner {
		rows, err = r.db.Query(ctx, salesSummarySQL()+`
			WHERE u.deleted_at IS NULL AND u.role IN ('EMPLOYEE','MANAGER')
			GROUP BY u.id, u.branch_id, b.code, u.name, u.email, u.role, u.status
			ORDER BY revenue DESC, u.name`)
	} else {
		rows, err = r.db.Query(ctx, salesSummarySQL()+`
			JOIN user_branches ub ON ub.branch_id=u.branch_id
			WHERE ub.user_id=$1 AND u.deleted_at IS NULL AND u.role='EMPLOYEE'
			GROUP BY u.id, u.branch_id, b.code, u.name, u.email, u.role, u.status
			ORDER BY revenue DESC, u.name`, actor.ID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	summaries := make([]domain.EmployeeSalesSummary, 0)
	for rows.Next() {
		var item domain.EmployeeSalesSummary
		if err := rows.Scan(&item.UserID, &item.BranchID, &item.BranchCode, &item.Name, &item.Email, &item.Role, &item.Status, &item.SalesCount, &item.Revenue); err != nil {
			return nil, err
		}
		summaries = append(summaries, item)
	}
	return summaries, rows.Err()
}

func salesSummarySQL() string {
	return `
		SELECT u.id, u.branch_id, b.code, u.name, u.email, u.role, u.status,
			COUNT(s.id) AS sales_count,
			COALESCE(SUM(s.total), 0) AS revenue
		FROM users u
		JOIN branches b ON b.id=u.branch_id
		LEFT JOIN sales s ON s.employee_id=u.id AND s.deleted_at IS NULL
		`
}

func (r *UserRepo) canManageUser(ctx context.Context, actor domain.User, user domain.User) error {
	if actor.Role == domain.RoleOwner {
		return nil
	}
	if actor.Role != domain.RoleManager {
		return errors.New("manager permission required")
	}
	if user.Role != domain.RoleEmployee {
		return errors.New("manager can manage employees only")
	}
	if user.BranchID == nil {
		return errors.New("branch_id is required")
	}
	if !r.canAccessBranch(ctx, actor, *user.BranchID) {
		return errors.New("branch access denied")
	}
	return nil
}

func (r *UserRepo) canAccessUser(ctx context.Context, actor domain.User, userID uuid.UUID) bool {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM users u
			JOIN user_branches ub ON ub.branch_id=u.branch_id
			WHERE u.id=$1 AND ub.user_id=$2 AND u.role='EMPLOYEE' AND u.deleted_at IS NULL
		)`, userID, actor.ID).Scan(&exists)
	return err == nil && exists
}

func (r *UserRepo) canAccessBranch(ctx context.Context, actor domain.User, branchID uuid.UUID) bool {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM user_branches WHERE user_id=$1 AND branch_id=$2)`, actor.ID, branchID).Scan(&exists)
	return err == nil && exists
}
