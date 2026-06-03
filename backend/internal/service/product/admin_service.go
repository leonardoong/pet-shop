package product

import (
	"errors"
	"fmt"
	"strings"

	productdto "petshop/internal/dto/product"
	productmodel "petshop/internal/model/product"
	productrepo "petshop/internal/repository/product"
	"petshop/pkg/response"
	"petshop/pkg/slug"

	"github.com/google/uuid"
)

var (
	ErrSKUExists          = errors.New("sku already exists")
	ErrCategoryHasProduct = errors.New("category still has products")
)

type AdminProductService interface {
	CreateProduct(req productdto.CreateProductRequest) (*productdto.AdminProductResponse, error)
	UpdateProduct(id string, req productdto.UpdateProductRequest) (*productdto.AdminProductResponse, error)
	DeleteProduct(id string) error
	GetProduct(id string) (*productdto.AdminProductResponse, error)
	ListProducts(f productdto.AdminProductFilter) (response.Paginated[productdto.AdminProductResponse], error)
}

type AdminCategoryService interface {
	CreateCategory(req productdto.CreateCategoryRequest) (*productmodel.Category, error)
	UpdateCategory(id string, req productdto.UpdateCategoryRequest) (*productmodel.Category, error)
	DeleteCategory(id string) error
}

type adminProductService struct {
	repo         productrepo.ProductRepository
	categoryRepo productrepo.CategoryRepository
}

func NewAdminProductService(repo productrepo.ProductRepository, categoryRepo productrepo.CategoryRepository) AdminProductService {
	return &adminProductService{repo: repo, categoryRepo: categoryRepo}
}

func (s *adminProductService) CreateProduct(req productdto.CreateProductRequest) (*productdto.AdminProductResponse, error) {
	existing, err := s.repo.FindProductBySKU(req.SKU)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrSKUExists
	}

	cat, err := s.categoryRepo.FindCategoryByID(req.CategoryID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, ErrCategoryNotFound
	}

	generatedSlug, err := slug.Unique(req.Name, s.repo.SlugExists)
	if err != nil {
		return nil, fmt.Errorf("slug generation: %w", err)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	p := &productmodel.Product{
		CategoryID:  cat.ID,
		Name:        strings.TrimSpace(req.Name),
		Slug:        generatedSlug,
		Description: strings.TrimSpace(req.Description),
		Price:       req.Price,
		CostPrice:   req.CostPrice,
		Stock:       req.Stock,
		SKU:         strings.TrimSpace(req.SKU),
		ImageURL:    req.ImageURL,
		IsActive:    isActive,
	}

	if err := s.repo.CreateProduct(p); err != nil {
		return nil, err
	}

	return toAdminResponse(p), nil
}

func (s *adminProductService) UpdateProduct(id string, req productdto.UpdateProductRequest) (*productdto.AdminProductResponse, error) {
	p, err := s.repo.FindProductByIDAdmin(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProductNotFound
	}

	if req.Name != nil {
		p.Name = strings.TrimSpace(*req.Name)
		if 		p.Slug != slug.Make(p.Name) {
			newSlug, err := slug.Unique(p.Name, func(candidate string) (bool, error) {
				ok, e := s.repo.SlugExists(candidate)
				if e != nil {
					return false, e
				}
				if ok && candidate == p.Slug {
					return false, nil
				}
				return ok, nil
			})
			if err != nil {
				return nil, fmt.Errorf("slug generation: %w", err)
			}
			p.Slug = newSlug
		}
	}
	if req.CategoryID != nil {
		catID, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			return nil, ErrCategoryNotFound
		}
		cat, err := s.categoryRepo.FindCategoryByID(*req.CategoryID)
		if err != nil {
			return nil, err
		}
		if cat == nil {
			return nil, ErrCategoryNotFound
		}
		p.CategoryID = catID
	}
	if req.Description != nil {
		p.Description = strings.TrimSpace(*req.Description)
	}
	if req.Price != nil {
		p.Price = *req.Price
	}
	if req.CostPrice != nil {
		p.CostPrice = *req.CostPrice
	}
	if req.Stock != nil {
		p.Stock = *req.Stock
	}
	if req.SKU != nil {
		existing, err := s.repo.FindProductBySKU(*req.SKU)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID.String() != p.ID.String() {
			return nil, ErrSKUExists
		}
		p.SKU = strings.TrimSpace(*req.SKU)
	}
	if req.ImageURL != nil {
		p.ImageURL = *req.ImageURL
	}
	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateProduct(p); err != nil {
		return nil, err
	}

	p, _ = s.repo.FindProductByIDAdmin(id)
	return toAdminResponse(p), nil
}

func (s *adminProductService) DeleteProduct(id string) error {
	p, err := s.repo.FindProductByIDAdmin(id)
	if err != nil {
		return err
	}
	if p == nil {
		return ErrProductNotFound
	}
	p.IsActive = false
	return s.repo.UpdateProduct(p)
}

func (s *adminProductService) GetProduct(id string) (*productdto.AdminProductResponse, error) {
	p, err := s.repo.FindProductByIDAdmin(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProductNotFound
	}
	return toAdminResponse(p), nil
}

func (s *adminProductService) ListProducts(f productdto.AdminProductFilter) (response.Paginated[productdto.AdminProductResponse], error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}

	products, total, err := s.repo.ListProductsAdmin(f)
	if err != nil {
		return response.Paginated[productdto.AdminProductResponse]{}, err
	}

	items := make([]productdto.AdminProductResponse, len(products))
	for i, p := range products {
		items[i] = *toAdminResponse(&p)
	}

	return response.NewPaginated(items, total, f.Page, f.Limit), nil
}

// --- AdminCategoryService ---

type adminCategoryService struct {
	repo         productrepo.CategoryRepository
	productRepo  productrepo.ProductRepository
}

func NewAdminCategoryService(repo productrepo.CategoryRepository, productRepo productrepo.ProductRepository) AdminCategoryService {
	return &adminCategoryService{repo: repo, productRepo: productRepo}
}

func (s *adminCategoryService) CreateCategory(req productdto.CreateCategoryRequest) (*productmodel.Category, error) {
	c := &productmodel.Category{
		Name:     strings.TrimSpace(req.Name),
		Slug:     slug.Make(req.Name),
		ImageURL: req.ImageURL,
	}
	if err := s.repo.CreateCategory(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *adminCategoryService) UpdateCategory(id string, req productdto.UpdateCategoryRequest) (*productmodel.Category, error) {
	cat, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, ErrCategoryNotFound
	}

	if req.Name != nil {
		cat.Name = strings.TrimSpace(*req.Name)
		cat.Slug = slug.Make(*req.Name)
	}
	if req.ImageURL != nil {
		cat.ImageURL = *req.ImageURL
	}

	if err := s.repo.UpdateCategory(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *adminCategoryService) DeleteCategory(id string) error {
	cat, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return err
	}
	if cat == nil {
		return ErrCategoryNotFound
	}

	count, err := s.productRepo.CountProductsByCategory(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrCategoryHasProduct
	}

	return s.repo.DeleteCategory(id)
}

func toAdminResponse(p *productmodel.Product) *productdto.AdminProductResponse {
	r := &productdto.AdminProductResponse{
		ID:          p.ID.String(),
		CategoryID:  p.CategoryID.String(),
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		CostPrice:   p.CostPrice,
		Stock:       p.Stock,
		SKU:         p.SKU,
		ImageURL:    p.ImageURL,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	if p.Category.ID != uuid.Nil {
		r.Category = p.Category.Name
	}
	return r
}
