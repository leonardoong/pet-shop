import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Search, RefreshCcw } from 'lucide-react'
import { inventoryApi } from '@/api/inventory'
import { categoriesApi } from '@/api/categories'
import type { InventoryItem, AdjustStock } from '@/types'

export default function InventoryPage() {
  const queryClient = useQueryClient()
  const [search, setSearch] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('')
  const [lowStockOnly, setLowStockOnly] = useState(false)
  const [page, setPage] = useState(1)
  const [adjusting, setAdjusting] = useState<InventoryItem | null>(null)
  const [operation, setOperation] = useState<'add' | 'subtract' | 'set'>('add')
  const [qty, setQty] = useState(1)
  const [costPrice, setCostPrice] = useState('')
  const [note, setNote] = useState('')
  const [error, setError] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['admin-inventory', { search, categoryFilter, lowStockOnly, page }],
    queryFn: () => inventoryApi.list({
      search,
      ...(categoryFilter && { category_id: categoryFilter }),
      ...(lowStockOnly && { low_stock: 'true' }),
      page: String(page),
      limit: '20',
    }),
  })

  const { data: categoriesData } = useQuery({
    queryKey: ['admin-categories'],
    queryFn: () => categoriesApi.list(),
  })

  const inv = data?.data?.data
  const categories = categoriesData?.data?.data || []

  const adjustMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: AdjustStock }) => inventoryApi.adjustStock(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-inventory'] })
      setAdjusting(null); setQty(1); setCostPrice(''); setNote(''); setError('')
    },
    onError: (err: any) => setError(err?.response?.data?.message || 'Failed to adjust stock'),
  })

  return (
    <div className="space-y-4">
      <div>
        <h2 className="text-2xl font-bold text-gray-900">Inventory</h2>
        <p className="text-sm text-gray-500 mt-1">Monitor and manage product stock</p>
      </div>

      <div className="flex gap-3 flex-wrap items-center">
        <div className="relative flex-1 min-w-[200px]">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input value={search} onChange={e => { setSearch(e.target.value); setPage(1) }} placeholder="Search by name or SKU..."
            className="w-full pl-9 pr-3 py-2 border border-gray-300 rounded-lg text-sm" />
        </div>
        <select value={categoryFilter} onChange={e => { setCategoryFilter(e.target.value); setPage(1) }}
          className="border border-gray-300 rounded-lg px-3 py-2 text-sm">
          <option value="">All Categories</option>
          {categories.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
        </select>
        <label className="flex items-center gap-2 text-sm text-gray-600">
          <input type="checkbox" checked={lowStockOnly} onChange={e => { setLowStockOnly(e.target.checked); setPage(1) }} className="rounded" />
          Low stock only
        </label>
      </div>

      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 text-gray-600 text-left">
            <tr>
              <th className="px-4 py-3 font-medium">Product</th>
              <th className="px-4 py-3 font-medium">SKU</th>
              <th className="px-4 py-3 font-medium">Category</th>
              <th className="px-4 py-3 font-medium text-right">Stock</th>
              <th className="px-4 py-3 font-medium text-right">Price</th>
              <th className="px-4 py-3 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {isLoading ? (
              <tr><td colSpan={6} className="px-4 py-8 text-center text-gray-400">Loading...</td></tr>
            ) : !inv || inv.items.length === 0 ? (
              <tr><td colSpan={6} className="px-4 py-8 text-center text-gray-400">No products found.</td></tr>
            ) : inv.items.map(p => (
              <tr key={p.id} className={`hover:bg-gray-50 ${p.stock === 0 ? 'bg-red-50' : p.stock <= 10 ? 'bg-amber-50' : ''}`}>
                <td className="px-4 py-3 font-medium text-gray-900">{p.name}</td>
                <td className="px-4 py-3 text-gray-600 font-mono text-xs">{p.sku}</td>
                <td className="px-4 py-3 text-gray-600">{p.category}</td>
                <td className="px-4 py-3 text-right">
                  <span className={p.stock === 0 ? 'text-red-600 font-bold' : p.stock <= 10 ? 'text-amber-600 font-bold' : 'text-gray-900'}>
                    {p.stock}
                  </span>
                </td>
                <td className="px-4 py-3 text-right text-gray-900">Rp {p.price.toLocaleString()}</td>
                <td className="px-4 py-3 text-right">
                  <button onClick={() => { setAdjusting(p); setOperation('add'); setQty(1); setCostPrice(''); setNote(''); setError('') }}
                    className="flex items-center gap-1 text-xs px-2 py-1 text-blue-600 hover:bg-blue-50 rounded">
                    <RefreshCcw className="w-3 h-3" /> Adjust
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {inv && inv.total_pages > 1 && (
          <div className="flex items-center justify-between px-4 py-3 border-t bg-gray-50">
            <span className="text-sm text-gray-500">Page {inv.page} of {inv.total_pages}</span>
            <div className="flex gap-2">
              <button disabled={page <= 1} onClick={() => setPage(p => p - 1)} className="px-3 py-1 text-sm border rounded-md disabled:opacity-40">Prev</button>
              <button disabled={page >= inv.total_pages} onClick={() => setPage(p => p + 1)} className="px-3 py-1 text-sm border rounded-md disabled:opacity-40">Next</button>
            </div>
          </div>
        )}
      </div>

      {adjusting && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div className="absolute inset-0 bg-black/40" onClick={() => setAdjusting(null)} />
          <div className="relative bg-white rounded-xl shadow-xl max-w-sm w-full mx-4 p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Adjust Stock — {adjusting.name}</h3>
            <p className="text-sm text-gray-500 mb-3">Current stock: {adjusting.stock}</p>
            <div className="space-y-3">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Operation</label>
                <select value={operation} onChange={e => setOperation(e.target.value as any)}
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm">
                  <option value="add">Add Stock</option>
                  <option value="subtract">Subtract Stock</option>
                  <option value="set">Set Stock</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Quantity</label>
                <input type="number" min={1} value={qty} onChange={e => setQty(Number(e.target.value))}
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm" />
              </div>
              {operation === 'add' && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Cost Price (per unit)</label>
                  <input type="number" step="0.01" value={costPrice} onChange={e => setCostPrice(e.target.value)}
                    placeholder="Harga beli per unit"
                    className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm" />
                  <p className="text-xs text-gray-400 mt-0.5">Diperlukan untuk kalkulasi COGS</p>
                </div>
              )}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Note</label>
                <input value={note} onChange={e => setNote(e.target.value)} placeholder="Reason for adjustment"
                  className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm" />
              </div>
              {error && <p className="text-sm text-red-600">{error}</p>}
              <div className="flex gap-3">
                <button onClick={() => setAdjusting(null)} className="flex-1 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg">Cancel</button>
                <button onClick={() => adjustMutation.mutate({ id: adjusting.id, data: { operation, quantity: qty, cost_price: Number(costPrice) || 0, note } })}
                  disabled={qty < 1 || adjustMutation.isPending}
                  className="flex-1 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg disabled:opacity-50">
                  {adjustMutation.isPending ? 'Adjusting...' : 'Confirm'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
