<script setup lang="ts">
import { drawCountryMap, type MapCountry } from '~/utils/map-renderer'

const props = defineProps<{
  country: MapCountry
  counts: Record<string, number>
  selected: string
  active: boolean
}>()

const emit = defineEmits<{
  select: [name: string]
}>()

const root = ref<HTMLElement | null>(null)
const svgEl = ref<SVGSVGElement | null>(null)
const { ready, load } = useMapAssets()
const tooltip = ref({ show: false, name: '', x: 0, y: 0 })

let cleanup: (() => void) | null = null
let resizeTimer = 0

const cssVar = (name: string, fallback: string) => {
  const el = root.value
  if (!el) return fallback
  return getComputedStyle(el).getPropertyValue(name).trim() || fallback
}

const emptyFill = () => cssVar('--vnfm-map-empty', '#ffdce9')
const lakeFill = () => cssVar('--vnfm-map-lake', '#ffffff')
const colorMode = useColorMode()

const render = () => {
  if (!import.meta.client || !ready.value || !root.value || !svgEl.value) return
  cleanup?.()
  const width = root.value.clientWidth
  const height = root.value.clientHeight
  if (width < 32 || height < 32) return
  cleanup = drawCountryMap(svgEl.value, {
    country: props.country,
    width,
    height,
    counts: props.counts,
    selected: props.selected,
    emptyFill: emptyFill(),
    lakeFill: lakeFill(),
    onSelect: (name) => emit('select', name),
    onHover: (name, x, y) => {
      tooltip.value = { show: true, name, x: x + 12, y: y + 12 }
    },
    onLeave: () => {
      tooltip.value.show = false
    }
  })
}

const scheduleRender = () => {
  window.clearTimeout(resizeTimer)
  resizeTimer = window.setTimeout(render, 160)
}

onMounted(async () => {
  await load()
  await nextTick()
  render()
  window.addEventListener('resize', scheduleRender)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', scheduleRender)
  window.clearTimeout(resizeTimer)
  cleanup?.()
})

watch(
  () =>
    [
      props.country,
      props.counts,
      props.selected,
      props.active,
      colorMode.value
    ] as const,
  async () => {
    if (!props.active) return
    await nextTick()
    render()
  }
)
</script>

<template>
  <div ref="root" class="absolute inset-0">
    <svg ref="svgEl" class="vnfm-map-svg" aria-label="同好会地图" />
    <div
      v-if="tooltip.show"
      class="vnfm-map-tooltip"
      :style="{ left: `${tooltip.x}px`, top: `${tooltip.y}px` }"
    >
      {{ tooltip.name
      }}{{ counts[tooltip.name] ? ` · ${counts[tooltip.name]}` : '' }}
    </div>
  </div>
</template>
