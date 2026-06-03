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

export interface Paginated<T> {
  items: T[]
  total: number
  page: number
  limit: number
  total_pages: number
}

// --- Auth ---

export interface Customer {
  id: string
  email: string
  full_name: string
  phone: string
  avatar_url?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
}

export interface CustomerAuthResponse {
  customer: Customer
  tokens: TokenPair
}

// --- Catalog ---

export interface Category {
  id: string
  name: string
  slug: string
  image_url: string
}

export interface Product {
  id: string
  category_id: string
  category: Category
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

export interface ProductFilter {
  category?: string
  search?: string
  min_price?: number
  max_price?: number
  sort?: 'newest' | 'oldest' | 'price_asc' | 'price_desc'
  page?: number
  limit?: number
}

// --- Cart ---

export interface CartItem {
  id: string
  product_id: string
  product: Product
  quantity: number
}

export interface Cart {
  id: string
  customer_id: string
  items: CartItem[] | null
  created_at: string
  updated_at: string
}

// --- Address ---

export interface Address {
  id: string
  customer_id: string
  label: string
  recipient_name: string
  phone: string
  street: string
  city: string
  province: string
  postal_code: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface AddressCreatePayload {
  label?: string
  recipient_name: string
  phone: string
  street: string
  city: string
  province: string
  postal_code: string
  is_default?: boolean
}

export interface AddressUpdatePayload {
  label?: string
  recipient_name?: string
  phone?: string
  street?: string
  city?: string
  province?: string
  postal_code?: string
  is_default?: boolean
}

// --- Order ---

export type OrderStatus = 'pending' | 'confirmed' | 'processing' | 'shipped' | 'delivered' | 'cancelled'

export interface OrderItem {
  id: string
  product_id: string
  product: Product
  quantity: number
  unit_price: number
  subtotal: number
}

export interface Order {
  id: string
  customer_id: string
  status: OrderStatus
  total_amount: number
  ship_name: string
  ship_phone: string
  ship_street: string
  ship_city: string
  ship_province: string
  ship_postal: string
  notes: string
  items: OrderItem[]
  created_at: string
  updated_at: string
}
