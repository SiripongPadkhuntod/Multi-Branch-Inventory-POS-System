package postgres

import (
	"context"
	"time"

	"pos-system/backend/internal/app/domain"

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

func (r *AuthRepo) StoreRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)`, userID, tokenHash, expiresAt)
	return err
}

func (r *AuthRepo) IsRefreshTokenActive(ctx context.Context, userID uuid.UUID, tokenHash string) (bool, error) {
	var active bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM refresh_tokens
			WHERE user_id=$1
			  AND token_hash=$2
			  AND revoked_at IS NULL
			  AND deleted_at IS NULL
			  AND expires_at > now()
		)`, userID, tokenHash).Scan(&active)
	return active, err
}

func (r *AuthRepo) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE refresh_tokens
		SET revoked_at=COALESCE(revoked_at, now()), updated_at=now()
		WHERE token_hash=$1 AND deleted_at IS NULL`, tokenHash)
	return err
}

type userScanner interface{ Scan(dest ...any) error }

func scanUser(row userScanner) (*domain.User, error) {
	var u domain.User
	if err := row.Scan(&u.ID, &u.BranchID, &u.Role, &u.Name, &u.Email, &u.PasswordHash, &u.Status, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
