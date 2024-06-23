package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-challenge/internal/api"
	"go-challenge/internal/auth"
	"go-challenge/internal/database/queries"
	"go-challenge/internal/models"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/uploadcare/uploadcare-go/ucare"
)

type UserHandler struct {
	userQueries      *queries.DatabaseService
	uploadcareClient ucare.Client
}

func NewUserHandler(userQueries *queries.DatabaseService, uploadcareClient ucare.Client) *UserHandler {
	return &UserHandler{userQueries: userQueries, uploadcareClient: uploadcareClient}
}

// RegisterHandler godoc
// @Summary Register a new user
// @Description Register a new user with the given email, password, name, address, cp, and ville
// @Tags users
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Param name formData string false "Name"
// @Param addressRue formData string false "Address"
// @Param cp formData string false "CP"
// @Param ville formData string false "Ville"
// @Success 200 {string} string "success"
// @Failure 400 {string} string "email and password are required"
// @Failure 500 {string} string "error creating user"
// @Router /register [post]
func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	var email, password, name, addressRue, cp, ville string

	if strings.Contains(contentType, "application/json") {
		var reqBody struct {
			Email      string `json:"email"`
			Password   string `json:"password"`
			Name       string `json:"name"`
			AddressRue string `json:"addressRue"`
			Cp         string `json:"cp"`
			Ville      string `json:"ville"`
		}

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		email = reqBody.Email
		password = reqBody.Password
		name = reqBody.Name
		addressRue = reqBody.AddressRue
		cp = reqBody.Cp
		ville = reqBody.Ville
	} else {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		email = r.FormValue("email")
		password = r.FormValue("password")
		name = r.FormValue("name")
		addressRue = r.FormValue("addressRue")
		cp = r.FormValue("cp")
		ville = r.FormValue("ville")
	}

	fmt.Println("email: " + email)
	fmt.Println("password: " + password)
	fmt.Println("name: " + name)
	fmt.Println("adresse: " + addressRue)
	fmt.Println("cp: " + cp)
	fmt.Println("ville: " + ville)

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
		AddressRue:    addressRue,
		Cp:            cp,
		Ville:         ville,
		Role:          models.Roles{Name: "user"},
		ProfilePicURL: "default",
	}

	err := h.userQueries.CreateUser(user)
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

	w.Header().Set("Content-Type", "application/json")
	response := fmt.Sprintf(`{"success": true, "token": "%s"}`, token)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// ModifyProfilePictureHandler godoc
// @Summary Modify profile picture
// @Description Modify the profile picture of the authenticated user
// @Tags users
// @Accept multipart/form-data
// @Param uploaded_file formData file true "Image"
// @Success 200 {string} string "Profile picture updated successfully"
// @Failure 500 {string} string "error updating user"
// @Router /profile/picture [post]
func (h *UserHandler) ModifyProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB is the max memory size
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("failed to parse form")
		return
	}

	file, header, err := r.FormFile("uploaded_file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("failed to get file")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)

	tempFile, err := os.CreateTemp(os.TempDir(), "upload-*"+ext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("failed to create temp file")
		return
	}
	defer os.Remove(tempFile.Name()) // clean up

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("failed to copy file")
		return
	}

	FileURL, _, err := api.UploadImage(h.uploadcareClient, tempFile.Name())
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
	user, err := h.userQueries.FindUserByID(userID)
	if err != nil {
		http.Error(w, "error finding user", http.StatusInternalServerError)
		return
	}

	user.ProfilePicURL = FileURL

	err = h.userQueries.UpdateUser(user)
	if err != nil {
		http.Error(w, "error updating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Profile picture updated successfully"))
}
