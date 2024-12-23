package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// End timer
		duration := time.Since(start)

		// Create log fields
		logFields := logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration,
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
		}

		// Include errors if they exist
		if len(c.Errors) > 0 {
			logFields["errors"] = c.Errors.ByType(gin.ErrorTypePrivate).String()
		}

		// Log the request details using logrus
		logger.WithFields(logFields).Info("Request")
	}
}
