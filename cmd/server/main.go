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

	// Validate Stripe configuration
	if err := config.ValidateStripeConfig(); err != nil {
		log.Printf("Warning: Stripe configuration issue: %v", err)
		log.Println("Some Stripe features may not work properly")
	}

	// Initialize stripe configuration
	config.InitStripe()

	// Connect to database
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

	// Initialize repositories
	repos := repositories.NewRepositoryContainer(DB)

	// Create the services dependency
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	dependency := services.Dependency{
		Repos:     repos,
		JwtSecret: jwtSecret,
	}

	// Initialize services
	services := services.NewServiceContainer(dependency)

	// Setup routes
	router := routes.SetupRoutes(services)

	// Start the server
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080" // Default port
	}

	log.Printf("Starting server on port: %s", PORT)
	log.Printf("Base URL: %s", os.Getenv("BASE_URL"))
	log.Printf("Registration available at: %s/register", os.Getenv("BASE_URL"))

	if err := http.ListenAndServe(":"+PORT, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
