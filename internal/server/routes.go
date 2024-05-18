package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go-challenge/internal/auth"

	_ "go-challenge/docs"

	_ "go-challenge/docs"
	"go-challenge/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	r.Group(func(r chi.Router) {
		// Apply JWT middleware to all routes within this group
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))

		r.Group(func(r chi.Router) {
			// Protected routes for admin users
			r.Use(AdminOnly)
		})

		r.Group(func(r chi.Router) {
			// Protected routes for personal user data
			r.Use(UserOnly)
		})

		//**	Annonces routes
		r.Get("/annonces", s.GetAllAnnoncesHandler)
		r.Get("/annonces/{id}", s.GetAnnonceByIDHandler)
		r.Post("/annonces", s.AnnonceCreationHandler)
		r.Put("/annonces/{id}", s.ModifyDescriptionAnnonceHandler)
		r.Delete("/annonces/{id}", s.DeleteAnnonceHandler)

		//**	Cats routes
		r.Get("/cats", s.GetAllCatsHandler)
		r.Get("/cats/{id}", s.GetCatByIDHandler)
		r.Post("/cats", s.CatCreationHandler)
		r.Delete("/cats/{id}", s.DeleteCatHandler)

		//** User routes
		r.Post("/profile/picture", s.ModifyProfilePictureHandler)

		//** Auth routes
		r.Get("/logout/{provider}", s.logoutProvider)
		r.Get("/logout", s.basicLogout)
	})

	// Public routes
	r.Get("/auth/{provider}/callback", s.getAuthCallbackFunction)
	r.Get("/auth/{provider}", s.beginAuthProviderCallback)
	r.Post("/login", s.LoginHandler)
	r.Post("/register", s.RegisterHandler)
	r.Get("/", s.HelloWorldHandler)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(os.Getenv("SERVER_URL")+"/swagger/doc.json"), //The url pointing to API definition
	))

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	utils.Logger("debug", "Accès route", "HelloWorld", "")

	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		utils.Logger("fatal", "Route", "HelloWorld", fmt.Sprintf("error handling JSON marshal. Err: %v", err))
	}

	_, _ = w.Write(jsonResp)
}
