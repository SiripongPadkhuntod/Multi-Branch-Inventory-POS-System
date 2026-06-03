package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"pos-system/backend/internal/delivery/http/middleware"
	"pos-system/backend/internal/domain"
	"pos-system/backend/internal/infrastructure/config"
	"pos-system/backend/internal/infrastructure/storage"
	"pos-system/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewRouter(cfg config.Config, services *usecase.Services) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	minioStorage, err := storage.NewMinioStorage(cfg)
	if err != nil {
		log.Printf("Warning: failed to initialize MinIO storage: %v", err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS())
	r.Static("/uploads", "./uploads")
	registerSwagger(r)

	v1 := r.Group("/api/v1")
	registerAuth(v1, cfg, services)

	protected := v1.Group("")
	protected.Use(middleware.Auth(services.Auth))
	registerProducts(protected, services, minioStorage)
	registerInventories(protected, services)
	registerSales(protected, services)
	registerUsers(protected, services)
	registerBranches(protected, services)
	registerCategories(protected, services)
	registerDashboard(protected, services)
	registerAuditLogs(protected, services)

	r.NoRoute(func(c *gin.Context) {
		fail(c, http.StatusNotFound, errors.New("route not found"))
	})

	return r
}

func registerAuth(r *gin.RouterGroup, cfg config.Config, services *usecase.Services) {
	r.POST("/auth/login", func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		access, refresh, user, err := services.Auth.Login(c.Request.Context(), input.Email, input.Password)
		if err != nil {
			fail(c, http.StatusUnauthorized, err)
			return
		}
		c.SetCookie("access_token", access, int(cfg.AccessTTL.Seconds()), "/", "", cfg.CookieSecure, true)
		c.SetCookie("refresh_token", refresh, int(cfg.RefreshTTL.Seconds()), "/", "", cfg.CookieSecure, true)
		ok(c, "Login success", gin.H{"access_token": access, "refresh_token": refresh, "user": user})
	})
	r.POST("/auth/logout", func(c *gin.Context) {
		c.SetCookie("access_token", "", -1, "/", "", cfg.CookieSecure, true)
		c.SetCookie("refresh_token", "", -1, "/", "", cfg.CookieSecure, true)
		ok(c, "Logout success", nil)
	})
	r.POST("/auth/refresh", func(c *gin.Context) { ok(c, "Refresh endpoint ready", nil) })
	r.GET("/auth/me", middleware.Auth(services.Auth), func(c *gin.Context) {
		ok(c, "Success", c.MustGet(middleware.CurrentUserKey))
	})
}

func registerProducts(r *gin.RouterGroup, services *usecase.Services, storage *storage.MinioStorage) {
	r.GET("/products", func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		products, err := services.Products.List(c.Request.Context(), c.Query("q"), limit, offset)
		if err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}
		ok(c, "Success", products)
	})
	r.GET("/products/barcode/:barcode", func(c *gin.Context) {
		product, err := services.Products.ByBarcode(c.Request.Context(), c.Param("barcode"))
		if err != nil {
			fail(c, http.StatusNotFound, err)
			return
		}
		ok(c, "Success", product)
	})
	r.GET("/products/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		product, err := services.Products.ByID(c.Request.Context(), id)
		if err != nil {
			fail(c, http.StatusNotFound, err)
			return
		}
		ok(c, "Success", product)
	})
	owner := r.Group("/products", middleware.OwnerOnly())
	owner.POST("/upload-image", func(c *gin.Context) {
		fileHeader, err := c.FormFile("image")
		if err != nil {
			fail(c, http.StatusBadRequest, errors.New("image file is required"))
			return
		}
		if fileHeader.Size > 5*1024*1024 {
			fail(c, http.StatusBadRequest, errors.New("image file must be 5MB or smaller"))
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		defer file.Close()

		header := make([]byte, 512)
		n, err := file.Read(header)
		if err != nil && err != io.EOF {
			fail(c, http.StatusBadRequest, err)
			return
		}
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}

		contentType := http.DetectContentType(header[:n])
		extensions := map[string]string{
			"image/jpeg": ".jpg",
			"image/png":  ".png",
			"image/webp": ".webp",
			"image/gif":  ".gif",
		}
		ext, supported := extensions[contentType]
		if !supported {
			fail(c, http.StatusBadRequest, errors.New("only JPG, PNG, WEBP, or GIF images are supported"))
			return
		}
		originalExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
		switch originalExt {
		case ".jpg", ".jpeg", ".png", ".webp", ".gif":
			ext = originalExt
		}

		if storage == nil {
			fail(c, http.StatusInternalServerError, errors.New("storage service not initialized"))
			return
		}

		filename := fmt.Sprintf("%s%s", uuid.NewString(), ext)
		url, err := storage.UploadFile(c.Request.Context(), filename, file, fileHeader.Size, contentType)
		if err != nil {
			fail(c, http.StatusInternalServerError, fmt.Errorf("failed to upload image: %w", err))
			return
		}
		ok(c, "Image uploaded", gin.H{"image_url": url})
	})
	owner.POST("", func(c *gin.Context) {
		var input domain.Product
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		product, err := services.Products.Create(c.Request.Context(), input)
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		created(c, "Product created", product)
	})
	owner.PUT("/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		var input domain.Product
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		product, err := services.Products.Update(c.Request.Context(), id, input)
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		ok(c, "Product updated", product)
	})
	owner.DELETE("/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		if err := services.Products.Delete(c.Request.Context(), id); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		ok(c, "Product deleted", nil)
	})
}

