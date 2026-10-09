<script setup lang="ts">
import { RouterLink } from 'vue-router'

/**
 * 面包屑导航。最后一项为当前页（纯文本），其余项渲染为 RouterLink。
 * 收敛 PersonDetail / WorkDetail / TermDetail 三处重复结构。
 */
defineProps<{
  items: { to?: string; label: string }[]
}>()
</script>

<template>
  <nav class="mb-5 font-mono text-[0.7rem] tracking-wider text-steel" aria-label="面包屑">
    <template v-for="(it, i) in items" :key="i">
      <span v-if="i > 0" class="mx-1.5">/</span>
      <RouterLink v-if="it.to" :to="it.to" class="hover:text-red-deep">{{ it.label }}</RouterLink>
      <span v-else>{{ it.label }}</span>
    </template>
  </nav>
</template>
