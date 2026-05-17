import api from './axios'
import type { ApiResponse, Category, Paginated, Product, ProductFilter } from '@/types'

export const categoriesApi = {
  list: () =>
    api.get<ApiResponse<Category[]>>('/categories'),
}

export const productsApi = {
  list: (params?: ProductFilter) =>
    api.get<ApiResponse<Paginated<Product>>>('/products', { params }),

  getBySlug: (slug: string) =>
    api.get<ApiResponse<Product>>(`/products/${slug}`),
}
