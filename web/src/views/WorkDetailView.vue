<script setup lang="ts">
import { toRef } from 'vue'
import { RouterLink } from 'vue-router'
import { useAsync } from '@/composables/useAsync'
import { getWork } from '@/api'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import AsyncBoundary from '@/components/AsyncBoundary.vue'

const props = defineProps<{ slug: string }>()
const slug = toRef(props, 'slug')
const { data: work, loading, error, reload } = useAsync((signal) => getWork(slug.value, signal), [slug])
</script>

<template>
  <div>
    <AsyncBoundary
      :loading="loading"
      :error="error"
      :empty="!work ? '未找到该著作' : ''"
      :on-retry="reload"
      state-class="mx-auto max-w-4xl px-4"
      empty-class="mx-auto my-10 max-w-2xl"
    >

    <template v-if="work">
      <header class="border-b-3 border-ink bg-paper-dim">
        <div class="mx-auto max-w-4xl px-4 py-10 sm:px-6">
          <Breadcrumbs :items="[{ to: '/works', label: '著作' }, { label: work.title }]" />
          <div class="mb-3 flex flex-wrap gap-2">
            <span class="stamp stamp-ink">{{ work.year }}</span>
            <span class="stamp">{{ work.category }}</span>
            <span class="stamp stamp-red">{{ work.person_name }}</span>
          </div>
          <h1 class="display-xl">《{{ work.title }}》</h1>
          <p class="mt-4 border-l-[6px] border-red pl-4 leading-relaxed">{{ work.summary }}</p>
        </div>
      </header>

      <div class="mx-auto max-w-4xl px-4 py-10 sm:px-6">
        <h2 class="mb-5 border-b-3 border-ink pb-2 text-2xl font-black">原文摘录</h2>
        <div class="space-y-5">
          <figure v-for="(e, i) in work.excerpts" :key="i" class="brutal-card border-l-[10px] border-l-red p-5">
            <span class="mb-2 block font-mono text-[0.65rem] tracking-widest text-steel">
              摘录 {{ String(i + 1).padStart(2, '0') }}
            </span>
            <blockquote class="font-display text-lg leading-relaxed font-bold">{{ e.text }}</blockquote>
            <figcaption class="mt-3 border-t-2 border-ink pt-2 text-xs text-steel">{{ e.note }}</figcaption>
          </figure>
        </div>

        <div class="brutal-card mt-8 p-5">
          <h2 class="mb-3 font-mono text-xs font-bold tracking-widest uppercase">书目与出处</h2>
          <dl class="space-y-2 text-sm">
            <div class="flex flex-col gap-0.5 sm:flex-row sm:gap-3">
              <dt class="shrink-0 font-mono text-xs text-steel sm:w-20">版本</dt>
              <dd class="font-medium">{{ work.source }}</dd>
            </div>
            <div v-if="work.source_url" class="flex flex-col gap-0.5 sm:flex-row sm:gap-3">
              <dt class="shrink-0 font-mono text-xs text-steel sm:w-20">外部阅读</dt>
              <dd>
                <a
                  :href="work.source_url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="font-bold text-red-deep underline underline-offset-4"
                  >前往阅读全文 →</a
                >
              </dd>
            </div>
          </dl>
          <p class="mt-4 border-t-2 border-ink pt-3 text-xs leading-relaxed text-steel">
            本站仅收录摘录用于学习研究，全部著作权归原作者及权利人所有。毛泽东著作在中国境内的版权保护期至 2026
            年底，故本站不收录其全文。
          </p>
        </div>

        <p class="mt-8">
          <RouterLink to="/works" class="font-mono text-sm font-bold hover:text-red-deep">← 返回著作库</RouterLink>
        </p>
      </div>
    </template>
    </AsyncBoundary>
  </div>
</template>
