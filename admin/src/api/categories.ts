import api from './axios'
import type { ApiResponse, Category, CreateCategory } from '@/types'

export const categoriesApi = {
  list: () =>
    api.get<ApiResponse<Category[]>>('/categories'),

  create: (data: CreateCategory) =>
    api.post<ApiResponse<Category>>('/admin/categories', data),

  update: (id: string, data: Partial<CreateCategory>) =>
    api.put<ApiResponse<Category>>(`/admin/categories/${id}`, data),

  delete: (id: string) =>
    api.delete<ApiResponse<null>>(`/admin/categories/${id}`),
}
