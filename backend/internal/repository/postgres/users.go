package postgres

import (
	"context"
	"errors"

	"pos-system/backend/internal/app/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct{ db *pgxpool.Pool }

func (r *UserRepo) List(ctx context.Context, actor domain.User) ([]domain.User, error) {
	var users []domain.User
	var err error
	if actor.Role == domain.RoleOwner {
		users, err = r.listOwner(ctx)
	} else {
		users, err = r.listManaged(ctx, actor)
	}
	if err != nil {
		return nil, err
	}
	if err := r.attachBranchIDs(ctx, users); err != nil {
		return nil, err
	}
	return users, nil
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
	normalizeUserBranches(&user)
	if err := r.canManageUser(ctx, actor, user); err != nil {
		return nil, err
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	created, err := scanUser(tx.QueryRow(ctx, `
		INSERT INTO users (branch_id, role, name, email, password_hash, status)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, branch_id, role, name, email, password_hash, status, created_at`,
		user.BranchID, user.Role, user.Name, user.Email, user.PasswordHash, user.Status))
	if err != nil {
		return nil, err
	}
	created.BranchIDs = user.BranchIDs
	if err := syncManagedBranches(ctx, tx, *created); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return created, nil
}

func (r *UserRepo) Update(ctx context.Context, actor domain.User, id uuid.UUID, user domain.User) (*domain.User, error) {
	normalizeUserBranches(&user)
	if err := r.canManageUser(ctx, actor, user); err != nil {
		return nil, err
	}
	if actor.Role == domain.RoleManager && !r.canAccessUser(ctx, actor, id) {
		return nil, errors.New("employee access denied")
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	updated, err := scanUser(tx.QueryRow(ctx, `
		UPDATE users SET branch_id=$2, role=$3, name=$4, email=$5, status=$6, updated_at=now()
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING id, branch_id, role, name, email, password_hash, status, created_at`,
		id, user.BranchID, user.Role, user.Name, user.Email, user.Status))
	if err != nil {
		return nil, err
	}
	updated.BranchIDs = user.BranchIDs
	if err := syncManagedBranches(ctx, tx, *updated); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return updated, nil
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
		if user.Role == domain.RoleEmployee && user.BranchID == nil {
			return errors.New("branch_id is required")
		}
		if user.Role == domain.RoleManager && len(user.BranchIDs) == 0 {
			return errors.New("manager must be assigned to at least one branch")
		}
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

func (r *UserRepo) attachBranchIDs(ctx context.Context, users []domain.User) error {
	if len(users) == 0 {
		return nil
	}
	userIDs := make([]uuid.UUID, 0, len(users))
	indexByID := make(map[uuid.UUID]int, len(users))
	for index, user := range users {
		userIDs = append(userIDs, user.ID)
		indexByID[user.ID] = index
		if user.BranchID != nil {
			users[index].BranchIDs = []uuid.UUID{*user.BranchID}
		}
	}
	rows, err := r.db.Query(ctx, `
		SELECT user_id, branch_id
		FROM user_branches
		WHERE user_id=ANY($1::uuid[])
		ORDER BY created_at`, userIDs)
	if err != nil {
		return err
	}
	defer rows.Close()

	seen := make(map[uuid.UUID]map[uuid.UUID]bool, len(users))
	for rows.Next() {
		var userID uuid.UUID
		var branchID uuid.UUID
		if err := rows.Scan(&userID, &branchID); err != nil {
			return err
		}
		index, ok := indexByID[userID]
		if !ok {
			continue
		}
		if seen[userID] == nil {
			seen[userID] = make(map[uuid.UUID]bool)
			users[index].BranchIDs = make([]uuid.UUID, 0)
		}
		if seen[userID][branchID] {
			continue
		}
		users[index].BranchIDs = append(users[index].BranchIDs, branchID)
		seen[userID][branchID] = true
	}
	return rows.Err()
}

func normalizeUserBranches(user *domain.User) {
	if user.Role != domain.RoleManager {
		if user.BranchID != nil {
			user.BranchIDs = []uuid.UUID{*user.BranchID}
		} else {
			user.BranchIDs = nil
		}
		return
	}

	branchIDs := make([]uuid.UUID, 0, len(user.BranchIDs)+1)
	seen := make(map[uuid.UUID]bool)
	for _, branchID := range user.BranchIDs {
		if seen[branchID] {
			continue
		}
		branchIDs = append(branchIDs, branchID)
		seen[branchID] = true
	}
	if user.BranchID != nil && !seen[*user.BranchID] {
		branchIDs = append([]uuid.UUID{*user.BranchID}, branchIDs...)
	}
	user.BranchIDs = branchIDs
	if len(branchIDs) > 0 {
		user.BranchID = &user.BranchIDs[0]
	}
}

func syncManagedBranches(ctx context.Context, tx pgx.Tx, user domain.User) error {
	if _, err := tx.Exec(ctx, `DELETE FROM user_branches WHERE user_id=$1`, user.ID); err != nil {
		return err
	}
	if user.Role != domain.RoleManager {
		return nil
	}
	branchIDs := user.BranchIDs
	if len(branchIDs) == 0 && user.BranchID != nil {
		branchIDs = []uuid.UUID{*user.BranchID}
	}
	for _, branchID := range branchIDs {
		if _, err := tx.Exec(ctx, `
			INSERT INTO user_branches (user_id, branch_id)
			VALUES ($1,$2)
			ON CONFLICT (user_id, branch_id) DO NOTHING`, user.ID, branchID); err != nil {
			return err
		}
	}
	return nil
}
