package main

import (
	"fmt"
	"go-challenge/internal/auth"
	"go-challenge/internal/server"
	"net/http"
	"os"

	"github.com/rs/cors"
)

// @title           GO-challenge-PurrfectMatch
// @version         1.0
// @description     Swagger de PurrfectMatch
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /

// Pour lancer le swagger : swag init --parseDependency -d internal/handlers -g ../../cmd/api/main.go
func main() {
	auth.NewAuth()
	server, err := server.NewServer()
	if err != nil {
		panic(fmt.Sprintf("cannot create server: %s", err))
	}

	// Création d'un nouveau ServeMux
	mux := http.NewServeMux()

	// Handle CORS for the entire ServeMux
    corsHandler := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: false,
        MaxAge:           300,
    })

    handler := corsHandler.Handler(mux)

	// Définir le gestionnaire pour la racine du ServeMux
	mux.Handle("/", server.Handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println(port)

	// Lancement du serveur
	fmt.Println("Server is running on port" + port)
	
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
