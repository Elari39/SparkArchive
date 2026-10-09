<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import { NAV, SITE } from '@/config/site'

const open = ref(false)
</script>

<template>
  <header class="sticky top-0 z-50 border-b-3 border-ink bg-paper">
    <div class="mx-auto flex max-w-6xl items-center gap-3 px-4 py-3 sm:px-6">
      <RouterLink to="/" class="group flex items-center gap-2.5" aria-label="返回首页">
        <span class="flex h-9 w-9 shrink-0 items-center justify-center border-3 border-ink bg-red">
          <span class="font-display text-lg font-black text-paper">星</span>
        </span>
        <span class="leading-none">
          <span class="block font-display text-lg font-black tracking-tight">{{ SITE.name }}</span>
          <span class="block font-mono text-[0.6rem] tracking-[0.18em] text-steel">{{ SITE.nameEn }}</span>
        </span>
      </RouterLink>

      <nav class="ml-auto hidden items-center gap-1 lg:flex" aria-label="主导航">
        <RouterLink
          v-for="item in NAV"
          :key="item.to"
          :to="item.to"
          class="border-2 border-transparent px-2.5 py-1.5 text-sm font-bold transition-colors hover:border-ink hover:bg-ink hover:text-paper"
          active-class="border-ink bg-red text-paper"
        >
          {{ item.label }}
        </RouterLink>
      </nav>

      <button
        type="button"
        class="ml-auto border-3 border-ink bg-paper px-3 py-1.5 font-mono text-xs font-bold lg:hidden"
        :aria-expanded="open"
        aria-controls="mobile-nav"
        @click="open = !open"
      >
        {{ open ? '关闭' : '菜单' }}
      </button>
    </div>

    <nav v-if="open" id="mobile-nav" class="border-t-3 border-ink bg-paper-dim lg:hidden" aria-label="移动导航">
      <div class="mx-auto grid max-w-6xl grid-cols-2 gap-2 px-4 py-3 sm:px-6">
        <RouterLink
          v-for="item in NAV"
          :key="item.to"
          :to="item.to"
          class="border-2 border-ink bg-paper px-3 py-2 text-center text-sm font-bold"
          active-class="bg-red text-paper"
          @click="open = false"
        >
          {{ item.label }}
        </RouterLink>
      </div>
    </nav>
  </header>
</template>
