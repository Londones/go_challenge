package server

import (
	"go-challenge/internal/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockQueries struct {
	mock.Mock
}

var testUser = &models.User{
	Email:    "test@example.com",
	Password: "password123",
	Name:     "Test User",
	Address:  strPtr("123 Main St"),
}

func strPtr(s string) *string { return &s }

func TestRegisterHandler(t *testing.T) {
	// Set up database mock
	mockQueries := &MockQueries{}
	mockQueries.On("CreateUser", testUser).Return(nil)

	// Set up server
	s := &Server{}

	// Create a new HTTP request
	form := url.Values{}
	form.Set("email", testUser.Email)
	form.Set("password", testUser.Password)
	form.Set("name", testUser.Name)
	form.Set("address", *testUser.Address)
	req, err := http.NewRequest("POST", "/register", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a new HTTP response recorder
	w := httptest.NewRecorder()

	// Call the RegisterHandler function
	s.RegisterHandler(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusFound, w.Code)

	// Check the response headers
	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, "jwt", cookies[0].Name)
	assert.Equal(t, http.SameSiteLaxMode, cookies[0].SameSite)

	// Check the response redirect URL
	assert.Equal(t, "http://localhost:8000/auth/success", w.Header().Get("Location"))
}

func TestRegisterHandler2(t *testing.T) {
	s := &Server{}

	server := httptest.NewServer(http.HandlerFunc(s.RegisterHandler))
	defer server.Close()

	form := url.Values{}
	form.Set("email", "test@example.com")
	form.Set("password", "password123")
	form.Set("name", "Test User")
	form.Set("address", "123 Main St")
	resp, err := http.Post(server.URL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)

	cookies := resp.Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, "jwt", cookies[0].Name)
	assert.Equal(t, http.SameSiteLaxMode, cookies[0].SameSite)
	assert.Equal(t, "http://localhost:8000/auth/success", resp.Header.Get("Location"))
}
