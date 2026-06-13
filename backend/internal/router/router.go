package router

import (
	"pos-system/backend/internal/app/port"
	httpdelivery "pos-system/backend/internal/delivery/http"
	"pos-system/backend/internal/property"
	"pos-system/backend/internal/usecase"
	"pos-system/backend/pkg/v1/dto"

	"github.com/gin-gonic/gin"
)

const ServiceBaseURL = "/pos-service"

func SetupRouter(cfg property.AppConfig, services *usecase.Services, svc port.Service, h port.Handler) *gin.Engine {
	engine := httpdelivery.NewRouter(cfg, services)
	SetupTemplateRouter(engine, svc, h)
	return engine
}

func SetupTemplateRouter(engine *gin.Engine, svc port.Service, h port.Handler) {
	v1 := engine.Group(ServiceBaseURL + "/v1")
	v1.GET("/health", func(c *gin.Context) {
		BindReqJson200Resp[dto.EmptyStruct, *dto.HealthResponse](c, h.HealthHandler)
	})
	user := v1.Group("/user")
	user.POST("/login", func(c *gin.Context) {
		BindReqJson200Resp[dto.LoginRequest, *dto.LoginResponse](c, h.LoginHandler)
	})
	user.POST("/refresh", func(c *gin.Context) {
		BindReqJson200Resp[dto.RefreshRequest, *dto.LoginResponse](c, h.RefreshHandler)
	})
	user.POST("/logout", func(c *gin.Context) {
		BindReqJson200Resp[dto.LogoutRequest, *dto.SuccessResponse](c, h.LogoutHandler)
	})

	protected := v1.Group("")
	protected.Use(AuthMiddleware(svc))

	products := protected.Group("/products")
	products.GET("", func(c *gin.Context) {
		BindReqJson200Resp[dto.ProductListRequest, *dto.ProductListResponse](c, h.ProductListHandler)
	})
	products.GET("/barcode/:barcode", func(c *gin.Context) {
		BindReqJson200Resp[dto.ProductBarcodeRequest, *dto.ProductResponse](c, h.ProductBarcodeHandler)
	})

	inventories := protected.Group("/inventories")
	inventories.GET("", func(c *gin.Context) {
		BindReqJson200Resp[dto.InventoryListRequest, *dto.InventoryListResponse](c, h.InventoryListHandler)
	})
	inventories.GET("/movements", func(c *gin.Context) {
		BindReqJson200Resp[dto.InventoryMovementsRequest, *dto.InventoryMovementsResponse](c, h.InventoryMovementsHandler)
	})
	inventories.GET("/all-stock", func(c *gin.Context) {
		BindReqJson200Resp[dto.InventoryAllStockRequest, *dto.InventoryAllStockResponse](c, h.InventoryAllStockHandler)
	})
	inventories.POST("/adjust", func(c *gin.Context) {
		BindReqJson200Resp[dto.InventoryAdjustRequest, *dto.SuccessResponse](c, h.InventoryAdjustHandler)
	})
	inventories.POST("/receive", func(c *gin.Context) {
		BindReqJson200Resp[dto.InventoryAdjustRequest, *dto.SuccessResponse](c, h.InventoryReceiveHandler)
	})
}
