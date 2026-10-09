<script setup lang="ts">
import { computed } from 'vue'
import { marked } from 'marked'
import { sanitizeHtml } from '@/utils/sanitize'

const props = defineProps<{ source: string }>()

marked.setOptions({ gfm: true, breaks: false })

// marked 输出经白名单净化后再交给 v-html，避免上游内容引入 XSS。
const html = computed(() => sanitizeHtml(marked.parse(props.source, { async: false }) as string))
</script>

<template>
  <div class="prose-archive" v-html="html" />
</template>
