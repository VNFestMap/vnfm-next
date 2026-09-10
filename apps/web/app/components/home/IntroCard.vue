<script setup lang="ts">
const collapsed = defineModel<boolean>('collapsed', { default: false })
</script>

<template>
  <div
    class="vnfm-intro-card absolute top-4 left-4 z-20 w-[min(300px,calc(100vw-32px))] transition-transform duration-300"
    :class="collapsed ? 'translate-x-[calc(-100%+58px)]' : ''"
  >
    <KunCard class="vnfm-float-card p-4">
      <div class="mb-2 flex items-start justify-between gap-2">
        <h2 class="text-primary text-lg font-semibold">
          全国Galgame同好会地图
        </h2>
        <KunButton size="sm" variant="flat" @click="collapsed = !collapsed">
          {{ collapsed ? '简介' : '收起' }}
        </KunButton>
      </div>
      <div :class="collapsed ? 'pointer-events-none opacity-0' : ''">
        <img
          src="/images/VNF.png"
          alt="VNFest 视觉小说高校学园祭"
          width="640"
          height="196"
          class="mb-3 w-full rounded-lg"
        />
        <p class="text-default-500 text-xs leading-relaxed">
          本网站用于聚合展示全国各省、高校及海外地区的 Galgame /
          视觉小说同好组织信息，支持地图缩放、拖拽、分省查看与切换分类，帮助同好快速找到组织。
        </p>
        <p class="text-default-500 mt-2 text-xs leading-relaxed">
          数据来自各高校同好会及公开信息，欢迎提交新的同好会资料。
        </p>
        <div class="mt-3 flex flex-col gap-2">
          <KunButton size="sm" variant="flat" @click="navigateTo('/clubs/new')">
            创建同好会
          </KunButton>
          <a
            href="https://github.com/VNFestMap/vnfm-next"
            target="_blank"
            rel="noopener noreferrer"
          >
            <KunButton size="sm" variant="flat" class="w-full"
              >开源仓库</KunButton
            >
          </a>
        </div>
        <p class="text-default-400 mt-3 text-[10px] leading-relaxed">
          基于
          <a
            class="text-primary"
            href="https://github.com/HELPMEEADICE/china-bandori-maps"
            target="_blank"
            rel="noopener noreferrer"
            >china-bandori-maps</a
          >
          二次开发，遵循 GPLv3
        </p>
        <p class="text-default-400 mt-1 text-[10px]">
          友情链接：
          <a
            class="text-primary"
            href="https://www.cngal.org/"
            target="_blank"
            rel="noopener noreferrer"
            >CnGal</a
          >
        </p>
        <p class="text-default-400 mt-1 text-[10px]">反馈交流群 185057101</p>
      </div>
    </KunCard>
  </div>
</template>
