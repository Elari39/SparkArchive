<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { useAsync } from '@/composables/useAsync'
import { getGraph, getPeople } from '@/api'
import { getPortrait } from '@/config/portraits'
import type { PersonSummary } from '@/api/types'
import PageHeader from '@/components/PageHeader.vue'
import AsyncBoundary from '@/components/AsyncBoundary.vue'

/**
 * 人物关系图 —— 世代分栏布局。
 *
 * 为什么不沿用单环：15 人时环上相邻圆盘的净空只剩 6px，且任意两点的直线弦必穿圆内，
 * 29 条边因此在中心叠成一团，边标签（全部放在弦中点）随之糊成色块。
 *
 * 现在的做法：横轴 = 出生世代，左源右流；带内按邻居重心排序压缩交叉；
 * 连线改绕行弧线，端点缩进到圆盘外沿，标签只在聚焦某个人物时按避让规则出现。
 * 全部参数由人数推导，新增人物无需手工调坐标。
 */

// ── 设计令牌（与 styles/theme.css 同源） ──────────────────────────
const C = {
  ink: '#111111',
  paper: '#f2ede4',
  paperDim: '#e3dbcd',
  red: '#d62828',
  redDeep: '#a4161a',
  steel: '#4a4a4a',
} as const

/** 单色处理的滤镜 id：站点 .portrait-photo 用 grayscale+contrast+sepia，这里用等价的 SVG 滤镜。 */
const MONO_FILTER = 'graph-mono'

// ── 关系编码：7 种 kind 归 4 个族，颜色 / 线型 / 箭头 / 粗细四个通道 ──
interface KindStyle {
  family: '源流' | '论战' | '共事'
  stroke: string
  width: number
  /** 线型；缺省为实线 */
  dash?: string
  /** 双向箭头 */
  both?: boolean
  /** 图例与关系列表里的说明 */
  note: string
}
const KIND: Record<string, KindStyle> = {
  传承: { family: '源流', stroke: C.ink, width: 2.4, note: '实质性的思想继承' },
  支援: { family: '源流', stroke: C.ink, width: 1.8, dash: '2 5', note: '单向的支援与声援' },
  借鉴: { family: '源流', stroke: C.steel, width: 1.8, dash: '8 5', note: '借鉴他处的经验' },
  溯源: { family: '源流', stroke: C.redDeep, width: 1.7, dash: '1.5 5', note: '回溯更早的源头' },
  论战: { family: '论战', stroke: C.red, width: 2.6, dash: '12 6', both: true, note: '长期争论、彼此批评' },
  合作: { family: '共事', stroke: C.steel, width: 3.4, both: true, note: '共同行动、彼此支撑' },
  战友: { family: '共事', stroke: C.steel, width: 2.2, both: true, note: '同一阵营并肩' },
}
const FALLBACK_KIND: KindStyle = { family: '源流', stroke: C.ink, width: 2, dash: '6 4', note: '其他关系' }
const kindStyle = (kind: string): KindStyle => KIND[kind] ?? FALLBACK_KIND

/** 关系列表与图例的排序：主脉在前，孤例在后。 */
const KIND_ORDER = ['传承', '支援', '借鉴', '溯源', '论战', '合作', '战友']
const kindRank = (kind: string) => {
  const i = KIND_ORDER.indexOf(kind)
  return i === -1 ? KIND_ORDER.length : i
}

/** 颜色 → 箭头 marker id（marker 不能继承 stroke，只能逐个定义）。 */
const MARKER: Record<string, string> = {
  [C.ink]: 'gm-ink',
  [C.red]: 'gm-red',
  [C.steel]: 'gm-steel',
  [C.redDeep]: 'gm-reddeep',
}
const MARKS = Object.entries(MARKER).map(([color, id]) => ({ color, id }))

/** 图例分组（同一族的多个 kind 合成一行展示）。 */
const LEGEND: { label: string; desc: string; kinds: string[]; stroke: string; width: number; dash?: string; both?: boolean }[] = [
  { label: '传承 · 单向', desc: '实质性的思想继承，实线中粗，构成主脉。', kinds: ['传承'], stroke: C.ink, width: 2.4 },
  { label: '支援 / 借鉴 / 溯源', desc: '单向的局部影响，细的虚线或点线；溯源用深红回指。', kinds: ['支援', '借鉴', '溯源'], stroke: C.ink, width: 1.8, dash: '2 5' },
  { label: '论战 · 双向', desc: '长期争论、彼此批评，红色长虚线。', kinds: ['论战'], stroke: C.red, width: 2.6, dash: '12 6', both: true },
  { label: '合作 / 战友 · 双向', desc: '共同行动或同一阵营，钢色双向，靠粗细区分。', kinds: ['合作', '战友'], stroke: C.steel, width: 3.4, both: true },
]

// ── 布局常量 ────────────────────────────────────────────────────
/** 出生年间隔达到该值即断代（本数据最小跨代 17 年、最大同代 11 年，阈值取两者之间）。 */
const BAND_GAP_YEARS = 15
const MARGIN_X = 100
const MARGIN_Y = 92
const BAND_TOP = 62
/** 姓名牌字号（viewBox 单位，非屏幕像素）；有效字号守卫以它为准。 */
const NAME_FS = 15
const ROLE_FS = 12
const DEG_FS = 12
/** 屏幕上姓名牌允许的最小有效字号，低于它就把默认缩放抬上去。 */
const NAME_FS_FLOOR = 12
const ZOOM_MAX = 6

const clamp = (v: number, lo: number, hi: number) => Math.min(Math.max(v, lo), hi)

/** 近似文本宽度：CJK 按 1em、其余按 0.56em，仅用于给文字底板留位。 */
const textWidth = (text: string, fs: number) =>
  [...text].reduce((sum, ch) => sum + (ch.codePointAt(0)! > 0x2e80 ? fs : fs * 0.56), 0)

/** 显示名取最后一个「·」之后的部分：「菲德尔·卡斯特罗」→「卡斯特罗」。
 *  原先的 name.slice(-3) 会把它截成无法辨认的「斯特罗」。 */
const shortName = (name: string) => {
  const tail = name.split('·').pop()?.trim()
  return tail || name
}

