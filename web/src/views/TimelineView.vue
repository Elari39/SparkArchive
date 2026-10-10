<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useAsync } from '@/composables/useAsync'
import { getEvents } from '@/api'
import { PEOPLE_FILTERS } from '@/config/site'
import PageHeader from '@/components/PageHeader.vue'
import AsyncBoundary from '@/components/AsyncBoundary.vue'

const person = ref('')
const category = ref('')
const { data: events, loading, error, reload } = useAsync(
  (signal) => getEvents({ person: person.value, category: category.value }, signal),
  [person, category],
)

const CATEGORIES = ['理论', '革命', '组织', '战争', '建设', '文化', '牺牲', '生平'] as const
const PEOPLE = PEOPLE_FILTERS

const visible = computed(() => events.value ?? [])
const catColor = (c: string) =>
  c === '革命' || c === '牺牲' ? 'stamp-red' : c === '理论' || c === '文化' ? 'stamp-ink' : ''
</script>

<template>
  <div>
    <PageHeader
      no="03"
      title="国际共运大事年表"
      sub="Timeline"
      lead="从 1848 年《共产党宣言》发表到 2016 年卡斯特罗逝世，记录国际共产主义运动的理论发展与实践历程。可按范畴与人物筛选。"
    />
    <section class="mx-auto max-w-5xl px-4 py-10 sm:px-6">
      <div class="mb-4 flex flex-wrap gap-2">
        <button
          v-for="p in PEOPLE"
          :key="p.slug"
          type="button"
          class="border-2 border-ink px-3 py-1 font-mono text-xs font-bold"
          :class="person === p.slug ? 'bg-red text-paper' : 'bg-paper hover:bg-paper-dim'"
          :aria-pressed="person === p.slug"
          @click="person = p.slug"
        >
          {{ p.name }}
        </button>
      </div>
      <div class="mb-8 flex flex-wrap gap-2">
        <button
          type="button"
          class="border-2 border-ink px-2.5 py-0.5 font-mono text-[0.7rem] font-bold"
          :class="category === '' ? 'bg-ink text-paper' : 'bg-paper hover:bg-paper-dim'"
          :aria-pressed="category === ''"
          @click="category = ''"
        >
          全部范畴
        </button>
        <button
          v-for="c in CATEGORIES"
          :key="c"
          type="button"
          class="border-2 border-ink px-2.5 py-0.5 font-mono text-[0.7rem] font-bold"
          :class="category === c ? 'bg-ink text-paper' : 'bg-paper hover:bg-paper-dim'"
          :aria-pressed="category === c"
          @click="category = c"
        >
          {{ c }}
        </button>
      </div>

      <AsyncBoundary
        :loading="loading"
        :error="error"
        :empty="visible.length ? '' : '没有符合条件的事件。'"
        :on-retry="reload"
      >
        <ol class="relative border-l-3 border-ink pl-6">
          <li v-for="e in visible" :key="e.id" class="relative mb-8 last:mb-0">
            <span class="absolute top-2 -left-[1.94rem] h-4 w-4 border-2 border-ink bg-red" aria-hidden="true" />
            <article class="brutal-card p-5">
              <div class="mb-2 flex flex-wrap items-center gap-2">
                <span class="border-2 border-ink bg-ink px-2 py-0.5 font-mono text-sm font-black text-paper">
                  {{ e.year }}
                </span>
                <span v-if="e.date_text" class="font-mono text-xs text-steel">{{ e.date_text }}</span>
                <span class="stamp" :class="catColor(e.category)">{{ e.category }}</span>
                <span v-if="e.location" class="ml-auto font-mono text-[0.7rem] text-steel">{{ e.location }}</span>
              </div>
              <h2 class="font-display text-lg leading-tight font-black">{{ e.title }}</h2>
              <p class="mt-2 text-sm leading-relaxed text-ink-soft">{{ e.body }}</p>
              <div v-if="e.people.length" class="mt-3 flex flex-wrap gap-2 border-t-2 border-ink pt-2.5">
                <RouterLink
                  v-for="p in e.people"
                  :key="p.slug"
                  :to="`/people/${p.slug}`"
                  class="border-2 border-ink bg-paper-dim px-2 py-0.5 font-mono text-[0.7rem] font-bold hover:bg-red hover:text-paper"
                >
                  {{ p.name }}
                </RouterLink>
              </div>
            </article>
          </li>
        </ol>
      </AsyncBoundary>
    </section>
  </div>
</template>
