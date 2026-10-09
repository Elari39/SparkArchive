<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useAsync } from '@/composables/useAsync'
import { useDebouncedRef } from '@/composables/useDebouncedRef'
import { getTerms } from '@/api'
import PageHeader from '@/components/PageHeader.vue'
import AsyncBoundary from '@/components/AsyncBoundary.vue'

// 名称筛选加防抖，避免逐字符请求后端
const q = useDebouncedRef('')
const { data: terms, loading, error, reload } = useAsync((signal) => getTerms({ q: q.value }, signal), [q])
const visible = computed(() => terms.value ?? [])
</script>

<template>
  <div>
    <PageHeader
      no="04"
      title="术语表·概念卡"
      sub="Glossary"
      lead="核心概念的释义与互链。每张概念卡给出定义、详细阐释、相关人物、关联术语与出处著作。"
    />
    <section class="mx-auto max-w-6xl px-4 py-10 sm:px-6">
      <label class="mb-6 block max-w-md">
        <span class="sr-only">按名称筛选术语</span>
        <input
          v-model="q"
          type="search"
          placeholder="按名称筛选…"
          class="w-full border-3 border-ink bg-paper px-3 py-2 text-sm"
        />
      </label>

      <AsyncBoundary
        :loading="loading"
        :error="error"
        :empty="visible.length ? '' : '没有符合条件的术语。'"
        :on-retry="reload"
      >
        <div class="grid gap-5 md:grid-cols-2 lg:grid-cols-3">
          <RouterLink v-for="t in visible" :key="t.slug" :to="`/glossary/${t.slug}`" class="block">
            <article class="brutal-card brutal-interactive flex h-full flex-col p-5">
              <h2 class="font-display text-xl leading-tight font-black">{{ t.term }}</h2>
              <p v-if="t.aliases.length" class="mt-1 font-mono text-[0.68rem] text-steel">
                别名：{{ t.aliases.join('、') }}
              </p>
              <p class="mt-2.5 flex-1 text-sm leading-relaxed text-ink-soft">{{ t.definition }}</p>
            </article>
          </RouterLink>
        </div>
      </AsyncBoundary>
    </section>
  </div>
</template>