func registerInventories(r *gin.RouterGroup, services *usecase.Services) {
	r.GET("/inventories/movements", middleware.OwnerOnly(), func(c *gin.Context) {
		var branchID *uuid.UUID
		if c.Query("branch_id") != "" {
			id, err := uuid.Parse(c.Query("branch_id"))
			if err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}
			branchID = &id
		}
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "150"))
		items, err := services.Inventory.Movements(c.Request.Context(), branchID, c.Query("q"), limit)
		if err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}
		ok(c, "Success", items)
	})
	r.GET("/inventories/all-stock", middleware.OwnerOnly(), func(c *gin.Context) {
		items, err := services.Inventory.AllStock(c.Request.Context(), c.Query("q"))
		if err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}
		ok(c, "Success", items)
	})
	r.GET("/inventories", func(c *gin.Context) {
		user := c.MustGet(middleware.CurrentUserKey).(domain.User)
		var branchID *uuid.UUID
		if user.Role == domain.RoleEmployee {
			branchID = user.BranchID
		} else if c.Query("branch_id") != "" {
			id, err := uuid.Parse(c.Query("branch_id"))
			if err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}
			branchID = &id
		}
		if user.Role == domain.RoleManager {
			if branchID == nil {
				branchID = user.BranchID
			}
			if branchID == nil || !canAccessBranch(c.Request.Context(), services, user, *branchID) {
				fail(c, http.StatusForbidden, errors.New("manager cannot view this branch inventory"))
				return
			}
		}
		var categoryID *uuid.UUID
		if c.Query("category_id") != "" {
			id, err := uuid.Parse(c.Query("category_id"))
			if err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}
			categoryID = &id
		}
		items, err := services.Inventory.List(c.Request.Context(), branchID, categoryID, c.Query("q"))
		if err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}
		ok(c, "Success", items)
	})
	r.POST("/inventories/adjust", middleware.OwnerOnly(), adjustInventory(services, "Stock adjusted"))
	r.POST("/inventories/receive", middleware.OwnerOrManager(), receiveInventory(services))
	r.POST("/inventories/reorder-threshold", middleware.OwnerOrManager(), func(c *gin.Context) {
		user := c.MustGet(middleware.CurrentUserKey).(domain.User)
		var input struct {
			BranchID         uuid.UUID `json:"branch_id" binding:"required"`
			ProductID        uuid.UUID `json:"product_id" binding:"required"`
			ReorderThreshold int64     `json:"reorder_threshold" binding:"min=0"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		if user.Role == domain.RoleManager && !canAccessBranch(c.Request.Context(), services, user, input.BranchID) {
			fail(c, http.StatusForbidden, errors.New("manager cannot update this branch threshold"))
			return
		}
		if err := services.Inventory.SetReorderThreshold(c.Request.Context(), input.BranchID, input.ProductID, input.ReorderThreshold); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		ok(c, "Reorder threshold updated", nil)
	})
	r.POST("/inventories/transfer", middleware.OwnerOnly(), func(c *gin.Context) {
		user := c.MustGet(middleware.CurrentUserKey).(domain.User)
		var input struct {
			FromBranchID uuid.UUID `json:"from_branch_id" binding:"required"`
			ToBranchID   uuid.UUID `json:"to_branch_id" binding:"required"`
			ProductID    uuid.UUID `json:"product_id" binding:"required"`
			Quantity     int64     `json:"quantity" binding:"required,min=1"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		if err := services.Inventory.Transfer(c.Request.Context(), input.FromBranchID, input.ToBranchID, input.ProductID, user.ID, input.Quantity); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		ok(c, "Stock transferred", nil)
	})
}

