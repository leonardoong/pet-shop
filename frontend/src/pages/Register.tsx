import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link, useNavigate } from 'react-router-dom'
import { useState } from 'react'
import { authApi } from '@/api/auth'
import { useAuthStore } from '@/store/authStore'

const PHONE_REGEX = /^(\+62|62|0)8[1-9][0-9]{7,10}$/

const schema = z.object({
  full_name: z.string().min(2, 'Nama minimal 2 karakter'),
  email:     z.string().email('Format email tidak valid'),
  password:  z.string().min(8, 'Password minimal 8 karakter'),
  phone:     z.string()
               .min(1, 'Nomor HP wajib diisi')
               .regex(PHONE_REGEX, 'Format tidak valid. Contoh: 08123456789 atau +6281234567890'),
})

type FormData = z.infer<typeof schema>

export default function Register() {
  const navigate = useNavigate()
  const setAuth = useAuthStore((s) => s.setAuth)
  const [serverError, setServerError] = useState('')

  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<FormData>({
    resolver: zodResolver(schema),
  })

  const onSubmit = async (data: FormData) => {
    setServerError('')
    try {
      const res = await authApi.register(data)
      const { customer, tokens } = res.data.data
      setAuth(customer, tokens.access_token, tokens.access_token)
      navigate('/')
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { message?: string } } })
        ?.response?.data?.message ?? 'Pendaftaran gagal, coba lagi'
      setServerError(msg)
    }
  }

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4 py-10">
      <div className="w-full max-w-md">
        <div className="card p-8">
          <div className="text-center mb-6">
            <span className="text-4xl">🐾</span>
            <h1 className="text-2xl font-bold text-gray-900 mt-2">Buat Akun Baru</h1>
            <p className="text-sm text-gray-500 mt-1">Daftar dan mulai belanja untuk hewan peliharaanmu</p>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Nama Lengkap</label>
              <input
                {...register('full_name')}
                className="input"
                placeholder="Budi Santoso"
              />
              {errors.full_name && <p className="text-xs text-red-500 mt-1">{errors.full_name.message}</p>}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Email</label>
              <input
                {...register('email')}
                type="email"
                autoComplete="email"
                className="input"
                placeholder="nama@email.com"
              />
              {errors.email && <p className="text-xs text-red-500 mt-1">{errors.email.message}</p>}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Nomor HP</label>
              <input
                {...register('phone')}
                type="tel"
                className="input"
                placeholder="08123456789"
              />
              {errors.phone && <p className="text-xs text-red-500 mt-1">{errors.phone.message}</p>}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Password</label>
              <input
                {...register('password')}
                type="password"
                autoComplete="new-password"
                className="input"
                placeholder="Minimal 8 karakter"
              />
              {errors.password && <p className="text-xs text-red-500 mt-1">{errors.password.message}</p>}
            </div>

            {serverError && (
              <div className="text-sm text-red-700 bg-red-50 border border-red-200 rounded-xl px-4 py-3">
                {serverError}
              </div>
            )}

            <button type="submit" disabled={isSubmitting} className="btn-primary w-full">
              {isSubmitting ? 'Memproses...' : 'Daftar Sekarang'}
            </button>
          </form>

          <p className="text-sm text-gray-500 text-center mt-5">
            Sudah punya akun?{' '}
            <Link to="/masuk" className="text-primary-600 font-semibold hover:underline">
              Masuk di sini
            </Link>
          </p>
        </div>
      </div>
    </div>
  )
}
