export interface Admin {
  id: string
  email: string
  full_name: string
  is_active: boolean
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
}

export interface AdminAuthResponse {
  admin: Admin
  tokens: TokenPair
  permissions: string[]
}

export interface ApiResponse<T> {
  success: boolean
  message: string
  data: T
}

export interface Role {
  id: string
  name: string
  description: string
  permissions: Permission[]
}

export interface Permission {
  id: string
  name: string
  description: string
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  limit: number
  total_pages: number
}

export interface Product {
  id: string
  category_id: string
  category: string
  name: string
  slug: string
  description: string
  price: number
  cost_price: number
  stock: number
  sku: string
  image_url: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateProduct {
  name: string
  category_id: string
  description: string
  price: number
  cost_price?: number
  stock: number
  sku: string
  image_url: string
  is_active?: boolean
}

export interface UpdateProduct {
  name?: string
  category_id?: string
  description?: string
  price?: number
  cost_price?: number
  stock?: number
  sku?: string
  image_url?: string
  is_active?: boolean
}

export interface Category {
  id: string
  name: string
  slug: string
  image_url: string
  created_at: string
  updated_at: string
}

export interface CreateCategory {
  name: string
  image_url: string
}

export interface AdminOrder {
  id: string
  customer_id: string
  customer_name: string
  status: string
  total_amount: number
  ship_name: string
  ship_phone: string
  ship_street: string
  ship_city: string
  ship_province: string
  ship_postal: string
  notes: string
  items: AdminOrderItem[]
  created_at: string
  updated_at: string
}

export interface AdminOrderItem {
  id: string
  product_id: string
  product_name: string
  quantity: number
  unit_price: number
  subtotal: number
}

export interface UpdateOrderStatus {
  status: string
  note?: string
}

export interface InventoryItem {
  id: string
  name: string
  sku: string
  category: string
  stock: number
  price: number
  is_active: boolean
  image_url: string
}

export interface AdjustStock {
  operation: 'add' | 'subtract' | 'set'
  quantity: number
  cost_price?: number
  note: string
}

export interface Customer {
  id: string
  full_name: string
  email: string
  phone: string
  is_active: boolean
  order_count: number
  total_spent: number
  created_at: string
}

export interface DashboardStats {
  total_revenue: number
  total_cogs: number
  net_income: number
  profit_margin: number
  total_orders: number
  total_products: number
  total_customers: number
  inventory_value: number
  orders_by_status: Record<string, number>
  recent_orders: {
    id: string
    customer_name: string
    status: string
    total_amount: number
    created_at: string
  }[]
  low_stock_products: {
    id: string
    name: string
    sku: string
    stock: number
  }[]
  revenue_by_month: {
    month: string
    revenue: number
  }[]
}
