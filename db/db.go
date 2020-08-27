package db

import (
	"fmt"
	"log"
	"os"

	"github.com/emanuelcondo/api-todo/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB

//ConnectDB function: Make database connection
func ConnectDB() *gorm.DB {
	//Load environmenatal variables
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")

	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, databasePort, username, databaseName, password)

	_db, err := gorm.Open("postgres", dbURI)
	db = _db
	if err != nil {
		log.Panicf("Connection Error Database: %s\n", err.Error())
	}

	db.AutoMigrate(
		&models.User{},
		&models.Todo{},
	)

	fmt.Println("Successfully connected!", db)

	return db
}

// Returns DB connection
func GetDBConnection() *gorm.DB {
	return db
}
