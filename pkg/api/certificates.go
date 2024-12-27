package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koodeyo/koodnet/pkg/database"
	"github.com/koodeyo/koodnet/pkg/models"
)

// FindCertificates godoc
// @Summary Get all certificates
// @Description Get a list of all certificates with optional pagination
// @Tags certificates
// @Produce json
// @Param page query int false "page for pagination" default(1)
// @Param pageSize query int false "pageSize for pagination" default(10)
// @Success 200 {object} api.paginatedResponse[models.Certificate]
// @Router /certificates [get]
func FindCertificates(c *gin.Context) {
	var certificates []models.Certificate

	// Fetch data from the database
	database.Conn.Model(&models.Certificate{}).Scopes(models.Paginate(c)).Find(&certificates)

	response := paginated(certificates, c)

	// Return the response using the response struct
	c.JSON(http.StatusOK, response)
}
