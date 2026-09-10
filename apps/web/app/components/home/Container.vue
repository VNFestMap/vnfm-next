<script setup lang="ts">
const config = useRuntimeConfig()
</script>

<template>
  <div class="mx-auto max-w-3xl space-y-6">
    <div>
      <h1 class="text-foreground text-2xl font-bold">同好会目录</h1>
      <p class="text-default-500 mt-2 text-sm">
        中日高校视觉小说同好会导航。列表将在接入数据后显示。
      </p>
    </div>

    <KunCard class="p-6">
      <p class="text-default-500 text-sm">
        API
        <code class="text-foreground">{{ config.public.apiBase }}</code>
      </p>
    </KunCard>
  </div>
</template>
