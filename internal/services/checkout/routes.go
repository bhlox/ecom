package checkout

import (
	"net/http"
	"sync"

	"github.com/bhlox/ecom/internal/middleware"
)

type Handler struct {
	queries CheckoutQueries
	mu      *sync.RWMutex
}

func NewHandler(queries CheckoutQueries) *Handler {
	return &Handler{queries: queries, mu: &sync.RWMutex{}}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.Handle("POST /checkout", middleware.VerifyJWT(http.HandlerFunc(h.checkout)))
}
