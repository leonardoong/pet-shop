package order

import (
	"errors"

	ordermodel "petshop/internal/model/order"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(tx *gorm.DB, o *ordermodel.Order) error
	List(customerID uuid.UUID, status string, page, limit int) ([]ordermodel.Order, int64, error)
	FindByID(id, customerID uuid.UUID) (*ordermodel.Order, error)
	DB() *gorm.DB
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) DB() *gorm.DB { return r.db }

func (r *repository) Create(tx *gorm.DB, o *ordermodel.Order) error {
	return tx.Create(o).Error
}

func (r *repository) List(customerID uuid.UUID, status string, page, limit int) ([]ordermodel.Order, int64, error) {
	q := r.db.Model(&ordermodel.Order{}).Where("customer_id = ?", customerID)
	if status != "" {
		q = q.Where("status = ?", status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var orders []ordermodel.Order
	err := q.Preload("Items.Product").
		Order("created_at desc").
		Limit(limit).Offset(offset).
		Find(&orders).Error
	return orders, total, err
}

func (r *repository) FindByID(id, customerID uuid.UUID) (*ordermodel.Order, error) {
	var o ordermodel.Order
	err := r.db.
		Preload("Items.Product").
		Where("id = ? AND customer_id = ?", id, customerID).
		First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &o, err
}
