package product

import (
	"errors"

	dto "petshop/internal/dto/product"
	productmodel "petshop/internal/model/product"

	"gorm.io/gorm"
)

var ErrInsufficientStock = errors.New("insufficient stock")

// CategoryRepository handles category persistence.
type CategoryRepository interface {
	ListCategories() ([]productmodel.Category, error)
	FindCategoryBySlug(slug string) (*productmodel.Category, error)
}

// ProductRepository handles product persistence.
type ProductRepository interface {
	ListProducts(f dto.Filter) ([]productmodel.Product, int64, error)
	FindProductBySlug(slug string) (*productmodel.Product, error)
	FindProductByID(id string) (*productmodel.Product, error)
	DecrementStock(tx *gorm.DB, productID string, qty int) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) (CategoryRepository, ProductRepository) {
	r := &repository{db: db}
	return r, r
}

// --- Category ---

func (r *repository) ListCategories() ([]productmodel.Category, error) {
	var cats []productmodel.Category
	err := r.db.Select("id, name, slug, image_url, created_at, updated_at").
		Order("name asc").Find(&cats).Error
	return cats, err
}

func (r *repository) FindCategoryBySlug(slug string) (*productmodel.Category, error) {
	var cat productmodel.Category
	err := r.db.Select("id, name, slug, image_url, created_at, updated_at").
		Where("slug = ?", slug).First(&cat).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &cat, err
}

// --- Product ---

func (r *repository) ListProducts(f dto.Filter) ([]productmodel.Product, int64, error) {
	q := r.db.Model(&productmodel.Product{}).
		Select("products.id, products.category_id, products.name, products.slug, products.description, products.price, products.stock, products.sku, products.image_url, products.is_active, products.created_at, products.updated_at").
		Preload("Category").
		Where("products.is_active = true")

	if f.CategorySlug != "" {
		q = q.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.slug = ?", f.CategorySlug)
	}
	if f.Search != "" {
		q = q.Where("products.name ILIKE ?", "%"+f.Search+"%")
	}
	if f.MinPrice != nil {
		q = q.Where("products.price >= ?", *f.MinPrice)
	}
	if f.MaxPrice != nil {
		q = q.Where("products.price <= ?", *f.MaxPrice)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch f.Sort {
	case "price_asc":
		q = q.Order("products.price asc")
	case "price_desc":
		q = q.Order("products.price desc")
	case "oldest":
		q = q.Order("products.created_at asc")
	default:
		q = q.Order("products.created_at desc")
	}

	offset := (f.Page - 1) * f.Limit
	var products []productmodel.Product
	err := q.Limit(f.Limit).Offset(offset).Find(&products).Error
	return products, total, err
}

func (r *repository) FindProductBySlug(slug string) (*productmodel.Product, error) {
	var p productmodel.Product
	err := r.db.Preload("Category").Where("slug = ? AND is_active = true", slug).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &p, err
}

func (r *repository) FindProductByID(id string) (*productmodel.Product, error) {
	var p productmodel.Product
	err := r.db.Where("id = ? AND is_active = true", id).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &p, err
}

func (r *repository) DecrementStock(tx *gorm.DB, productID string, qty int) error {
	result := tx.Model(&productmodel.Product{}).
		Where("id = ? AND stock >= ?", productID, qty).
		UpdateColumn("stock", gorm.Expr("stock - ?", qty))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrInsufficientStock
	}
	return nil
}
