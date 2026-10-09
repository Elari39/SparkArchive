<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<template>
  <section class="mx-auto max-w-3xl px-4 py-20 text-center sm:px-6">
    <p class="font-mono text-8xl font-black text-red">404</p>
    <h1 class="mt-4 display-xl">档案未找到</h1>
    <p class="mt-4 text-ink-soft">你访问的档案条目不存在，或已被移动。</p>
    <div class="mt-8 flex flex-wrap justify-center gap-3">
      <RouterLink
        to="/"
        class="border-3 border-ink bg-red px-5 py-2 font-mono text-xs font-bold tracking-widest text-paper uppercase"
      >
        返回首页
      </RouterLink>
      <RouterLink
        to="/search"
        class="border-3 border-ink bg-paper px-5 py-2 font-mono text-xs font-bold tracking-widest uppercase"
      >
        去检索
      </RouterLink>
    </div>
  </section>
</template>
