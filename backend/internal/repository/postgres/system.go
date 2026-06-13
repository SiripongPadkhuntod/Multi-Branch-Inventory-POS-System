package postgres

import (
	"context"

	"pos-system/backend/internal/app/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SystemRepo struct{ db *pgxpool.Pool }

func (r *SystemRepo) ListCategories(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, description
		FROM categories
		WHERE deleted_at IS NULL
		ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]domain.Category, 0)
	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (r *SystemRepo) CreateCategory(ctx context.Context, category domain.Category) (*domain.Category, error) {
	return scanCategory(r.db.QueryRow(ctx, `
		INSERT INTO categories (name, description)
		VALUES ($1,$2)
		RETURNING id, name, description`, category.Name, category.Description))
}

func (r *SystemRepo) UpdateCategory(ctx context.Context, id uuid.UUID, category domain.Category) (*domain.Category, error) {
	return scanCategory(r.db.QueryRow(ctx, `
		UPDATE categories
		SET name=$2, description=$3, updated_at=now()
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING id, name, description`, id, category.Name, category.Description))
}

func (r *SystemRepo) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE categories SET deleted_at=now(), updated_at=now() WHERE id=$1`, id)
	return err
}

func (r *SystemRepo) CreateBranch(ctx context.Context, branch domain.Branch) (*domain.Branch, error) {
	return scanBranch(r.db.QueryRow(ctx, `
		INSERT INTO branches (code, name, address, phone, status)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, code, name, address, phone, status, created_at`,
		branch.Code, branch.Name, branch.Address, branch.Phone, branch.Status))
}

func (r *SystemRepo) UpdateBranch(ctx context.Context, id uuid.UUID, branch domain.Branch) (*domain.Branch, error) {
	return scanBranch(r.db.QueryRow(ctx, `
		UPDATE branches
		SET code=$2, name=$3, address=$4, phone=$5, status=$6, updated_at=now()
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING id, code, name, address, phone, status, created_at`,
		id, branch.Code, branch.Name, branch.Address, branch.Phone, branch.Status))
}

type categoryScanner interface{ Scan(dest ...any) error }

func scanCategory(row categoryScanner) (*domain.Category, error) {
	var category domain.Category
	err := row.Scan(&category.ID, &category.Name, &category.Description)
	return &category, err
}

func scanBranch(row categoryScanner) (*domain.Branch, error) {
	var branch domain.Branch
	err := row.Scan(&branch.ID, &branch.Code, &branch.Name, &branch.Address, &branch.Phone, &branch.Status, &branch.CreatedAt)
	return &branch, err
}
