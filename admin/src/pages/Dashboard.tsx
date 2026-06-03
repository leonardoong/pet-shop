import { useQuery } from '@tanstack/react-query'
import { ShoppingCart, Package, AlertTriangle, DollarSign, PiggyBank, Percent, Archive } from 'lucide-react'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import { dashboardApi } from '@/api/dashboard'
import StatusBadge from '@/components/StatusBadge'
import { useNavigate } from 'react-router-dom'

export default function Dashboard() {
  const navigate = useNavigate()
  const { data, isLoading } = useQuery({
    queryKey: ['admin-dashboard'],
    queryFn: () => dashboardApi.getStats(),
    refetchInterval: 30000,
  })

  const stats = data?.data?.data

  const kpiCards = [
    {
      label: 'Revenue', subtitle: 'Total penjualan',
      value: isLoading ? '...' : `Rp ${(stats?.total_revenue ?? 0).toLocaleString()}`,
      icon: DollarSign, color: 'bg-blue-50 text-blue-600',
    },
    {
      label: 'COGS', subtitle: 'Harga pokok penjualan',
      value: isLoading ? '...' : `Rp ${(stats?.total_cogs ?? 0).toLocaleString()}`,
      icon: ShoppingCart, color: 'bg-amber-50 text-amber-600',
    },
    {
      label: 'Net Income', subtitle: 'Laba bersih',
      value: isLoading ? '...' : `Rp ${(stats?.net_income ?? 0).toLocaleString()}`,
      icon: PiggyBank, color: 'bg-green-50 text-green-600',
    },
    {
      label: 'Profit Margin', subtitle: 'Margin keuntungan',
      value: isLoading ? '...' : `${(stats?.profit_margin ?? 0).toFixed(1)}%`,
      icon: Percent, color: 'bg-purple-50 text-purple-600',
    },
    {
      label: 'Total Products',
      value: isLoading ? '...' : (stats?.total_products ?? 0).toLocaleString(),
      icon: Package, color: 'bg-indigo-50 text-indigo-600',
    },
    {
      label: 'Inventory Value', subtitle: 'Nilai stok saat ini',
      value: isLoading ? '...' : `Rp ${(stats?.inventory_value ?? 0).toLocaleString()}`,
      icon: Archive, color: 'bg-teal-50 text-teal-600',
    },
  ]

  const chartData = stats?.revenue_by_month?.map(r => ({
    month: r.month,
    revenue: r.revenue,
  })) || []

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-gray-900">Dashboard</h2>
        <p className="text-sm text-gray-500 mt-1">Welcome back. Here's what's happening.</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {kpiCards.map(({ label, value, subtitle, icon: Icon, color }) => (
          <div key={label} className="bg-white rounded-xl border border-gray-200 p-5 flex items-center gap-4">
            <div className={`w-11 h-11 rounded-lg flex items-center justify-center ${color}`}>
              <Icon className="w-5 h-5" />
            </div>
            <div>
              <p className="text-sm text-gray-500">{label}</p>
              <p className="text-xl font-bold text-gray-900">{value}</p>
              {subtitle && <p className="text-xs text-gray-400">{subtitle}</p>}
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 bg-white rounded-xl border border-gray-200 p-6">
          <h3 className="text-base font-semibold text-gray-800 mb-4">Monthly Revenue</h3>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={chartData} margin={{ top: 0, right: 0, left: 0, bottom: 0 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                <XAxis dataKey="month" tick={{ fontSize: 12 }} />
                <YAxis tick={{ fontSize: 12 }} />
                <Tooltip />
                <Bar dataKey="revenue" fill="#3b82f6" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h3 className="text-base font-semibold text-gray-800 mb-4">Orders by Status</h3>
          <div className="space-y-2">
            {stats?.orders_by_status && Object.entries(stats.orders_by_status).map(([status, count]) => (
              <div key={status} className="flex justify-between items-center">
                <StatusBadge status={status} />
                <span className="text-sm font-semibold text-gray-900">{count}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h3 className="text-base font-semibold text-gray-800 mb-4">Recent Orders</h3>
          <div className="space-y-3">
            {stats?.recent_orders?.length === 0 ? (
              <p className="text-sm text-gray-400">No orders yet.</p>
            ) : stats?.recent_orders?.map(o => (
              <div key={o.id} className="flex items-center justify-between text-sm border-b pb-2 cursor-pointer hover:bg-gray-50 -mx-2 px-2 rounded"
                onClick={() => navigate('/orders')}>
                <div>
                  <span className="text-gray-900 font-medium">{o.customer_name}</span>
                  <span className="text-gray-400 ml-2 font-mono text-xs">#{o.id.slice(0, 8)}</span>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-gray-900">Rp {o.total_amount.toLocaleString()}</span>
                  <StatusBadge status={o.status} />
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h3 className="text-base font-semibold text-gray-800 mb-4 flex items-center gap-2">
            <AlertTriangle className="w-4 h-4 text-amber-500" /> Low Stock Alerts
          </h3>
          <div className="space-y-3">
            {stats?.low_stock_products?.length === 0 ? (
              <p className="text-sm text-gray-400">All products well stocked.</p>
            ) : stats?.low_stock_products?.map(p => (
              <div key={p.id} className="flex items-center justify-between text-sm border-b pb-2">
                <div>
                  <span className="text-gray-900 font-medium">{p.name}</span>
                  <span className="text-gray-400 ml-2 font-mono text-xs">{p.sku}</span>
                </div>
                <span className={`font-bold ${p.stock === 0 ? 'text-red-600' : 'text-amber-600'}`}>
                  {p.stock} left
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
