<script setup lang="ts">
import { RouterLink } from 'vue-router'
import Portrait from './Portrait.vue'
import type { PersonSummary } from '@/api/types'

defineProps<{ person: PersonSummary; index: number }>()
</script>

<template>
  <RouterLink :to="`/people/${person.slug}`" class="group block">
    <article class="brutal-card brutal-interactive flex h-full flex-col">
      <div class="flex items-stretch border-b-3 border-ink">
        <Portrait
          :person="person.portrait_key"
          :size="96"
          :name="person.name"
          loading="lazy"
          class="shrink-0 border-r-3 border-ink"
        />
        <div class="flex min-w-0 flex-1 flex-col justify-between p-3">
          <span class="font-mono text-[0.65rem] tracking-widest text-steel">
            {{ String(index + 1).padStart(2, '0') }} / {{ person.nationality }}
          </span>
          <div>
            <h3 class="truncate font-display text-lg leading-tight font-black">{{ person.name }}</h3>
            <p class="truncate font-mono text-[0.65rem] tracking-wider text-steel uppercase">
              {{ person.name_en }}
            </p>
          </div>
        </div>
      </div>
      <div class="flex flex-1 flex-col p-3">
        <p class="font-mono text-[0.68rem] text-steel">
          {{ person.birth.slice(0, 4) }} — {{ person.death.slice(0, 4) }}
        </p>
        <p class="mt-1 text-sm font-bold text-red-deep">{{ person.epithet }}</p>
        <p class="mt-2 line-clamp-3 text-[0.82rem] leading-relaxed text-ink-soft">{{ person.summary }}</p>
      </div>
    </article>
  </RouterLink>
</template>
