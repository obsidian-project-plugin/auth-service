package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

func LoggingMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Writer.Header().Set("X-Request-ID", reqID)

		entry := logger.With("request_id", reqID)
		c.Set("logger", entry)

		c.Next()

		latency := time.Since(start)

		logArgs := []interface{}{
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
			"latency", latency,
			"user_agent", c.Request.UserAgent(),
		}
		if raw := c.Request.URL.RawQuery; raw != "" {
			logArgs = append(logArgs, "query", raw)
		}

		if len(c.Errors) > 0 {
			logArgs = append(logArgs, "errors", c.Errors.String())
			entry.Errorw("request completed with errors", logArgs...)
		} else {
			entry.Infow("request completed successfully", logArgs...)
		}
	}
}
