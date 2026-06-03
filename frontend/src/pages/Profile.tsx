import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import api from '@/api/axios'
import { useAuthStore } from '@/store/authStore'
import type { ApiResponse, Customer } from '@/types'

const profileSchema = z.object({
  full_name: z.string().min(2, 'Minimal 2 karakter'),
  phone: z.string().min(9, 'Minimal 9 digit'),
})

const emailSchema = z.object({
  email: z.string().email('Format email tidak valid'),
  password: z.string().min(1, 'Password wajib diisi'),
})

const passwordSchema = z.object({
  old_password: z.string().min(1, 'Password saat ini wajib diisi'),
  new_password: z.string().min(8, 'Minimal 8 karakter'),
})

type ProfileData = z.infer<typeof profileSchema>
type EmailData = z.infer<typeof emailSchema>
type PasswordData = z.infer<typeof passwordSchema>

export default function Profile() {
  const customer = useAuthStore(s => s.customer)
  const [msg, setMsg] = useState({ text: '', type: '' })

  const { data, refetch } = useQuery({
    queryKey: ['customer-profile'],
    queryFn: () => api.get<ApiResponse<Customer>>('/customer/me'),
  })

  const profile = data?.data?.data

  const {
    register: regProfile, handleSubmit: handleProfile,
    formState: { errors: errProfile, isSubmitting: subProfile },
  } = useForm<ProfileData>({ resolver: zodResolver(profileSchema), defaultValues: { full_name: '', phone: '' } })

  const {
    register: regEmail, handleSubmit: handleEmail,
    formState: { errors: errEmail, isSubmitting: subEmail },
  } = useForm<EmailData>({ resolver: zodResolver(emailSchema) })

  const {
    register: regPw, handleSubmit: handlePw,
    formState: { errors: errPw, isSubmitting: subPw },
  } = useForm<PasswordData>({ resolver: zodResolver(passwordSchema) })

  const showMsg = (text: string, type: string) => { setMsg({ text, type }); setTimeout(() => setMsg({ text: '', type: '' }), 4000) }

  const onProfile = async (d: ProfileData) => {
    try {
      await api.put('/customer/me', d)
      refetch()
      showMsg('Profil berhasil diperbarui', 'success')
    } catch { showMsg('Gagal memperbarui profil', 'error') }
  }

  const onEmail = async (d: EmailData) => {
    try {
      await api.put('/customer/me/email', d)
      showMsg('Email berhasil diubah', 'success')
    } catch (e: any) {
      showMsg(e?.response?.data?.message || 'Gagal mengubah email', 'error')
    }
  }

  const onPassword = async (d: PasswordData) => {
    try {
      await api.put('/customer/me/password', d)
      showMsg('Password berhasil diubah', 'success')
    } catch (e: any) {
      showMsg(e?.response?.data?.message || 'Password saat ini salah', 'error')
    }
  }

  return (
    <div className="max-w-2xl mx-auto px-4 py-10">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Profil Saya</h1>

      {msg.text && (
        <div className={`mb-4 px-4 py-3 rounded-xl text-sm ${msg.type === 'success' ? 'bg-green-50 text-green-700 border border-green-200' : 'bg-red-50 text-red-700 border border-red-200'}`}>
          {msg.text}
        </div>
      )}

      <div className="space-y-6">
        <div className="card p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Foto Profil</h2>
          <div className="flex items-center gap-4">
            <div className="w-16 h-16 rounded-full bg-primary-100 flex items-center justify-center overflow-hidden">
              {profile?.avatar_url ? (
                <img src={profile.avatar_url} alt="Avatar" className="w-full h-full object-cover" />
              ) : (
                <span className="text-2xl">🐾</span>
              )}
            </div>
            <label className="cursor-pointer text-sm text-primary-600 hover:underline">
              Ubah Foto
              <input type="file" accept="image/*" className="hidden" onChange={async (e) => {
                const file = e.target.files?.[0]
                if (!file) return
                const f = new FormData()
                f.append('file', file)
                try {
                  await api.post('/customer/me/avatar', f, { headers: { 'Content-Type': 'multipart/form-data' } })
                  refetch()
                  showMsg('Foto profil diperbarui', 'success')
                } catch {
                  showMsg('Gagal upload foto', 'error')
                }
              }} />
            </label>
          </div>
        </div>

        <div className="card p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Informasi Pribadi</h2>
          <form onSubmit={handleProfile(onProfile)} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Nama Lengkap</label>
              <input {...regProfile('full_name')} defaultValue={profile?.full_name || customer?.full_name} className="input" />
              {errProfile.full_name && <p className="text-xs text-red-500 mt-1">{errProfile.full_name.message}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Nomor HP</label>
              <input {...regProfile('phone')} defaultValue={profile?.phone || customer?.phone} className="input" />
              {errProfile.phone && <p className="text-xs text-red-500 mt-1">{errProfile.phone.message}</p>}
            </div>
            <p className="text-sm text-gray-500">Email: {profile?.email || customer?.email}</p>
            <button type="submit" disabled={subProfile} className="btn-primary">Simpan Perubahan</button>
          </form>
        </div>

        <div className="card p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Ubah Email</h2>
          <form onSubmit={handleEmail(onEmail)} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Email Baru</label>
              <input {...regEmail('email')} type="email" className="input" placeholder="email@baru.com" />
              {errEmail.email && <p className="text-xs text-red-500 mt-1">{errEmail.email.message}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Password Saat Ini (verifikasi)</label>
              <input {...regEmail('password')} type="password" className="input" />
              {errEmail.password && <p className="text-xs text-red-500 mt-1">{errEmail.password.message}</p>}
            </div>
            <button type="submit" disabled={subEmail} className="btn-primary">Ubah Email</button>
          </form>
        </div>

        <div className="card p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Ubah Password</h2>
          <form onSubmit={handlePw(onPassword)} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Password Saat Ini</label>
              <input {...regPw('old_password')} type="password" className="input" />
              {errPw.old_password && <p className="text-xs text-red-500 mt-1">{errPw.old_password.message}</p>}
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Password Baru</label>
              <input {...regPw('new_password')} type="password" className="input" placeholder="Minimal 8 karakter" />
              {errPw.new_password && <p className="text-xs text-red-500 mt-1">{errPw.new_password.message}</p>}
            </div>
            <button type="submit" disabled={subPw} className="btn-primary">Ubah Password</button>
          </form>
        </div>
      </div>
    </div>
  )
}
