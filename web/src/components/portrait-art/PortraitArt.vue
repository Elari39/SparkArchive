<script setup lang="ts">
import { computed } from 'vue'

/**
 * 构成主义矢量肖像插画（Portrait.vue 的降级分支）。
 * 各人物图形与配色规格集中在此，便于独立调整而无需改动加载逻辑。
 */
const props = defineProps<{
  /** 人物键：marx / engels / plekhanov / zetkin / katayama / lenin / luxemburg / kollontai /
   *  li-dazhao / ho-chi-minh / gramsci / mao / castro / guevara / sankara（未知键回退 marx） */
  variant: string
  /** 显示名，用于 aria-label */
  name: string
  /** 是否撑满父容器（响应式照片墙） */
  fill?: boolean
  /** 边长（像素），非 fill 模式生效 */
  size?: number
}>()

type Variant =
  | 'marx'
  | 'engels'
  | 'plekhanov'
  | 'zetkin'
  | 'katayama'
  | 'lenin'
  | 'luxemburg'
  | 'kollontai'
  | 'li-dazhao'
  | 'ho-chi-minh'
  | 'gramsci'
  | 'mao'
  | 'castro'
  | 'guevara'
  | 'sankara'
type Spec = { bg: string; accent: string; coat: string }

const PAPER = '#f2ede4'
const INK = '#111111'

const SPECS: Record<Variant, Spec> = {
  marx: { bg: '#d62828', accent: PAPER, coat: '#1f1f1f' },
  engels: { bg: '#4a4a4a', accent: '#f4b400', coat: '#1f1f1f' },
  plekhanov: { bg: '#4a4a4a', accent: '#f4b400', coat: '#1f1f1f' },
  zetkin: { bg: '#4a4a4a', accent: '#d62828', coat: '#1f1f1f' },
  katayama: { bg: '#111111', accent: '#f4b400', coat: '#4a4a4a' },
  lenin: { bg: '#a4161a', accent: '#f4b400', coat: '#1f1f1f' },
  luxemburg: { bg: '#a4161a', accent: PAPER, coat: '#1f1f1f' },
  kollontai: { bg: '#a4161a', accent: PAPER, coat: '#1f1f1f' },
  'li-dazhao': { bg: '#a4161a', accent: '#f4b400', coat: '#1f1f1f' },
  'ho-chi-minh': { bg: '#111111', accent: '#f4b400', coat: '#4a4a4a' },
  gramsci: { bg: '#4a4a4a', accent: '#d62828', coat: '#1f1f1f' },
  mao: { bg: '#d62828', accent: '#f4b400', coat: '#4a4a4a' },
  castro: { bg: '#1f1f1f', accent: '#d62828', coat: '#4a4a4a' },
  guevara: { bg: '#111111', accent: '#d62828', coat: '#4a4a4a' },
  sankara: { bg: '#d62828', accent: '#f4b400', coat: '#1f1f1f' },
}

const variant = computed<Variant>(() => (props.variant as Variant) in SPECS ? (props.variant as Variant) : 'marx')
const spec = computed<Spec>(() => SPECS[variant.value])
const isFill = computed(() => props.fill === true)
</script>

