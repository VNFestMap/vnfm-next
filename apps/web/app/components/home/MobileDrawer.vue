<script setup lang="ts">
const open = defineModel<boolean>({ default: false })
</script>

<template>
  <div v-if="open" class="fixed inset-0 z-40 md:hidden">
    <button
      type="button"
      class="absolute inset-0 bg-black/40"
      aria-label="关闭菜单"
      @click="open = false"
    />
    <div
      class="bg-content1 absolute inset-y-0 left-0 flex w-[min(320px,85vw)] flex-col overflow-y-auto p-5"
    >
      <div class="mb-3 flex items-center justify-between">
        <h2 class="text-foreground text-lg font-semibold">
          全国Galgame同好会地图
        </h2>
        <KunButton
          size="sm"
          is-icon-only
          variant="light"
          aria-label="关闭"
          @click="open = false"
        >
          <KunIcon name="lucide:x" class="size-4" />
        </KunButton>
      </div>
      <img
        src="/images/VNF.png"
        alt="VNFest"
        class="mb-3 w-full rounded-xl"
        width="640"
        height="196"
      />
      <p class="text-default-500 text-sm leading-relaxed">
        聚合展示各省、高校的 Galgame /
        视觉小说同好组织，支持地图选省与列表检索。
      </p>
      <div class="mt-4 flex flex-col gap-2">
        <KunButton color="primary" @click="navigateTo('/clubs/new')"
          >创建同好会</KunButton
        >
        <KunButton variant="flat" @click="navigateTo('/events')"
          >活动日历</KunButton
        >
        <a
          href="https://github.com/VNFestMap/vnfm-next"
          target="_blank"
          rel="noopener noreferrer"
        >
          <KunButton variant="flat" class="w-full">开源仓库</KunButton>
        </a>
      </div>
      <p class="text-default-400 mt-4 text-xs">反馈交流群 185057101</p>
    </div>
  </div>
</template>
