<script setup lang="ts">
import { useAsync } from '@/composables/useAsync'
import { getPeople } from '@/api'
import SectionHeading from '@/components/SectionHeading.vue'
import PersonCard from '@/components/PersonCard.vue'
import AsyncBoundary from '@/components/AsyncBoundary.vue'

const people = useAsync((signal) => getPeople(signal))
</script>

<template>
  <!-- 人物 -->
  <section class="border-t-3 border-ink bg-paper-dim">
    <div class="mx-auto max-w-6xl px-4 py-12 sm:px-6">
      <SectionHeading no="01" title="收录人物" sub="Figures" />
      <AsyncBoundary
        :loading="people.loading.value"
        :error="people.error.value"
        :on-retry="people.reload"
      >
        <div class="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
          <PersonCard
            v-for="(p, i) in people.data.value ?? []"
            :key="p.slug"
            :person="p"
            :index="i"
          />
        </div>
      </AsyncBoundary>
    </div>
  </section>
</template>
