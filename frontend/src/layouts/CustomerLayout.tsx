import { Outlet, Link, useNavigate, NavLink } from 'react-router-dom'
import { ShoppingCart, User, LogOut, Menu, X, Search } from 'lucide-react'
import { useState } from 'react'
import { useAuthStore } from '@/store/authStore'
import { cn } from '@/lib/utils'

const navLinks = [
  { to: '/shop',                 label: 'Semua Produk' },
  { to: '/shop?category=anjing', label: 'Anjing' },
  { to: '/shop?category=kucing', label: 'Kucing' },
]

export default function CustomerLayout() {
  const { isAuthenticated, customer, logout } = useAuthStore()
  const navigate = useNavigate()
  const [menuOpen, setMenuOpen] = useState(false)

  const handleLogout = () => {
    logout()
    navigate('/')
  }

  return (
    <div className="min-h-screen flex flex-col bg-warm-50">
      {/* Header merges into hero — warm green background */}
      <header className="sticky top-0 z-50 bg-primary-500 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          {/* Top row: logo, search, actions */}
          <div className="flex items-center gap-4 h-16">
            {/* Logo */}
            <Link to="/" className="flex items-center gap-2 font-bold text-xl text-white shrink-0">
              🐾 <span>PetShop</span>
            </Link>

            {/* Search bar */}
            <div className="flex-1 hidden sm:flex items-center max-w-lg relative">
              <Search className="absolute left-3 w-4 h-4 text-primary-400 pointer-events-none" />
              <input
                type="text"
                placeholder="Cari produk untuk hewan peliharaan..."
                className="w-full bg-white/20 placeholder-white/70 text-white border border-white/30 rounded-xl pl-9 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-white/50 focus:bg-white/30 transition-colors"
              />
            </div>

            {/* Actions */}
            <div className="flex items-center gap-2 ml-auto">
              <Link
                to="/keranjang"
                className="p-2 rounded-xl hover:bg-white/20 transition-colors relative"
                aria-label="Keranjang belanja"
              >
                <ShoppingCart className="w-5 h-5 text-white" />
              </Link>

              {isAuthenticated ? (
                <>
                  <Link
                    to="/akun"
                    className="hidden md:flex items-center gap-1.5 text-sm text-white/90 hover:text-white px-3 py-2 rounded-xl hover:bg-white/20 transition-colors"
                  >
                    <User className="w-4 h-4" />
                    {customer?.full_name.split(' ')[0]}
                  </Link>
                  <button
                    onClick={handleLogout}
                    className="hidden md:flex items-center gap-1 p-2 rounded-xl hover:bg-white/20 text-white/80 hover:text-white transition-colors"
                    aria-label="Keluar"
                  >
                    <LogOut className="w-4 h-4" />
                  </button>
                </>
              ) : (
                <div className="hidden md:flex items-center gap-2">
                  <Link
                    to="/masuk"
                    className="text-sm text-white/90 hover:text-white font-medium px-3 py-2 rounded-xl hover:bg-white/20 transition-colors"
                  >
                    Masuk
                  </Link>
                  <Link
                    to="/daftar"
                    className="text-sm bg-white text-primary-600 hover:bg-primary-50 px-4 py-2 rounded-xl font-semibold transition-colors"
                  >
                    Daftar
                  </Link>
                </div>
              )}

              <button
                className="md:hidden p-2 rounded-xl hover:bg-white/20 text-white"
                onClick={() => setMenuOpen((o) => !o)}
                aria-label="Menu"
              >
                {menuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
              </button>
            </div>
          </div>

          {/* Category nav row */}
          <nav className="hidden md:flex items-center gap-1 pb-3">
            {navLinks.map(({ to, label }) => (
              <NavLink
                key={to}
                to={to}
                className={({ isActive }) =>
                  cn(
                    'text-sm px-4 py-1.5 rounded-lg font-medium whitespace-nowrap transition-colors',
                    isActive
                      ? 'bg-white/30 text-white'
                      : 'text-white/80 hover:bg-white/20 hover:text-white',
                  )
                }
              >
                {label}
              </NavLink>
            ))}
          </nav>
        </div>

        {/* Mobile menu */}
        {menuOpen && (
          <div className="md:hidden border-t border-white/20 bg-primary-600 px-4 py-3 space-y-1 text-sm">
            <div className="flex items-center relative mb-3">
              <Search className="absolute left-3 w-4 h-4 text-white/50" />
              <input
                type="text"
                placeholder="Cari produk..."
                className="w-full bg-white/20 placeholder-white/60 text-white border border-white/30 rounded-xl pl-9 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-white/40"
              />
            </div>
            {navLinks.map(({ to, label }) => (
              <Link key={to} to={to} onClick={() => setMenuOpen(false)}
                className="block px-3 py-2 rounded-lg text-white/90 hover:bg-white/20"
              >
                {label}
              </Link>
            ))}
            <div className="border-t border-white/20 pt-2 mt-2">
              {isAuthenticated ? (
                <>
                  <Link to="/akun" onClick={() => setMenuOpen(false)} className="block px-3 py-2 rounded-lg text-white/90 hover:bg-white/20">Akun Saya</Link>
                  <button onClick={handleLogout} className="block w-full text-left px-3 py-2 rounded-lg text-red-200 hover:bg-white/10">Keluar</button>
                </>
              ) : (
                <>
                  <Link to="/masuk" onClick={() => setMenuOpen(false)} className="block px-3 py-2 rounded-lg text-white/90 hover:bg-white/20">Masuk</Link>
                  <Link to="/daftar" onClick={() => setMenuOpen(false)} className="block px-3 py-2 rounded-lg text-white font-semibold hover:bg-white/20">Daftar</Link>
                </>
              )}
            </div>
          </div>
        )}
      </header>

      <main className="flex-1">
        <Outlet />
      </main>

      <footer className="bg-gray-900 text-gray-400 text-sm mt-16">
        <div className="max-w-7xl mx-auto px-4 py-10 grid grid-cols-1 sm:grid-cols-3 gap-8">
          <div>
            <p className="text-white font-bold text-lg mb-2">🐾 PetShop</p>
            <p className="text-gray-500 text-xs leading-relaxed">
              Toko kebutuhan hewan peliharaan terpercaya. Produk berkualitas, harga terjangkau.
            </p>
          </div>
          <div>
            <p className="text-white font-semibold mb-2">Kategori</p>
            <ul className="space-y-1 text-xs text-gray-500">
              <li><Link to="/shop?category=anjing" className="hover:text-white transition-colors">Anjing</Link></li>
              <li><Link to="/shop?category=kucing" className="hover:text-white transition-colors">Kucing</Link></li>
            </ul>
          </div>
          <div>
            <p className="text-white font-semibold mb-2">Bantuan</p>
            <ul className="space-y-1 text-xs text-gray-500">
              <li><Link to="/tentang" className="hover:text-white transition-colors">Tentang Kami</Link></li>
              <li><Link to="/kontak" className="hover:text-white transition-colors">Hubungi Kami</Link></li>
              <li><Link to="/kebijakan-privasi" className="hover:text-white transition-colors">Kebijakan Privasi</Link></li>
            </ul>
          </div>
        </div>
        <div className="border-t border-gray-800 py-4 text-center text-xs text-gray-600">
          © {new Date().getFullYear()} PetShop. Hak cipta dilindungi.
        </div>
      </footer>
    </div>
  )
}
