package server

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func AdminOnly(next http.Handler, token *jwtauth.JWTAuth) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		role := claims["role"]

		if role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
