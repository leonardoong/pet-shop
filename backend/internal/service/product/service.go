package product

import (
	"errors"

	productdto "petshop/internal/dto/product"
	productmodel "petshop/internal/model/product"
	productrepo "petshop/internal/repository/product"
	"petshop/pkg/response"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrProductNotFound  = errors.New("product not found")
)

// --- Category service ---

type CategoryService interface {
	ListCategories() ([]productmodel.Category, error)
	GetCategoryBySlug(slug string) (*productmodel.Category, error)
}

type categoryService struct{ repo productrepo.CategoryRepository }

func NewCategoryService(repo productrepo.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) ListCategories() ([]productmodel.Category, error) {
	return s.repo.ListCategories()
}

func (s *categoryService) GetCategoryBySlug(slug string) (*productmodel.Category, error) {
	cat, err := s.repo.FindCategoryBySlug(slug)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, ErrCategoryNotFound
	}
	return cat, nil
}

// --- Product service ---

type ProductService interface {
	ListProducts(f productdto.Filter) (response.Paginated[productmodel.Product], error)
	GetProductBySlug(slug string) (*productmodel.Product, error)
}

type productService struct{ repo productrepo.ProductRepository }

func NewProductService(repo productrepo.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) ListProducts(f productdto.Filter) (response.Paginated[productmodel.Product], error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}

	products, total, err := s.repo.ListProducts(f)
	if err != nil {
		return response.Paginated[productmodel.Product]{}, err
	}

	return response.NewPaginated(products, total, f.Page, f.Limit), nil
}

func (s *productService) GetProductBySlug(slug string) (*productmodel.Product, error) {
	p, err := s.repo.FindProductBySlug(slug)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProductNotFound
	}
	return p, nil
}
