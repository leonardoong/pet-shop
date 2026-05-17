import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { ChevronLeft, Minus, Plus, ShoppingCart, Check, Package } from 'lucide-react'
import { productsApi } from '@/api/products'
import { useCartStore } from '@/store/cartStore'
import { useAuthStore } from '@/store/authStore'
import { formatPrice } from '@/lib/utils'

export default function ProductDetail() {
  const { slug } = useParams<{ slug: string }>()
  const navigate = useNavigate()
  const [quantity, setQuantity] = useState(1)
  const [added, setAdded] = useState(false)
  const addItem = useCartStore((s) => s.addItem)
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)

  const { data, isLoading, isError } = useQuery({
    queryKey: ['product', slug],
    queryFn: () => productsApi.getBySlug(slug!).then((r) => r.data.data),
    enabled: !!slug,
  })

  const product = data

  const handleAddToCart = async () => {
    if (!isAuthenticated) {
      navigate('/masuk')
      return
    }
    if (!product || product.stock === 0) return
    try {
      await addItem(product.id, quantity)
      setAdded(true)
      setTimeout(() => setAdded(false), 2000)
    } catch {
      // silent
    }
  }

  if (isLoading) {
    return (
      <div className="max-w-5xl mx-auto px-4 py-8">
        <div className="animate-pulse grid grid-cols-1 md:grid-cols-2 gap-8">
          <div className="aspect-square bg-gray-100 rounded-2xl" />
          <div className="space-y-4 py-4">
            <div className="h-3 bg-gray-100 rounded w-24" />
            <div className="h-7 bg-gray-100 rounded w-3/4" />
            <div className="h-8 bg-gray-100 rounded w-32" />
            <div className="h-20 bg-gray-100 rounded" />
          </div>
        </div>
      </div>
    )
  }

  if (isError || !product) {
    return (
      <div className="max-w-5xl mx-auto px-4 py-16 text-center">
        <p className="text-5xl mb-4">😕</p>
        <p className="text-lg font-semibold text-gray-700">Produk tidak ditemukan</p>
        <Link to="/shop" className="mt-4 inline-block text-primary-600 hover:underline text-sm">
          Kembali ke katalog
        </Link>
      </div>
    )
  }

  const imageUrl =
    product.image_url ||
    `https://placehold.co/600x600/e5f4da/3a6820?text=${encodeURIComponent(product.name)}`

  return (
    <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      {/* Breadcrumb */}
      <nav className="flex items-center gap-1.5 text-sm text-gray-500 mb-6">
        <Link to="/" className="hover:text-primary-600 transition-colors">Beranda</Link>
        <span>/</span>
        <Link to="/shop" className="hover:text-primary-600 transition-colors">Produk</Link>
        {product.category && (
          <>
            <span>/</span>
            <Link
              to={`/shop?category=${product.category.slug}`}
              className="hover:text-primary-600 transition-colors"
            >
              {product.category.name}
            </Link>
          </>
        )}
        <span>/</span>
        <span className="text-gray-800 font-medium truncate max-w-[160px]">{product.name}</span>
      </nav>

      <button
        onClick={() => navigate(-1)}
        className="flex items-center gap-1 text-sm text-gray-500 hover:text-primary-600 mb-6 transition-colors"
      >
        <ChevronLeft className="w-4 h-4" /> Kembali
      </button>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        {/* Image */}
        <div className="relative rounded-2xl overflow-hidden bg-primary-50">
          <img
            src={imageUrl}
            alt={product.name}
            className="w-full aspect-square object-cover"
          />
          {product.stock === 0 && (
            <div className="absolute inset-0 bg-black/40 flex items-center justify-center">
              <span className="text-white font-semibold bg-black/60 px-4 py-2 rounded-full">
                Stok Habis
              </span>
            </div>
          )}
        </div>

        {/* Info */}
        <div className="flex flex-col gap-4">
          {product.category && (
            <Link
              to={`/shop?category=${product.category.slug}`}
              className="text-sm text-primary-600 font-medium hover:underline w-fit"
            >
              {product.category.name}
            </Link>
          )}

          <h1 className="text-2xl font-bold text-gray-900 leading-tight">{product.name}</h1>

          <p className="text-3xl font-bold text-gray-900">{formatPrice(product.price)}</p>

          {/* Stock indicator */}
          <div className="flex items-center gap-2">
            <Package className="w-4 h-4 text-gray-400" />
            {product.stock === 0 ? (
              <span className="text-sm text-red-500 font-medium">Stok habis</span>
            ) : product.stock <= 10 ? (
              <span className="text-sm text-orange-500 font-medium">Sisa {product.stock} unit</span>
            ) : (
              <span className="text-sm text-green-600 font-medium">Tersedia ({product.stock} unit)</span>
            )}
          </div>

          {/* Description */}
          {product.description && (
            <div className="text-sm text-gray-600 leading-relaxed border-t border-gray-100 pt-4">
              <p className="font-semibold text-gray-800 mb-1">Deskripsi</p>
              <p className="whitespace-pre-line">{product.description}</p>
            </div>
          )}

          {/* SKU */}
          {product.sku && (
            <p className="text-xs text-gray-400">SKU: {product.sku}</p>
          )}

          {/* Quantity selector */}
          {product.stock > 0 && (
            <div className="flex items-center gap-3 mt-2">
              <p className="text-sm text-gray-600 font-medium">Jumlah:</p>
              <div className="flex items-center border border-gray-200 rounded-xl overflow-hidden">
                <button
                  onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                  className="w-9 h-9 flex items-center justify-center text-gray-600 hover:bg-gray-50 transition-colors"
                >
                  <Minus className="w-3.5 h-3.5" />
                </button>
                <span className="w-10 text-center text-sm font-semibold">{quantity}</span>
                <button
                  onClick={() => setQuantity((q) => Math.min(product.stock, q + 1))}
                  className="w-9 h-9 flex items-center justify-center text-gray-600 hover:bg-gray-50 transition-colors"
                >
                  <Plus className="w-3.5 h-3.5" />
                </button>
              </div>
            </div>
          )}

          {/* Add to cart */}
          <button
            onClick={handleAddToCart}
            disabled={product.stock === 0}
            className={`mt-2 py-3 px-6 rounded-2xl font-semibold flex items-center justify-center gap-2 transition-all text-base
              ${added
                ? 'bg-green-500 text-white'
                : product.stock === 0
                  ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                  : 'btn-primary'
              }`}
          >
            {added ? (
              <><Check className="w-5 h-5" /> Ditambahkan ke Keranjang!</>
            ) : (
              <><ShoppingCart className="w-5 h-5" /> Tambah ke Keranjang</>
            )}
          </button>

          {!isAuthenticated && product.stock > 0 && (
            <p className="text-xs text-gray-500 text-center">
              <Link to="/masuk" className="text-primary-600 hover:underline">Masuk</Link>{' '}
              untuk menambahkan ke keranjang
            </p>
          )}
        </div>
      </div>
    </div>
  )
}
