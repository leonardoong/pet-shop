import axios from 'axios'
import { useAdminAuthStore } from '@/store/authStore'

const api = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use((config) => {
  const token = useAdminAuthStore.getState().accessToken
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  (res) => res,
  async (error) => {
    const original = error.config
    if (error.response?.status === 401 && !original._retry) {
      original._retry = true
      try {
        await useAdminAuthStore.getState().refresh()
        const newToken = useAdminAuthStore.getState().accessToken
        original.headers.Authorization = `Bearer ${newToken}`
        return api(original)
      } catch {
        useAdminAuthStore.getState().logout()
      }
    }
    return Promise.reject(error)
  },
)

export default api
