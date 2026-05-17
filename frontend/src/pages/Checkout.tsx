import { useState, useEffect } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useQuery, useMutation } from '@tanstack/react-query'
import { MapPin, Plus, CheckCircle2, ChevronRight } from 'lucide-react'
import { addressesApi } from '@/api/addresses'
import { ordersApi } from '@/api/orders'
import { useCartStore } from '@/store/cartStore'
import { formatPrice } from '@/lib/utils'

export default function Checkout() {
  const navigate = useNavigate()
  const { cart, clear } = useCartStore()
  const [selectedAddressId, setSelectedAddressId] = useState<string>('')
  const [notes, setNotes] = useState('')
  const [error, setError] = useState('')

  const { data: addresses, isLoading: loadingAddresses } = useQuery({
    queryKey: ['addresses'],
    queryFn: () => addressesApi.list().then((r) => r.data.data),
  })

  // Auto-select default address when addresses load
  useEffect(() => {
    if (addresses && !selectedAddressId) {
      const def = addresses.find((a) => a.is_default)
      if (def) setSelectedAddressId(def.id)
    }
  }, [addresses, selectedAddressId])

  const items = cart?.items ?? []
  const total = items.reduce((sum, item) => sum + item.product.price * item.quantity, 0)

  const { mutate: placeOrder, isPending: placing } = useMutation({
    mutationFn: () => ordersApi.checkout(selectedAddressId, notes || undefined).then((r) => r.data.data),
    onSuccess: (order) => {
      clear()
      navigate(`/pesanan/${order.id}`, { replace: true })
    },
    onError: (err: { response?: { data?: { message?: string } } }) => {
      setError(err.response?.data?.message ?? 'Terjadi kesalahan. Silakan coba lagi.')
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    if (!selectedAddressId) {
      setError('Pilih alamat pengiriman terlebih dahulu.')
      return
    }
    if (items.length === 0) {
      setError('Keranjangmu kosong.')
      return
    }
    placeOrder()
  }

  if (items.length === 0) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-20 text-center">
        <p className="text-5xl mb-3">🛒</p>
        <p className="text-lg font-semibold text-gray-700">Keranjangmu kosong</p>
        <Link to="/shop" className="mt-4 inline-block text-primary-600 hover:underline text-sm">
          Kembali belanja
        </Link>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 py-6">
      <h1 className="text-xl font-bold text-gray-900 mb-6">Checkout</h1>

      <form onSubmit={handleSubmit} className="flex flex-col lg:flex-row gap-6">
        {/* Left: address + notes */}
        <div className="flex-1 space-y-4">

          {/* Delivery address */}
          <div className="card p-5">
            <div className="flex items-center justify-between mb-4">
              <p className="font-semibold text-gray-900 flex items-center gap-2">
                <MapPin className="w-4 h-4 text-primary-500" /> Alamat Pengiriman
              </p>
              <Link
                to="/akun/alamat"
                className="text-xs text-primary-600 hover:underline flex items-center gap-0.5"
              >
                Kelola alamat <ChevronRight className="w-3 h-3" />
              </Link>
            </div>

            {loadingAddresses ? (
              <div className="space-y-3 animate-pulse">
                {[1, 2].map((i) => (
                  <div key={i} className="h-16 bg-gray-100 rounded-xl" />
                ))}
              </div>
            ) : !addresses || addresses.length === 0 ? (
              <div className="text-center py-6 text-gray-500">
                <p className="text-sm mb-3">Belum ada alamat tersimpan.</p>
                <Link
                  to="/akun/alamat"
                  className="inline-flex items-center gap-1 text-sm text-primary-600 hover:underline"
                >
                  <Plus className="w-3.5 h-3.5" /> Tambah alamat
                </Link>
              </div>
            ) : (
              <div className="space-y-2">
                {addresses.map((addr) => (
                  <label
                    key={addr.id}
                    className={`flex items-start gap-3 p-3 rounded-xl border cursor-pointer transition-colors ${
                      selectedAddressId === addr.id
                        ? 'border-primary-400 bg-primary-50'
                        : 'border-gray-200 hover:border-primary-200'
                    }`}
                  >
                    <input
                      type="radio"
                      name="address"
                      value={addr.id}
                      checked={selectedAddressId === addr.id}
                      onChange={() => setSelectedAddressId(addr.id)}
                      className="mt-0.5 accent-primary-500"
                    />
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 flex-wrap">
                        <p className="text-sm font-semibold text-gray-800">{addr.label}</p>
                        {addr.is_default && (
                          <span className="text-xs bg-primary-100 text-primary-700 px-1.5 py-0.5 rounded-full font-medium">
                            Utama
                          </span>
                        )}
                      </div>
                      <p className="text-xs text-gray-700 mt-0.5">{addr.recipient_name} · {addr.phone}</p>
                      <p className="text-xs text-gray-500 mt-0.5 leading-relaxed">
                        {addr.street}, {addr.city}, {addr.province} {addr.postal_code}
                      </p>
                    </div>
                    {selectedAddressId === addr.id && (
                      <CheckCircle2 className="w-4 h-4 text-primary-500 shrink-0 mt-0.5" />
                    )}
                  </label>
                ))}
              </div>
            )}
          </div>

          {/* Notes */}
          <div className="card p-5">
            <p className="font-semibold text-gray-900 mb-3">Catatan (opsional)</p>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="Contoh: Titipkan ke satpam jika tidak ada di rumah"
              rows={3}
              className="w-full border border-gray-200 rounded-xl px-3 py-2 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-primary-400"
            />
          </div>
        </div>

        {/* Right: order summary */}
        <div className="lg:w-72 shrink-0">
          <div className="card p-5 sticky top-24">
            <p className="font-bold text-gray-900 mb-4">Ringkasan Pesanan</p>

            <div className="space-y-2 text-sm max-h-48 overflow-y-auto">
              {items.map((item) => (
                <div key={item.id} className="flex justify-between text-gray-600">
                  <span className="truncate mr-2 flex-1">
                    {item.product.name} ×{item.quantity}
                  </span>
                  <span className="shrink-0">{formatPrice(item.product.price * item.quantity)}</span>
                </div>
              ))}
            </div>

            <div className="border-t border-gray-100 mt-4 pt-4 flex justify-between font-bold text-gray-900">
              <span>Total</span>
              <span>{formatPrice(total)}</span>
            </div>

            {error && (
              <p className="mt-3 text-xs text-red-500 bg-red-50 rounded-lg px-3 py-2">{error}</p>
            )}

            <button
              type="submit"
              disabled={placing || !selectedAddressId}
              className="mt-4 w-full btn-primary py-3 rounded-2xl font-semibold disabled:opacity-60 disabled:cursor-not-allowed"
            >
              {placing ? 'Memproses...' : 'Buat Pesanan'}
            </button>

            <Link to="/keranjang" className="block mt-3 text-center text-sm text-gray-500 hover:text-primary-600 transition-colors">
              ← Kembali ke keranjang
            </Link>
          </div>
        </div>
      </form>
    </div>
  )
}
