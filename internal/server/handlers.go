package server

import (
	"context"
	"encoding/json"
	"fmt"
	"go-challenge/internal/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go-challenge/internal/api"
	"go-challenge/internal/auth"
	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"

	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
)

// getAuthCallbackFunction godoc
// @Summary Authentication callback
// @Description Completes user authentication with the specified provider
// @Tags auth
// @Accept  json
// @Produce  json
// @Param provider path string true "Authentication Provider"
// @Success 200 {object} models.User "Authenticated user"
// @Failure 500 {string} string "Error message"
// @Router /auth/{provider}/callback [get]
func (s *Server) getAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	type contextKey string

	queriesService := queries.NewQueriesService(s.dbService)

	const providerKey contextKey = "provider"

	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(context.Background(), providerKey, provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fmt.Println(user)
	// check if user with this google id exists
	existingUser, err := queriesService.FindUserByGoogleID(user.UserID)
	if err != nil {
		// check if user with this email exists
		existingUser, err = queriesService.FindUserByEmail(user.Email)
		if err != nil {
			// create user
			newUser := &models.User{
				ID:       uuid.New().String(),
				Email:    user.Email,
				Name:     user.Name,
				GoogleID: user.UserID,
				Role:     models.Roles{Name: "user"},
			}

			err := queriesService.CreateUser(newUser)
			if err != nil {
				http.Error(w, "error creating user", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "An account has already been registered with this email", http.StatusConflict)
			return
		}
	}

	token := auth.MakeToken(existingUser.ID, string(existingUser.Role.Name))

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		Name:     "jwt",
		Value:    token,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, os.Getenv("CLIENT_URL")+"/auth/success", http.StatusFound)
}

// logoutProvider godoc
// @Summary Logout from provider
// @Description Logout from the authentication provider and remove the JWT cookie
// @Tags auth
// @Produce  json
// @Success 307 {header} string "Location" "Redirect location"
// @Router /logout/{provider} [get]
func (s *Server) logoutProvider(res http.ResponseWriter, req *http.Request) {
	gothic.Logout(res, req)

	//remove the cookie
	http.SetCookie(res, &http.Cookie{
		Name:   "jwt",
		MaxAge: -1,
	})

	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusTemporaryRedirect)
}

// basicLogout godoc
// @Summary Basic logout
// @Description Remove the JWT cookie and redirect to the success page
// @Tags auth
// @Produce  json
// @Success 302 {header} string "Location" "Redirect location"
// @Router /logout [get]
func (s *Server) basicLogout(w http.ResponseWriter, r *http.Request) {
	utils.Logger("debug", "Accès route", "basicLogout", "")
	http.SetCookie(w, &http.Cookie{
		Name:   "jwt",
		MaxAge: -1,
	})
	http.Redirect(w, r, os.Getenv("CLIENT_URL")+"/auth/success", http.StatusFound)
}

// beginAuthProviderCallback godoc
// @Summary Begin authentication provider callback
// @Description Start the authentication process with the specified provider
// @Tags auth
// @Accept  json
// @Produce  json
// @Param provider path string true "Authentication Provider"
// @Success 200 {string} string "Authentication process started"
// @Failure 500 {string} string "Error message"
// @Router /auth/{provider} [get]
func (s *Server) beginAuthProviderCallback(w http.ResponseWriter, r *http.Request) {
	type contextKey string

	const providerKey contextKey = "provider"

	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(context.Background(), providerKey, provider))

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
// @Success 302 {header} string "Location" "Redirect location"
// @Header 302 {string} Set-Cookie "jwt={token}; HttpOnly; SameSite=Lax; Expires={expiry}"
// @Failure 400 {string} string "email and password are required"
// @Failure 404 {string} string "user not found"
// @Failure 401 {string} string "invalid password"
// @Router /login [post]
func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	utils.Logger("debug", "Accès route", "Login", "")

	queriesService := queries.NewQueriesService(s.dbService)

	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	user, err := queriesService.FindUserByEmail(email)

	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	fmt.Println(user)
	token := auth.MakeToken(user.ID, string(user.Role.Name))

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		Name:     "jwt",
		Value:    token,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, os.Getenv("CLIENT_URL")+"/auth/success", http.StatusFound)
}

