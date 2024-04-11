package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "password123"
	hash, err := HashPassword(password)

	assert.Nil(t, err)
	assert.NotEmpty(t, hash)
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password123"
	hash, _ := HashPassword(password)

	assert.True(t, CheckPasswordHash(password, hash))
}
