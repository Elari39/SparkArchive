<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useAsync } from '@/composables/useAsync'
import { getGraph } from '@/api'
import PageHeader from '@/components/PageHeader.vue'
import AsyncBoundary from '@/components/AsyncBoundary.vue'

const { data: graph, loading, error, reload } = useAsync((signal) => getGraph(signal))
const active = ref<string | null>(null)

/** 固定环形布局：人物均匀分布，保证图谱稳定可预期。
 *
 *  节点为直径 80（r=40）的圆盘，内部居中排布人名，因此相邻节点圆心的弦长
 *  必须大于直径，否则圆盘相互重叠、标签被遮住。收录到 16 人后，原固定半径
 *  R=150 的弦长只有约 58，已明显重叠，故半径改为按人数推导：
 *
 *    R = max(基准半径, 不重叠所需的最小半径)
 *
 *  画布高度随之增长；人数不超过 9 时结果与原先完全一致（R=150、画布 720×460），
 *  既有观感不变。 */
const W = 720
const CX = W / 2
/** 基准半径与画布高（9 人时的原始取值） */
const BASE_R = 150
const BASE_H = 460
/** 节点半径与相邻圆盘之间保留的间隙 */
const NODE_R = 40
const NODE_GAP = 6

const R = computed(() => {
  const n = graph.value?.nodes.length ?? 0
  if (n < 2) return BASE_R
  // 弦长 = 2R·sin(π/n)，令其 ≥ 2·NODE_R + NODE_GAP
  const need = (NODE_R + NODE_GAP / 2) / Math.sin(Math.PI / n)
  return Math.max(BASE_R, need)
})
const H = computed(() => Math.max(BASE_H, Math.ceil((R.value + NODE_R + 10) * 2)))
const CY = computed(() => H.value / 2)

const layout = computed(() => {
  const nodes = graph.value?.nodes ?? []
  const r = R.value
  const cy = CY.value
  const pos = new Map<string, { x: number; y: number }>()
  nodes.forEach((n, i) => {
    const a = (i / Math.max(nodes.length, 1)) * Math.PI * 2 - Math.PI / 2
    pos.set(n.slug, { x: CX + r * Math.cos(a), y: cy + r * Math.sin(a) })
  })
  return { nodes, pos }
})

const edges = computed(() => {
  const g = graph.value
  if (!g) return []
  const { pos } = layout.value
  return g.edges
    .map((e, i) => {
      const a = pos.get(e.from)
      const b = pos.get(e.to)
      if (!a || !b) return null
      return { ...e, i, a, b }
    })
    .filter((x): x is NonNullable<typeof x> => x !== null)
})

const highlight = computed(() => {
  if (!active.value) return { nodes: new Set<string>(), edges: new Set<number>() }
  const nodes = new Set<string>([active.value])
  const es = new Set<number>()
  edges.value.forEach((e) => {
    if (e.from === active.value || e.to === active.value) {
      nodes.add(e.from)
      nodes.add(e.to)
      es.add(e.i)
    }
  })
  return { nodes, edges: es }
})

const nameOf = (slug: string) => layout.value.nodes.find((n) => n.slug === slug)?.name ?? slug
const activeEdges = computed(() => edges.value.filter((e) => e.from === active.value || e.to === active.value))
</script>

<template>
  <div>
    <PageHeader
      no="05"
      title="人物关系图"
      sub="Relations"
      lead="思想传承与国际联系的图谱。点击节点可高亮该人物的全部关系；下方同时提供无需图形的等价列表，便于键盘与读屏访问。"
    />
    <section class="mx-auto max-w-6xl px-4 py-10 sm:px-6">
      <AsyncBoundary :loading="loading" :error="error" :on-retry="reload">
        <div class="brutal-card overflow-x-auto p-2">
          <svg :viewBox="`0 0 ${W} ${H}`" class="mx-auto block h-auto w-full min-w-[560px]" role="img"
            aria-label="人物关系图谱">
            <rect :width="W" :height="H" fill="#f2ede4" />
            <g opacity="0.06">
              <path v-for="i in 14" :key="i" :d="`M0 ${i * 46} L${W} ${i * 46}`" stroke="#111" stroke-width="2" />
            </g>

            <!-- 连线 -->
            <g v-for="e in edges" :key="`e${e.i}`">
              <line
                :x1="e.a.x"
                :y1="e.a.y"
                :x2="e.b.x"
                :y2="e.b.y"
                :stroke="highlight.edges.has(e.i) ? '#d62828' : '#111111'"
                :stroke-width="highlight.edges.has(e.i) ? 5 : 2.5"
                :opacity="active && !highlight.edges.has(e.i) ? 0.15 : 1"
                stroke-dasharray="none"
              />
              <text
                :x="(e.a.x + e.b.x) / 2"
                :y="(e.a.y + e.b.y) / 2 - 5"
                text-anchor="middle"
                font-size="11"
                font-weight="700"
                fill="#4a4a4a"
                :opacity="active && !highlight.edges.has(e.i) ? 0.2 : 1"
              >
                {{ e.kind }}
              </text>
            </g>

            <!-- 节点 -->
            <g
              v-for="n in layout.nodes"
              :key="n.slug"
              class="cursor-pointer"
              :opacity="active && !highlight.nodes.has(n.slug) ? 0.3 : 1"
              @click="active = active === n.slug ? null : n.slug"
            >
              <circle
                :cx="layout.pos.get(n.slug)?.x ?? 0"
                :cy="layout.pos.get(n.slug)?.y ?? 0"
                r="40"
                :fill="active === n.slug ? '#d62828' : '#111111'"
                stroke="#f2ede4"
                stroke-width="3"
              />
              <text
                :x="layout.pos.get(n.slug)?.x ?? 0"
                :y="(layout.pos.get(n.slug)?.y ?? 0) + 5"
                text-anchor="middle"
                font-size="15"
                font-weight="900"
                fill="#f2ede4"
              >
                {{ n.name.slice(-3) }}
              </text>
            </g>
          </svg>
        </div>

        <p v-if="active" class="mt-4 flex flex-wrap items-center gap-3">
          <span class="stamp stamp-red">当前聚焦 {{ nameOf(active) }}</span>
          <button type="button" class="font-mono text-xs font-bold underline" @click="active = null">
            清除聚焦
          </button>
        </p>

        <!-- 等价列表（可访问性） -->
        <h2 class="mt-10 mb-4 border-b-3 border-ink pb-2 text-2xl font-black">关系列表</h2>
        <div class="grid gap-4 md:grid-cols-2">
          <article
            v-for="e in active ? activeEdges : edges"
            :key="`l${e.i}`"
            class="brutal-card border-l-[10px] border-l-red p-4"
          >
            <div class="mb-2 flex flex-wrap items-center gap-2">
              <RouterLink :to="`/people/${e.from}`" class="font-bold hover:text-red-deep">
                {{ nameOf(e.from) }}
              </RouterLink>
              <span class="font-mono text-xs text-red">→</span>
              <RouterLink :to="`/people/${e.to}`" class="font-bold hover:text-red-deep">
                {{ nameOf(e.to) }}
              </RouterLink>
              <span class="stamp ml-auto">{{ e.kind }}</span>
            </div>
            <p class="text-sm font-bold">{{ e.label }}</p>
            <p class="mt-1.5 text-xs leading-relaxed text-ink-soft">{{ e.note }}</p>
          </article>
        </div>
      </AsyncBoundary>
    </section>
  </div>
</template>
