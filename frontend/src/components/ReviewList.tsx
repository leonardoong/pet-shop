import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import axios from 'axios'
import StarRating from './StarRating'
import { useState } from 'react'
import { useAuthStore } from '@/store/authStore'
import type { ApiResponse, Paginated } from '@/types'

interface ReviewItem {
  id: string
  customer_name: string
  rating: number
  comment: string
  created_at: string
}

export default function ReviewList({ productSlug }: { productSlug: string }) {
  const queryClient = useQueryClient()
  const isAuth = useAuthStore(s => s.isAuthenticated)
  const [page, setPage] = useState(1)
  const [showForm, setShowForm] = useState(false)
  const [rating, setRating] = useState(5)
  const [comment, setComment] = useState('')
  const [productId, setProductId] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['reviews', productSlug, page],
    queryFn: async () => {
      const pRes = await axios.get(`/api/v1/products/${productSlug}`)
      const pid = pRes.data?.data?.id
      if (pid) setProductId(pid)
      const res = await axios.get<ApiResponse<Paginated<ReviewItem>>>(`/api/v1/products/${productSlug}/reviews`, { params: { page, limit: 10 } })
      return res.data.data
    },
    enabled: !!productSlug,
  })

  const submitMutation = useMutation({
    mutationFn: () => axios.post('/api/v1/customer/reviews', { rating, comment }, { params: { product_id: productId } }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reviews', productSlug] })
      setShowForm(false)
      setComment('')
    },
    onError: (err: any) => alert(err?.response?.data?.message || 'Gagal mengirim review'),
  })

  const avgRating = data?.items?.length
    ? (data.items.reduce((s, r) => s + r.rating, 0) / data.items.length).toFixed(1)
    : null

  return (
    <div className="mt-10">
      <div className="flex items-center gap-3 mb-4">
        <h3 className="text-xl font-bold text-gray-900">Ulasan</h3>
        {avgRating && (
          <div className="flex items-center gap-2 text-sm">
            <StarRating rating={Math.round(Number(avgRating))} />
            <span className="text-gray-500">({avgRating} / {data?.total})</span>
          </div>
        )}
      </div>

      {isAuth && productId && !showForm && (
        <button onClick={() => setShowForm(true)} className="text-sm text-primary-600 font-medium mb-4 hover:underline">
          Tulis Ulasan
        </button>
      )}

      {showForm && (
        <div className="card p-4 mb-4">
          <div className="flex gap-1 mb-3">
            {[1,2,3,4,5].map(i => (
              <button key={i} onClick={() => setRating(i)} className="text-2xl">
                {i <= rating ? '★' : '☆'}
              </button>
            ))}
          </div>
          <textarea value={comment} onChange={e => setComment(e.target.value)} rows={3}
            placeholder="Bagikan pengalaman Anda..." className="w-full border rounded-lg px-3 py-2 text-sm mb-3" />
          <div className="flex gap-2">
            <button onClick={() => submitMutation.mutate()} disabled={submitMutation.isPending}
              className="btn-primary text-sm">Kirim</button>
            <button onClick={() => setShowForm(false)} className="text-sm text-gray-500 hover:underline">Batal</button>
          </div>
        </div>
      )}

      {isLoading ? (
        <p className="text-gray-400 text-sm">Memuat ulasan...</p>
      ) : !data || data.items.length === 0 ? (
        <p className="text-gray-400 text-sm">Belum ada ulasan.</p>
      ) : (
        <div className="space-y-4">
          {data.items.map(r => (
            <div key={r.id} className="border-b pb-3">
              <div className="flex items-center gap-2 mb-1">
                <span className="font-medium text-sm text-gray-900">{r.customer_name}</span>
                <StarRating rating={r.rating} />
                <span className="text-xs text-gray-400">{r.created_at}</span>
              </div>
              {r.comment && <p className="text-sm text-gray-600">{r.comment}</p>}
            </div>
          ))}
          {data.total_pages > 1 && (
            <div className="flex gap-2 text-sm">
              <button disabled={page <= 1} onClick={() => setPage(p => p - 1)} className="text-primary-600 disabled:opacity-30">Sebelumnya</button>
              <button disabled={page >= data.total_pages} onClick={() => setPage(p => p + 1)} className="text-primary-600 disabled:opacity-30">Berikutnya</button>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
