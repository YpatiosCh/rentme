package routes

import (
	"html/template"
	"net/http"

	"github.com/YpatiosCh/rentme/internal/handlers"
	"github.com/YpatiosCh/rentme/internal/middleware"
	"github.com/YpatiosCh/rentme/internal/services"
)

func SetupRoutes(services services.Services) http.Handler {
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve HTML templates
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	// initialize middleware
	middleware := middleware.NewMiddleware(services)

	// initialize handlers with services and templates
	handlers := handlers.NewHandlerContainer(services, tmpl)

	// Home endpoint
	mux.HandleFunc("/", handlers.Home().Home)

	// Authentication endpoints
	mux.HandleFunc("/register", handlers.Auth().ShowRegistrationForm)
	mux.HandleFunc("/logout", handlers.Auth().Logout)

	mux.HandleFunc("/add-item", middleware.RequireUser(handlers.Item().CreateItemForm))

	// Payment endpoints
	mux.HandleFunc("/create-subscription", handlers.Auth().CreateSubscription)
	mux.HandleFunc("/complete-registration", handlers.Auth().CompleteRegistration)
	mux.HandleFunc("/stripe/config", handlers.Auth().GetStripeConfig)
	mux.HandleFunc("/webhook/stripe", handlers.Auth().StripeWebhook)

	// User endpoints
	mux.HandleFunc("/users", handlers.User().GetAllUsers)

	return middleware.AuthMiddleware(mux)
}
