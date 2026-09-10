<script setup lang="ts">
import Cookies from 'js-cookie'
import { apiFetch } from '~/utils/api'
import type { UserProfile } from '~/store/user'

definePageMeta({ layout: false })
useKunDisableSeo('OAuth 登录回调')

const route = useRoute()
const error = ref('')
const { store } = useSession()

onMounted(async () => {
  const code = route.query.code as string
  const returnedState = route.query.state as string
  const savedState = Cookies.get('oauth_state')
  const codeVerifier = Cookies.get('oauth_code_verifier')

  Cookies.remove('oauth_state', { path: '/' })
  Cookies.remove('oauth_code_verifier', { path: '/' })

  if (!code) {
    error.value = '未收到授权码'
    return
  }
  if (returnedState !== savedState) {
    error.value = 'State 不匹配，可能存在安全风险'
    return
  }
  if (!codeVerifier) {
    error.value = 'PKCE 验证器丢失，请重新登录'
    return
  }

  try {
    const profile = await apiFetch<UserProfile>('/auth/oauth/callback', {
      method: 'POST',
      body: { code, code_verifier: codeVerifier }
    })
    store.setProfile(profile)
    await navigateTo(consumeOAuthReturnTo() ?? '/user')
  } catch (e) {
    error.value = e instanceof Error ? e.message : '登录失败，请重试'
  }
})
</script>

<template>
  <div
    class="bg-background flex min-h-dvh w-full flex-col items-center justify-center gap-5 px-6 text-center"
  >
    <template v-if="!error">
      <KunIcon
        class="text-primary size-10 animate-spin"
        name="lucide:loader-circle"
      />
      <p class="text-foreground text-lg font-medium">正在登录...</p>
    </template>
    <template v-else>
      <KunIcon class="text-danger size-10" name="lucide:circle-alert" />
      <p class="text-danger text-lg font-medium">{{ error }}</p>
      <KunButton class="mt-2" @click="navigateTo('/auth/login')">
        返回登录
      </KunButton>
    </template>
  </div>
</template>
