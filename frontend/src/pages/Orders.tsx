import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { ChevronLeft, ChevronRight, Package } from 'lucide-react'
import { ordersApi } from '@/api/orders'
import { formatPrice, formatDate } from '@/lib/utils'
import type { OrderStatus } from '@/types'

const STATUS_CONFIG: Record<OrderStatus, { label: string; color: string }> = {
  pending:     { label: 'Menunggu',      color: 'bg-yellow-100 text-yellow-700' },
  confirmed:   { label: 'Dikonfirmasi',  color: 'bg-blue-100 text-blue-700' },
  processing:  { label: 'Diproses',      color: 'bg-indigo-100 text-indigo-700' },
  shipped:     { label: 'Dikirim',       color: 'bg-purple-100 text-purple-700' },
  delivered:   { label: 'Selesai',       color: 'bg-green-100 text-green-700' },
  cancelled:   { label: 'Dibatalkan',    color: 'bg-red-100 text-red-700' },
}

const STATUS_OPTIONS: { value: OrderStatus | ''; label: string }[] = [
  { value: '',           label: 'Semua Status' },
  { value: 'pending',    label: 'Menunggu' },
  { value: 'confirmed',  label: 'Dikonfirmasi' },
  { value: 'processing', label: 'Diproses' },
  { value: 'shipped',    label: 'Dikirim' },
  { value: 'delivered',  label: 'Selesai' },
  { value: 'cancelled',  label: 'Dibatalkan' },
]

export default function Orders() {
  const [status, setStatus] = useState<OrderStatus | ''>('')
  const [page, setPage] = useState(1)

  const { data, isLoading } = useQuery({
    queryKey: ['orders', status, page],
    queryFn: () =>
      ordersApi
        .list({ status: status || undefined, page, limit: 10 })
        .then((r) => r.data.data),
  })

  const orders     = data?.items ?? []
  const totalPages = data?.total_pages ?? 1

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 py-6">
      <h1 className="text-xl font-bold text-gray-900 mb-6 flex items-center gap-2">
        <Package className="w-5 h-5 text-primary-500" /> Pesanan Saya
      </h1>

      {/* Status filter */}
      <div className="flex gap-2 flex-wrap mb-5">
        {STATUS_OPTIONS.map((opt) => (
          <button
            key={opt.value}
            onClick={() => { setStatus(opt.value); setPage(1) }}
            className={`text-xs px-3 py-1.5 rounded-full border transition-colors ${
              status === opt.value
                ? 'bg-primary-500 text-white border-primary-500'
                : 'border-gray-200 text-gray-600 hover:border-primary-300'
            }`}
          >
            {opt.label}
          </button>
        ))}
      </div>

      {isLoading ? (
        <div className="space-y-3 animate-pulse">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="card p-4 space-y-2">
              <div className="flex justify-between">
                <div className="h-4 bg-gray-100 rounded w-24" />
                <div className="h-4 bg-gray-100 rounded w-16" />
              </div>
              <div className="h-3 bg-gray-100 rounded w-1/3" />
              <div className="h-3 bg-gray-100 rounded w-1/2" />
            </div>
          ))}
        </div>
      ) : orders.length === 0 ? (
        <div className="text-center py-16 text-gray-400">
          <p className="text-4xl mb-3">📦</p>
          <p className="font-medium text-gray-600">Belum ada pesanan</p>
          {status && (
            <button onClick={() => setStatus('')} className="mt-2 text-sm text-primary-600 hover:underline">
              Lihat semua pesanan
            </button>
          )}
          {!status && (
            <Link to="/shop" className="mt-3 inline-block text-sm text-primary-600 hover:underline">
              Mulai belanja
            </Link>
          )}
        </div>
      ) : (
        <div className="space-y-3">
          {orders.map((order) => {
            const cfg = STATUS_CONFIG[order.status]
            return (
              <Link
                key={order.id}
                to={`/pesanan/${order.id}`}
                className="card block p-4 hover:shadow-md transition-shadow"
              >
                <div className="flex items-start justify-between gap-2 mb-2">
                  <div>
                    <p className="text-xs text-gray-400 mb-0.5">
                      #{order.id.slice(0, 8).toUpperCase()}
                    </p>
                    <p className="text-xs text-gray-500">{formatDate(order.created_at)}</p>
                  </div>
                  <span className={`text-xs font-semibold px-2.5 py-1 rounded-full shrink-0 ${cfg.color}`}>
                    {cfg.label}
                  </span>
                </div>

                <div className="flex items-center justify-between">
                  <div className="text-sm text-gray-600">
                    {order.items?.length ?? 0} produk · Dikirim ke {order.ship_city}
                  </div>
                  <p className="text-sm font-bold text-gray-900">{formatPrice(order.total_amount)}</p>
                </div>

                {/* First 2 item names preview */}
                {order.items && order.items.length > 0 && (
                  <p className="text-xs text-gray-400 mt-1.5 truncate">
                    {order.items
                      .slice(0, 2)
                      .map((i) => i.product?.name ?? `Produk`)
                      .join(', ')}
                    {order.items.length > 2 && ` +${order.items.length - 2} lainnya`}
                  </p>
                )}
              </Link>
            )
          })}
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-2 mt-8">
          <button
            disabled={page <= 1}
            onClick={() => setPage((p) => p - 1)}
            className="p-2 rounded-lg border border-gray-200 disabled:opacity-40 hover:bg-gray-50 transition-colors"
          >
            <ChevronLeft className="w-4 h-4" />
          </button>
          <span className="text-sm text-gray-600">
            Halaman {page} dari {totalPages}
          </span>
          <button
            disabled={page >= totalPages}
            onClick={() => setPage((p) => p + 1)}
            className="p-2 rounded-lg border border-gray-200 disabled:opacity-40 hover:bg-gray-50 transition-colors"
          >
            <ChevronRight className="w-4 h-4" />
          </button>
        </div>
      )}
    </div>
  )
}
