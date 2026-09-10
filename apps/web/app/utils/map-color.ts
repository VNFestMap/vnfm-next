export const colorByCount = (
  count: number,
  maxCount: number,
  empty: string
) => {
  if (!count) return empty
  const ratio = Math.max(0, Math.min(1, count / Math.max(1, maxCount)))
  if (ratio > 0.75) return '#c2185b'
  if (ratio > 0.5) return '#d94f84'
  if (ratio > 0.25) return '#ec78a5'
  return '#f59cc0'
}

export const ensurePointInsideProvince = (
  pathNode: SVGPathElement,
  box: { x: number; y: number; width: number; height: number },
  preferred: { cx: number; cy: number }
) => {
  const svg = pathNode.ownerSVGElement
  if (!svg || typeof pathNode.isPointInFill !== 'function') return preferred

  const test = (x: number, y: number) => {
    const pt = svg.createSVGPoint()
    pt.x = x
    pt.y = y
    return pathNode.isPointInFill(pt)
  }

  const candidates = [
    [preferred.cx, preferred.cy],
    [box.x + box.width * 0.5, box.y + box.height * 0.62],
    [box.x + box.width * 0.35, box.y + box.height * 0.62],
    [box.x + box.width * 0.65, box.y + box.height * 0.62]
  ] as const

  for (const [x, y] of candidates) {
    if (test(x, y)) return { cx: x, cy: y }
  }
  return preferred
}
