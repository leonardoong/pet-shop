package order

import (
	"errors"
	"fmt"

	orderdto "petshop/internal/dto/order"
	ordermodel "petshop/internal/model/order"
	orderrepo "petshop/internal/repository/order"
	"petshop/pkg/response"

	"github.com/google/uuid"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

var validTransitions = map[ordermodel.Status][]ordermodel.Status{
	ordermodel.StatusPending:    {ordermodel.StatusConfirmed, ordermodel.StatusCancelled},
	ordermodel.StatusConfirmed:  {ordermodel.StatusProcessing, ordermodel.StatusCancelled},
	ordermodel.StatusProcessing: {ordermodel.StatusShipped, ordermodel.StatusCancelled},
	ordermodel.StatusShipped:    {ordermodel.StatusDelivered},
}

type AdminService interface {
	ListOrders(f orderdto.AdminOrderFilter) (response.Paginated[orderdto.AdminOrderResponse], error)
	GetOrder(id string) (*orderdto.AdminOrderResponse, error)
	UpdateStatus(id string, req orderdto.UpdateStatusRequest) (*orderdto.AdminOrderResponse, error)
}

type adminService struct{ repo orderrepo.AdminRepository }

func NewAdminService(repo orderrepo.AdminRepository) AdminService {
	return &adminService{repo: repo}
}

func (s *adminService) ListOrders(f orderdto.AdminOrderFilter) (response.Paginated[orderdto.AdminOrderResponse], error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}

	orders, total, err := s.repo.ListAll(f)
	if err != nil {
		return response.Paginated[orderdto.AdminOrderResponse]{}, err
	}

	items := make([]orderdto.AdminOrderResponse, len(orders))
	for i, o := range orders {
		items[i] = toAdminResponse(&o)
	}

	return response.NewPaginated(items, total, f.Page, f.Limit), nil
}

func (s *adminService) GetOrder(id string) (*orderdto.AdminOrderResponse, error) {
	o, err := s.repo.FindByIDAdmin(id)
	if err != nil {
		return nil, err
	}
	if o == nil {
		return nil, ErrOrderNotFound
	}
	r := toAdminResponse(o)
	return &r, nil
}

func (s *adminService) UpdateStatus(id string, req orderdto.UpdateStatusRequest) (*orderdto.AdminOrderResponse, error) {
	o, err := s.repo.FindByIDAdmin(id)
	if err != nil {
		return nil, err
	}
	if o == nil {
		return nil, ErrOrderNotFound
	}

	newStatus := ordermodel.Status(req.Status)
	if !isValidTransition(o.Status, newStatus) {
		return nil, fmt.Errorf("%w: cannot transition from %s to %s", ErrInvalidStatusTransition, o.Status, newStatus)
	}

	if err := s.repo.UpdateStatus(id, newStatus); err != nil {
		return nil, err
	}

	o, _ = s.repo.FindByIDAdmin(id)
	r := toAdminResponse(o)
	return &r, nil
}

func isValidTransition(current, next ordermodel.Status) bool {
	allowed, ok := validTransitions[current]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == next {
			return true
		}
	}
	return false
}

func toAdminResponse(o *ordermodel.Order) orderdto.AdminOrderResponse {
	r := orderdto.AdminOrderResponse{
		ID:           o.ID.String(),
		CustomerID:   o.CustomerID.String(),
		Status:       string(o.Status),
		TotalAmount:  o.TotalAmount,
		ShipName:     o.ShipName,
		ShipPhone:    o.ShipPhone,
		ShipStreet:   o.ShipStreet,
		ShipCity:     o.ShipCity,
		ShipProvince: o.ShipProvince,
		ShipPostal:   o.ShipPostal,
		Notes:        o.Notes,
		CreatedAt:    o.CreatedAt,
		UpdatedAt:    o.UpdatedAt,
	}

	for _, item := range o.Items {
		ri := orderdto.AdminOrderItemResponse{
			ID:        item.ID.String(),
			ProductID: item.ProductID.String(),
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
		}
		if item.Product.ID != uuid.Nil {
			ri.Product = item.Product.Name
		}
		r.Items = append(r.Items, ri)
	}

	return r
}
