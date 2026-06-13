package router

import (
	"pos-system/backend/internal/app/domain"
	"pos-system/backend/internal/app/port"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(svc port.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := svc.Authenticate(c.Request.Context(), c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(statusFromError(err), APIError{Code: "ERROR", Description: err.Error()})
			return
		}
		ctx := domain.ContextWithUser(c.Request.Context(), *user)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
