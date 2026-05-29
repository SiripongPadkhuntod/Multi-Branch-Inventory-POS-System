package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpdelivery "pos-system/backend/internal/delivery/http"
	"pos-system/backend/internal/infrastructure/config"
	"pos-system/backend/internal/infrastructure/database"
	"pos-system/backend/internal/repository/postgres"
	"pos-system/backend/internal/usecase"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(context.Background(), cfg)
	if err != nil {
		log.Fatalf("database connect: %v", err)
	}
	defer db.Close()

	repos := postgres.NewRepositories(db)
	services := usecase.NewServices(repos, cfg)
	router := httpdelivery.NewRouter(cfg, services)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("api listening on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}
