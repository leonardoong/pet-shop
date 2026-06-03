package dashboard

import "time"

type Stats struct {
	TotalRevenue     int64            `json:"total_revenue"`
	TotalCOGS        int64            `json:"total_cogs"`
	NetIncome        int64            `json:"net_income"`
	ProfitMargin     float64          `json:"profit_margin"`
	TotalOrders      int64            `json:"total_orders"`
	TotalProducts    int64            `json:"total_products"`
	TotalCustomers   int64            `json:"total_customers"`
	InventoryValue   int64            `json:"inventory_value"`
	OrdersByStatus   map[string]int64 `json:"orders_by_status"`
	RecentOrders     []RecentOrder    `json:"recent_orders"`
	LowStockProducts []LowStockItem   `json:"low_stock_products"`
	RevenueByMonth   []MonthlyRevenue `json:"revenue_by_month"`
}

type RecentOrder struct {
	ID           string    `json:"id"`
	CustomerName string    `json:"customer_name"`
	Status       string    `json:"status"`
	TotalAmount  float64   `json:"total_amount"`
	CreatedAt    time.Time `json:"created_at"`
}

type LowStockItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SKU      string `json:"sku"`
	Stock    int    `json:"stock"`
}

type MonthlyRevenue struct {
	Month   string  `json:"month"`
	Revenue float64 `json:"revenue"`
}