<template>
  <svg
    :width="isFill ? '100%' : (props.size ?? 160)"
    :height="isFill ? '100%' : (props.size ?? 160)"
    viewBox="0 0 120 120"
    role="img"
    :aria-label="`${name} 的构成主义风格插画肖像`"
    class="block"
    :class="isFill ? 'h-full w-full' : ''"
    :data-portrait="variant"
  >
    <rect width="120" height="120" :fill="spec.bg" />
    <!-- 构成主义斜切装饰 -->
    <polygon points="0,0 120,0 120,34 0,64" :fill="spec.accent" opacity="0.20" />
    <polygon points="0,120 120,120 120,92 0,112" :fill="PAPER" opacity="0.10" />

    <!-- ============ 马克思：浓密卷发 + 大络腮胡 ============ -->
    <g v-if="variant === 'marx'">
      <!-- 肩与外套 -->
      <path d="M12 120 Q14 92 40 86 L60 96 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <path d="M60 96 L48 120 L60 120 L72 120 Z" :fill="PAPER" opacity="0.9" />
      <!-- 颈部 -->
      <rect x="51" y="74" width="18" height="18" :fill="PAPER" />
      <!-- 头发（蓬松外扩） -->
      <path d="M28 52 Q22 20 60 18 Q98 20 92 52 Q90 34 78 30 Q60 26 42 30 Q30 34 28 52 Z" :fill="INK" />
      <path d="M28 52 Q24 40 30 34 Q28 48 34 56 Z" :fill="INK" />
      <path d="M92 52 Q96 40 90 34 Q92 48 86 56 Z" :fill="INK" />
      <!-- 面部 -->
      <path d="M38 46 Q38 30 60 30 Q82 30 82 46 L82 60 Q82 78 60 78 Q38 78 38 60 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <!-- 大络腮胡（标志性） -->
      <path d="M34 52 Q34 84 60 86 Q86 84 86 52 Q84 66 76 70 Q68 62 60 62 Q52 62 44 70 Q36 66 34 52 Z" :fill="INK" />
      <path d="M44 78 Q60 90 76 78 Q68 88 60 88 Q52 88 44 78 Z" :fill="INK" />
      <!-- 眉与眼 -->
      <rect x="45" y="46" width="12" height="3.5" :fill="INK" />
      <rect x="63" y="46" width="12" height="3.5" :fill="INK" />
      <circle cx="51" cy="53" r="2.6" :fill="INK" />
      <circle cx="69" cy="53" r="2.6" :fill="INK" />
      <!-- 鼻 -->
      <path d="M60 54 L60 63 L55 63" fill="none" :stroke="INK" stroke-width="2" />
    </g>

    <!-- ============ 恩格斯：高额头 + 络腮胡 + 领结 ============ -->
    <g v-else-if="variant === 'engels'">
      <path d="M12 120 Q14 92 40 86 L60 96 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <!-- 衬衫领与领结 -->
      <path d="M52 92 L60 104 L68 92 L60 96 Z" :fill="PAPER" />
      <path d="M56 100 L52 108 L60 106 Z" :fill="spec.accent" />
      <path d="M64 100 L68 108 L60 106 Z" :fill="spec.accent" />
      <rect x="51" y="74" width="18" height="18" :fill="PAPER" />
      <!-- 高额头后梳发 -->
      <path d="M31 48 Q31 19 60 19 Q89 19 89 48 Q84 30 60 30 Q36 30 31 48 Z" :fill="INK" />
      <path d="M31 48 Q27 58 31 66 Q29 54 35 48 Z" :fill="INK" />
      <path d="M89 48 Q93 58 89 66 Q91 54 85 48 Z" :fill="INK" />
      <path d="M38 46 Q38 30 60 30 Q82 30 82 46 L82 60 Q82 78 60 78 Q38 78 38 60 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <!-- 络腮胡（较马克思短） -->
      <path d="M36 54 Q36 82 60 84 Q84 82 84 54 Q82 68 74 71 Q67 63 60 63 Q53 63 46 71 Q38 68 36 54 Z" :fill="INK" />
      <rect x="45" y="46" width="12" height="3.5" :fill="INK" />
      <rect x="63" y="46" width="12" height="3.5" :fill="INK" />
      <circle cx="51" cy="53" r="2.6" :fill="INK" />
      <circle cx="69" cy="53" r="2.6" :fill="INK" />
      <path d="M60 54 L60 63 L55 63" fill="none" :stroke="INK" stroke-width="2" />
    </g>

    <!-- ============ 列宁：光头 + 山羊胡 + 西装 ============ -->
    <g v-else-if="variant === 'lenin'">
      <path d="M12 120 Q14 92 40 86 L60 96 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <path d="M52 92 L60 104 L68 92 L60 96 Z" :fill="PAPER" />
      <rect x="51" y="74" width="18" height="18" :fill="PAPER" />
      <!-- 秃顶：仅在两侧与后方留少量头发 -->
      <path d="M31 56 Q30 36 44 28 Q36 40 37 56 Z" :fill="INK" opacity="0.85" />
      <path d="M89 56 Q90 36 76 28 Q84 40 83 56 Z" :fill="INK" opacity="0.85" />
      <path d="M38 46 Q38 30 60 30 Q82 30 82 46 L82 60 Q82 78 60 78 Q38 78 38 60 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <!-- 高额头（发际线靠后，以浅弧示意） -->
      <path d="M42 38 Q60 32 78 38" fill="none" :stroke="INK" stroke-width="1.6" opacity="0.5" />
      <!-- 标志性山羊胡 -->
      <path d="M52 68 Q60 88 68 68 Q64 74 60 74 Q56 74 52 68 Z" :fill="INK" />
      <path d="M56 62 Q60 66 64 62 Q60 70 56 62 Z" :fill="INK" />
      <!-- 眯眼 -->
      <path d="M45 52 L57 51" :stroke="INK" stroke-width="2.4" stroke-linecap="round" fill="none" />
      <path d="M63 51 L75 52" :stroke="INK" stroke-width="2.4" stroke-linecap="round" fill="none" />
      <rect x="45" y="44" width="12" height="3" :fill="INK" />
      <rect x="63" y="44" width="12" height="3" :fill="INK" />
      <path d="M60 52 L60 62 L55 62" fill="none" :stroke="INK" stroke-width="2" />
    </g>

    <!-- ============ 毛泽东：后梳背头 + 中山装 ============ -->
    <g v-else-if="variant === 'mao'">
      <!-- 中山装立领 -->
      <path d="M12 120 Q14 92 40 86 L60 98 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <path d="M50 88 L60 100 L70 88 L60 94 Z" :fill="PAPER" />
      <!-- 中山装四袋（几何示意） -->
      <rect x="26" y="102" width="16" height="12" fill="none" :stroke="PAPER" stroke-width="1.4" opacity="0.55" />
      <rect x="78" y="102" width="16" height="12" fill="none" :stroke="PAPER" stroke-width="1.4" opacity="0.55" />
      <!-- 立领 -->
      <path d="M46 82 L60 90 L74 82 L70 78 L60 84 L50 78 Z" :fill="PAPER" :stroke="INK" stroke-width="1.5" />
      <rect x="51" y="72" width="18" height="16" :fill="PAPER" />
      <!-- 饱满的后梳背头：高发际线 + 宽额 -->
      <path d="M30 50 Q28 17 60 16 Q92 17 90 50 Q88 32 74 27 Q60 23 46 27 Q32 32 30 50 Z" :fill="INK" />
      <!-- 两侧鬓角 -->
      <path d="M30 50 Q28 60 32 68 Q31 56 36 50 Z" :fill="INK" />
      <path d="M90 50 Q92 60 88 68 Q89 56 84 50 Z" :fill="INK" />
      <path d="M36 46 Q36 28 60 28 Q84 28 84 46 L84 60 Q84 79 60 79 Q36 79 36 60 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <!-- 眉眼：平直浓眉 -->
      <rect x="43" y="45" width="14" height="3.6" :fill="INK" />
      <rect x="63" y="45" width="14" height="3.6" :fill="INK" />
      <circle cx="50" cy="53" r="2.6" :fill="INK" />
      <circle cx="70" cy="53" r="2.6" :fill="INK" />
      <path d="M60 54 L60 64 L55 64" fill="none" :stroke="INK" stroke-width="2" />
    </g>

    <!-- ============ 罗莎·卢森堡：蓬松卷发 + 开领衬衫（无胡须） ============ -->
    <g v-else-if="variant === 'luxemburg'">
      <path d="M12 120 Q14 92 40 86 L60 96 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <!-- 开领衬衫 -->
      <path d="M52 90 L60 106 L68 90 L60 96 Z" :fill="PAPER" />
      <rect x="51" y="74" width="18" height="18" :fill="PAPER" />
      <!-- 浓密外扩卷发 -->
      <path d="M26 60 Q20 16 60 14 Q100 16 94 60 Q92 36 78 28 Q60 22 42 28 Q28 36 26 60 Z" :fill="INK" />
      <path d="M26 60 Q22 44 28 34 Q26 52 34 62 Z" :fill="INK" />
      <path d="M94 60 Q98 44 92 34 Q94 52 86 62 Z" :fill="INK" />
      <path d="M30 56 Q28 30 44 22 Q34 34 34 58 Z" :fill="INK" />
      <path d="M90 56 Q92 30 76 22 Q86 34 86 58 Z" :fill="INK" />
      <path d="M38 46 Q38 30 60 30 Q82 30 82 46 L82 60 Q82 78 60 78 Q38 78 38 60 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <!-- 眉与眼（无胡须，靠发型与开领辨认） -->
      <rect x="45" y="46" width="12" height="3.2" :fill="INK" />
      <rect x="63" y="46" width="12" height="3.2" :fill="INK" />
      <circle cx="51" cy="53" r="2.6" :fill="INK" />
      <circle cx="69" cy="53" r="2.6" :fill="INK" />
      <path d="M60 54 L60 63 L55 63" fill="none" :stroke="INK" stroke-width="2" />
      <!-- 嘴 -->
      <path d="M54 69 L66 69" :stroke="INK" stroke-width="2" stroke-linecap="round" fill="none" />
    </g>

    <!-- ============ 克拉拉·蔡特金：后梳发髻 + 高领（无胡须） ============ -->
    <g v-else-if="variant === 'zetkin'">
      <path d="M12 120 Q14 92 40 86 L60 96 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <!-- 高领 -->
      <path d="M48 86 L60 96 L72 86 L68 82 L60 88 L52 82 Z" :fill="PAPER" :stroke="INK" stroke-width="1.5" />
      <rect x="51" y="72" width="18" height="16" :fill="PAPER" />
      <!-- 脑后的发髻 -->
      <circle cx="88" cy="46" r="11" :fill="INK" />
      <!-- 中分后梳的头发 -->
      <path d="M30 48 Q30 18 60 18 Q90 18 90 48 Q84 28 60 28 Q36 28 30 48 Z" :fill="INK" />
      <path d="M30 48 Q27 58 31 66 Q29 54 35 48 Z" :fill="INK" />
      <path d="M90 48 Q93 58 89 66 Q91 54 85 48 Z" :fill="INK" />
      <path d="M38 46 Q38 30 60 30 Q82 30 82 46 L82 60 Q82 78 60 78 Q38 78 38 60 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <rect x="45" y="46" width="12" height="3.2" :fill="INK" />
      <rect x="63" y="46" width="12" height="3.2" :fill="INK" />
      <circle cx="51" cy="53" r="2.6" :fill="INK" />
      <circle cx="69" cy="53" r="2.6" :fill="INK" />
      <path d="M60 54 L60 63 L55 63" fill="none" :stroke="INK" stroke-width="2" />
      <path d="M54 69 L66 69" :stroke="INK" stroke-width="2" stroke-linecap="round" fill="none" />
    </g>

    <!-- ============ 胡志明：稀疏后梳发 + 长山羊胡 + 立领 ============ -->
    <g v-else-if="variant === 'ho-chi-minh'">
      <path d="M12 120 Q14 92 40 86 L60 98 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <!-- 立领 -->
      <path d="M46 82 L60 90 L74 82 L70 78 L60 84 L50 78 Z" :fill="PAPER" :stroke="INK" stroke-width="1.5" />
      <rect x="51" y="72" width="18" height="16" :fill="PAPER" />
      <!-- 稀疏的头发（高发际线） -->
      <path d="M33 48 Q33 22 60 22 Q87 22 87 48 Q82 32 60 32 Q38 32 33 48 Z" :fill="INK" opacity="0.75" />
      <path d="M33 48 Q30 56 34 62 Q32 52 38 46 Z" :fill="INK" opacity="0.75" />
      <path d="M87 48 Q90 56 86 62 Q88 52 82 46 Z" :fill="INK" opacity="0.75" />
      <path d="M38 46 Q38 30 60 30 Q82 30 82 46 L82 60 Q82 78 60 78 Q38 78 38 60 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <path d="M42 40 Q60 34 78 40" fill="none" :stroke="INK" stroke-width="1.6" opacity="0.45" />
      <!-- 长而稀疏的山羊胡 -->
      <path d="M54 66 Q60 92 66 66 Q63 74 60 74 Q57 74 54 66 Z" :fill="INK" />
      <!-- 眉眼 -->
      <rect x="44" y="46" width="13" height="3.2" :fill="INK" />
      <rect x="63" y="46" width="13" height="3.2" :fill="INK" />
      <circle cx="51" cy="54" r="2.6" :fill="INK" />
      <circle cx="69" cy="54" r="2.6" :fill="INK" />
      <path d="M60 54 L60 63 L55 63" fill="none" :stroke="INK" stroke-width="2" />
    </g>

    <!-- ============ 普列汉诺夫：高额头后梳发 + 短络腮胡 + 领结 ============ -->
    <g v-else-if="variant === 'plekhanov'">
      <path d="M12 120 Q14 92 40 86 L60 96 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <path d="M52 92 L60 104 L68 92 L60 96 Z" :fill="PAPER" />
      <path d="M56 100 L52 108 L60 106 Z" :fill="spec.accent" />
      <path d="M64 100 L68 108 L60 106 Z" :fill="spec.accent" />
      <rect x="51" y="74" width="18" height="18" :fill="PAPER" />
      <path d="M31 48 Q31 19 60 19 Q89 19 89 48 Q84 30 60 30 Q36 30 31 48 Z" :fill="INK" />
      <path d="M31 48 Q27 58 31 66 Q29 54 35 48 Z" :fill="INK" />
      <path d="M89 48 Q93 58 89 66 Q91 54 85 48 Z" :fill="INK" />
      <path d="M38 46 Q38 30 60 30 Q82 30 82 46 L82 60 Q82 78 60 78 Q38 78 38 60 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <!-- 高额头（发际线靠后） -->
      <path d="M42 38 Q60 32 78 38" fill="none" :stroke="INK" stroke-width="1.6" opacity="0.5" />
      <!-- 短络腮胡 -->
      <path d="M36 56 Q36 84 60 86 Q84 84 84 56 Q82 70 74 72 Q67 64 60 64 Q53 64 46 72 Q38 70 36 56 Z" :fill="INK" />
      <rect x="45" y="46" width="12" height="3.4" :fill="INK" />
      <rect x="63" y="46" width="12" height="3.4" :fill="INK" />
      <circle cx="51" cy="53" r="2.6" :fill="INK" />
      <circle cx="69" cy="53" r="2.6" :fill="INK" />
      <path d="M60 54 L60 63 L55 63" fill="none" :stroke="INK" stroke-width="2" />
    </g>

    <!-- ============ 片山潜：后梳短发 + 圆框眼镜 + 八字胡 ============ -->
    <g v-else-if="variant === 'katayama'">
      <path d="M12 120 Q14 92 40 86 L60 96 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <path d="M52 92 L60 104 L68 92 L60 96 Z" :fill="PAPER" />
      <rect x="51" y="74" width="18" height="18" :fill="PAPER" />
      <path d="M31 46 Q31 18 60 18 Q89 18 89 46 Q84 28 60 28 Q36 28 31 46 Z" :fill="INK" />
      <path d="M31 46 Q28 56 32 64 Q30 52 36 46 Z" :fill="INK" />
      <path d="M89 46 Q92 56 88 64 Q90 52 84 46 Z" :fill="INK" />
      <path d="M38 46 Q38 30 60 30 Q82 30 82 46 L82 60 Q82 78 60 78 Q38 78 38 60 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <!-- 圆框眼镜 -->
      <circle cx="50" cy="54" r="8" fill="none" :stroke="INK" stroke-width="2.2" />
      <circle cx="70" cy="54" r="8" fill="none" :stroke="INK" stroke-width="2.2" />
      <path d="M58 54 L62 54" :stroke="INK" stroke-width="2.2" />
      <rect x="43" y="44" width="14" height="3.2" :fill="INK" />
      <rect x="63" y="44" width="14" height="3.2" :fill="INK" />
      <circle cx="50" cy="54" r="2.2" :fill="INK" />
      <circle cx="70" cy="54" r="2.2" :fill="INK" />
      <path d="M60 56 L60 64 L55 64" fill="none" :stroke="INK" stroke-width="1.8" />
      <!-- 八字胡 -->
      <path d="M48 68 Q60 63 72 68 Q66 73 60 71 Q54 73 48 68 Z" :fill="INK" />
    </g>

    <!-- ============ 李大钊：平头短发 + 浓密八字胡 + 中式立领 ============ -->
    <g v-else-if="variant === 'li-dazhao'">
      <path d="M12 120 Q14 92 40 86 L60 98 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <path d="M46 82 L60 92 L74 82 L70 78 L60 86 L50 78 Z" :fill="PAPER" :stroke="INK" stroke-width="1.5" />
      <rect x="51" y="72" width="18" height="16" :fill="PAPER" />
      <!-- 平头（顶部平直） -->
      <path d="M30 44 Q30 21 60 21 Q90 21 90 44 Q84 30 60 30 Q36 30 30 44 Z" :fill="INK" />
      <path d="M30 44 Q27 54 31 62 Q29 50 35 44 Z" :fill="INK" />
      <path d="M90 44 Q93 54 89 62 Q91 50 85 44 Z" :fill="INK" />
      <path d="M38 46 Q38 30 60 30 Q82 30 82 46 L82 60 Q82 78 60 78 Q38 78 38 60 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <rect x="43" y="45" width="14" height="3.8" :fill="INK" />
      <rect x="63" y="45" width="14" height="3.8" :fill="INK" />
      <circle cx="51" cy="53" r="2.6" :fill="INK" />
      <circle cx="69" cy="53" r="2.6" :fill="INK" />
      <path d="M60 54 L60 63 L55 63" fill="none" :stroke="INK" stroke-width="2" />
      <!-- 浓密八字胡 -->
      <path d="M46 66 Q60 59 74 66 Q67 73 60 70 Q53 73 46 66 Z" :fill="INK" />
      <path d="M46 66 Q41 68 39 73 Q46 70 48 68 Z" :fill="INK" />
      <path d="M74 66 Q79 68 81 73 Q74 70 72 68 Z" :fill="INK" />
    </g>

    <!-- ============ 葛兰西：圆顶礼帽 + 圆框眼镜 + 高领 ============ -->
    <g v-else-if="variant === 'gramsci'">
      <path d="M12 120 Q14 92 40 86 L60 96 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <path d="M48 86 L60 96 L72 86 L68 82 L60 88 L52 82 Z" :fill="PAPER" :stroke="INK" stroke-width="1.5" />
      <rect x="51" y="72" width="18" height="16" :fill="PAPER" />
      <!-- 圆顶礼帽与帽檐 -->
      <path d="M34 40 Q34 14 60 13 Q86 14 86 40 Z" :fill="INK" />
      <path d="M24 40 Q60 33 96 40 Q60 47 24 40 Z" :fill="INK" />
      <path d="M34 40 Q60 35 86 40 Q86 43 84 45 L36 45 Q34 43 34 40 Z" :fill="spec.accent" opacity="0.6" />
      <path d="M38 48 Q38 32 60 32 Q82 32 82 48 L82 62 Q82 80 60 80 Q38 80 38 62 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <!-- 圆框眼镜 -->
      <circle cx="50" cy="56" r="7.5" fill="none" :stroke="INK" stroke-width="2.2" />
      <circle cx="70" cy="56" r="7.5" fill="none" :stroke="INK" stroke-width="2.2" />
      <path d="M57.5 56 L62.5 56" :stroke="INK" stroke-width="2.2" />
      <rect x="43" y="46" width="14" height="3.2" :fill="INK" />
      <rect x="63" y="46" width="14" height="3.2" :fill="INK" />
      <circle cx="50" cy="56" r="2" :fill="INK" />
      <circle cx="70" cy="56" r="2" :fill="INK" />
      <path d="M60 58 L60 66 L55 66" fill="none" :stroke="INK" stroke-width="1.8" />
      <path d="M54 71 L66 71" :stroke="INK" stroke-width="2" stroke-linecap="round" fill="none" />
    </g>

    <!-- ============ 柯伦泰：齐耳短发 + 高领（无胡须） ============ -->
    <g v-else-if="variant === 'kollontai'">
      <path d="M12 120 Q14 92 40 86 L60 96 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <path d="M48 86 L60 96 L72 86 L68 82 L60 88 L52 82 Z" :fill="PAPER" :stroke="INK" stroke-width="1.5" />
      <rect x="51" y="72" width="18" height="16" :fill="PAPER" />
      <!-- 齐耳短发（波波头） -->
      <path d="M27 50 Q25 16 60 15 Q95 16 93 50 L93 74 Q88 62 88 46 Q83 28 60 26 Q37 28 32 46 Q32 62 27 74 Z" :fill="INK" />
      <path d="M38 46 Q38 30 60 30 Q82 30 82 46 L82 60 Q82 78 60 78 Q38 78 38 60 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <rect x="45" y="46" width="12" height="3.2" :fill="INK" />
      <rect x="63" y="46" width="12" height="3.2" :fill="INK" />
      <circle cx="51" cy="53" r="2.6" :fill="INK" />
      <circle cx="69" cy="53" r="2.6" :fill="INK" />
      <path d="M60 54 L60 63 L55 63" fill="none" :stroke="INK" stroke-width="2" />
      <path d="M54 69 L66 69" :stroke="INK" stroke-width="2" stroke-linecap="round" fill="none" />
    </g>

    <!-- ============ 卡斯特罗：军便帽 + 大胡子 + 军装 ============ -->
    <g v-else-if="variant === 'castro'">
      <path d="M12 120 Q14 92 40 86 L60 96 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <path d="M50 88 L60 100 L70 88 L66 84 L60 92 L54 84 Z" :fill="spec.accent" opacity="0.85" />
      <rect x="51" y="74" width="18" height="18" :fill="PAPER" />
      <!-- 军便帽 -->
      <path d="M32 40 Q36 16 60 16 Q84 16 88 40 Z" :fill="INK" />
      <path d="M30 40 Q60 32 90 40 Q60 47 30 40 Z" :fill="INK" />
      <path d="M32 40 Q60 34 88 40 L86 44 Q60 38 34 44 Z" :fill="spec.accent" opacity="0.6" />
      <path d="M38 48 Q38 32 60 32 Q82 32 82 48 L82 62 Q82 80 60 80 Q38 80 38 62 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <!-- 浓密大胡子 -->
      <path d="M38 58 Q38 88 60 90 Q82 88 82 58 Q78 72 68 74 Q62 66 60 66 Q58 66 52 74 Q42 72 38 58 Z" :fill="INK" />
      <rect x="44" y="47" width="13" height="3.4" :fill="INK" />
      <rect x="63" y="47" width="13" height="3.4" :fill="INK" />
      <circle cx="51" cy="55" r="2.8" :fill="INK" />
      <circle cx="69" cy="55" r="2.8" :fill="INK" />
      <path d="M60 56 L60 64 L55 64" fill="none" :stroke="INK" stroke-width="2" />
    </g>

    <!-- ============ 桑卡拉：贝雷帽 + 八字胡 + 军装 ============ -->
    <g v-else-if="variant === 'sankara'">
      <path d="M12 120 Q14 92 40 86 L60 96 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <path d="M52 92 L60 104 L68 92 L60 96 Z" :fill="PAPER" />
      <path d="M58 96 L62 96 L64 116 L60 120 L56 116 Z" :fill="spec.accent" />
      <rect x="51" y="74" width="18" height="18" :fill="PAPER" />
      <!-- 贝雷帽（右倾，带金属帽徽） -->
      <path d="M28 40 Q32 12 64 13 Q96 15 96 40 Z" :fill="INK" />
      <path d="M28 40 Q32 36 64 35 Q96 34 96 40 Q64 44 28 40 Z" :fill="INK" />
      <path d="M86 16 q11 -4 10 7 q-6 -5 -11 -3 Z" :fill="spec.accent" />
      <path d="M31 44 Q28 62 33 70 Q31 54 36 46 Z" :fill="INK" />
      <path d="M89 44 Q92 62 87 70 Q89 54 84 46 Z" :fill="INK" />
      <path d="M38 48 Q38 32 60 32 Q82 32 82 48 L82 62 Q82 80 60 80 Q38 80 38 62 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <rect x="45" y="47" width="12" height="3.4" :fill="INK" />
      <rect x="63" y="47" width="12" height="3.4" :fill="INK" />
      <circle cx="51" cy="55" r="2.6" :fill="INK" />
      <circle cx="69" cy="55" r="2.6" :fill="INK" />
      <path d="M60 56 L60 64 L55 64" fill="none" :stroke="INK" stroke-width="2" />
      <!-- 八字胡 -->
      <path d="M46 68 Q60 62 74 68 Q67 74 60 71 Q53 74 46 68 Z" :fill="INK" />
    </g>

    <!-- ============ 切·格瓦拉：贝雷帽 + 长发 + 胡须 ============ -->
    <g v-else>
      <path d="M12 120 Q14 92 40 86 L60 96 L80 86 Q106 92 108 120 Z" :fill="spec.coat" />
      <path d="M52 92 L60 104 L68 92 L60 96 Z" :fill="PAPER" />
      <rect x="51" y="74" width="18" height="18" :fill="PAPER" />
      <!-- 标志性贝雷帽（略歪，带红星） -->
      <path d="M26 40 Q30 12 62 13 Q95 14 96 40 Z" :fill="INK" />
      <path d="M26 40 Q30 36 62 35 Q95 34 96 40 Q62 44 26 40 Z" :fill="INK" />
      <path d="M92 16 l12 -6 l-4 11 Z" :fill="spec.accent" />
      <!-- 帽下头发 -->
      <path d="M31 44 Q28 62 33 70 Q31 54 36 46 Z" :fill="INK" />
      <path d="M89 44 Q92 62 87 70 Q89 54 84 46 Z" :fill="INK" />
      <path d="M38 48 Q38 32 60 32 Q82 32 82 48 L82 62 Q82 80 60 80 Q38 80 38 62 Z" :fill="PAPER" :stroke="INK" stroke-width="2.5" />
      <!-- 长发（垂至肩） -->
      <path d="M33 50 Q30 74 38 84 Q34 66 38 54 Z" :fill="INK" />
      <path d="M87 50 Q90 74 82 84 Q86 66 82 54 Z" :fill="INK" />
      <!-- 胡须 -->
      <path d="M42 60 Q44 80 60 82 Q76 80 78 60 Q74 72 60 74 Q46 72 42 60 Z" :fill="INK" />
      <!-- 眉眼 -->
      <rect x="44" y="46" width="13" height="3.4" :fill="INK" />
      <rect x="63" y="46" width="13" height="3.4" :fill="INK" />
      <circle cx="51" cy="54" r="2.8" :fill="INK" />
      <circle cx="69" cy="54" r="2.8" :fill="INK" />
      <path d="M60 54 L60 64 L55 64" fill="none" :stroke="INK" stroke-width="2" />
    </g>
  </svg>
</template>
