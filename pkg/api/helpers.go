package api

import (
	"github.com/gin-gonic/gin"
	"github.com/koodeyo/koodnet/pkg/models"
)

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Errors []apiError `json:"errors"`
}

// ResponseMetadata holds metadata about the response, such as the total count of items.
type metadata struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"` // Total represents the total number of items.
	Total      int `json:"total"`       // Total represents the total number of items.
}

type ResourceModel interface {
	models.Network
}

// CollectionResponse is a generic structure for a paginated response.
// It contains the data and metadata associated with the response.
type collectionResponse[T ResourceModel] struct {
	Data     []T      `json:"data"`     // Data contains the actual collection of items.
	Metadata metadata `json:"metadata"` // Metadata contains additional info like the total count.
}

type repository interface {
	*networkRepository
}

func contextMiddleware[T repository](repository T) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("appCtx", repository)
		c.Next()
	}
}

func paginated[T ResourceModel](data []T, c *gin.Context) collectionResponse[T] {
	return collectionResponse[T]{
		Data: data,
		Metadata: metadata{
			Page:       c.GetInt("page"),
			Total:      c.GetInt("total"),
			PageSize:   c.GetInt("pageSize"),
			TotalPages: c.GetInt("totalPages"),
		},
	}
}

// Middleware for centralized error handling
// func errorHandler(c *gin.Context) {
// 	c.Next()

// 	if len(c.Errors) > 0 {
// 		var errors []apiError

// 		for _, err := range c.Errors {
// 			errors = append(errors, apiError{
// 				Code:    "ERR_UNKNOWN", // Default error code
// 				Message: err.Error(),   // Use the error message from Gin
// 			})
// 		}

// 		// Respond with a JSON error response
// 		c.JSON(http.StatusBadRequest, errorResponse{Errors: errors})
// 	}
// }