func receiveInventory(services *usecase.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet(middleware.CurrentUserKey).(domain.User)
		var input struct {
			BranchID      uuid.UUID `json:"branch_id"`
			ProductID     uuid.UUID `json:"product_id" binding:"required"`
			QuantityDelta int64     `json:"quantity_delta" binding:"required,min=1"`
			Reason        string    `json:"reason"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		branchID := input.BranchID
		if user.Role == domain.RoleManager {
			if branchID == uuid.Nil && user.BranchID != nil {
				branchID = *user.BranchID
			}
			if branchID == uuid.Nil || !canAccessBranch(c.Request.Context(), services, user, branchID) {
				fail(c, http.StatusForbidden, errors.New("manager cannot receive stock into this branch"))
				return
			}
		}
		if branchID == uuid.Nil {
			fail(c, http.StatusBadRequest, errors.New("branch_id is required"))
			return
		}
		if err := services.Inventory.Adjust(c.Request.Context(), branchID, input.ProductID, user.ID, input.QuantityDelta, input.Reason); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		ok(c, "Stock received", nil)
	}
}

func adjustInventory(services *usecase.Services, message string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet(middleware.CurrentUserKey).(domain.User)
		var input struct {
			BranchID      uuid.UUID `json:"branch_id" binding:"required"`
			ProductID     uuid.UUID `json:"product_id" binding:"required"`
			QuantityDelta int64     `json:"quantity_delta" binding:"required"`
			Reason        string    `json:"reason"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		if err := services.Inventory.Adjust(c.Request.Context(), input.BranchID, input.ProductID, user.ID, input.QuantityDelta, input.Reason); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		ok(c, message, nil)
	}
}

func registerSales(r *gin.RouterGroup, services *usecase.Services) {
	r.POST("/sales", func(c *gin.Context) {
		user := c.MustGet(middleware.CurrentUserKey).(domain.User)
		var input domain.CreateSaleInput
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		sale, err := services.Sales.Create(c.Request.Context(), user, input)
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		created(c, "Sale created", sale)
	})
	r.GET("/sales", func(c *gin.Context) {
		user := c.MustGet(middleware.CurrentUserKey).(domain.User)
		var branchID *uuid.UUID
		if c.Query("branch_id") != "" {
			id, err := uuid.Parse(c.Query("branch_id"))
			if err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}
			branchID = &id
		}
		dateFrom, dateTo, err := parseDateRange(c.Query("date_from"), c.Query("date_to"))
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		sales, err := services.Sales.List(c.Request.Context(), user, branchID, dateFrom, dateTo)
		if err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}
		ok(c, "Success", sales)
	})
	r.GET("/sales/branch", middleware.OwnerOrManager(), func(c *gin.Context) {
		user := c.MustGet(middleware.CurrentUserKey).(domain.User)
		var branchID *uuid.UUID
		if c.Query("branch_id") != "" {
			id, err := uuid.Parse(c.Query("branch_id"))
			if err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}
			branchID = &id
		}
		sales, err := services.Sales.BranchList(c.Request.Context(), user, branchID)
		if err != nil {
			fail(c, http.StatusForbidden, err)
			return
		}
		ok(c, "Success", sales)
	})
	r.GET("/sales/:id", func(c *gin.Context) {
		user := c.MustGet(middleware.CurrentUserKey).(domain.User)
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		sale, err := services.Sales.Detail(c.Request.Context(), user, id)
		if err != nil {
			fail(c, http.StatusNotFound, err)
			return
		}
		ok(c, "Success", sale)
	})
	r.POST("/sales/refund", func(c *gin.Context) {
		user := c.MustGet(middleware.CurrentUserKey).(domain.User)
		var input struct {
			SaleID uuid.UUID              `json:"sale_id" binding:"required"`
			Items  []domain.CartItemInput `json:"items" binding:"required,min=1"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		if err := services.Sales.Refund(c.Request.Context(), user, input.SaleID, input.Items); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		ok(c, "Refund completed", nil)
	})
}

func canAccessBranch(ctx context.Context, services *usecase.Services, user domain.User, branchID uuid.UUID) bool {
	branches, err := services.Dashboard.AccessibleBranches(ctx, user)
	if err != nil {
		return false
	}
	for _, branch := range branches {
		if branch.ID == branchID {
			return true
		}
	}
	return false
}

func parseDateRange(from, to string) (*time.Time, *time.Time, error) {
	var dateFrom *time.Time
	var dateTo *time.Time
	if from != "" {
		parsed, err := time.Parse("2006-01-02", from)
		if err != nil {
			return nil, nil, err
		}
		dateFrom = &parsed
	}
	if to != "" {
		parsed, err := time.Parse("2006-01-02", to)
		if err != nil {
			return nil, nil, err
		}
		parsed = parsed.Add(24 * time.Hour)
		dateTo = &parsed
	}
	return dateFrom, dateTo, nil
}

func registerUsers(r *gin.RouterGroup, services *usecase.Services) {
	users := r.Group("/users", middleware.OwnerOrManager())
	users.GET("", func(c *gin.Context) {
		user := c.MustGet(middleware.CurrentUserKey).(domain.User)
		data, err := services.Users.List(c.Request.Context(), user)
		if err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}
		ok(c, "Success", data)
	})
	users.GET("/sales-summary", func(c *gin.Context) {
		user := c.MustGet(middleware.CurrentUserKey).(domain.User)
		data, err := services.Users.SalesSummary(c.Request.Context(), user)
		if err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}
		ok(c, "Success", data)
	})
	users.POST("", func(c *gin.Context) {
		actor := c.MustGet(middleware.CurrentUserKey).(domain.User)
		var input struct {
			domain.User
			Password string `json:"password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		user, err := services.Users.Create(c.Request.Context(), actor, input.User, input.Password)
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		created(c, "User created", user)
	})
	users.PUT("/:id", func(c *gin.Context) {
		actor := c.MustGet(middleware.CurrentUserKey).(domain.User)
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		var input domain.User
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		user, err := services.Users.Update(c.Request.Context(), actor, id, input)
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		ok(c, "User updated", user)
	})
}

