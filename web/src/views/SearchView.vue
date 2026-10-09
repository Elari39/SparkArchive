<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useSearchStore } from '@/stores/search'
import PageHeader from '@/components/PageHeader.vue'
import { escapeHtml, sanitizeHighlight } from '@/utils/sanitize'

const route = useRoute()
const router = useRouter()
const store = useSearchStore()
const input = ref('')

const KIND_LABELS: Record<string, string> = {
  person: '人物',
  person_section: '人物·贡献',
  work: '著作',
  event: '事件',
  term: '术语',
}
const KINDS = [
  { key: '', label: '全部' },
  { key: 'person', label: '人物' },
  { key: 'work', label: '著作' },
  { key: 'event', label: '事件' },
  { key: 'term', label: '术语' },
]

function routePath(hit: { kind: string; ref: string }) {
  if (hit.kind === 'person' || hit.kind === 'person_section') return `/people/${hit.ref}`
  if (hit.kind === 'work') return `/works/${hit.ref}`
  if (hit.kind === 'term') return `/glossary/${hit.ref}`
  return '/timeline'
}

/** 把命中词条高亮：先转义（含单引号）再替换，最后经白名单净化，防止 XSS */
function highlight(text: string, q: string) {
  const safe = escapeHtml(text)
  if (!q.trim()) return safe
  const esc = q.trim().replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  try {
    const marked = safe.replace(new RegExp(`(${esc})`, 'gi'), '<mark class="hit">$1</mark>')
    return sanitizeHighlight(marked)
  } catch {
    return safe
  }
}

const q = computed(() => String(route.query.q ?? ''))
const kind = computed(() => String(route.query.type ?? ''))

function submit() {
  void router.push({ name: 'search', query: { q: input.value.trim(), type: kind.value || undefined } })
}

onMounted(() => {
  input.value = q.value
  void store.run(q.value, kind.value)
})

watch([q, kind], () => {
  input.value = q.value
  void store.run(q.value, kind.value)
})
</script>

<template>
  <div>
    <PageHeader
      no="06"
      title="全文检索"
      sub="Search"
      lead="跨人物、生平、贡献、著作、事件与术语的全文检索。支持中文子串匹配与结果高亮。"
    />
    <section class="mx-auto max-w-4xl px-4 py-10 sm:px-6">
      <form class="mb-4 flex" role="search" @submit.prevent="submit">
        <label for="q" class="sr-only">检索词</label>
        <input
          id="q"
          v-model="input"
          type="search"
          placeholder="输入关键词，如：无产阶级、十月革命、游击"
          class="min-w-0 flex-1 border-3 border-ink bg-paper px-3 py-2.5 font-medium"
        />
        <button
          type="submit"
          class="shrink-0 border-3 border-l-0 border-ink bg-red px-5 font-mono text-xs font-bold tracking-widest text-paper uppercase"
        >
          检索
        </button>
      </form>

      <div class="mb-6 flex flex-wrap gap-2">
        <RouterLink
          v-for="k in KINDS"
          :key="k.key"
          :to="{ name: 'search', query: { q, type: k.key || undefined } }"
          class="border-2 border-ink px-3 py-1 font-mono text-xs font-bold"
          :class="kind === k.key ? 'bg-ink text-paper' : 'bg-paper hover:bg-paper-dim'"
        >
          {{ k.label }}
        </RouterLink>
      </div>

      <p v-if="store.loading" class="font-mono text-sm">检索中…</p>
      <p v-else-if="store.error" class="brutal-card border-l-[10px] border-l-red p-4 font-bold">
        {{ store.error }}
      </p>
      <template v-else-if="q">
        <p class="mb-5 font-mono text-xs tracking-wider text-steel">
          命中 <strong class="text-red">{{ store.total }}</strong> 条 · 关键词「{{ q }}」
        </p>
        <p v-if="!store.hits.length" class="brutal-card border-l-[10px] border-l-steel p-5 font-bold">
          没有找到相关内容。可尝试更短的关键词。
        </p>
        <ol v-else class="space-y-4">
          <li v-for="(h, i) in store.hits" :key="i">
            <RouterLink :to="routePath(h)" class="block">
              <article class="brutal-card brutal-interactive p-5">
                <div class="mb-2 flex flex-wrap items-center gap-2">
                  <span class="stamp stamp-red">{{ KIND_LABELS[h.kind] ?? h.kind }}</span>
                  <span v-if="h.person" class="stamp">{{ h.person }}</span>
                </div>
                <h2 class="font-display text-lg leading-tight font-black">{{ h.title }}</h2>
                <p class="mt-2 text-sm leading-relaxed text-ink-soft" v-html="highlight(h.snippet, q)" />
              </article>
            </RouterLink>
          </li>
        </ol>
      </template>
      <p v-else class="font-mono text-sm text-steel">输入关键词开始检索。</p>
    </section>
  </div>
</template>
