package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koodeyo/koodnet/pkg/database"
	"github.com/koodeyo/koodnet/pkg/models"
)

// FindNetworks godoc
// @Summary Get all networks
// @Description Get a list of all networks with optional pagination
// @Tags networks
// @Produce json
// @Param page query int false "page for pagination" default(1)
// @Param page_size query int false "page_size for pagination" default(10)
// @Success 200 {object} api.paginatedResponse[models.Network]
// @Router /networks [get]
func FindNetworks(c *gin.Context) {
	var networks []models.Network

	// Fetch data from the database
	database.Conn.Scopes(models.Paginate(c)).Find(&networks)

	response := paginated(networks, c)

	// Return the response using the response struct
	c.JSON(http.StatusOK, response)
}

func CreateNetwork(c *gin.Context) {
	network := models.Network{
		ID: uuid.New().String(),
	}

	database.Conn.Create(&network)

	panic("unimplemented")
}

func DeleteNetwork(c *gin.Context) {
	panic("unimplemented")
}

func FindNetwork(c *gin.Context) {
	panic("unimplemented")
}

func UpdateNetwork(c *gin.Context) {
	panic("unimplemented")
}
