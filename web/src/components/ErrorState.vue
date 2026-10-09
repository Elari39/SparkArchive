<script setup lang="ts">
defineProps<{ message: string }>()
</script>

<template>
  <div class="brutal-card border-l-[10px] border-l-red p-6" role="alert">
    <p class="mb-1 font-mono text-xs tracking-widest text-red-deep uppercase">ERROR</p>
    <p class="font-bold">{{ message }}</p>
    <slot />
  </div>
</template>
