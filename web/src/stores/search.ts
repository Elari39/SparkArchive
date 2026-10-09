import { defineStore } from 'pinia'
import { ref } from 'vue'
import { search as apiSearch } from '@/api'
import type { SearchHit } from '@/api/types'

export const useSearchStore = defineStore('search', () => {
  const query = ref('')
  const kind = ref('')
  const hits = ref<SearchHit[]>([])
  const total = ref(0)
  const loading = ref(false)
  const error = ref('')
  let seq = 0

  async function run(q: string, k = '') {
    query.value = q
    kind.value = k
    error.value = ''
    const trimmed = q.trim()
    if (!trimmed) {
      hits.value = []
      total.value = 0
      return
    }
    const mine = ++seq
    loading.value = true
    try {
      const res = await apiSearch({ q: trimmed, type: k || undefined, limit: 100 })
      // 丢弃过期响应，避免快速输入时旧结果覆盖新结果
      if (mine !== seq) return
      hits.value = res.items
      total.value = res.total
    } catch (e) {
      if (mine !== seq) return
      error.value = e instanceof Error ? e.message : '检索失败'
      hits.value = []
      total.value = 0
    } finally {
      if (mine === seq) loading.value = false
    }
  }

  return { query, kind, hits, total, loading, error, run }
})
