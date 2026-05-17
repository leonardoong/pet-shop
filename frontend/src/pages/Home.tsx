import { Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { ChevronRight, Truck, ShieldCheck, Star, Tag } from 'lucide-react'
import { categoriesApi, productsApi } from '@/api/products'
import ProductCard from '@/components/ProductCard'

const TRUST_BADGES = [
  { icon: Truck,       label: 'Gratis Ongkir',  desc: 'Pembelian ≥ Rp 150.000' },
  { icon: ShieldCheck, label: 'Produk Asli',    desc: '100% original & bergaransi' },
  { icon: Star,        label: '4.9/5 Rating',   desc: 'Dari 10.000+ pelanggan' },
  { icon: Tag,         label: 'Harga Terbaik',  desc: 'Selalu ada promo spesial' },
]

const CATEGORY_EMOJI: Record<string, string> = {
  anjing: '🐕',
  kucing: '🐈',
  ikan:   '🐟',
  burung: '🦜',
  hamster:'🐹',
}

export default function Home() {
  const { data: categories } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoriesApi.list().then((r) => r.data.data),
    staleTime: 5 * 60 * 1000,
  })

  const { data: featuredData, isLoading } = useQuery({
    queryKey: ['products', { page: 1, limit: 8, sort: 'newest' }],
    queryFn: () =>
      productsApi.list({ page: 1, limit: 8, sort: 'newest' }).then((r) => r.data.data),
    staleTime: 2 * 60 * 1000,
  })

  const products = featuredData?.items ?? []

  return (
    <div className="pb-10">

      {/* ── HERO ── */}
      <section className="bg-primary-500 pb-10 pt-4">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex flex-col sm:flex-row items-center gap-6">
          <div className="flex-1 text-white">
            <motion.p
              initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}
              className="text-primary-100 text-sm font-medium mb-1"
            >
              Selamat datang di PetShop 🐾
            </motion.p>
            <motion.h1
              initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.05 }}
              className="text-2xl sm:text-3xl font-bold mb-3 leading-snug"
            >
              Semua Kebutuhan<br />Hewan Peliharaanmu
            </motion.h1>
            <motion.div
              initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.1 }}
              className="flex gap-3 flex-wrap"
            >
              <Link to="/shop" className="bg-white text-primary-700 hover:bg-primary-50 font-semibold px-5 py-2 rounded-xl text-sm transition-colors">
                Belanja Sekarang
              </Link>
              <Link to="/shop?sort=newest" className="border border-white/40 hover:bg-white/10 text-white font-semibold px-5 py-2 rounded-xl text-sm transition-colors">
                Produk Baru
              </Link>
            </motion.div>
          </div>
          <div className="text-8xl select-none hidden sm:block">🐾</div>
        </div>
      </section>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-6">

        {/* ── TRUST BADGES ── */}
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 my-6">
          {TRUST_BADGES.map(({ icon: Icon, label, desc }) => (
            <div key={label} className="card flex items-center gap-3 p-3">
              <div className="w-9 h-9 bg-primary-100 rounded-lg flex items-center justify-center shrink-0">
                <Icon className="w-4 h-4 text-primary-600" />
              </div>
              <div>
                <p className="text-xs font-semibold text-gray-800">{label}</p>
                <p className="text-xs text-gray-500">{desc}</p>
              </div>
            </div>
          ))}
        </div>

        {/* ── KATEGORI ── */}
        <section className="mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-base font-bold text-gray-900">Kategori</h2>
            <Link to="/shop" className="text-sm text-primary-600 hover:underline flex items-center gap-0.5">
              Lihat semua <ChevronRight className="w-3.5 h-3.5" />
            </Link>
          </div>
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
            {categories
              ? categories.map((cat) => (
                  <Link
                    key={cat.id}
                    to={`/shop?category=${cat.slug}`}
                    className="flex items-center gap-3 p-4 rounded-2xl bg-white border border-gray-100 hover:border-primary-300 hover:bg-primary-50 transition-colors group"
                  >
                    <span className="text-3xl group-hover:scale-110 transition-transform">
                      {CATEGORY_EMOJI[cat.slug] ?? '🐾'}
                    </span>
                    <span className="text-sm font-semibold text-gray-700">{cat.name}</span>
                  </Link>
                ))
              : Array.from({ length: 4 }).map((_, i) => (
                  <div key={i} className="h-16 rounded-2xl bg-gray-100 animate-pulse" />
                ))}
          </div>
        </section>

        {/* ── PROMO BANNER ── */}
        <section className="mb-8">
          <div className="bg-gradient-to-r from-amber-400 to-orange-400 rounded-2xl p-5 flex items-center justify-between">
            <div>
              <p className="text-white/80 text-xs font-medium mb-0.5">Promo Spesial Hari Ini</p>
              <p className="text-white font-bold text-lg leading-tight">Diskon hingga 40%<br />untuk makanan kucing</p>
              <Link to="/shop?category=kucing" className="inline-block mt-3 bg-white text-orange-600 hover:bg-orange-50 font-semibold px-4 py-1.5 rounded-xl text-sm transition-colors">
                Lihat Promo
              </Link>
            </div>
            <div className="text-6xl select-none">🐈</div>
          </div>
        </section>

        {/* ── PRODUK UNGGULAN ── */}
        <section>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-base font-bold text-gray-900">Produk Terbaru</h2>
            <Link to="/shop" className="text-sm text-primary-600 hover:underline flex items-center gap-0.5">
              Lihat semua <ChevronRight className="w-3.5 h-3.5" />
            </Link>
          </div>

          {isLoading ? (
            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
              {Array.from({ length: 8 }).map((_, i) => (
                <div key={i} className="card overflow-hidden animate-pulse">
                  <div className="aspect-square bg-gray-100" />
                  <div className="p-3 space-y-2">
                    <div className="h-3 bg-gray-100 rounded w-16" />
                    <div className="h-3 bg-gray-100 rounded w-full" />
                    <div className="h-5 bg-gray-100 rounded w-24 mt-2" />
                    <div className="h-8 bg-gray-100 rounded-xl w-full mt-1" />
                  </div>
                </div>
              ))}
            </div>
          ) : products.length === 0 ? (
            <p className="text-gray-400 text-sm text-center py-10">
              Belum ada produk. Silakan tambahkan produk melalui dashboard admin.
            </p>
          ) : (
            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
              {products.map((product, i) => (
                <ProductCard key={product.id} product={product} index={i} />
              ))}
            </div>
          )}
        </section>

      </div>
    </div>
  )
}
