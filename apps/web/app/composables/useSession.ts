import { apiFetch } from '~/utils/api'
import { useUserStore, type UserProfile } from '~/store/user'

export const useSession = () => {
  const store = useUserStore()

  const fetchMe = async () => {
    try {
      const profile = await apiFetch<UserProfile>('/auth/me')
      store.setProfile(profile)
      return profile
    } catch {
      store.reset()
      return null
    }
  }

  const logout = async () => {
    try {
      await apiFetch('/auth/logout', { method: 'POST' })
    } finally {
      store.reset()
    }
  }

  return { fetchMe, logout, store }
}
