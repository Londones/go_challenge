package auth

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"golang.org/x/crypto/bcrypt"
)

const (
	maxAge = 86400 * 60
	isProd = false
)

var TokenAuth *jwtauth.JWTAuth

var secret = os.Getenv("JWT_SECRET")

func MakeToken(id string, role string) string {
	_, tokenString, _ := TokenAuth.Encode(map[string]interface{}{"id": id, "role": role})
	return tokenString
}

func GetTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("jwt") // replace with your cookie name
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func NewAuth() {
	var key = os.Getenv("SESSION_KEY")
	TokenAuth = jwtauth.New("HS256", []byte(secret), nil)

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, "http://localhost:8080/auth/google/callback"),
	)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
