package server

import (
	"context"
	"net/http"
	"time"

	"pos-system/backend/internal/app/handler"
	"pos-system/backend/internal/app/port"
	apprepository "pos-system/backend/internal/app/repository"
	postgresrepository "pos-system/backend/internal/app/repository/postgres-repository"
	appservice "pos-system/backend/internal/app/service"
	dbclient "pos-system/backend/internal/infrastructure/db-client"
	ginclient "pos-system/backend/internal/infrastructure/gin-client"
	middlewareclient "pos-system/backend/internal/infrastructure/middleware-client"
	migraterunner "pos-system/backend/internal/infrastructure/migrate"
	"pos-system/backend/internal/property"
	"pos-system/backend/internal/router"
	"pos-system/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg    property.AppConfig
	server *http.Server
	close  func()
}

func Init(ctx context.Context, cfg property.AppConfig) (*App, error) {
	db, err := initPostgresRepository(ctx, cfg)
	if err != nil {
		return nil, err
	}
	repo := initRepository(ctx, db)
	compatServices := usecase.NewServices(repo.SQL, cfg)
	appServices := initService(*repo)
	appHandler := handler.New(appServices)
	engine := initServer(cfg, compatServices, appHandler, appServices)

	return &App{
		cfg: cfg,
		server: &http.Server{
			Addr:              ":" + cfg.Port,
			Handler:           engine,
			ReadHeaderTimeout: 5 * time.Second,
		},
		close: func() { db.Close() },
	}, nil
}

func initRepository(_ context.Context, db *pgxpool.Pool) *port.Repository {
	sql := postgresrepository.NewRepositories(db)
	return apprepository.New(sql)
}

func initService(repo port.Repository) port.Service {
	return appservice.New(repo)
}

func initServer(cfg property.AppConfig, services *usecase.Services, h port.Handler, svc port.Service) *gin.Engine {
	engine := ginclient.GinInit()
	middlewareclient.ServicesMiddleware(engine)
	return router.SetupRouter(cfg, services, svc, h)
}

func initPostgresRepository(ctx context.Context, cfg property.AppConfig) (*pgxpool.Pool, error) {
	db, err := dbclient.ConnectPostgres(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if cfg.RunMigrations {
		if err := migraterunner.RunSQLMigrations(ctx, db, cfg.MigrationsPath); err != nil {
			db.Close()
			return nil, err
		}
	}
	return db, nil
}
