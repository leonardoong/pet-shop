import api from './axios'
import type { ApiResponse, PaginatedResponse, Product, CreateProduct, UpdateProduct } from '@/types'

export const productsApi = {
  list: (params?: Record<string, string>) =>
    api.get<ApiResponse<PaginatedResponse<Product>>>('/admin/products', { params }),

  get: (id: string) =>
    api.get<ApiResponse<Product>>(`/admin/products/${id}`),

  create: (data: CreateProduct) =>
    api.post<ApiResponse<Product>>('/admin/products', data),

  update: (id: string, data: UpdateProduct) =>
    api.put<ApiResponse<Product>>(`/admin/products/${id}`, data),

  delete: (id: string) =>
    api.delete<ApiResponse<null>>(`/admin/products/${id}`),
}
