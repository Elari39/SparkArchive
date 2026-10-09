import { computed } from 'vue'
import { useSiteStore } from '@/stores/site'
import { SITE } from '@/config/site'

/**
 * 站点文案的统一读取入口。
 *
 * 文案以 `content/site.yaml` 为**单一来源**，由后端经 `/api/meta` 提供；
 * 在 meta 尚未加载（或加载失败）时回退到 `config/site.ts` 中的同值常量，避免首屏闪烁。
 */
export function useSiteCopy() {
  const site = useSiteStore()
  const s = computed(() => site.meta?.site)

  return {
    name: computed(() => s.value?.name ?? SITE.name),
    nameEn: computed(() => s.value?.name_en ?? SITE.nameEn),
    tagline: computed(() => s.value?.tagline ?? SITE.tagline),
    subtitle: computed(() => s.value?.subtitle ?? SITE.subtitle),
    description: computed(() => s.value?.description ?? SITE.description),
    footerNote: computed(() => s.value?.footer_note ?? SITE.footerNote),
    licenseNote: computed(() => s.value?.license_note ?? SITE.licenseNote),
  }
}
