package utils

import (
	"os"

	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
)

var tokenAuth *jwtauth.JWTAuth
var Secret string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	Secret = os.Getenv("JWT_SECRET")
}
