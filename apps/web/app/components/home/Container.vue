<script setup lang="ts">
import { apiFetch } from '~/utils/api'

type Club = {
  id: number
  country: string
  name: string
  school: string
  province: string
  prefecture: string
  city: string
  type: string
}

const q = ref('')
const country = ref('all')
const items = ref<Club[]>([])
const total = ref(0)
const loading = ref(true)

const load = async () => {
  loading.value = true
  try {
    const data = await apiFetch<{ items: Club[]; total: number }>(
      `/clubs?q=${encodeURIComponent(q.value)}&country=${country.value}`
    )
    items.value = data.items || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

onMounted(load)

const region = (c: Club) => {
  if (c.country === 'japan') return c.prefecture || '日本'
  return c.province || c.city || '中国'
}
</script>

<template>
  <div class="mx-auto max-w-3xl space-y-6">
    <div class="flex flex-wrap items-end justify-between gap-3">
      <div>
        <h1 class="text-foreground text-2xl font-bold">同好会目录</h1>
        <p class="text-default-500 mt-2 text-sm">中日高校视觉小说同好会导航</p>
      </div>
      <KunButton color="primary" @click="navigateTo('/clubs/new')">
        创建同好会
      </KunButton>
    </div>

    <div class="flex flex-wrap gap-3">
      <KunInput v-model="q" placeholder="搜索名称、学校、地区" @keyup.enter="load" />
      <KunTab
        v-model="country"
        :items="[
          { textValue: '全部', value: 'all' },
          { textValue: '中国', value: 'china' },
          { textValue: '日本', value: 'japan' }
        ]"
        size="sm"
        variant="solid"
        @update:model-value="load"
      />
      <KunButton variant="flat" @click="load">搜索</KunButton>
    </div>

    <p v-if="loading" class="text-default-500 text-sm">加载中…</p>
    <p v-else-if="!items.length" class="text-default-500 text-sm">暂无同好会。</p>
    <div v-else class="space-y-3">
      <NuxtLink
        v-for="c in items"
        :key="c.id"
        :to="`/clubs/${c.id}`"
        class="block"
      >
        <KunCard class="hover:border-primary-200 p-4 transition-colors">
          <p class="text-foreground font-semibold">{{ c.name }}</p>
          <p class="text-default-500 mt-1 text-sm">
            {{ region(c) }} · {{ c.school || '未填写学校' }}
          </p>
        </KunCard>
      </NuxtLink>
      <p class="text-default-400 text-xs">共 {{ total }} 个</p>
    </div>
  </div>
</template>
