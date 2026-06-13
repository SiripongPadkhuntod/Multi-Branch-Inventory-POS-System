package property

import (
	"sync"

	"pos-system/backend/internal/constant"
	"pos-system/backend/internal/infrastructure/config"
)

// AppConfig is the application settings contract used by the server layer.
// It currently wraps the existing config package while the project migrates
// toward envconfig-style properties.
type AppConfig = config.Config

func Load() AppConfig {
	return config.Load()
}

var (
	schemaOnce sync.Once
	schema     Schema
)

func Get() Schema {
	schemaOnce.Do(func() {
		cfg := Load()
		schema = Schema{
			Server: ServerProperty{
				ServiceName:        constant.ServiceName,
				ServiceDescription: constant.ServiceDescription,
				ServiceVersion:     constant.ServiceVersion,
				Port:               cfg.Port,
				AppEnv:             cfg.AppEnv,
			},
			Swagger: SwaggerProperty{
				ApiDocs:        true,
				ApiDocsVersion: constant.ServiceVersion,
				ApiDocsHost:    "",
				ApiDocsSchema:  "http",
			},
		}
	})
	return schema
}
