package dashboard

import (
	"time"

	dashdto "petshop/internal/dto/dashboard"

	"gorm.io/gorm"
)

type Repository interface {
	GetStats() (*dashdto.Stats, error)
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) GetStats() (*dashdto.Stats, error) {
	s := &dashdto.Stats{
		OrdersByStatus: make(map[string]int64),
	}

	r.db.Table("customers").Where("is_active = true").Count(&s.TotalCustomers)
	r.db.Table("products").Count(&s.TotalProducts)

	var totalRevenue float64
	r.db.Table("orders").
		Where("status != ?", "cancelled").
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalRevenue)
	s.TotalRevenue = int64(totalRevenue)

	var totalCOGS float64
	r.db.Table("stock_logs").
		Where("type = ?", "sale").
		Select("COALESCE(SUM(total_cost), 0)").
		Scan(&totalCOGS)
	s.TotalCOGS = int64(totalCOGS)
	s.NetIncome = s.TotalRevenue - s.TotalCOGS

	if s.TotalRevenue > 0 {
		s.ProfitMargin = float64(s.NetIncome) / float64(s.TotalRevenue) * 100
	}

	var inventoryValue float64
	r.db.Table("products").
		Select("COALESCE(SUM(cost_price * stock), 0)").
		Scan(&inventoryValue)
	s.InventoryValue = int64(inventoryValue)

	r.db.Table("orders").Count(&s.TotalOrders)

	type statusRow struct {
		Status string
		Count  int64
	}
	var statusRows []statusRow
	r.db.Table("orders").
		Select("status::text AS status, COUNT(*) AS count").
		Group("status").
		Find(&statusRows)
	for _, row := range statusRows {
		s.OrdersByStatus[row.Status] = row.Count
	}

	type recentRow struct {
		ID          string
		Status      string
		TotalAmount float64
		ShipName    string
		CreatedAt   time.Time
	}
	var recentRows []recentRow
	r.db.Table("orders").
		Select("id::text, status::text, total_amount, ship_name, created_at").
		Order("created_at desc").
		Limit(5).
		Find(&recentRows)
	for _, row := range recentRows {
		s.RecentOrders = append(s.RecentOrders, dashdto.RecentOrder{
			ID:           row.ID,
			CustomerName: row.ShipName,
			Status:       row.Status,
			TotalAmount:  row.TotalAmount,
			CreatedAt:    row.CreatedAt,
		})
	}

	type lowStockRow struct {
		ID    string
		Name  string
		SKU   string
		Stock int
	}
	var lowRows []lowStockRow
	r.db.Table("products").
		Select("id::text, name, sku, stock").
		Where("stock <= 10 AND is_active = true").
		Order("stock asc").
		Limit(5).
		Find(&lowRows)
	for _, row := range lowRows {
		s.LowStockProducts = append(s.LowStockProducts, dashdto.LowStockItem{
			ID: row.ID, Name: row.Name, SKU: row.SKU, Stock: row.Stock,
		})
	}

	type monthlyRow struct {
		Month   string
		Revenue float64
	}
	var monthly []monthlyRow
	r.db.Raw(`SELECT TO_CHAR(created_at, 'YYYY-MM') as month, COALESCE(SUM(total_amount), 0) as revenue
		FROM orders WHERE status != ? AND created_at >= ? 
		GROUP BY month ORDER BY month ASC`, "cancelled", time.Now().AddDate(0, -6, 0)).Scan(&monthly)
	for _, m := range monthly {
		s.RevenueByMonth = append(s.RevenueByMonth, dashdto.MonthlyRevenue{
			Month: m.Month, Revenue: m.Revenue,
		})
	}

	return s, nil
}