func registerBranches(r *gin.RouterGroup, services *usecase.Services) {
	r.GET("/branches/my", func(c *gin.Context) {
		user := c.MustGet(middleware.CurrentUserKey).(domain.User)
		branches, err := services.Dashboard.AccessibleBranches(c.Request.Context(), user)
		if err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}
		ok(c, "Success", branches)
	})
	owner := r.Group("/branches", middleware.OwnerOnly())
	owner.POST("", func(c *gin.Context) {
		var input domain.Branch
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		branch, err := services.System.CreateBranch(c.Request.Context(), input)
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		created(c, "Branch created", branch)
	})
	owner.PUT("/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		var input domain.Branch
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		branch, err := services.System.UpdateBranch(c.Request.Context(), id, input)
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		ok(c, "Branch updated", branch)
	})
}

func registerCategories(r *gin.RouterGroup, services *usecase.Services) {
	r.GET("/categories", func(c *gin.Context) {
		categories, err := services.System.ListCategories(c.Request.Context())
		if err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}
		ok(c, "Success", categories)
	})
	owner := r.Group("/categories", middleware.OwnerOnly())
	owner.POST("", func(c *gin.Context) {
		var input domain.Category
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		category, err := services.System.CreateCategory(c.Request.Context(), input)
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		created(c, "Category created", category)
	})
	owner.PUT("/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		var input domain.Category
		if err := c.ShouldBindJSON(&input); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		category, err := services.System.UpdateCategory(c.Request.Context(), id, input)
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		ok(c, "Category updated", category)
	})
	owner.DELETE("/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		if err := services.System.DeleteCategory(c.Request.Context(), id); err != nil {
			fail(c, http.StatusBadRequest, err)
			return
		}
		ok(c, "Category deleted", nil)
	})
}

func registerDashboard(r *gin.RouterGroup, services *usecase.Services) {
	r.GET("/dashboard/summary", middleware.OwnerOrManager(), func(c *gin.Context) {
		user := c.MustGet(middleware.CurrentUserKey).(domain.User)
		var branchID *uuid.UUID
		if c.Query("branch_id") != "" {
			id, err := uuid.Parse(c.Query("branch_id"))
			if err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}
			branchID = &id
		}
		summary, err := services.Dashboard.Summary(c.Request.Context(), user, branchID)
		if err != nil {
			fail(c, http.StatusForbidden, err)
			return
		}
		ok(c, "Success", summary)
	})
}

func registerAuditLogs(r *gin.RouterGroup, services *usecase.Services) {
	r.GET("/audit-logs", middleware.OwnerOnly(), func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "150"))
		logs, err := services.Audit.List(c.Request.Context(), c.Query("action"), c.Query("entity_type"), c.Query("q"), limit)
		if err != nil {
			fail(c, http.StatusInternalServerError, err)
			return
		}
		ok(c, "Success", logs)
	})
}
