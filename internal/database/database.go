package database

import (
	"fmt"
	"os"

	"go-challenge/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	FindUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
}

type service struct {
	db *gorm.DB
}

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func New(config *Config) (*service, error) {
	if config.Username == "" {
		config.Username = os.Getenv("DB_USERNAME")
	}
	if config.Password == "" {
		config.Password = os.Getenv("DB_PASSWORD")
	}
	if config.Host == "" {
		config.Host = os.Getenv("DB_HOST")
	}
	if config.Port == "" {
		config.Port = os.Getenv("DB_PORT")
	}
	if config.Database == "" {
		config.Database = os.Getenv("DB_DATABASE")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	s := &service{db: db}
	return s, nil
}

func (s *service) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}
