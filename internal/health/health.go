package health

import (
	"net/http"

	"github.com/bhlox/ecom/internal/response"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /", h.success)
	router.HandleFunc("GET /error", h.error)
}

func (h *Handler) success(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) error(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusInternalServerError, "error test")
}
