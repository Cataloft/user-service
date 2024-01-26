package middlewares

import (
	"github.com/gin-gonic/gin"
	requestid "github.com/sumit-tembe/gin-requestid"
	"log/slog"
	"test-task/internal/utils"
	"time"
)

func LogMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := utils.GetDurationInMilliseconds(start)
		entry := log.With(
			"client_ip", c.ClientIP(),
			"duration", duration,
			"method", c.Request.Method,
			"path", c.Request.RequestURI,
			"status", c.Writer.Status(),
			"request_id", requestid.GetRequestIDFromHeaders(c),
		)

		if c.Writer.Status() >= 500 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info("success")
		}
	}
}
