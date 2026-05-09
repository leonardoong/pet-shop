import { createBrowserRouter, Navigate } from 'react-router-dom'
import CustomerLayout from '@/layouts/CustomerLayout'
import Home from '@/pages/Home'
import Login from '@/pages/Login'
import Register from '@/pages/Register'
import { useAuthStore } from '@/store/authStore'

function RequireAuth({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  if (!isAuthenticated) return <Navigate to="/masuk" replace />
  return <>{children}</>
}

const placeholder = (label: string) => (
  <div className="max-w-7xl mx-auto px-4 py-16 text-center text-gray-400">{label}</div>
)

export const router = createBrowserRouter([
  {
    path: '/',
    element: <CustomerLayout />,
    children: [
      { index: true,        element: <Home /> },
      { path: 'masuk',      element: <Login /> },
      { path: 'daftar',     element: <Register /> },
      { path: 'shop',       element: placeholder('Katalog produk — hadir di Phase 2') },
      { path: 'produk/:id', element: placeholder('Detail produk — hadir di Phase 2') },
      {
        path: 'keranjang',
        element: <RequireAuth>{placeholder('Keranjang belanja — hadir di Phase 2')}</RequireAuth>,
      },
      {
        path: 'akun',
        element: <RequireAuth>{placeholder('Akun saya — hadir di Phase 4')}</RequireAuth>,
      },
    ],
  },
])
