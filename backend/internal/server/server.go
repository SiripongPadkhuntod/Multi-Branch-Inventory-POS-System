package server

import (
	"context"

	"pos-system/backend/internal/infrastructure/logger-client"
	"pos-system/backend/internal/property"
)

func RunServer(ctx context.Context) error {
	loggerclient.InitGlobalLogger(true, true)
	defer loggerclient.FromContext(ctx).Sync()

	cfg := property.Load()
	ctx = loggerclient.WithContext(ctx, "service", property.Get().Server.ServiceName)
	loggerclient.GetLogger(ctx, "startup").Infow("service starting", "version", property.Get().Server.ServiceVersion)

	app, err := Init(ctx, cfg)
	if err != nil {
		return err
	}
	defer app.Close()
	return app.Listen(ctx)
}

func (a *App) Close() {
	if a.close != nil {
		a.close()
	}
}