// ── 取数：/relations 与 /people 并行，后者失败时降级为只显示姓名 ──
const { data, loading, error, reload } = useAsync(async (signal) => {
  const [graph, people] = await Promise.all([
    getGraph(signal),
    getPeople(signal).catch(() => [] as PersonSummary[]),
  ])
  return { graph, index: new Map(people.map((p) => [p.slug, p])) }
})

// ── 图数据 ──────────────────────────────────────────────────────
interface NodeView {
  slug: string
  name: string
  /** 姓名牌上的显示名 */
  short: string
  epithet: string
  nationality: string
  /** 「1870–1924」，缺人物详情时为空 */
  life: string
  /** 出生年；缺失为 0 */
  year: number
  degree: number
  /** 头像 webp；缺失为空串（回退为姓氏首字） */
  portrait: string
}

const nodes = computed<NodeView[]>(() => {
  const g = data.value?.graph
  if (!g) return []
  const index = data.value?.index
  const degree = new Map<string, number>()
  for (const e of g.edges) {
    degree.set(e.from, (degree.get(e.from) ?? 0) + 1)
    degree.set(e.to, (degree.get(e.to) ?? 0) + 1)
  }
  return g.nodes.map((n) => {
    const p = index?.get(n.slug)
    const y = Number((p?.birth ?? '').slice(0, 4))
    return {
      slug: n.slug,
      name: n.name,
      short: shortName(n.name),
      epithet: n.epithet || p?.epithet || '',
      nationality: p?.nationality ?? '',
      life: p?.birth && p.death ? `${p.birth.slice(0, 4)}–${p.death.slice(0, 4)}` : '',
      year: Number.isFinite(y) && y > 0 ? y : 0,
      degree: degree.get(n.slug) ?? 0,
      portrait: getPortrait(p?.portrait_key ?? n.slug)?.webp ?? '',
    }
  })
})

/** 无向邻接表，用于带内重心排序。 */
const adjacency = computed(() => {
  const m = new Map<string, Set<string>>()
  for (const e of data.value?.graph.edges ?? []) {
    if (!m.has(e.from)) m.set(e.from, new Set())
    if (!m.has(e.to)) m.set(e.to, new Set())
    m.get(e.from)!.add(e.to)
    m.get(e.to)!.add(e.from)
  }
  return m
})

/**
 * 世代分带 + 带内排序。
 *
 * 1) 按出生年升序（缺失记 0，自成一「未系年」带）；
 * 2) 相邻出生年间隔 ≥ 15 年即断带；
 * 3) 带内初序按度数降序（把枢纽人物排到带内中段）；
 * 4) 8 轮「下扫 + 上扫」重心排序：每带按其邻居在相邻带中的平均序号重排，
 *    把跨带连线尽量拉成近水平的短弧。
 */
const bands = computed<NodeView[][]>(() => {
  const sorted = [...nodes.value].sort(
    (a, b) => a.year - b.year || b.degree - a.degree || a.name.localeCompare(b.name),
  )
  const out: NodeView[][] = []
  for (const n of sorted) {
    const cur = out[out.length - 1]
    if (!cur || n.year - cur[cur.length - 1].year >= BAND_GAP_YEARS) out.push([n])
    else cur.push(n)
  }
  for (const b of out) {
    b.sort((x, y) => y.degree - x.degree || x.year - y.year || x.name.localeCompare(y.name))
  }

  const bandOf = new Map<string, number>()
  out.forEach((b, i) => b.forEach((n) => bandOf.set(n.slug, i)))
  const idx = new Map<string, number>()
  const sync = () => out.forEach((b) => b.forEach((n, i) => idx.set(n.slug, i)))
  sync()

  const bary = (slug: string, ref: number) => {
    const nbrs = [...(adjacency.value.get(slug) ?? [])].filter((s) => bandOf.get(s) === ref)
    if (!nbrs.length) return (idx.get(slug) ?? 0) + 0.5 // 无邻居则原地不动
    return nbrs.reduce((sum, s) => sum + (idx.get(s) ?? 0), 0) / nbrs.length
  }

  for (let it = 0; it < 8; it++) {
    const forward = it % 2 === 0
    const order = out.map((_, i) => i)
    if (!forward) order.reverse()
    for (const bi of order) {
      const ref = forward ? bi - 1 : bi + 1
      if (ref < 0 || ref >= out.length) continue
      out[bi].sort((x, y) => bary(x.slug, ref) - bary(y.slug, ref))
      sync()
    }
  }
  return out
})

/** 画布与节距：全部由人数、带数、最大带人数推导，不写死。 */
const metrics = computed(() => {
  const bs = bands.value
  const n = nodes.value.length
  const bandCount = Math.max(bs.length, 1)
  const maxPerBand = Math.max(1, ...bs.map((b) => b.length))
  // 人越多圆盘略收，把空间还给排布；否则逐条边长会互相压住。
  const d = clamp(76 - 0.8 * n, 48, 76)
  const vPitch = d + 58
  const hPitch = d + 168
  const W = Math.max(560, (bandCount - 1) * hPitch + 2 * MARGIN_X)
  const H = maxPerBand * vPitch + 2 * MARGIN_Y
  return { d, r: d / 2, vPitch, hPitch, W, H, bandCount, maxPerBand }
})

/** 节点坐标：带内垂直居中，带横坐标等距。half 为姓名牌／角色行的半宽，
 *  连线端点要靠它绕开圆盘下方的文字区。 */
const position = computed(() => {
  const { vPitch, hPitch, maxPerBand } = metrics.value
  const pos = new Map<string, { x: number; y: number; half: number }>()
  bands.value.forEach((band, bi) => {
    const x = MARGIN_X + bi * hPitch
    const offset = (maxPerBand - band.length) / 2
    band.forEach((n, i) => {
      const half =
        Math.max(
          textWidth(n.short, NAME_FS) + 18,
          textWidth(n.nationality, ROLE_FS),
        ) / 2
      pos.set(n.slug, { x, y: MARGIN_Y + (i + offset + 0.5) * vPitch, half })
    })
  })
  return pos
})

/** 带世代标题（出生年区间）。注意带内已按重心重排，这里必须取极值而非首尾。 */
const bandMeta = computed(() =>
  bands.value.map((band) => {
    const years = band.map((n) => n.year).filter((y) => y > 0)
    const lo = years.length ? Math.min(...years) : 0
    const hi = years.length ? Math.max(...years) : 0
    const label = years.length ? (lo === hi ? String(lo) : `${lo}–${hi}`) : '未系年'
    return { label, count: band.length }
  }),
)

