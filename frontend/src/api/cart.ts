import api from './axios'
import type { ApiResponse, Cart } from '@/types'

export const cartApi = {
  get: () =>
    api.get<ApiResponse<Cart>>('/customer/cart'),

  addItem: (product_id: string, quantity: number) =>
    api.post<ApiResponse<Cart>>('/customer/cart/items', { product_id, quantity }),

  updateItem: (productId: string, quantity: number) =>
    api.put<ApiResponse<Cart>>(`/customer/cart/items/${productId}`, { quantity }),

  removeItem: (productId: string) =>
    api.delete<ApiResponse<null>>(`/customer/cart/items/${productId}`),
}
