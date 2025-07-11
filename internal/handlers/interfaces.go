package handlers

import "net/http"

type Handlers interface {
	Home() HomeHandler
	User() UserHandler
	Auth() AuthHandler
	Item() ItemHandler
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
	StripeWebhook(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type ItemHandler interface {
	CreateItemForm(w http.ResponseWriter, r *http.Request)
}
