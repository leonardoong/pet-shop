package dashboard

import (
	dashsvc "petshop/internal/service/dashboard"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct{ svc dashsvc.Service }

func NewHandler(svc dashsvc.Service) *Handler {
	return &Handler{svc: svc}
}

// GetStats godoc
// @Summary      Dashboard analytics
// @Description  Returns aggregated dashboard statistics
// @Tags         Admin Dashboard
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrorResponse
// @Router       /admin/dashboard [get]
func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats()
	if err != nil {
		response.InternalError(c, "Failed to fetch dashboard stats")
		return
	}
	response.OK(c, "Success", stats)
}
