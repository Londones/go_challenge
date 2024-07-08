package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go-challenge/internal/fixtures"
	"go-challenge/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

type Database interface {
	DB() *gorm.DB
}

type Service struct {
	Db *gorm.DB
}

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func New(config *Config) (*Service, error) {
	// Get the root directory of the project.
	var root string
	var err error

	root, err = filepath.Abs("./")
	//root, err = filepath.Abs("../..")

	if err != nil {
		log.Fatal(err)
	}

	// Construct the path to the .env file.
	envPath := filepath.Join(root, ".env")

	// Load the .env file.
	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Variable root: " + root)
		//log.Fatal("Error loading .env file")
	}

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

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/?sslmode=disable", config.Username, config.Password, config.Host, config.Port)
	dbTemp, err := gorm.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("failed to connect to server: %v", err)
	}

	err = createDbIfNotExists(dbTemp, config.Database)
	if err != nil {
		fmt.Printf("failed to create db: %v", err)
	}

	db, err := gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Username, config.Password, config.Host, config.Port, config.Database))
	if err != nil {
		fmt.Printf("failed to connect to database: %v", err)
	}

	err = migrateAllModels(db)
	if err != nil {
		fmt.Printf("failed to migrate models: %v", err)
	}

	// Exécution des fixtures de chats avec un ID utilisateur statique
	cats, err := fixtures.CreateCatFixtures(db, 10)
	if err != nil {
		fmt.Printf("failed to create cat fixtures: %v", err)
	}

	// Exécution des fixtures d'annonces avec les chats existants
	if err := fixtures.CreateAnnonceFixtures(db, cats); err != nil {
		fmt.Printf("failed to create annonce fixtures: %v", err)
	}

	s := &Service{Db: db}

	// print that the database is connected
	fmt.Printf("Connected to database %s\n", config.Database)

	return s, nil
}

func createDbIfNotExists(db *gorm.DB, dbName string) error {
	defer db.Close()
	var count int
	err := db.Raw("SELECT count(*) FROM pg_database WHERE datname = ?", dbName).Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check if db exists: %w", err)
	}
	if count == 0 {
		err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
		if err != nil {
			return fmt.Errorf("failed to create db: %w", err)
		}
	}
	return nil
}

func migrateAllModels(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Annonce{},
		&models.Association{},
		&models.Cats{},
		&models.Favorite{},
		&models.Rating{},
		&models.Roles{},
		&models.User{},
	).Error
	if err != nil {
		fmt.Printf("AutoMigrate error: %v\n", err)
		return err
	}

	// Insert roles
	roles := []models.Roles{
		{Name: models.AdminRole},
		{Name: models.UserRole},
		{Name: models.AssoRole},
	}

	for _, role := range roles {
		var existingRole models.Roles
		if db.Where("name = ?", role.Name).First(&existingRole).RecordNotFound() {
			if err := db.Create(&role).Error; err != nil {
				fmt.Printf("Error creating role: %v\n", err)
				return err
			}
		}
	}

	fmt.Println("Migrated models and inserted roles successfully")
	return nil
}

func (s *Service) DB() *gorm.DB {
	return s.Db
}
