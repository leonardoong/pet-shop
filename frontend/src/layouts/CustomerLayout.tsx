import { Outlet, Link, useNavigate, NavLink } from 'react-router-dom'
import { ShoppingCart, User, LogOut, Menu, X, Search } from 'lucide-react'
import { useState, useEffect, useRef } from 'react'
import { useAuthStore } from '@/store/authStore'
import { useCartStore } from '@/store/cartStore'
import { useQuery } from '@tanstack/react-query'
import { categoriesApi } from '@/api/products'
import { cn } from '@/lib/utils'

export default function CustomerLayout() {
  const { isAuthenticated, customer, logout } = useAuthStore()
  const { fetch: fetchCart, cart } = useCartStore()
  const navigate = useNavigate()
  const [menuOpen, setMenuOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  const searchRef = useRef<HTMLInputElement>(null)

  const { data: catData } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoriesApi.list(),
    staleTime: 5 * 60_000,
  })

  const raw = catData?.data?.data
  const categories = Array.isArray(raw) ? raw : []
  const navLinks = [
    { to: '/shop', label: 'Semua Produk' },
    ...categories.slice(0, 6).map((c: { slug: string; name: string }) => ({
      to: `/shop?category=${c.slug}`,
      label: c.name,
    })),
  ]

  const itemCount = cart?.items?.reduce((sum, i) => sum + i.quantity, 0) ?? 0

  // Fetch cart when user is authenticated
  useEffect(() => {
    if (isAuthenticated) fetchCart()
  }, [isAuthenticated, fetchCart])

  const handleLogout = () => {
    logout()
    useCartStore.getState().clear()
    navigate('/')
  }

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchQuery.trim()) {
      navigate(`/shop?search=${encodeURIComponent(searchQuery.trim())}`)
      setMenuOpen(false)
    }
  }

  return (
    <div className="min-h-screen flex flex-col bg-warm-50">
      <header className="sticky top-0 z-50 bg-primary-500 shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          {/* Top row */}
          <div className="flex items-center gap-4 h-16">
            {/* Logo */}
            <Link to="/" className="flex items-center gap-2 font-bold text-xl text-white shrink-0">
              🐾 <span>PetShop</span>
            </Link>

            {/* Search bar */}
            <form onSubmit={handleSearch} className="flex-1 hidden sm:flex items-center max-w-lg relative">
              <Search className="absolute left-3 w-4 h-4 text-primary-400 pointer-events-none" />
              <input
                ref={searchRef}
                type="search"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Cari produk untuk hewan peliharaan..."
                className="w-full bg-white/20 placeholder-white/70 text-white border border-white/30 rounded-xl pl-9 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-white/50 focus:bg-white/30 transition-colors"
              />
            </form>

            {/* Actions */}
            <div className="flex items-center gap-2 ml-auto">
              {/* Cart with badge */}
              <Link
                to="/keranjang"
                className="p-2 rounded-xl hover:bg-white/20 transition-colors relative"
                aria-label="Keranjang belanja"
              >
                <ShoppingCart className="w-5 h-5 text-white" />
                {itemCount > 0 && (
                  <span className="absolute -top-1 -right-1 min-w-[18px] h-[18px] flex items-center justify-center bg-red-500 text-white text-[10px] font-bold rounded-full px-1 leading-none">
                    {itemCount > 99 ? '99+' : itemCount}
                  </span>
                )}
              </Link>

              {isAuthenticated ? (
                <>
                  <Link
                    to="/akun/alamat"
                    className="hidden md:flex items-center gap-1.5 text-sm text-white/90 hover:text-white px-3 py-2 rounded-xl hover:bg-white/20 transition-colors"
                  >
                    <User className="w-4 h-4" />
                    {customer?.full_name.split(' ')[0]}
                  </Link>
                  <Link
                    to="/pesanan"
                    className="hidden md:block text-sm text-white/80 hover:text-white px-3 py-2 rounded-xl hover:bg-white/20 transition-colors"
                  >
                    Pesanan
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
          
        </div>

        {/* Mobile menu */}
        {menuOpen && (
          <div className="md:hidden border-t border-white/20 bg-primary-600 px-4 py-3 space-y-1 text-sm">
            <form onSubmit={handleSearch} className="flex items-center relative mb-3">
              <Search className="absolute left-3 w-4 h-4 text-white/50" />
              <input
                type="search"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Cari produk..."
                className="w-full bg-white/20 placeholder-white/60 text-white border border-white/30 rounded-xl pl-9 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-white/40"
              />
            </form>
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
                  <Link to="/pesanan" onClick={() => setMenuOpen(false)} className="block px-3 py-2 rounded-lg text-white/90 hover:bg-white/20">Pesanan Saya</Link>
                  <Link to="/akun/alamat" onClick={() => setMenuOpen(false)} className="block px-3 py-2 rounded-lg text-white/90 hover:bg-white/20">Buku Alamat</Link>
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
              {categories.slice(0, 5).map((c: { slug: string; name: string }) => (
                <li key={c.slug}><Link to={`/shop?category=${c.slug}`} className="hover:text-white transition-colors">{c.name}</Link></li>
              ))}
            </ul>
          </div>
          <div>
            <p className="text-white font-semibold mb-2">Bantuan</p>
            <ul className="space-y-1 text-xs text-gray-500">
              <li><span className="text-gray-600">Tentang Kami</span></li>
              <li><span className="text-gray-600">Hubungi Kami</span></li>
              <li><span className="text-gray-600">Kebijakan Privasi</span></li>
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
