<script setup lang="ts">
import { apiFetch } from '~/utils/api'

type Ev = {
  id: number
  title: string
  description: string
  club_name: string
  location: string
  starts_at: string
  registered: boolean
  locked: boolean
}

const route = useRoute()
const { store, fetchMe } = useSession()
const ev = ref<Ev | null>(null)
const message = ref('')

const load = async () => {
  ev.value = await apiFetch<Ev>(`/events/${route.params.id}`)
}

onMounted(async () => {
  if (!store.profile) await fetchMe()
  await load()
})

const toggle = async () => {
  if (!ev.value) return
  const path = ev.value.registered ? 'unregister' : 'register'
  await apiFetch(`/events/${ev.value.id}/${path}`, { method: 'POST' })
  message.value = ev.value.registered ? '已取消报名' : '报名成功'
  await load()
}
</script>

<template>
  <div v-if="ev" class="mx-auto max-w-3xl space-y-6">
    <h1 class="text-foreground text-2xl font-bold">{{ ev.title }}</h1>
    <p class="text-default-500 text-sm">
      {{ ev.club_name }} · {{ ev.starts_at.slice(0, 16).replace('T', ' ') }} ·
      {{ ev.location || '地点待定' }}
    </p>
    <KunCard class="p-6">
      <p class="text-foreground whitespace-pre-wrap text-sm">
        {{ ev.description || '暂无介绍' }}
      </p>
    </KunCard>
    <p v-if="message" class="text-success text-sm">{{ message }}</p>
    <KunButton
      v-if="store.isLoggedIn"
      :color="ev.registered ? 'default' : 'primary'"
      :disabled="ev.locked && !ev.registered"
      @click="toggle"
    >
      {{ ev.locked && !ev.registered ? '报名已截止' : ev.registered ? '取消报名' : '报名参加' }}
    </KunButton>
    <p v-else class="text-default-500 text-sm">
      <KunLink to="/auth/login">登录</KunLink> 后可报名。
    </p>
  </div>
</template>
