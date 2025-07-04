package middleware

import (
	"backend2/internal/auth"
	"fmt"
	"net/http"
	"strings"
)

func AuthMiddleware(auth *auth.AuthUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			fmt.Fprintf(w, authHeader)
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "unauthorized: missing token", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			_, valid := auth.ValidateToken(tokenStr)
			if !valid {
				http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
