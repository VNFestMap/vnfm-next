import * as d3 from 'd3'

type MapDrawer = {
  width: (...args: number[]) => MapDrawer
  height: (...args: number[]) => MapDrawer
  scale: (...args: number[]) => MapDrawer
  language: (...args: string[]) => MapDrawer
  colorDefault: (...args: string[]) => MapDrawer
  colorLake: (...args: string[]) => MapDrawer
  draw: (target: string | Element, options?: Record<string, unknown>) => unknown
}

declare global {
  interface Window {
    d3?: typeof d3
    china?: (() => MapDrawer) & MapDrawer
    japan?: (() => MapDrawer) & MapDrawer
  }
}

const loadScript = (src: string) =>
  new Promise<void>((resolve, reject) => {
    const existing = document.querySelector(`script[src="${src}"]`)
    if (existing) {
      resolve()
      return
    }
    const script = document.createElement('script')
    script.src = src
    script.async = false
    script.onload = () => resolve()
    script.onerror = () => reject(new Error(`failed to load ${src}`))
    document.head.appendChild(script)
  })

export const useMapAssets = () => {
  const ready = useState('vnfm-map-assets-ready', () => false)

  const load = async () => {
    if (!import.meta.client || ready.value) return
    window.d3 = d3
    await loadScript('/maps/china.js')
    await loadScript('/maps/japan.js')
    ready.value = true
  }

  return { ready, load }
}
