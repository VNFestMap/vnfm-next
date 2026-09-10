<script setup lang="ts">
import { apiFetch } from '~/utils/api'

type Ev = {
  id: number
  title: string
  club_name: string
  starts_at: string
  location: string
  locked: boolean
}

const items = ref<Ev[]>([])
const { store } = useSession()

onMounted(async () => {
  const data = await apiFetch<{ items: Ev[]; total: number }>('/events')
  items.value = data.items || []
})
</script>

<template>
  <div class="mx-auto max-w-3xl space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-foreground text-2xl font-bold">活动</h1>
      <KunButton v-if="store.isLoggedIn" color="primary" @click="navigateTo('/events/new')">
        发布活动
      </KunButton>
    </div>
    <p v-if="!items.length" class="text-default-500 text-sm">暂无活动。</p>
    <NuxtLink v-for="e in items" :key="e.id" :to="`/events/${e.id}`" class="block">
      <KunCard class="p-4">
        <p class="text-foreground font-semibold">{{ e.title }}</p>
        <p class="text-default-500 mt-1 text-sm">
          {{ e.club_name }} · {{ e.starts_at.slice(0, 16).replace('T', ' ') }}
          · {{ e.location || '地点待定' }}
        </p>
      </KunCard>
    </NuxtLink>
  </div>
</template>
