package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func unauthResponse(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error": "` + msg + `"}`))
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			unauthResponse(w, "no token found")
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			unauthResponse(w, err.Error())
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			unauthResponse(w, "error")
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			unauthResponse(w, "incorrect claim")
			return
		}

		ctx := context.WithValue(r.Context(), "role", role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Role(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Context().Value("role").(string)

			if !slices.Contains(requiredRoles, role) {
				response := fmt.Sprintf("Only [%s] can access! Your role is [%s]", strings.Join(requiredRoles, ", "), role)
				unauthResponse(w, response)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
