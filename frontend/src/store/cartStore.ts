import { create } from 'zustand'
import { cartApi } from '@/api/cart'
import type { Cart } from '@/types'

interface CartState {
  cart: Cart | null
  isLoading: boolean
  fetch: () => Promise<void>
  addItem: (productId: string, quantity: number) => Promise<void>
  updateItem: (productId: string, quantity: number) => Promise<void>
  removeItem: (productId: string) => Promise<void>
  clear: () => void
  itemCount: number
}

export const useCartStore = create<CartState>((set, get) => ({
  cart: null,
  isLoading: false,

  get itemCount() {
    return get().cart?.items?.reduce((sum, i) => sum + i.quantity, 0) ?? 0
  },

  fetch: async () => {
    set({ isLoading: true })
    try {
      const res = await cartApi.get()
      set({ cart: res.data.data })
    } finally {
      set({ isLoading: false })
    }
  },

  addItem: async (productId, quantity) => {
    const res = await cartApi.addItem(productId, quantity)
    set({ cart: res.data.data })
  },

  updateItem: async (productId, quantity) => {
    const res = await cartApi.updateItem(productId, quantity)
    set({ cart: res.data.data })
  },

  removeItem: async (productId) => {
    await cartApi.removeItem(productId)
    const res = await cartApi.get()
    set({ cart: res.data.data })
  },

  clear: () => set({ cart: null }),
}))
