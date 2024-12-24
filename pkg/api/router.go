package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/koodeyo/koodnet/docs"
	"github.com/koodeyo/koodnet/pkg/middleware"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

func NewRouter(l *logrus.Logger) *gin.Engine {
	r := gin.Default()
	// Use centralized error handling middleware
	r.Use(errorHandler)

	// Apply the logging middleware
	r.Use(middleware.Logger(l))

	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.Security())
		r.Use(middleware.Xss())
	}

	r.Use(middleware.Cors())
	r.Use(middleware.RateLimiter(rate.Every(1*time.Minute), 60)) // 60 requests per minute

	docs.SwaggerInfo.BasePath = "/api/v1"

	v1 := r.Group("/api/v1")
	{
		v1.GET("/", healthCheckHandler)
		v1.GET("/networks", FindNetworks)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Handle 404 (unmatched routes)
	r.NoRoute(notFoundHandler)

	return r
}
