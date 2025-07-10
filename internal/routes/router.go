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
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.Auth().ShowRegistrationForm(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Payment endpoints
	mux.HandleFunc("/create-subscription", handlers.Auth().CreateSubscription)
	mux.HandleFunc("/complete-registration", handlers.Auth().CompleteRegistration)
	mux.HandleFunc("/stripe/config", handlers.Auth().GetStripeConfig)

	// User endpoints
	mux.HandleFunc("/users", handlers.User().GetAllUsers)

	return mux
}
