import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { authApi } from '@/api/auth'
import type { Admin } from '@/types'

interface AdminAuthState {
  admin: Admin | null
  accessToken: string | null
  refreshToken: string | null
  permissions: string[]
  isAuthenticated: boolean

  setAuth: (admin: Admin, accessToken: string, refreshToken: string, permissions: string[]) => void
  hasPermission: (permission: string) => boolean
  refresh: () => Promise<void>
  logout: () => void
}

export const useAdminAuthStore = create<AdminAuthState>()(
  persist(
    (set, get) => ({
      admin: null,
      accessToken: null,
      refreshToken: null,
      permissions: [],
      isAuthenticated: false,

      setAuth: (admin, accessToken, refreshToken, permissions) =>
        set({ admin, accessToken, refreshToken, permissions, isAuthenticated: true }),

      hasPermission: (permission) => get().permissions.includes(permission),

      refresh: async () => {
        const { refreshToken } = get()
        if (!refreshToken) throw new Error('No refresh token')
        const res = await authApi.refresh(refreshToken)
        const { access_token, refresh_token } = res.data.data
        set({
          accessToken: access_token,
          refreshToken: refresh_token || refreshToken,
        })
      },

      logout: () => {
        const { refreshToken } = get()
        if (refreshToken) authApi.logout(refreshToken).catch(() => {})
        set({ admin: null, accessToken: null, refreshToken: null, permissions: [], isAuthenticated: false })
      },
    }),
    {
      name: 'petshop-admin-auth',
      partialize: (state) => ({ refreshToken: state.refreshToken, admin: state.admin }),
    },
  ),
)
