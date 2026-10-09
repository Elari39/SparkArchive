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
    const args = (m.params.args ?? []) as { value?: unknown; description?: string }[]
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
    }))()`,
  )) as { appLen: number; h1: string; textLen: number; cards: number }

  const mounted = info.appLen > 400 && info.textLen > 200 && info.h1.length > 0
  const noErr = errors.length === 0

  if (mounted && noErr) {
    pass++
    console.log(`  [PASS] ${route}  h1="${info.h1}"  卡片=${info.cards}  文本=${info.textLen}`)
  } else {
    fail++
    console.log(`  [FAIL] ${route}`)
    console.log(`         mounted=${mounted} appLen=${info.appLen} textLen=${info.textLen} h1="${info.h1}"`)
    if (errors.length) console.log("         控制台错误: " + errors.slice(0, 3).join(" | "))
  }

  savePng(OUT, `${name}.png`, await cdp.screenshot(sid))
  await cdp.closeTab(sid)
}

console.log("")
console.log(fail === 0 ? `渲染全部通过: ${pass} 个路由` : `通过 ${pass} 失败 ${fail}`)
cdp.dispose()
process.exit(fail === 0 ? 0 : 1)