const placedNodes = computed(() =>
  bands.value.flatMap((band) =>
    band.map((n) => ({ ...n, ...(position.value.get(n.slug) ?? { x: 0, y: 0, half: 0 }) })),
  ),
)

// ── 连线几何：绕行弧线 ───────────────────────────────────────────
interface EdgeView {
  i: number
  from: string
  to: string
  kind: string
  label: string
  note: string
  style: KindStyle
  /** 起点 / 终点（已按圆盘半径缩进） */
  x1: number
  y1: number
  x2: number
  y2: number
  /** 控制点 */
  cx: number
  cy: number
  /** 弦的法向单位向量，标签试位时沿它偏移 */
  nx: number
  ny: number
}

/** 姓名牌下沿到圆心的距离（牌高 25 + 上下留白）；角色行再往下一点。 */
const PLATE_TOP = 4
const PLATE_H = 25
const ROLE_H = 19

/**
 * 沿单位方向 (ox, oy) 从节点圆心出发，走到「圆盘 ∪ 姓名牌 ∪ 角色行」轮廓之外所需的最小距离。
 * 只按圆盘半径缩进会让向下的连线起点落在姓名牌／角色行上，箭头压在文字上。
 */
const footprintExit = (ox: number, oy: number, r: number, half: number) => {
  let t = r
  // 下方矩形：[−half, half] × [r+PLATE_TOP, r+PLATE_TOP+PLATE_H+ROLE_H]
  const y0 = r + PLATE_TOP
  const y1 = y0 + PLATE_H + ROLE_H
  const slab = (o: number, lo: number, hi: number): [number, number] | null => {
    if (Math.abs(o) < 1e-6) return lo <= 0 && 0 <= hi ? [-Infinity, Infinity] : null
    const a = lo / o
    const b = hi / o
    return a < b ? [a, b] : [b, a]
  }
  const sx = slab(ox, -half, half)
  const sy = slab(oy, y0, y1)
  if (sx && sy) {
    const enter = Math.max(sx[0], sy[0], 0)
    const exit = Math.min(sx[1], sy[1])
    if (exit >= enter && exit > 0) t = Math.max(t, exit)
  }
  return t
}

const edges = computed<EdgeView[]>(() => {
  const g = data.value?.graph
  if (!g) return []
  const { W, H, r } = metrics.value
  const cxCanvas = W / 2
  const cyCanvas = H / 2
  const out: EdgeView[] = []
  g.edges.forEach((e, i) => {
    const a = position.value.get(e.from)
    const b = position.value.get(e.to)
    if (!a || !b) return
    const vx = b.x - a.x
    const vy = b.y - a.y
    const len = Math.hypot(vx, vy)
    if (len < 1) return
    const ux = vx / len
    const uy = vy / len
    const nx = -uy
    const ny = ux

    // 端点位置与弧线控制点
    let x1: number
    let y1: number
    let x2: number
    let y2: number
    let ctrlX: number
    let ctrlY: number
    let lnx: number
    let lny: number

    if (Math.abs(a.x - b.x) < 1) {
      // 同带（竖直线）：正下方被姓名牌与角色行占着，改成从圆盘侧面出入，
      // 再统一向「未来」一侧鼓出，与跨带弧线在视觉上分开。
      const side = r + 4
      x1 = a.x + side
      y1 = a.y
      x2 = b.x + side
      y2 = b.y
      const bulge = Math.min(72, Math.max(34, Math.abs(y2 - y1) * 0.32))
      ctrlX = Math.max(x1, x2) + bulge
      ctrlY = (y1 + y2) / 2
      lnx = 1
      lny = 0
    } else {
      // 两端各自按「朝外方向」绕开自己的圆盘与文字区，再按弦长等比收缩
      let padA = footprintExit(ux, uy, r, a.half) + 14
      let padB = footprintExit(-ux, -uy, r, b.half) + 14
      const budget = Math.max(len - 10, 2)
      if (padA + padB > budget) {
        const k = budget / (padA + padB)
        padA *= k
        padB *= k
      }
      x1 = a.x + ux * padA
      y1 = a.y + uy * padA
      x2 = b.x - ux * padB
      y2 = b.y - uy * padB
      const mx = (x1 + x2) / 2
      const my = (y1 + y2) / 2
      const chord = Math.hypot(x2 - x1, y2 - y1)
      const off = Math.min(64, 16 + chord * 0.1)
      // 控制点取「远离画布中心」的一侧，同区域的弧线自然呈扇形散开、不再穿过图心
      const d1 = (mx + nx * off - cxCanvas) ** 2 + (my + ny * off - cyCanvas) ** 2
      const d2 = (mx - nx * off - cxCanvas) ** 2 + (my - ny * off - cyCanvas) ** 2
      const sign = d1 >= d2 ? 1 : -1
      ctrlX = mx + nx * off * sign
      ctrlY = my + ny * off * sign
      lnx = nx
      lny = ny
    }

    out.push({
      i,
      from: e.from,
      to: e.to,
      kind: e.kind,
      label: e.label,
      note: e.note,
      style: kindStyle(e.kind),
      x1,
      y1,
      x2,
      y2,
      cx: ctrlX,
      cy: ctrlY,
      nx: lnx,
      ny: lny,
    })
  })
  return out
})

const edgePath = (e: EdgeView) =>
  `M${e.x1.toFixed(1)} ${e.y1.toFixed(1)} Q${e.cx.toFixed(1)} ${e.cy.toFixed(1)} ${e.x2.toFixed(1)} ${e.y2.toFixed(1)}`

// 聚焦时相关边转红加粗；方向语义保持不变，所以双向边才补起始箭头。
const edgeHot = (e: EdgeView) => !!active.value && focus.value.edges.has(e.i)
const edgeStroke = (e: EdgeView) => (edgeHot(e) ? C.red : e.style.stroke)
const edgeWidth = (e: EdgeView) => (edgeHot(e) ? e.style.width + 1.6 : e.style.width)
const edgeMarker = (e: EdgeView) => `url(#${MARKER[edgeStroke(e)]})`

// ── 聚焦与选中 ──────────────────────────────────────────────────
const active = ref<string | null>(null)
const activeEdge = ref<number | null>(null)

