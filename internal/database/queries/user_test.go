package queries

import (
	queries "go-challenge/internal/database/queries/mocks"
	"go-challenge/internal/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
