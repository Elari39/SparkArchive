<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'
import { useSiteStore } from '@/stores/site'

const site = useSiteStore()
onMounted(() => void site.load())
</script>

<template>
  <div class="flex min-h-screen flex-col bg-paper">
    <a
      href="#main"
      class="sr-only focus:not-sr-only focus:absolute focus:top-2 focus:left-2 focus:z-[60] focus:border-3 focus:border-ink focus:bg-brass focus:px-3 focus:py-2 focus:font-bold"
      >跳到主内容</a
    >
    <AppHeader />
    <main id="main" class="flex-1">
      <RouterView />
    </main>
    <AppFooter />
  </div>
</template>
