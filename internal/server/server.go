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
	"go-challenge/internal/database/queries"
)

type Server struct {
	port             int
	db               *database.Service
	uploadcareClient ucare.Client
	dbService        *queries.DatabaseService
}

func NewServer() (*http.Server, error) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db, err := database.New(&database.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %v", err)
	}
	ucClient, err := api.CreateUCClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create uploadcare client: %v", err)
	}
	dbService := queries.NewQueriesService(db)

	newServer := &Server{
		port:             port,
		db:               db,
		uploadcareClient: ucClient,
		dbService:        dbService,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("Server is running on port %s\n", server.Addr)

	return server, nil
}
