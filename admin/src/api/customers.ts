import api from './axios'
import type { ApiResponse, PaginatedResponse, Customer } from '@/types'

export const customersApi = {
  list: (params?: Record<string, string>) =>
    api.get<ApiResponse<PaginatedResponse<Customer>>>('/admin/customers', { params }),

  get: (id: string) =>
    api.get<ApiResponse<Customer>>(`/admin/customers/${id}`),

  toggleActive: (id: string, is_active: boolean) =>
    api.patch<ApiResponse<Customer>>(`/admin/customers/${id}`, { is_active }),
}
