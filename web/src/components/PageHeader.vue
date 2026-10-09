<script setup lang="ts">
defineProps<{ no: string; title: string; sub?: string; lead?: string }>()
</script>

<template>
  <header class="border-b-3 border-ink bg-paper">
    <div class="mx-auto max-w-6xl px-4 py-10 sm:px-6 sm:py-14">
      <div class="mb-3 flex flex-wrap items-center gap-2">
        <span class="stamp stamp-ink">{{ no }}</span>
        <span v-if="sub" class="stamp">{{ sub }}</span>
      </div>
      <h1 class="display-hero">{{ title }}</h1>
      <p v-if="lead" class="mt-5 max-w-3xl border-l-[6px] border-red pl-4 text-base leading-relaxed text-ink-soft">
        {{ lead }}
      </p>
    </div>
  </header>
</template>
