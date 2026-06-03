package order

import (
	"errors"
	"strconv"

	orderdto "petshop/internal/dto/order"
	"petshop/internal/middleware"
	ordersvc "petshop/internal/service/order"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct{ svc ordersvc.Service }

func NewHandler(svc ordersvc.Service) *Handler { return &Handler{svc: svc} }

// Checkout godoc
// @Summary      Place an order
// @Description  Create an order from the current cart. Decrements stock atomically within a transaction.
// @Tags         Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      orderdto.CreateRequest  true  "Checkout payload"
// @Success      201      {object}  response.Response
// @Failure      400      {object}  response.ErrorResponse
// @Failure      401      {object}  response.ErrorResponse
// @Failure      404      {object}  response.ErrorResponse
// @Failure      409      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /customer/orders [post]
func (h *Handler) Checkout(c *gin.Context) {
	var req orderdto.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	customerID := middleware.MustGetCustomerID(c)
	resp, err := h.svc.Checkout(customerID, req)
	if err != nil {
		switch {
		case errors.Is(err, ordersvc.ErrEmptyCart):
			response.BadRequest(c, "Cart is empty")
		case errors.Is(err, ordersvc.ErrAddressNotFound):
			response.NotFound(c, "Address not found")
		case errors.Is(err, ordersvc.ErrInsufficientStock):
			response.Conflict(c, "One or more items have insufficient stock")
		default:
			response.InternalError(c, "Checkout failed")
		}
		return
	}
	response.Created(c, "Order placed", resp)
}

// ListOrders godoc
// @Summary      List orders
// @Description  Returns paginated order history for the authenticated customer
// @Tags         Orders
// @Security     BearerAuth
// @Produce      json
// @Param        status  query     string  false  "Filter by status"
// @Param        page    query     int     false  "Page number (default: 1)"
// @Param        limit   query     int     false  "Items per page (default: 20)"
// @Success      200     {object}  response.Response{data=OrdersPage}
// @Failure      401     {object}  response.ErrorResponse
// @Failure      500     {object}  response.ErrorResponse
// @Router       /customer/orders [get]
func (h *Handler) List(c *gin.Context) {
	f := orderdto.ListFilter{
		Status: c.Query("status"),
		Page:   parseIntQuery(c, "page", 1),
		Limit:  parseIntQuery(c, "limit", 20),
	}

	customerID := middleware.MustGetCustomerID(c)
	paginated, err := h.svc.ListOrders(customerID, f)
	if err != nil {
		response.InternalError(c, "Failed to fetch orders")
		return
	}
	response.OK(c, "Success", paginated)
}

// GetOrderByID godoc
// @Summary      Get order detail
// @Description  Returns a single order with items for the authenticated customer
// @Tags         Orders
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "Order ID (UUID)"
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /customer/orders/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid order ID")
		return
	}

	customerID := middleware.MustGetCustomerID(c)
	o, err := h.svc.GetOrderByID(id, customerID)
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

func parseIntQuery(c *gin.Context, key string, defaultVal int) int {
	v := c.Query(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return defaultVal
	}
	return n
}