/** 聚焦集合：节点与边的下标。 */
const focus = computed(() => {
  const nodeSet = new Set<string>()
  const edgeSet = new Set<number>()
  if (active.value) {
    nodeSet.add(active.value)
    for (const e of edges.value) {
      if (e.from === active.value || e.to === active.value) {
        nodeSet.add(e.from)
        nodeSet.add(e.to)
        edgeSet.add(e.i)
      }
    }
  }
  return { nodes: nodeSet, edges: edgeSet }
})

const dimNode = (slug: string) => !!active.value && !focus.value.nodes.has(slug)
const dimEdge = (i: number) => !!active.value && !focus.value.edges.has(i)
const isFocusNode = (slug: string) => active.value === slug

const nameOf = (slug: string) => nodes.value.find((n) => n.slug === slug)?.name ?? slug
const personOf = (slug: string) => nodes.value.find((n) => n.slug === slug)

/** 聚焦态需要显示的边标签（含避让试位结果）。 */
const edgeLabels = computed(() => {
  if (!active.value) return []
  const { r } = metrics.value
  const verts = placedNodes.value.map((n) => ({ x: n.x, y: n.y, half: n.half }))
  const placed: { x: number; y: number; w: number; h: number }[] = []
  const qp = (e: EdgeView, t: number) => {
    const u = 1 - t
    return {
      x: u * u * e.x1 + 2 * u * t * e.cx + t * t * e.x2,
      y: u * u * e.y1 + 2 * u * t * e.cy + t * t * e.y2,
    }
  }
  const fits = (x: number, y: number, w: number, h: number) => {
    const bottom = r + PLATE_TOP + PLATE_H + ROLE_H + 4
    for (const v of verts) {
      const dx = Math.max(Math.abs(x - v.x) - w / 2, 0)
      const dy = Math.max(Math.abs(y - v.y) - h / 2, 0)
      if (Math.hypot(dx, dy) < r + 5) return false // 压住圆盘
      if (
        Math.abs(x - v.x) < w / 2 + v.half + 4 &&
        y - v.y > -r - 8 &&
        y - v.y < bottom
      ) {
        return false // 压住姓名牌或角色行
      }
    }
    for (const q of placed) {
      if (Math.abs(x - q.x) < (w + q.w) / 2 + 8 && Math.abs(y - q.y) < (h + q.h) / 2 + 6) return false
    }
    return true
  }

  const out: { i: number; x: number; y: number; w: number; h: number; text: string }[] = []
  for (const e of edges.value) {
    if (!focus.value.edges.has(e.i)) continue
    const text = e.kind
    const w = textWidth(text, 12) + 14
    const h = 22
    const cands: { x: number; y: number }[] = []
    const TS = [0.5, 0.45, 0.55, 0.4, 0.6, 0.35, 0.65, 0.3, 0.7, 0.25, 0.75, 0.2, 0.8, 0.12, 0.88]
    const OFFSETS = [0, 20, -20, 34, -34, 48, -48, 62, -62]
    for (const t of TS) {
      const p = qp(e, t)
      for (const k of OFFSETS) cands.push({ x: p.x + e.nx * k, y: p.y + e.ny * k })
    }
    const spot = cands.find((c) => fits(c.x, c.y, w, h)) ?? qp(e, 0.5)
    placed.push({ x: spot.x, y: spot.y, w, h })
    out.push({ i: e.i, x: spot.x, y: spot.y, w, h, text })
  }
  return out
})

// ── 关系列表（等价物）：按关系族分组 ─────────────────────────────
const groupedEdges = computed(() => {
  const list = active.value ? edges.value.filter((e) => focus.value.edges.has(e.i)) : edges.value
  const map = new Map<string, EdgeView[]>()
  for (const e of list) {
    if (!map.has(e.kind)) map.set(e.kind, [])
    map.get(e.kind)!.push(e)
  }
  return [...map.entries()]
    .sort((a, b) => kindRank(a[0]) - kindRank(b[0]) || b[1].length - a[1].length)
    .map(([kind, items]) => ({ kind, items, style: kindStyle(kind) }))
})

const selectedEdge = computed(() =>
  activeEdge.value === null ? null : (edges.value.find((e) => e.i === activeEdge.value) ?? null),
)

/** 聚焦人物的关系按族计数，用于详情面板。 */
const activeGroups = computed(() => {
  if (!active.value) return []
  const list = edges.value.filter((e) => focus.value.edges.has(e.i))
  const map = new Map<string, number>()
  for (const e of list) map.set(e.kind, (map.get(e.kind) ?? 0) + 1)
  return [...map.entries()]
    .sort((a, b) => kindRank(a[0]) - kindRank(b[0]) || b[1] - a[1])
    .map(([kind, count]) => ({ kind, count, style: kindStyle(kind) }))
})
const activeOut = computed(() =>
  active.value ? edges.value.filter((e) => focus.value.edges.has(e.i) && e.from === active.value).length : 0,
)
const activeNbrs = computed(() => {
  if (!active.value) return []
  const others = new Set<string>()
  for (const e of edges.value) {
    if (!focus.value.edges.has(e.i)) continue
    others.add(e.from === active.value ? e.to : e.from)
  }
  return placedNodes.value.filter((n) => others.has(n.slug)).sort((a, b) => b.degree - a.degree)
})

function clearAll() {
  active.value = null
  activeEdge.value = null
}
function toggleNode(slug: string) {
  active.value = active.value === slug ? null : slug
  activeEdge.value = null
}
function selectEdge(i: number) {
  activeEdge.value = activeEdge.value === i ? null : i
}

// ── 相机：viewBox 缩放 + 平移 ────────────────────────────────────
const graphBox = ref<HTMLElement | null>(null)
const containerWidth = ref(0)
const zoom = ref(1)
const center = ref<{ x: number; y: number } | null>(null)

let observer: ResizeObserver | undefined
watch(graphBox, (el) => {
  observer?.disconnect()
  if (!el) return
  containerWidth.value = Math.max(240, el.clientWidth)
  observer = new ResizeObserver((entries) => {
    const w = entries[0]?.contentRect.width ?? 0
    if (w > 0) containerWidth.value = Math.max(240, w)
  })
  observer.observe(el)
})
onBeforeUnmount(() => observer?.disconnect())

