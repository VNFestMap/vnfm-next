<script setup lang="ts">
import {
  clubRegion,
  clubTypeLabel,
  CLUB_TYPES,
  type HomeClub
} from '~/constants/regions'
import type { MapCountry } from '~/utils/map-renderer'

const props = defineProps<{
  mode: 'map' | 'list'
  country: MapCountry
  regions: { key: string; count: number }[]
  selected: string
  clubs: HomeClub[]
  total: number
  loading: boolean
  q: string
  typeFilter: string
  sort: string
}>()

const emit = defineEmits<{
  'update:mode': [value: 'map' | 'list']
  'update:country': [value: MapCountry]
  'update:selected': [value: string]
  'update:q': [value: string]
  'update:typeFilter': [value: string]
  'update:sort': [value: string]
  'open-drawer': []
}>()

const typeOptions = CLUB_TYPES.map((item) => ({
  value: item.value,
  label: item.label
}))

const title = computed(() => {
  if (props.selected) return props.selected
  return props.country === 'japan' ? '日本同好会' : '全部同好会'
})

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
    class="bg-background absolute inset-0 z-50 flex flex-col overflow-hidden"
  >
    <div class="relative">
      <KunButton
        class="absolute top-3 left-3 z-10 md:hidden"
        size="sm"
        is-icon-only
        variant="flat"
        aria-label="菜单"
        @click="emit('open-drawer')"
      >
        <KunIcon name="lucide:menu" class="size-4" />
      </KunButton>
      <HomeTopBar
        :mode="mode"
        :country="country"
        variant="bar"
        @update:mode="emit('update:mode', $event)"
        @update:country="emit('update:country', $event)"
      />
    </div>

    <div class="min-h-0 flex-1 overflow-hidden p-3 md:flex md:gap-3">
      <aside class="hidden w-[260px] shrink-0 overflow-y-auto md:block">
        <KunCard class="h-full space-y-3 p-4">
          <p class="text-default-400 text-xs font-semibold">站点信息</p>
          <h2 class="text-foreground text-base font-semibold">
            全国Galgame同好会地图
          </h2>
          <p class="text-default-500 text-xs leading-relaxed">
            本网站用于聚合展示全国各省、高校的 Galgame /
            视觉小说同好组织信息，支持地图缩放、拖拽、分省查看与切换分类。
          </p>
          <img
            src="/images/VNF.png"
            alt="VNFest"
            class="bg-content2 w-full rounded-lg p-2"
            width="640"
            height="196"
          />
          <KunButton size="sm" class="w-full" @click="navigateTo('/clubs/new')">
            创建同好会
          </KunButton>
          <KunButton
            size="sm"
            variant="flat"
            class="w-full"
            @click="navigateTo('/events')"
          >
            活动日历
          </KunButton>
          <p class="text-default-400 text-[10px] leading-relaxed">
            基于 china-bandori-maps 二次开发，遵循 GPLv3
          </p>
        </KunCard>
      </aside>

      <nav
        class="mb-3 max-h-28 overflow-x-auto md:mb-0 md:h-full md:max-h-none md:w-[200px] md:overflow-y-auto"
      >
        <KunCard class="h-full p-3">
          <div class="mb-2 flex items-baseline justify-between gap-2">
            <span class="text-foreground text-sm font-bold">地区索引</span>
            <small class="text-default-400 text-[10px]">按收录数量排序</small>
          </div>
          <div class="flex gap-1 md:flex-col">
            <button
              type="button"
              class="rounded-lg px-2.5 py-2 text-left text-sm whitespace-nowrap"
              :class="
                !selected
                  ? 'bg-primary-50 text-primary'
                  : 'text-foreground hover:bg-default-100'
              "
              @click="emit('update:selected', '')"
            >
              <span>全部</span>
              <span class="text-default-400 ml-2 text-xs">{{ total }}</span>
            </button>
            <button
              v-for="region in regions"
              :key="region.key"
              type="button"
              class="flex items-center justify-between gap-2 rounded-lg px-2.5 py-2 text-left text-sm whitespace-nowrap"
              :class="
                selected === region.key
                  ? 'bg-primary-50 text-primary'
                  : 'text-foreground hover:bg-default-100'
              "
              @click="emit('update:selected', region.key)"
            >
              <span>{{ region.key }}</span>
              <span class="text-default-400 text-xs">{{ region.count }}</span>
            </button>
          </div>
        </KunCard>
      </nav>

      <main class="min-h-0 min-w-0 flex-1 overflow-hidden">
        <KunCard class="flex h-full flex-col overflow-hidden p-0">
          <div
            class="border-default-200 flex flex-wrap items-end gap-3 border-b p-4"
          >
            <div class="min-w-48 flex-1">
              <div class="flex items-baseline gap-2">
                <h2 class="text-foreground text-lg font-bold">{{ title }}</h2>
                <span
                  class="bg-primary-50 text-primary rounded-full px-2 py-0.5 text-xs"
                >
                  {{ total }} 个
                </span>
              </div>
              <p class="text-default-400 mt-1 text-[11px]">
                选择地区后查看同好会，也可以按名称或类型筛选。
              </p>
            </div>
            <KunInput
              :model-value="q"
              placeholder="搜索组织名 / 学校"
              class="w-48"
              @update:model-value="emit('update:q', String($event ?? ''))"
            />
            <KunSelect
              :model-value="typeFilter"
              :options="typeOptions"
              size="sm"
              @update:model-value="
                emit('update:typeFilter', String($event ?? 'all'))
              "
            />
          </div>
          <div class="min-h-0 flex-1 overflow-y-auto p-4">
            <p v-if="loading" class="text-default-500 text-sm">加载中…</p>
            <p v-else-if="!sorted.length" class="text-default-500 text-sm">
              暂无同好会。
            </p>
            <div
              v-else
              class="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-3"
            >
              <NuxtLink
                v-for="club in sorted"
                :key="club.id"
                :to="`/clubs/${club.id}`"
                class="block"
              >
                <KunCard class="hover:border-primary-200 h-full p-4">
                  <p class="text-foreground font-semibold">{{ club.name }}</p>
                  <p class="text-default-500 mt-1 text-sm">
                    {{ clubRegion(club) }} · {{ club.school || '未填写学校' }}
                  </p>
                  <p class="text-default-400 mt-2 text-xs">
                    {{ clubTypeLabel(club.type) }}
                  </p>
                </KunCard>
              </NuxtLink>
            </div>
          </div>
        </KunCard>
      </main>
    </div>
  </div>
</template>
