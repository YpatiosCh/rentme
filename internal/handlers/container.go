package handlers

import (
	"html/template"

	"github.com/YpatiosCh/rentme/internal/services"
)

type HandlerContainer struct {
	homeHandler HomeHandler
	userHandler UserHandler
	authHandler AuthHandler
}

func NewHandlerContainer(services services.Services, tmpl *template.Template) Handlers {
	return &HandlerContainer{
		homeHandler: NewHomeHandler(services, tmpl),
		userHandler: NewUserHandler(services, tmpl),
		authHandler: NewAuthHandler(services, tmpl),
	}
}

func (h *HandlerContainer) Home() HomeHandler {
	return h.homeHandler
}

func (h *HandlerContainer) User() UserHandler {
	return h.userHandler
}

func (h *HandlerContainer) Auth() AuthHandler {
	return h.authHandler
}
