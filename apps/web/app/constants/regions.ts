export const CHINA_PROVINCES = [
  '北京',
  '天津',
  '河北',
  '山西',
  '内蒙古',
  '辽宁',
  '吉林',
  '黑龙江',
  '上海',
  '江苏',
  '浙江',
  '安徽',
  '福建',
  '江西',
  '山东',
  '河南',
  '湖北',
  '湖南',
  '广东',
  '广西',
  '海南',
  '重庆',
  '四川',
  '贵州',
  '云南',
  '西藏',
  '陕西',
  '甘肃',
  '青海',
  '宁夏',
  '新疆',
  '香港',
  '澳门',
  '台湾'
] as const

export const JAPAN_PREFECTURES = [
  '北海道',
  '青森县',
  '岩手县',
  '宫城县',
  '秋田县',
  '山形县',
  '福岛县',
  '茨城县',
  '栃木县',
  '群马县',
  '埼玉县',
  '千叶县',
  '东京都',
  '神奈川县',
  '新潟县',
  '富山县',
  '石川县',
  '福井县',
  '山梨县',
  '长野县',
  '岐阜县',
  '静冈县',
  '爱知县',
  '三重县',
  '滋贺县',
  '京都府',
  '大阪府',
  '兵库县',
  '奈良县',
  '和歌山县',
  '鸟取县',
  '岛根县',
  '冈山县',
  '广岛县',
  '山口县',
  '德岛县',
  '香川县',
  '爱媛县',
  '高知县',
  '福冈县',
  '佐贺县',
  '长崎县',
  '熊本县',
  '大分县',
  '宫崎县',
  '鹿儿岛县',
  '冲绳县'
] as const

export const CHINA_ID_TO_NAME: Record<string, string> = {
  hlj: '黑龙江',
  jl: '吉林',
  ln: '辽宁',
  hb: '河北',
  sd: '山东',
  js: '江苏',
  zj: '浙江',
  ah: '安徽',
  hn: '河南',
  sx: '山西',
  snx: '陕西',
  gs: '甘肃',
  hub: '湖北',
  jx: '江西',
  hun: '湖南',
  gz: '贵州',
  sc: '四川',
  yn: '云南',
  qh: '青海',
  han: '海南',
  cq: '重庆',
  tj: '天津',
  bj: '北京',
  nx: '宁夏',
  im: '内蒙古',
  gx: '广西',
  xj: '新疆',
  tb: '西藏',
  sh: '上海',
  fj: '福建',
  gd: '广东',
  hk: '香港',
  mc: '澳门',
  tw: '台湾'
}

export const JAPAN_ID_TO_NAME: Record<string, string> = {
  'JP-01': '北海道',
  'JP-02': '青森县',
  'JP-03': '岩手县',
  'JP-04': '宫城县',
  'JP-05': '秋田县',
  'JP-06': '山形县',
  'JP-07': '福岛县',
  'JP-08': '茨城县',
  'JP-09': '栃木县',
  'JP-10': '群马县',
  'JP-11': '埼玉县',
  'JP-12': '千叶县',
  'JP-13': '东京都',
  'JP-14': '神奈川县',
  'JP-15': '新潟县',
  'JP-16': '富山县',
  'JP-17': '石川县',
  'JP-18': '福井县',
  'JP-19': '山梨县',
  'JP-20': '长野县',
  'JP-21': '岐阜县',
  'JP-22': '静冈县',
  'JP-23': '爱知县',
  'JP-24': '三重县',
  'JP-25': '滋贺县',
  'JP-26': '京都府',
  'JP-27': '大阪府',
  'JP-28': '兵库县',
  'JP-29': '奈良县',
  'JP-30': '和歌山县',
  'JP-31': '鸟取县',
  'JP-32': '岛根县',
  'JP-33': '冈山县',
  'JP-34': '广岛县',
  'JP-35': '山口县',
  'JP-36': '德岛县',
  'JP-37': '香川县',
  'JP-38': '爱媛县',
  'JP-39': '高知县',
  'JP-40': '福冈县',
  'JP-41': '佐贺县',
  'JP-42': '长崎县',
  'JP-43': '熊本县',
  'JP-44': '大分县',
  'JP-45': '宫崎县',
  'JP-46': '鹿儿岛县',
  'JP-47': '冲绳县'
}

export const CHINA_BASE = { width: 960, height: 700 }

export const BADGE_OFFSET: Record<string, { dx: number; dy: number }> = {
  sh: { dx: 16, dy: -10 },
  hk: { dx: 20, dy: -12 },
  mc: { dx: -18, dy: 10 },
  hb: { dx: 0, dy: 20 },
  im: { dx: 0, dy: 0 }
}

export const CLUB_TYPES = [
  { value: 'all', label: '全部类型' },
  { value: 'region', label: '地区高校联合' },
  { value: 'school', label: '高校同好会' },
  { value: 'vnfest', label: '视觉小说学园祭' }
] as const

export type HomeClub = {
  id: number
  country: string
  name: string
  school: string
  province: string
  prefecture: string
  city: string
  type: string
  created_at?: string
}

export const clubTypeLabel = (type: string) => {
  if (type === 'region') return '地区高校联合'
  if (type === 'vnfest') return '视觉小说学园祭'
  if (type === 'school') return '高校同好会'
  return type || '未分类'
}

export const clubRegion = (club: {
  country: string
  province: string
  prefecture: string
  city: string
}) => {
  if (club.country === 'japan') return club.prefecture || '日本'
  return club.province || club.city || '中国'
}
