import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Plus, Pencil, Trash2, Star, Check, X } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { addressesApi } from '@/api/addresses'
import type { Address } from '@/types'

const schema = z.object({
  label:          z.string().min(1, 'Wajib diisi'),
  recipient_name: z.string().min(1, 'Wajib diisi'),
  phone:          z.string().min(8, 'Nomor telepon tidak valid'),
  street:         z.string().min(1, 'Wajib diisi'),
  city:           z.string().min(1, 'Wajib diisi'),
  province:       z.string().min(1, 'Wajib diisi'),
  postal_code:    z.string().min(4, 'Kode pos tidak valid'),
  is_default:     z.boolean().optional(),
})

type FormData = z.infer<typeof schema>

interface AddressFormProps {
  initial?: Address
  onDone: () => void
}

function AddressForm({ initial, onDone }: AddressFormProps) {
  const qc = useQueryClient()

  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: initial
      ? {
          label:          initial.label,
          recipient_name: initial.recipient_name,
          phone:          initial.phone,
          street:         initial.street,
          city:           initial.city,
          province:       initial.province,
          postal_code:    initial.postal_code,
          is_default:     initial.is_default,
        }
      : { is_default: false },
  })

  const [apiError, setApiError] = useState('')

  const { mutate, isPending: isLoading } = useMutation({
    mutationFn: (data: FormData) =>
      initial
        ? addressesApi.update(initial.id, data).then((r) => r.data)
        : addressesApi.create(data).then((r) => r.data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['addresses'] })
      onDone()
    },
    onError: (err: { response?: { data?: { message?: string } } }) => {
      setApiError(err.response?.data?.message ?? 'Terjadi kesalahan.')
    },
  })

  const field = (label: string, name: keyof FormData, placeholder?: string) => (
    <div>
      <label className="block text-xs font-medium text-gray-600 mb-1">{label}</label>
      <input
        {...register(name)}
        placeholder={placeholder}
        className="w-full border border-gray-200 rounded-xl px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary-400"
      />
      {errors[name] && (
        <p className="text-xs text-red-500 mt-0.5">{errors[name]?.message as string}</p>
      )}
    </div>
  )

  return (
    <form onSubmit={handleSubmit((d) => mutate(d))} className="space-y-3">
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        {field('Label Alamat', 'label', 'Rumah, Kantor, dll.')}
        {field('Nama Penerima', 'recipient_name')}
        {field('Nomor Telepon', 'phone', '08xxxxxxxxxx')}
        {field('Kode Pos', 'postal_code')}
      </div>
      {field('Alamat Lengkap', 'street', 'Jl. ..., No. ...')}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
        {field('Kota', 'city')}
        {field('Provinsi', 'province')}
      </div>

      <label className="flex items-center gap-2 cursor-pointer">
        <input type="checkbox" {...register('is_default')} className="accent-primary-500" />
        <span className="text-sm text-gray-600">Jadikan alamat utama</span>
      </label>

      {apiError && (
        <p className="text-xs text-red-500 bg-red-50 rounded-lg px-3 py-2">{apiError}</p>
      )}

      <div className="flex gap-2 pt-1">
        <button
          type="submit"
          disabled={isLoading}
          className="btn-primary px-5 py-2 rounded-xl text-sm font-semibold disabled:opacity-60"
        >
          {isLoading ? 'Menyimpan...' : initial ? 'Simpan Perubahan' : 'Tambah Alamat'}
        </button>
        <button
          type="button"
          onClick={onDone}
          className="px-5 py-2 rounded-xl text-sm font-semibold border border-gray-200 text-gray-600 hover:bg-gray-50 transition-colors"
        >
          Batal
        </button>
      </div>
    </form>
  )
}

