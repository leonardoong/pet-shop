package product

import (
	"errors"

	productdto "petshop/internal/dto/product"
	"petshop/internal/middleware"
	productsvc "petshop/internal/service/product"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct{ svc productsvc.ReviewService }

func NewReviewHandler(svc productsvc.ReviewService) *ReviewHandler {
	return &ReviewHandler{svc: svc}
}

// ListReviews godoc
// @Summary      List product reviews
// @Description  Returns approved reviews for a product
// @Tags         Reviews
// @Produce      json
// @Param        slug   path      string  true  "Product slug"
// @Param        page   query     int     false "Page number"
// @Param        limit  query     int     false "Items per page"
// @Success      200    {object}  response.Response
// @Router       /products/{slug}/reviews [get]
func (h *ReviewHandler) ListByProduct(c *gin.Context) {
	slug := c.Param("slug")
	page := parseIntQuery(c, "page", 1)
	limit := parseIntQuery(c, "limit", 10)

	p, err := h.svc.ListReviews(slug, page, limit)
	if err != nil {
		response.InternalError(c, "Failed to fetch reviews")
		return
	}
	response.OK(c, "Success", p)
}

// CreateReview godoc
// @Summary      Create product review
// @Description  Rate and review a purchased product
// @Tags         Reviews
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        review  body      productdto.CreateReviewRequest  true  "Review"
// @Success      201     {object}  response.Response
// @Router       /customer/reviews [post]
func (h *ReviewHandler) Create(c *gin.Context) {
	var req productdto.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	customerID := middleware.MustGetCustomerID(c).String()
	productID := c.Query("product_id")
	if productID == "" {
		response.BadRequest(c, "product_id query param required")
		return
	}

	r, err := h.svc.CreateReview(customerID, productID, req)
	if err != nil {
		if errors.Is(err, productsvc.ErrNotPurchased) {
			response.BadRequest(c, "You can only review products you have purchased")
			return
		}
		if errors.Is(err, productsvc.ErrAlreadyReviewed) {
			response.Conflict(c, "You have already reviewed this product")
			return
		}
		response.InternalError(c, "Failed to create review")
		return
	}
	response.Created(c, "Review submitted for approval", r)
}

// AdminListReviews godoc
// @Summary      List all reviews (admin)
// @Description  Admin moderation queue for reviews
// @Tags         Admin Reviews
// @Produce      json
// @Param        pending  query     bool  false "Show only pending reviews"
// @Param        page     query     int   false "Page"
// @Param        limit    query     int   false "Limit"
// @Success      200      {object}  response.Response
// @Router       /admin/reviews [get]
func (h *ReviewHandler) AdminList(c *gin.Context) {
	f := productdto.AdminReviewFilter{
		Page:  parseIntQuery(c, "page", 1),
		Limit: parseIntQuery(c, "limit", 20),
	}

	if v := c.Query("pending"); v == "true" {
		pending := false
		f.IsApproved = &pending
	}

	paginated, err := h.svc.ListAllReviews(f)
	if err != nil {
		response.InternalError(c, "Failed to fetch reviews")
		return
	}
	response.OK(c, "Success", paginated)
}

// AdminToggleApproval godoc
// @Summary      Toggle review approval
// @Description  Approve or unapprove a review
// @Tags         Admin Reviews
// @Produce      json
// @Param        id  path  string  true  "Review ID"
// @Success      200  {object}  response.Response
// @Router       /admin/reviews/{id}/approve [patch]
func (h *ReviewHandler) AdminToggle(c *gin.Context) {
	r, err := h.svc.ToggleApproval(c.Param("id"))
	if err != nil {
		if errors.Is(err, productsvc.ErrReviewNotFound) {
			response.NotFound(c, "Review not found")
			return
		}
		response.InternalError(c, "Failed to toggle approval")
		return
	}
	response.OK(c, "Review updated", r)
}
