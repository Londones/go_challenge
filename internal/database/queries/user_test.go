package queries

import (
	queries "go-challenge/internal/database/queries/mocks"
	"go-challenge/internal/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := queries.NewMockUserQueries(ctrl)
	mockService.EXPECT().CreateUser(&models.User{Email: "test@example.com"}).Return(nil)

	err := mockService.CreateUser(&models.User{Email: "test@example.com"})
	assert.NoError(t, err)
}

func TestFindUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := queries.NewMockUserQueries(ctrl)
	mockService.EXPECT().FindUserByEmail("test@example.com").Return(&models.User{Email: "test@example.com"}, nil)

	user, err := mockService.FindUserByEmail("test@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
}