/** 窄屏：默认收起图形，改走下方的关系列表。 */
const narrow = computed(() => containerWidth.value > 0 && containerWidth.value < 640)
const showGraph = ref(false)
let visibilityDecided = false
watch(containerWidth, (w) => {
  if (visibilityDecided || w <= 0) return
  visibilityDecided = true
  showGraph.value = !narrow.value
})

/** 视图尺寸 = 画布尺寸 / 缩放。zoom = 1 即「全图可见」。 */
const viewSize = computed(() => {
  const { W, H } = metrics.value
  return { w: W / zoom.value, h: H / zoom.value }
})

const viewBox = computed(() => {
  const { W, H } = metrics.value
  const { w, h } = viewSize.value
  const cx = clamp(center.value?.x ?? W / 2, w / 2, Math.max(w / 2, W - w / 2))
  const cy = clamp(center.value?.y ?? H / 2, h / 2, Math.max(h / 2, H - h / 2))
  return `${(cx - w / 2).toFixed(1)} ${(cy - h / 2).toFixed(1)} ${w.toFixed(1)} ${h.toFixed(1)}`
})

/** 屏幕上的姓名牌有效字号（px）。低于下限时提示改用列表或放大。 */
const effectiveNamePx = computed(() =>
  containerWidth.value > 0 ? (NAME_FS * zoom.value * containerWidth.value) / metrics.value.W : NAME_FS,
)

/** 保守的下限缩放：屏幕有效字号不低于 NAME_FS_FLOOR。容器宽度未知时不设限。 */
const minZoom = computed(() => {
  if (containerWidth.value <= 0) return 1
  return clamp((NAME_FS_FLOOR * metrics.value.W) / (NAME_FS * containerWidth.value), 1, ZOOM_MAX)
})

/** 用户是否自己调过缩放；调过就不再用下限去覆盖他的选择。 */
const zoomedByUser = ref(false)

// 画布尺寸变化（数据量变了）→ 视野回中
watch([() => metrics.value.W, () => metrics.value.H], () => {
  center.value = { x: metrics.value.W / 2, y: metrics.value.H / 2 }
})

// 画布或容器变化时，若用户没动过缩放，就落到可读下限
watch([() => metrics.value.W, containerWidth], () => {
  if (!zoomedByUser.value) zoom.value = minZoom.value
}, { immediate: true })

watch(active, (slug) => {
  const p = slug ? position.value.get(slug) : null
  if (p) center.value = { x: p.x, y: p.y }
})

const zoomBy = (factor: number, anchor?: { x: number; y: number }) => {
  const { W, H } = metrics.value
  const before = { w: W / zoom.value, h: H / zoom.value }
  const c = center.value ?? { x: W / 2, y: H / 2 }
  const next = clamp(zoom.value * factor, 1, ZOOM_MAX)
  const after = { w: W / next, h: H / next }
  if (anchor) {
    // 保持光标下的点不动
    const bx = c.x - before.w / 2 + anchor.x * before.w
    const by = c.y - before.h / 2 + anchor.y * before.h
    center.value = { x: bx + (0.5 - anchor.x) * after.w, y: by + (0.5 - anchor.y) * after.h }
  }
  zoom.value = next
  zoomedByUser.value = true
}

const fitAll = () => {
  zoom.value = 1
  center.value = { x: metrics.value.W / 2, y: metrics.value.H / 2 }
  zoomedByUser.value = true
}

const zoomPercent = computed(() => Math.round(zoom.value * 100))

const drag = { on: false, x: 0, y: 0, el: null as HTMLElement | null }

const localAnchor = (clientX: number, clientY: number) => {
  const el = graphBox.value
  if (!el) return { x: 0.5, y: 0.5 }
  const rect = el.getBoundingClientRect()
  return { x: (clientX - rect.left) / rect.width, y: (clientY - rect.top) / rect.height }
}

/**
 * 滚轮缩放。缩放态或按住 Ctrl（触控板捏合）时才拦截默认行为；
 * 全图状态下让滚轮照常滚动页面，避免整页被图卡住。
 */
function onWheel(e: WheelEvent) {
  if (zoom.value <= 1 && !e.ctrlKey) return
  e.preventDefault()
  const anchor = localAnchor(e.clientX, e.clientY)
  zoomBy(Math.exp(-e.deltaY * 0.0015), anchor)
}

function onPointerDown(e: PointerEvent) {
  if (zoom.value <= 1) return
  const el = graphBox.value
  if (!el) return
  drag.on = true
  drag.x = e.clientX
  drag.y = e.clientY
  drag.el = el
  el.setPointerCapture(e.pointerId)
}

function onPointerMove(e: PointerEvent) {
  if (!drag.on) return
  const rect = graphBox.value?.getBoundingClientRect()
  if (!rect) return
  const { w, h } = viewSize.value
  const c = center.value ?? { x: metrics.value.W / 2, y: metrics.value.H / 2 }
  center.value = {
    x: c.x - ((e.clientX - drag.x) * w) / rect.width,
    y: c.y - ((e.clientY - drag.y) * h) / rect.height,
  }
  drag.x = e.clientX
  drag.y = e.clientY
}

function onPointerUp(e: PointerEvent) {
  if (!drag.on) return
  drag.on = false
  drag.el?.releasePointerCapture?.(e.pointerId)
  drag.el = null
}
</script>

