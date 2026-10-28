package handlers

import (
	"log/slog"

	"github.com/rvxt21/queueing-system/internal/service/hub"
)

type Handler struct {
	logger   *slog.Logger
	services ServicesService
	hub      *hub.Hub
}

func NewHandler(logger *slog.Logger, services ServicesService, h *hub.Hub) *Handler {
	return &Handler{
		logger:   logger,
		services: services,
		hub:      h,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}
