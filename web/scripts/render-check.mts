// 用 Chrome DevTools Protocol 渲染页面，验证 Vue 真的挂载并捕获截图。
import { launchCdp, savePng, sleep } from "./lib/cdp.mts"

const BASE = process.env.API_BASE ?? "http://127.0.0.1:12026"
const PORT = 9333
const OUT = "shots"

const cdp = await launchCdp({ port: PORT, windowSize: { width: 1440, height: 1000 } })

// 收集控制台错误
const errors: string[] = []
cdp.on((m) => {
  if (m.method === "Runtime.exceptionThrown") {
    const details = m.params?.exceptionDetails as { exception?: { description?: string } } | undefined
    errors.push((details?.exception?.description ?? "unknown").slice(0, 200))
  }
  if (m.method === "Runtime.consoleAPICalled" && m.params?.type === "error") {
    const args = (m.params?.args ?? []) as { value?: unknown; description?: string }[]
    errors.push(args.map((a) => a.value ?? a.description ?? "").join(" ").slice(0, 200))
  }
})

const ROUTES: [string, string][] = [
  ["/", "home"],
  ["/people", "people"],
  ["/people/mao", "person-mao"],
  ["/works", "works"],
  ["/works/imperialism", "work-imperialism"],
  ["/timeline", "timeline"],
  ["/glossary", "glossary"],
  ["/glossary/proletariat", "term-proletariat"],
  ["/graph", "graph"],
  ["/search?q=" + encodeURIComponent("无产阶级"), "search"],
  ["/about", "about"],
]

/**
 * 关系图专项断言：世代分栏布局的几何不变量。
 *
 * 这些不变量是「改版前那版环形布局」真正缺的东西——节点重叠、边标签在圆心叠成一团，
 * 都不是靠肉眼发现，而是靠这里量出来的。回归背景见 README 技术决策记录。
 */
async function checkGraph(sid: string): Promise<string[]> {
  const raw = (await cdp.evaluate(
    sid,
    `(async () => {
      const fail = []
      const svg = document.querySelector('svg[role=img]')
      if (!svg) return ['未渲染关系图 SVG']

      const graph = await (await fetch('/api/relations')).json()
      const groups = [...svg.querySelectorAll('g.graph-node')]
      const nodes = groups.map((g) => {
        const c = g.querySelector('circle')
        const plate = g.querySelector('rect[height="25"]')
        const texts = [...g.querySelectorAll('text')]
        const plateText = texts.find((t) => +t.getAttribute('font-size') === 15)
        const role = texts.find((t) => t.getAttribute('paint-order') === 'stroke')
        const label = g.getAttribute('aria-label') || ''
        return {
          label,
          full: label.split('，')[0],
          plate: plateText ? plateText.textContent.trim() : '',
          plateFont: plateText ? +plateText.getAttribute('font-size') : 0,
          roleBox: role ? role.getBBox() : null,
          plateBox: {
            x: +plate.getAttribute('x'), y: +plate.getAttribute('y'), w: +plate.getAttribute('width'), h: +plate.getAttribute('height'),
          },
          x: +c.getAttribute('cx'), y: +c.getAttribute('cy'), r: +c.getAttribute('r'),
        }
      })

      // ① 节点数 = 接口返回的人物数
      if (nodes.length !== graph.nodes.length) {
        fail.push(\`节点数 \${nodes.length} ≠ 接口人物数 \${graph.nodes.length}\`)
      }

      // ② 姓名牌必须等于「最后一个 · 之后」的显示名（旧版 slice(-3) 会截出「斯特罗」）
      for (const n of nodes) {
        const want = n.full.split('·').pop().trim()
        if (n.plate !== want) fail.push(\`姓名牌「\${n.plate}」应为「\${want}」\`)
      }

      // ③ 圆盘两两不重叠：中心距 ≥ 直径 + 8
      let minDist = Infinity
      for (let i = 0; i < nodes.length; i++) {
        for (let j = i + 1; j < nodes.length; j++) {
          minDist = Math.min(minDist, Math.hypot(nodes[i].x - nodes[j].x, nodes[i].y - nodes[j].y))
        }
      }
      if (nodes.length > 1 && minDist < nodes[0].r * 2 + 8) {
        fail.push(\`节点净空不足：最小中心距 \${minDist.toFixed(1)} < \${nodes[0].r * 2 + 8}\`)
      }

      // ④ 姓名牌的有效字号不低于 13px
      const vb = svg.viewBox.baseVal
      const scale = svg.getBoundingClientRect().width / (vb.width || 1)
      const effFont = nodes[0] ? nodes[0].plateFont * scale : 0
      if (effFont < 13) fail.push(\`姓名牌有效字号仅 \${effFont.toFixed(1)}px\`)

      // ⑤ 连线端点不得落在任何节点的姓名牌／角色行矩形内
      const boxes = nodes.map((n) => {
        const b = n.roleBox
        return {
          x0: Math.min(n.plateBox.x, b ? b.x : n.plateBox.x),
          x1: Math.max(n.plateBox.x + n.plateBox.w, b ? b.x + b.width : n.plateBox.x + n.plateBox.w),
          y0: n.plateBox.y,
          y1: b ? b.y + b.height : n.plateBox.y + n.plateBox.h,
        }
      })
      const hitPaths = [...svg.querySelectorAll('g[fill="none"] > path[stroke="transparent"]')]
      if (hitPaths.length !== graph.edges.length) {
        fail.push(\`连线数 \${hitPaths.length} ≠ 接口关系数 \${graph.edges.length}\`)
      }
      let endpointHits = 0
      for (const p of hitPaths) {
        const total = p.getTotalLength()
        for (const t of [0, total]) {
          const pt = p.getPointAtLength(t)
          if (boxes.some((b) => pt.x > b.x0 - 2 && pt.x < b.x1 + 2 && pt.y > b.y0 - 2 && pt.y < b.y1 + 2)) {
            endpointHits++
          }
        }
      }
      if (endpointHits > 0) fail.push(\`\${endpointHits} 个连线端点压在姓名牌／角色行上\`)

      // ⑥ 聚焦态：边标签不与任何节点圆盘相交、彼此不重叠
      const top = nodes.reduce((best, n, i) =>
        (n.label.match(/关系 (\\d+) 条/) || [0, 0])[1] > (nodes[best]?.label.match(/关系 (\\d+) 条/) || [0, 0])[1] ? i : best, 0)
      groups[top]?.dispatchEvent(new MouseEvent('click', { bubbles: true }))
      await new Promise((r) => setTimeout(r, 500))
      const labels = [...svg.querySelectorAll('g.edge-label rect')].map((r) => ({
        x: +r.getAttribute('x'), y: +r.getAttribute('y'), w: +r.getAttribute('width'), h: +r.getAttribute('height'),
      }))
      if (labels.length === 0) fail.push('聚焦态未渲染任何边标签')
      let clips = 0
      for (const l of labels) {
        for (const n of nodes) {
          const cx = l.x + l.w / 2, cy = l.y + l.h / 2
          const dx = Math.max(Math.abs(cx - n.x) - l.w / 2, 0)
          const dy = Math.max(Math.abs(cy - n.y) - l.h / 2, 0)
          if (Math.hypot(dx, dy) < n.r) clips++
        }
      }
      if (clips > 0) fail.push(\`\${clips} 个边标签压住节点圆盘\`)
      for (let i = 0; i < labels.length; i++) {
        for (let j = i + 1; j < labels.length; j++) {
          const a = labels[i], b = labels[j]
          if (a.x < b.x + b.w && b.x < a.x + a.w && a.y < b.y + b.h && b.y < a.y + a.h) {
            fail.push('边标签互相重叠')
            i = labels.length
            break
          }
        }
      }
      return fail
    })()`,
  )) as string[]

  return Array.isArray(raw) ? raw : ["关系图断言返回异常"]
}

