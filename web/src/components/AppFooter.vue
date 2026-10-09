<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { NAV, SITE } from '@/config/site'
import { useSiteCopy } from '@/composables/useSiteCopy'

const { name, nameEn, tagline, footerNote, licenseNote } = useSiteCopy()
</script>

<template>
  <footer class="mt-16 border-t-3 border-ink bg-ink text-paper">
    <div class="mx-auto max-w-6xl px-4 py-10 sm:px-6">
      <div class="grid gap-8 md:grid-cols-[2fr_1fr]">
        <div>
          <p class="font-display text-2xl font-black">{{ name }} · {{ nameEn }}</p>
          <p class="mt-1 font-mono text-xs tracking-[0.18em] text-brass">{{ tagline }}</p>
          <p class="mt-4 max-w-xl text-sm leading-relaxed text-paper/75">
            {{ footerNote }}
          </p>
        </div>
        <nav class="grid grid-cols-2 gap-x-4 gap-y-1.5 self-start font-mono text-xs" aria-label="页脚导航">
          <RouterLink v-for="item in NAV" :key="item.to" :to="item.to" class="text-paper/70 hover:text-brass">
            / {{ item.label }}
          </RouterLink>
        </nav>
      </div>
      <div
        class="mt-8 flex flex-col gap-2 border-t-2 border-paper/25 pt-5 font-mono text-[0.7rem] text-paper/55 sm:flex-row sm:justify-between"
      >
        <span>{{ SITE.archiveNo }} · {{ SITE.est }}</span>
        <span>{{ licenseNote }}</span>
      </div>
    </div>
  </footer>
</template>