<template>
  <div>
    <PageHeader
      no="05"
      title="人物关系图"
      sub="Relations"
      lead="按出生世代分栏的思想传承谱系：横轴是世代，左为源、右为流。点击人物可聚焦其全部关系，点击连线可读该条关系的说明；下方同时提供按关系族分组、无需图形的等价列表，便于键盘与读屏访问。"
    />

    <section class="mx-auto max-w-6xl px-4 py-10 sm:px-6">
      <AsyncBoundary :loading="loading" :error="error" :on-retry="reload">
        <div ref="graphBox">
          <!-- 工具条 -->
          <div class="mb-3 flex flex-wrap items-center gap-2">
            <span class="stamp stamp-ink">世代分栏</span>
            <span class="stamp">{{ nodes.length }} 人 · {{ edges.length }} 条关系</span>
            <span class="stamp">{{ metrics.bandCount }} 个世代</span>
            <div class="ml-auto flex items-center gap-2">
              <button
                type="button"
                class="border-2 border-ink bg-paper px-2.5 py-1 font-mono text-xs font-bold hover:bg-paper-dim disabled:opacity-40"
                :disabled="zoom <= 1"
                aria-label="缩小"
                @click="zoomBy(1 / 1.3)"
              >
                −
              </button>
              <span class="font-mono text-xs font-bold tabular-nums">{{ zoomPercent }}%</span>
              <button
                type="button"
                class="border-2 border-ink bg-paper px-2.5 py-1 font-mono text-xs font-bold hover:bg-paper-dim disabled:opacity-40"
                :disabled="zoom >= ZOOM_MAX"
                aria-label="放大"
                @click="zoomBy(1.3)"
              >
                ＋
              </button>
              <button
                type="button"
                class="border-2 border-ink bg-paper px-2.5 py-1 font-mono text-xs font-bold hover:bg-paper-dim"
                @click="fitAll()"
              >
                全图
              </button>
              <button
                type="button"
                class="border-2 border-ink bg-paper px-2.5 py-1 font-mono text-xs font-bold hover:bg-paper-dim sm:hidden"
                @click="showGraph = !showGraph"
              >
                {{ showGraph ? '收起图形' : '显示图形' }}
              </button>
            </div>
          </div>

          <!-- 窄屏提示：图形默认收起 -->
          <p v-if="narrow && !showGraph" class="brutal-card border-l-[10px] border-l-steel p-5">
            <span class="font-bold">屏幕较窄，关系图默认收起。</span>
            下方的关系列表提供完全等价的信息，可按关系族逐条阅读；也可以点右上角「显示图形」放大查看。
          </p>

          <div v-show="showGraph" class="brutal-card p-2">
            <p
              v-if="effectiveNamePx < 11"
              class="mb-2 border-l-[6px] border-brass bg-paper-dim/70 px-3 py-2 font-mono text-xs font-bold text-ink-soft"
            >
              当前宽度下姓名只有约 {{ effectiveNamePx.toFixed(1) }}px，点「＋」放大后拖动查看，或改用下方关系列表。
            </p>
            <div
              class="touch-none select-none"
              :class="zoom > 1 ? 'cursor-grab active:cursor-grabbing' : ''"
              @wheel="onWheel"
              @pointerdown="onPointerDown"
              @pointermove="onPointerMove"
              @pointerup="onPointerUp"
              @pointercancel="onPointerUp"
              @dblclick="fitAll()"
              @keydown.esc="clearAll()"
            >
              <svg
                :viewBox="viewBox"
                class="mx-auto block h-auto w-full"
                role="img"
                aria-label="人物关系图谱：按出生世代分栏的思想传承谱系"
              >
                <defs>
                  <filter :id="MONO_FILTER" x="-10%" y="-10%" width="120%" height="120%">
                    <feColorMatrix type="saturate" values="0" />
                    <feComponentTransfer>
                      <feFuncR type="linear" slope="1.06" intercept="0.06" />
                      <feFuncG type="linear" slope="1.02" intercept="0.03" />
                      <feFuncB type="linear" slope="0.92" />
                    </feComponentTransfer>
                  </filter>
                  <marker
                    v-for="m in MARKS"
                    :id="m.id"
                    :key="m.id"
                    viewBox="0 0 10 8"
                    refX="9"
                    refY="4"
                    markerWidth="9"
                    markerHeight="7"
                    orient="auto-start-reverse"
                    markerUnits="userSpaceOnUse"
                  >
                    <path d="M0,0 L10,4 L0,8 Z" :fill="m.color" />
                  </marker>
                  <clipPath v-for="n in placedNodes" :id="`gclip-${n.slug}`" :key="n.slug">
                    <circle :cx="n.x" :cy="n.y" :r="metrics.r - 3" />
                  </clipPath>
                </defs>

                <!-- 底：色纸 + 横线（按画布高度动态生成，不再写死条数） -->
                <rect :width="metrics.W" :height="metrics.H" :fill="C.paper" @click="clearAll()" />
                <g opacity="0.05" pointer-events="none">
                  <path
                    v-for="y in Math.floor(metrics.H / 46) + 1"
                    :key="`rule${y}`"
                    :d="`M0 ${(y - 1) * 46} L${metrics.W} ${(y - 1) * 46}`"
                    :stroke="C.ink"
                    stroke-width="2"
                  />
                </g>

                <!-- 世代带 -->
                <g pointer-events="none">
                  <template v-for="(meta, bi) in bandMeta" :key="`band${bi}`">
                    <rect
                      :x="MARGIN_X + bi * metrics.hPitch - metrics.hPitch / 2 + 16"
                      :y="BAND_TOP"
                      :width="metrics.hPitch - 32"
                      :height="metrics.H - BAND_TOP * 2"
                      :fill="C.paperDim"
                      fill-opacity="0.55"
                      :stroke="C.ink"
                      stroke-width="2"
                      stroke-dasharray="7 6"
                      opacity="0.85"
                    />
                    <rect
                      :x="MARGIN_X + bi * metrics.hPitch - metrics.hPitch / 2 + 16"
                      :y="BAND_TOP - 34"
                      :width="metrics.hPitch - 32"
                      height="26"
                      :fill="C.ink"
                    />
                    <text
                      :x="MARGIN_X + bi * metrics.hPitch"
                      :y="BAND_TOP - 16"
                      text-anchor="middle"
                      font-size="12"
                      font-weight="700"
                      letter-spacing="0.12em"
                      :fill="C.paper"
                      class="font-mono"
                    >
                      {{ meta.label }}
                    </text>
                    <text
                      :x="MARGIN_X + bi * metrics.hPitch"
                      :y="metrics.H - BAND_TOP + 22"
                      text-anchor="middle"
                      font-size="11"
                      font-weight="700"
                      :fill="C.steel"
                      class="font-mono"
                    >
                      第 {{ bi + 1 }} 世代 · {{ meta.count }} 人
                    </text>
                  </template>
                </g>

                <!-- 连线 -->
                <g fill="none">
                  <template v-for="e in edges" :key="`e${e.i}`">
                    <path
                      :d="edgePath(e)"
                      :stroke="edgeStroke(e)"
                      :stroke-width="edgeWidth(e).toFixed(1)"
                      :stroke-dasharray="e.style.dash"
                      stroke-linecap="round"
                      :marker-end="edgeMarker(e)"
                      :marker-start="e.style.both ? edgeMarker(e) : undefined"
                      :opacity="dimEdge(e.i) ? 0.16 : 1"
                      pointer-events="none"
                    />
                    <!-- 细线也点得中的热区 -->
                    <path
                      :d="edgePath(e)"
                      stroke="transparent"
                      stroke-width="14"
                      pointer-events="stroke"
                      class="cursor-pointer"
                      @click.stop="selectEdge(e.i)"
                    >
                      <title>{{ nameOf(e.from) }} — {{ nameOf(e.to) }}（{{ e.kind }}）：{{ e.label }}</title>
                    </path>
                  </template>
                </g>

                <!-- 聚焦态边标签（已在脚本里做过避让试位） -->
                <g pointer-events="none">
                  <g v-for="l in edgeLabels" :key="`lb${l.i}`" class="edge-label">
                    <rect
                      :x="l.x - l.w / 2"
                      :y="l.y - l.h / 2"
                      :width="l.w"
                      :height="l.h"
                      :fill="C.paper"
                      :stroke="C.ink"
                      stroke-width="2"
                    />
                    <text
                      :x="l.x"
                      :y="l.y + 4"
                      text-anchor="middle"
                      font-size="12"
                      font-weight="700"
                      :fill="C.ink"
                      class="font-mono"
                    >
                      {{ l.text }}
                    </text>
                  </g>
                </g>

                <!-- 节点 -->
                <g
                  v-for="n in placedNodes"
                  :key="n.slug"
                  role="button"
                  tabindex="0"
                  class="graph-node cursor-pointer"
                  :opacity="dimNode(n.slug) ? 0.18 : 1"
                  :aria-pressed="isFocusNode(n.slug)"
                  :aria-label="`${n.name}${n.nationality ? `，${n.nationality}` : ''}${n.epithet ? `，${n.epithet}` : ''}，关系 ${n.degree} 条`"
                  @click.stop="toggleNode(n.slug)"
                  @keydown.enter.prevent="toggleNode(n.slug)"
                  @keydown.space.prevent="toggleNode(n.slug)"
                >
                  <!-- 4px 硬阴影，呼应 brutal-card -->
                  <circle :cx="n.x + 4" :cy="n.y + 4" :r="metrics.r" :fill="C.ink" />
                  <circle :cx="n.x" :cy="n.y" :r="metrics.r" :fill="C.paperDim" />
                  <!-- 缺图时露出的姓氏首字 -->
                  <text
                    :x="n.x"
                    :y="n.y + metrics.r * 0.34"
                    text-anchor="middle"
                    :font-size="metrics.r * 0.9"
                    font-weight="900"
                    :fill="C.steel"
                  >
                    {{ n.short.charAt(0) }}
                  </text>
                  <image
                    v-if="n.portrait"
                    :x="n.x - metrics.r"
                    :y="n.y - metrics.r"
                    :width="metrics.r * 2"
                    :height="metrics.r * 2"
                    :href="n.portrait"
                    :clip-path="`url(#gclip-${n.slug})`"
                    preserveAspectRatio="xMidYMid slice"
                    :filter="`url(#${MONO_FILTER})`"
                  />
                  <circle
                    :cx="n.x"
                    :cy="n.y"
                    :r="metrics.r"
                    fill="none"
                    :stroke="isFocusNode(n.slug) ? C.red : C.ink"
                    :stroke-width="isFocusNode(n.slug) ? 4 : 3"
                  />
                  <circle
                    v-if="isFocusNode(n.slug)"
                    :cx="n.x"
                    :cy="n.y"
                    :r="metrics.r + 7"
                    fill="none"
                    :stroke="C.red"
                    stroke-width="2"
                    stroke-dasharray="4 4"
                  />
                  <!-- 键盘焦点环：SVG 元素上 outline 不可靠，改用一层显隐的描边 -->
                  <circle
                    class="graph-focus-ring"
                    :cx="n.x"
                    :cy="n.y"
                    :r="metrics.r + 11"
                    fill="none"
                    :stroke="C.red"
                    stroke-width="3"
                    pointer-events="none"
                  />
                  <!-- 度数角标 -->
                  <rect
                    :x="n.x + metrics.r - 9"
                    :y="n.y - metrics.r - 16"
                    width="27"
                    height="23"
                    :fill="isFocusNode(n.slug) ? C.red : C.ink"
                  />
                  <text
                    :x="n.x + metrics.r + 4.5"
                    :y="n.y - metrics.r + 0.5"
                    text-anchor="middle"
                    :font-size="DEG_FS"
                    font-weight="900"
                    :fill="C.paper"
                    class="font-mono"
                  >
                    {{ n.degree }}
                  </text>
                  <!-- 姓名牌 + 角色行 -->
                  <rect
                    :x="n.x - (textWidth(n.short, NAME_FS) + 18) / 2"
                    :y="n.y + metrics.r + 4"
                    :width="textWidth(n.short, NAME_FS) + 18"
                    height="25"
                    :fill="isFocusNode(n.slug) ? C.red : C.paper"
                    :stroke="C.ink"
                    stroke-width="2"
                  />
                  <text
                    :x="n.x"
                    :y="n.y + metrics.r + 22"
                    text-anchor="middle"
                    :font-size="NAME_FS"
                    font-weight="900"
                    :fill="isFocusNode(n.slug) ? C.paper : C.ink"
                  >
                    {{ n.short }}
                  </text>
                  <text
                    :x="n.x"
                    :y="n.y + metrics.r + 42"
                    text-anchor="middle"
                    :font-size="ROLE_FS"
                    font-weight="700"
                    :fill="C.steel"
                    :stroke="C.paper"
                    stroke-width="3.5"
                    paint-order="stroke"
                    class="font-mono"
                  >
                    {{ n.nationality }}<template v-if="n.nationality && n.year"> · </template>{{ n.year || '' }}
                  </text>
                  <title>{{ n.name }}｜{{ n.epithet }}｜关系 {{ n.degree }} 条</title>
                </g>
              </svg>
            </div>
          </div>

          <!-- 图例 -->
          <h2 class="mt-8 mb-3 border-b-3 border-ink pb-2 text-xl font-black">图例</h2>
          <ul class="grid gap-2 sm:grid-cols-2">
            <li
              v-for="lg in LEGEND"
              :key="lg.label"
              class="flex items-start gap-3 border-2 border-ink bg-paper p-3"
            >
              <svg width="58" height="14" viewBox="0 0 58 14" class="mt-1 shrink-0" aria-hidden="true">
                <path
                  d="M6 7 L52 7"
                  :stroke="lg.stroke"
                  :stroke-width="lg.width"
                  :stroke-dasharray="lg.dash"
                />
                <path :d="lg.both ? 'M14 7 L14 3 L20 7 L14 11 Z' : 'M44 2 L54 7 L44 12 Z'" :fill="lg.stroke" />
                <path v-if="lg.both" d="M44 7 L44 3 L38 7 L44 11 Z" :fill="lg.stroke" />
              </svg>
              <span>
                <span class="block text-sm font-black">{{ lg.label }}</span>
                <span class="block text-xs leading-relaxed text-ink-soft">{{ lg.desc }}</span>
              </span>
            </li>
          </ul>

          <!-- 聚焦详情 -->
          <h2 class="mt-10 mb-4 border-b-3 border-ink pb-2 text-2xl font-black">
            {{ active ? '聚焦详情' : '关系说明' }}
          </h2>

          <article v-if="active" class="brutal-card">
            <div class="flex flex-wrap items-center gap-2 border-b-3 border-ink bg-red px-4 py-2.5 text-paper">
              <h3 class="text-xl font-black">{{ nameOf(active) }}</h3>
              <button
                type="button"
                class="ml-auto border-2 border-paper px-2 py-0.5 font-mono text-xs font-bold"
                @click="clearAll()"
              >
                清除聚焦
              </button>
            </div>
            <div class="p-4">
              <p class="font-mono text-xs text-steel">
                {{ personOf(active)?.nationality }}
                <template v-if="personOf(active)?.life"> · {{ personOf(active)?.life }}</template>
              </p>
              <p v-if="personOf(active)?.epithet" class="mt-1 text-sm font-bold text-red-deep">
                {{ personOf(active)?.epithet }}
              </p>
              <div class="mt-3 flex flex-wrap gap-1.5">
                <span class="stamp stamp-ink">关系 {{ focus.edges.size }} 条</span>
                <span class="stamp">出 {{ activeOut }}</span>
                <span class="stamp">入 {{ focus.edges.size - activeOut }}</span>
                <span v-for="g in activeGroups" :key="g.kind" class="stamp">{{ g.kind }} {{ g.count }}</span>
              </div>
              <ul class="mt-4 flex flex-wrap gap-1.5">
                <li v-for="n in activeNbrs" :key="n.slug">
                  <button
                    type="button"
                    class="border-2 border-ink bg-paper px-2 py-0.5 text-xs font-bold hover:bg-paper-dim"
                    @click="toggleNode(n.slug)"
                  >
                    {{ n.short }}
                  </button>
                </li>
              </ul>
            </div>
          </article>

          <article v-else-if="selectedEdge" class="brutal-card border-l-[10px] border-l-red p-4">
            <div class="flex flex-wrap items-center gap-2">
              <RouterLink :to="`/people/${selectedEdge.from}`" class="font-bold hover:text-red-deep">
                {{ nameOf(selectedEdge.from) }}
              </RouterLink>
              <span class="font-mono text-xs text-red">
                {{ kindStyle(selectedEdge.kind).both ? '↔' : '→' }}
              </span>
              <RouterLink :to="`/people/${selectedEdge.to}`" class="font-bold hover:text-red-deep">
                {{ nameOf(selectedEdge.to) }}
              </RouterLink>
              <span class="stamp ml-auto">{{ selectedEdge.kind }}</span>
              <button
                type="button"
                class="font-mono text-xs font-bold underline"
                @click="activeEdge = null"
              >
                关闭
              </button>
            </div>
            <p class="mt-2 text-sm font-bold">{{ selectedEdge.label }}</p>
            <p class="mt-1.5 text-xs leading-relaxed text-ink-soft">{{ selectedEdge.note }}</p>
            <p class="mt-2 font-mono text-[0.68rem] text-steel">
              {{ selectedEdge.style.family }}族 · {{ selectedEdge.style.note }}
            </p>
          </article>

          <p v-else class="border-l-[6px] border-steel bg-paper-dim/60 px-4 py-3 text-sm text-ink-soft">
            点击图中的人物聚焦其关系网，或点击任意一条连线查看这条关系是什么。当前显示全部
            {{ edges.length }} 条关系。
          </p>

          <!-- 等价列表：按关系族分组 -->
          <h2 class="mt-10 mb-4 border-b-3 border-ink pb-2 text-2xl font-black">
            关系列表
            <span v-if="active" class="font-mono text-sm font-bold text-steel">
              —— 仅 {{ nameOf(active) }}
            </span>
          </h2>
          <div class="grid gap-4 md:grid-cols-2">
            <article
              v-for="g in groupedEdges"
              :key="g.kind"
              class="brutal-card border-l-[10px] border-l-red p-4"
            >
              <div class="mb-3 flex flex-wrap items-center gap-2 border-b-2 border-ink pb-2">
                <h3 class="text-base font-black">{{ g.kind }}</h3>
                <span class="stamp">{{ g.style.family }}族</span>
                <span class="ml-auto font-mono text-xs text-steel">{{ g.items.length }} 条</span>
              </div>
              <ul class="grid gap-3">
                <li v-for="e in g.items" :key="`l${e.i}`">
                  <div class="flex flex-wrap items-center gap-2">
                    <RouterLink :to="`/people/${e.from}`" class="text-sm font-bold hover:text-red-deep">
                      {{ nameOf(e.from) }}
                    </RouterLink>
                    <span class="font-mono text-xs text-red">
                      {{ g.style.both ? '↔' : '→' }}
                    </span>
                    <RouterLink :to="`/people/${e.to}`" class="text-sm font-bold hover:text-red-deep">
                      {{ nameOf(e.to) }}
                    </RouterLink>
                  </div>
                  <p class="text-sm font-bold">{{ e.label }}</p>
                  <p class="mt-1 text-xs leading-relaxed text-ink-soft">{{ e.note }}</p>
                </li>
              </ul>
            </article>
          </div>
        </div>
      </AsyncBoundary>
    </section>
  </div>
</template>

<style scoped>
/* 节点键盘焦点环：SVG 元素上的 outline 在各浏览器表现不一，改用显隐描边。 */
.graph-focus-ring {
  opacity: 0;
}
.graph-node:focus-visible .graph-focus-ring {
  opacity: 1;
}
</style>
