<script setup lang="ts">
defineProps<{ label?: string }>()
</script>

<template>
  <div class="flex items-center gap-3 py-16" role="status" aria-live="polite">
    <span class="hatch-red h-8 w-8 border-3 border-ink" />
    <span class="font-mono text-sm tracking-widest uppercase">{{ label ?? '载入中' }}…</span>
  </div>
</template>
