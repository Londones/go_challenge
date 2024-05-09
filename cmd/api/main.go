package main

import (
	"fmt"
	"go-challenge/internal/auth"
	"go-challenge/internal/server"
	"net/http"

	"github.com/rs/cors"
)

//    @title            GO-challenge-PurrfectMatch
//    @version        1.0
//    @description    Swagger de PurrfectMatch
//    @termsOfService    http://swagger.io/terms/

//    @contact.name    API Support
//    @contact.url    http://www.swagger.io/support
//    @contact.email    support@swagger.io

//    @license.name    Apache 2.0
//    @license.url    http://www.apache.org/licenses/LICENSE-2.0.html

// @host        localhost:8080
// @BasePath    /

// Pour lancer le swagger : swag init --parseDependency -d ./internal/server -g ../../cmd/api/main.go
// puis supprimer les lignes
func main() {
	auth.NewAuth()
	server, err := server.NewServer()
	if err != nil {
		panic(fmt.Sprintf("cannot create server: %s", err))
	}

	mux := http.NewServeMux()
	handler := cors.Default().Handler(mux)
	mux.Handle("/", server.Handler)

	fmt.Println("Server is running on port 8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
