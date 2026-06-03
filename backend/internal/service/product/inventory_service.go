package product

import (
	"errors"
	"fmt"
	"math"

	productdto "petshop/internal/dto/product"
	productmodel "petshop/internal/model/product"
	productrepo "petshop/internal/repository/product"
	"petshop/pkg/response"
)

var (
	ErrInvalidOperation = errors.New("invalid operation")
)

type InventoryService interface {
	ListInventory(f productdto.InventoryFilter) (response.Paginated[productdto.InventoryItem], error)
	AdjustStock(productID string, req productdto.AdjustStockRequest) (*productdto.InventoryItem, error)
}

type inventoryService struct {
	repo         productrepo.ProductRepository
	categoryRepo productrepo.CategoryRepository
}

func NewInventoryService(repo productrepo.ProductRepository, categoryRepo productrepo.CategoryRepository) InventoryService {
	return &inventoryService{repo: repo, categoryRepo: categoryRepo}
}

func (s *inventoryService) ListInventory(f productdto.InventoryFilter) (response.Paginated[productdto.InventoryItem], error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 20
	}

	var isActive *bool
	var stockMax *int
	if f.LowStock != nil && *f.LowStock {
		active := true
		isActive = &active
		threshold := 10
		stockMax = &threshold
	}

	adminF := productdto.AdminProductFilter{
		Search:     f.Search,
		CategoryID: f.CategoryID,
		IsActive:   isActive,
		StockMax:   stockMax,
		Sort:       "stock_asc",
		Page:       f.Page,
		Limit:      f.Limit,
	}

	products, total, err := s.repo.ListProductsAdmin(adminF)
	if err != nil {
		return response.Paginated[productdto.InventoryItem]{}, err
	}

	var items []productdto.InventoryItem
	for _, p := range products {
		catName := ""
		if p.Category.ID.String() != "" {
			catName = p.Category.Name
		}
		items = append(items, productdto.InventoryItem{
			ID:       p.ID.String(),
			Name:     p.Name,
			SKU:      p.SKU,
			Category: catName,
			Stock:    p.Stock,
			Price:    p.Price,
			IsActive: p.IsActive,
			ImageURL: p.ImageURL,
		})
	}

	return response.NewPaginated(items, total, f.Page, f.Limit), nil
}

func (s *inventoryService) AdjustStock(productID string, req productdto.AdjustStockRequest) (*productdto.InventoryItem, error) {
	p, err := s.repo.FindProductByIDAdmin(productID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProductNotFound
	}

	log := &productmodel.StockLog{
		ProductID: p.ID,
		Quantity:  req.Quantity,
		CostPrice: req.CostPrice,
		Note:      req.Note,
	}

	switch req.Operation {
	case "add":
		p.Stock += req.Quantity
		log.Type = productmodel.StockLogPurchase
		log.TotalCost = req.CostPrice * float64(req.Quantity)
		if req.CostPrice > 0 {
			totalValue := (p.CostPrice * float64(p.Stock-req.Quantity)) + (req.CostPrice * float64(req.Quantity))
			p.CostPrice = math.Floor(totalValue / float64(p.Stock))
		}
	case "subtract":
		if p.Stock < req.Quantity {
			return nil, fmt.Errorf("insufficient stock: have %d, want %d", p.Stock, req.Quantity)
		}
		p.Stock -= req.Quantity
		log.Type = productmodel.StockLogAdjustment
		log.Quantity = -req.Quantity
	case "set":
		diff := req.Quantity - p.Stock
		p.Stock = req.Quantity
		log.Type = productmodel.StockLogAdjustment
		log.Quantity = diff
	default:
		return nil, ErrInvalidOperation
	}

	if err := s.repo.UpdateProduct(p); err != nil {
		return nil, err
	}
	s.repo.CreateStockLog(log)

	p, _ = s.repo.FindProductByIDAdmin(productID)
	catName := ""
	if p.Category.ID.String() != "" {
		catName = p.Category.Name
	}
	return &productdto.InventoryItem{
		ID: p.ID.String(), Name: p.Name, SKU: p.SKU, Category: catName,
		Stock: p.Stock, Price: p.Price, IsActive: p.IsActive, ImageURL: p.ImageURL,
	}, nil
}
