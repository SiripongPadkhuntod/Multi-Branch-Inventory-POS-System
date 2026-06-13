package middleware

import (
	"net/http"
	"strings"

	"pos-system/backend/internal/app/domain"
	"pos-system/backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

const CurrentUserKey = "current_user"

func Auth(auth *usecase.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := bearer(c.GetHeader("Authorization"))
		if token == "" {
			if cookie, err := c.Cookie("access_token"); err == nil {
				token = cookie
			}
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "missing token"})
			return
		}

		claims, err := auth.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid token"})
			return
		}
		if claims.TokenUse != "" && claims.TokenUse != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid token"})
			return
		}
		user, err := auth.FindUser(c.Request.Context(), claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "user not found"})
			return
		}
		c.Set(CurrentUserKey, *user)
		c.Next()
	}
}

func OwnerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet(CurrentUserKey).(domain.User)
		if user.Role != domain.RoleOwner {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "message": "owner permission required"})
			return
		}
		c.Next()
	}
}

func OwnerOrManager() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet(CurrentUserKey).(domain.User)
		if user.Role != domain.RoleOwner && user.Role != domain.RoleManager {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "message": "manager permission required"})
			return
		}
		c.Next()
	}
}

func bearer(header string) string {
	if strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return strings.TrimSpace(header[7:])
	}
	return ""
}
