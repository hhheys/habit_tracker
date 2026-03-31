package router

import (
	"habit-tracker/internal/handler"

	"github.com/gin-gonic/gin"
)

// NewUserRouter creates a new router for user endpoints.
func NewUserRouter(r *gin.Engine, h handler.Handler) {
	api := r.Group("/api")
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)
}
