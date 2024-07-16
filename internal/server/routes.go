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
	r.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("/internal/uploads/"))))
	r.Handle("/.well-known/*", http.StripPrefix("/.well-known/", http.FileServer(http.Dir("assets"))))

	authHandler := handlers.NewAuthHandler(s.dbService)
	userHandler := handlers.NewUserHandler(s.dbService, s.uploadcareClient)
	annonceHandler := handlers.NewAnnonceHandler(s.dbService, s.dbService, s.dbService)
	catHandler := handlers.NewCatHandler(s.dbService, s.uploadcareClient)
	favoriteHandler := handlers.NewFavoriteHandler(s.dbService, s.dbService)
	raceHandler := handlers.NewRaceHandler(s.dbService, s.uploadcareClient)
	associationHandler := handlers.NewAssociationHandler(s.dbService, s.uploadcareClient)
	ratingHandler := handlers.NewRatingHandler(s.dbService, s.dbService)
	roomHandler := handlers.NewRoomHandler(s.dbService)

	roomHandler.LoadRooms()

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

		r.Group(func(r chi.Router) {
			// Protected routes for personal user data
			r.Use(UserOnly)
			// User specific routes
		})

		//**	Rating routes
		r.Get("/ratings", ratingHandler.FetchAllRatingsHandler)
		r.Get("/ratings/{id}", ratingHandler.GetRatingByIDHandler)
		r.Post("/ratings", ratingHandler.CreateRatingHandler)
		r.Put("/ratings/{id}", ratingHandler.UpdateRatingHandler)
		r.Delete("/ratings/{id}", ratingHandler.DeleteRatingHandler)
		r.Get("/ratings/user/{userID}", ratingHandler.GetUserRatingsHandler)
		r.Get("/ratings/author/{authorID}", ratingHandler.GetAuthorsRatingsHandler)

		//**	Annonces routes
		r.Get("/annonces", annonceHandler.GetAllAnnoncesHandler)
		r.Get("/annonces/{id}", annonceHandler.GetAnnonceByIDHandler)
		r.Post("/annonces", annonceHandler.AnnonceCreationHandler)
		r.Put("/annonces/{id}", annonceHandler.ModifyAnnonceHandler)
		r.Delete("/annonces/{id}", annonceHandler.DeleteAnnonceHandler)
		r.Get("/annonces/cats/{catID}", annonceHandler.FetchAnnonceByCatIDHandler)
		r.Get("annonce/address/{id}", annonceHandler.GetAddressFromUserID)

		//**	Cats routes
		r.Get("/cats", catHandler.GetAllCatsHandler)
		r.Get("/cats/{id}", catHandler.GetCatByIDHandler)
		r.Put("/cats/{id}", catHandler.UpdateCatHandler)
		r.Post("/cats", catHandler.CatCreationHandler)
		r.Delete("/cats/{id}", catHandler.DeleteCatHandler)
		r.Get("/cats/", catHandler.FindCatsByFilterHandler)
		r.Get("/cats/user/{userID}", catHandler.GetCatsByUserHandler)

		//** Race routes
		r.Get("/races", raceHandler.GetAllRaceHandler)
		r.Get("/race/{id}", raceHandler.GetRaceByIDHandler)
		r.Post("/races", raceHandler.CreateRaceHandler)
		r.Put("/races/{id}", raceHandler.UpdateRaceHandler)
		r.Delete("/races/{id}", raceHandler.DeleteRaceHandler)

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

		//** Association routes
		r.Post("/associations", associationHandler.CreateAssociationHandler)
		r.Get("/associations", associationHandler.GetAllAssociationsHandler)
		r.Get("/users/{userId}/associations", associationHandler.GetUserAssociationsHandler)
		r.Get("/associations/{id}", associationHandler.GetAssociationByIdHandler)
		r.Put("/associations/{id}/verify", associationHandler.UpdateAssociationVerifyStatusHandler)
		r.Delete("/associations/{id}", associationHandler.DeleteAssociationHandler)

		//** Chat routes
		r.Get("/rooms", roomHandler.GetUserRooms)
		r.Get("/rooms/{roomID}", roomHandler.GetRoomMessages)
		r.Get("/ws/{roomID}", roomHandler.HandleWebSocket)

	})

	// Public routes
	// r.Handle("/auth/success/", http.StripPrefix("/assets/", http.FileServer(http.Dir("/assets/success.html"))))
	r.Handle("/auth/success", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/success.html")
	}))
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
