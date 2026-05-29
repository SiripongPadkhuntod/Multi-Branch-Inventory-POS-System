package postgres

import (
	"context"

	"pos-system/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct{ db *pgxpool.Pool }

func (r *AuthRepo) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, branch_id, role, name, email, password_hash, status, created_at
		FROM users WHERE email=$1 AND deleted_at IS NULL`, email)
	return scanUser(row)
}

func (r *AuthRepo) FindUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, branch_id, role, name, email, password_hash, status, created_at
		FROM users WHERE id=$1 AND deleted_at IS NULL`, id)
	return scanUser(row)
}

type userScanner interface{ Scan(dest ...any) error }

func scanUser(row userScanner) (*domain.User, error) {
	var u domain.User
	if err := row.Scan(&u.ID, &u.BranchID, &u.Role, &u.Name, &u.Email, &u.PasswordHash, &u.Status, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
