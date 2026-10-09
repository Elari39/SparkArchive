<script setup lang="ts">
import { computed } from 'vue'
import { useAsync } from '@/composables/useAsync'
import { getPeople } from '@/api'
import type { PersonSummary } from '@/api/types'
import { SITE, HERO_WALL } from '@/config/site'
import { useSiteCopy } from '@/composables/useSiteCopy'
import { useSiteStore } from '@/stores/site'
import Portrait from '@/components/Portrait.vue'
import SearchBox from '@/components/SearchBox.vue'

const { subtitle, description } = useSiteCopy()

// 复用 site store 的 /api/meta（App.vue 已在挂载时请求），避免重复发同一请求
const site = useSiteStore()
const meta = computed(() => site.meta)
const people = useAsync((signal) => getPeople(signal))

/** 门楣档案墙：按 HERO_WALL 的策展顺序填满 3×2 网格；人数不足时才补构成主义色块 */
type Cell =
  | { kind: 'person'; person: PersonSummary }
  | { kind: 'decor'; decor: 'ink' | 'red' | 'brass' }

// 档案墙网格参数：每行 3 格、共 2 行，正好容纳 6 位人物
const GRID_COLS = 3
const GRID_ROWS = 2
const WALL_SLOTS = GRID_COLS * GRID_ROWS

// 兜底色块：仅在人物数量不足 6 位时才出现。
// 顺序刻意以深色打头——工业黄是整块版面唯一的亮饱和色，一旦出现在缺口处会抢走视线，
// 因此把它排在最后，确保降级状态也不会重现「缺口被黄块占据」的观感问题。
const DECORS = ['ink', 'red', 'brass'] as const

const wallCells = computed<Cell[]>(() => {
  const all = people.data.value ?? []
  const bySlug = new Map(all.map((person) => [person.slug, person]))

  // 先按策展顺序取，再按后端顺序（出生日期升序）补足，最后才用色块兜底。
  const picked: PersonSummary[] = []
  const chosen = new Set<string>()
  const take = (person: PersonSummary | undefined) => {
    if (!person || chosen.has(person.slug) || picked.length >= WALL_SLOTS) return
    picked.push(person)
    chosen.add(person.slug)
  }

  HERO_WALL.forEach((slug) => take(bySlug.get(slug)))
  all.forEach(take)

  const cells: Cell[] = picked.map((person) => ({ kind: 'person', person }))
  for (let i = 0; cells.length < WALL_SLOTS; i++) {
    cells.push({ kind: 'decor', decor: DECORS[i % DECORS.length] })
  }
  return cells
})

const DECOR_BG: Record<'ink' | 'red' | 'brass', string> = {
  ink: 'bg-ink',
  red: 'bg-red',
  brass: 'bg-brass',
}
</script>

<template>
  <!-- 构成主义门楣 -->
  <section class="relative overflow-hidden border-b-3 border-ink bg-paper">
    <div class="absolute inset-0 hatch opacity-[0.05]" aria-hidden="true" />
    <div class="relative mx-auto max-w-6xl px-4 py-14 sm:px-6 sm:py-20">
      <div class="grid gap-10 lg:grid-cols-[1.4fr_1fr] lg:items-center">
        <div>
          <div class="mb-4 flex flex-wrap gap-2">
            <span class="stamp stamp-red">{{ SITE.archiveNo }}</span>
            <span class="stamp">{{ SITE.est }}</span>
            <span class="stamp">中文</span>
          </div>
          <h1 class="display-hero">
            星火<span class="text-red">档案</span>馆
          </h1>
          <p class="mt-3 font-mono text-xs tracking-[0.3em] text-steel uppercase sm:text-sm">
            {{ SITE.nameEn }}
          </p>
          <p class="mt-6 max-w-2xl border-l-[6px] border-red pl-4 text-base leading-relaxed sm:text-lg">
            {{ subtitle }}
          </p>
          <p class="mt-4 max-w-2xl text-sm leading-relaxed text-ink-soft">
            {{ description }}
          </p>
          <div class="mt-7 flex flex-col gap-3 sm:flex-row sm:items-center">
            <SearchBox />
          </div>
          <p
            v-if="meta"
            class="mt-5 flex flex-wrap gap-x-5 gap-y-1 font-mono text-[0.7rem] tracking-wider text-steel"
          >
            <span>人物 {{ meta.counts.people }}</span>
            <span>著作 {{ meta.counts.works }}</span>
            <span>事件 {{ meta.counts.events }}</span>
            <span>术语 {{ meta.counts.terms }}</span>
          </p>
        </div>

        <!-- 真实照片档案墙 -->
        <div class="relative hidden lg:block" aria-hidden="true">
          <div class="grid grid-cols-3 gap-3" data-hero-wall>
            <template v-for="(c, i) in wallCells" :key="i">
              <figure
                v-if="c.kind === 'person'"
                class="border-3 border-ink bg-paper shadow-[var(--shadow-hard-sm)]"
              >
                <div class="aspect-square">
                  <Portrait
                    :person="c.person.portrait_key"
                    :name="c.person.name"
                    fill
                    loading="lazy"
                  />
                </div>
              </figure>
              <div
                v-else
                class="aspect-square border-3 border-ink"
                :class="DECOR_BG[c.decor]"
              />
            </template>
          </div>
          <div class="mt-3 flex items-center justify-center border-3 border-ink bg-ink px-4 py-4">
            <p class="text-center font-display text-lg leading-tight font-black text-paper">
              全世界无产者，联合起来！
            </p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
