import api from './axios'
import type { ApiResponse, PaginatedResponse, InventoryItem, AdjustStock } from '@/types'

export const inventoryApi = {
  list: (params?: Record<string, string>) =>
    api.get<ApiResponse<PaginatedResponse<InventoryItem>>>('/admin/inventory', { params }),

  adjustStock: (productId: string, data: AdjustStock) =>
    api.patch<ApiResponse<InventoryItem>>(`/admin/inventory/${productId}`, data),
}
