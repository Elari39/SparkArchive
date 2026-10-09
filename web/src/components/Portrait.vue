<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { getPortrait } from '@/config/portraits'
import PortraitArt from './portrait-art/PortraitArt.vue'

/**
 * 人物肖像：优先使用**公有领域历史照片**（见 `config/portraits.ts` 与《凡例》「图片来源」）。
 * 映射缺失或图片加载失败时，自动回退为构成主义矢量插画（`PortraitArt.vue`）。
 */
const props = withDefaults(
  defineProps<{
    /** 人物键（= 后端 portrait_key / slug） */
    person: string
    /** 边长（像素） */
    size?: number
    /** 人物中文名，用于图片 alt；缺省时回退映射表内的名称 */
    name?: string
    /** 图片加载策略：卡片缩略图用 lazy，详情页首图用 eager */
    loading?: 'lazy' | 'eager'
    /** 填充模式：不锁定像素尺寸，撑满父容器（用于响应式照片墙） */
    fill?: boolean
  }>(),
  { size: 160, loading: 'lazy', fill: false },
)

const isFill = computed(() => props.fill)

const failed = ref(false)
const meta = computed(() => getPortrait(props.person))
const showPhoto = computed(() => !!meta.value && !failed.value)
const displayName = computed(() => props.name ?? meta.value?.name ?? props.person)
// 始终为字符串，避免模板中出现 TS 专属语法
const photoWebp = computed(() => meta.value?.webp ?? '')
const photoJpg = computed(() => meta.value?.jpg ?? '')

// 切换人物时重置加载失败状态
watch(
  () => props.person,
  () => {
    failed.value = false
  },
)
</script>

<template>
  <!-- 真实历史照片（WebP 主格式 + JPEG 回退） -->
  <picture
    v-if="showPhoto"
    class="block"
    :class="isFill ? 'h-full w-full' : ''"
    :style="isFill ? undefined : { width: `${size}px`, height: `${size}px` }"
    :data-portrait="person"
  >
    <source :srcset="photoWebp" type="image/webp" />
    <img
      :src="photoJpg"
      :width="size"
      :height="size"
      :alt="`${displayName} 历史照片`"
      :loading="loading"
      :fetchpriority="loading === 'eager' ? 'high' : undefined"
      decoding="async"
      class="portrait-photo block h-full w-full object-cover"
      @error="failed = true"
    />
  </picture>

  <!-- 构成主义矢量插画（降级） -->
  <PortraitArt
    v-else
    :variant="person"
    :name="displayName"
    :fill="isFill"
    :size="size"
  />
</template>
