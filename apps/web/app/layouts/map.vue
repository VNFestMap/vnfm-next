<script setup lang="ts">
onMounted(() => {
  document.documentElement.classList.add('vnfm-map-lock')
})

onBeforeUnmount(() => {
  document.documentElement.classList.remove('vnfm-map-lock')
})
</script>

<template>
  <div class="bg-background h-dvh overflow-hidden">
    <slot />
  </div>
</template>
