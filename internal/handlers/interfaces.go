package handlers

import "net/http"

type Handlers interface {
	Home() HomeHandler
	User() UserHandler
	Auth() AuthHandler
}

type HomeHandler interface {
	Home(w http.ResponseWriter, r *http.Request)
}

type UserHandler interface {
	GetAllUsers(w http.ResponseWriter, r *http.Request)
}

type AuthHandler interface {
	ShowRegistrationForm(w http.ResponseWriter, r *http.Request)
	CreateSubscription(w http.ResponseWriter, r *http.Request)
	GetStripeConfig(w http.ResponseWriter, r *http.Request)
	CompleteRegistration(w http.ResponseWriter, r *http.Request)
}
