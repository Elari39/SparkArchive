<script setup lang="ts">
import { RouterLink } from 'vue-router'
import SectionHeading from '@/components/SectionHeading.vue'

const MODULES = [
  { to: '/people', no: '01', title: '人物', desc: '九位革命者的生平、简介与政治经济文化贡献' },
  { to: '/works', no: '02', title: '著作', desc: '经典文献的原文摘录与书目索引' },
  { to: '/timeline', no: '03', title: '年表', desc: '国际共产主义运动大事年表' },
  { to: '/glossary', no: '04', title: '术语', desc: '核心概念卡与互链释义' },
  { to: '/graph', no: '05', title: '关系', desc: '思想传承与国际联系图谱' },
]
</script>

<template>
  <!-- 模块导航 -->
  <section class="mx-auto max-w-6xl px-4 py-12 sm:px-6">
    <SectionHeading no="00" title="档案分区" sub="Sections" />
    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
      <RouterLink v-for="m in MODULES" :key="m.to" :to="m.to" class="block">
        <div class="brutal-card brutal-interactive h-full p-4">
          <span class="font-mono text-2xl font-black text-red">{{ m.no }}</span>
          <h3 class="mt-1 font-display text-xl font-black">{{ m.title }}</h3>
          <p class="mt-1.5 text-[0.8rem] leading-relaxed text-ink-soft">{{ m.desc }}</p>
        </div>
      </RouterLink>
    </div>
  </section>
</template>