export default function AddressBook() {
  const qc = useQueryClient()
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<Address | null>(null)
  const [deletingId, setDeletingId] = useState<string | null>(null)

  const { data: addresses, isLoading } = useQuery({
    queryKey: ['addresses'],
    queryFn: () => addressesApi.list().then((r) => r.data.data),
  })

  const { mutate: setDefault } = useMutation({
    mutationFn: (id: string) => addressesApi.setDefault(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['addresses'] }),
  })

  const { mutate: deleteAddress } = useMutation({
    mutationFn: (id: string) => addressesApi.delete(id),
    onSuccess: () => {
      setDeletingId(null)
      qc.invalidateQueries({ queryKey: ['addresses'] })
    },
  })

  return (
    <div className="max-w-2xl mx-auto px-4 sm:px-6 py-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-xl font-bold text-gray-900">Buku Alamat</h1>
        {!showForm && !editing && (
          <button
            onClick={() => setShowForm(true)}
            className="btn-primary px-4 py-2 rounded-xl text-sm font-semibold flex items-center gap-1.5"
          >
            <Plus className="w-4 h-4" /> Tambah Alamat
          </button>
        )}
      </div>

      {/* Add form */}
      {showForm && (
        <div className="card p-5 mb-4">
          <p className="font-semibold text-gray-900 mb-4">Alamat Baru</p>
          <AddressForm onDone={() => setShowForm(false)} />
        </div>
      )}

      {isLoading ? (
        <div className="space-y-3 animate-pulse">
          {[1, 2].map((i) => (
            <div key={i} className="card p-5 space-y-2">
              <div className="h-4 bg-gray-100 rounded w-24" />
              <div className="h-3 bg-gray-100 rounded w-full" />
              <div className="h-3 bg-gray-100 rounded w-2/3" />
            </div>
          ))}
        </div>
      ) : !addresses || addresses.length === 0 ? (
        <div className="text-center py-16 text-gray-400">
          <p className="text-4xl mb-3">📍</p>
          <p className="font-medium text-gray-600">Belum ada alamat tersimpan</p>
          <p className="text-sm mt-1 text-gray-400">Tambahkan alamat untuk mempermudah checkout</p>
        </div>
      ) : (
        <div className="space-y-3">
          {addresses.map((addr) => (
            <div key={addr.id} className="card p-5">
              {editing?.id === addr.id ? (
                <>
                  <p className="font-semibold text-gray-900 mb-4">Edit Alamat</p>
                  <AddressForm initial={addr} onDone={() => setEditing(null)} />
                </>
              ) : (
                <>
                  <div className="flex items-start justify-between gap-2 mb-2">
                    <div className="flex items-center gap-2 flex-wrap">
                      <p className="font-semibold text-gray-900 text-sm">{addr.label}</p>
                      {addr.is_default && (
                        <span className="inline-flex items-center gap-0.5 text-xs bg-primary-100 text-primary-700 px-1.5 py-0.5 rounded-full font-medium">
                          <Check className="w-3 h-3" /> Utama
                        </span>
                      )}
                    </div>
                    <div className="flex items-center gap-1 shrink-0">
                      {!addr.is_default && (
                        <button
                          onClick={() => setDefault(addr.id)}
                          title="Jadikan utama"
                          className="p-1.5 text-gray-400 hover:text-primary-500 hover:bg-primary-50 rounded-lg transition-colors"
                        >
                          <Star className="w-3.5 h-3.5" />
                        </button>
                      )}
                      <button
                        onClick={() => setEditing(addr)}
                        className="p-1.5 text-gray-400 hover:text-primary-500 hover:bg-primary-50 rounded-lg transition-colors"
                      >
                        <Pencil className="w-3.5 h-3.5" />
                      </button>
                      <button
                        onClick={() => setDeletingId(addr.id)}
                        className="p-1.5 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors"
                      >
                        <Trash2 className="w-3.5 h-3.5" />
                      </button>
                    </div>
                  </div>

                  <p className="text-sm font-medium text-gray-800">{addr.recipient_name}</p>
                  <p className="text-sm text-gray-600">{addr.phone}</p>
                  <p className="text-sm text-gray-500 leading-relaxed">
                    {addr.street}, {addr.city}, {addr.province} {addr.postal_code}
                  </p>

                  {/* Delete confirmation */}
                  {deletingId === addr.id && (
                    <div className="mt-3 flex items-center gap-3 bg-red-50 rounded-xl px-3 py-2">
                      <p className="text-sm text-red-700 flex-1">Hapus alamat ini?</p>
                      <button
                        onClick={() => deleteAddress(addr.id)}
                        className="text-xs bg-red-500 text-white px-3 py-1 rounded-lg hover:bg-red-600 transition-colors"
                      >
                        Hapus
                      </button>
                      <button
                        onClick={() => setDeletingId(null)}
                        className="text-gray-500 hover:text-gray-700"
                      >
                        <X className="w-4 h-4" />
                      </button>
                    </div>
                  )}
                </>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
