package user

import (
	"testing"
)

func TestVerifyPW(t *testing.T) {
	storedPW := "JDJhJDEwJFJsOHYycWU1OVdYUkVyUGtKQTZsQ09Oa1hKaDFDdEYuaDk5cVJ3MTFqb3ZCZWpJSm5CNDg2"
	payloadPW := "password123"

	t.Run("test valid password", func(t *testing.T) {
		if err := verifyPW(payloadPW, storedPW); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	// Test with an invalid password
	t.Run("test invalid password", func(t *testing.T) {
		if err := verifyPW("wrongPassword", storedPW); err == nil {
			t.Errorf("Expected an error, got nil")
		}
	})
}
