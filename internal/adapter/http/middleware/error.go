package middleware

import (
	"errors"
	"habit-tracker/internal/adapter/http/v1/dto/response"
	"habit-tracker/internal/domain"
	authuc "habit-tracker/internal/usecase/auth"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type ErrorInfo struct {
	Code    string
	Status  int
	Message string
}

var errorMap = map[error]ErrorInfo{
	domain.ErrUserNotFound:          {"NOT_FOUND", http.StatusNotFound, ""},
	domain.ErrUserAlreadyExists:     {"INVALID_REQUEST", http.StatusConflict, ""},
	domain.ErrHabitNotFound:         {"NOT_FOUND", http.StatusNotFound, ""},
	domain.ErrHabitAlreadyExists:    {"INVALID_REQUEST", http.StatusConflict, ""},
	domain.ErrHabitAlreadyAdded:     {"INVALID_REQUEST", http.StatusConflict, ""},
	domain.ErrHabitNotAdded:         {"INVALID_REQUEST", http.StatusConflict, ""},
	domain.ErrHabitAlreadyConfirmed: {"INVALID_REQUEST", http.StatusConflict, ""},
	domain.ErrUnauthorized:          {"UNAUTHORIZED", http.StatusUnauthorized, ""},
	domain.ErrInvalidRole:           {"FORBIDDEN", http.StatusForbidden, ""},
	domain.ErrNoPermissions:         {"FORBIDDEN", http.StatusForbidden, ""},
	authuc.ErrInvalidCredentials:    {"UNAUTHORIZED", http.StatusUnauthorized, ""},
	domain.ErrSessionNotFound:       {"UNAUTHORIZED", http.StatusUnauthorized, ""},
	domain.ErrSessionRevoked:        {"UNAUTHORIZED", http.StatusUnauthorized, ""},
	domain.ErrTokenExpired:          {"UNAUTHORIZED", http.StatusUnauthorized, ""},
}

func ErrorHandler(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}
		err := c.Errors.Last().Err
		log.Error("request failed", zap.Error(err))

		info, ok := errorMap[err]
		var validationErrors validator.ValidationErrors
		var numberError *strconv.NumError
		if errors.As(err, &validationErrors) || errors.As(err, &numberError) {
			info = ErrorInfo{"INVALID_REQUEST", http.StatusBadRequest, "invalid request"}
			ok = true
		}
		if !ok {
			info = ErrorInfo{"INTERNAL_ERROR", http.StatusInternalServerError, "internal server error"}
		}
		if info.Message == "" {
			info.Message = err.Error()
		}
		c.JSON(info.Status, response.NewErrorResponse(info.Code, info.Message))
	}
}
