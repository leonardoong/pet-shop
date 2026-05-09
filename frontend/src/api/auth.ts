import axios from 'axios'
import type { ApiResponse, CustomerAuthResponse, TokenPair } from '@/types'

const base = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
})

export const authApi = {
  register: (data: { full_name: string; email: string; password: string; phone?: string }) =>
    base.post<ApiResponse<CustomerAuthResponse>>('/customer/auth/register', data),

  login: (data: { email: string; password: string }) =>
    base.post<ApiResponse<CustomerAuthResponse>>('/customer/auth/login', data),

  refresh: (refreshToken: string) =>
    base.post<ApiResponse<TokenPair>>('/customer/auth/refresh', { refresh_token: refreshToken }),

  logout: (refreshToken: string) =>
    base.post('/customer/auth/logout', { refresh_token: refreshToken }),
}
