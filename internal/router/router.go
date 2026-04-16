package router

import (
	"habit-tracker/internal/auth"
	"habit-tracker/internal/handler"
	"habit-tracker/internal/middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewRouter creates a new router for the application.
func NewRouter(h handler.Handler, jwtService auth.JWTService, log *zap.Logger) *gin.Engine {
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowAllOrigins:  true,
		AllowCredentials: false,
	}))

	r.Static("../images", "./images")

	r.Use(gin.Logger())
	r.Use(middleware.ErrorHandler(log))

	r.GET("/_info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	NewUserRouter(r, h)
	NewHabitRouter(r, h, jwtService)

	return r
}
