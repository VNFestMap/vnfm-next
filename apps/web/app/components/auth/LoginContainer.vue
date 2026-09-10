<script setup lang="ts">
import { startOAuthLogin, startOAuthRegister } from '~/utils/oauth-auth'

const route = useRoute()
const { store, fetchMe } = useSession()

const returnTo = computed(() => {
  const raw = String(route.query.redirect || '/user')
  if (raw.startsWith('/') && !raw.startsWith('//')) return raw
  return '/user'
})

onMounted(async () => {
  if (store.isLoggedIn) {
    await navigateTo(returnTo.value, { replace: true })
    return
  }
  const me = await fetchMe()
  if (me) {
    await navigateTo(returnTo.value, { replace: true })
  }
})

const login = () => {
  startOAuthLogin(returnTo.value)
}

const register = () => {
  startOAuthRegister(returnTo.value)
}
</script>

<template>
  <div class="bg-background mx-auto grid min-h-[calc(100dvh-4rem)] max-w-5xl gap-8 lg:grid-cols-2">
    <section class="hidden flex-col justify-center lg:flex">
      <p class="text-primary text-sm font-medium">VNFest</p>
      <h1 class="text-foreground mt-3 text-3xl font-bold">回到你的同好会。</h1>
      <p class="text-default-500 mt-3 text-sm leading-6">
        登录后继续管理社团资料、活动报名与个人账号。账号由 NextMoe·未萌
        统一提供。
      </p>
      <KunLink to="/" class="text-primary mt-8 text-sm">先看看同好会目录</KunLink>
    </section>

    <KunCard class="self-center p-8">
      <p class="text-default-400 text-xs tracking-wide">账号入口</p>
      <h2 class="text-foreground mt-2 text-xl font-bold">欢迎回来</h2>
      <p class="text-default-500 mt-2 text-sm">
        登录后继续管理社团资料、活动报名与个人账号。
      </p>
      <div class="mt-6 flex flex-col gap-3">
        <KunButton color="primary" size="lg" @click="login">
          用 NextMoe 登录
        </KunButton>
        <KunButton variant="flat" size="lg" @click="register">
          注册账号
        </KunButton>
      </div>
      <p class="text-default-400 mt-6 text-xs">
        改名、头像、邮箱与密码请在账户中心完成。
      </p>
    </KunCard>
  </div>
</template>
