<script setup lang="ts">
import type { MapCountry } from '~/utils/map-renderer'
import { apiFetch } from '~/utils/api'

defineProps<{
  mode: 'map' | 'list'
  country: MapCountry
  variant?: 'float' | 'bar'
}>()

const emit = defineEmits<{
  'update:mode': [value: 'map' | 'list']
  'update:country': [value: MapCountry]
}>()

const { store } = useSession()
const colorMode = useColorMode()
const mobileOpen = ref(false)
const unread = ref(0)

const colorModeOptions = [
  { value: 'light', label: '浅色', icon: 'lucide:sun' },
  { value: 'dark', label: '深色', icon: 'lucide:moon' },
  { value: 'system', label: '跟随系统', icon: 'lucide:monitor' }
] as const

const initial = computed(() => store.profile?.name?.slice(0, 1) || '?')

const loadUnread = async () => {
  if (!store.isLoggedIn) {
    unread.value = 0
    return
  }
  try {
    const data = await apiFetch<{ unread: number }>('/notifications/unread')
    unread.value = data.unread || 0
  } catch {
    unread.value = 0
  }
}

onMounted(loadUnread)
watch(() => store.isLoggedIn, loadUnread)

const setMode = (value: 'map' | 'list') => emit('update:mode', value)
const setCountry = (value: MapCountry) => emit('update:country', value)
</script>

<template>
  <div
    :class="variant === 'bar' ? 'w-full' : 'vnfm-top-card absolute top-2 z-30'"
  >
    <KunCard
      :class="
        cn(
          'vnfm-float-card p-3',
          variant === 'bar'
            ? 'rounded-none border-x-0 border-t-0 shadow-none'
            : ''
        )
      "
    >
      <div class="flex items-center gap-2 max-md:flex-wrap">
        <div
          class="bg-default-100 text-default-600 flex size-8 shrink-0 items-center justify-center overflow-hidden rounded-full text-xs font-semibold"
        >
          <img
            v-if="store.profile?.avatar"
            :src="store.profile.avatar"
            :alt="store.profile.name"
            class="size-full object-cover"
          />
          <span v-else>{{ initial }}</span>
        </div>
        <p class="text-foreground min-w-0 flex-1 truncate text-sm font-medium">
          {{ store.profile?.name || '访客' }}
        </p>
        <KunButton
          v-if="unread > 0"
          size="sm"
          is-icon-only
          variant="light"
          aria-label="通知"
          @click="navigateTo('/user')"
        >
          <KunIcon name="lucide:bell" class="size-4" />
        </KunButton>
        <KunButton
          v-if="store.isLoggedIn"
          size="sm"
          variant="flat"
          @click="navigateTo('/user')"
        >
          账号
        </KunButton>
        <KunButton
          v-else
          size="sm"
          color="primary"
          @click="navigateTo('/auth/login')"
        >
          登录 / 注册
        </KunButton>
        <KunPopover position="bottom-end">
          <template #trigger>
            <KunButton
              size="sm"
              is-icon-only
              variant="light"
              aria-label="切换主题"
            >
              <KunIcon name="lucide:sun-moon" class="size-4" />
            </KunButton>
          </template>
          <div class="w-36 py-1">
            <button
              v-for="option in colorModeOptions"
              :key="option.value"
              class="flex w-full items-center gap-3 px-3 py-2 text-sm"
              :class="
                colorMode.preference === option.value
                  ? 'bg-primary-50 text-primary'
                  : 'text-default-500 hover:bg-default-100 hover:text-foreground'
              "
              @click="colorMode.preference = option.value"
            >
              <KunIcon :name="option.icon" class="size-4" />
              <span>{{ option.label }}</span>
            </button>
          </div>
        </KunPopover>
        <div
          class="bg-default-100 flex shrink-0 rounded-lg p-0.5 max-md:min-w-40 max-md:flex-1"
        >
          <button
            type="button"
            class="flex flex-1 items-center justify-center gap-1 rounded-md px-3 py-1.5 text-xs font-medium"
            :class="
              mode === 'map'
                ? 'bg-content1 text-foreground shadow-sm'
                : 'text-default-500'
            "
            @click="setMode('map')"
          >
            <KunIcon name="lucide:map" class="size-3.5" />
            地图
          </button>
          <button
            type="button"
            class="flex flex-1 items-center justify-center gap-1 rounded-md px-3 py-1.5 text-xs font-medium"
            :class="
              mode === 'list'
                ? 'bg-content1 text-foreground shadow-sm'
                : 'text-default-500'
            "
            @click="setMode('list')"
          >
            <KunIcon name="lucide:list" class="size-3.5" />
            列表
          </button>
        </div>
        <KunButton
          class="md:hidden"
          size="sm"
          is-icon-only
          variant="light"
          aria-label="展开导航"
          @click="mobileOpen = !mobileOpen"
        >
          <KunIcon name="lucide:chevron-down" class="size-4" />
        </KunButton>
      </div>
      <div
        :class="
          cn(
            'border-default-200 mt-2 grid-cols-3 gap-1 border-t pt-2',
            mobileOpen ? 'grid md:flex' : 'hidden md:flex'
          )
        "
      >
        <KunButton
          size="sm"
          class="flex-1"
          :color="country === 'china' ? 'primary' : 'default'"
          :variant="country === 'china' ? 'solid' : 'flat'"
          @click="setCountry('china')"
        >
          中国同好会
        </KunButton>
        <KunButton
          size="sm"
          class="flex-1"
          :color="country === 'japan' ? 'primary' : 'default'"
          :variant="country === 'japan' ? 'solid' : 'flat'"
          @click="setCountry('japan')"
        >
          日本同好会
        </KunButton>
        <KunButton
          size="sm"
          class="flex-1"
          variant="flat"
          @click="navigateTo('/events')"
        >
          活动日历
        </KunButton>
      </div>
    </KunCard>
  </div>
</template>
