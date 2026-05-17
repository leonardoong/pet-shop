import api from './axios'
import type { ApiResponse, Order, OrderStatus, Paginated } from '@/types'

export const ordersApi = {
  checkout: (address_id: string, notes?: string) =>
    api.post<ApiResponse<Order>>('/customer/orders', { address_id, notes }),

  list: (params?: { status?: OrderStatus; page?: number; limit?: number }) =>
    api.get<ApiResponse<Paginated<Order>>>('/customer/orders', { params }),

  getById: (id: string) =>
    api.get<ApiResponse<Order>>(`/customer/orders/${id}`),
}
