<script setup lang="ts">
import { RouterLink } from 'vue-router'
import PageHeader from '@/components/PageHeader.vue'
import { PORTRAITS } from '@/config/portraits'

const RULES = [
  {
    no: '01',
    title: '内容体例',
    body: '每个人物档案固定包含：生平年表、简介、政治贡献、经济贡献、文化贡献、思想体系、争议与评价七个部分。四域贡献与思想体系采用统一的编号结构，便于横向对照。',
  },
  {
    no: '02',
    title: '史料来源',
    body: '史实、时间与引文均尽量标注出处。著作摘录标注所据版本与章节，并附外部阅读链接。本站仅收录摘录用于学习研究，不收录受版权保护的全文。',
  },
  {
    no: '03',
    title: '争议标注规则',
    body: '对存在学术争论或历史评价分歧的内容，一律设「争议与评价」专节如实说明，包括对失误与错误的记载。这既是史实要求，也是档案可信度的来源。',
  },
  {
    no: '04',
    title: '版权说明',
    body: '本站原创文字采用 CC BY-NC-SA 4.0 许可。所引原著摘录版权归各自权利人所有。毛泽东著作在中国境内的版权保护期至 2026 年底，故本站仅收录摘录与书目索引。',
  },
  {
    no: '05',
    title: '肖像说明',
    body: '本站人物肖像取自公开的历史照片与史料图片，逐张标注来源、摄影者与许可状态（见下方「图片来源」）。其中维基共享资源来源的 7 张为公有领域史料，百度百科来源的 8 张原图未标注许可、依站点运营方的声明使用，两者许可状态不同，不作混同。照片仅用于学习与研究，不作任何商业用途；若某张照片缺失或无法加载，页面将回退为本站绘制的构成主义矢量插画。',
  },
  {
    no: '06',
    title: '技术说明',
    body: '前端 Vue 3 + Tailwind CSS 4，后端 Go + SQLite（FTS5 全文索引），Docker 一键部署，服务端口 12026。',
  },
]

/** 人物肖像来源清单，由 config/portraits.ts 渲染 */
const portraitList = Object.entries(PORTRAITS).map(([key, meta]) => ({ key, meta }))
</script>

<template>
  <div>
    <PageHeader
      no="07"
      title="凡例"
      sub="About"
      lead="本页说明档案馆的收录范围、内容体例、史料标注规则与版权政策。"
    />
    <section class="mx-auto max-w-4xl px-4 py-10 sm:px-6">
      <div class="space-y-5">
        <article v-for="r in RULES" :key="r.no" class="brutal-card p-5">
          <div class="mb-2 flex items-center gap-3">
            <span class="border-3 border-ink bg-ink px-2 py-0.5 font-mono text-sm font-black text-paper">
              {{ r.no }}
            </span>
            <h2 class="font-display text-xl font-black">{{ r.title }}</h2>
          </div>
          <p class="text-sm leading-relaxed text-ink-soft">{{ r.body }}</p>
        </article>
      </div>

      <div class="mt-10 border-t-3 border-ink pt-6">
        <h2 class="mb-3 font-display text-xl font-black">收录范围</h2>
        <ul class="grid gap-2 text-sm sm:grid-cols-2">
          <li><RouterLink to="/people" class="font-bold text-red-deep hover:underline">人物</RouterLink> — 马克思、恩格斯、普列汉诺夫、蔡特金、片山潜、列宁、卢森堡、柯伦泰、李大钊、胡志明、葛兰西、毛泽东、卡斯特罗、切·格瓦拉、桑卡拉</li>
          <li><RouterLink to="/works" class="font-bold text-red-deep hover:underline">著作</RouterLink> — 经典文献摘录与书目索引</li>
          <li><RouterLink to="/timeline" class="font-bold text-red-deep hover:underline">年表</RouterLink> — 1848—2016 国际共运大事</li>
          <li><RouterLink to="/glossary" class="font-bold text-red-deep hover:underline">术语</RouterLink> — 核心概念卡与互链</li>
          <li><RouterLink to="/graph" class="font-bold text-red-deep hover:underline">关系</RouterLink> — 思想传承与国际联系</li>
          <li><RouterLink to="/search" class="font-bold text-red-deep hover:underline">检索</RouterLink> — 全文跨类检索</li>
        </ul>
      </div>

      <div class="mt-10 border-t-3 border-ink pt-6">
        <h2 class="mb-3 font-display text-xl font-black">图片来源</h2>
        <p class="mb-4 text-sm leading-relaxed text-ink-soft">
          本页列出全部 15 张人物肖像的来源与许可状态。其中 7 张取自维基共享资源，为公有领域史料；
          另 8 张取自百度百科词条摘要图，原图未标注许可，依站点运营方的声明使用。逐张署名如下：
        </p>
        <ul class="space-y-4">
          <li v-for="p in portraitList" :key="p.key" class="border-l-[6px] border-red pl-3">
            <p class="font-display text-base font-black">
              {{ p.meta.name }}
              <span class="ml-1 font-mono text-xs font-normal text-steel">{{ p.meta.credit }}</span>
            </p>
            <p class="mt-0.5 text-xs leading-relaxed text-ink-soft">
              {{ p.meta.license }} ·
              <a
                :href="p.meta.source"
                target="_blank"
                rel="noopener noreferrer"
                class="font-bold text-red-deep underline underline-offset-2"
                >原始文件页</a
              >
            </p>
            <p v-if="p.meta.note" class="mt-0.5 text-xs leading-relaxed text-steel">{{ p.meta.note }}</p>
          </li>
        </ul>
      </div>

      <blockquote class="mt-10 border-l-[10px] border-red bg-ink p-6 font-display text-xl leading-snug font-black text-paper">
        「全世界无产者，联合起来！」
      </blockquote>
    </section>
  </div>
</template>
