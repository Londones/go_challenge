package server

import (
	"context"
	"database/sql"
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

	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
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

// **ANNONCES
// GetAllAnnoncesHandler godoc
// @Summary Get all annonces
// @Description Retrieve all annonces from the database
// @Tags annonces
// @Produce json
// @Success 200 {array} models.Annonce "List of annonces"
// @Failure 500 {string} string "Internal server error"
// @Router /annonces [get]
func (s *Server) GetAllAnnoncesHandler(w http.ResponseWriter, r *http.Request) {
	queriesService := queries.NewQueriesService(s.dbService)

	annonces, err := queriesService.GetAllAnnonces()
	if err != nil {
		http.Error(w, "error getting annonces", http.StatusInternalServerError)
		return
	}

	// Renvoie les annonces sous forme de réponse JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(annonces)
}

// GetAnnonceByIDHandler godoc
// @Summary Get an annonce by ID
// @Description Retrieve an annonce from the database by its ID
// @Tags annonces
// @Produce json
// @Param id path string true "ID of the annonce to retrieve"
// @Success 200 {object} models.Annonce "Annonce details"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Annonce not found"
// @Failure 500 {string} string "Internal server error"
// @Router /annonce/{id} [get]
func (s *Server) GetAnnonceByIDHandler(w http.ResponseWriter, r *http.Request) {
	queriesService := queries.NewQueriesService(s.dbService)

	annonceID := chi.URLParam(r, "id")
	if annonceID == "" {
		http.Error(w, "ID of the annonce is required", http.StatusBadRequest)
		return
	}

	annonce, err := queriesService.FindAnnonceByID(annonceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Annonce not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving annonce", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(annonce)
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

	fmt.Println(user.ID)

	// Crée une nouvelle annonce
	annonce := &models.Annonce{
		Description: &description,
		UserID:      user.ID,
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

// ModifyDescriptionAnnonceHandler godoc
// @Summary Modify annonce description
// @Description Modify the description of an existing annonce
// @Tags annonces
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Param id path string true "ID of the annonce to modify"
// @Param description formData string true "New description of the annonce"
// @Security ApiKeyAuth
// @Success 200 {object} models.Annonce "Updated annonce"
// @Failure 400 {string} string "Missing or invalid fields in the request"
// @Failure 403 {string} string "User is not authorized to modify this annonce"
// @Failure 404 {string} string "Annonce not found"
// @Failure 500 {string} string "Internal server error"
// @Router /annonce/{id} [put]
func (s *Server) ModifyDescriptionAnnonceHandler(w http.ResponseWriter, r *http.Request) {
	queriesService := queries.NewQueriesService(s.dbService)

	r.ParseForm()

	// Récup l'id de l'annonce à modifier depuis les params de la requête
	annonceID := chi.URLParam(r, "id")

	// Get updated description
	description := r.FormValue("description")

	if description == "" {
		http.Error(w, "description is required", http.StatusBadRequest)
		return
	}

	// Get user ID
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

	existingAnnonce, err := queriesService.FindAnnonceByID(annonceID)
	if err != nil {
		http.Error(w, "error finding annonce", http.StatusNotFound)
		return
	}

	if existingAnnonce.UserID != user.ID {
		http.Error(w, "user is not authorized to modify this annonce", http.StatusForbidden)
		return
	}

	// Update description
	updatedAnnonce, err := queriesService.UpdateAnnonceDescription(annonceID, description)
	if err != nil {
		http.Error(w, "Error updating annonce", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedAnnonce)
}

// DeleteAnnonceHandler godoc
// @Summary Delete annonce
// @Description Delete an existing annonce
// @Tags annonces
// @Security ApiKeyAuth
// @Param id path string true "ID of the annonce to delete"
// @Success 204 {string} string "Annonce deleted successfully"
// @Failure 403 {string} string "User is not authorized to delete this annonce"
// @Failure 404 {string} string "Annonce not found"
// @Failure 500 {string} string "Internal server error"
// @Router /annonce/{id} [delete]
func (s *Server) DeleteAnnonceHandler(w http.ResponseWriter, r *http.Request) {
	queriesService := queries.NewQueriesService(s.dbService)

	// Get l'id de l'annonce à supprimer depuis les params de la requête
	annonceID := chi.URLParam(r, "id")

	annonce, err := queriesService.FindAnnonceByID(annonceID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "annonce not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error finding annonce", http.StatusInternalServerError)
		return
	}

	// Get user ID
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "error getting claims", http.StatusInternalServerError)
		return
	}
	userID := claims["id"].(string)

	if annonce.UserID != userID {
		http.Error(w, "user is not authorized to modify this annonce", http.StatusForbidden)
		return
	}

	// Supprimer l'annonce
	if err := queriesService.DeleteAnnonce(annonceID); err != nil {
		http.Error(w, "error deleting annonce", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//**ANNONCES

// **CHATS

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

// GetAllCatsHandler godoc
// @Summary Get all cats
// @Description Retrieve a list of all cats
// @Tags cats
// @Produce  json
// @Success 200 {array} models.Cats "List of cats"
// @Failure 500 {string} string "error fetching cats"
// @Router /cats [get]
func (s *Server) GetAllCatsHandler(w http.ResponseWriter, r *http.Request) {
	queriesService := queries.NewQueriesService(s.dbService)

	cats, err := queriesService.GetAllCats()
	if err != nil {
		http.Error(w, "error fetching cats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cats)
	if err != nil {
		http.Error(w, "error encoding cats to JSON", http.StatusInternalServerError)
		return
	}
}

// GetCatByIDHandler godoc
// @Summary Get cat by ID
// @Description Retrieve a cat by its ID
// @Tags cats
// @Param id query string true "Cat ID"
// @Produce  json
// @Success 200 {object} models.Cats "Found cat"
// @Failure 400 {string} string "cat ID is required"
// @Failure 404 {string} string "cat not found"
// @Failure 500 {string} string "error fetching cat"
// @Router /cat/{id} [get]
func (s *Server) GetCatByIDHandler(w http.ResponseWriter, r *http.Request) {
	queriesService := queries.NewQueriesService(s.dbService)
	params := r.URL.Query()
	id := params.Get("id")

	if id == "" {
		http.Error(w, "cat ID is required", http.StatusBadRequest)
		return
	}

	cat, err := queriesService.FindCatByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("cat with ID %s not found", id), http.StatusNotFound)
			return
		}
		http.Error(w, "error fetching cat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cat)
	if err != nil {
		http.Error(w, "error encoding cat to JSON", http.StatusInternalServerError)
		return
	}
}

// DeleteCatByIDHandler godoc
// @Summary Delete cat by ID
// @Description Delete a cat by its ID
// @Tags cats
// @Param id query string true "Cat ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "cat ID is required"
// @Failure 404 {string} string "cat not found"
// @Failure 500 {string} string "error deleting cat"
// @Router /cat/{id} [delete]
func (s *Server) DeleteCatByIDHandler(w http.ResponseWriter, r *http.Request) {
	queriesService := queries.NewQueriesService(s.dbService)
	params := r.URL.Query()
	id := params.Get("id")

	if id == "" {
		http.Error(w, "cat ID is required", http.StatusBadRequest)
		return
	}

	err := queriesService.DeleteCatByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("cat with ID %s not found", id), http.StatusNotFound)
			return
		}
		http.Error(w, "error deleting cat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// **CHATS
