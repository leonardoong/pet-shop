package product

import (
	"errors"
	"fmt"

	productdto "petshop/internal/dto/product"
	productmodel "petshop/internal/model/product"
	productrepo "petshop/internal/repository/product"
	"petshop/pkg/response"

	"github.com/google/uuid"
)

var (
	ErrAlreadyReviewed  = errors.New("customer already reviewed this product")
	ErrNotPurchased      = errors.New("customer has not purchased this product")
	ErrReviewNotFound    = errors.New("review not found")
)

type ReviewService interface {
	CreateReview(customerID, productID string, req productdto.CreateReviewRequest) (*productdto.ReviewResponse, error)
	ListReviews(productID string, page, limit int) (response.Paginated[productdto.ReviewResponse], error)
	ListAllReviews(f productdto.AdminReviewFilter) (response.Paginated[productmodel.Review], error)
	ToggleApproval(id string) (*productmodel.Review, error)
}

type reviewService struct {
	repo productrepo.ProductRepository
}

func NewReviewService(repo productrepo.ProductRepository) ReviewService {
	return &reviewService{repo: repo}
}

func (s *reviewService) CreateReview(customerID, productID string, req productdto.CreateReviewRequest) (*productdto.ReviewResponse, error) {
	ordered, err := s.repo.HasCustomerOrderedProduct(customerID, productID)
	if err != nil {
		return nil, err
	}
	if !ordered {
		return nil, ErrNotPurchased
	}

	existing, err := s.repo.FindCustomerReview(customerID, productID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrAlreadyReviewed
	}

	cUID, _ := uuid.Parse(customerID)
	pUID, _ := uuid.Parse(productID)

	r := &productmodel.Review{
		ProductID:  pUID,
		CustomerID: cUID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		IsApproved: false,
	}
	if err := s.repo.CreateReview(r); err != nil {
		return nil, err
	}

	return toReviewResponse(r), nil
}

func (s *reviewService) ListReviews(productSlug string, page, limit int) (response.Paginated[productdto.ReviewResponse], error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 50 { limit = 10 }

	p, err := s.repo.FindProductBySlug(productSlug)
	if err != nil {
		return response.Paginated[productdto.ReviewResponse]{}, err
	}
	if p == nil {
		return response.Paginated[productdto.ReviewResponse]{}, ErrProductNotFound
	}

	reviews, total, err := s.repo.ListReviews(p.ID.String(), true, page, limit)
	if err != nil {
		return response.Paginated[productdto.ReviewResponse]{}, err
	}

	items := make([]productdto.ReviewResponse, len(reviews))
	for i, r := range reviews {
		items[i] = *toReviewResponse(&r)
	}
	return response.NewPaginated(items, total, page, limit), nil
}

func (s *reviewService) ListAllReviews(f productdto.AdminReviewFilter) (response.Paginated[productmodel.Review], error) {
	if f.Page < 1 { f.Page = 1 }
	if f.Limit < 1 || f.Limit > 50 { f.Limit = 20 }

	reviews, total, err := s.repo.ListAllReviews(f)
	if err != nil {
		_ = fmt.Errorf("list reviews: %w", err)
		return response.Paginated[productmodel.Review]{}, err
	}
	return response.NewPaginated(reviews, total, f.Page, f.Limit), nil
}

func (s *reviewService) ToggleApproval(id string) (*productmodel.Review, error) {
	r, err := s.repo.FindReviewByID(id)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, ErrReviewNotFound
	}

	r.IsApproved = !r.IsApproved
	if err := s.repo.UpdateReview(r); err != nil {
		return nil, err
	}
	return r, nil
}

func toReviewResponse(r *productmodel.Review) *productdto.ReviewResponse {
	name := r.Customer.FullName
	if name == "" {
		name = "Pelanggan"
	}
	return &productdto.ReviewResponse{
		ID:           r.ID.String(),
		CustomerName: name,
		Rating:       r.Rating,
		Comment:      r.Comment,
		IsApproved:   r.IsApproved,
		CreatedAt:    r.CreatedAt.Format("2006-01-02"),
	}
}
