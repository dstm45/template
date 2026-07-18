// Package middlewares contient tous les middlewares
package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/dstm45/template/pkg/services"
)

type CtxKey string

const (
	UUIDKey CtxKey = "uuid"
	RoleKey CtxKey = "role"
)

func AuthMiddleware(next http.HandlerFunc, authService services.IAuthService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value

		claims, err := authService.ParseAccessToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UUIDKey, claims.UUID)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func IsAdminMiddleware(next http.HandlerFunc, authService services.IAuthService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value

		claims, err := authService.ParseAccessToken(tokenString)
		if err != nil || claims.Role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			log.Println("Access denied: Admin role required")
			return
		}
		ctx := context.WithValue(r.Context(), UUIDKey, claims.UUID)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
