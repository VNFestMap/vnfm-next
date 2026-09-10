import * as d3 from 'd3'
import {
  BADGE_OFFSET,
  CHINA_BASE,
  CHINA_ID_TO_NAME,
  JAPAN_ID_TO_NAME
} from '~/constants/regions'
import { colorByCount, ensurePointInsideProvince } from '~/utils/map-color'

export type MapCountry = 'china' | 'japan'

type DrawOpts = {
  country: MapCountry
  width: number
  height: number
  counts: Record<string, number>
  selected: string
  emptyFill: string
  lakeFill: string
  onSelect: (name: string) => void
  onHover: (name: string, x: number, y: number) => void
  onLeave: () => void
}

type Zoomable = d3.Selection<SVGGElement, unknown, null, undefined>

const maxOf = (counts: Record<string, number>) => {
  const values = Object.values(counts)
  return values.length ? Math.max(...values) : 1
}

const nameOf = (country: MapCountry, id: string) =>
  country === 'japan' ? JAPAN_ID_TO_NAME[id] : CHINA_ID_TO_NAME[id]

const bindProvince = (
  node: SVGPathElement,
  name: string,
  opts: DrawOpts,
  g: Zoomable
) => {
  const el = d3.select(node)
  el.on('mouseover', null)
  el.on('mouseout', null)
  el.style('cursor', 'pointer')
  el.on('click', (event: MouseEvent) => {
    event.stopPropagation()
    opts.onSelect(name)
    g.selectAll<SVGPathElement, unknown>('.province').classed(
      'selected',
      function (this: SVGPathElement) {
        return nameOf(opts.country, this.id) === name
      }
    )
  })
  el.on('mouseover', (event: MouseEvent) => {
    opts.onHover(name, event.clientX, event.clientY)
  })
  el.on('mousemove', (event: MouseEvent) => {
    opts.onHover(name, event.clientX, event.clientY)
  })
  el.on('mouseout', () => opts.onLeave())
}

const addBadge = (
  layer: Zoomable,
  name: string,
  count: number,
  cx: number,
  cy: number,
  radius: number,
  fontSize: number,
  opts: DrawOpts,
  g: Zoomable
) => {
  const badge = layer
    .append('g')
    .attr('class', 'count-badge')
    .attr('transform', `translate(${cx},${cy})`)

  badge
    .append('circle')
    .attr('r', radius)
    .attr('fill', 'var(--color-primary)')
    .attr('stroke', '#ffffff')
    .attr('stroke-width', 1.5)

  badge
    .append('text')
    .attr('text-anchor', 'middle')
    .attr('dy', '0.35em')
    .attr('font-size', `${fontSize}px`)
    .attr('fill', '#ffffff')
    .attr('font-weight', 'bold')
    .text(count > 99 ? '99+' : count)

  badge.on('click', (event: MouseEvent) => {
    event.stopPropagation()
    opts.onSelect(name)
    g.selectAll<SVGPathElement, unknown>('.province').classed(
      'selected',
      function (this: SVGPathElement) {
        return nameOf(opts.country, this.id) === name
      }
    )
  })
}

const paintAndBadge = (g: Zoomable, opts: DrawOpts, badgeScale: number) => {
  const maxCount = maxOf(opts.counts)
  const badgeLayer = g.append('g').attr('class', 'count-layer')

  g.selectAll<SVGPathElement, unknown>('.province').each(function (
    this: SVGPathElement
  ) {
    const name = nameOf(opts.country, this.id)
    if (!name) return
    const count = opts.counts[name] || 0
    d3.select(this).style('fill', colorByCount(count, maxCount, opts.emptyFill))
    this.classList.toggle('selected', opts.selected === name)
    bindProvince(this, name, opts, g)
    if (!count) return

    const box = this.getBBox()
    if (!box.width || !box.height) return

    if (opts.country === 'china') {
      const preferred = {
        cx: box.x + box.width / (this.id === 'im' ? 2.8 : 2),
        cy: box.y + box.height / (this.id === 'im' ? 1.5 : 2)
      }
      const inside = ensurePointInsideProvince(this, box, preferred)
      const offset = BADGE_OFFSET[this.id] || { dx: 0, dy: 0 }
      const cx = Math.max(
        14,
        Math.min(CHINA_BASE.width - 14, inside.cx + offset.dx)
      )
      const cy = Math.max(
        14,
        Math.min(CHINA_BASE.height - 14, inside.cy + offset.dy)
      )
      addBadge(
        badgeLayer,
        name,
        count,
        cx,
        cy,
        count > 99 ? 13 : 11,
        count > 99 ? 10 : 12,
        opts,
        g
      )
      return
    }

    const radius = Math.max(6, Math.min(25, 8 / badgeScale))
    addBadge(
      badgeLayer,
      name,
      count,
      box.x + box.width / 2,
      box.y + box.height / 2,
      radius,
      Math.max(3, Math.min(18, radius * 0.7)),
      opts,
      g
    )
  })

  return badgeLayer
}

