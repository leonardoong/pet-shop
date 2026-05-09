package cart

import (
	"errors"

	cartdto "petshop/internal/dto/cart"
	cartmodel "petshop/internal/model/cart"
	cartrepo "petshop/internal/repository/cart"
	productrepo "petshop/internal/repository/product"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrItemNotFound      = errors.New("cart item not found")
)

type Service interface {
	GetCart(customerID uuid.UUID) (*cartmodel.Cart, error)
	AddItem(customerID uuid.UUID, req cartdto.AddItemRequest) (*cartmodel.Cart, error)
	UpdateItem(customerID uuid.UUID, productID uuid.UUID, req cartdto.UpdateItemRequest) (*cartmodel.Cart, error)
	RemoveItem(customerID uuid.UUID, productID uuid.UUID) error
}

type service struct {
	repo        cartrepo.Repository
	productRepo productrepo.ProductRepository
}

func NewService(repo cartrepo.Repository, productRepo productrepo.ProductRepository) Service {
	return &service{repo: repo, productRepo: productRepo}
}

func (s *service) getOrCreateCart(customerID uuid.UUID) (*cartmodel.Cart, error) {
	c, err := s.repo.GetByCustomerID(customerID)
	if err != nil {
		return nil, err
	}
	if c != nil {
		return c, nil
	}
	c = &cartmodel.Cart{CustomerID: customerID}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *service) GetCart(customerID uuid.UUID) (*cartmodel.Cart, error) {
	return s.getOrCreateCart(customerID)
}

func (s *service) AddItem(customerID uuid.UUID, req cartdto.AddItemRequest) (*cartmodel.Cart, error) {
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	p, err := s.productRepo.FindProductByID(req.ProductID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProductNotFound
	}
	if p.Stock < req.Quantity {
		return nil, ErrInsufficientStock
	}

	c, err := s.getOrCreateCart(customerID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.AddOrUpdateItem(c.ID, productID, req.Quantity); err != nil {
		return nil, err
	}

	return s.repo.GetByCustomerID(customerID)
}

func (s *service) UpdateItem(customerID uuid.UUID, productID uuid.UUID, req cartdto.UpdateItemRequest) (*cartmodel.Cart, error) {
	p, err := s.productRepo.FindProductByID(productID.String())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProductNotFound
	}
	if p.Stock < req.Quantity {
		return nil, ErrInsufficientStock
	}

	c, err := s.repo.GetByCustomerID(customerID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrItemNotFound
	}

	if err := s.repo.UpdateItemQuantity(c.ID, productID, req.Quantity); err != nil {
		return nil, err
	}

	return s.repo.GetByCustomerID(customerID)
}

func (s *service) RemoveItem(customerID uuid.UUID, productID uuid.UUID) error {
	c, err := s.repo.GetByCustomerID(customerID)
	if err != nil {
		return err
	}
	if c == nil {
		return ErrItemNotFound
	}
	return s.repo.RemoveItem(c.ID, productID)
}
