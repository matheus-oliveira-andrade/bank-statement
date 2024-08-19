package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func DefaultStructuredLogger() gin.HandlerFunc {
	return StructuredLogger()
}

func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method

		c.Next()

		timeStamp := time.Now()
		latency := timeStamp.Sub(start)
		if latency > time.Minute {
			latency = latency.Truncate(time.Second)
		}

		statusCode := c.Writer.Status()
		if raw != "" {
			path = path + "?" + raw
		}

		slog.Info(fmt.Sprintf("called %v %v", method, path),
			"method", method,
			"responseStatusCode", statusCode,
			"path", path,
			"latency", latency.String())
	}
}
