package postgresrepository

import (
	"pos-system/backend/internal/app/port"
	postgresadapter "pos-system/backend/internal/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRepositories(db *pgxpool.Pool) port.Repositories {
	return postgresadapter.NewRepositories(db)
}
