<script setup lang="ts">
import { apiFetch } from '~/utils/api'

type Membership = {
  club_id: number
  club_name: string
  role: string
  status: string
}

const { store, fetchMe } = useSession()
const clubs = ref<Membership[]>([])
const form = reactive({
  club_id: 0,
  title: '',
  description: '',
  location: '',
  starts_at: ''
})
const error = ref('')

onMounted(async () => {
  if (!store.isLoggedIn) {
    const me = await fetchMe()
    if (!me) {
      await navigateTo('/auth/login?redirect=/events/new')
      return
    }
  }
  const rows = await apiFetch<Membership[]>('/me/memberships')
  clubs.value = (rows || []).filter(
    (m) => m.status === 'active' && (m.role === 'manager' || m.role === 'representative')
  )
  if (clubs.value[0]) form.club_id = clubs.value[0].club_id
})

const submit = async () => {
  error.value = ''
  try {
    const starts = form.starts_at ? new Date(form.starts_at).toISOString() : ''
    const ev = await apiFetch<{ id: number }>('/events', {
      method: 'POST',
      body: {
        club_id: Number(form.club_id),
        title: form.title,
        description: form.description,
        location: form.location,
        starts_at: starts
      }
    })
    await navigateTo(`/events/${ev.id}`)
  } catch (e) {
    error.value = e instanceof Error ? e.message : '发布失败'
  }
}
</script>

<template>
  <div class="mx-auto max-w-xl space-y-6">
    <h1 class="text-foreground text-2xl font-bold">发布活动</h1>
    <KunCard class="space-y-4 p-6">
      <p v-if="!clubs.length" class="text-default-500 text-sm">
        只有同好会负责人或管理员可以发布活动。
      </p>
      <template v-else>
        <label class="text-default-500 block text-sm">
          发起同好会
          <select v-model.number="form.club_id" class="border-default-200 mt-1 w-full rounded-md border p-2">
            <option v-for="c in clubs" :key="c.club_id" :value="c.club_id">
              {{ c.club_name }}
            </option>
          </select>
        </label>
        <KunInput v-model="form.title" label="标题" />
        <KunInput v-model="form.location" label="地点" />
        <label class="text-default-500 block text-sm">
          开始时间
          <input v-model="form.starts_at" type="datetime-local" class="border-default-200 mt-1 w-full rounded-md border p-2" />
        </label>
        <KunInput v-model="form.description" label="介绍" />
        <p v-if="error" class="text-danger text-sm">{{ error }}</p>
        <KunButton color="primary" @click="submit">发布</KunButton>
      </template>
    </KunCard>
  </div>
</template>