let fail = 0
let pass = 0

for (const [route, name] of ROUTES) {
  errors.length = 0
  const sid = await cdp.newTab(BASE + route)
  await sleep(1400)

  // 挂载检查：Vue 渲染后 #app 必须有实际内容
  const info = (await cdp.evaluate(
    sid,
    `(() => ({
      appLen: (document.querySelector("#app")?.innerHTML ?? "").length,
      h1: document.querySelector("h1")?.textContent?.trim() ?? "",
      textLen: (document.body.innerText ?? "").length,
      nav: document.querySelectorAll("header nav a").length,
      cards: document.querySelectorAll(".brutal-card").length,
      wall: document.querySelectorAll("[data-hero-wall] [data-portrait]").length,
    }))()`,
  )) as { appLen: number; h1: string; textLen: number; cards: number; wall: number }

  const mounted = info.appLen > 400 && info.textLen > 200 && info.h1.length > 0
  const noErr = errors.length === 0
  // 首页门楣照片墙必须恰好填满 3×2 六格。
  // 回归背景：照片上限曾写死为 5，第 6 格被高饱和色块占据，五张照片反而成了陪衬。
  const wallOk = route !== "/" || info.wall === 6

  let graphIssues: string[] = []
  if (route === "/graph" && mounted) {
    graphIssues = await checkGraph(sid)
  }

  if (mounted && noErr && wallOk && graphIssues.length === 0) {
    pass++
    console.log(
      `  [PASS] ${route}  h1="${info.h1}"  卡片=${info.cards}  文本=${info.textLen}` +
        (route === "/" ? `  门楣墙=${info.wall}` : "") +
        (route === "/graph" ? "  关系图几何不变量通过" : ""),
    )
  } else {
    fail++
    console.log(`  [FAIL] ${route}`)
    console.log(`         mounted=${mounted} appLen=${info.appLen} textLen=${info.textLen} h1="${info.h1}"`)
    if (!wallOk) console.log(`         首页门楣照片墙应为 6 格，实际 ${info.wall}`)
    for (const g of graphIssues) console.log(`         关系图：${g}`)
    if (errors.length) console.log("         控制台错误: " + errors.slice(0, 3).join(" | "))
  }

  savePng(OUT, `${name}.png`, await cdp.screenshot(sid))
  await cdp.closeTab(sid)
}

console.log("")
console.log(fail === 0 ? `渲染全部通过: ${pass} 个路由` : `通过 ${pass} 失败 ${fail}`)
cdp.dispose()
process.exit(fail === 0 ? 0 : 1)
