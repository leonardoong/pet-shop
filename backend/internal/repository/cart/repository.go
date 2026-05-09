package cart

import (
	"errors"

	cartmodel "petshop/internal/model/cart"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetByCustomerID(customerID uuid.UUID) (*cartmodel.Cart, error)
	Create(cart *cartmodel.Cart) error
	AddOrUpdateItem(cartID, productID uuid.UUID, qty int) error
	UpdateItemQuantity(cartID, productID uuid.UUID, qty int) error
	RemoveItem(cartID, productID uuid.UUID) error
	ClearItems(tx *gorm.DB, cartID uuid.UUID) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) GetByCustomerID(customerID uuid.UUID) (*cartmodel.Cart, error) {
	var c cartmodel.Cart
	err := r.db.
		Preload("Items.Product.Category").
		Where("customer_id = ?", customerID).
		First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *repository) Create(cart *cartmodel.Cart) error {
	return r.db.Create(cart).Error
}

// AddOrUpdateItem atomically increments quantity when item exists, inserts otherwise.
func (r *repository) AddOrUpdateItem(cartID, productID uuid.UUID, qty int) error {
	return r.db.Exec(`
		INSERT INTO cart_items (id, cart_id, product_id, quantity, created_at, updated_at)
		VALUES (gen_random_uuid(), ?, ?, ?, NOW(), NOW())
		ON CONFLICT (cart_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = NOW()
	`, cartID, productID, qty).Error
}

func (r *repository) UpdateItemQuantity(cartID, productID uuid.UUID, qty int) error {
	return r.db.
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		Updates(map[string]any{"quantity": qty}).Error
}

func (r *repository) RemoveItem(cartID, productID uuid.UUID) error {
	return r.db.
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		Delete(&cartmodel.CartItem{}).Error
}

func (r *repository) ClearItems(tx *gorm.DB, cartID uuid.UUID) error {
	return tx.Where("cart_id = ?", cartID).Delete(&cartmodel.CartItem{}).Error
}
