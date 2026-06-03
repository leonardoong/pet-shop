package user

import (
	"errors"
	"strconv"

	userdto "petshop/internal/dto/user"
	usersvc "petshop/internal/service/user"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminCustomerHandler struct{ svc usersvc.AdminCustomerService }

func NewAdminCustomerHandler(svc usersvc.AdminCustomerService) *AdminCustomerHandler {
	return &AdminCustomerHandler{svc: svc}
}

// ListCustomersAdmin godoc
// @Summary      List customers (admin)
// @Description  Returns paginated list of customers with search
// @Tags         Admin Customers
// @Produce      json
// @Param        search  query     string  false  "Search by name, email, or phone"
// @Param        page    query     int     false  "Page number (default: 1)"
// @Param        limit   query     int     false  "Items per page (default: 20)"
// @Success      200     {object}  response.Response
// @Failure      500     {object}  response.ErrorResponse
// @Router       /admin/customers [get]
func (h *AdminCustomerHandler) List(c *gin.Context) {
	f := userdto.AdminCustomerFilter{
		Search: c.Query("search"),
		Page:   parseIntQuery(c, "page", 1),
		Limit:  parseIntQuery(c, "limit", 20),
	}

	paginated, err := h.svc.ListCustomers(f)
	if err != nil {
		response.InternalError(c, "Failed to fetch customers")
		return
	}
	response.OK(c, "Success", paginated)
}

// GetCustomerAdmin godoc
// @Summary      Get customer detail (admin)
// @Description  Returns a single customer by ID
// @Tags         Admin Customers
// @Produce      json
// @Param        id   path      string  true  "Customer ID"
// @Success      200  {object}  response.Response{data=userdto.CustomerResponse}
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /admin/customers/{id} [get]
func (h *AdminCustomerHandler) Get(c *gin.Context) {
	customer, err := h.svc.GetCustomer(c.Param("id"))
	if err != nil {
		if errors.Is(err, usersvc.ErrCustomerNotFound) {
			response.NotFound(c, "Customer not found")
			return
		}
		response.InternalError(c, "Failed to fetch customer")
		return
	}
	response.OK(c, "Success", customer)
}

// ToggleActive godoc
// @Summary      Toggle customer active status
// @Description  Activate or deactivate a customer account
// @Tags         Admin Customers
// @Accept       json
// @Produce      json
// @Param        id       path      string                      true  "Customer ID"
// @Param        request  body      userdto.ToggleActiveRequest  true  "Active status"
// @Success      200      {object}  response.Response
// @Failure      400      {object}  response.ErrorResponse
// @Failure      404      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /admin/customers/{id} [patch]
func (h *AdminCustomerHandler) ToggleActive(c *gin.Context) {
	var req userdto.ToggleActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	customer, err := h.svc.ToggleActive(c.Param("id"), req.IsActive)
	if err != nil {
		if errors.Is(err, usersvc.ErrCustomerNotFound) {
			response.NotFound(c, "Customer not found")
			return
		}
		response.InternalError(c, "Failed to update customer")
		return
	}
	response.OK(c, "Customer updated", customer)
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
