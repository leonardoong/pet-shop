import { useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Minus, Plus, Trash2, ShoppingCart, ArrowRight } from 'lucide-react'
import { useCartStore } from '@/store/cartStore'
import { formatPrice } from '@/lib/utils'

export default function Cart() {
  const navigate = useNavigate()
  const { cart, isLoading, fetch, updateItem, removeItem } = useCartStore()

  useEffect(() => { fetch() }, [fetch])

  const items = cart?.items ?? []
  const total = items.reduce((sum, item) => sum + item.product.price * item.quantity, 0)

  if (isLoading) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-8 animate-pulse space-y-4">
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i} className="card flex gap-4 p-4">
            <div className="w-20 h-20 bg-gray-100 rounded-xl shrink-0" />
            <div className="flex-1 space-y-2">
              <div className="h-4 bg-gray-100 rounded w-1/2" />
              <div className="h-3 bg-gray-100 rounded w-1/4" />
              <div className="h-3 bg-gray-100 rounded w-16" />
            </div>
          </div>
        ))}
      </div>
    )
  }

  if (items.length === 0) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-20 text-center">
        <div className="text-6xl mb-4">🛒</div>
        <h2 className="text-xl font-bold text-gray-800 mb-2">Keranjangmu masih kosong</h2>
        <p className="text-gray-500 text-sm mb-6">Yuk mulai belanja kebutuhan hewan peliharaanmu!</p>
        <Link to="/shop" className="btn-primary px-8 py-3 rounded-2xl inline-flex items-center gap-2">
          <ShoppingCart className="w-4 h-4" /> Mulai Belanja
        </Link>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 py-6">
      <h1 className="text-xl font-bold text-gray-900 mb-6 flex items-center gap-2">
        <ShoppingCart className="w-5 h-5 text-primary-500" />
        Keranjang Belanja
        <span className="text-sm font-normal text-gray-500">({items.length} produk)</span>
      </h1>

      <div className="flex flex-col lg:flex-row gap-6">
        {/* Item list */}
        <div className="flex-1 space-y-3">
          {items.map((item) => {
            const imageUrl =
              item.product.image_url ||
              `https://placehold.co/80x80/e5f4da/3a6820?text=${encodeURIComponent(item.product.name)}`

            return (
              <div key={item.id} className="card flex gap-4 p-4">
                <Link to={`/produk/${item.product.slug}`} className="shrink-0">
                  <img
                    src={imageUrl}
                    alt={item.product.name}
                    className="w-20 h-20 object-cover rounded-xl bg-primary-50"
                  />
                </Link>

                <div className="flex-1 min-w-0">
                  <Link
                    to={`/produk/${item.product.slug}`}
                    className="text-sm font-semibold text-gray-800 hover:text-primary-600 transition-colors line-clamp-2"
                  >
                    {item.product.name}
                  </Link>
                  <p className="text-xs text-primary-600 mt-0.5">{item.product.category?.name}</p>
                  <p className="text-sm font-bold text-gray-900 mt-1">{formatPrice(item.product.price)}</p>

                  <div className="flex items-center justify-between mt-2">
                    {/* Quantity controls */}
                    <div className="flex items-center border border-gray-200 rounded-lg overflow-hidden">
                      <button
                        onClick={() => {
                          if (item.quantity <= 1) removeItem(item.product_id)
                          else updateItem(item.product_id, item.quantity - 1)
                        }}
                        className="w-8 h-8 flex items-center justify-center text-gray-600 hover:bg-gray-50 transition-colors"
                      >
                        <Minus className="w-3 h-3" />
                      </button>
                      <span className="w-8 text-center text-sm font-semibold">{item.quantity}</span>
                      <button
                        onClick={() => updateItem(item.product_id, item.quantity + 1)}
                        disabled={item.quantity >= item.product.stock}
                        className="w-8 h-8 flex items-center justify-center text-gray-600 hover:bg-gray-50 transition-colors disabled:opacity-40"
                      >
                        <Plus className="w-3 h-3" />
                      </button>
                    </div>

                    <div className="flex items-center gap-3">
                      <p className="text-sm font-bold text-gray-900">
                        {formatPrice(item.product.price * item.quantity)}
                      </p>
                      <button
                        onClick={() => removeItem(item.product_id)}
                        className="p-1.5 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                        aria-label="Hapus item"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            )
          })}
        </div>

        {/* Order summary */}
        <div className="lg:w-72 shrink-0">
          <div className="card p-5 sticky top-24">
            <p className="font-bold text-gray-900 mb-4">Ringkasan Pesanan</p>

            <div className="space-y-2 text-sm">
              {items.map((item) => (
                <div key={item.id} className="flex justify-between text-gray-600">
                  <span className="truncate mr-2 flex-1">{item.product.name} ×{item.quantity}</span>
                  <span className="shrink-0">{formatPrice(item.product.price * item.quantity)}</span>
                </div>
              ))}
            </div>

            <div className="border-t border-gray-100 mt-4 pt-4 flex justify-between font-bold text-gray-900">
              <span>Total</span>
              <span>{formatPrice(total)}</span>
            </div>

            <button
              onClick={() => navigate('/checkout')}
              className="mt-4 w-full btn-primary py-3 rounded-2xl font-semibold flex items-center justify-center gap-2"
            >
              Lanjut Checkout <ArrowRight className="w-4 h-4" />
            </button>

            <Link to="/shop" className="block mt-3 text-center text-sm text-primary-600 hover:underline">
              Lanjut belanja
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}
