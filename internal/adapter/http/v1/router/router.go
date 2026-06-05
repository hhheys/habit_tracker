package router

import (
	authadapter "habit-tracker/internal/adapter/auth"
	"habit-tracker/internal/adapter/http/middleware"
	"habit-tracker/internal/adapter/http/v1/handler"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(h handler.Handler, jwt *authadapter.JwtService, log *zap.Logger) *gin.Engine {
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
		AllowAllOrigins: true,
	}))
	r.Static("/images", "./images")
	r.Use(gin.Logger(), gin.Recovery(), middleware.ErrorHandler(log))
	r.GET("/_info", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) })
	NewUserRouter(r, h)
	NewHabitRouter(r, h, jwt)
	NewAchievementRouter(r, h, jwt)
	return r
}
