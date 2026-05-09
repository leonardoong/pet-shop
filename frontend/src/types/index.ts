export interface Customer {
  id: string
  email: string
  full_name: string
  phone: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface TokenPair {
  access_token: string
  token_type: string
  expires_in: number
}

export interface CustomerAuthResponse {
  customer: Customer
  tokens: TokenPair
}

export interface ApiResponse<T> {
  success: boolean
  message: string
  data: T
}

export interface ApiError {
  success: false
  message: string
  errors?: string[]
}

export interface Product {
  id: string
  category_id: string
  name: string
  slug: string
  description: string
  price: number
  stock: number
  sku: string
  image_url: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Category {
  id: string
  name: string
  slug: string
  image_url: string
}

export interface CartItem {
  id: string
  product: Product
  quantity: number
}

export interface Cart {
  id: string
  items: CartItem[]
}

export interface Order {
  id: string
  status: 'pending' | 'confirmed' | 'processing' | 'shipped' | 'delivered' | 'cancelled'
  total_amount: number
  created_at: string
  updated_at: string
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  limit: number
  total_pages: number
}
