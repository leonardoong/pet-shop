package user

import (
	"errors"

	userdto "petshop/internal/dto/user"
	usermodel "petshop/internal/model/user"
	userrepo "petshop/internal/repository/user"
	"petshop/pkg/response"
)

var (
	ErrCustomerNotFound = errors.New("customer not found")
)

type AdminCustomerService interface {
	ListCustomers(f userdto.AdminCustomerFilter) (response.Paginated[userdto.CustomerResponse], error)
	GetCustomer(id string) (*userdto.CustomerResponse, error)
	ToggleActive(id string, active bool) (*usermodel.Customer, error)
}

type adminCustomerService struct{ repo userrepo.AdminRepository }

func NewAdminCustomerService(repo userrepo.AdminRepository) AdminCustomerService {
	return &adminCustomerService{repo: repo}
}

func (s *adminCustomerService) ListCustomers(f userdto.AdminCustomerFilter) (response.Paginated[userdto.CustomerResponse], error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}

	customers, total, err := s.repo.ListCustomers(f)
	if err != nil {
		return response.Paginated[userdto.CustomerResponse]{}, err
	}

	items := make([]userdto.CustomerResponse, len(customers))
	for i, c := range customers {
		items[i] = userdto.CustomerResponse{
			ID:        c.ID.String(),
			FullName:  c.FullName,
			Email:     c.Email,
			Phone:     c.Phone,
			IsActive:  c.IsActive,
			CreatedAt: c.CreatedAt,
		}
	}

	return response.NewPaginated(items, total, f.Page, f.Limit), nil
}

func (s *adminCustomerService) GetCustomer(id string) (*userdto.CustomerResponse, error) {
	c, err := s.repo.FindCustomerByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCustomerNotFound
	}

	return &userdto.CustomerResponse{
		ID:        c.ID.String(),
		FullName:  c.FullName,
		Email:     c.Email,
		Phone:     c.Phone,
		IsActive:  c.IsActive,
		CreatedAt: c.CreatedAt,
	}, nil
}

func (s *adminCustomerService) ToggleActive(id string, active bool) (*usermodel.Customer, error) {
	c, err := s.repo.FindCustomerByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCustomerNotFound
	}

	c.IsActive = active
	if err := s.repo.UpdateCustomer(c); err != nil {
		return nil, err
	}
	return c, nil
}
