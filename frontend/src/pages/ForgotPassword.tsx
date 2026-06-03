import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useState } from 'react'
import axios from 'axios'
import { Link } from 'react-router-dom'

const schema = z.object({
  email: z.string().email('Format email tidak valid'),
})

type FormData = z.infer<typeof schema>

export default function ForgotPassword() {
  const [sent, setSent] = useState(false)
  const [error, setError] = useState('')

  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<FormData>({
    resolver: zodResolver(schema),
  })

  const onSubmit = async (data: FormData) => {
    setError('')
    try {
      await axios.post('/api/v1/auth/forgot-password', data)
      setSent(true)
    } catch {
      setError('Gagal mengirim permintaan, coba lagi')
    }
  }

  if (sent) {
    return (
      <div className="min-h-[70vh] flex items-center justify-center px-4 py-10">
        <div className="w-full max-w-md card p-8 text-center">
          <span className="text-4xl">📧</span>
          <h1 className="text-2xl font-bold text-gray-900 mt-3">Cek Email Anda</h1>
          <p className="text-sm text-gray-500 mt-2">
            Jika email terdaftar, kami akan mengirimkan link untuk mereset password.
          </p>
          <Link to="/masuk" className="inline-block mt-5 text-primary-600 font-semibold hover:underline text-sm">
            Kembali ke login
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-[70vh] flex items-center justify-center px-4 py-10">
      <div className="w-full max-w-md card p-8">
        <h1 className="text-2xl font-bold text-gray-900 text-center mb-2">Lupa Password</h1>
        <p className="text-sm text-gray-500 text-center mb-6">Masukkan email Anda dan kami akan mengirim link reset password</p>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
            <input {...register('email')} type="email" className="input" placeholder="nama@email.com" />
            {errors.email && <p className="text-xs text-red-500 mt-1">{errors.email.message}</p>}
          </div>
          {error && <p className="text-sm text-red-600 bg-red-50 border border-red-200 rounded-lg px-3 py-2">{error}</p>}
          <button type="submit" disabled={isSubmitting} className="btn-primary w-full">
            {isSubmitting ? 'Mengirim...' : 'Kirim Link Reset'}
          </button>
        </form>

        <p className="text-sm text-gray-500 text-center mt-5">
          <Link to="/masuk" className="text-primary-600 font-semibold hover:underline">Kembali ke login</Link>
        </p>
      </div>
    </div>
  )
}
