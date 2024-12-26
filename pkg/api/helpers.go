package api

import (
	"net/http"

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

type metadata struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
	Total      int `json:"total"` // Total represents the total number of items.
}

// for a paginated response.
// It contains the data and metadata associated with the response.
type paginatedResponse[T interface{}] struct {
	Data     []T      `json:"data"`     // Data contains the actual collection of items.
	Metadata metadata `json:"metadata"` // Metadata contains additional info like the total count.
}

func paginated[T interface{}](data []T, c *gin.Context) paginatedResponse[T] {
	return paginatedResponse[T]{
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
func errorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		var errors []apiError

		for _, err := range c.Errors {
			errors = append(errors, apiError{
				Code:    "ERR_UNKNOWN", // Default error code
				Message: err.Error(),   // Use the error message from Gin
			})
		}

		// Respond with a JSON error response
		c.JSON(http.StatusBadRequest, errorResponse{Errors: errors})
	}
}

// Custom 404 Handler
func notFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, errorResponse{
		Errors: []apiError{
			{
				Code:    "ERR_NOTFOUND",
				Message: "The requested resource was not found",
			},
		},
	})
}

func dbErrorHandler(err error, c *gin.Context) {
	// Look up the error in the map
	if errInfo, found := models.Errors[err]; found {
		c.JSON(errInfo.Status, errorResponse{
			Errors: []apiError{
				{
					Code:    errInfo.Code,
					Message: errInfo.Message + " Details: " + err.Error(),
				},
			},
		})
		return
	}

	// Default case for unexpected errors
	c.JSON(http.StatusInternalServerError, errorResponse{
		Errors: []apiError{
			{
				Code:    "ERR_INTERNAL",
				Message: "An internal server error occurred. " + err.Error(),
			},
		},
	})
}
