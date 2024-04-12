package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go-challenge/internal/auth"
	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/markbates/goth/gothic"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth))

		r.Get("/", s.HelloWorldHandler)
		r.Get("/logout/{provider}", s.logoutProvider)
		r.Get("/logout", s.basicLogout)
	})

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

func (s *Server) getAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fmt.Println(user)

	token := auth.MakeToken(user.Email)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		Name:     "jwt",
		Value:    token,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "http://localhost:8000/auth/success", http.StatusFound)
}

func (s *Server) logoutProvider(res http.ResponseWriter, req *http.Request) {
	gothic.Logout(res, req)

	//remove the cookie
	http.SetCookie(res, &http.Cookie{
		Name:   "jwt",
		MaxAge: -1,
	})

	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) basicLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "jwt",
		MaxAge: -1,
	})
	http.Redirect(w, r, "http://localhost:8000/auth/success", http.StatusFound)
}

func (s *Server) beginAuthProviderCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	gothic.BeginAuthHandler(w, r)
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	user, err := queries.FindUserByEmail(email)

	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	fmt.Println(user)
	token := auth.MakeToken(email)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		Name:     "jwt",
		Value:    token,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "http://localhost:8000/auth/success", http.StatusFound)
}

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	name := r.FormValue("name")
	address := r.FormValue("address")

	if email == "" || password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, passwordError := auth.HashPassword(password)
	if passwordError != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Email:    email,
		Password: hashedPassword,
		Name:     name,
		Address:  &address,
		Roles:    []models.Roles{{Name: "user"}},
	}

	err := queries.CreateUser(user)

	if err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	token := auth.MakeToken(email)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		Name:     "jwt",
		Value:    token,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "http://localhost:8000/auth/success", http.StatusFound)
}
