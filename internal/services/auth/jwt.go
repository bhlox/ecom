package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateJWT(secret string, userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"userID":  userID,
		"expires": time.Now().AddDate(0, 0, 30).Unix(),
		// "role": "admin"
	}
	claims.Valid()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("error encountered when signing token")
	}
	return signedToken, nil
}
