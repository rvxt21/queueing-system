package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rvxt21/queueing-system/pkg/dto"
)

// GetServices godoc
//
//	@Summary		List services
//	@Description	Returns all services, including inactive ones, ordered by name.
//	@Tags			services
//	@Produce		json
//	@Success		200	{array}		dto.ServicesResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/services [get]
func (h *Handler) GetServices(c *gin.Context) {
	ctx := c.Request.Context()

	services, err := h.services.ListServices(ctx)
	if err != nil {
		h.logger.ErrorContext(ctx, "list services", "error", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "internal server error"})

		return
	}

	c.JSON(http.StatusOK, dto.NewServicesResponse(services))
}
