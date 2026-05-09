package order

import (
	"errors"

	orderdto "petshop/internal/dto/order"
	ordermodel "petshop/internal/model/order"
	cartrepo "petshop/internal/repository/cart"
	orderrepo "petshop/internal/repository/order"
	productrepo "petshop/internal/repository/product"
	userrepo "petshop/internal/repository/user"
	"petshop/pkg/response"

	"github.com/google/uuid"
)

var (
	ErrEmptyCart         = errors.New("cart is empty")
	ErrAddressNotFound   = errors.New("address not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrOrderNotFound     = errors.New("order not found")
)

type Service interface {
	Checkout(customerID uuid.UUID, req orderdto.CreateRequest) (*ordermodel.Order, error)
	ListOrders(customerID uuid.UUID, f orderdto.ListFilter) (response.Paginated[ordermodel.Order], error)
	GetOrderByID(id, customerID uuid.UUID) (*ordermodel.Order, error)
}

type service struct {
	repo        orderrepo.Repository
	cartRepo    cartrepo.Repository
	addrRepo    userrepo.Repository
	productRepo productrepo.ProductRepository
}

func NewService(
	repo orderrepo.Repository,
	cartRepo cartrepo.Repository,
	addrRepo userrepo.Repository,
	productRepo productrepo.ProductRepository,
) Service {
	return &service{
		repo:        repo,
		cartRepo:    cartRepo,
		addrRepo:    addrRepo,
		productRepo: productRepo,
	}
}

func (s *service) Checkout(customerID uuid.UUID, req orderdto.CreateRequest) (*ordermodel.Order, error) {
	addrID, err := uuid.Parse(req.AddressID)
	if err != nil {
		return nil, ErrAddressNotFound
	}

	addr, err := s.addrRepo.FindAddressByID(addrID, customerID)
	if err != nil {
		return nil, err
	}
	if addr == nil {
		return nil, ErrAddressNotFound
	}

	cart, err := s.cartRepo.GetByCustomerID(customerID)
	if err != nil {
		return nil, err
	}
	if cart == nil || len(cart.Items) == 0 {
		return nil, ErrEmptyCart
	}

	var total float64
	items := make([]ordermodel.OrderItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		subtotal := item.Product.Price * float64(item.Quantity)
		total += subtotal
		items = append(items, ordermodel.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.Product.Price,
			Subtotal:  subtotal,
		})
	}

	o := &ordermodel.Order{
		CustomerID:   customerID,
		Status:       ordermodel.StatusPending,
		TotalAmount:  total,
		ShipName:     addr.RecipientName,
		ShipPhone:    addr.Phone,
		ShipStreet:   addr.Street,
		ShipCity:     addr.City,
		ShipProvince: addr.Province,
		ShipPostal:   addr.PostalCode,
		Notes:        req.Notes,
		Items:        items,
	}

	db := s.repo.DB()
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	for _, item := range cart.Items {
		if err := s.productRepo.DecrementStock(tx, item.ProductID.String(), item.Quantity); err != nil {
			tx.Rollback()
			if errors.Is(err, productrepo.ErrInsufficientStock) {
				return nil, ErrInsufficientStock
			}
			return nil, err
		}
	}

	if err := s.repo.Create(tx, o); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.cartRepo.ClearItems(tx, cart.ID); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return o, nil
}

func (s *service) ListOrders(customerID uuid.UUID, f orderdto.ListFilter) (response.Paginated[ordermodel.Order], error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}

	orders, total, err := s.repo.List(customerID, f.Status, f.Page, f.Limit)
	if err != nil {
		return response.Paginated[ordermodel.Order]{}, err
	}
	return response.NewPaginated(orders, total, f.Page, f.Limit), nil
}

func (s *service) GetOrderByID(id, customerID uuid.UUID) (*ordermodel.Order, error) {
	o, err := s.repo.FindByID(id, customerID)
	if err != nil {
		return nil, err
	}
	if o == nil {
		return nil, ErrOrderNotFound
	}
	return o, nil
}
