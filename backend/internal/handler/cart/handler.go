package cart

import (
	"errors"

	cartdto "petshop/internal/dto/cart"
	"petshop/internal/middleware"
	cartsvc "petshop/internal/service/cart"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct{ svc cartsvc.Service }

func NewHandler(svc cartsvc.Service) *Handler { return &Handler{svc: svc} }

// GetCart godoc
// @Summary      Get cart
// @Description  Returns the authenticated customer's cart with all items
// @Tags         Cart
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /customer/cart [get]
func (h *Handler) Get(c *gin.Context) {
	customerID := middleware.MustGetCustomerID(c)
	cart, err := h.svc.GetCart(customerID)
	if err != nil {
		response.InternalError(c, "Failed to fetch cart")
		return
	}
	response.OK(c, "Success", cart)
}

// AddItem godoc
// @Summary      Add item to cart
// @Description  Add a product to the cart. If the product already exists, quantity is incremented.
// @Tags         Cart
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      cartdto.AddItemRequest  true  "Item payload"
// @Success      200      {object}  response.Response
// @Failure      400      {object}  response.ErrorResponse
// @Failure      401      {object}  response.ErrorResponse
// @Failure      404      {object}  response.ErrorResponse  "Product not found"
// @Failure      409      {object}  response.ErrorResponse  "Insufficient stock"
// @Failure      500      {object}  response.ErrorResponse
// @Router       /customer/cart/items [post]
func (h *Handler) AddItem(c *gin.Context) {
	var req cartdto.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	customerID := middleware.MustGetCustomerID(c)
	cart, err := h.svc.AddItem(customerID, req)
	if err != nil {
		switch {
		case errors.Is(err, cartsvc.ErrProductNotFound):
			response.NotFound(c, "Product not found")
		case errors.Is(err, cartsvc.ErrInsufficientStock):
			response.Conflict(c, "Insufficient stock")
		default:
			response.InternalError(c, "Failed to add item")
		}
		return
	}
	response.OK(c, "Item added", cart)
}

// UpdateItem godoc
// @Summary      Update cart item quantity
// @Description  Set the quantity of a specific product in the cart
// @Tags         Cart
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        productId  path      string                  true  "Product ID (UUID)"
// @Param        request    body      cartdto.UpdateItemRequest  true  "Quantity payload"
// @Success      200        {object}  response.Response
// @Failure      400        {object}  response.ErrorResponse
// @Failure      401        {object}  response.ErrorResponse
// @Failure      404        {object}  response.ErrorResponse
// @Failure      409        {object}  response.ErrorResponse
// @Failure      500        {object}  response.ErrorResponse
// @Router       /customer/cart/items/{productId} [put]
func (h *Handler) UpdateItem(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		response.BadRequest(c, "Invalid product ID")
		return
	}

	var req cartdto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	customerID := middleware.MustGetCustomerID(c)
	cart, err := h.svc.UpdateItem(customerID, productID, req)
	if err != nil {
		switch {
		case errors.Is(err, cartsvc.ErrProductNotFound), errors.Is(err, cartsvc.ErrItemNotFound):
			response.NotFound(c, "Item not found")
		case errors.Is(err, cartsvc.ErrInsufficientStock):
			response.Conflict(c, "Insufficient stock")
		default:
			response.InternalError(c, "Failed to update item")
		}
		return
	}
	response.OK(c, "Cart updated", cart)
}

// RemoveItem godoc
// @Summary      Remove item from cart
// @Description  Remove a specific product from the cart
// @Tags         Cart
// @Security     BearerAuth
// @Produce      json
// @Param        productId  path      string  true  "Product ID (UUID)"
// @Success      200        {object}  response.Response
// @Failure      400        {object}  response.ErrorResponse
// @Failure      401        {object}  response.ErrorResponse
// @Failure      500        {object}  response.ErrorResponse
// @Router       /customer/cart/items/{productId} [delete]
func (h *Handler) RemoveItem(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		response.BadRequest(c, "Invalid product ID")
		return
	}

	customerID := middleware.MustGetCustomerID(c)
	if err := h.svc.RemoveItem(customerID, productID); err != nil {
		response.InternalError(c, "Failed to remove item")
		return
	}
	response.OK(c, "Item removed", nil)
}
