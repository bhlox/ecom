package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthChecks(t *testing.T) {
	router := http.NewServeMux()

	healthHandler := NewHandler()
	router.HandleFunc("GET /health/", healthHandler.success)
	router.HandleFunc("GET /health/error", healthHandler.error)

	t.Run("health should pass", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health/", nil)
		if err != nil {
			t.Errorf("error occured :%v", err.Error())
		}
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code, fmt.Sprintf("should have a status code of %v", http.StatusOK))
	})

	t.Run("error test", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/health/error", nil)
		if err != nil {
			t.Errorf("error occured :%v", err.Error())
		}
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code, "it should be an error code of 500")

		var resp map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NotNil(t, resp["error"], "Response should contain an error field")
	})
}
