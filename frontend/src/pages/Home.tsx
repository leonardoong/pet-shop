import { Link } from 'react-router-dom'
import { motion } from 'framer-motion'
import { ChevronRight, Tag, Truck, ShieldCheck, Star } from 'lucide-react'

// Placeholder product data — akan diganti data API pada Phase 2
const PLACEHOLDER_PRODUCTS = [
  { id: '1', name: 'Royal Canin Adult Dog', category: 'Anjing', price: 185000, image: 'https://placehold.co/300x300/e5f4da/3a6820?text=Makanan+Anjing', badge: 'Terlaris' },
  { id: '2', name: 'Whiskas Ocean Fish', category: 'Kucing', price: 45000, image: 'https://placehold.co/300x300/e5f4da/3a6820?text=Makanan+Kucing', badge: 'Promo' },
  { id: '3', name: 'Tempat Makan Anti-Tumpah', category: 'Aksesoris', price: 65000, image: 'https://placehold.co/300x300/e5f4da/3a6820?text=Aksesoris', badge: null },
  { id: '4', name: 'Kandang Kucing Lipat', category: 'Kucing', price: 320000, image: 'https://placehold.co/300x300/e5f4da/3a6820?text=Kandang', badge: 'Baru' },
  { id: '5', name: 'Pakan Ikan Mas Koki', category: 'Ikan', price: 25000, image: 'https://placehold.co/300x300/e5f4da/3a6820?text=Pakan+Ikan', badge: null },
  { id: '6', name: 'Vitamin Anjing Senior', category: 'Anjing', price: 98000, image: 'https://placehold.co/300x300/e5f4da/3a6820?text=Vitamin', badge: 'Promo' },
  { id: '7', name: 'Sangkar Burung Kayu', category: 'Burung', price: 210000, image: 'https://placehold.co/300x300/e5f4da/3a6820?text=Sangkar', badge: null },
  { id: '8', name: 'Pasir Kucing Wangi 5kg', category: 'Kucing', price: 72000, image: 'https://placehold.co/300x300/e5f4da/3a6820?text=Pasir+Kucing', badge: 'Terlaris' },
]

const CATEGORIES = [
  { label: 'Anjing', emoji: '🐕', to: '/shop?category=anjing' },
  { label: 'Kucing', emoji: '🐈', to: '/shop?category=kucing' },
]

const TRUST_BADGES = [
  { icon: Truck,       label: 'Gratis Ongkir', desc: 'Pembelian ≥ Rp 150.000' },
  { icon: ShieldCheck, label: 'Produk Asli',   desc: '100% original & bergaransi' },
  { icon: Star,        label: '4.9/5 Rating',  desc: 'Dari 10.000+ pelanggan' },
  { icon: Tag,         label: 'Harga Terbaik', desc: 'Selalu ada promo spesial' },
]

function formatHarga(amount: number) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount)
}

export default function Home() {
  return (
    <div className="pb-10">

      {/* ── HERO — flows visually from the green header ── */}
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
              <Link to="/shop?sort=new" className="border border-white/40 hover:bg-white/10 text-white font-semibold px-5 py-2 rounded-xl text-sm transition-colors">
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
          <div className="grid grid-cols-2 gap-3">
            {CATEGORIES.map(({ label, emoji, to }) => (
              <Link
                key={label}
                to={to}
                className="flex items-center gap-3 p-4 rounded-2xl bg-white border border-gray-100 hover:border-primary-300 hover:bg-primary-50 transition-colors group"
              >
                <span className="text-3xl group-hover:scale-110 transition-transform">{emoji}</span>
                <span className="text-sm font-semibold text-gray-700">{label}</span>
              </Link>
            ))}
          </div>
        </section>

        {/* ── PROMO BANNER ── */}
        <section className="mb-8">
          <div className="bg-gradient-to-r from-amber-400 to-orange-400 rounded-2xl p-5 flex items-center justify-between">
            <div>
              <p className="text-white/80 text-xs font-medium mb-0.5">Promo Spesial Hari Ini</p>
              <p className="text-white font-bold text-lg leading-tight">Diskon hingga 40%<br />untuk makanan kucing</p>
              <Link to="/shop?category=kucing&promo=true" className="inline-block mt-3 bg-white text-orange-600 hover:bg-orange-50 font-semibold px-4 py-1.5 rounded-xl text-sm transition-colors">
                Lihat Promo
              </Link>
            </div>
            <div className="text-6xl select-none">🐈</div>
          </div>
        </section>

        {/* ── PRODUK UNGGULAN ── */}
        <section>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-base font-bold text-gray-900">Produk Unggulan</h2>
            <Link to="/shop" className="text-sm text-primary-600 hover:underline flex items-center gap-0.5">
              Lihat semua <ChevronRight className="w-3.5 h-3.5" />
            </Link>
          </div>

          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
            {PLACEHOLDER_PRODUCTS.map((product, i) => (
              <motion.div
                key={product.id}
                initial={{ opacity: 0, y: 16 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.05 }}
              >
                <Link to={`/produk/${product.id}`} className="card block overflow-hidden hover:shadow-md transition-shadow group">
                  <div className="relative overflow-hidden bg-primary-50">
                    <img
                      src={product.image}
                      alt={product.name}
                      className="w-full aspect-square object-cover group-hover:scale-105 transition-transform duration-300"
                    />
                    {product.badge && (
                      <span className={`absolute top-2 left-2 text-xs font-semibold px-2 py-0.5 rounded-full ${
                        product.badge === 'Promo'   ? 'bg-orange-400 text-white' :
                        product.badge === 'Baru'    ? 'bg-primary-500 text-white' :
                        'bg-red-500 text-white'
                      }`}>
                        {product.badge}
                      </span>
                    )}
                  </div>
                  <div className="p-3">
                    <p className="text-xs text-primary-600 font-medium mb-0.5">{product.category}</p>
                    <p className="text-sm font-semibold text-gray-800 line-clamp-2 leading-snug mb-2">{product.name}</p>
                    <p className="text-base font-bold text-gray-900">{formatHarga(product.price)}</p>
                    <button className="mt-2 w-full btn-primary text-sm py-1.5">
                      + Keranjang
                    </button>
                  </div>
                </Link>
              </motion.div>
            ))}
          </div>
        </section>

      </div>
    </div>
  )
}
