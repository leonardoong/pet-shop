import { useState, useCallback } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Plus, Search, Edit, Trash2, ToggleLeft, ToggleRight } from 'lucide-react'
import { productsApi } from '@/api/products'
import { categoriesApi } from '@/api/categories'
import ProductForm from '@/components/ProductForm'
import ConfirmDialog from '@/components/ConfirmDialog'
import type { Product, UpdateProduct } from '@/types'

export default function ProductsPage() {
  const queryClient = useQueryClient()
  const [search, setSearch] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [page, setPage] = useState(1)
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [editing, setEditing] = useState<Product | null>(null)
  const [deleting, setDeleting] = useState<Product | null>(null)

  const { data, isLoading } = useQuery({
    queryKey: ['admin-products', { search, categoryFilter, statusFilter, page }],
    queryFn: () => productsApi.list({
      search,
      ...(categoryFilter && { category_id: categoryFilter }),
      ...(statusFilter && { is_active: statusFilter }),
      page: String(page),
      limit: '15',
    }),
  })

  const { data: categoriesData } = useQuery({
    queryKey: ['admin-categories'],
    queryFn: () => categoriesApi.list(),
  })

  const products = data?.data?.data
  const categories = categoriesData?.data?.data || []

  const createMutation = useMutation({
    mutationFn: productsApi.create,
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['admin-products'] }); setDrawerOpen(false) },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateProduct }) => productsApi.update(id, data),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['admin-products'] }); setDrawerOpen(false); setEditing(null) },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => productsApi.delete(id),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['admin-products'] }); setDeleting(null) },
  })

  const toggleActiveMutation = useMutation({
    mutationFn: ({ id, is_active }: { id: string; is_active: boolean }) =>
      productsApi.update(id, { is_active }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin-products'] }),
  })

  const handleSubmit = useCallback(async (formData: any) => {
    if (editing) {
      const updateData: UpdateProduct = {}
      if (formData.name !== editing.name) updateData.name = formData.name
      if (formData.category_id !== editing.category_id) updateData.category_id = formData.category_id
      if (formData.description !== editing.description) updateData.description = formData.description || ''
      if (formData.price !== editing.price) updateData.price = formData.price
      if (formData.stock !== editing.stock) updateData.stock = formData.stock
      if (formData.sku !== editing.sku) updateData.sku = formData.sku
      if (formData.image_url !== editing.image_url) updateData.image_url = formData.image_url || ''
      if (formData.is_active !== editing.is_active) updateData.is_active = formData.is_active
      await updateMutation.mutateAsync({ id: editing.id, data: updateData })
    } else {
      await createMutation.mutateAsync(formData)
    }
  }, [editing, createMutation, updateMutation])

  const openCreate = () => { setEditing(null); setDrawerOpen(true) }
  const openEdit = (p: Product) => { setEditing(p); setDrawerOpen(true) }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Products</h2>
          <p className="text-sm text-gray-500 mt-1">Manage your product catalog</p>
        </div>
        <button onClick={openCreate} className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700">
          <Plus className="w-4 h-4" /> Add Product
        </button>
      </div>

      <div className="flex gap-3 flex-wrap">
        <div className="relative flex-1 min-w-[200px]">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input value={search} onChange={e => { setSearch(e.target.value); setPage(1) }} placeholder="Search by name or SKU..."
            className="w-full pl-9 pr-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:outline-none" />
        </div>
        <select value={categoryFilter} onChange={e => { setCategoryFilter(e.target.value); setPage(1) }}
          className="border border-gray-300 rounded-lg px-3 py-2 text-sm">
          <option value="">All Categories</option>
          {categories.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
        </select>
        <select value={statusFilter} onChange={e => { setStatusFilter(e.target.value); setPage(1) }}
          className="border border-gray-300 rounded-lg px-3 py-2 text-sm">
          <option value="">All Status</option>
          <option value="true">Active</option>
          <option value="false">Inactive</option>
        </select>
      </div>

      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 text-gray-600 text-left">
            <tr>
              <th className="px-4 py-3 font-medium">Product</th>
              <th className="px-4 py-3 font-medium">SKU</th>
              <th className="px-4 py-3 font-medium">Category</th>
              <th className="px-4 py-3 font-medium text-right">Price (Sell)</th>
              <th className="px-4 py-3 font-medium text-right">Cost</th>
              <th className="px-4 py-3 font-medium text-right">Stock</th>
              <th className="px-4 py-3 font-medium">Status</th>
              <th className="px-4 py-3 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {isLoading ? (
              <tr><td colSpan={8} className="px-4 py-8 text-center text-gray-400">Loading...</td></tr>
            ) : !products || products.items.length === 0 ? (
              <tr><td colSpan={8} className="px-4 py-8 text-center text-gray-400">No products found.</td></tr>
            ) : products.items.map(p => (
              <tr key={p.id} className="hover:bg-gray-50">
                <td className="px-4 py-3">
                  <div className="flex items-center gap-3">
                    {p.image_url ? (
                      <img src={p.image_url} alt={p.name} className="w-9 h-9 rounded object-cover" onError={e => { (e.target as HTMLImageElement).style.display = 'none' }} />
                    ) : <div className="w-9 h-9 rounded bg-gray-100 flex items-center justify-center text-gray-400 text-xs">N/A</div>}
                    <span className="font-medium text-gray-900">{p.name}</span>
                  </div>
                </td>
                <td className="px-4 py-3 text-gray-600 font-mono text-xs">{p.sku}</td>
                <td className="px-4 py-3 text-gray-600">{p.category}</td>
                <td className="px-4 py-3 text-right text-gray-900">Rp {p.price.toLocaleString()}</td>
                <td className="px-4 py-3 text-right text-gray-500 text-xs">Rp {p.cost_price?.toLocaleString() || '—'}</td>
                <td className="px-4 py-3 text-right">
                  <span className={p.stock === 0 ? 'text-red-600' : p.stock <= 10 ? 'text-amber-600' : 'text-gray-900'}>
                    {p.stock}
                  </span>
                </td>
                <td className="px-4 py-3">
                  <button
                    onClick={() => toggleActiveMutation.mutate({ id: p.id, is_active: !p.is_active })}
                    className={`inline-flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                      p.is_active
                        ? 'bg-green-50 text-green-700 hover:bg-green-100 border border-green-200'
                        : 'bg-gray-50 text-gray-500 hover:bg-gray-100 border border-gray-200'
                    }`}
                  >
                    {p.is_active ? (
                      <><ToggleRight className="w-5 h-5" /> Active</>
                    ) : (
                      <><ToggleLeft className="w-5 h-5" /> Inactive</>
                    )}
                  </button>
                </td>
                <td className="px-4 py-3 text-right">
                  <div className="flex items-center gap-1 justify-end">
                    <button onClick={() => openEdit(p)} className="p-1.5 rounded hover:bg-gray-100 text-gray-400 hover:text-blue-600" title="Edit">
                      <Edit className="w-4 h-4" />
                    </button>
                    <button onClick={() => setDeleting(p)} className="p-1.5 rounded hover:bg-gray-100 text-gray-400 hover:text-red-600" title="Delete">
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {products && products.total_pages > 1 && (
          <div className="flex items-center justify-between px-4 py-3 border-t bg-gray-50">
            <span className="text-sm text-gray-500">Page {products.page} of {products.total_pages} ({products.total} items)</span>
            <div className="flex gap-2">
              <button disabled={page <= 1} onClick={() => setPage(p => p - 1)}
                className="px-3 py-1 text-sm border rounded-md disabled:opacity-40 hover:bg-gray-100">Prev</button>
              <button disabled={page >= products.total_pages} onClick={() => setPage(p => p + 1)}
                className="px-3 py-1 text-sm border rounded-md disabled:opacity-40 hover:bg-gray-100">Next</button>
            </div>
          </div>
        )}
      </div>

      <ProductForm
        key={editing?.id || 'new'}
        open={drawerOpen}
        onClose={() => { setDrawerOpen(false); setEditing(null) }}
        onSubmit={handleSubmit}
        product={editing}
        categories={categories}
      />

      <ConfirmDialog
        open={!!deleting}
        onClose={() => setDeleting(null)}
        onConfirm={() => deleting && deleteMutation.mutate(deleting.id)}
        title="Delete Product"
        message={`Are you sure you want to deactivate "${deleting?.name}"? It will no longer be visible to customers.`}
        confirmLabel="Deactivate"
      />
    </div>
  )
}
