import api from './axios'
import type { ApiResponse, PaginatedResponse, AdminOrder, UpdateOrderStatus } from '@/types'

export const ordersApi = {
  list: (params?: Record<string, string>) =>
    api.get<ApiResponse<PaginatedResponse<AdminOrder>>>('/admin/orders', { params }),

  get: (id: string) =>
    api.get<ApiResponse<AdminOrder>>(`/admin/orders/${id}`),

  updateStatus: (id: string, data: UpdateOrderStatus) =>
    api.patch<ApiResponse<AdminOrder>>(`/admin/orders/${id}/status`, data),
}
