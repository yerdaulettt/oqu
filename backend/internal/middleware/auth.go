package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"oqu/internal/auth"
)

var (
	noTokenErr = errors.New("No token found")
)

func unauthResponse(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error": "` + msg + `"}`))
}

func JWTAuthMiddleware(jwtService *auth.JwtAuth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")

			if tokenString == "" {
				unauthResponse(w, http.StatusUnauthorized, noTokenErr.Error())
				return
			}
			tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

			claims, err := jwtService.ParseToken(tokenString)
			if err != nil {
				unauthResponse(w, http.StatusUnauthorized, err.Error())
				return
			}

			if claims.TokenType != "access" {
				unauthResponse(w, http.StatusBadRequest, "No access token")
				return
			}

			ctx := context.WithValue(r.Context(), "role", claims.Role)
			ctx = context.WithValue(ctx, "userId", claims.UserId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Role(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Context().Value("role").(string)

			if !slices.Contains(requiredRoles, role) {
				response := fmt.Sprintf("Only [%s] can access! Your role is [%s]", strings.Join(requiredRoles, ", "), role)
				unauthResponse(w, http.StatusForbidden, response)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
