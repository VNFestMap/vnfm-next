<script setup lang="ts">
import {
  clubRegion,
  clubTypeLabel,
  CLUB_TYPES,
  type HomeClub
} from '~/constants/regions'

const props = defineProps<{
  title: string
  subtitle: string
  meta: string
  clubs: HomeClub[]
  loading: boolean
  q: string
  typeFilter: string
  sort: string
}>()

const emit = defineEmits<{
  'update:q': [value: string]
  'update:typeFilter': [value: string]
  'update:sort': [value: string]
  search: []
}>()

const collapsed = defineModel<boolean>('collapsed', { default: false })

const typeOptions = CLUB_TYPES.map((item) => ({
  value: item.value,
  label: item.label
}))

const sorted = computed(() => {
  const rows = [...props.clubs]
  if (props.sort === 'name_asc') {
    rows.sort((a, b) => a.name.localeCompare(b.name, 'zh'))
  } else if (props.sort === 'type_asc') {
    rows.sort((a, b) => a.type.localeCompare(b.type))
  } else if (props.sort === 'time_desc') {
    rows.sort((a, b) => (b.created_at || '').localeCompare(a.created_at || ''))
  }
  return rows
})
</script>

<template>
  <div
    class="vnfm-selected-card absolute top-4 right-4 z-20 flex max-h-[calc(100dvh-32px)] w-[min(360px,calc(100vw-32px))] flex-col"
    :class="collapsed ? 'w-16 min-w-16' : ''"
  >
    <KunCard
      class="vnfm-float-card flex min-h-0 flex-1 flex-col overflow-hidden p-4"
    >
      <button
        v-if="collapsed"
        type="button"
        class="text-primary absolute inset-0 z-10 flex flex-col items-center justify-center gap-1 text-sm tracking-widest"
        @click="collapsed = false"
      >
        <KunIcon name="lucide:chevron-left" class="size-4" />
        展开
      </button>
      <div
        :class="
          collapsed
            ? 'pointer-events-none opacity-0'
            : 'flex min-h-0 flex-1 flex-col'
        "
      >
        <div class="mb-3 flex items-start justify-between gap-2">
          <div class="min-w-0 text-right">
            <h2 class="text-primary text-base font-semibold">{{ title }}</h2>
            <p class="text-primary mt-1 text-xl font-bold">{{ subtitle }}</p>
            <p class="text-default-500 mt-1 text-xs">{{ meta }}</p>
          </div>
          <KunButton size="sm" variant="flat" @click="collapsed = true"
            >收起</KunButton
          >
        </div>
        <KunInput
          :model-value="q"
          placeholder="搜索组织名 / 学校"
          @update:model-value="emit('update:q', String($event ?? ''))"
          @keyup.enter="emit('search')"
        />
        <div class="mt-2 flex flex-wrap gap-2">
          <KunSelect
            :model-value="typeFilter"
            :options="typeOptions"
            size="sm"
            @update:model-value="
              emit('update:typeFilter', String($event ?? 'all'))
            "
          />
        </div>
        <div class="mt-2 flex flex-wrap gap-1">
          <KunButton
            size="sm"
            :variant="sort === 'default' ? 'solid' : 'flat'"
            :color="sort === 'default' ? 'primary' : 'default'"
            @click="emit('update:sort', 'default')"
          >
            默认
          </KunButton>
          <KunButton
            size="sm"
            :variant="sort === 'time_desc' ? 'solid' : 'flat'"
            :color="sort === 'time_desc' ? 'primary' : 'default'"
            @click="emit('update:sort', 'time_desc')"
          >
            登记时间
          </KunButton>
          <KunButton
            size="sm"
            :variant="sort === 'name_asc' ? 'solid' : 'flat'"
            :color="sort === 'name_asc' ? 'primary' : 'default'"
            @click="emit('update:sort', 'name_asc')"
          >
            名称
          </KunButton>
        </div>
        <div
          class="border-default-200 mt-3 min-h-0 flex-1 overflow-y-auto border-t pt-3"
        >
          <p v-if="loading" class="text-default-500 text-sm">加载中…</p>
          <p v-else-if="!sorted.length" class="text-default-500 text-sm">
            暂无同好会。
          </p>
          <NuxtLink
            v-for="club in sorted"
            :key="club.id"
            :to="`/clubs/${club.id}`"
            class="border-default-200 hover:border-primary mb-2 block rounded-2xl border p-3"
          >
            <p class="text-foreground font-semibold">{{ club.name }}</p>
            <p class="text-default-500 mt-1 text-xs">
              {{ clubRegion(club) }} · {{ club.school || '未填写学校' }} ·
              {{ clubTypeLabel(club.type) }}
            </p>
          </NuxtLink>
        </div>
      </div>
    </KunCard>
  </div>
</template>
