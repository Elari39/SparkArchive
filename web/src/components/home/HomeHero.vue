<script setup lang="ts">
import { computed } from 'vue'
import { useAsync } from '@/composables/useAsync'
import { getPeople } from '@/api'
import type { PersonSummary } from '@/api/types'
import { SITE } from '@/config/site'
import { useSiteCopy } from '@/composables/useSiteCopy'
import { useSiteStore } from '@/stores/site'
import Portrait from '@/components/Portrait.vue'
import SearchBox from '@/components/SearchBox.vue'

const { subtitle, description } = useSiteCopy()

// 复用 site store 的 /api/meta（App.vue 已在挂载时请求），避免重复发同一请求
const site = useSiteStore()
const meta = computed(() => site.meta)
const people = useAsync((signal) => getPeople(signal))

/** 门楣档案墙：最多 5 张人物照片，其余补构成主义色块，凑满 3×2 网格 */
type Cell =
  | { kind: 'person'; person: PersonSummary }
  | { kind: 'decor'; decor: 'brass' | 'red' | 'ink' }

// 档案墙网格参数：最多 5 张照片、每行 3 格、共 2 行
const MAX_PHOTOS = 5
const GRID_COLS = 3
const GRID_ROWS = 2
const WALL_SLOTS = GRID_COLS * GRID_ROWS
const DECORS = ['brass', 'red', 'ink'] as const

const wallCells = computed<Cell[]>(() => {
  const list = (people.data.value ?? []).slice(0, MAX_PHOTOS)
  const cells: Cell[] = list.map((person) => ({ kind: 'person', person }))
  for (let i = 0; cells.length < WALL_SLOTS; i++) {
    cells.push({ kind: 'decor', decor: DECORS[i % DECORS.length] })
  }
  return cells
})

const DECOR_BG: Record<'brass' | 'red' | 'ink', string> = {
  brass: 'bg-brass',
  red: 'bg-red',
  ink: 'bg-ink',
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
          <div class="grid grid-cols-3 gap-3">
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
