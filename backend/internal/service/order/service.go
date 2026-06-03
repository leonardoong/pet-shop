package order

import (
	"errors"
	"fmt"

	orderdto "petshop/internal/dto/order"
	ordermodel "petshop/internal/model/order"
	productmodel "petshop/internal/model/product"
	cartrepo "petshop/internal/repository/cart"
	orderrepo "petshop/internal/repository/order"
	productrepo "petshop/internal/repository/product"
	userrepo "petshop/internal/repository/user"
	"petshop/pkg/payment"
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
	Checkout(customerID uuid.UUID, req orderdto.CreateRequest) (*orderdto.CheckoutResponse, error)
	ListOrders(customerID uuid.UUID, f orderdto.ListFilter) (response.Paginated[ordermodel.Order], error)
	GetOrderByID(id, customerID uuid.UUID) (*ordermodel.Order, error)
}

type service struct {
	repo        orderrepo.Repository
	cartRepo    cartrepo.Repository
	addrRepo    userrepo.Repository
	productRepo productrepo.ProductRepository
	paymentDrv  payment.Driver
}

func NewService(
	repo orderrepo.Repository,
	cartRepo cartrepo.Repository,
	addrRepo userrepo.Repository,
	productRepo productrepo.ProductRepository,
	paymentDrv payment.Driver,
) Service {
	return &service{
		repo:        repo,
		cartRepo:    cartRepo,
		addrRepo:    addrRepo,
		productRepo: productRepo,
		paymentDrv:  paymentDrv,
	}
}

func (s *service) Checkout(customerID uuid.UUID, req orderdto.CreateRequest) (*orderdto.CheckoutResponse, error) {
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
	payItems := make([]payment.ItemDetail, 0, len(cart.Items))
	for _, item := range cart.Items {
		subtotal := item.Product.Price * float64(item.Quantity)
		total += subtotal
		items = append(items, ordermodel.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.Product.Price,
			Subtotal:  subtotal,
		})
		payItems = append(payItems, payment.ItemDetail{
			ID: item.ProductID.String(), Name: item.Product.Name,
			Price: item.Product.Price, Quantity: item.Quantity,
		})
	}

	o := &ordermodel.Order{
		CustomerID:    customerID,
		Status:        ordermodel.StatusPending,
		PaymentStatus: ordermodel.PaymentStatusPending,
		TotalAmount:   total,
		ShipName:      addr.RecipientName,
		ShipPhone:     addr.Phone,
		ShipStreet:    addr.Street,
		ShipCity:      addr.City,
		ShipProvince:  addr.Province,
		ShipPostal:    addr.PostalCode,
		Notes:         req.Notes,
		Items:         items,
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

	for _, item := range cart.Items {
		costPrice := max(item.Product.CostPrice, 0)
		sl := &productmodel.StockLog{
			ProductID: item.ProductID,
			Type:      productmodel.StockLogSale,
			Quantity:  item.Quantity,
			CostPrice: costPrice,
			TotalCost: costPrice * float64(item.Quantity),
			Note:      fmt.Sprintf("Order #%s", o.ID.String()[:8]),
		}
		s.productRepo.CreateStockLog(sl)
	}

	resp := &orderdto.CheckoutResponse{Order: *o}

	if s.paymentDrv != nil {
		payResp, err := s.paymentDrv.CreateTransaction(payment.CreateRequest{
			OrderID:       o.ID.String(),
			GrossAmount:   total,
			CustomerName:  addr.RecipientName,
			CustomerEmail: req.CustomerEmail,
			CustomerPhone: addr.Phone,
			Items:         payItems,
		})
		if err == nil && payResp != nil {
			o.PaymentTransactionID = payResp.TransactionID
			o.PaymentURL = payResp.PaymentURL
			s.repo.UpdateOrder(o)
			resp.PaymentURL = payResp.PaymentURL
			resp.PaymentToken = payResp.PaymentToken
		}
	}

	return resp, nil
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
