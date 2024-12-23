package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koodeyo/koodnet/pkg/models"
	"gorm.io/gorm"
)

type NetworkRepository interface {
	FindNetworks(c *gin.Context)
	CreateNetwork(c *gin.Context)
	FindNetwork(c *gin.Context)
	UpdateNetwork(c *gin.Context)
	DeleteNetwork(c *gin.Context)
}

// networkRepository holds shared resources like database and context
type networkRepository struct {
	db  *gorm.DB
	ctx *context.Context
}

// NewAppContext creates a new AppContext
func newNetworkRepository(db *gorm.DB, ctx *context.Context) *networkRepository {
	return &networkRepository{
		db:  db,
		ctx: ctx,
	}
}

// FindNetworks godoc
// @Summary Get all networks
// @Description Get a list of all networks with optional pagination
// @Tags networks
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {object} api.collectionResponse[models.Network]
// @Router /networks [get]
func (r *networkRepository) FindNetworks(c *gin.Context) {
	var networks []models.Network

	// Fetch data from the database
	r.db.Scopes(models.Paginate(c)).Find(&networks)

	response := paginated(networks, c)

	// Return the response using the response struct
	c.JSON(http.StatusOK, response)
}

// CreateNetwork implements NetworkRepository.
func (r *networkRepository) CreateNetwork(c *gin.Context) {
	panic("unimplemented")
}

// DeleteNetwork implements NetworkRepository.
func (r *networkRepository) DeleteNetwork(c *gin.Context) {
	panic("unimplemented")
}

// FindNetwork implements NetworkRepository.
func (r *networkRepository) FindNetwork(c *gin.Context) {
	panic("unimplemented")
}

// UpdateNetwork implements NetworkRepository.
func (r *networkRepository) UpdateNetwork(c *gin.Context) {
	panic("unimplemented")
}
