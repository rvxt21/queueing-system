package handlers

import "log/slog"

type Handler struct {
	logger   *slog.Logger
	services ServicesService
}

func NewHandler(logger *slog.Logger, services ServicesService) *Handler {
	return &Handler{
		logger:   logger,
		services: services,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}
