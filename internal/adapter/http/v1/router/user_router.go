package router

import (
	"habit-tracker/internal/adapter/http/v1/handler"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(r *gin.Engine, h handler.Handler) {
	api := r.Group("/api")
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)
	api.POST("/refresh", h.Refresh)
}
