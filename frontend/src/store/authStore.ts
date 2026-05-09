import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { authApi } from '@/api/auth'
import type { Customer } from '@/types'

interface AuthState {
  customer: Customer | null
  accessToken: string | null
  refreshToken: string | null
  isAuthenticated: boolean

  setAuth: (customer: Customer, accessToken: string, refreshToken: string) => void
  refresh: () => Promise<void>
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      customer: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,

      setAuth: (customer, accessToken, refreshToken) =>
        set({ customer, accessToken, refreshToken, isAuthenticated: true }),

      refresh: async () => {
        const { refreshToken } = get()
        if (!refreshToken) throw new Error('No refresh token')
        const res = await authApi.refresh(refreshToken)
        set({
          accessToken: res.data.data.access_token,
          // refresh token rotation: backend returns a new one in the next phase
        })
      },

      logout: () => {
        const { refreshToken } = get()
        if (refreshToken) {
          authApi.logout(refreshToken).catch(() => {})
        }
        set({ customer: null, accessToken: null, refreshToken: null, isAuthenticated: false })
      },
    }),
    {
      name: 'petshop-auth',
      // Only persist the refresh token; access token is kept in memory per session
      partialize: (state) => ({ refreshToken: state.refreshToken, customer: state.customer }),
    },
  ),
)
