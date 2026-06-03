import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import {
  LayoutDashboard, Package, ShoppingCart, Users,
  Warehouse, Settings, LogOut, ChevronRight, Tags,
} from 'lucide-react'
import { useAdminAuthStore } from '@/store/authStore'
import { cn } from '@/lib/utils'

const navItems = [
  { to: '/dashboard',  label: 'Dashboard',  icon: LayoutDashboard, permission: 'dashboard:read' },
  { to: '/products',   label: 'Products',   icon: Package,         permission: 'products:read' },
  { to: '/categories', label: 'Categories', icon: Tags,            permission: 'categories:read' },
  { to: '/orders',     label: 'Orders',     icon: ShoppingCart,    permission: 'orders:read' },
  { to: '/inventory',  label: 'Inventory',  icon: Warehouse,       permission: 'inventory:read' },
  { to: '/customers',  label: 'Customers',  icon: Users,           permission: 'customers:read' },
  { to: '/settings',   label: 'Settings',   icon: Settings,        permission: 'admins:read' },
]

export default function AdminLayout() {
  const { admin, hasPermission, logout } = useAdminAuthStore()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const visibleItems = navItems.filter((item) => hasPermission(item.permission))

  return (
    <div className="flex h-screen bg-gray-50">
      {/* Sidebar */}
      <aside className="w-60 bg-gray-900 text-gray-300 flex flex-col shrink-0">
        {/* Logo */}
        <div className="px-6 py-5 border-b border-gray-800">
          <span className="text-white font-bold text-lg">🐾 PetShop</span>
          <span className="ml-2 text-xs text-gray-500 font-medium uppercase tracking-wider">Admin</span>
        </div>

        {/* Nav */}
        <nav className="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
          {visibleItems.map(({ to, label, icon: Icon }) => (
            <NavLink
              key={to}
              to={to}
              className={({ isActive }) =>
                cn(
                  'flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors',
                  isActive
                    ? 'bg-gray-800 text-white'
                    : 'hover:bg-gray-800 hover:text-white',
                )
              }
            >
              <Icon className="w-4 h-4 shrink-0" />
              {label}
              <ChevronRight className="w-3 h-3 ml-auto opacity-40" />
            </NavLink>
          ))}
        </nav>

        {/* Admin info + logout */}
        <div className="px-4 py-4 border-t border-gray-800">
          <div className="text-xs text-gray-500 truncate mb-1">{admin?.email}</div>
          <div className="text-sm text-gray-300 font-medium truncate mb-3">{admin?.full_name}</div>
          <button
            onClick={handleLogout}
            className="flex items-center gap-2 text-sm text-gray-400 hover:text-red-400 transition-colors"
          >
            <LogOut className="w-4 h-4" />
            Sign out
          </button>
        </div>
      </aside>

      {/* Main content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        <header className="bg-white border-b border-gray-200 px-6 py-4 shrink-0">
          <h1 className="text-lg font-semibold text-gray-800">Admin Dashboard</h1>
        </header>
        <main className="flex-1 overflow-y-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
