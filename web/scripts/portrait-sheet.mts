// 把全部人物肖像（真实照片）拼成对比图，核对裁切与单色处理是否统一。
// 直接读取 public 下的资产（无需站点运行），用无头 Chrome 截图输出 shots/portraits.png。
// 15 位人物现已全部配图，故逐个列出；缺图／加载失败时前端仍会回退为矢量插画，但不在本图内。
import { launchCdp, savePng, sleep, writeTempHtml } from "./lib/cdp.mts"

const PORT = 9335
const W = 1420
const H = 470

const people: [string, string][] = [
  ["marx", "马克思"],
  ["engels", "恩格斯"],
  ["plekhanov", "普列汉诺夫"],
  ["zetkin", "蔡特金"],
  ["katayama", "片山潜"],
  ["lenin", "列宁"],
  ["luxemburg", "卢森堡"],
  ["kollontai", "柯伦泰"],
  ["li-dazhao", "李大钊"],
  ["ho-chi-minh", "胡志明"],
  ["gramsci", "葛兰西"],
  ["mao", "毛泽东"],
  ["castro", "卡斯特罗"],
  ["guevara", "切·格瓦拉"],
  ["sankara", "桑卡拉"],
]

const figs = people
  .map(
    ([k, name]) =>
      `<figure><img src="../public/assets/portraits/${k}.webp" alt="${name}"><figcaption>${name}</figcaption></figure>`,
  )
  .join("")

const page =
  '<!doctype html><html><head><meta charset="utf-8"><style>' +
  'body{background:#e3dbcd;font-family:"Microsoft YaHei",sans-serif;display:flex;flex-wrap:wrap;gap:16px;padding:26px;margin:0}' +
  'figure{margin:0;border:3px solid #111;background:#f2ede4;box-shadow:7px 7px 0 #111}' +
  'figcaption{font-weight:900;text-align:center;padding:7px;font-size:14px;border-top:3px solid #111}' +
  'img{display:block;width:150px;height:150px;object-fit:cover;' +
  'filter:grayscale(1) contrast(1.06) sepia(.18)}' +
  '</style></head><body>' +
  figs +
  '</body></html>'

const url = await writeTempHtml("shots", "portraits.html", page)
const cdp = await launchCdp({ port: PORT, windowSize: { width: W, height: H } })
const sid = await cdp.newTab(url)
await cdp.send(
  "Emulation.setDeviceMetricsOverride",
  { width: W, height: H, deviceScaleFactor: 1, mobile: false },
  sid,
)
await sleep(1200)

savePng("shots", "portraits.png", await cdp.screenshot(sid, { captureBeyondViewport: true }))
console.log("已生成 portraits.png")
cdp.dispose()
