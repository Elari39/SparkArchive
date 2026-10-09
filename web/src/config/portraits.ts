/**
 * 人物肖像映射表 —— 内容 `portrait_key`（= 人物 slug）到静态照片资产的连接键。
 *
 * 照片为**公有领域**历史照片，逐张署名见 `public/assets/portraits/CREDITS.md`，
 * 亦在《凡例》页「图片来源」区块展示。图片随前端构建打包于 `/assets/portraits/`。
 * 映射缺失或图片加载失败时，`Portrait.vue` 自动回退为本站绘制的构成主义矢量插画。
 */

export interface PortraitMeta {
  /** 人物中文名（用于 img 的 alt） */
  name: string
  /** WebP 主格式路径 */
  webp: string
  /** JPEG 回退路径 */
  jpg: string
  /** 摄影者与年代 */
  credit: string
  /** 原始文件来源页 */
  source: string
  /** 许可 */
  license: string
  /** 可选备注（如格瓦拉的道德权说明） */
  note?: string
}

const dir = '/assets/portraits'

export const PORTRAITS: Record<string, PortraitMeta> = {
  marx: {
    name: '卡尔·马克思',
    webp: `${dir}/marx.webp`,
    jpg: `${dir}/marx.jpg`,
    credit: 'John Jabez Edwin Mayall，1875 年',
    source: 'https://commons.wikimedia.org/wiki/File:Karl_Marx_001.jpg',
    license: '公有领域',
  },
  engels: {
    name: '弗里德里希·恩格斯',
    webp: `${dir}/engels.webp`,
    jpg: `${dir}/engels.jpg`,
    credit: '作者未知，1879 年',
    source: 'https://commons.wikimedia.org/wiki/File:Friedrich_Engels_portrait.jpg',
    license: '公有领域',
  },
  lenin: {
    name: '列宁',
    webp: `${dir}/lenin.webp`,
    jpg: `${dir}/lenin.jpg`,
    credit: 'Pavel S. Zhukov，1920 年',
    source: 'https://commons.wikimedia.org/wiki/File:Vladimir_Lenin.jpg',
    license: '公有领域',
  },
  luxemburg: {
    name: '罗莎·卢森堡',
    webp: `${dir}/luxemburg.webp`,
    jpg: `${dir}/luxemburg.jpg`,
    credit: '摄影者与年代未载明',
    source: 'https://commons.wikimedia.org/wiki/File:Rosa_Luxemburg.jpg',
    license: '公有领域（PD-old）',
  },
  dimitrov: {
    name: '格奥尔基·季米特洛夫',
    webp: `${dir}/dimitrov.webp`,
    jpg: `${dir}/dimitrov.jpg`,
    credit: '摄影者与年代未载明',
    source: 'https://commons.wikimedia.org/wiki/File:Georgi_Dimitrov.jpg',
    license: '公有领域',
    note: '原始照片分辨率较低（232×299），放大至 480×480 输出，站点以单色滤镜呈现。',
  },
  'ho-chi-minh': {
    name: '胡志明',
    webp: `${dir}/ho-chi-minh.webp`,
    jpg: `${dir}/ho-chi-minh.jpg`,
    credit: '摄影者未载明，1946 年',
    source: 'https://commons.wikimedia.org/wiki/File:Ho_Chi_Minh_1946.jpg',
    license: '公有领域（PD-Vietnam：作者去世逾 50 年）',
    note: '原始照片分辨率较低（282×383），放大至 480×480 输出，站点以单色滤镜呈现。',
  },
  mao: {
    name: '毛泽东',
    webp: `${dir}/mao.webp`,
    jpg: `${dir}/mao.jpg`,
    credit: '孟庆彪、侯波（新华社），1959 年',
    source: 'https://commons.wikimedia.org/wiki/File:Mao_Zedong_1959.jpg',
    license: '公有领域',
  },
  guevara: {
    name: '切·格瓦拉',
    webp: `${dir}/guevara.webp`,
    jpg: `${dir}/guevara.jpg`,
    credit: 'Alberto Korda，1960 年',
    source: 'https://commons.wikimedia.org/wiki/File:CheHigh.jpg',
    license: '公有领域（道德权争议见凡例）',
    note: '《英勇的游击队员》：维基共享资源依古巴 1994 年第 156 号法令与美国未履行手续出版标注为公有领域。',
  },
}

export function getPortrait(key: string): PortraitMeta | undefined {
  return PORTRAITS[key]
}
