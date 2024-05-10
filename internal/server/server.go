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
	dbService        *database.Service
}

func NewServer() (*http.Server, error) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db, err := database.New(&database.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %v", err)
	}
	client, err := api.CreateUCClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create uploadcare client: %v", err)
	}
	NewServer := &Server{
		port:             port,
		db:               db,
		uploadcareClient: client,
		dbService:        db,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("Server is running on port %s\n", server.Addr)

	return server, nil
}
