import { createBrowserRouter, Navigate } from 'react-router-dom'
import AdminLayout from '@/layouts/AdminLayout'
import Login from '@/pages/Login'
import Dashboard from '@/pages/Dashboard'
import { useAdminAuthStore } from '@/store/authStore'

function RequireAdmin({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAdminAuthStore((s) => s.isAuthenticated)
  if (!isAuthenticated) return <Navigate to="/login" replace />
  return <>{children}</>
}

function RequirePermission({ permission, children }: { permission: string; children: React.ReactNode }) {
  const hasPermission = useAdminAuthStore((s) => s.hasPermission)
  if (!hasPermission(permission)) {
    return <div className="p-8 text-center text-gray-500">You don't have permission to view this page.</div>
  }
  return <>{children}</>
}

const placeholder = (label: string) => (
  <div className="p-8 text-center text-gray-400">{label} — coming in Phase 3</div>
)

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/',
    element: <RequireAdmin><AdminLayout /></RequireAdmin>,
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      {
        path: 'dashboard',
        element: <RequirePermission permission="dashboard:read"><Dashboard /></RequirePermission>,
      },
      {
        path: 'products',
        element: <RequirePermission permission="products:read">{placeholder('Product Management')}</RequirePermission>,
      },
      {
        path: 'orders',
        element: <RequirePermission permission="orders:read">{placeholder('Order Management')}</RequirePermission>,
      },
      {
        path: 'inventory',
        element: <RequirePermission permission="inventory:read">{placeholder('Inventory Management')}</RequirePermission>,
      },
      {
        path: 'customers',
        element: <RequirePermission permission="customers:read">{placeholder('Customer List')}</RequirePermission>,
      },
      {
        path: 'settings',
        element: <RequirePermission permission="admins:read">{placeholder('Settings & Role Management')}</RequirePermission>,
      },
    ],
  },
])
