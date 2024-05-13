package middleware

import (
	"context"
	"net/http"

	"github.com/bhlox/ecom/internal/configs"
	"github.com/bhlox/ecom/internal/response"
	"github.com/bhlox/ecom/internal/types"
	"github.com/bhlox/ecom/internal/utils"
	"github.com/golang-jwt/jwt/v4"
)

func VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := utils.GetAuthHeader(r, "Bearer")
		if tokenHeader == "" {
			response.Error(w, http.StatusUnauthorized, "unauthorized user")
			return
		}
		claims := &types.CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.Envs.JWTSECRET), nil
		})
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid token")
			return
		}
		if claims, ok := token.Claims.(*types.CustomClaims); ok {
			if claims.Valid() == nil {
				ctx := context.WithValue(r.Context(), types.UserIDKey, claims.UserID)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			} else {
				response.Error(w, http.StatusBadRequest, "token has expired")
				return
			}
		} else {
			response.Error(w, http.StatusInternalServerError, "unknown claims type")
			return
		}
	})
}
