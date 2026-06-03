package order

import (
	"errors"

	orderdto "petshop/internal/dto/order"
	ordersvc "petshop/internal/service/order"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct{ svc ordersvc.AdminService }

func NewAdminOrderHandler(svc ordersvc.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

// ListOrdersAdmin godoc
// @Summary      List all orders (admin)
// @Description  Returns a paginated list of all orders with filters
// @Tags         Admin Orders
// @Produce      json
// @Param        status    query     string  false  "Filter by status"
// @Param        search    query     string  false  "Search by customer name or order ID"
// @Param        date_from query     string  false  "From date (YYYY-MM-DD)"
// @Param        date_to   query     string  false  "To date (YYYY-MM-DD)"
// @Param        page      query     int     false  "Page number (default: 1)"
// @Param        limit     query     int     false  "Items per page (default: 20, max: 100)"
// @Success      200       {object}  response.Response
// @Failure      500       {object}  response.ErrorResponse
// @Router       /admin/orders [get]
func (h *AdminHandler) List(c *gin.Context) {
	f := orderdto.AdminOrderFilter{
		Status:   c.Query("status"),
		Search:   c.Query("search"),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Page:     parseIntQuery(c, "page", 1),
		Limit:    parseIntQuery(c, "limit", 20),
	}

	paginated, err := h.svc.ListOrders(f)
	if err != nil {
		response.InternalError(c, "Failed to fetch orders")
		return
	}
	response.OK(c, "Success", paginated)
}

// GetOrderAdmin godoc
// @Summary      Get order detail (admin)
// @Description  Returns a single order with items (admin, any customer)
// @Tags         Admin Orders
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  response.Response{data=orderdto.AdminOrderResponse}
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /admin/orders/{id} [get]
func (h *AdminHandler) Get(c *gin.Context) {
	o, err := h.svc.GetOrder(c.Param("id"))
	if err != nil {
		if errors.Is(err, ordersvc.ErrOrderNotFound) {
			response.NotFound(c, "Order not found")
			return
		}
		response.InternalError(c, "Failed to fetch order")
		return
	}
	response.OK(c, "Success", o)
}

// UpdateOrderStatus godoc
// @Summary      Update order status (admin)
// @Description  Update the status of an order with valid transitions
// @Tags         Admin Orders
// @Accept       json
// @Produce      json
// @Param        id       path      string                           true  "Order ID"
// @Param        request  body      orderdto.UpdateStatusRequest     true  "New status"
// @Success      200      {object}  response.Response{data=orderdto.AdminOrderResponse}
// @Failure      400      {object}  response.ErrorResponse
// @Failure      404      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /admin/orders/{id}/status [patch]
func (h *AdminHandler) UpdateStatus(c *gin.Context) {
	var req orderdto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	o, err := h.svc.UpdateStatus(c.Param("id"), req)
	if err != nil {
		if errors.Is(err, ordersvc.ErrOrderNotFound) {
			response.NotFound(c, "Order not found")
			return
		}
		if errors.Is(err, ordersvc.ErrInvalidStatusTransition) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "Failed to update order status")
		return
	}
	response.OK(c, "Order status updated", o)
}
