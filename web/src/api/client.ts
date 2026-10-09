import { z } from 'zod'

const BASE = '/api'

export class ApiError extends Error {
  constructor(
    message: string,
    readonly status: number,
    /** zod 校验失败时携带的路径级问题摘要，便于定位契约漂移。 */
    readonly issues?: string[],
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

/** 把 zod 错误压缩成「path: message」列表，只取前几条避免刷屏。 */
function summarizeIssues(error: z.ZodError): string[] {
  return error.issues.slice(0, 5).map((i) => {
    const path = i.path.length > 0 ? i.path.join('.') : '(root)'
    return `${path}: ${i.message}`
  })
}

/** 统一请求 + zod 校验。校验失败会抛出，避免把畸形数据渲染进组件。 */
export async function apiGet<T>(path: string, schema: z.ZodType<T>, signal?: AbortSignal): Promise<T> {
  const res = await fetch(BASE + path, { signal, headers: { Accept: 'application/json' } })
  if (!res.ok) {
    throw new ApiError(`请求失败：${res.status} ${path}`, res.status)
  }

  const text = await res.text()
  let raw: unknown
  try {
    raw = text ? JSON.parse(text) : null
  } catch {
    throw new ApiError(`响应不是合法 JSON：${path}`, res.status)
  }

  const parsed = schema.safeParse(raw)
  if (!parsed.success) {
    throw new ApiError(`响应格式不符合预期：${path}`, res.status, summarizeIssues(parsed.error))
  }
  return parsed.data
}

export function qs(params: Record<string, string | number | undefined>): string {
  const sp = new URLSearchParams()
  for (const [k, v] of Object.entries(params)) {
    if (v !== undefined && v !== '') sp.set(k, String(v))
  }
  const s = sp.toString()
  return s ? `?${s}` : ''
}
