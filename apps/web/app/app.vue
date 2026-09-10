<script setup lang="ts">
useHead({
  htmlAttrs: { lang: 'zh-Hans' },
  titleTemplate: (title) => (title ? `${title} - VNFest` : 'VNFest 同好会地图')
})

useSeoMeta({
  description: '中日高校视觉小说同好会导航',
  ogSiteName: 'VNFest',
  ogLocale: 'zh_CN',
  ogType: 'website'
})
</script>

<template>
  <div>
    <NuxtRouteAnnouncer />
    <NuxtLoadingIndicator color="var(--color-primary)" />
    <KunMessageProvider />
    <KunAlertProvider />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>
