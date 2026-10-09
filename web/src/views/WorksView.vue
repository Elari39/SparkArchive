<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { useAsync } from '@/composables/useAsync'
import { useDebouncedRef } from '@/composables/useDebouncedRef'
import { getWorks } from '@/api'
import { PEOPLE_FILTERS } from '@/config/site'
import PageHeader from '@/components/PageHeader.vue'
import AsyncBoundary from '@/components/AsyncBoundary.vue'

const person = ref('')
// 标题筛选加防抖，避免逐字符请求后端
const q = useDebouncedRef('')
const { data: works, loading, error, reload } = useAsync(
  (signal) => getWorks({ person: person.value, q: q.value }, signal),
  [person, q],
)

const peopleOptions = PEOPLE_FILTERS

const visible = computed(() => works.value ?? [])
</script>

<template>
  <div>
    <PageHeader
      no="02"
      title="著作原文库"
      sub="Works"
      lead="收录经典文献的精选原文摘录，每段摘录均标注出处，并附完整书目与外部阅读链接。受版权保护的作品不收录全文。"
    />
    <section class="mx-auto max-w-6xl px-4 py-10 sm:px-6">
      <div class="mb-6 flex flex-col gap-3 sm:flex-row">
        <div class="flex flex-wrap gap-2">
          <button
            v-for="o in peopleOptions"
            :key="o.slug"
            type="button"
            class="border-2 border-ink px-3 py-1 font-mono text-xs font-bold transition-colors"
            :class="person === o.slug ? 'bg-red text-paper' : 'bg-paper hover:bg-paper-dim'"
            :aria-pressed="person === o.slug"
            @click="person = o.slug"
          >
            {{ o.name }}
          </button>
        </div>
        <label class="sm:ml-auto">
          <span class="sr-only">按标题筛选著作</span>
          <input
            v-model="q"
            type="search"
            placeholder="按标题筛选…"
            class="w-full border-3 border-ink bg-paper px-3 py-1.5 text-sm sm:w-64"
          />
        </label>
      </div>

      <AsyncBoundary
        :loading="loading"
        :error="error"
        :empty="visible.length ? '' : '没有符合条件的著作。'"
        :on-retry="reload"
      >
        <div class="grid gap-5 md:grid-cols-2">
          <RouterLink v-for="w in visible" :key="w.slug" :to="`/works/${w.slug}`" class="block">
            <article class="brutal-card brutal-interactive flex h-full flex-col p-5">
              <div class="mb-2 flex flex-wrap items-center gap-2">
                <span class="stamp stamp-ink">{{ w.year }}</span>
                <span class="stamp">{{ w.category }}</span>
              </div>
              <h2 class="font-display text-xl leading-tight font-black">{{ w.title }}</h2>
              <p class="mt-1 font-mono text-[0.7rem] tracking-wider text-steel">{{ w.person_name }}</p>
              <p class="mt-2.5 flex-1 text-sm leading-relaxed text-ink-soft">{{ w.summary }}</p>
            </article>
          </RouterLink>
        </div>
      </AsyncBoundary>
    </section>
  </div>
</template>
