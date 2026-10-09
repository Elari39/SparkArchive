<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const q = ref('')

function submit() {
  const v = q.value.trim()
  if (!v) return
  void router.push({ name: 'search', query: { q: v } })
}
</script>

<template>
  <form class="flex w-full max-w-xl" role="search" @submit.prevent="submit">
    <label for="site-search" class="sr-only">检索档案馆</label>
    <input
      id="site-search"
      v-model="q"
      type="search"
      placeholder="检索人物、著作、事件、术语…"
      class="min-w-0 flex-1 border-3 border-ink bg-paper px-3 py-2 text-sm font-medium outline-none placeholder:text-steel-light"
    />
    <button
      type="submit"
      class="shrink-0 border-3 border-l-0 border-ink bg-red px-4 py-2 font-mono text-xs font-bold tracking-widest text-paper uppercase hover:bg-red-deep"
    >
      检索
    </button>
  </form>
</template>
