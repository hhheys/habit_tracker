package middleware

import (
	"errors"
	"habit-tracker/internal/dto/response"
	appErrors "habit-tracker/internal/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorInfo stores error information
type ErrorInfo struct {
	Code    string
	Status  int
	Message string
}

var errorMap = map[error]ErrorInfo{
	appErrors.ErrUserNotFound:      {"NOT_FOUND", http.StatusNotFound, ""},
	appErrors.ErrUserAlreadyExists: {"INVALID_REQUEST", http.StatusBadRequest, ""},
	appErrors.ErrWrongPassword:     {"INVALID_REQUEST", http.StatusUnauthorized, ""},
}

// ErrorHandler returns a middleware that handles errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		var info ErrorInfo

		err := c.Errors.Last().Err
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			info = ErrorInfo{
				Code:    "INVALID_REQUEST",
				Status:  http.StatusBadRequest,
				Message: "invalid request",
			}
		} else {
			var ok bool
			info, ok = errorMap[err]
			if !ok {
				info = ErrorInfo{
					Code:    "INTERNAL_ERROR",
					Status:  http.StatusInternalServerError,
					Message: "internal server error",
				}
			}
		}

		if info.Message == "" {
			info.Message = err.Error()
		}

		c.JSON(info.Status, response.NewErrorResponse(info.Code, info.Message))
	}
}
