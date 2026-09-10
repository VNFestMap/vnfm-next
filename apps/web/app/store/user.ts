import { defineStore } from 'pinia'

export type UserProfile = {
  id: number
  sub: string
  name: string
  avatar: string
  email?: string
  language_preference: string
  theme_preference: string
  display_membership_id: number | null
  account_center_url: string
}

export const useUserStore = defineStore('vnfm-user', {
  state: () => ({
    profile: null as UserProfile | null
  }),
  getters: {
    isLoggedIn: (s) => Boolean(s.profile?.id)
  },
  actions: {
    setProfile(profile: UserProfile | null) {
      this.profile = profile
    },
    reset() {
      this.profile = null
    }
  },
  persist: true
})
