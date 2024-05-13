package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	id := uuid.New()
	t.Run("generate JWT", func(t *testing.T) {
		token, err := GenerateJWT("custom_secret", id)
		if err != nil {
			t.Errorf("error creating JWT: %v", err)
		}
		if token == "" {
			t.Error("Generated token is empty")
		}
	})
}