// RegisterHandler godoc
// @Summary Register a new user
// @Description Register a new user with the given email, password, name, address, cp, and city
// @Tags users
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Param name formData string false "Name"
// @Param address formData string false "Address"
// @Param cp formData string false "CP"
// @Param city formData string false "City"
// @Success 201 {string} string
// @Header 201 {string} Set-Cookie "jwt={token}; HttpOnly; SameSite=Lax; Expires={expiry}"
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Failure 500 {string} string
// @Router /register [post]
func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	utils.Logger("debug", "Accès route", "Register", "")
	queriesService := queries.NewQueriesService(s.dbService)

	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	name := r.FormValue("name")
	address := r.FormValue("address")
	cp := r.FormValue("cp")
	city := r.FormValue("city")

	if email == "" || password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, passwordError := auth.HashPassword(password)
	if passwordError != nil {
		http.Error(w, "error hashing password", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		ID:            uuid.New().String(),
		Email:         email,
		Password:      hashedPassword,
		Name:          name,
		AddressRue:    address,
		Cp:            cp,
		Ville:         city,
		Role:          models.Roles{Name: "user"},
		ProfilePicURL: "default",
	}

	err := queriesService.CreateUser(user)

	if err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	token := auth.MakeToken(user.ID, "user")

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		Name:     "jwt",
		Value:    token,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, os.Getenv("CLIENT_URL")+"/auth/success", http.StatusCreated)
}

// CatCreationHandler godoc
// @Summary Create cat
// @Description Create a new cat with the provided details
// @Tags cats
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param name formData string true "Name"
// @Param uploaded_file formData file true "Image"
// @Success 201 {object} models.Cats "Created cat"
// @Failure 400 {string} string "all fields are required"
// @Failure 500 {string} string "error creating cat"
// @Router /cat [post]
func (s *Server) CatCreationHandler(w http.ResponseWriter, r *http.Request) {
	queriesService := queries.NewQueriesService(s.dbService)
	fileURLs := make([]string, 0)

	err := r.ParseMultipartForm(10 << 20) // 10 MB is the max memory size
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")

	if name == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	files := r.MultipartForm.File["uploaded_file"]
	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Get the extension of the uploaded file
		ext := filepath.Ext(header.Filename)

		// Create a temporary file with the same extension
		tempFile, err := os.CreateTemp(os.TempDir(), "upload-*"+ext)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name()) // clean up

		// Copy the uploaded file to the temporary file
		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Upload the file to uploadcare
		FileURL, _, err := api.UploadImage(s.uploadcareClient, tempFile.Name())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fileURLs = append(fileURLs, FileURL)
	}

	cat := &models.Cats{
		Name:        name,
		PicturesURL: fileURLs,
	}

	_, err = queriesService.CreateCat(cat)
	if err != nil {
		http.Error(w, "error creating cat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cat)
	if err != nil {
		http.Error(w, "error encoding cat to JSON", http.StatusInternalServerError)
		return
	}

}

// ModifyProfilePictureHandler godoc
// @Summary Modify profile picture
// @Description Modify the profile picture of the authenticated user
// @Tags users
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param uploaded_file formData file true "Image"
// @Success 200 {string} string "Profile picture updated successfully"
// @Failure 500 {string} string "error getting claims"
// @Failure 500 {string} string "error finding user"
// @Failure 500 {string} string "error updating user"
// @Router /profile/picture [post]
func (s *Server) ModifyProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	queriesService := queries.NewQueriesService(s.dbService)

	err := r.ParseMultipartForm(10 << 20) // 10 MB is the max memory size
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("failed to parse form")
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("uploaded_file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("failed to get file")
		return
	}
	defer file.Close()

	// Get the extension of the uploaded file
	ext := filepath.Ext(header.Filename)

	// Create a temporary file with the same extension
	tempFile, err := os.CreateTemp(os.TempDir(), "upload-*"+ext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("failed to create temp file")
		return
	}
	defer os.Remove(tempFile.Name()) // clean up

	// Copy the uploaded file to the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("failed to copy file")
		return
	}

	// Upload the file to uploadcare
	FileURL, _, err := api.UploadImage(s.uploadcareClient, tempFile.Name())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		errorMsg := fmt.Errorf("failed to upload file to uploadcare: %v", err)
		fmt.Println(errorMsg)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "error getting claims", http.StatusInternalServerError)
		return
	}

	userID := claims["id"].(string)
	user, err := queriesService.FindUserByID(userID)
	if err != nil {
		http.Error(w, "error finding user", http.StatusInternalServerError)
		return
	}

	user.ProfilePicURL = FileURL

	err = queriesService.UpdateUser(user)
	if err != nil {
		http.Error(w, "error updating user", http.StatusInternalServerError)
		return
	}

	// return success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Profile picture updated successfully"))
}

