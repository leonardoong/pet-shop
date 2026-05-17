import { useParams, Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { ChevronLeft, MapPin, Package } from 'lucide-react'
import { ordersApi } from '@/api/orders'
import { formatPrice, formatDate } from '@/lib/utils'
import type { OrderStatus } from '@/types'

const STATUS_CONFIG: Record<OrderStatus, { label: string; color: string }> = {
  pending:     { label: 'Menunggu Konfirmasi', color: 'bg-yellow-100 text-yellow-700' },
  confirmed:   { label: 'Dikonfirmasi',        color: 'bg-blue-100 text-blue-700' },
  processing:  { label: 'Sedang Diproses',     color: 'bg-indigo-100 text-indigo-700' },
  shipped:     { label: 'Sedang Dikirim',      color: 'bg-purple-100 text-purple-700' },
  delivered:   { label: 'Selesai',             color: 'bg-green-100 text-green-700' },
  cancelled:   { label: 'Dibatalkan',          color: 'bg-red-100 text-red-700' },
}

export default function OrderDetail() {
  const { id } = useParams<{ id: string }>()

  const { data, isLoading, isError } = useQuery({
    queryKey: ['order', id],
    queryFn: () => ordersApi.getById(id!).then((r) => r.data.data),
    enabled: !!id,
  })

  if (isLoading) {
    return (
      <div className="max-w-3xl mx-auto px-4 py-8 animate-pulse space-y-4">
        <div className="h-6 bg-gray-100 rounded w-48" />
        <div className="card p-5 space-y-3">
          <div className="h-4 bg-gray-100 rounded w-32" />
          <div className="h-3 bg-gray-100 rounded w-full" />
          <div className="h-3 bg-gray-100 rounded w-3/4" />
        </div>
      </div>
    )
  }

  if (isError || !data) {
    return (
      <div className="max-w-3xl mx-auto px-4 py-16 text-center">
        <p className="text-5xl mb-4">😕</p>
        <p className="text-lg font-semibold text-gray-700">Pesanan tidak ditemukan</p>
        <Link to="/pesanan" className="mt-4 inline-block text-primary-600 hover:underline text-sm">
          Kembali ke daftar pesanan
        </Link>
      </div>
    )
  }

  const order = data
  const cfg   = STATUS_CONFIG[order.status]
  const items = order.items ?? []

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 py-6">
      <Link to="/pesanan" className="flex items-center gap-1 text-sm text-gray-500 hover:text-primary-600 mb-6 transition-colors">
        <ChevronLeft className="w-4 h-4" /> Semua Pesanan
      </Link>

      {/* Header */}
      <div className="flex items-start justify-between flex-wrap gap-3 mb-5">
        <div>
          <h1 className="text-lg font-bold text-gray-900">
            Pesanan #{order.id.slice(0, 8).toUpperCase()}
          </h1>
          <p className="text-sm text-gray-500 mt-0.5">{formatDate(order.created_at)}</p>
        </div>
        <span className={`text-sm font-semibold px-3 py-1.5 rounded-full ${cfg.color}`}>
          {cfg.label}
        </span>
      </div>

      <div className="space-y-4">
        {/* Shipping address */}
        <div className="card p-5">
          <p className="font-semibold text-gray-900 mb-3 flex items-center gap-2">
            <MapPin className="w-4 h-4 text-primary-500" /> Alamat Pengiriman
          </p>
          <p className="text-sm font-medium text-gray-800">{order.ship_name}</p>
          <p className="text-sm text-gray-600 mt-0.5">{order.ship_phone}</p>
          <p className="text-sm text-gray-600 mt-1 leading-relaxed">
            {order.ship_street}, {order.ship_city}, {order.ship_province} {order.ship_postal}
          </p>
        </div>

        {/* Items */}
        <div className="card p-5">
          <p className="font-semibold text-gray-900 mb-4 flex items-center gap-2">
            <Package className="w-4 h-4 text-primary-500" /> Produk ({items.length})
          </p>
          <div className="space-y-4">
            {items.map((item) => {
              const imageUrl =
                item.product?.image_url ||
                `https://placehold.co/64x64/e5f4da/3a6820?text=Produk`

              return (
                <div key={item.id} className="flex gap-3">
                  <img
                    src={imageUrl}
                    alt={item.product?.name}
                    className="w-16 h-16 object-cover rounded-xl bg-primary-50 shrink-0"
                  />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-semibold text-gray-800 leading-snug">
                      {item.product?.name ?? 'Produk'}
                    </p>
                    <p className="text-xs text-gray-500 mt-0.5">
                      {formatPrice(item.unit_price)} × {item.quantity}
                    </p>
                  </div>
                  <p className="text-sm font-bold text-gray-900 shrink-0">
                    {formatPrice(item.subtotal)}
                  </p>
                </div>
              )
            })}
          </div>
        </div>

        {/* Notes */}
        {order.notes && (
          <div className="card p-5">
            <p className="font-semibold text-gray-900 mb-1 text-sm">Catatan</p>
            <p className="text-sm text-gray-600">{order.notes}</p>
          </div>
        )}

        {/* Total */}
        <div className="card p-5">
          <div className="flex justify-between items-center">
            <p className="font-bold text-gray-900">Total Pembayaran</p>
            <p className="text-xl font-bold text-gray-900">{formatPrice(order.total_amount)}</p>
          </div>
        </div>
      </div>
    </div>
  )
}
