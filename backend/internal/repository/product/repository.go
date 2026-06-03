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
	FindCategoryByID(id string) (*productmodel.Category, error)
	CreateCategory(c *productmodel.Category) error
	UpdateCategory(c *productmodel.Category) error
	DeleteCategory(id string) error
}

// ProductRepository handles product persistence.
type ProductRepository interface {
	ListProducts(f dto.Filter) ([]productmodel.Product, int64, error)
	FindProductBySlug(slug string) (*productmodel.Product, error)
	FindProductByID(id string) (*productmodel.Product, error)
	DecrementStock(tx *gorm.DB, productID string, qty int) error

	CreateProduct(p *productmodel.Product) error
	UpdateProduct(p *productmodel.Product) error
	FindProductByIDAdmin(id string) (*productmodel.Product, error)
	FindProductBySKU(sku string) (*productmodel.Product, error)
	SlugExists(slug string) (bool, error)
	ListProductsAdmin(f dto.AdminProductFilter) ([]productmodel.Product, int64, error)
	CountProductsByCategory(categoryID string) (int64, error)

	CreateStockLog(sl *productmodel.StockLog) error
	TotalCOGS() (float64, error)

	CreateReview(r *productmodel.Review) error
	ListReviews(productID string, approved bool, page, limit int) ([]productmodel.Review, int64, error)
	ListAllReviews(f dto.AdminReviewFilter) ([]productmodel.Review, int64, error)
	FindReviewByID(id string) (*productmodel.Review, error)
	UpdateReview(r *productmodel.Review) error
	HasCustomerOrderedProduct(customerID, productID string) (bool, error)
	FindCustomerReview(customerID, productID string) (*productmodel.Review, error)
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

func (r *repository) FindCategoryByID(id string) (*productmodel.Category, error) {
	var cat productmodel.Category
	err := r.db.Select("id, name, slug, image_url, created_at, updated_at").
		Where("id = ?", id).First(&cat).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &cat, err
}

func (r *repository) CreateCategory(c *productmodel.Category) error {
	return r.db.Create(c).Error
}

func (r *repository) UpdateCategory(c *productmodel.Category) error {
	return r.db.Save(c).Error
}

func (r *repository) DeleteCategory(id string) error {
	return r.db.Delete(&productmodel.Category{}, "id = ?", id).Error
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

func (r *repository) CreateProduct(p *productmodel.Product) error {
	return r.db.Create(p).Error
}

func (r *repository) UpdateProduct(p *productmodel.Product) error {
	return r.db.Save(p).Error
}

func (r *repository) FindProductByIDAdmin(id string) (*productmodel.Product, error) {
	var p productmodel.Product
	err := r.db.Preload("Category").Where("id = ?", id).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &p, err
}

func (r *repository) FindProductBySKU(sku string) (*productmodel.Product, error) {
	var p productmodel.Product
	err := r.db.Where("sku = ?", sku).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &p, err
}

func (r *repository) SlugExists(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&productmodel.Product{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

func (r *repository) ListProductsAdmin(f dto.AdminProductFilter) ([]productmodel.Product, int64, error) {
	q := r.db.Model(&productmodel.Product{}).
		Select("products.id, products.category_id, products.name, products.slug, products.description, products.price, products.cost_price, products.stock, products.sku, products.image_url, products.is_active, products.created_at, products.updated_at").
		Preload("Category")

	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("products.name ILIKE ? OR products.sku ILIKE ?", like, like)
	}
	if f.CategoryID != "" {
		q = q.Where("products.category_id = ?", f.CategoryID)
	}
	if f.IsActive != nil {
		q = q.Where("products.is_active = ?", *f.IsActive)
	}
	if f.StockMax != nil {
		q = q.Where("products.stock <= ?", *f.StockMax)
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
	case "stock_asc":
		q = q.Order("products.stock asc")
	case "stock_desc":
		q = q.Order("products.stock desc")
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

func (r *repository) CountProductsByCategory(categoryID string) (int64, error) {
	var count int64
	err := r.db.Model(&productmodel.Product{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

func (r *repository) CreateStockLog(sl *productmodel.StockLog) error {
	return r.db.Create(sl).Error
}

func (r *repository) TotalCOGS() (float64, error) {
	var cogs float64
	err := r.db.Model(&productmodel.StockLog{}).
		Where("type = ?", productmodel.StockLogSale).
		Select("COALESCE(SUM(total_cost), 0)").
		Scan(&cogs).Error
	return cogs, err
}

func (r *repository) CreateReview(rv *productmodel.Review) error {
	return r.db.Create(rv).Error
}

func (r *repository) ListReviews(productID string, approved bool, page, limit int) ([]productmodel.Review, int64, error) {
	q := r.db.Model(&productmodel.Review{}).
		Preload("Customer").
		Where("product_id = ?", productID)
	if approved {
		q = q.Where("is_approved = true")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var reviews []productmodel.Review
	err := q.Order("created_at desc").Limit(limit).Offset(offset).Find(&reviews).Error
	return reviews, total, err
}

func (r *repository) ListAllReviews(f dto.AdminReviewFilter) ([]productmodel.Review, int64, error) {
	q := r.db.Model(&productmodel.Review{}).Preload("Customer")

	if f.IsApproved != nil {
		q = q.Where("is_approved = ?", *f.IsApproved)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Limit
	var reviews []productmodel.Review
	err := q.Order("created_at desc").Limit(f.Limit).Offset(offset).Find(&reviews).Error
	return reviews, total, err
}

func (r *repository) FindReviewByID(id string) (*productmodel.Review, error) {
	var rv productmodel.Review
	err := r.db.Preload("Customer").Where("id = ?", id).First(&rv).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rv, err
}

func (r *repository) UpdateReview(rv *productmodel.Review) error {
	return r.db.Save(rv).Error
}

func (r *repository) HasCustomerOrderedProduct(customerID, productID string) (bool, error) {
	var count int64
	err := r.db.Table("order_items").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.customer_id = ? AND order_items.product_id = ? AND orders.status = ?",
			customerID, productID, "delivered").
		Count(&count).Error
	return count > 0, err
}

func (r *repository) FindCustomerReview(customerID, productID string) (*productmodel.Review, error) {
	var rv productmodel.Review
	err := r.db.Where("customer_id = ? AND product_id = ?", customerID, productID).First(&rv).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rv, err
}
