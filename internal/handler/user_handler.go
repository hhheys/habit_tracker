package handler

import (
	"habit-tracker/internal/dto/request"
	"habit-tracker/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type userHandler struct {
	log     *zap.Logger
	service service.Service
}

// NewUserHandler returns a new instance of UserHandler.
func NewUserHandler(log *zap.Logger, service service.Service) UserHandler {
	return &userHandler{
		log:     log,
		service: service,
	}
}

func (h *userHandler) Register(c *gin.Context) {
	var req request.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Bind json error", zap.Error(err))
		_ = c.Error(err)
		return
	}

	resp, err := h.service.RegisterUser(&req)
	if err != nil {
		h.log.Error("User register error", zap.Error(err))
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *userHandler) Login(c *gin.Context) {
	var req request.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Bind json error", zap.Error(err))
		_ = c.Error(err)
		return
	}

	resp, err := h.service.AuthUser(req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
