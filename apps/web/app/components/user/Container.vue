<script setup lang="ts">
import type { KunTabItem } from '@kungal/ui-vue'
import { apiFetch } from '~/utils/api'
import type { UserProfile } from '~/store/user'

const { store, fetchMe, logout } = useSession()
const loading = ref(true)
const loadError = ref('')
const tab = ref('overview')

const tabs: KunTabItem[] = [
  { textValue: '总览', value: 'overview' },
  { textValue: '账户', value: 'account' },
  { textValue: '偏好设置', value: 'preferences' },
  { textValue: '同好会', value: 'clubs' },
  { textValue: '通知', value: 'notifications' }
]

const reload = async () => {
  loading.value = true
  loadError.value = ''
  const me = await fetchMe()
  loading.value = false
  if (!me) {
    await navigateTo('/auth/login?redirect=/user')
  }
}

onMounted(reload)

const handleLogout = async () => {
  await logout()
  await navigateTo('/')
}

const language = ref('zh')
const theme = ref('system')
const saving = ref(false)

watch(
  () => store.profile,
  (p) => {
    if (!p) return
    language.value = p.language_preference || 'zh'
    theme.value = p.theme_preference || 'system'
  },
  { immediate: true }
)

const savePrefs = async () => {
  saving.value = true
  try {
    const profile = await apiFetch<UserProfile>('/auth/me', {
      method: 'PATCH',
      body: {
        language_preference: language.value,
        theme_preference: theme.value
      }
    })
    store.setProfile(profile)
    const colorMode = useColorMode()
    colorMode.preference = theme.value
  } finally {
    saving.value = false
  }
}

const initial = computed(() => {
  const name = store.profile?.name || '用户'
  return name.slice(0, 1)
})
</script>

<template>
  <div class="mx-auto max-w-5xl space-y-6">
    <div class="flex items-center justify-between gap-3">
      <div>
        <p class="text-default-400 text-xs">VNFest</p>
        <h1 class="text-foreground text-2xl font-bold">用户中心</h1>
      </div>
      <KunButton variant="flat" @click="handleLogout">退出登录</KunButton>
    </div>

    <KunCard v-if="loading" class="p-8 text-center">
      <p class="text-default-500 text-sm">正在连接用户中心后端…</p>
    </KunCard>

    <div v-else-if="store.profile" class="grid gap-6 md:grid-cols-[220px_minmax(0,1fr)]">
      <aside class="space-y-4">
        <KunTab
          v-model="tab"
          :items="tabs"
          orientation="vertical"
          variant="underlined"
          color="primary"
        />
        <KunButton variant="light" class="w-full" @click="navigateTo('/')">
          返回目录
        </KunButton>
      </aside>

      <div class="min-w-0 space-y-4">
        <section v-if="tab === 'overview'" class="space-y-4">
          <KunCard class="flex items-center gap-4 p-6">
            <div
              class="bg-primary-100 text-primary flex size-20 shrink-0 items-center justify-center rounded-xl text-2xl font-bold"
            >
              <img
                v-if="store.profile.avatar"
                :src="store.profile.avatar"
                alt=""
                class="size-20 rounded-xl object-cover"
              />
              <span v-else>{{ initial }}</span>
            </div>
            <div class="min-w-0">
              <h2 class="text-foreground truncate text-xl font-bold">
                {{ store.profile.name || '用户' }}
              </h2>
              <p class="text-default-400 mt-1 text-xs">ID {{ store.profile.id }}</p>
            </div>
          </KunCard>
          <div class="grid gap-4 sm:grid-cols-3">
            <KunCard class="p-4">
              <p class="text-default-400 text-xs">我的同好会</p>
              <p class="text-foreground mt-1 text-2xl font-semibold">0</p>
            </KunCard>
            <KunCard class="p-4">
              <p class="text-default-400 text-xs">未读通知</p>
              <p class="text-foreground mt-1 text-2xl font-semibold">0</p>
            </KunCard>
            <KunCard class="p-4">
              <p class="text-default-400 text-xs">已报名活动</p>
              <p class="text-foreground mt-1 text-2xl font-semibold">0</p>
            </KunCard>
          </div>
        </section>

        <section v-else-if="tab === 'account'" class="space-y-4">
          <KunCard class="space-y-3 p-6">
            <h3 class="text-foreground font-semibold">账户资料</h3>
            <p class="text-default-500 text-sm">
              名称、头像、邮箱与密码由 NextMoe·未萌 账户中心管理。
            </p>
            <p class="text-foreground text-sm">{{ store.profile.name }}</p>
            <a
              class="inline-flex"
              :href="store.profile.account_center_url"
              target="_blank"
              rel="noreferrer"
            >
              <KunButton>在账户中心修改</KunButton>
            </a>
          </KunCard>
          <KunCard class="space-y-3 p-6">
            <h3 class="text-foreground font-semibold">代表同好会</h3>
            <p class="text-default-500 text-sm">
              选择一个用于公开展示的正式隶属同好会。加入同好会后可在此设置。
            </p>
            <p class="text-default-400 text-sm">暂无可展示的正式活跃会籍。</p>
          </KunCard>
        </section>

        <section v-else-if="tab === 'preferences'" class="space-y-4">
          <KunCard class="space-y-4 p-6">
            <h3 class="text-foreground font-semibold">偏好设置</h3>
            <p class="text-default-500 text-sm">语言偏好会跟随本站账号；主题保存在此浏览器。</p>
            <div>
              <p class="text-default-500 mb-2 text-xs">语言</p>
              <KunTab
                v-model="language"
                :items="[
                  { textValue: '中文', value: 'zh' },
                  { textValue: '日本語', value: 'ja' }
                ]"
                variant="solid"
                size="sm"
              />
            </div>
            <div>
              <p class="text-default-500 mb-2 text-xs">外观</p>
              <KunTab
                v-model="theme"
                :items="[
                  { textValue: '浅色', value: 'light' },
                  { textValue: '深色', value: 'dark' },
                  { textValue: '跟随系统', value: 'system' }
                ]"
                variant="solid"
                size="sm"
              />
            </div>
            <KunButton color="primary" :disabled="saving" @click="savePrefs">
              保存偏好
            </KunButton>
          </KunCard>
        </section>

        <section v-else-if="tab === 'clubs'">
          <KunCard class="p-6">
            <h3 class="text-foreground font-semibold">我的同好会</h3>
            <p class="text-default-500 mt-2 text-sm">暂无同好会。从目录页申请加入，或使用绑定码。</p>
          </KunCard>
        </section>

        <section v-else>
          <KunCard class="p-6">
            <h3 class="text-foreground font-semibold">通知中心</h3>
            <p class="text-default-500 mt-2 text-sm">
              暂无通知。审核结果、绑定反馈会出现在这里。
            </p>
          </KunCard>
        </section>
      </div>
    </div>
  </div>
</template>
