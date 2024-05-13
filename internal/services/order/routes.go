package order

import (
	"net/http"
	"sync"

	"github.com/bhlox/ecom/internal/middleware"
)

type Handler struct {
	queries OrderQueries
	mu      *sync.RWMutex
}

func NewHandler(queries OrderQueries) *Handler {
	return &Handler{queries: queries, mu: &sync.RWMutex{}}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.Handle("GET /order/{orderId}", middleware.VerifyJWT(http.HandlerFunc(h.getOrder)))
}
