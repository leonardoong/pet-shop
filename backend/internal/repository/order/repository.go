package order

import (
	"errors"
	"fmt"
	"time"

	orderdto "petshop/internal/dto/order"
	ordermodel "petshop/internal/model/order"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(tx *gorm.DB, o *ordermodel.Order) error
	List(customerID uuid.UUID, status string, page, limit int) ([]ordermodel.Order, int64, error)
	FindByID(id, customerID uuid.UUID) (*ordermodel.Order, error)
	UpdateOrder(o *ordermodel.Order) error
	DB() *gorm.DB
}

type AdminRepository interface {
	ListAll(f orderdto.AdminOrderFilter) ([]ordermodel.Order, int64, error)
	FindByIDAdmin(id string) (*ordermodel.Order, error)
	UpdateStatus(id string, status ordermodel.Status) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) (Repository, AdminRepository) {
	r := &repository{db: db}
	return r, r
}

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

func (r *repository) UpdateOrder(o *ordermodel.Order) error {
	return r.db.Save(o).Error
}

func (r *repository) ListAll(f orderdto.AdminOrderFilter) ([]ordermodel.Order, int64, error) {
	q := r.db.Model(&ordermodel.Order{})

	if f.Status != "" {
		q = q.Where("orders.status = ?", f.Status)
	}
	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("orders.ship_name ILIKE ? OR CAST(orders.id AS TEXT) ILIKE ?", like, like)
	}
	if f.DateFrom != "" {
		if t, err := time.Parse("2006-01-02", f.DateFrom); err == nil {
			q = q.Where("orders.created_at >= ?", t)
		}
	}
	if f.DateTo != "" {
		if t, err := time.Parse("2006-01-02", f.DateTo); err == nil {
			q = q.Where("orders.created_at <= ?", t.Add(24*time.Hour))
		}
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Limit
	var orders []ordermodel.Order
	err := q.Preload("Items.Product").
		Order("orders.created_at desc").
		Limit(f.Limit).Offset(offset).
		Find(&orders).Error
	return orders, total, err
}

func (r *repository) FindByIDAdmin(id string) (*ordermodel.Order, error) {
	var o ordermodel.Order
	err := r.db.
		Preload("Items.Product").
		Where("id = ?", id).
		First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &o, err
}

func (r *repository) UpdateStatus(id string, status ordermodel.Status) error {
	result := r.db.Model(&ordermodel.Order{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}
