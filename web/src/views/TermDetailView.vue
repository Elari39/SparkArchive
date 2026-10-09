<script setup lang="ts">
import { toRef } from 'vue'
import { RouterLink } from 'vue-router'
import { useAsync } from '@/composables/useAsync'
import { getTerm } from '@/api'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import AsyncBoundary from '@/components/AsyncBoundary.vue'

const props = defineProps<{ slug: string }>()
const slug = toRef(props, 'slug')
const { data: term, loading, error, reload } = useAsync((signal) => getTerm(slug.value, signal), [slug])
</script>

<template>
  <div>
    <AsyncBoundary
      :loading="loading"
      :error="error"
      :empty="!term ? '未找到该术语' : ''"
      :on-retry="reload"
      state-class="mx-auto max-w-3xl px-4"
      empty-class="mx-auto my-10 max-w-2xl"
    >

    <template v-if="term">
      <header class="border-b-3 border-ink bg-paper-dim">
        <div class="mx-auto max-w-3xl px-4 py-10 sm:px-6">
          <Breadcrumbs :items="[{ to: '/glossary', label: '术语' }, { label: term.term }]" />
          <h1 class="display-hero">{{ term.term }}</h1>
          <p v-if="term.aliases.length" class="mt-3 font-mono text-xs tracking-wider text-steel">
            别名 / {{ term.aliases.join(' · ') }}
          </p>
          <p class="mt-5 border-l-[6px] border-red pl-4 text-base leading-relaxed font-medium">
            {{ term.definition }}
          </p>
        </div>
      </header>

      <div class="mx-auto max-w-3xl px-4 py-10 sm:px-6">
        <p class="prose-archive">{{ term.body }}</p>

        <div class="mt-10 grid gap-5 sm:grid-cols-3">
          <div v-if="term.people.length" class="brutal-card p-4">
            <h2 class="mb-3 border-b-2 border-ink pb-1.5 font-mono text-xs font-bold tracking-widest uppercase">
              相关人物
            </h2>
            <ul class="space-y-1.5">
              <li v-for="p in term.people" :key="p.slug">
                <RouterLink :to="`/people/${p.slug}`" class="text-sm font-bold hover:text-red-deep">
                  {{ p.name }}
                </RouterLink>
              </li>
            </ul>
          </div>

          <div v-if="term.related.length" class="brutal-card p-4">
            <h2 class="mb-3 border-b-2 border-ink pb-1.5 font-mono text-xs font-bold tracking-widest uppercase">
              关联术语
            </h2>
            <ul class="space-y-1.5">
              <li v-for="r in term.related" :key="r.slug">
                <RouterLink :to="`/glossary/${r.slug}`" class="text-sm font-bold hover:text-red-deep">
                  {{ r.term }}
                </RouterLink>
              </li>
            </ul>
          </div>

          <div v-if="term.sources.length" class="brutal-card p-4">
            <h2 class="mb-3 border-b-2 border-ink pb-1.5 font-mono text-xs font-bold tracking-widest uppercase">
              出处著作
            </h2>
            <ul class="space-y-1.5">
              <li v-for="s in term.sources" :key="s.slug">
                <RouterLink :to="`/works/${s.slug}`" class="text-sm font-bold hover:text-red-deep">
                  {{ s.title }}
                </RouterLink>
              </li>
            </ul>
          </div>
        </div>

        <p class="mt-8">
          <RouterLink to="/glossary" class="font-mono text-sm font-bold hover:text-red-deep">← 返回术语表</RouterLink>
        </p>
      </div>
    </template>
    </AsyncBoundary>
  </div>
</template>
