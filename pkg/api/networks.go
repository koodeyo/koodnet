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
	database.Conn.Model(&models.Network{}).Scopes(models.Paginate(c)).Preload("Ca").Find(&networks)

	response := paginated(networks, c)

	// Return the response using the response struct
	c.JSON(http.StatusOK, response)
}

// CreateNetwork godoc
// @Summary Create a new network
// @Description Create a network with the provided details
// @Tags networks
// @Accept json
// @Produce json
// @Param network body models.NetworkDto true "Network Payload"
// @Success 201 {object} models.Network
// @Failure 400 {object} api.errorResponse
// @Router /networks [post]
func CreateNetwork(c *gin.Context) {
	var p models.NetworkDto

	// Validate the p
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Errors: []apiError{
				{
					Code:    "INVALID_DATA",
					Message: err.Error(),
				},
			},
		})
		return
	}

	n := models.Network{
		ID:               uuid.New().String(), // Unique identifier for the network
		Name:             p.Name,              // Name of the network
		IPs:              p.IPs,               // List of IP ranges
		Subnets:          p.Subnets,           // List of subnets
		Groups:           p.Groups,            // Associated groups
		Duration:         p.Duration,          // Duration in seconds
		Encrypt:          p.Encrypt,           // Whether encryption is enabled
		Passphrase:       p.Passphrase,        // Encryption passphrase
		ArgonMemory:      p.ArgonMemory,       // Memory usage for Argon2
		ArgonIterations:  p.ArgonIterations,   // Iterations for Argon2
		ArgonParallelism: p.ArgonParallelism,  // Parallelism for Argon2
		Curve:            p.Curve,             // Cryptographic curve (e.g., 25519)
	}

	// Save to the database
	if err := database.Conn.Create(&n).Error; err != nil {
		dbErrorHandler(err, c)
		return
	}

	// Respond with the created network
	c.JSON(http.StatusCreated, n)
}

// DeleteNetwork godoc
// @Summary Delete a network
// @Description Delete a network by ID
// @Tags networks
// @Param id path string true "Network ID"
// @Success 204 {string} string "Network deleted"
// @Failure 404 {object} api.errorResponse
// @Router /networks/{id} [delete]
func DeleteNetwork(c *gin.Context) {
	id := c.Param("id")

	// Attempt to delete the network
	if err := database.Conn.Delete(&models.Network{}, "id = ?", id).Error; err != nil {
		dbErrorHandler(err, c)
		return
	}

	// Respond with no content
	c.Status(http.StatusNoContent)
}

// FindNetwork godoc
// @Summary Get a network by ID
// @Description Retrieve details of a single network
// @Tags networks
// @Param id path string true "Network ID"
// @Produce json
// @Success 200 {object} models.Network
// @Failure 404 {object} api.errorResponse
// @Router /networks/{id} [get]
func FindNetwork(c *gin.Context) {
	id := c.Param("id")
	var network models.Network

	// Attempt to find the network
	if err := database.Conn.Preload("Ca").First(&network, "id = ?", id).Error; err != nil {
		dbErrorHandler(err, c)
		return
	}

	// Respond with the found network
	c.JSON(http.StatusOK, network)
}

// UpdateNetwork godoc
// @Summary Update a network
// @Description Update the details of an existing network
// @Tags networks
// @Accept json
// @Produce json
// @Param id path string true "Network ID"
// @Param network body models.NetworkDto true "Updated network details"
// @Success 200 {object} models.Network
// @Failure 400 {object} api.errorResponse
// @Failure 404 {object} api.errorResponse
// @Router /networks/{id} [patch]
func UpdateNetwork(c *gin.Context) {
	id := c.Param("id")
	var n models.Network

	// Attempt to find the existing network
	if err := database.Conn.First(&n, "id = ?", id).Error; err != nil {
		dbErrorHandler(err, c)
		return
	}

	// Bind the payload JSON to a new network struct
	var u models.NetworkDto
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Errors: []apiError{
				{
					Code:    "INVALID_INPUT",
					Message: err.Error(),
				},
			},
		})
		return
	}

	// Save the updated network to the database
	if err := database.Conn.Model(&n).Updates(u).Error; err != nil {
		dbErrorHandler(err, c)
		return
	}

	// Respond with the updated network
	c.JSON(http.StatusOK, n)
}
