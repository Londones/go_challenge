package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/uploadcare/uploadcare-go/ucare"

	"go-challenge/internal/api"
	"go-challenge/internal/database"
)

type Server struct {
	port             int
	db               database.Database
	uploadcareClient ucare.Client
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db, _ := database.New(&database.Config{})
	client, _ := api.CreateUCClient()
	NewServer := &Server{
		port:             port,
		db:               db,
		uploadcareClient: client,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("Server is running on port %s", server.Addr)

	return server
}
