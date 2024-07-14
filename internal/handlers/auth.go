package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"go-challenge/internal/auth"
	"go-challenge/internal/database/queries"
	// "go-challenge/internal/models"
	// "go-challenge/internal/utils"

	"github.com/go-chi/chi/v5"
	// "github.com/google/uuid"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	userQueries *queries.DatabaseService
}

func NewAuthHandler(userQueries *queries.DatabaseService) *AuthHandler {
	return &AuthHandler{userQueries: userQueries}
}

// GetAuthCallbackFunction godoc
// @Summary Authentication callback
// @Description Completes user authentication with the specified provider
// @Tags auth
// @Accept  json
// @Produce  json
// @Param provider path string true "Authentication Provider"
// @Success 200 {object} models.User "Authenticated user"
// @Failure 500 {string} string "Error message"
// @Router /auth/{provider}/callback [get]
func (h *AuthHandler) GetAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetAuthCallbackFunction")
    // type contextKey string

    // const providerKey contextKey = "provider"
    provider := chi.URLParam(r, "provider")
    r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

    user, err := gothic.CompleteUserAuth(w, r)

	if (err != nil) {
		fmt.Println("error: ", err)
		return
	}

	fmt.Println("user: ", user)

	http.Redirect(w, r, os.Getenv("CLIENT_URL")+"/auth/success", http.StatusFound)

    // utils.Logger("info", "User", fmt.Sprintf("%+v", user), "")
    // if err != nil {
    //     utils.Logger("error", "Complete User Auth:", "Failed to complete user authentication", fmt.Sprintf("Error: %v", err))
    //     fmt.Fprintln(w, err)
    //     return
    // }

    // userRole, err := h.userQueries.GetRoleByName(models.UserRole)
    // if err != nil {
    //     http.Error(w, "error fetching user role", http.StatusInternalServerError)
    //     return
    // }

    // check if user with this google id exists
    // _, err = h.userQueries.FindUserByGoogleID(user.UserID)
    // if err != nil {
    //     // check if user with this email exists
    //     _, err = h.userQueries.FindUserByEmail(user.Email)
    //     if err != nil {
    //         // create user
    //         newUser := &models.User{
    //             ID:            uuid.New().String(),
    //             Email:         user.Email,
    //             Name:          user.Name,
    //             GoogleID:      user.UserID,
    //             ProfilePicURL: user.AvatarURL,
    //         }

    //         err = h.userQueries.CreateUser(newUser, userRole)
    //         if err != nil {
    //             http.Error(w, "error creating user", http.StatusInternalServerError)
    //             return
    //         }

    //         token := auth.MakeToken(newUser.ID, string(newUser.Roles[0].Name))

    //         http.SetCookie(w, &http.Cookie{
    //             HttpOnly: true,
    //             Expires:  time.Now().Add(24 * time.Hour),
    //             Name:     "jwt",
    //             Value:    token,
    //             SameSite: http.SameSiteLaxMode,
    //         })

	// 		fmt.Println("new user: ", newUser)

    //         http.Redirect(w, r, os.Getenv("CLIENT_URL")+"/auth/success", http.StatusFound)
    //     } else {
    //         http.Redirect(w, r, os.Getenv("CLIENT_URL")+"/auth/login", http.StatusFound)
    //     }
    // } else {
    //     http.Redirect(w, r, os.Getenv("CLIENT_URL")+"/auth/login", http.StatusFound)
    // }
}

// LogoutProvider godoc
// @Summary Logout from provider
// @Description Logout from the authentication provider and remove the JWT cookie
// @Tags auth
// @Success 307 {string} string "Redirect location"
// @Router /logout/{provider} [get]
func (h *AuthHandler) LogoutProvider(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)

	http.SetCookie(w, &http.Cookie{
		Name:   "jwt",
		MaxAge: -1,
	})

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// BasicLogout godoc
// @Summary Basic logout
// @Description Remove the JWT cookie and redirect to the success page
// @Tags auth
// @Success 302 {string} string "Redirect location"
// @Router /logout [get]
func (h *AuthHandler) BasicLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "jwt",
		MaxAge: -1,
	})
	http.Redirect(w, r, os.Getenv("CLIENT_URL")+"/auth/success", http.StatusFound)
}

// BeginAuthProviderCallback godoc
// @Summary Begin authentication provider callback
// @Description Start the authentication process with the specified provider
// @Tags auth
// @Accept  json
// @Produce  json
// @Param provider path string true "Authentication Provider"
// @Success 200 {string} string "Authentication process started"
// @Failure 500 {string} string "Error message"
// @Router /auth/{provider} [get]
func (h *AuthHandler) BeginAuthProviderCallback(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
	q.Add("provider", chi.URLParam(r, "provider"))
	r.URL.RawQuery = q.Encode()

	fmt.Println("BeginAuthProviderCallback", chi.URLParam(r, "provider"))

	session, _ := gothic.Store.Get(r, "goth")
	session.Values["provider"] = chi.URLParam(r, "provider")
	session.Save(r, w)

    gothic.BeginAuthHandler(w, r)
}

// LoginHandler godoc
// @Summary Login
// @Description Login with the given email and password
// @Tags auth
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Success 200 {string} string "Login successful"
// @Failure 400 {string} string "Email and password are required"
// @Failure 401 {string} string "Invalid password"
// @Failure 404 {string} string "User not found"
// @Router /login [post]
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	var email, password string

	if strings.Contains(contentType, "application/json") {
		var creds struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		email = creds.Email
		password = creds.Password
	} else {
		r.ParseForm()
		email = r.FormValue("email")
		password = r.FormValue("password")
	}

	if email == "" || password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.userQueries.FindUserByEmail(email)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	token := auth.MakeToken(user.ID, "USER")
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		Name:     "jwt",
		Value:    token,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	response := fmt.Sprintf(`{"success": true, "token": "%s"}`, token)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
