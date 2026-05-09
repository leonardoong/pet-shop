package user

import (
	"errors"

	userdto "petshop/internal/dto/user"
	"petshop/internal/middleware"
	usersvc "petshop/internal/service/user"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddressHandler struct{ svc usersvc.AddressService }

func NewAddressHandler(svc usersvc.AddressService) *AddressHandler {
	return &AddressHandler{svc: svc}
}

// ListAddresses godoc
// @Summary      List addresses
// @Description  Returns all shipping addresses for the authenticated customer
// @Tags         Addresses
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /customer/addresses [get]
func (h *AddressHandler) List(c *gin.Context) {
	customerID := middleware.MustGetCustomerID(c)
	addrs, err := h.svc.ListAddresses(customerID)
	if err != nil {
		response.InternalError(c, "Failed to fetch addresses")
		return
	}
	response.OK(c, "Success", addrs)
}

// CreateAddress godoc
// @Summary      Create address
// @Description  Add a new shipping address for the authenticated customer
// @Tags         Addresses
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      userdto.AddressCreateRequest  true  "Address payload"
// @Success      201      {object}  response.Response
// @Failure      400      {object}  response.ErrorResponse
// @Failure      401      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /customer/addresses [post]
func (h *AddressHandler) Create(c *gin.Context) {
	var req userdto.AddressCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	customerID := middleware.MustGetCustomerID(c)
	addr, err := h.svc.CreateAddress(customerID, req)
	if err != nil {
		response.InternalError(c, "Failed to create address")
		return
	}
	response.Created(c, "Address created", addr)
}

// UpdateAddress godoc
// @Summary      Update address
// @Description  Update an existing shipping address (partial update)
// @Tags         Addresses
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                        true  "Address ID (UUID)"
// @Param        request  body      userdto.AddressUpdateRequest  true  "Fields to update"
// @Success      200      {object}  response.Response
// @Failure      400      {object}  response.ErrorResponse
// @Failure      401      {object}  response.ErrorResponse
// @Failure      404      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /customer/addresses/{id} [put]
func (h *AddressHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid address ID")
		return
	}

	var req userdto.AddressUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	customerID := middleware.MustGetCustomerID(c)
	addr, err := h.svc.UpdateAddress(id, customerID, req)
	if err != nil {
		if errors.Is(err, usersvc.ErrAddressNotFound) {
			response.NotFound(c, "Address not found")
			return
		}
		response.InternalError(c, "Failed to update address")
		return
	}
	response.OK(c, "Address updated", addr)
}

// DeleteAddress godoc
// @Summary      Delete address
// @Description  Remove a shipping address belonging to the authenticated customer
// @Tags         Addresses
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "Address ID (UUID)"
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /customer/addresses/{id} [delete]
func (h *AddressHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid address ID")
		return
	}

	customerID := middleware.MustGetCustomerID(c)
	if err := h.svc.DeleteAddress(id, customerID); err != nil {
		if errors.Is(err, usersvc.ErrAddressNotFound) {
			response.NotFound(c, "Address not found")
			return
		}
		response.InternalError(c, "Failed to delete address")
		return
	}
	response.OK(c, "Address deleted", nil)
}

// SetDefaultAddress godoc
// @Summary      Set default address
// @Description  Mark an address as the default shipping address
// @Tags         Addresses
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "Address ID (UUID)"
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /customer/addresses/{id}/default [patch]
func (h *AddressHandler) SetDefault(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid address ID")
		return
	}

	customerID := middleware.MustGetCustomerID(c)
	addr, err := h.svc.SetDefaultAddress(id, customerID)
	if err != nil {
		if errors.Is(err, usersvc.ErrAddressNotFound) {
			response.NotFound(c, "Address not found")
			return
		}
		response.InternalError(c, "Failed to set default address")
		return
	}
	response.OK(c, "Default address updated", addr)
}
