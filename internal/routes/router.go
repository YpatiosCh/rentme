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

	// Initialize subscription handler
	subscriptionHandler := handlers.NewSubscriptionHandler(services.Subscription())

	// Basic test page
	mux.HandleFunc("/", testPageHandler)

	// User endpoints
	mux.HandleFunc("/users", handlers.GetAllUsersHandler(services.User()))

	// Subscription endpoints
	mux.HandleFunc("/checkout", subscriptionHandler.CreateCheckoutSessionHandler)
	mux.HandleFunc("/subscription/success", subscriptionHandler.SubscriptionSuccessHandler)
	mux.HandleFunc("/subscription/cancel", subscriptionHandler.SubscriptionCancelHandler)

	return mux
}

func testPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/test.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, nil)
}
