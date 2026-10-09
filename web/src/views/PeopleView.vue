<script setup lang="ts">
import { useAsync } from '@/composables/useAsync'
import { getPeople } from '@/api'
import PageHeader from '@/components/PageHeader.vue'
import PersonCard from '@/components/PersonCard.vue'
import AsyncBoundary from '@/components/AsyncBoundary.vue'

const { data: people, loading, error, reload } = useAsync((signal) => getPeople(signal))
</script>

<template>
  <div>
    <PageHeader
      no="01"
      title="收录人物"
      sub="Figures"
      lead="本站收录马克思、恩格斯、列宁、毛泽东、切·格瓦拉五个人的完整档案。每人档案包含生平年表、简介、政治／经济／文化三域贡献、思想体系，以及如实的争议与评价。"
    />
    <section class="mx-auto max-w-6xl px-4 py-10 sm:px-6">
      <AsyncBoundary :loading="loading" :error="error" :on-retry="reload">
        <div class="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
          <PersonCard v-for="(p, i) in people ?? []" :key="p.slug" :person="p" :index="i" />
        </div>
      </AsyncBoundary>
    </section>
  </div>
</template>
