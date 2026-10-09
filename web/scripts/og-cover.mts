// 渲染一张 1200×630 的社交分享封面（og:image）。
// 用无头 Chrome 打开本地 HTML（引用 /public 下的真实肖像），截图输出到 public/og-cover.png。
import { launchCdp, sleep, writeTempHtml } from "./lib/cdp.mts"

const PORT = 9336
const W = 1200
const H = 630

const people: [string, string][] = [
  ["marx", "卡尔·马克思"],
  ["engels", "弗里德里希·恩格斯"],
  ["lenin", "列宁"],
  ["mao", "毛泽东"],
  ["guevara", "切·格瓦拉"],
]

const figs = people
  .map(([k, name]) => `<figure><img src="../public/assets/portraits/${k}.webp" alt="${name}"></figure>`)
  .join("")

const html = `<!doctype html><html lang="zh-CN"><head><meta charset="utf-8"><style>
  *{box-sizing:border-box;margin:0;padding:0}
  html,body{width:${W}px;height:${H}px;overflow:hidden}
  body{background:#f2ede4;color:#111;display:flex;flex-direction:column;
    font-family:'Noto Sans SC','Source Han Sans SC','Microsoft YaHei',sans-serif}
  header{position:relative;padding:46px 48px 0;flex:0 0 auto}
  header::after{content:"";position:absolute;top:0;right:0;width:220px;height:64px;background:#d62828}
  .chips{display:flex;gap:8px;margin-bottom:20px}
  .chip{border:2px solid #111;background:#f2ede4;font-family:Consolas,monospace;font-size:13px;
    font-weight:700;padding:4px 9px;letter-spacing:.08em}
  .chip.red{background:#d62828;color:#f2ede4}
  h1{font-size:90px;font-weight:900;line-height:.9;letter-spacing:-.03em}
  h1 .r{color:#d62828}
  .en{margin-top:10px;font-family:Consolas,monospace;font-size:15px;letter-spacing:.34em;color:#4a4a4a}
  .sub{margin-top:20px;border-left:8px solid #d62828;padding-left:16px;font-size:20px;font-weight:700}
  .desc{margin-top:10px;font-size:15px;color:#1f1f1f;line-height:1.7;max-width:640px}
  .strip{margin-top:auto;padding:26px 48px 38px;display:grid;grid-template-columns:repeat(5,1fr);gap:14px}
  figure{border:3px solid #111;background:#e3dbcd;box-shadow:6px 6px 0 #111;overflow:hidden;aspect-ratio:1}
  figure img{width:100%;height:100%;object-fit:cover;display:block;
    filter:grayscale(1) contrast(1.06) sepia(.18)}
</style></head><body>
  <header>
    <div class="chips"><span class="chip red">NO. 12026</span><span class="chip">EST. 2026</span><span class="chip">中文</span></div>
    <h1>星火<span class="r">档案</span>馆</h1>
    <div class="en">SPARK ARCHIVE</div>
    <div class="sub">全世界无产者，联合起来！</div>
    <div class="desc">无产阶级革命理论与实践的文献档案 —— 收录马克思、恩格斯、列宁、毛泽东、切·格瓦拉五个人的生平、著作与思想。</div>
  </header>
  <div class="strip">${figs}</div>
</body></html>`

const url = await writeTempHtml("shots", "og.html", html)
const cdp = await launchCdp({ port: PORT, windowSize: { width: W, height: H } })
const sid = await cdp.newTab(url)
await cdp.send(
  "Emulation.setDeviceMetricsOverride",
  { width: W, height: H, deviceScaleFactor: 1, mobile: false },
  sid,
)
await sleep(1500)

const { mkdirSync, writeFileSync } = await import("node:fs")
mkdirSync("public", { recursive: true })
writeFileSync("public/og-cover.png", await cdp.screenshot(sid))
console.log("已生成 public/og-cover.png")
cdp.dispose()
