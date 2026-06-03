import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useState } from 'react'
import axios from 'axios'
import { useSearchParams, Link } from 'react-router-dom'

const schema = z.object({
  new_password: z.string().min(8, 'Minimal 8 karakter'),
  confirm: z.string(),
}).refine(d => d.new_password === d.confirm, { message: 'Password tidak cocok', path: ['confirm'] })

type FormData = z.infer<typeof schema>

export default function ResetPassword() {
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token') || ''
  const [done, setDone] = useState(false)
  const [error, setError] = useState('')

  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<FormData>({
    resolver: zodResolver(schema),
  })

  const onSubmit = async (data: FormData) => {
    setError('')
    try {
      await axios.post('/api/v1/auth/reset-password', { token, new_password: data.new_password })
      setDone(true)
    } catch (e: any) {
      setError(e?.response?.data?.message || 'Token tidak valid atau sudah kadaluarsa')
    }
  }

  if (!token) {
    return <div className="min-h-[70vh] flex items-center justify-center"><p className="text-gray-500">Link tidak valid.</p></div>
  }

  if (done) {
    return (
      <div className="min-h-[70vh] flex items-center justify-center px-4 py-10">
        <div className="w-full max-w-md card p-8 text-center">
          <span className="text-4xl">✅</span>
          <h1 className="text-2xl font-bold text-gray-900 mt-3">Password Berhasil Diubah!</h1>
          <Link to="/masuk" className="inline-block mt-5 text-primary-600 font-semibold hover:underline text-sm">Login dengan password baru</Link>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-[70vh] flex items-center justify-center px-4 py-10">
      <div className="w-full max-w-md card p-8">
        <h1 className="text-2xl font-bold text-gray-900 text-center mb-2">Reset Password</h1>
        <p className="text-sm text-gray-500 text-center mb-6">Masukkan password baru Anda</p>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Password Baru</label>
            <input {...register('new_password')} type="password" className="input" placeholder="Minimal 8 karakter" />
            {errors.new_password && <p className="text-xs text-red-500 mt-1">{errors.new_password.message}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Konfirmasi Password</label>
            <input {...register('confirm')} type="password" className="input" placeholder="Ulangi password" />
            {errors.confirm && <p className="text-xs text-red-500 mt-1">{errors.confirm.message}</p>}
          </div>
          {error && <p className="text-sm text-red-600 bg-red-50 border border-red-200 rounded-lg px-3 py-2">{error}</p>}
          <button type="submit" disabled={isSubmitting} className="btn-primary w-full">
            {isSubmitting ? 'Memproses...' : 'Reset Password'}
          </button>
        </form>
      </div>
    </div>
  )
}
