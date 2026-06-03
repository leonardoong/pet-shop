import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Search, Eye, ChevronDown } from 'lucide-react'
import { ordersApi } from '@/api/orders'
import StatusBadge from '@/components/StatusBadge'
import type { AdminOrder } from '@/types'

const STATUSES = ['pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled']
const validTransitions: Record<string, string[]> = {
  pending:    ['confirmed', 'cancelled'],
  confirmed:  ['processing', 'cancelled'],
  processing: ['shipped', 'cancelled'],
  shipped:    ['delivered'],
}

export default function OrdersPage() {
  const queryClient = useQueryClient()
  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [page, setPage] = useState(1)
  const [detail, setDetail] = useState<AdminOrder | null>(null)
  const [openDropdown, setOpenDropdown] = useState<string | null>(null)

  const { data, isLoading } = useQuery({
    queryKey: ['admin-orders', { search, statusFilter, page }],
    queryFn: () => ordersApi.list({ search, ...(statusFilter && { status: statusFilter }), page: String(page), limit: '15' }),
  })

  const orders = data?.data?.data

  const statusMutation = useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) =>
      ordersApi.updateStatus(id, { status }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-orders'] })
      if (detail) queryClient.invalidateQueries({ queryKey: ['admin-orders', detail.id] })
    },
  })

  return (
    <div className="space-y-4">
      <div>
        <h2 className="text-2xl font-bold text-gray-900">Orders</h2>
        <p className="text-sm text-gray-500 mt-1">Manage customer orders</p>
      </div>

      <div className="flex gap-3 flex-wrap">
        <div className="relative flex-1 min-w-[200px]">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input value={search} onChange={e => { setSearch(e.target.value); setPage(1) }} placeholder="Search by customer or order ID..."
            className="w-full pl-9 pr-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:outline-none" />
        </div>
        <select value={statusFilter} onChange={e => { setStatusFilter(e.target.value); setPage(1) }}
          className="border border-gray-300 rounded-lg px-3 py-2 text-sm">
          <option value="">All Status</option>
          {STATUSES.map(s => <option key={s} value={s}>{s}</option>)}
        </select>
      </div>

      <div className="bg-white rounded-xl border border-gray-200 overflow-visible">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 text-gray-600 text-left rounded-t-xl overflow-hidden">
            <tr>
              <th className="px-4 py-3 font-medium">Order ID</th>
              <th className="px-4 py-3 font-medium">Customer</th>
              <th className="px-4 py-3 font-medium">Date</th>
              <th className="px-4 py-3 font-medium text-right">Total</th>
              <th className="px-4 py-3 font-medium">Status</th>
              <th className="px-4 py-3 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {isLoading ? (
              <tr><td colSpan={6} className="px-4 py-8 text-center text-gray-400">Loading...</td></tr>
            ) : !orders || orders.items.length === 0 ? (
              <tr><td colSpan={6} className="px-4 py-8 text-center text-gray-400">No orders found.</td></tr>
            ) : orders.items.map(o => (
              <tr key={o.id} className="hover:bg-gray-50">
                <td className="px-4 py-3 font-mono text-xs text-gray-600">{o.id.slice(0, 8)}...</td>
                <td className="px-4 py-3 text-gray-900">{o.ship_name}</td>
                <td className="px-4 py-3 text-gray-500">{new Date(o.created_at).toLocaleDateString()}</td>
                <td className="px-4 py-3 text-right text-gray-900">Rp {o.total_amount.toLocaleString()}</td>
                <td className="px-4 py-3"><StatusBadge status={o.status} /></td>
                <td className="px-4 py-3 text-right">
                  <div className="flex items-center gap-1 justify-end">
                    <button onClick={() => setDetail(o)}
                      className="inline-flex items-center gap-1 px-2.5 py-1.5 text-xs font-medium text-blue-600 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors"
                      title="View detail">
                      <Eye className="w-3.5 h-3.5" /> View
                    </button>
                    {o.status !== 'delivered' && o.status !== 'cancelled' && (
                      <div className="relative">
                        <button
                          onClick={() => setOpenDropdown(openDropdown === o.id ? null : o.id)}
                          className="inline-flex items-center gap-1 px-2.5 py-1.5 text-xs font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors">
                          <ChevronDown className="w-3.5 h-3.5" /> Update
                        </button>
                        {openDropdown === o.id && (
                          <>
                            <div className="fixed inset-0 z-10" onClick={() => setOpenDropdown(null)} />
                            <div className="absolute right-0 mt-1 w-36 bg-white border rounded-lg shadow-lg z-20 py-1">
                              {(validTransitions[o.status] || []).map(s => (
                                <button key={s} onClick={() => { statusMutation.mutate({ id: o.id, status: s }); setOpenDropdown(null) }}
                                  className="block w-full text-left px-3 py-1.5 text-sm hover:bg-gray-100 capitalize">
                                  Mark {s}
                                </button>
                              ))}
                            </div>
                          </>
                        )}
                      </div>
                    )}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {orders && orders.total_pages > 1 && (
          <div className="flex items-center justify-between px-4 py-3 border-t bg-gray-50">
            <span className="text-sm text-gray-500">Page {orders.page} of {orders.total_pages}</span>
            <div className="flex gap-2">
              <button disabled={page <= 1} onClick={() => setPage(p => p - 1)} className="px-3 py-1 text-sm border rounded-md disabled:opacity-40">Prev</button>
              <button disabled={page >= orders.total_pages} onClick={() => setPage(p => p + 1)} className="px-3 py-1 text-sm border rounded-md disabled:opacity-40">Next</button>
            </div>
          </div>
        )}
      </div>

      {detail && (
        <div className="fixed inset-0 z-50 flex justify-end">
          <div className="absolute inset-0 bg-black/40" onClick={() => setDetail(null)} />
          <div className="relative w-full max-w-lg bg-white shadow-xl overflow-y-auto p-6">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-lg font-semibold">Order Detail</h3>
              <button onClick={() => setDetail(null)} className="text-gray-400 hover:text-gray-600 text-xl">&times;</button>
            </div>
            <div className="space-y-4">
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Order ID</span>
                <span className="font-mono">{detail.id.slice(0, 8)}...</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Status</span>
                <StatusBadge status={detail.status} />
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Date</span>
                <span>{new Date(detail.created_at).toLocaleString()}</span>
              </div>
              <hr />
              <h4 className="text-sm font-semibold text-gray-900">Shipping Info</h4>
              <p className="text-sm">{detail.ship_name}</p>
              <p className="text-sm text-gray-500">{detail.ship_phone}</p>
              <p className="text-sm text-gray-500">{detail.ship_street}, {detail.ship_city}, {detail.ship_province} {detail.ship_postal}</p>
              {detail.notes && <p className="text-sm text-gray-400 italic">Notes: {detail.notes}</p>}
              <hr />
              <h4 className="text-sm font-semibold text-gray-900">Items</h4>
              <table className="w-full text-sm">
                <thead><tr className="text-gray-500 text-left"><th className="py-1">Item</th><th className="text-right py-1">Qty</th><th className="text-right py-1">Total</th></tr></thead>
                <tbody>
                  {detail.items.map(i => (
                    <tr key={i.id}>
                      <td className="py-1">{i.product_name}</td>
                      <td className="text-right py-1">{i.quantity}</td>
                      <td className="text-right py-1">Rp {i.subtotal.toLocaleString()}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
              <div className="flex justify-between font-semibold text-sm pt-2 border-t">
                <span>Total</span>
                <span>Rp {detail.total_amount.toLocaleString()}</span>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