// AssociationCreationHandler godoc
// @Summary Create association
// @Description Create a new association with the provided details
// @Tags associations
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param name formData string true "Name"
// @Param address formData string true "Address"
// @Param cp formData string true "Postal Code"
// @Param city formData string true "City"
// @Param phone formData string true "Phone"
// @Param email formData string true "Email"
// @Success 201 {header} string "Location" "Redirect location"
// @Failure 400 {string} string "all fields are required"
// @Failure 500 {string} string "error getting claims"
// @Failure 500 {string} string "error finding user"
// @Failure 500 {string} string "error creating association"
// @Router /association [post]
func (s *Server) AssociationCreationHandler(w http.ResponseWriter, r *http.Request) {
	queriesService := queries.NewQueriesService(s.dbService)

	r.ParseForm()
	name := r.FormValue("name")
	address := r.FormValue("address")
	cp := r.FormValue("cp")
	city := r.FormValue("city")
	phone := r.FormValue("phone")
	email := r.FormValue("email")

	if name == "" || address == "" || cp == "" || city == "" || phone == "" || email == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "error getting claims", http.StatusInternalServerError)
		return
	}

	userID := claims["id"].(string)
	user, err := queriesService.FindUserByID(userID)
	if err != nil {
		http.Error(w, "error finding user", http.StatusInternalServerError)
		return
	}

	association := &models.Association{
		Name:       name,
		AddressRue: address,
		Cp:         cp,
		Ville:      city,
		Phone:      phone,
		Email:      email,
		MemberIDs:  []string{user.ID},
		Verified:   false,
	}

	id, err := queriesService.CreateAssociation(association)
	if err != nil {
		http.Error(w, "error creating association", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf(os.Getenv("CLIENT_URL")+"/association/%d", id), http.StatusCreated)
}

// AnnonceCreationHandler godoc
// @Summary Create annonce
// @Description Create a new annonce with the provided details
// @Tags annonces
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param description formData string true "Description of the annonce"
// @Param userID formData string true "User ID"
// @Param cats formData []string true "Categories of the annonce"
// @Param favorite formData []string true "Favorite list of the annonce"
// @Param rating formData []string true "Rating of the annonce"
// @Success 201 {string} string "Location of the created annonce"
// @Failure 400 {string} string "Missing or invalid fields in the request"
// @Failure 500 {string} string "Internal server error"
// @Router /annonce [post]
func (s *Server) AnnonceCreationHandler(w http.ResponseWriter, r *http.Request) {
	queriesService := queries.NewQueriesService(s.dbService)

	r.ParseForm()
	description := r.FormValue("description")

	if description == "" {
		http.Error(w, "description is required", http.StatusBadRequest)
		return
	}

	// Obtient l'ID de l'utilisateur à partir des revendications
	userID := r.FormValue("userID")
	if userID == "" {
		http.Error(w, "user id is required", http.StatusInternalServerError)
		return
	}

	// Trouve l'utilisateur dans la base de données
	user, err := queriesService.FindUserByID(userID)
	if err != nil {
		http.Error(w, "error finding user", http.StatusInternalServerError)
		return
	}

	cats := []models.Cats{}
	favorites := []models.Favorite{}
	ratings := []models.Rating{}

	// Crée une nouvelle annonce
	annonce := &models.Annonce{
		Description: &description,
		UserID:      user.ID,
		Cats:        cats,
		Favorite:    favorites,
		Rating:      ratings,
	}

	// Crée l'annonce dans la base de données
	createAnnonce, err := queriesService.CreateAnnonce(annonce)
	if err != nil {
		http.Error(w, "error creating annonce", http.StatusInternalServerError)
		return
	}

	// Renvoie l'annonce créée au format JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createAnnonce)
}
