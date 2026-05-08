package router

import (
	"habit-tracker/internal/auth"
	"habit-tracker/internal/handler"
	"habit-tracker/internal/middleware"
	"habit-tracker/internal/models"

	"github.com/gin-gonic/gin"
)

func NewHabitRouter(r *gin.Engine, h handler.Handler, jwtService auth.JWTService) {
	g := r.Group("/api/habit")

	g.Use(middleware.JWTAuthMiddleware(jwtService))

	g.GET("", h.GetAllHabits)
	g.POST("", middleware.RoleAuthMiddleware(string(models.UserRoleAdmin)), h.CreateHabit)
	g.PUT("/:id", middleware.RoleAuthMiddleware(string(models.UserRoleAdmin)), h.UpdateHabit)
	g.DELETE("/:id", middleware.RoleAuthMiddleware(string(models.UserRoleAdmin)), h.DeleteHabit)

	g.POST("/:id/add", middleware.RoleAuthMiddleware(string(models.UserRoleDefault)), h.AddUserHabit)

	g.GET("/tag/all", h.GetAllTags)
	g.POST("/tag", middleware.RoleAuthMiddleware(string(models.UserRoleAdmin)), h.CreateTag)
	g.PUT("/tag/:id", middleware.RoleAuthMiddleware(string(models.UserRoleAdmin)), h.EditTag)
	g.DELETE("/tag/:id", middleware.RoleAuthMiddleware(string(models.UserRoleAdmin)), h.DeleteTag)

	g.GET("/my", h.GetAllUserHabits)
	g.POST("/confirm/:id", h.CreateDailyConfirmation)

	g.GET("/heatmap", h.GetHeatMap)
}
