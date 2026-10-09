// 无头 Chrome + CDP 的公共样板：启动、等待、建标签页、求值、截图、收尾。
// render-check / og-cover / portrait-sheet 三个脚本共用此模块，避免各自复制整套样板。
//
// Chrome 路径优先取 CHROME_PATH 环境变量，其次按平台探测常见安装位置。
import { spawn } from "node:child_process"
import { existsSync, mkdirSync, writeFileSync } from "node:fs"

export const sleep = (ms: number): Promise<void> => new Promise((r) => setTimeout(r, ms))

/** 按平台候选顺序探测 Chrome 可执行文件；找不到时返回 env 值（可能为空串）。 */
export function resolveChrome(): string {
  const env = process.env.CHROME_PATH
  if (env) return env
  const candidates =
    process.platform === "win32"
      ? [
          "C:/Program Files/Google/Chrome/Application/chrome.exe",
          "C:/Program Files (x86)/Google/Chrome/Application/chrome.exe",
          `${process.env.LOCALAPPDATA ?? ""}/Google/Chrome/Application/chrome.exe`,
        ]
      : process.platform === "darwin"
        ? ["/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"]
        : ["/usr/bin/google-chrome", "/usr/bin/chromium", "/usr/bin/chromium-browser"]
  for (const p of candidates) if (p && existsSync(p)) return p
  throw new Error("未找到 Chrome，请设置 CHROME_PATH 环境变量指向可执行文件")
}

export interface CDPMessage {
  id?: number
  method?: string
  params?: Record<string, unknown>
  result?: Record<string, unknown>
  sessionId?: string
}

type Listener = (msg: CDPMessage) => void

export interface CdpClient {
  /** 发一条 CDP 命令并等待响应。 */
  send(method: string, params?: Record<string, unknown>, sessionId?: string): Promise<CDPMessage>
  /** 新建标签页、开启 Page/Runtime、导航到 url，返回 sessionId。 */
  newTab(url: string): Promise<string>
  /** 在指定 session 里求值表达式（按值返回，等待 Promise）。 */
  evaluate(sessionId: string, expression: string): Promise<unknown>
  /** 截图，返回 PNG Buffer。 */
  screenshot(
    sessionId: string,
    opts?: { captureBeyondViewport?: boolean },
  ): Promise<Buffer>
  /** 关闭指定 session 对应的标签页。 */
  closeTab(sessionId: string): Promise<void>
  /** 订阅 CDP 事件（如 Runtime.exceptionThrown）。返回取消订阅函数。 */
  on(fn: Listener): () => void
  /** 关闭 WebSocket 并杀掉 Chrome 进程。 */
  dispose(): void
}

/** 启动无头 Chrome 并连接 CDP，返回客户端。 */
export async function launchCdp(opts: {
  port: number
  windowSize: { width: number; height: number }
  extraArgs?: string[]
}): Promise<CdpClient> {
  const chrome = resolveChrome()
  const { port, windowSize, extraArgs = [] } = opts
  const proc = spawn(
    chrome,
    [
      "--headless=new",
      "--disable-gpu",
      "--no-sandbox",
      "--hide-scrollbars",
      "--force-device-scale-factor=1",
      `--window-size=${windowSize.width},${windowSize.height}`,
      `--remote-debugging-port=${port}`,
      ...extraArgs,
      "about:blank",
    ],
    { stdio: "ignore" },
  )

  let wsUrl: string | undefined
  for (let i = 0; i < 40; i++) {
    try {
      const res = await fetch(`http://127.0.0.1:${port}/json/version`)
      const j = (await res.json()) as { webSocketDebuggerUrl?: string }
      if (j.webSocketDebuggerUrl) {
        wsUrl = j.webSocketDebuggerUrl
        break
      }
    } catch {
      /* 尚未就绪 */
    }
    await sleep(250)
  }
  if (!wsUrl) {
    proc.kill()
    throw new Error("Chrome CDP 未就绪")
  }

  const ws = new WebSocket(wsUrl)
  await new Promise<void>((res, rej) => {
    ws.onopen = () => res()
    ws.onerror = (e) => rej(e)
  })

  let id = 0
  const pending = new Map<number, (m: CDPMessage) => void>()
  const listeners = new Set<Listener>()
  const targets = new Map<string, string>()

  ws.onmessage = (ev) => {
    const msg = JSON.parse(ev.data as string) as CDPMessage
    if (msg.id && pending.has(msg.id)) {
      pending.get(msg.id)!(msg)
      pending.delete(msg.id)
      return
    }
    for (const fn of listeners) fn(msg)
  }

  const send = (method: string, params: Record<string, unknown> = {}, sessionId?: string) =>
    new Promise<CDPMessage>((res) => {
      const myId = ++id
      pending.set(myId, res)
      ws.send(JSON.stringify({ id: myId, method, params, sessionId }))
    })

  const newTab = async (url: string) => {
    const t = await send("Target.createTarget", { url: "about:blank" })
    const targetId = t.result!.targetId as string
    const att = await send("Target.attachToTarget", { targetId, flatten: true })
    const sid = att.result!.sessionId as string
    targets.set(sid, targetId)
    await send("Page.enable", {}, sid)
    await send("Runtime.enable", {}, sid)
    await send("Page.navigate", { url }, sid)
    return sid
  }

  const evaluate = async (sessionId: string, expression: string) => {
    const r = await send(
      "Runtime.evaluate",
      { expression, returnByValue: true, awaitPromise: true },
      sessionId,
    )
    if (r.result?.exceptionDetails) {
      throw new Error(JSON.stringify(r.result.exceptionDetails).slice(0, 300))
    }
    return r.result?.result?.value
  }

  const screenshot = async (
    sessionId: string,
    o: { captureBeyondViewport?: boolean } = {},
  ) => {
    const r = await send(
      "Page.captureScreenshot",
      { format: "png", captureBeyondViewport: o.captureBeyondViewport ?? false },
      sessionId,
    )
    return Buffer.from(r.result!.data as string, "base64")
  }

  const closeTab = async (sessionId: string) => {
    const targetId = targets.get(sessionId)
    if (targetId) await send("Target.closeTarget", { targetId })
    targets.delete(sessionId)
  }

  const on = (fn: Listener) => {
    listeners.add(fn)
    return () => listeners.delete(fn)
  }

  const dispose = () => {
    ws.close()
    proc.kill()
  }

  return { send, newTab, evaluate, screenshot, closeTab, on, dispose }
}

/** 确保输出目录存在并写入文件。 */
export function savePng(dir: string, name: string, data: Buffer): string {
  mkdirSync(dir, { recursive: true })
  const path = `${dir}/${name}`
  writeFileSync(path, data)
  return path
}

/** 写一个临时 HTML 文件到 shots 目录，返回其 file:// URL。 */
export async function writeTempHtml(
  dir: string,
  name: string,
  html: string,
): Promise<string> {
  const { pathToFileURL } = await import("node:url")
  const { resolve } = await import("node:path")
  mkdirSync(dir, { recursive: true })
  const file = resolve(`${dir}/${name}`)
  writeFileSync(file, html, "utf8")
  return pathToFileURL(file).href
}
