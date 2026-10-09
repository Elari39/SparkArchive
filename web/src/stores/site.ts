import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getMeta } from '@/api'
import type { SiteMeta } from '@/api/types'

export const useSiteStore = defineStore('site', () => {
  const meta = ref<SiteMeta | null>(null)
  const error = ref('')

  async function load() {
    if (meta.value) return
    try {
      meta.value = await getMeta()
    } catch (e) {
      error.value = e instanceof Error ? e.message : '未知错误'
    }
  }

  return { meta, error, load }
})
