package queries

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go-challenge/internal/database"
	queries "go-challenge/internal/database/queries/mocks"
	"go-challenge/internal/models"
	"go-challenge/internal/utils"
	"testing"
)

func TestHandlers(t *testing.T) {
	s, err := database.TestDatabaseInit()
	if err != nil {
		return
	}
	utils.Logger("debug", "GO-Tests", "Database", "Test database created")

	// Auth

	// User
	TestCreateUser(t)
	TestFindUserByEmail(t)
	TestUpdateUser(t)
	TestDeleteUser(t)
	TestGetUser(t)
	TestGetAllUsers(t)

	// Annonce
	TestCreateAnnonce(t)
	TestUpdateAnnonce(t)
	TestDeleteAnnonce(t)
	TestGetAnnonce(t)
	TestGetAllAnnonces(t)

	// Association
	TestCreateAssociation(t)
	TestUpdateAssociation(t)
	TestDeleteAssociation(t)
	TestGetAssociation(t)
	TestGetAllAssociations(t)

	// Cat
	TestCreateCat(t)
	TestUpdateCat(t)
	TestDeleteCat(t)
	TestGetCat(t)
	TestGetAllCats(t)

	// Favorites
	TestCreateFavorites(t)
	TestGetFavorites(t)
	TestGetAllFavorites(t)

	// Race
	TestCreateRace(t)
	TestGetAllRaces(t)
	TestUpdateRace(t)
	TestDeleteRace(t)
	TestGetRace(t)

	// Rating
	TestCreateRating(t)
	TestUpdateRating(t)
	TestGetRatings(t)
	TestGetAllRatings(t)
	TestDeleteRatings(t)

	// Room

	destroy, err := database.TestDatabaseDestroy(s.Db)
	if err != nil {
		return
	}
	utils.Logger("debug", "GO-Tests", "Database", destroy)
	utils.Logger("debug", "GO-Tests", "General", "End of Tests")
}

// Auth

// User
var testUser = &models.User{
	Email:    "test@example.com",
	Name:     "Test User",
	Password: "password123",
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := queries.NewMockUserQueries(ctrl)
	mockService.EXPECT().CreateUser(testUser).Return(nil)

	err := mockService.CreateUser(testUser)
	assert.NoError(t, err)
}

func TestFindUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := queries.NewMockUserQueries(ctrl)
	mockService.EXPECT().FindUserByEmail("test@example.com").Return(testUser, nil)

	user, err := mockService.FindUserByEmail("test@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUser.Email, user.Email)
}

func TestUpdateUser(t *testing.T) {}

func TestDeleteUser(t *testing.T) {}

func TestGetUser(t *testing.T) {}

func TestGetAllUsers(t *testing.T) {}

// Annonce
func TestCreateAnnonce(t *testing.T) {}

func TestUpdateAnnonce(t *testing.T) {}

func TestDeleteAnnonce(t *testing.T) {}

func TestGetAnnonce(t *testing.T) {}

func TestGetAllAnnonces(t *testing.T) {}

// Association
func TestCreateAssociation(t *testing.T) {}

func TestUpdateAssociation(t *testing.T) {}

func TestDeleteAssociation(t *testing.T) {}

func TestGetAssociation(t *testing.T) {}

func TestGetAllAssociations(t *testing.T) {}

// Cat
func TestCreateCat(t *testing.T) {}

func TestUpdateCat(t *testing.T) {}

func TestDeleteCat(t *testing.T) {}

func TestGetCat(t *testing.T) {}

func TestGetAllCats(t *testing.T) {}

// Favorites
func TestCreateFavorites(t *testing.T) {}

func TestGetFavorites(t *testing.T) {}

func TestGetAllFavorites(t *testing.T) {}

// Race
func TestCreateRace(t *testing.T) {}

func TestGetAllRaces(t *testing.T) {}

func TestUpdateRace(t *testing.T) {}

func TestDeleteRace(t *testing.T) {}

func TestGetRace(t *testing.T) {}

// Rating
func TestCreateRating(t *testing.T) {}

func TestUpdateRating(t *testing.T) {}

func TestGetRatings(t *testing.T) {}

func TestGetAllRatings(t *testing.T) {}

func TestDeleteRatings(t *testing.T) {}

// Rooms
