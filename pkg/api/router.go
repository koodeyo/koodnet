package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koodeyo/koodnet/pkg/middleware"
	"github.com/sirupsen/logrus"
)

func NewRouter(l *logrus.Logger, ctx *context.Context) *gin.Engine {
	r := gin.Default()

	// Apply the logging middleware
	r.Use(middleware.Logger(l))

	// Health check endpoint to verify if the service is running
	// Responds with a 200 OK status and a message indicating the service is healthy
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Service is healthy",
		})
	})

	return r
}
