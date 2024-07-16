package handlers

import (
	"go-challenge/internal/database"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	s, err := database.TestDatabaseInit()
	if err != nil {
		return
	}

	// Annonce

	ret := m.Run()

	database.TestDatabaseDestroy(s.Db)
	os.Exit(ret)
}
