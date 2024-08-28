package handlers

import "mywebsite.tv/name/cmd/service"

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}
