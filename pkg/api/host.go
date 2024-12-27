package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koodeyo/koodnet/pkg/database"
	"github.com/koodeyo/koodnet/pkg/models"
	"gorm.io/gorm/clause"
)

// FindHosts godoc
// @Summary Get all hosts
// @Description Get a list of all hosts with optional pagination
// @Tags hosts
// @Produce json
// @Param page query int false "page for pagination" default(1)
// @Param page_size query int false "page_size for pagination" default(10)
// @Success 200 {object} api.paginatedResponse[models.Host]
// @Router /hosts [get]
func FindHosts(c *gin.Context) {
	var hosts []models.Host

	// Fetch data from the database with pagination
	database.Conn.Model(&models.Host{}).Scopes(models.Paginate(c)).
		Preload("Network").
		Preload("Certificates").
		Find(&hosts)

	response := paginated(hosts, c)

	c.JSON(http.StatusOK, response)
}

// CreateHost godoc
// @Summary Create a new host
// @Description Create a host with the provided details
// @Tags hosts
// @Accept json
// @Produce json
// @Param host body models.HostDto true "Host Payload"
// @Success 201 {object} models.Host
// @Failure 400 {object} api.errorResponse
// @Router /hosts [post]
func CreateHost(c *gin.Context) {
	var dto models.HostDto

	// Validate the request body
	if err := c.ShouldBindJSON(&dto); err != nil {
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

	// Create new host instance
	host := models.Host{
		Name:            dto.Name,
		IP:              dto.IP,
		StaticAddresses: dto.StaticAddresses,
		Subnets:         dto.Subnets,
		Groups:          dto.Groups,
		ListenPort:      dto.ListenPort,
		IsLighthouse:    dto.IsLighthouse,
		NetworkID:       dto.NetworkID,
	}

	// Save to database
	if err := database.Conn.Create(&host).Error; err != nil {
		dbErrorHandler(err, c)
		return
	}

	// Reload the host with associations
	database.Conn.Preload("Network").Preload("Certificates").First(&host, "id = ?", host.ID)

	c.JSON(http.StatusCreated, host)
}

// DeleteHost godoc
// @Summary Delete a host
// @Description Delete a host by ID
// @Tags hosts
// @Param id path string true "Host ID"
// @Success 200 {object} map[string]bool "Delete status"
// @Failure 404 {object} api.errorResponse
// @Router /hosts/{id} [delete]
func DeleteHost(c *gin.Context) {
	id := c.Param("id")

	if err := database.Conn.Delete(&models.Host{}, "id = ?", id).Error; err != nil {
		dbErrorHandler(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"delete": true})
}

// FindHost godoc
// @Summary Get a host by ID
// @Description Retrieve details of a single host
// @Tags hosts
// @Param id path string true "Host ID"
// @Produce json
// @Success 200 {object} models.Host
// @Failure 404 {object} api.errorResponse
// @Router /hosts/{id} [get]
func FindHost(c *gin.Context) {
	id := c.Param("id")
	var host models.Host

	if err := database.Conn.Preload("Network").Preload("Cert").First(&host, "id = ?", id).Error; err != nil {
		dbErrorHandler(err, c)
		return
	}

	c.JSON(http.StatusOK, host)
}

// UpdateHost godoc
// @Summary Update a host
// @Description Update the details of an existing host
// @Tags hosts
// @Accept json
// @Produce json
// @Param id path string true "Host ID"
// @Param host body models.HostDto true "Updated host details"
// @Success 200 {object} models.Host
// @Failure 400 {object} api.errorResponse
// @Failure 404 {object} api.errorResponse
// @Router /hosts/{id} [patch]
func UpdateHost(c *gin.Context) {
	id := c.Param("id")
	var host models.Host

	// Find existing host
	if err := database.Conn.Omit(clause.Associations).First(&host, "id = ?", id).Error; err != nil {
		dbErrorHandler(err, c)
		return
	}

	// Bind the update payload
	var dto models.HostDto
	if err := c.ShouldBindJSON(&dto); err != nil {
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

	// Update the host
	if err := database.Conn.Model(&host).Updates(dto).Error; err != nil {
		dbErrorHandler(err, c)
		return
	}

	// Reload the host with associations
	database.Conn.Preload("Network").Preload("Certificates").First(&host, "id = ?", id)

	c.JSON(http.StatusOK, host)
}
