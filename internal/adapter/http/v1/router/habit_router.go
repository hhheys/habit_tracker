package router

import (
	"habit-tracker/internal/adapter/auth"
	"habit-tracker/internal/adapter/http/middleware"
	"habit-tracker/internal/adapter/http/v1/handler"
	"habit-tracker/internal/domain"

	"github.com/gin-gonic/gin"
)

func NewHabitRouter(r *gin.Engine, h handler.Handler, jwtService *auth.JwtService) {
	g := r.Group("/api/habit")

	g.Use(middleware.JWTAuthMiddleware(jwtService))

	g.GET("", h.GetAllHabits)
	g.GET("/:id", h.GetHabitByID)
	g.POST("", middleware.RoleAuthMiddleware(domain.UserRoleAdmin), h.CreateHabit)
	g.PUT("/:id", middleware.RoleAuthMiddleware(domain.UserRoleAdmin), h.UpdateHabit)
	g.DELETE("/:id", middleware.RoleAuthMiddleware(domain.UserRoleAdmin), h.DeleteHabit)

	g.POST("/:id/add", middleware.RoleAuthMiddleware(domain.UserRoleDefault), h.AddUserHabit)

	g.GET("/tag/all", h.GetAllTags)
	g.GET("/tag/:id", h.GetTagByID)
	g.POST("/tag", middleware.RoleAuthMiddleware(domain.UserRoleAdmin), h.CreateTag)
	g.PUT("/tag/:id", middleware.RoleAuthMiddleware(domain.UserRoleAdmin), h.EditTag)
	g.DELETE("/tag/:id", middleware.RoleAuthMiddleware(domain.UserRoleAdmin), h.DeleteTag)

	g.GET("/my", h.GetAllUserHabits)
	g.POST("/confirm/:id", h.CreateDailyConfirmation)

	g.GET("/heatmap", h.GetHeatMap)
}
