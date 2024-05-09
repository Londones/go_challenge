package main

import (
	"fmt"
	"go-challenge/internal/auth"
	"go-challenge/internal/server"
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
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

	// print that the server is running
	fmt.Printf("Server is running on port %s", server.Addr)
}
