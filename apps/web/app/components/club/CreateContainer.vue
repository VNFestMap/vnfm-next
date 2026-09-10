<script setup lang="ts">
import { apiFetch } from '~/utils/api'

const { store, fetchMe } = useSession()
const error = ref('')
const saving = ref(false)
const form = reactive({
  name: '',
  country: 'china',
  school: '',
  province: '',
  prefecture: '',
  city: '',
  type: 'school',
  info: '',
  contact_hidden: true
})

onMounted(async () => {
  if (!store.isLoggedIn) {
    const me = await fetchMe()
    if (!me) {
      await navigateTo('/auth/login?redirect=/clubs/new')
    }
  }
})

const submit = async () => {
  error.value = ''
  saving.value = true
  try {
    const club = await apiFetch<{ id: number }>('/clubs', {
      method: 'POST',
      body: { ...form }
    })
    await navigateTo(`/clubs/${club.id}`)
  } catch (e) {
    error.value = e instanceof Error ? e.message : '创建失败'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-xl space-y-6">
    <h1 class="text-foreground text-2xl font-bold">创建同好会</h1>
    <KunCard class="space-y-4 p-6">
      <KunInput v-model="form.name" label="名称" placeholder="同好会名称" />
      <KunTab
        v-model="form.country"
        :items="[
          { textValue: '中国', value: 'china' },
          { textValue: '日本', value: 'japan' }
        ]"
        size="sm"
        variant="solid"
      />
      <KunInput v-model="form.school" label="学校" />
      <KunInput
        v-if="form.country === 'china'"
        v-model="form.province"
        label="省份"
      />
      <KunInput v-else v-model="form.prefecture" label="都道府县" />
      <KunInput v-model="form.city" label="城市" />
      <KunInput v-model="form.info" label="介绍 / 联系方式" />
      <label class="text-default-500 flex items-center gap-2 text-sm">
        <input v-model="form.contact_hidden" type="checkbox" />
        联系方式仅成员可见
      </label>
      <p v-if="error" class="text-danger text-sm">{{ error }}</p>
      <KunButton color="primary" :disabled="saving" @click="submit">
        创建并成为负责人
      </KunButton>
    </KunCard>
  </div>
</template>
