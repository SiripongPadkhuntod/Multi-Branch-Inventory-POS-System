package migrate

import (
	"context"

	"pos-system/backend/internal/infrastructure/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunSQLMigrations(ctx context.Context, db *pgxpool.Pool, migrationsPath string) error {
	return database.RunMigrations(ctx, db, migrationsPath)
}
