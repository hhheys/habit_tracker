package router

import (
	"habit-tracker/internal/auth"
	"habit-tracker/internal/handler"
	"habit-tracker/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewHabitRouter(r *gin.Engine, h handler.Handler, jwtService auth.JWTService) {
	g := r.Group("/api/habit")

	g.Use(middleware.JWTAuthMiddleware(jwtService))

	g.GET("", h.GetAllHabits)
}
