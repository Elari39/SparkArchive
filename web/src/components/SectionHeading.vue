<script setup lang="ts">
defineProps<{ no?: string; title: string; sub?: string }>()
</script>

<template>
  <div class="mb-6 flex items-end gap-4">
    <span
      v-if="no"
      class="shrink-0 border-3 border-ink bg-ink px-2 py-1 font-mono text-lg font-extrabold text-paper"
      >{{ no }}</span
    >
    <div class="min-w-0">
      <h2 class="display-xl uppercase">{{ title }}</h2>
      <p v-if="sub" class="mt-1 font-mono text-xs tracking-widest text-steel uppercase">{{ sub }}</p>
    </div>
    <div class="mb-2 hidden h-1 flex-1 bg-ink sm:block" />
  </div>
</template>
