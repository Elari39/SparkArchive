<script setup lang="ts">
import LoadingState from '@/components/LoadingState.vue'
import ErrorState from '@/components/ErrorState.vue'

/**
 * 列表/详情页的统一状态骨架：loading → error → empty → 内容。
 *
 * 收敛原先在 7 个视图里重复的「LoadingState / ErrorState / 空态文案」三段式；
 * 调用方只需传状态与空态文案，内容通过默认插槽渲染。
 */
withDefaults(
  defineProps<{
    loading?: boolean
    error?: string
    /** 数据为空时展示的文案；不传则不渲染空态分支 */
    empty?: string
    /** 重试回调：传入则在错误态显示「重试」按钮 */
    onRetry?: () => void
    /** 空态容器附加类名（不同页面宽度不同） */
    emptyClass?: string
    /** 错误 / 加载态容器附加类名 */
    stateClass?: string
  }>(),
  { loading: false, error: '', empty: '', onRetry: undefined, emptyClass: '', stateClass: '' },
)
</script>

<template>
  <LoadingState v-if="loading" :class="stateClass" />
  <ErrorState v-else-if="error" :message="error" :class="stateClass">
    <button
      v-if="onRetry"
      type="button"
      class="mt-3 border-2 border-ink bg-paper px-3 py-1 font-mono text-xs font-bold hover:bg-paper-dim"
      @click="onRetry()"
    >
      重试
    </button>
  </ErrorState>
  <p
    v-else-if="empty"
    class="brutal-card border-l-[10px] border-l-steel p-5 font-bold"
    :class="emptyClass"
  >
    {{ empty }}
  </p>
  <slot v-else />
</template>
