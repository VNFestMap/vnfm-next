<script setup lang="ts">
import { type HomeClub } from '~/constants/regions'
import { apiFetch } from '~/utils/api'
import type { MapCountry } from '~/utils/map-renderer'

type RegionCount = { key: string; count: number }

const mode = ref<'map' | 'list'>('map')
const country = ref<MapCountry>('china')
const selected = ref('')
const introCollapsed = ref(false)
const selectedCollapsed = ref(false)
const drawerOpen = ref(false)
const q = ref('')
const typeFilter = ref('all')
const sort = ref('default')
const clubs = ref<HomeClub[]>([])
const regions = ref<RegionCount[]>([])
const total = ref(0)
const loading = ref(true)

const counts = computed(() =>
  Object.fromEntries(regions.value.map((row) => [row.key, row.count]))
)

const selectedTitle = computed(() =>
  selected.value ? `${selected.value}同好会` : '全国Galgame同好会数据'
)
const selectedSubtitle = computed(() => `${total.value} 个组织`)
const selectedMeta = computed(() => {
  const label = country.value === 'japan' ? '日本' : '中国'
  return selected.value ? `${label} · ${selected.value}` : `${label} · 全部`
})

const loadRegions = async () => {
  const data = await apiFetch<{ items: RegionCount[] }>(
    `/clubs/regions?country=${country.value}`
  )
  regions.value = data.items || []
}

const loadClubs = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams({
      country: country.value,
      q: q.value,
      type: typeFilter.value,
      limit: '200'
    })
    if (selected.value) params.set('province', selected.value)
    const data = await apiFetch<{ items: HomeClub[]; total: number }>(
      `/clubs?${params.toString()}`
    )
    clubs.value = data.items || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

const reload = async () => {
  await Promise.all([loadRegions(), loadClubs()])
}

watch(country, async () => {
  selected.value = ''
  await reload()
})

watch([selected, typeFilter], loadClubs)

let searchTimer = 0
watch(q, () => {
  if (!import.meta.client) return
  window.clearTimeout(searchTimer)
  searchTimer = window.setTimeout(loadClubs, 220)
})

onMounted(reload)

const onSelectRegion = (name: string) => {
  selected.value = selected.value === name ? '' : name
  selectedCollapsed.value = false
}
</script>

<template>
  <div class="vnfm-map-page">
    <ClientOnly>
      <HomeMapCanvas
        v-show="mode === 'map'"
        :active="mode === 'map'"
        :country="country"
        :counts="counts"
        :selected="selected"
        @select="onSelectRegion"
      />
    </ClientOnly>

    <template v-if="mode === 'map'">
      <HomeTopBar v-model:mode="mode" v-model:country="country" />
      <HomeIntroCard v-model:collapsed="introCollapsed" />
      <HomeSelectedCard
        v-model:collapsed="selectedCollapsed"
        :title="selectedTitle"
        :subtitle="selectedSubtitle"
        :meta="selectedMeta"
        :clubs="clubs"
        :loading="loading"
        :q="q"
        :type-filter="typeFilter"
        :sort="sort"
        @update:q="q = $event"
        @update:type-filter="typeFilter = $event"
        @update:sort="sort = $event"
        @search="loadClubs"
      />
      <div class="absolute top-36 left-3 z-30 md:hidden">
        <KunButton
          size="sm"
          is-icon-only
          variant="flat"
          aria-label="菜单"
          @click="drawerOpen = true"
        >
          <KunIcon name="lucide:menu" class="size-4" />
        </KunButton>
      </div>
    </template>

    <HomeListMode
      v-else
      v-model:mode="mode"
      v-model:country="country"
      :regions="regions"
      :selected="selected"
      :clubs="clubs"
      :total="total"
      :loading="loading"
      :q="q"
      :type-filter="typeFilter"
      :sort="sort"
      @update:selected="selected = $event"
      @update:q="q = $event"
      @update:type-filter="typeFilter = $event"
      @update:sort="sort = $event"
      @open-drawer="drawerOpen = true"
    />

    <HomeMobileDrawer v-model="drawerOpen" />
  </div>
</template>
