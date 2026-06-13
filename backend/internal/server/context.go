package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pos-system/backend/internal/constant"
)

func (a *App) Listen(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		log.Printf("Start server on port %s at %s\n", a.cfg.Port, time.Now().Local().Format(constant.TimeLayoutDateTimeFormat))
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-quit:
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return a.server.Shutdown(shutdownCtx)
}
