package product

import (
	"errors"
	"strconv"

	productdto "petshop/internal/dto/product"
	productsvc "petshop/internal/service/product"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct{ svc productsvc.InventoryService }

func NewInventoryHandler(svc productsvc.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

// ListInventory godoc
// @Summary      List inventory
// @Description  Returns paginated product inventory list
// @Tags         Admin Inventory
// @Produce      json
// @Param        search      query     string  false  "Search by name or SKU"
// @Param        category_id query     string  false  "Filter by category ID"
// @Param        low_stock   query     bool    false  "Show only low stock items"
// @Param        page        query     int     false  "Page number (default: 1)"
// @Param        limit       query     int     false  "Items per page (default: 20)"
// @Success      200         {object}  response.Response
// @Failure      500         {object}  response.ErrorResponse
// @Router       /admin/inventory [get]
func (h *InventoryHandler) List(c *gin.Context) {
	f := productdto.InventoryFilter{
		Search:     c.Query("search"),
		CategoryID: c.Query("category_id"),
		Page:       parseIntQuery(c, "page", 1),
		Limit:      parseIntQuery(c, "limit", 20),
	}

	if v := c.Query("low_stock"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			f.LowStock = &b
		}
	}

	paginated, err := h.svc.ListInventory(f)
	if err != nil {
		response.InternalError(c, "Failed to fetch inventory")
		return
	}
	response.OK(c, "Success", paginated)
}

// AdjustStock godoc
// @Summary      Adjust stock
// @Description  Add, subtract, or set product stock
// @Tags         Admin Inventory
// @Accept       json
// @Produce      json
// @Param        productId  path      string                        true  "Product ID"
// @Param        request    body      productdto.AdjustStockRequest true  "Stock adjustment"
// @Success      200        {object}  response.Response{data=productdto.InventoryItem}
// @Failure      400        {object}  response.ErrorResponse
// @Failure      404        {object}  response.ErrorResponse
// @Failure      500        {object}  response.ErrorResponse
// @Router       /admin/inventory/{productId} [patch]
func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	var req productdto.AdjustStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	item, err := h.svc.AdjustStock(c.Param("productId"), req)
	if err != nil {
		if errors.Is(err, productsvc.ErrProductNotFound) {
			response.NotFound(c, "Product not found")
			return
		}
		if errors.Is(err, productsvc.ErrInvalidOperation) {
			response.BadRequest(c, "Invalid operation, must be add, subtract, or set")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, "Stock adjusted", item)
}
