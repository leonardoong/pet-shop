package product

import (
	"errors"
	"strconv"

	productdto "petshop/internal/dto/product"
	productsvc "petshop/internal/service/product"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
)

// --- Category handler ---

type CategoryHandler struct{ svc productsvc.CategoryService }

func NewCategoryHandler(svc productsvc.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// ListCategories godoc
// @Summary      List categories
// @Description  Returns all active product categories
// @Tags         Categories
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrorResponse
// @Router       /categories [get]
func (h *CategoryHandler) List(c *gin.Context) {
	cats, err := h.svc.ListCategories()
	if err != nil {
		response.InternalError(c, "Failed to fetch categories")
		return
	}
	response.OK(c, "Success", cats)
}

// GetCategoryBySlug godoc
// @Summary      Get category by slug
// @Description  Returns a single category by its slug
// @Tags         Categories
// @Produce      json
// @Param        slug  path      string  true  "Category slug"
// @Success      200   {object}  response.Response
// @Failure      404   {object}  response.ErrorResponse
// @Failure      500   {object}  response.ErrorResponse
// @Router       /categories/{slug} [get]
func (h *CategoryHandler) GetBySlug(c *gin.Context) {
	cat, err := h.svc.GetCategoryBySlug(c.Param("slug"))
	if err != nil {
		if errors.Is(err, productsvc.ErrCategoryNotFound) {
			response.NotFound(c, "Category not found")
			return
		}
		response.InternalError(c, "Failed to fetch category")
		return
	}
	response.OK(c, "Success", cat)
}

// --- Product handler ---

type ProductHandler struct{ svc productsvc.ProductService }

func NewProductHandler(svc productsvc.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// ListProducts godoc
// @Summary      List products
// @Description  Returns a paginated list of active products with optional filters
// @Tags         Products
// @Produce      json
// @Param        category  query     string   false  "Filter by category slug"
// @Param        search    query     string   false  "Search by product name"
// @Param        min_price query     number   false  "Minimum price"
// @Param        max_price query     number   false  "Maximum price"
// @Param        sort      query     string   false  "Sort: newest|oldest|price_asc|price_desc"
// @Param        page      query     int      false  "Page number (default: 1)"
// @Param        limit     query     int      false  "Items per page (default: 20, max: 100)"
// @Success      200       {object}  response.Response{data=ProductsPage}
// @Failure      500       {object}  response.ErrorResponse
// @Router       /products [get]
func (h *ProductHandler) List(c *gin.Context) {
	f := productdto.Filter{
		CategorySlug: c.Query("category"),
		Search:       c.Query("search"),
		Sort:         c.Query("sort"),
		Page:         parseIntQuery(c, "page", 1),
		Limit:        parseIntQuery(c, "limit", 20),
	}

	if v := c.Query("min_price"); v != "" {
		if f64, err := strconv.ParseFloat(v, 64); err == nil {
			f.MinPrice = &f64
		}
	}
	if v := c.Query("max_price"); v != "" {
		if f64, err := strconv.ParseFloat(v, 64); err == nil {
			f.MaxPrice = &f64
		}
	}

	paginated, err := h.svc.ListProducts(f)
	if err != nil {
		response.InternalError(c, "Failed to fetch products")
		return
	}
	response.OK(c, "Success", paginated)
}

// GetProductBySlug godoc
// @Summary      Get product by slug
// @Description  Returns a single product by its slug
// @Tags         Products
// @Produce      json
// @Param        slug  path      string  true  "Product slug"
// @Success      200   {object}  response.Response
// @Failure      404   {object}  response.ErrorResponse
// @Failure      500   {object}  response.ErrorResponse
// @Router       /products/{slug} [get]
func (h *ProductHandler) GetBySlug(c *gin.Context) {
	p, err := h.svc.GetProductBySlug(c.Param("slug"))
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
