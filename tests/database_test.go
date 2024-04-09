package tests

import (
	"testing"

	"go-challenge/internal/database"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Call the New function
	service, err := database.New(&database.Config{
		Database: "testgochallenge",
	})

	// Assert that the error is nil
	assert.Nil(t, err)

	// Assert that the returned service is not nil
	assert.NotNil(t, service)

}
