package ginclient

import (
	"pos-system/backend/internal/delivery/http/middleware"
	"pos-system/backend/internal/property"

	"github.com/gin-gonic/gin"
)

func NewEngine(cfg property.AppConfig) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery(), middleware.CORS())
	return engine
}

func GinInit() *gin.Engine {
	return NewEngine(property.Load())
}
