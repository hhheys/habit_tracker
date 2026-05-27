package middleware

import (
	authadapter "habit-tracker/internal/adapter/auth"
	"habit-tracker/internal/adapter/http/v1/dto/response"
	"habit-tracker/internal/domain"
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenValidator interface {
	ValidateToken(token string) (*authadapter.Claims, error)
}

func JWTAuthMiddleware(tokens TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(401, response.NewErrorResponse("UNAUTHORIZED", "authorization token is required"))
			return
		}
		claims, err := tokens.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, response.NewErrorResponse("UNAUTHORIZED", "invalid token"))
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("role", domain.UserRole(claims.Role))
		c.Next()
	}
}

func RoleAuthMiddleware(allowedRoles ...domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("role")
		role, ok := value.(domain.UserRole)
		if !exists || !ok {
			_ = c.Error(domain.ErrInvalidRole)
			c.Abort()
			return
		}
		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				return
			}
		}
		_ = c.Error(domain.ErrNoPermissions)
		c.Abort()
	}
}

func extractToken(c *gin.Context) string {
	parts := strings.Fields(c.GetHeader("Authorization"))
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}
