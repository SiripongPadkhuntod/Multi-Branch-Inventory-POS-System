package dbclient

import (
	"context"

	"pos-system/backend/internal/infrastructure/database"
	"pos-system/backend/internal/property"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres(ctx context.Context, cfg property.AppConfig) (*pgxpool.Pool, error) {
	return database.Connect(ctx, cfg)
}
