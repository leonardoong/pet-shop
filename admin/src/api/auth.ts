import axios from 'axios'
import type { ApiResponse, AdminAuthResponse, TokenPair } from '@/types'

const base = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
})

export const authApi = {
  login: (data: { email: string; password: string }) =>
    base.post<ApiResponse<AdminAuthResponse>>('/admin/auth/login', data),

  refresh: (refreshToken: string) =>
    base.post<ApiResponse<TokenPair>>('/admin/auth/refresh', { refresh_token: refreshToken }),

  logout: (refreshToken: string) =>
    base.post('/admin/auth/logout', { refresh_token: refreshToken }),
}
