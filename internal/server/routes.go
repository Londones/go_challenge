package server

import (
	"encoding/json"
	"log"
	"net/http"

	"go-challenge/internal/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		// Protected routes for any authenticated user
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))

		r.Get("/", s.HelloWorldHandler)
		r.Get("/logout/{provider}", s.logoutProvider)
		r.Get("/logout", s.basicLogout)
	})

	r.Group(func(r chi.Router) {
		// Protected routes for admin users
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))
		r.Use(AdminOnly)

	})

	r.Group(func(r chi.Router) {
		// Protected routes for personal user data
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))
		r.Use(UserOnly)

	})

	// Public routes
	r.Get("/auth/{provider}/callback", s.getAuthCallbackFunction)
	r.Get("/auth/{provider}", s.beginAuthProviderCallback)
	r.Post("/login", s.LoginHandler)
	r.Post("/register", s.RegisterHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
