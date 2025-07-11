package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/YpatiosCh/rentme/internal/models"
	"github.com/YpatiosCh/rentme/internal/services"
)

type contextKey string

const UserKey contextKey = "user"

type Middleware struct {
	services services.Services
}

func NewMiddleware(services services.Services) *Middleware {
	return &Middleware{
		services: services,
	}
}

// RequireUser checks for valid JWT token and loads full user
func (m *Middleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Try to get token from cookie first
		token := getTokenFromCookie(r)

		// If no cookie, try Authorization header
		if token == "" {
			token = getTokenFromHeader(r)
		}

		// If no token found
		if token == "" {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		// Validate token and get user ID
		userID, err := m.services.Auth().ValidateToken(token)
		if err != nil {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		// Get full user from database
		user, userErr := m.services.User().GetUserByID(*userID)
		if userErr != nil {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		// Add full user to context
		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// OptionalAuthMiddleware adds user info to context if token exists, but doesn't require it
func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get token
		token := getTokenFromCookie(r)
		if token == "" {
			token = getTokenFromHeader(r)
		}

		// If token exists and is valid, add user to context
		if token != "" {
			if userID, err := m.services.Auth().ValidateToken(token); err == nil {
				if user, userErr := m.services.User().GetUserByID(*userID); userErr == nil {
					ctx := context.WithValue(r.Context(), UserKey, user)
					r = r.WithContext(ctx)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// Helper functions
func getTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func getTokenFromHeader(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// GetUserFromContext extracts full user from request context
func GetUserFromContext(r *http.Request) (*models.User, bool) {
	user, ok := r.Context().Value(UserKey).(*models.User)
	if !ok {
		return nil, false
	}
	return user, true
}
