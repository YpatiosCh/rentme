package main

import (
	"log"
	"net/http"
	"os"

	"github.com/YpatiosCh/rentme/internal/config"
	"github.com/YpatiosCh/rentme/internal/database"
	"github.com/YpatiosCh/rentme/internal/models"
	"github.com/YpatiosCh/rentme/internal/repositories"
	"github.com/YpatiosCh/rentme/internal/routes"
	"github.com/YpatiosCh/rentme/internal/services"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// initialize stripe configuration
	config.InitStripe()

	// connect to database
	DB, err = database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	// Get underlying sql.DB for closing
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB:", err)
	}
	defer sqlDB.Close()

	// Run migrations
	err = DB.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// initialize repositories
	repos := repositories.NewRepositoryContainer(DB)

	// create the services dependency
	jwtSecret := os.Getenv("JWT_SECRET")
	dependency := services.Dependency{
		Repos:     repos,
		JwtSecret: jwtSecret,
	}

	// initialize services
	services := services.NewServiceContainer(dependency)

	// var user models.User
	// user.Email = "ypatios@gmail.com"
	// user.FirstName = "Ypatios"
	// user.LastName = "Chanio"
	// _, err = services.Auth().RegisterUser(&user, "password123")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// setup routes
	router := routes.SetupRoutes(services)

	// Start the server
	PORT := os.Getenv("PORT")
	log.Println("Starting server on :", PORT)
	if err := http.ListenAndServe(":"+PORT, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
