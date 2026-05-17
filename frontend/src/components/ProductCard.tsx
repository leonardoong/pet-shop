import { useState } from 'react'
import { Link } from 'react-router-dom'
import { ShoppingCart, Check } from 'lucide-react'
import { motion } from 'framer-motion'
import { formatPrice } from '@/lib/utils'
import { useCartStore } from '@/store/cartStore'
import type { Product } from '@/types'

interface Props {
  product: Product
  index?: number
}

export default function ProductCard({ product, index = 0 }: Props) {
  const [added, setAdded] = useState(false)
  const addItem = useCartStore((s) => s.addItem)

  const handleAddToCart = async (e: React.MouseEvent) => {
    e.preventDefault()
    if (added || product.stock === 0) return
    try {
      await addItem(product.id, 1)
      setAdded(true)
      setTimeout(() => setAdded(false), 1500)
    } catch {
      // silent — user still sees the button revert
    }
  }

  const imageUrl =
    product.image_url ||
    `https://placehold.co/300x300/e5f4da/3a6820?text=${encodeURIComponent(product.name)}`

  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{ delay: Math.min(index * 0.05, 0.3) }}
      className="h-full"
    >
      <div className="card overflow-hidden hover:shadow-md transition-shadow group h-full flex flex-col">
        <Link to={`/produk/${product.slug}`} className="block relative overflow-hidden bg-primary-50">
          <img
            src={imageUrl}
            alt={product.name}
            className="w-full aspect-square object-cover group-hover:scale-105 transition-transform duration-300"
          />
          {product.stock === 0 && (
            <div className="absolute inset-0 bg-black/40 flex items-center justify-center">
              <span className="text-white text-xs font-semibold bg-black/60 px-3 py-1 rounded-full">
                Stok Habis
              </span>
            </div>
          )}
          {product.stock > 0 && product.stock <= 5 && (
            <span className="absolute top-2 left-2 text-xs font-semibold px-2 py-0.5 rounded-full bg-orange-400 text-white">
              Sisa {product.stock}
            </span>
          )}
        </Link>

        <div className="p-3 flex flex-col flex-1">
          <p className="text-xs text-primary-600 font-medium mb-0.5">{product.category?.name}</p>
          <Link to={`/produk/${product.slug}`} className="flex-1 hover:text-primary-700 transition-colors">
            <p className="text-sm font-semibold text-gray-800 line-clamp-2 leading-snug mb-2">
              {product.name}
            </p>
          </Link>
          <p className="text-base font-bold text-gray-900 mb-2">{formatPrice(product.price)}</p>
          <button
            onClick={handleAddToCart}
            disabled={product.stock === 0}
            className={`w-full text-sm py-1.5 rounded-xl font-semibold flex items-center justify-center gap-1.5 transition-all
              ${added
                ? 'bg-green-500 text-white'
                : product.stock === 0
                  ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                  : 'btn-primary'
              }`}
          >
            {added ? (
              <><Check className="w-3.5 h-3.5" /> Ditambahkan!</>
            ) : (
              <><ShoppingCart className="w-3.5 h-3.5" /> + Keranjang</>
            )}
          </button>
        </div>
      </div>
    </motion.div>
  )
}
