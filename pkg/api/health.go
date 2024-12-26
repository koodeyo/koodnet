package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// @BasePath /api/v1

// Healthcheck godoc
// @Summary Health check for the service
// @Schemes
// @Description This endpoint is used to verify the health and availability of the service.
// It can be used by monitoring tools or external systems to ensure that the service is running and responsive.
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} api.healthResponse "The service is operational and healthy"
// @Router / [get]
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, healthResponse{
		Status:  "ok",
		Message: "Service is healthy",
	})
}
