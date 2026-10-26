package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rvxt21/queueing-system/internal/domain"
	"github.com/rvxt21/queueing-system/pkg/dto"
)

// RequestService godoc
//
//	@Summary		Request a service
//	@Description	Customer chooses a service and receives a ticket number, e.g. PN-042. Numbers restart from 001 every day.
//	@Tags			tickets
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RequestServiceRequest	true	"Service to request"
//	@Success		201		{object}	dto.TicketResponse
//	@Failure		400		{object}	errorResponse	"Invalid request body"
//	@Failure		404		{object}	errorResponse	"Service not found"
//	@Failure		409		{object}	errorResponse	"Service is not active"
//	@Failure		500		{object}	errorResponse
//	@Router			/tickets [post]
func (h *Handler) RequestService(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.RequestServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid request body: service_code is required"})

		return
	}

	ticket, err := h.services.RequestService(ctx, req.ServiceCode)

	switch {
	case err == nil:
		c.JSON(http.StatusCreated, dto.NewTicketResponse(*ticket))
	case errors.Is(err, domain.ErrServiceNotFound):
		c.JSON(http.StatusNotFound, errorResponse{Error: "service not found"})
	case errors.Is(err, domain.ErrServiceInactive):
		c.JSON(http.StatusConflict, errorResponse{Error: "service is not active"})
	default:
		h.logger.ErrorContext(ctx, "request service", "error", err, "service_code", req.ServiceCode)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}
