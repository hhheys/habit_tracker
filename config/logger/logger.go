package logger

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewLogger creates a new logger.
func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

// WithGin adds gin zap middleware to the gin engine.
func WithGin(logger *zap.Logger, r *gin.Engine) {
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true)) // access‑логи [web:6][web:10]
	r.Use(ginzap.RecoveryWithZap(logger, true))      // panic → error с stacktrace [web:6][web:10]
}
