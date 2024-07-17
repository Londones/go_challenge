package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go-challenge/internal/fixtures"
	"go-challenge/internal/models"
	"go-challenge/internal/utils"

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
	Env      string
}

func New(config *Config) (*Service, error) {
	var db *gorm.DB

	config.Env = os.Getenv("APP_ENV")
	// Get the root directory of the project.
	var root string
	var err error

	if config.Env == "local" {

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
		fmt.Println("Try to connect with:", connStr)
		dbTemp, err := gorm.Open("postgres", connStr)
		if err != nil {
			fmt.Printf("failed to connect to server: %v", err)
		}

		err = createDbIfNotExists(dbTemp, config.Database)
		if err != nil {
			fmt.Printf("failed to create db: %v", err)
		}

		db, err = gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Username, config.Password, config.Host, config.Port, config.Database))
		if err != nil {
			fmt.Printf("failed to connect to database: %v", err)
		}
	} else {
		db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Printf("failed to connect to database: %v", err)
		}
		config.Database = os.Getenv("DATABASE_URL")
	}

	err = migrateAllModels(db)
	if err != nil {
		fmt.Printf("failed to migrate models: %v", err)
	}

	// Get the USER role
	/*var userRole models.Roles
	if err := db.Where("name = ?", models.UserRole).First(&userRole).Error; err != nil {
		fmt.Printf("failed to find user role: %v", err)
	}

	// Create 5 races
	err = fixtures.CreateRaceFixture(db)
	if err != nil {
		fmt.Printf("failed to create race fixture: %v", err)
	}

	// Create 5 users
	users, err := fixtures.CreateUserFixtures(db, 5, &userRole)
	if err != nil {
		fmt.Printf("failed to create user fixtures: %v", err)
	}

	// For each user, create 5 cats and 5 corresponding annonces
	for _, user := range users {
		cats, err := fixtures.CreateCatFixturesForUser(db, 5, user.ID)
		if err != nil {
			fmt.Printf("failed to create cat fixtures for user %s: %v", user.ID, err)
		}

		if err := fixtures.CreateAnnonceFixtures(db, cats); err != nil {
			fmt.Printf("failed to create annonce fixtures for user %s: %v", user.ID, err)
		}
	}

	// Création des fixtures pour les évaluations
	staticUserID := "38f5ca5d-0c87-425f-97fe-c84c3ee0997c"
	staticAuthorID := "5a7a8b69-6f8d-4818-ac15-b6a83b4fe518"
	err = fixtures.CreateRatingFixtures(db, staticUserID, staticAuthorID, 5)
	if err != nil {
		fmt.Printf("failed to create rating fixtures: %v", err)
	}*/

	//asso, err := fixtures.CreateAssociationFixtures(db, users[1].ID)
	//if err != nil {
	//	fmt.Printf("failed to create association with owner: %v", users[1].ID)
	//}
	//
	//var allUser = &users
	//
	//_, err = fixtures.AssociationAddMembers(db, asso, allUser)
	//if err != nil {
	//	fmt.Printf("failed to add member to association %c : %v", asso.ID, err)
	//}

	s := &Service{Db: db}

	// Print that the database is connected
	fmt.Printf("Connected to database %s\n", config.Database)

	return s, nil

}

