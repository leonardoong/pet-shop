package product

import (
	"errors"
	"strconv"

	productdto "petshop/internal/dto/product"
	productsvc "petshop/internal/service/product"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminProductHandler struct{ svc productsvc.AdminProductService }

func NewAdminProductHandler(svc productsvc.AdminProductService) *AdminProductHandler {
	return &AdminProductHandler{svc: svc}
}

// CreateProduct godoc
// @Summary      Create product
// @Description  Create a new product (admin)
// @Tags         Admin Products
// @Accept       json
// @Produce      json
// @Param        request  body      productdto.CreateProductRequest  true  "Product payload"
// @Success      201      {object}  response.Response{data=productdto.AdminProductResponse}
// @Failure      400      {object}  response.ErrorResponse
// @Failure      409      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /admin/products [post]
func (h *AdminProductHandler) Create(c *gin.Context) {
	var req productdto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	p, err := h.svc.CreateProduct(req)
	if err != nil {
		if errors.Is(err, productsvc.ErrSKUExists) {
			response.Conflict(c, "SKU already in use")
			return
		}
		if errors.Is(err, productsvc.ErrCategoryNotFound) {
			response.BadRequest(c, "Category not found")
			return
		}
		response.InternalError(c, "Failed to create product")
		return
	}

	response.Created(c, "Product created", p)
}

// UpdateProduct godoc
// @Summary      Update product
// @Description  Update an existing product (admin)
// @Tags         Admin Products
// @Accept       json
// @Produce      json
// @Param        id       path      string                           true  "Product ID"
// @Param        request  body      productdto.UpdateProductRequest  true  "Product payload"
// @Success      200      {object}  response.Response{data=productdto.AdminProductResponse}
// @Failure      400      {object}  response.ErrorResponse
// @Failure      404      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /admin/products/{id} [put]
func (h *AdminProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req productdto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	p, err := h.svc.UpdateProduct(id, req)
	if err != nil {
		if errors.Is(err, productsvc.ErrProductNotFound) {
			response.NotFound(c, "Product not found")
			return
		}
		if errors.Is(err, productsvc.ErrSKUExists) {
			response.Conflict(c, "SKU already in use")
			return
		}
		if errors.Is(err, productsvc.ErrCategoryNotFound) {
			response.BadRequest(c, "Category not found")
			return
		}
		response.InternalError(c, "Failed to update product")
		return
	}

	response.OK(c, "Product updated", p)
}

// DeleteProduct godoc
// @Summary      Delete product
// @Description  Soft-delete a product (set as inactive)
// @Tags         Admin Products
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  response.Response
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /admin/products/{id} [delete]
func (h *AdminProductHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteProduct(c.Param("id")); err != nil {
		if errors.Is(err, productsvc.ErrProductNotFound) {
			response.NotFound(c, "Product not found")
			return
		}
		response.InternalError(c, "Failed to delete product")
		return
	}
	response.OK(c, "Product deactivated", nil)
}

// GetProduct godoc
// @Summary      Get product
// @Description  Get a single product by ID (admin, includes inactive)
// @Tags         Admin Products
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  response.Response{data=productdto.AdminProductResponse}
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /admin/products/{id} [get]
func (h *AdminProductHandler) Get(c *gin.Context) {
	p, err := h.svc.GetProduct(c.Param("id"))
	if err != nil {
		if errors.Is(err, productsvc.ErrProductNotFound) {
			response.NotFound(c, "Product not found")
			return
		}
		response.InternalError(c, "Failed to fetch product")
		return
	}
	response.OK(c, "Success", p)
}

// ListProducts godoc
// @Summary      List products (admin)
// @Description  Returns a paginated list of all products with filters
// @Tags         Admin Products
// @Produce      json
// @Param        search      query     string  false  "Search by name or SKU"
// @Param        category_id query     string  false  "Filter by category ID"
// @Param        is_active   query     bool    false  "Filter active/inactive"
// @Param        sort        query     string  false  "Sort: newest|oldest|price_asc|price_desc|stock_asc|stock_desc"
// @Param        page        query     int     false  "Page number (default: 1)"
// @Param        limit       query     int     false  "Items per page (default: 20, max: 100)"
// @Success      200         {object}  response.Response
// @Failure      500         {object}  response.ErrorResponse
// @Router       /admin/products [get]
func (h *AdminProductHandler) List(c *gin.Context) {
	f := productdto.AdminProductFilter{
		Search:     c.Query("search"),
		CategoryID: c.Query("category_id"),
		Sort:       c.Query("sort"),
		Page:       parseIntQuery(c, "page", 1),
		Limit:      parseIntQuery(c, "limit", 20),
	}

	if v := c.Query("is_active"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			f.IsActive = &b
		}
	}

	paginated, err := h.svc.ListProducts(f)
	if err != nil {
		response.InternalError(c, "Failed to fetch products")
		return
	}
	response.OK(c, "Success", paginated)
}

// --- Admin Category Handler ---

type AdminCategoryHandler struct{ svc productsvc.AdminCategoryService }

func NewAdminCategoryHandler(svc productsvc.AdminCategoryService) *AdminCategoryHandler {
	return &AdminCategoryHandler{svc: svc}
}

// CreateCategory godoc
// @Summary      Create category
// @Description  Create a new product category (admin)
// @Tags         Admin Categories
// @Accept       json
// @Produce      json
// @Param        request  body      productdto.CreateCategoryRequest  true  "Category payload"
// @Success      201      {object}  response.Response
// @Failure      400      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /admin/categories [post]
func (h *AdminCategoryHandler) Create(c *gin.Context) {
	var req productdto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	cat, err := h.svc.CreateCategory(req)
	if err != nil {
		response.InternalError(c, "Failed to create category")
		return
	}

	response.Created(c, "Category created", cat)
}

// UpdateCategory godoc
// @Summary      Update category
// @Description  Update an existing category (admin)
// @Tags         Admin Categories
// @Accept       json
// @Produce      json
// @Param        id       path      string                           true  "Category ID"
// @Param        request  body      productdto.UpdateCategoryRequest  true  "Category payload"
// @Success      200      {object}  response.Response
// @Failure      400      {object}  response.ErrorResponse
// @Failure      404      {object}  response.ErrorResponse
// @Failure      500      {object}  response.ErrorResponse
// @Router       /admin/categories/{id} [put]
func (h *AdminCategoryHandler) Update(c *gin.Context) {
	var req productdto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	cat, err := h.svc.UpdateCategory(c.Param("id"), req)
	if err != nil {
		if errors.Is(err, productsvc.ErrCategoryNotFound) {
			response.NotFound(c, "Category not found")
			return
		}
		response.InternalError(c, "Failed to update category")
		return
	}

	response.OK(c, "Category updated", cat)
}

// DeleteCategory godoc
// @Summary      Delete category
// @Description  Delete a category (only if it has no products)
// @Tags         Admin Categories
// @Produce      json
// @Param        id   path      string  true  "Category ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /admin/categories/{id} [delete]
func (h *AdminCategoryHandler) Delete(c *gin.Context) {
	if err := h.svc.DeleteCategory(c.Param("id")); err != nil {
		if errors.Is(err, productsvc.ErrCategoryNotFound) {
			response.NotFound(c, "Category not found")
			return
		}
		if errors.Is(err, productsvc.ErrCategoryHasProduct) {
			response.BadRequest(c, "Cannot delete category that still has products")
			return
		}
		response.InternalError(c, "Failed to delete category")
		return
	}
	response.OK(c, "Category deleted", nil)
}
