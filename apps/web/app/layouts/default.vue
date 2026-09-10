<script setup lang="ts">
const colorMode = useColorMode()
const { store, fetchMe } = useSession()

onMounted(() => {
  if (!store.profile) {
    fetchMe()
  }
})

const colorModeOptions = [
  { value: 'light', label: '浅色', icon: 'lucide:sun' },
  { value: 'dark', label: '深色', icon: 'lucide:moon' },
  { value: 'system', label: '跟随系统', icon: 'lucide:monitor' }
] as const

const setColorMode = (mode: string) => {
  colorMode.preference = mode
}
</script>

<template>
  <div class="bg-background flex min-h-screen flex-col">
    <header
      class="border-default-200 bg-content1 sticky top-0 z-30 flex h-16 items-center justify-between gap-2 border-b px-4 shadow-sm md:px-6"
    >
      <div class="flex min-w-0 items-center gap-4">
        <NuxtLink to="/" class="flex min-w-0 items-center gap-2">
          <span class="text-primary truncate text-lg font-bold">VNFest</span>
        </NuxtLink>
        <NuxtLink to="/" class="text-default-500 hover:text-foreground hidden text-sm sm:inline">
          目录
        </NuxtLink>
        <NuxtLink to="/events" class="text-default-500 hover:text-foreground hidden text-sm sm:inline">
          活动
        </NuxtLink>
      </div>

      <div class="flex shrink-0 items-center gap-2">
        <KunButton
          v-if="store.isLoggedIn"
          variant="light"
          @click="navigateTo('/user')"
        >
          用户中心
        </KunButton>
        <KunButton v-else color="primary" @click="navigateTo('/auth/login')">
          登录 / 注册
        </KunButton>
        <KunPopover position="bottom-end">
          <template #trigger>
            <KunButton
              variant="light"
              size="md"
              is-icon-only
              aria-label="切换主题"
            >
              <KunIcon name="lucide:sun-moon" class="size-6" />
            </KunButton>
          </template>

          <div class="w-36 py-1">
            <button
              v-for="option in colorModeOptions"
              :key="option.value"
              class="flex w-full items-center gap-3 px-3 py-2 text-sm transition-colors"
              :class="
                colorMode.preference === option.value
                  ? 'bg-primary-50 text-primary'
                  : 'text-default-500 hover:bg-default-100 hover:text-foreground'
              "
              @click="setColorMode(option.value)"
            >
              <KunIcon :name="option.icon" class="size-4" />
              <span>{{ option.label }}</span>
            </button>
          </div>
        </KunPopover>
      </div>
    </header>

    <main class="flex-1 p-4 md:p-6">
      <slot />
    </main>
  </div>
</template>
