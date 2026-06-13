package postgres

import (
	"context"
	"errors"

	"pos-system/backend/internal/app/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct{ db *pgxpool.Pool }

func (r *ProductRepo) List(ctx context.Context, query string, limit, offset int) ([]domain.Product, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, sku, barcode, qr_code, name, description, category_id, image_url, cost_price, sell_price, status
		FROM products
		WHERE deleted_at IS NULL
		  AND ($1='' OR sku ILIKE '%'||$1||'%' OR barcode ILIKE '%'||$1||'%' OR name ILIKE '%'||$1||'%')
		ORDER BY name
		LIMIT $2 OFFSET $3`, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}
	return products, rows.Err()
}

func (r *ProductRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return scanProduct(r.db.QueryRow(ctx, `
		SELECT id, sku, barcode, qr_code, name, description, category_id, image_url, cost_price, sell_price, status
		FROM products WHERE id=$1 AND deleted_at IS NULL`, id))
}

func (r *ProductRepo) FindByBarcode(ctx context.Context, barcode string) (*domain.Product, error) {
	return scanProduct(r.db.QueryRow(ctx, `
		SELECT id, sku, barcode, qr_code, name, description, category_id, image_url, cost_price, sell_price, status
		FROM products WHERE barcode=$1 AND deleted_at IS NULL`, barcode))
}

func (r *ProductRepo) Create(ctx context.Context, p domain.Product) (*domain.Product, error) {
	if p.Status == "" {
		p.Status = "ACTIVE"
	}
	if p.CostPrice > p.SellPrice {
		return nil, errors.New("cost price should not exceed sell price")
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	product, err := scanProduct(tx.QueryRow(ctx, `
		INSERT INTO products (sku, barcode, qr_code, name, description, category_id, image_url, cost_price, sell_price, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id, sku, barcode, qr_code, name, description, category_id, image_url, cost_price, sell_price, status`,
		p.SKU, p.Barcode, p.QRCode, p.Name, p.Description, p.CategoryID, p.ImageURL, p.CostPrice, p.SellPrice, p.Status))
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO inventories (branch_id, product_id, quantity)
		SELECT id, $1, 0 FROM branches WHERE deleted_at IS NULL
		ON CONFLICT DO NOTHING`, product.ID); err != nil {
		return nil, err
	}
	return product, tx.Commit(ctx)
}

func (r *ProductRepo) Update(ctx context.Context, id uuid.UUID, p domain.Product) (*domain.Product, error) {
	if p.Status == "" {
		p.Status = "ACTIVE"
	}
	if p.CostPrice > p.SellPrice {
		return nil, errors.New("cost price should not exceed sell price")
	}
	return scanProduct(r.db.QueryRow(ctx, `
		UPDATE products SET sku=$2, barcode=$3, qr_code=$4, name=$5, description=$6, category_id=$7,
			image_url=$8, cost_price=$9, sell_price=$10, status=$11, updated_at=now()
		WHERE id=$1 AND deleted_at IS NULL
		RETURNING id, sku, barcode, qr_code, name, description, category_id, image_url, cost_price, sell_price, status`,
		id, p.SKU, p.Barcode, p.QRCode, p.Name, p.Description, p.CategoryID, p.ImageURL, p.CostPrice, p.SellPrice, p.Status))
}

func (r *ProductRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE products SET deleted_at=now(), updated_at=now() WHERE id=$1`, id)
	return err
}

type productScanner interface{ Scan(dest ...any) error }

func scanProduct(row productScanner) (*domain.Product, error) {
	var p domain.Product
	err := row.Scan(&p.ID, &p.SKU, &p.Barcode, &p.QRCode, &p.Name, &p.Description, &p.CategoryID, &p.ImageURL, &p.CostPrice, &p.SellPrice, &p.Status)
	return &p, err
}
