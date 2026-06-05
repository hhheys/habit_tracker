package handler

import (
	"habit-tracker/internal/adapter/http/v1/dto/request"
	"habit-tracker/internal/adapter/http/v1/dto/response"
	authuc "habit-tracker/internal/usecase/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Refresh(c *gin.Context)
}

type authHandler struct {
	service *authuc.Service
	log     *zap.Logger
}

func NewAuthHandler(service *authuc.Service, log *zap.Logger) AuthHandler {
	return &authHandler{service: service, log: log}
}

func (h *authHandler) Register(c *gin.Context) {
	var req request.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	output, err := h.service.Register(c.Request.Context(), &authuc.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Timezone: req.Timezone,
	}, sessionInfo(c))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, response.NewUserRegisterResponse(output))
}

func (h *authHandler) Login(c *gin.Context) {
	var req request.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	if req.Username == "" && req.Email == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse("INVALID_REQUEST", "username or email is required"))
		return
	}
	output, err := h.service.Login(c.Request.Context(), &authuc.LoginInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}, sessionInfo(c))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.NewUserLoginResponse(output))
}

func (h *authHandler) Refresh(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}
	output, err := h.service.RefreshToken(c.Request.Context(), &authuc.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
	}, sessionInfo(c))
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.NewAuthResponse(*output))
}

func sessionInfo(c *gin.Context) authuc.SessionInfoInput {
	return authuc.SessionInfoInput{UserIP: c.ClientIP(), UserAgent: c.Request.UserAgent()}
}
