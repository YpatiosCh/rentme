package routes

import (
	"html/template"
	"net/http"

	"github.com/YpatiosCh/rentme/internal/handlers"
	"github.com/YpatiosCh/rentme/internal/services"
)

func SetupRoutes(services services.Services) http.Handler {
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve HTML templates
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	// initialize handlers with services and templates
	handlers := handlers.NewHandlerContainer(services, tmpl)

	// Home endpoint
	mux.HandleFunc("/", handlers.Home().Home)

	// Authentication endpoints
	mux.HandleFunc("/register", handlers.Auth().ShowRegistrationForm)

	// User endpoints
	mux.HandleFunc("/users", handlers.User().GetAllUsers)

	return mux
}
