import api from './axios'
import type { Address, AddressCreatePayload, AddressUpdatePayload, ApiResponse } from '@/types'

export const addressesApi = {
  list: () =>
    api.get<ApiResponse<Address[]>>('/customer/addresses'),

  create: (data: AddressCreatePayload) =>
    api.post<ApiResponse<Address>>('/customer/addresses', data),

  update: (id: string, data: AddressUpdatePayload) =>
    api.put<ApiResponse<Address>>(`/customer/addresses/${id}`, data),

  delete: (id: string) =>
    api.delete<ApiResponse<null>>(`/customer/addresses/${id}`),

  setDefault: (id: string) =>
    api.patch<ApiResponse<Address>>(`/customer/addresses/${id}/default`),
}
