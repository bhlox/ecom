package types

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type CustomKey string

const UserIDKey CustomKey = "userID"

type CustomClaims struct {
	jwt.MapClaims
	UserID  uuid.UUID `json:"userID"`
	Expires int64     `json:"expires"`
}

func (c *CustomClaims) Valid() error {
	if c.Expires <= time.Now().Unix() {
		return fmt.Errorf("token is expired")
	}
	return nil
}

type OrderItemPayload struct {
	ProductID int32 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}

type OrderPayload struct {
	OrderItems []OrderItemPayload `json:"orderItems" validate:"required"`
	Address    string             `json:"address" validate:"required"`
}

type CreateProductPayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
}

type RegisterPayload struct {
	Firstname string `json:"firstName" validate:"required"`
	Lastname  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required,min=6,max=50"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

