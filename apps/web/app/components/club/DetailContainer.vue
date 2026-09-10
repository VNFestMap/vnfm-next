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
  info?: string
  contact_hidden: boolean
  can_apply: boolean
  my_role?: string
  my_status?: string
}

type Member = {
  id: number
  user_id: number
  user_name: string
  role: string
  status: string
  apply_reason?: string
  apply_role?: string
}

type Code = {
  id: number
  code: string
  max_uses: number
  use_count: number
  is_active: boolean
}

const route = useRoute()
const { store, fetchMe } = useSession()
const club = ref<Club | null>(null)
const error = ref('')
const reason = ref('')
const contact = ref('')
const members = ref<Member[]>([])
const codes = ref<Code[]>([])
const canManage = computed(
  () => club.value?.my_role === 'manager' || club.value?.my_role === 'representative'
)

const id = computed(() => String(route.params.id))

const load = async () => {
  club.value = await apiFetch<Club>(`/clubs/${id.value}`)
  if (canManage.value) {
    members.value = await apiFetch<Member[]>(`/clubs/${id.value}/members`)
    codes.value = await apiFetch<Code[]>(`/clubs/${id.value}/codes`)
  }
}

onMounted(async () => {
  if (!store.profile) await fetchMe()
  try {
    await load()
  } catch (e) {
    error.value = e instanceof Error ? e.message : '加载失败'
  }
})

const apply = async () => {
  await apiFetch(`/clubs/${id.value}/apply`, {
    method: 'POST',
    body: { apply_reason: reason.value, contact_account: contact.value, apply_role: 'member' }
  })
  await load()
}

const review = async (membershipId: number, ok: boolean) => {
  await apiFetch(`/memberships/${membershipId}/${ok ? 'approve' : 'reject'}`, {
    method: 'POST'
  })
  await load()
}

const makeCode = async () => {
  await apiFetch(`/clubs/${id.value}/codes`, { method: 'POST', body: { max_uses: 20 } })
  await load()
}

const region = computed(() => {
  if (!club.value) return ''
  return club.value.country === 'japan'
    ? club.value.prefecture
    : club.value.province || club.value.city
})
</script>

<template>
  <div class="mx-auto max-w-3xl space-y-6">
    <p v-if="error" class="text-danger">{{ error }}</p>
    <template v-else-if="club">
      <div>
        <h1 class="text-foreground text-2xl font-bold">{{ club.name }}</h1>
        <p class="text-default-500 mt-1 text-sm">
          {{ region }} · {{ club.school || '未填写学校' }}
        </p>
        <KunChip v-if="club.my_status" class="mt-2" size="sm">
          {{ club.my_status }} / {{ club.my_role }}
        </KunChip>
      </div>

      <KunCard class="p-6">
        <h2 class="text-foreground font-semibold">介绍</h2>
        <p class="text-default-500 mt-2 whitespace-pre-wrap text-sm">
          {{ club.info || (club.contact_hidden ? '申请绑定后可见联系方式' : '暂无介绍') }}
        </p>
      </KunCard>

      <KunCard v-if="club.can_apply && store.isLoggedIn" class="space-y-3 p-6">
        <h2 class="text-foreground font-semibold">申请加入</h2>
        <KunInput v-model="contact" label="联系方式" />
        <KunInput v-model="reason" label="申请理由" />
        <KunButton color="primary" @click="apply">提交申请</KunButton>
      </KunCard>
      <p v-else-if="!store.isLoggedIn" class="text-default-500 text-sm">
        <KunLink to="/auth/login">登录</KunLink> 后可申请加入。
      </p>

      <KunCard v-if="canManage" class="space-y-4 p-6">
        <div class="flex items-center justify-between">
          <h2 class="text-foreground font-semibold">绑定码</h2>
          <KunButton size="sm" @click="makeCode">生成绑定码</KunButton>
        </div>
        <p v-for="c in codes" :key="c.id" class="text-foreground font-mono text-sm">
          {{ c.code }} · {{ c.use_count }}/{{ c.max_uses }}
          <span class="text-default-400">{{ c.is_active ? '' : '已吊销' }}</span>
        </p>
      </KunCard>

      <KunCard v-if="canManage" class="space-y-3 p-6">
        <h2 class="text-foreground font-semibold">成员与申请</h2>
        <div v-for="m in members" :key="m.id" class="border-default-200 flex items-center justify-between gap-3 border-b py-2 text-sm">
          <div>
            <p class="text-foreground">{{ m.user_name || m.user_id }} · {{ m.role }} · {{ m.status }}</p>
            <p v-if="m.apply_reason" class="text-default-400">{{ m.apply_reason }}</p>
          </div>
          <div v-if="m.status === 'pending'" class="flex gap-2">
            <KunButton size="sm" color="success" @click="review(m.id, true)">通过</KunButton>
            <KunButton size="sm" variant="flat" @click="review(m.id, false)">拒绝</KunButton>
          </div>
        </div>
      </KunCard>
    </template>
  </div>
</template>