// TestDatabaseInit DATABASE FOR TESTS
func TestDatabaseInit() (*Service, error) {

	var config Config
	var db *gorm.DB

	config.Env = os.Getenv("APP_ENV")
	// Get the root directory of the project.
	var err error

	config.Username = "postgres"
	config.Password = "postgres"
	config.Host = "postgres"
	config.Port = "5432"
	config.Database = "go_purrfectmatch_test"
	config.Env = "local"

	//config.Username = "macbook"
	//config.Password = "postgres"
	//config.Host = "localhost"
	//config.Port = "5432"
	//config.Database = "go_purrfectmatch_test"
	//config.Env = "local"

	if config.Env == "local" {

		fmt.Println("Config Username:", config.Username)
		fmt.Println("Config Password:", config.Password)
		fmt.Println("Config Host:", config.Host)
		fmt.Println("Config Port:", config.Port)
		fmt.Println("Config Database:", config.Database)
		fmt.Println("Config Env:", config.Env)

		//connStr, _ := fmt.Printf("postgres://", config.Username, ":", config.Password, "@", config.Host, ":", config.Port, "/?sslmode=disable")
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/?sslmode=disable", config.Username, config.Password, config.Host, config.Port)
		fmt.Println("Try to connect with:", connStr)
		dbTemp, err := gorm.Open("postgres", connStr)
		fmt.Printf("DBTemp: %v\n", dbTemp)
		if err != nil {
			fmt.Printf("failed to connect to server: %v", err)
		}
		//
		//var count int
		//var dbFounded = dbTemp.Raw("SELECT count(*) FROM pg_database WHERE datname = ?", config.Database)
		//fmt.Println("Founded DB is:", dbFounded.Count(&count))
		//fmt.Println("Nombre de DB:", count)
		//err = dbTemp.Raw("SELECT count(*) FROM pg_database WHERE datname = ?", config.Database).Count(&count).Error
		//if err != nil {
		//	fmt.Errorf("failed to check if db exists: %w", err)
		//}
		//if count == 0 {
		//	err = dbTemp.Exec(fmt.Sprintf("CREATE DATABASE %s", config.Database)).Error
		//	if err != nil {
		//		fmt.Errorf("failed to create db: %w", err)
		//	}
		//}

		err = dbTemp.Exec(fmt.Sprintf("CREATE DATABASE %s", config.Database)).Error
		fmt.Println(err)
		if err != nil {
			fmt.Println("failed to create db:", err)
		}

		db, err = gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Username, config.Password, config.Host, config.Port, config.Database))
		if err != nil {
			fmt.Printf("failed to connect to database: %v", err)
		}
	} else {
		db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Printf("failed to connect to database: %v", err)
		}
		config.Database = os.Getenv("DATABASE_URL")
	}

	err = migrateAllModels(db)
	if err != nil {
		fmt.Printf("failed to migrate models: %v", err)
	}

	// Get the USER role
	var userRole models.Roles
	if err := db.Where("name = ?", models.UserRole).First(&userRole).Error; err != nil {
		fmt.Printf("failed to find user role: %v", err)
	}

	users, err := fixtures.CreateUserFixtures(db, 5, &userRole)
	if err != nil {
		fmt.Printf("failed to create user fixtures: %v", err)
	}
	// Create 5 races
	err = fixtures.CreateRaceFixture(db)
	if err != nil {
		fmt.Printf("failed to create race fixture: %v", err)
	}
	// For each user, create 5 cats and 5 corresponding annonces
	for _, user := range users {
		cats, err := fixtures.CreateCatFixturesForUser(db, 5, user.ID)
		if err != nil {
			fmt.Printf("failed to create cat fixtures for user %s: %v", user.ID, err)
		}
		if err := fixtures.CreateAnnonceFixtures(db, cats); err != nil {
			fmt.Printf("failed to create annonce fixtures for user %s: %v", user.ID, err)
		}
	}

	//asso, err := fixtures.CreateAssociationFixtures(db, users[1].ID)
	//if err != nil {
	//	fmt.Printf("failed to create association with owner: %v", users[1].ID)
	//}
	//
	//var allUser = &users
	//
	//_, err = fixtures.AssociationAddMembers(db, asso, allUser)
	//if err != nil {
	//	fmt.Printf("failed to add member to association %c : %v", asso.ID, err)
	//}

	s := &Service{Db: db}

	// Print that the database is connected
	fmt.Printf("Connected to database %s\n", config.Database)

	return s, nil

}

func TestDatabaseDestroy(db *gorm.DB) (string, error) {
	fmt.Println("attemp destroy")
	db.Close()

	connection, err := sql.Open("postgres", "user=macbook")
	_, err = connection.Exec("DROP DATABASE go_purrfectmatch_test")

	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return "Test Database Destroyed", nil
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
		&models.Races{},
		&models.Favorite{},
		&models.Rating{},
		&models.Roles{},
		&models.User{},
		&models.Message{},
		&models.Room{},
	).Error
	if err != nil {
		utils.Logger("debug", "AutoMigrate:", "Failed to migrate models", fmt.Sprintf("Error: %v", err))
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
				utils.Logger("debug", "Create Role:", "Failed to create role", fmt.Sprintf("Error: %v", err))
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
