import { onScopeDispose, ref, watch, type Ref } from 'vue'

/** 极简数据加载：返回 data/loading/error，依赖变化时自动重取。 */
export function useAsync<T>(loader: (signal?: AbortSignal) => Promise<T>, deps: Ref<unknown>[] = []) {
  const data = ref<T | null>(null) as Ref<T | null>
  const loading = ref(true)
  const error = ref('')

  // 序号守卫：依赖快速变化时只允许最新一次请求写回结果，
  // 避免旧请求的响应覆盖新结果（与 stores/search.ts 的策略一致）。
  let seq = 0
  let controller: AbortController | undefined

  async function run() {
    const mine = ++seq
    controller?.abort()
    controller = new AbortController()
    loading.value = true
    error.value = ''
    try {
      const result = await loader(controller.signal)
      if (mine !== seq) return
      data.value = result
    } catch (e) {
      if (mine !== seq) return
      // 主动 abort 不算错误
      if (e instanceof DOMException && e.name === 'AbortError') return
      error.value = e instanceof Error ? e.message : '加载失败'
      data.value = null
    } finally {
      if (mine === seq) loading.value = false
    }
  }

  if (deps.length) watch(deps, run, { immediate: true })
  else void run()

  // 组件卸载时中止未完成请求
  onScopeDispose(() => controller?.abort())

  return { data, loading, error, reload: run }
}
