package product

import (
	"net/http"
	"sync"
)

type Handler struct {
	queries StoreQueries
	mu      *sync.RWMutex
}

func NewHandler(queries StoreQueries) *Handler {
	return &Handler{queries: queries, mu: &sync.RWMutex{}}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /products", h.createProduct)
	router.HandleFunc("GET /products", h.getAllProducts)
	router.HandleFunc("GET /products/{id}", h.getProduct)
}
