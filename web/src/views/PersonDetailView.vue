<script setup lang="ts">
import { computed, toRef } from 'vue'
import { RouterLink } from 'vue-router'
import { useAsync } from '@/composables/useAsync'
import { getPerson } from '@/api'
import { SECTION_KINDS } from '@/config/site'
import Portrait from '@/components/Portrait.vue'
import MarkdownBlock from '@/components/MarkdownBlock.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import AsyncBoundary from '@/components/AsyncBoundary.vue'

const props = defineProps<{ slug: string }>()
const slug = toRef(props, 'slug')
const { data: person, loading, error, reload } = useAsync((signal) => getPerson(slug.value, signal), [slug])

/** 按 政治/经济/文化/思想 固定顺序分组，保证每页结构一致 */
const grouped = computed(() => {
  const sections = person.value?.sections ?? []
  return SECTION_KINDS.map((k) => ({
    ...k,
    items: sections.filter((s) => s.kind === k.key).sort((a, b) => a.ord - b.ord),
  })).filter((g) => g.items.length > 0)
})

const lifeSpan = computed(() => {
  if (!person.value) return ''
  return `${person.value.birth} — ${person.value.death}`
})
</script>

<template>
  <div>
    <AsyncBoundary
      :loading="loading"
      :error="error"
      :empty="!person ? '未找到该人物档案' : ''"
      :on-retry="reload"
      state-class="mx-auto max-w-6xl px-4"
      empty-class="mx-auto my-10 max-w-2xl"
    >

    <template v-if="person">
      <!-- 档案页头 -->
      <header class="border-b-3 border-ink bg-paper-dim">
        <div class="mx-auto max-w-6xl px-4 py-10 sm:px-6">
          <Breadcrumbs :items="[{ to: '/people', label: '人物' }, { label: person.name }]" />
          <div class="grid gap-7 sm:grid-cols-[auto_1fr]">
            <Portrait
              :person="person.portrait_key"
              :size="200"
              :name="person.name"
              loading="eager"
              class="border-3 border-ink shadow-[var(--shadow-hard)]"
            />
            <div class="min-w-0">
              <div class="mb-3 flex flex-wrap gap-2">
                <span class="stamp stamp-red">{{ person.nationality }}</span>
                <span class="stamp">{{ person.birth.slice(0, 4) }}–{{ person.death.slice(0, 4) }}</span>
              </div>
              <h1 class="display-hero">{{ person.name }}</h1>
              <p class="mt-2 font-mono text-xs tracking-[0.25em] text-steel uppercase">{{ person.name_en }}</p>
              <p class="mt-3 text-lg font-bold text-red-deep">{{ person.epithet }}</p>
              <dl class="mt-5 grid gap-2 font-mono text-xs sm:grid-cols-2">
                <div class="flex gap-2">
                  <dt class="text-steel">生卒</dt>
                  <dd class="font-bold">{{ lifeSpan }}</dd>
                </div>
                <div class="flex gap-2">
                  <dt class="text-steel">出生地</dt>
                  <dd class="font-bold">{{ person.birth_place }}</dd>
                </div>
              </dl>
            </div>
          </div>

          <blockquote
            v-if="person.thesis"
            class="mt-7 border-l-[10px] border-red bg-ink px-5 py-4 font-display text-xl leading-snug font-black text-paper sm:text-2xl"
          >
            「{{ person.thesis }}」
          </blockquote>
        </div>
      </header>

      <div class="mx-auto grid max-w-6xl gap-10 px-4 py-10 sm:px-6 lg:grid-cols-[1fr_300px]">
        <div class="min-w-0">
          <section class="mb-10">
            <h2 class="mb-3 border-b-3 border-ink pb-2 text-2xl font-black">简介</h2>
            <p class="prose-archive">{{ person.summary }}</p>
          </section>

          <!-- 生平年表 -->
          <section class="mb-12">
            <h2 class="mb-5 border-b-3 border-ink pb-2 text-2xl font-black">生平年表</h2>
            <ol class="relative border-l-3 border-ink pl-6">
              <li v-for="(t, i) in person.timeline" :key="i" class="relative mb-6 last:mb-0">
                <span
                  class="absolute top-1.5 -left-[1.9rem] h-3 w-3 border-2 border-ink bg-red"
                  aria-hidden="true"
                />
                <p class="font-mono text-xs font-bold tracking-wider text-red-deep">
                  {{ t.year }}<span v-if="t.date_text"> · {{ t.date_text }}</span>
                </p>
                <h3 class="mt-0.5 font-display text-base font-black">{{ t.title }}</h3>
                <p class="mt-1 text-sm leading-relaxed text-ink-soft">{{ t.body }}</p>
              </li>
            </ol>
          </section>

          <!-- 贡献 / 思想：四域分组 -->
          <section v-for="g in grouped" :key="g.key" class="mb-12">
            <div class="mb-5 flex items-end gap-3">
              <span class="border-3 border-ink bg-ink px-2 py-0.5 font-mono text-base font-black text-paper">
                {{ g.no }}
              </span>
              <h2 class="display-xl">{{ g.label }}</h2>
              <div class="mb-2.5 hidden h-1 flex-1 bg-ink sm:block" />
            </div>
            <article v-for="(s, i) in g.items" :key="i" class="brutal-card mb-5 p-5">
              <h3 class="mb-2 font-display text-lg font-black text-red-deep">{{ s.title }}</h3>
              <MarkdownBlock :source="s.body_md" />
            </article>
          </section>

          <!-- 争议与评价 -->
          <section v-if="(person.sections ?? []).some((s) => s.kind === 'controversy')" class="mb-12">
            <h2 class="mb-4 border-b-3 border-ink pb-2 text-2xl font-black">争议与评价</h2>
            <details class="brutal-card border-l-[10px] border-l-steel p-5" open>
              <summary class="cursor-pointer font-mono text-sm font-bold tracking-wider">
                展开 / 收起 — 历史评价中的争议
              </summary>
              <div class="mt-4">
                <MarkdownBlock
                  v-for="(s, i) in person.sections.filter((x) => x.kind === 'controversy')"
                  :key="i"
                  :source="s.body_md"
                />
              </div>
            </details>
          </section>
        </div>

        <!-- 侧栏 -->
        <aside class="space-y-6 lg:sticky lg:top-20 lg:self-start">
          <div v-if="person.works.length" class="brutal-card p-4">
            <h2 class="mb-3 border-b-2 border-ink pb-1.5 font-mono text-xs font-bold tracking-widest uppercase">
              相关著作
            </h2>
            <ul class="space-y-2">
              <li v-for="w in person.works" :key="w.slug">
                <RouterLink :to="`/works/${w.slug}`" class="group flex gap-2 text-sm hover:text-red-deep">
                  <span class="font-mono text-xs text-steel">{{ w.year }}</span>
                  <span class="font-bold group-hover:underline">{{ w.title }}</span>
                </RouterLink>
              </li>
            </ul>
          </div>

          <div v-if="person.terms.length" class="brutal-card p-4">
            <h2 class="mb-3 border-b-2 border-ink pb-1.5 font-mono text-xs font-bold tracking-widest uppercase">
              相关术语
            </h2>
            <ul class="flex flex-wrap gap-2">
              <li v-for="t in person.terms" :key="t.slug">
                <RouterLink
                  :to="`/glossary/${t.slug}`"
                  class="inline-block border-2 border-ink bg-paper-dim px-2 py-0.5 font-mono text-xs font-bold hover:bg-red hover:text-paper"
                >
                  {{ t.term }}
                </RouterLink>
              </li>
            </ul>
          </div>

          <div class="brutal-card border-l-[10px] border-l-red p-4">
            <h2 class="mb-2 font-mono text-xs font-bold tracking-widest uppercase">史料说明</h2>
            <p class="text-xs leading-relaxed text-ink-soft">
              本页史实与引文均标注出处；对存在争议的历史评价设专节如实说明，不作回避。
            </p>
          </div>
        </aside>
      </div>
    </template>
    </AsyncBoundary>
  </div>
</template>
