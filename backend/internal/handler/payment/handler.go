package payment

import (
	"time"

	ordermodel "petshop/internal/model/order"
	orderrepo "petshop/internal/repository/order"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	orderRepo  orderrepo.Repository
	adminRepo  orderrepo.AdminRepository
}

func NewHandler(orderRepo orderrepo.Repository, adminRepo orderrepo.AdminRepository) *Handler {
	return &Handler{orderRepo: orderRepo, adminRepo: adminRepo}
}

// MockPay godoc
// @Summary      Mock payment completion
// @Description  Simulates a successful payment. Sets order status to confirmed.
// @Tags         Payment
// @Produce      json
// @Param        order  query  string  true  "Order ID"
// @Success      200    {object} response.Response
// @Failure      404    {object} response.ErrorResponse
// @Router       /payment/mock/pay [get]
func (h *Handler) MockPay(c *gin.Context) {
	orderID := c.Query("order")
	if orderID == "" {
		response.BadRequest(c, "order query param required")
		return
	}

	o, err := h.adminRepo.FindByIDAdmin(orderID)
	if err != nil || o == nil {
		response.NotFound(c, "Order not found")
		return
	}

	now := time.Now()
	o.PaymentStatus = ordermodel.PaymentStatusSettlement
	o.Status = ordermodel.StatusConfirmed
	o.PaidAt = &now

	if err := h.orderRepo.UpdateOrder(o); err != nil {
		response.InternalError(c, "Failed to update order")
		return
	}

	response.OK(c, "Payment confirmed", nil)
}
