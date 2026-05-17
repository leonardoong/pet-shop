import { createBrowserRouter, Navigate } from 'react-router-dom'
import CustomerLayout from '@/layouts/CustomerLayout'
import Home from '@/pages/Home'
import Login from '@/pages/Login'
import Register from '@/pages/Register'
import Shop from '@/pages/Shop'
import ProductDetail from '@/pages/ProductDetail'
import Cart from '@/pages/Cart'
import Checkout from '@/pages/Checkout'
import Orders from '@/pages/Orders'
import OrderDetail from '@/pages/OrderDetail'
import AddressBook from '@/pages/AddressBook'
import { useAuthStore } from '@/store/authStore'

function RequireAuth({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  if (!isAuthenticated) return <Navigate to="/masuk" replace />
  return <>{children}</>
}

export const router = createBrowserRouter([
  {
    path: '/',
    element: <CustomerLayout />,
    children: [
      { index: true,              element: <Home /> },
      { path: 'masuk',            element: <Login /> },
      { path: 'daftar',           element: <Register /> },
      { path: 'shop',             element: <Shop /> },
      { path: 'produk/:slug',     element: <ProductDetail /> },
      {
        path: 'keranjang',
        element: <RequireAuth><Cart /></RequireAuth>,
      },
      {
        path: 'checkout',
        element: <RequireAuth><Checkout /></RequireAuth>,
      },
      {
        path: 'pesanan',
        element: <RequireAuth><Orders /></RequireAuth>,
      },
      {
        path: 'pesanan/:id',
        element: <RequireAuth><OrderDetail /></RequireAuth>,
      },
      {
        path: 'akun/alamat',
        element: <RequireAuth><AddressBook /></RequireAuth>,
      },
    ],
  },
])
