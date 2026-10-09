// 契约校验：用前端的 zod schema 校验后端真实响应。
// 这能捕捉到冒烟测试覆盖不到的集成问题——即 API 返回的数据前端无法解析。
import {
  PersonSummary, PersonDetail, WorkSummary, WorkDetail, EventItem,
  TermSummary, TermDetail, GraphData, SearchResult, SiteMeta,
} from '../src/api/types.ts'

// 默认地址统一用 127.0.0.1：Docker Desktop 端口代理只监听 IPv4，
// 而 localhost 优先解析到 ::1，会让 .NET/Node 侧等待 TCP 超时（详见 README）。
const BASE = process.env.API_BASE ?? 'http://127.0.0.1:12026/api'

let fail = 0
let pass = 0

async function check(name, path, schema, pick) {
  try {
    const res = await fetch(BASE + path)
    if (!res.ok) throw new Error('HTTP ' + res.status)
    const raw = await res.json()
    const data = pick ? pick(raw) : raw
    const parsed = schema.safeParse(data)
    if (parsed.success) {
      pass++
      console.log('  [PASS] ' + name)
    } else {
      fail++
      console.log('  [FAIL] ' + name)
      for (const issue of parsed.error.issues.slice(0, 6)) {
        console.log('         ' + issue.path.join('.') + ': ' + issue.message)
      }
    }
  } catch (e) {
    fail++
    console.log('  [FAIL] ' + name + ' -- ' + e.message)
  }
}

console.log('')
console.log('前后端契约校验 (前端 zod schema vs 后端响应)')
console.log('')

await check('GET /meta', '/meta', SiteMeta)
await check('GET /people', '/people', PersonSummary.array())
await check('GET /works', '/works', WorkSummary.array())
await check('GET /events', '/events', EventItem.array())
await check('GET /terms', '/terms', TermSummary.array())
await check('GET /relations', '/relations', GraphData)
await check('GET /search?q=无产阶级', '/search?q=' + encodeURIComponent('无产阶级'), SearchResult)

for (const slug of ['marx', 'engels', 'lenin', 'mao', 'guevara']) {
  await check('GET /people/' + slug, '/people/' + slug, PersonDetail)
}
for (const slug of ['imperialism', 'on-practice', 'guerrilla-warfare']) {
  await check('GET /works/' + slug, '/works/' + slug, WorkDetail)
}
for (const slug of ['proletariat', 'foquismo', 'mass-line']) {
  await check('GET /terms/' + slug, '/terms/' + slug, TermDetail)
}

console.log('')
console.log(fail === 0 ? '契约全部匹配: ' + pass + ' 项' : '通过 ' + pass + ' 项, 失败 ' + fail + ' 项')
process.exit(fail === 0 ? 0 : 1)