export const drawCountryMap = (
  svgNode: SVGSVGElement,
  opts: DrawOpts
): (() => void) => {
  const svg = d3.select(svgNode)
  svg.selectAll('*').remove()
  svg.attr('width', opts.width).attr('height', opts.height)

  if (opts.country === 'china') {
    if (!window.china) return () => {}
    window
      .china()
      .width(opts.width)
      .height(opts.height)
      .scale(1)
      .language('cn')
      .colorDefault(opts.emptyFill)
      .colorLake(opts.lakeFill)
      .draw(svgNode)
  } else {
    if (!window.japan) return () => {}
    const japanWidth = Math.max(opts.width, 1200)
    const japanHeight = Math.max(opts.height, 1100)
    window
      .japan()
      .width(japanWidth)
      .height(japanHeight)
      .scale(1)
      .language('cn')
      .colorDefault(opts.emptyFill)
      .colorLake(opts.lakeFill)
      .draw(svgNode)
  }

  const g = svg.select<SVGGElement>('g')
  if (g.empty()) return () => {}

  let fitScale = 1
  let offsetX = 0
  let offsetY = 0
  let minScale = 1
  let maxScale = 12

  if (opts.country === 'china') {
    fitScale =
      Math.min(opts.width / CHINA_BASE.width, opts.height / CHINA_BASE.height) *
      0.95
    offsetX = (opts.width - CHINA_BASE.width * fitScale) / 2
    offsetY = (opts.height - CHINA_BASE.height * fitScale) / 2
    minScale = fitScale
    maxScale = fitScale * 12
  } else {
    const japanWidth = Math.max(opts.width, 1200)
    const japanHeight = Math.max(opts.height, 1100)
    fitScale =
      Math.min(opts.width / japanWidth, opts.height / japanHeight) * 1.25
    offsetX =
      (opts.width - japanWidth * fitScale) / 2 + japanWidth * fitScale * 0.14
    offsetY =
      (opts.height - japanHeight * fitScale) / 2 + japanHeight * fitScale * 0.12
    minScale = fitScale * 0.6
    maxScale = fitScale * 20
  }

  const badgeLayer = paintAndBadge(g, opts, fitScale)

  const zoom = d3
    .zoom<SVGSVGElement, unknown>()
    .scaleExtent([minScale, maxScale])
    .on('zoom', (event) => {
      g.attr('transform', event.transform.toString())
      if (opts.country !== 'japan') return
      const currentScale = event.transform.k || 1
      badgeLayer.selectAll<SVGGElement, unknown>('.count-badge').each(function (
        this: SVGGElement
      ) {
        const badge = d3.select(this)
        const circle = badge.select('circle')
        const text = badge.select('text')
        const radius = Math.max(6, Math.min(25, 8 / currentScale))
        circle.attr('r', radius)
        text.attr('font-size', `${Math.max(2, Math.min(18, radius * 0.7))}px`)
      })
    })

  svg.call(zoom).on('dblclick.zoom', null)
  svg.call(
    zoom.transform,
    d3.zoomIdentity.translate(offsetX, offsetY).scale(fitScale)
  )

  return () => {
    svg.on('.zoom', null)
    svg.selectAll('*').remove()
  }
}
