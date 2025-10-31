package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giovannyptr/bookshelf/models"
)

// ---------- tiny response helpers ----------

type Envelope struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{Data: data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Envelope{Data: data})
}

func Fail(c *gin.Context, status int, msg string) {
	c.JSON(status, Envelope{Error: msg})
}

// ---------- Swagger DTOs (for nicer docs) ----------

// ErrorResponse is used in Swagger examples.
type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}

// PagedBooks is the paginated payload for GET /books (used in Swagger).
type PagedBooks struct {
	Items []models.Book `json:"items"`
	Total int64         `json:"total" example:"42"`
	Page  int           `json:"page"  example:"1"`
	Limit int           `json:"limit" example:"10"`
}
