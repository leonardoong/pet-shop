import { useSearchParams, Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { useState, useEffect, useCallback } from 'react'
import { SlidersHorizontal, ChevronLeft, ChevronRight, X } from 'lucide-react'
import { categoriesApi, productsApi } from '@/api/products'
import ProductCard from '@/components/ProductCard'
import type { ProductFilter } from '@/types'

const SORT_OPTIONS = [
  { value: 'newest',     label: 'Terbaru' },
  { value: 'oldest',     label: 'Terlama' },
  { value: 'price_asc',  label: 'Harga: Rendah ke Tinggi' },
  { value: 'price_desc', label: 'Harga: Tinggi ke Rendah' },
]

export default function Shop() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [filterOpen, setFilterOpen] = useState(false)
  const [localSearch, setLocalSearch] = useState(searchParams.get('search') ?? '')

  const category  = searchParams.get('category') ?? ''
  const search    = searchParams.get('search') ?? ''
  const sort      = (searchParams.get('sort') ?? 'newest') as ProductFilter['sort']
  const minPrice  = searchParams.get('min_price') ? Number(searchParams.get('min_price')) : undefined
  const maxPrice  = searchParams.get('max_price') ? Number(searchParams.get('max_price')) : undefined
  const page      = Number(searchParams.get('page') ?? '1')

  const updateParam = useCallback(
    (key: string, value: string | undefined) => {
      setSearchParams((prev) => {
        const next = new URLSearchParams(prev)
        if (value) next.set(key, value)
        else next.delete(key)
        next.delete('page')
        return next
      })
    },
    [setSearchParams],
  )

  // Debounce search input → URL
  useEffect(() => {
    const t = setTimeout(() => {
      updateParam('search', localSearch || undefined)
    }, 400)
    return () => clearTimeout(t)
  }, [localSearch, updateParam])

  const { data: categoriesData } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoriesApi.list().then((r) => r.data.data),
    staleTime: 5 * 60 * 1000,
  })

  const filter: ProductFilter = { category: category || undefined, search: search || undefined, sort, min_price: minPrice, max_price: maxPrice, page, limit: 12 }

  const { data: productsData, isLoading, isPlaceholderData } = useQuery({
    queryKey: ['products', filter],
    queryFn: () => productsApi.list(filter).then((r) => r.data.data),
    placeholderData: (prev) => prev,
  })

  const products    = productsData?.items ?? []
  const totalPages  = productsData?.total_pages ?? 1

  const clearFilters = () => {
    setLocalSearch('')
    setSearchParams({})
  }

  const hasActiveFilters = !!(category || search || minPrice || maxPrice || (sort && sort !== 'newest'))

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div className="flex flex-col lg:flex-row gap-6">

        {/* ── Sidebar filters (desktop) ── */}
        <aside className="hidden lg:block w-56 shrink-0">
          <div className="card p-4 sticky top-24 space-y-5">
            <p className="font-semibold text-gray-800 text-sm">Filter</p>

            {/* Kategori */}
            <div>
              <p className="text-xs font-semibold text-gray-500 uppercase mb-2">Kategori</p>
              <ul className="space-y-1">
                <li>
                  <button
                    onClick={() => updateParam('category', undefined)}
                    className={`w-full text-left text-sm px-2 py-1.5 rounded-lg transition-colors ${!category ? 'bg-primary-100 text-primary-700 font-medium' : 'text-gray-600 hover:bg-gray-50'}`}
                  >
                    Semua Produk
                  </button>
                </li>
                {categoriesData?.map((cat) => (
                  <li key={cat.id}>
                    <button
                      onClick={() => updateParam('category', cat.slug)}
                      className={`w-full text-left text-sm px-2 py-1.5 rounded-lg transition-colors ${category === cat.slug ? 'bg-primary-100 text-primary-700 font-medium' : 'text-gray-600 hover:bg-gray-50'}`}
                    >
                      {cat.name}
                    </button>
                  </li>
                ))}
              </ul>
            </div>

            {/* Harga */}
            <div>
              <p className="text-xs font-semibold text-gray-500 uppercase mb-2">Harga</p>
              <div className="flex gap-2">
                <input
                  type="number"
                  placeholder="Min"
                  defaultValue={minPrice}
                  onBlur={(e) => updateParam('min_price', e.target.value || undefined)}
                  className="w-full text-xs border border-gray-200 rounded-lg px-2 py-1.5 focus:outline-none focus:ring-1 focus:ring-primary-400"
                />
                <input
                  type="number"
                  placeholder="Max"
                  defaultValue={maxPrice}
                  onBlur={(e) => updateParam('max_price', e.target.value || undefined)}
                  className="w-full text-xs border border-gray-200 rounded-lg px-2 py-1.5 focus:outline-none focus:ring-1 focus:ring-primary-400"
                />
              </div>
            </div>

            {hasActiveFilters && (
              <button onClick={clearFilters} className="text-xs text-red-500 hover:underline flex items-center gap-1">
                <X className="w-3 h-3" /> Hapus semua filter
              </button>
            )}
          </div>
        </aside>

        {/* ── Main content ── */}
        <div className="flex-1 min-w-0">
          {/* Toolbar */}
          <div className="flex items-center gap-3 mb-4 flex-wrap">
            <div className="flex-1">
              <input
                type="search"
                value={localSearch}
                onChange={(e) => setLocalSearch(e.target.value)}
                placeholder="Cari produk..."
                className="w-full max-w-sm border border-gray-200 rounded-xl px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary-400"
              />
            </div>

            <button
              onClick={() => setFilterOpen((o) => !o)}
              className="lg:hidden flex items-center gap-1.5 text-sm border border-gray-200 px-3 py-2 rounded-xl hover:bg-gray-50 transition-colors"
            >
              <SlidersHorizontal className="w-4 h-4" /> Filter
            </button>

            <select
              value={sort}
              onChange={(e) => updateParam('sort', e.target.value)}
              className="border border-gray-200 rounded-xl px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary-400 bg-white"
            >
              {SORT_OPTIONS.map((o) => (
                <option key={o.value} value={o.value}>{o.label}</option>
              ))}
            </select>
          </div>

          {/* Mobile filter panel */}
          {filterOpen && (
            <div className="lg:hidden card p-4 mb-4 space-y-4">
              <div className="flex items-center justify-between">
                <p className="font-semibold text-sm">Filter</p>
                <button onClick={() => setFilterOpen(false)}><X className="w-4 h-4 text-gray-500" /></button>
              </div>
              <div>
                <p className="text-xs font-semibold text-gray-500 uppercase mb-2">Kategori</p>
                <div className="flex flex-wrap gap-2">
                  <button
                    onClick={() => { updateParam('category', undefined); setFilterOpen(false) }}
                    className={`text-xs px-3 py-1 rounded-full border transition-colors ${!category ? 'bg-primary-500 text-white border-primary-500' : 'border-gray-200 text-gray-600 hover:border-primary-300'}`}
                  >
                    Semua
                  </button>
                  {categoriesData?.map((cat) => (
                    <button
                      key={cat.id}
                      onClick={() => { updateParam('category', cat.slug); setFilterOpen(false) }}
                      className={`text-xs px-3 py-1 rounded-full border transition-colors ${category === cat.slug ? 'bg-primary-500 text-white border-primary-500' : 'border-gray-200 text-gray-600 hover:border-primary-300'}`}
                    >
                      {cat.name}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          )}

          {/* Active filter chips */}
          {hasActiveFilters && (
            <div className="flex flex-wrap gap-2 mb-4">
              {category && (
                <span className="inline-flex items-center gap-1 text-xs bg-primary-100 text-primary-700 px-2.5 py-1 rounded-full">
                  {categoriesData?.find((c) => c.slug === category)?.name ?? category}
                  <button onClick={() => updateParam('category', undefined)}><X className="w-3 h-3" /></button>
                </span>
              )}
              {search && (
                <span className="inline-flex items-center gap-1 text-xs bg-gray-100 text-gray-600 px-2.5 py-1 rounded-full">
                  "{search}"
                  <button onClick={() => { setLocalSearch(''); updateParam('search', undefined) }}><X className="w-3 h-3" /></button>
                </span>
              )}
            </div>
          )}

          {/* Result count */}
          {!isLoading && (
            <p className="text-xs text-gray-500 mb-3">
              {productsData?.total ?? 0} produk ditemukan
            </p>
          )}

          {/* Product grid */}
          {isLoading ? (
            <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
              {Array.from({ length: 8 }).map((_, i) => (
                <div key={i} className="card overflow-hidden animate-pulse">
                  <div className="aspect-square bg-gray-100" />
                  <div className="p-3 space-y-2">
                    <div className="h-3 bg-gray-100 rounded w-16" />
                    <div className="h-3 bg-gray-100 rounded w-full" />
                    <div className="h-3 bg-gray-100 rounded w-3/4" />
                    <div className="h-5 bg-gray-100 rounded w-24 mt-1" />
                    <div className="h-8 bg-gray-100 rounded-xl w-full mt-2" />
                  </div>
                </div>
              ))}
            </div>
          ) : products.length === 0 ? (
            <div className="text-center py-20 text-gray-400">
              <p className="text-4xl mb-3">🔍</p>
              <p className="font-medium text-gray-600">Produk tidak ditemukan</p>
              <p className="text-sm mt-1">Coba ubah filter atau kata kunci pencarian</p>
              {hasActiveFilters && (
                <button onClick={clearFilters} className="mt-4 text-sm text-primary-600 hover:underline">
                  Hapus semua filter
                </button>
              )}
            </div>
          ) : (
            <div className={`grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4 transition-opacity ${isPlaceholderData ? 'opacity-60' : ''}`}>
              {products.map((product, i) => (
                <ProductCard key={product.id} product={product} index={i} />
              ))}
            </div>
          )}

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-center gap-2 mt-8">
              <button
                disabled={page <= 1}
                onClick={() => setSearchParams((prev) => { const next = new URLSearchParams(prev); next.set('page', String(page - 1)); return next })}
                className="p-2 rounded-lg border border-gray-200 disabled:opacity-40 hover:bg-gray-50 transition-colors"
              >
                <ChevronLeft className="w-4 h-4" />
              </button>

              {Array.from({ length: totalPages }, (_, i) => i + 1)
                .filter((p) => p === 1 || p === totalPages || Math.abs(p - page) <= 1)
                .reduce<(number | '...')[]>((acc, p, idx, arr) => {
                  if (idx > 0 && p - (arr[idx - 1] as number) > 1) acc.push('...')
                  acc.push(p)
                  return acc
                }, [])
                .map((p, i) =>
                  p === '...' ? (
                    <span key={`dots-${i}`} className="px-2 text-gray-400 text-sm">…</span>
                  ) : (
                    <button
                      key={p}
                      onClick={() => setSearchParams((prev) => { const next = new URLSearchParams(prev); next.set('page', String(p)); return next })}
                      className={`w-9 h-9 rounded-lg text-sm font-medium transition-colors ${page === p ? 'bg-primary-500 text-white' : 'border border-gray-200 text-gray-700 hover:bg-gray-50'}`}
                    >
                      {p}
                    </button>
                  ),
                )}

              <button
                disabled={page >= totalPages}
                onClick={() => setSearchParams((prev) => { const next = new URLSearchParams(prev); next.set('page', String(page + 1)); return next })}
                className="p-2 rounded-lg border border-gray-200 disabled:opacity-40 hover:bg-gray-50 transition-colors"
              >
                <ChevronRight className="w-4 h-4" />
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Mobile: link back to top of page for categories */}
      <div className="lg:hidden mt-6 text-center text-sm text-gray-500">
        <Link to="/shop" className="text-primary-600 hover:underline">Lihat semua kategori</Link>
      </div>
    </div>
  )
}
