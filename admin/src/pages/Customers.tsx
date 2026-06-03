import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Search, Eye, ToggleLeft, ToggleRight } from 'lucide-react'
import { customersApi } from '@/api/customers'
import type { Customer } from '@/types'

export default function CustomersPage() {
  const queryClient = useQueryClient()
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(1)
  const [detail, setDetail] = useState<Customer | null>(null)

  const { data, isLoading } = useQuery({
    queryKey: ['admin-customers', { search, page }],
    queryFn: () => customersApi.list({ search, page: String(page), limit: '15' }),
  })

  const customers = data?.data?.data

  const toggleMutation = useMutation({
    mutationFn: ({ id, is_active }: { id: string; is_active: boolean }) =>
      customersApi.toggleActive(id, is_active),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin-customers'] }),
  })

  return (
    <div className="space-y-4">
      <div>
        <h2 className="text-2xl font-bold text-gray-900">Customers</h2>
        <p className="text-sm text-gray-500 mt-1">Manage customer accounts</p>
      </div>

      <div className="relative flex-1 max-w-md">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
        <input value={search} onChange={e => { setSearch(e.target.value); setPage(1) }} placeholder="Search by name, email, or phone..."
          className="w-full pl-9 pr-3 py-2 border border-gray-300 rounded-lg text-sm" />
      </div>

      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 text-gray-600 text-left">
            <tr>
              <th className="px-4 py-3 font-medium">Name</th>
              <th className="px-4 py-3 font-medium">Email</th>
              <th className="px-4 py-3 font-medium">Phone</th>
              <th className="px-4 py-3 font-medium">Status</th>
              <th className="px-4 py-3 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {isLoading ? (
              <tr><td colSpan={5} className="px-4 py-8 text-center text-gray-400">Loading...</td></tr>
            ) : !customers || customers.items.length === 0 ? (
              <tr><td colSpan={5} className="px-4 py-8 text-center text-gray-400">No customers found.</td></tr>
            ) : customers.items.map(c => (
              <tr key={c.id} className="hover:bg-gray-50">
                <td className="px-4 py-3 font-medium text-gray-900">{c.full_name}</td>
                <td className="px-4 py-3 text-gray-600">{c.email}</td>
                <td className="px-4 py-3 text-gray-600">{c.phone}</td>
                <td className="px-4 py-3">
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${c.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                    {c.is_active ? 'Active' : 'Inactive'}
                  </span>
                </td>
                <td className="px-4 py-3 text-right">
                  <div className="flex items-center gap-1 justify-end">
                    <button onClick={() => setDetail(c)} className="p-1.5 rounded hover:bg-gray-100 text-gray-400 hover:text-blue-600"><Eye className="w-4 h-4" /></button>
                    <button onClick={() => toggleMutation.mutate({ id: c.id, is_active: !c.is_active })}
                      className="p-1.5 rounded hover:bg-gray-100 text-gray-400 hover:text-gray-600" title={c.is_active ? 'Deactivate' : 'Activate'}>
                      {c.is_active ? <ToggleRight className="w-4 h-4 text-green-600" /> : <ToggleLeft className="w-4 h-4" />}
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {customers && customers.total_pages > 1 && (
          <div className="flex items-center justify-between px-4 py-3 border-t bg-gray-50">
            <span className="text-sm text-gray-500">Page {customers.page} of {customers.total_pages}</span>
            <div className="flex gap-2">
              <button disabled={page <= 1} onClick={() => setPage(p => p - 1)} className="px-3 py-1 text-sm border rounded-md disabled:opacity-40">Prev</button>
              <button disabled={page >= customers.total_pages} onClick={() => setPage(p => p + 1)} className="px-3 py-1 text-sm border rounded-md disabled:opacity-40">Next</button>
            </div>
          </div>
        )}
      </div>

      {detail && (
        <div className="fixed inset-0 z-50 flex justify-end">
          <div className="absolute inset-0 bg-black/40" onClick={() => setDetail(null)} />
          <div className="relative w-full max-w-sm bg-white shadow-xl overflow-y-auto p-6">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-lg font-semibold">Customer Detail</h3>
              <button onClick={() => setDetail(null)} className="text-gray-400 hover:text-gray-600 text-xl">&times;</button>
            </div>
            <div className="space-y-3 text-sm">
              <div><span className="text-gray-500">Name</span><p className="text-gray-900 font-medium">{detail.full_name}</p></div>
              <div><span className="text-gray-500">Email</span><p className="text-gray-900">{detail.email}</p></div>
              <div><span className="text-gray-500">Phone</span><p className="text-gray-900">{detail.phone}</p></div>
              <div><span className="text-gray-500">Status</span><p>{detail.is_active ? <span className="text-green-600">Active</span> : <span className="text-red-600">Inactive</span>}</p></div>
              <div><span className="text-gray-500">Joined</span><p className="text-gray-900">{new Date(detail.created_at).toLocaleDateString()}</p></div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
