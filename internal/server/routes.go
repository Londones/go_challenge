package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go-challenge/internal/auth"
	"go-challenge/internal/handlers"
	"go-challenge/internal/utils"

	_ "go-challenge/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	authHandler := handlers.NewAuthHandler(s.dbService)
	userHandler := handlers.NewUserHandler(s.dbService, s.uploadcareClient)
	annonceHandler := handlers.NewAnnonceHandler(s.dbService, s.dbService, s.dbService)
	catHandler := handlers.NewCatHandler(s.dbService, s.uploadcareClient)
	favoriteHandler := handlers.NewFavoriteHandler(s.dbService, s.dbService)
	raceHandler := handlers.NewRaceHandler(s.dbService, s.uploadcareClient)

	r.Group(func(r chi.Router) {
		// Apply JWT middleware to all routes within this group
		r.Use(jwtauth.Verifier(auth.TokenAuth))
		r.Use(jwtauth.Authenticator(auth.TokenAuth))

		r.Group(func(r chi.Router) {
			// Protected routes for admin users
			r.Use(AdminOnly)
			// Admin specific routes
		})
		//** Race routes for admin
		r.Put("/race/{id}", raceHandler.UpdateRaceHandler)
		r.Post("/race", raceHandler.RaceCreationHandler)
		r.Delete("/race/{id}", raceHandler.DeleteRaceHandler)

		r.Group(func(r chi.Router) {
			// Protected routes for personal user data
			r.Use(UserOnly)
			// User specific routes
		})

		//**	Annonces routes
		r.Get("/annonces", annonceHandler.GetAllAnnoncesHandler)
		r.Get("/annonces/{id}", annonceHandler.GetAnnonceByIDHandler)
		r.Post("/annonces", annonceHandler.AnnonceCreationHandler)
		r.Put("/annonces/{id}", annonceHandler.ModifyDescriptionAnnonceHandler)
		r.Delete("/annonces/{id}", annonceHandler.DeleteAnnonceHandler)
		//r.Get("/annonces/cats/{catID}", annonceHandler.FetchAnnonceByCatIDHandler)

		//**	Cats routes
		r.Get("/cats", catHandler.GetAllCatsHandler)
		r.Get("/cats/{id}", catHandler.GetCatByIDHandler)
		r.Put("/cats/{id}", catHandler.UpdateCatHandler)
		r.Post("/cats", catHandler.CatCreationHandler)
		r.Delete("/cats/{id}", catHandler.DeleteCatHandler)
		//r.Get("/cats/", catHandler.FindCatsByFilterHandler)

		//** Race routes
		r.Get("/races", raceHandler.GetAllRaceHandler)
		r.Get("/race/{id}", raceHandler.GetRaceByIDHandler)

		//** User routes
		r.Get("/users", userHandler.GetAllUsersHandler)
		r.Get("/users/annonces/{id}", annonceHandler.GetUserAnnoncesHandler)
		r.Get("/users/{id}", userHandler.GetUserByIDHandler)
		r.Get("/users/current", userHandler.GetCurrentUserHandler)
		r.Post("/users", userHandler.CreateUserHandler)
		r.Post("/profile/picture", userHandler.ModifyProfilePictureHandler)
		r.Put("/users/{id}", userHandler.UpdateUserHandler)
		r.Delete("/users/{id}", userHandler.DeleteUserHandler)

		//** Favorite routes
		r.Post("/favorites", favoriteHandler.FavoriteCreationHandler)
		r.Get("/favorites/users/{userID}", favoriteHandler.GetFavoritesByUserHandler)

		//** Auth routes
		r.Get("/logout/{provider}", authHandler.LogoutProvider)
		r.Get("/logout", authHandler.BasicLogout)
	})

	// Public routes
	r.Get("/annonces/cats/{catID}", annonceHandler.FetchAnnonceByCatIDHandler)
	r.Get("/cats/", catHandler.FindCatsByFilterHandler)
	r.Get("/auth/{provider}/callback", authHandler.GetAuthCallbackFunction)
	r.Get("/auth/{provider}", authHandler.BeginAuthProviderCallback)
	r.Post("/login", authHandler.LoginHandler)
	r.Post("/register", userHandler.RegisterHandler)
	r.Get("/", s.HelloWorldHandler)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(os.Getenv("SERVER_URL")+"/swagger/doc.json"),
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
