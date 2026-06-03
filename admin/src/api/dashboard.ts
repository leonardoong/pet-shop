import api from './axios'
import type { ApiResponse, DashboardStats } from '@/types'

export const dashboardApi = {
  getStats: () =>
    api.get<ApiResponse<DashboardStats>>('/admin/dashboard'),
}
