/**
 * 站点结构性常量。
 *
 * 文案（站点名、标语、简介、版权声明等）以 `content/site.yaml` 为**单一来源**，
 * 由后端经 `/api/meta` 提供，前端通过 `useSiteCopy()` 读取；
 * 此处保留同名同值的**回退常量**，仅在 meta 尚未加载时使用，避免首屏闪烁。
 * 导航与分区结构属于前端结构，仍在此处定义。
 */
export const SITE = {
  name: '星火档案馆',
  nameEn: 'SPARK ARCHIVE',
  tagline: '全世界无产者，联合起来！',
  subtitle: '无产阶级革命理论与实践的文献档案',
  description:
    '星火档案馆收录马克思、恩格斯、列宁、毛泽东、切·格瓦拉五个人的生平、著作与思想，记录国际共产主义运动的理论源流与实践历程。',
  footerNote: '本站内容用于学习与研究，史料均标注出处。争议性评价如实标注。',
  licenseNote: '本站原创文字采用 CC BY-NC-SA 4.0 许可；所引原著摘录版权归各自权利人所有。',
  est: 'EST. 2026',
  archiveNo: 'NO. 12026',
} as const

export const NAV = [
  { to: '/people', label: '人物' },
  { to: '/works', label: '著作' },
  { to: '/timeline', label: '年表' },
  { to: '/glossary', label: '术语' },
  { to: '/graph', label: '关系' },
  { to: '/search', label: '检索' },
  { to: '/about', label: '凡例' },
] as const

/** 贡献领域：政治 / 经济 / 文化 / 思想 */
export const SECTION_KINDS = [
  { key: 'politics', label: '政治贡献', no: '01' },
  { key: 'economy', label: '经济贡献', no: '02' },
  { key: 'culture', label: '文化贡献', no: '03' },
  { key: 'thought', label: '思想体系', no: '04' },
] as const

export type SectionKind = (typeof SECTION_KINDS)[number]['key']

/** 人物筛选下拉/按钮的统一数据源（「全部」+ 5 位收录人物）。 */
export const PEOPLE_FILTERS = [
  { slug: '', name: '全部' },
  { slug: 'marx', name: '马克思' },
  { slug: 'engels', name: '恩格斯' },
  { slug: 'lenin', name: '列宁' },
  { slug: 'mao', name: '毛泽东' },
  { slug: 'guevara', name: '切·格瓦拉' },
] as const
