package user

import (
	"net/http"
	"sync"
)

type Handler struct {
	queries UserQueries
	mu      *sync.RWMutex
}

func NewHandler(queries UserQueries) *Handler {
	return &Handler{queries: queries, mu: &sync.RWMutex{}}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /login", h.login)
	router.HandleFunc("POST /register", h.register)
	router.HandleFunc("GET /users", h.getAllUsers)
	router.HandleFunc("DELETE /users/{id}", h.deleteUser)
}
