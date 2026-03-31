package middleware

import (
	"habit-tracker/internal/auth"
	"habit-tracker/internal/dto/response"
	appErrors "habit-tracker/internal/errors"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware validates JWT token
func JWTAuthMiddleware(jwtService auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)

		if tokenString == "" {
			c.AbortWithStatusJSON(401, response.NewErrorResponse("INVALID_REQUEST", "invalid request"))
			return
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, response.NewErrorResponse("INVALID_REQUEST", "invalid request"))
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RoleAuthMiddleware checks if the user has the required role
func RoleAuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			_ = c.Error(appErrors.ErrInvalidRole)
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			_ = c.Error(appErrors.ErrInvalidRole)
			return
		}

		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				c.Abort()
				return
			}
		}

		_ = c.Error(appErrors.ErrNoPermissions)
		c.Abort()
	}
}

func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")

	if bearerToken == "" {
		return ""
	}

	// "Bearer <token>"
	parts := strings.Split(bearerToken, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}

	return ""
}
