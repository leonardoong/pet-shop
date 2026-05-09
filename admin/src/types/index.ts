export interface Admin {
  id: string
  email: string
  full_name: string
  is_active: boolean
}

export interface TokenPair {
  access_token: string
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
