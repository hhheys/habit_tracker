package router

import (
	"habit-tracker/internal/adapter/auth"
	"habit-tracker/internal/adapter/http/middleware"
	"habit-tracker/internal/adapter/http/v1/handler"

	"github.com/gin-gonic/gin"
)

func NewAchievementRouter(r *gin.Engine, h handler.Handler, jwtService *auth.JwtService) {
	g := r.Group("")
	g.Use(middleware.JWTAuthMiddleware(jwtService))

	g.GET("/achievements", h.GetAchievements)
}
