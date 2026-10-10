/**
 * 人物肖像映射表 —— 内容 `portrait_key`（= 人物 slug）到静态照片资产的连接键。
 *
 * 照片取自公开的历史照片与史料图片，**逐张标注来源与许可状态**
 * （见 `public/assets/portraits/CREDITS.md`），亦在《凡例》页「图片来源」区块展示。
 * 图片随前端构建打包于 `/assets/portraits/`。
 * 映射缺失或图片加载失败时，`Portrait.vue` 自动回退为本站绘制的构成主义矢量插画。
 *
 * 注意：**并非全部照片都处于公有领域**。维基共享资源来源的 8 张为公有领域史料；
 * 百度百科来源的 8 张原图未标注许可，依站点运营方的声明使用。
 * 新增照片时请如实填写 license，不要一律写「公有领域」。
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
  zetkin: {
    name: '克拉拉·蔡特金',
    webp: `${dir}/zetkin.webp`,
    jpg: `${dir}/zetkin.jpg`,
    credit: '摄影者未载明（百度百科词条图）',
    source: 'https://baike.baidu.com/item/%E5%85%8B%E6%8B%89%E6%8B%89%C2%B7%E8%94%A1%E7%89%B9%E9%87%91',
    license: '未标注许可（来源：百度百科）',
    note: '原图 350×513，裁切放大至 480×480 输出；来源为词条摘要图，非维基共享资源的公有领域档案照。',
  },
  plekhanov: {
    name: '格奥尔基·普列汉诺夫',
    webp: `${dir}/plekhanov.webp`,
    jpg: `${dir}/plekhanov.jpg`,
    credit: '摄影者未载明（百度百科词条图）',
    source: 'https://baike.baidu.com/item/%E6%99%AE%E5%88%97%E6%B1%89%E8%AF%BA%E5%A4%AB',
    license: '未标注许可（来源：百度百科）',
  },
  katayama: {
    name: '片山潜',
    webp: `${dir}/katayama.webp`,
    jpg: `${dir}/katayama.jpg`,
    credit: '摄影者未载明（百度百科词条图）',
    source: 'https://baike.baidu.com/item/%E7%89%87%E5%B1%B1%E6%BD%9C',
    license: '未标注许可（来源：百度百科）',
    note: '原图为现代数字修复／上色版本，非原始档案照；已裁切至 480×480。',
  },
  kollontai: {
    name: '亚历山德拉·柯伦泰',
    webp: `${dir}/kollontai.webp`,
    jpg: `${dir}/kollontai.jpg`,
    credit: '摄影者未载明（百度百科词条图）',
    source: 'https://baike.baidu.com/item/%E6%9F%AF%E4%BC%A6%E6%B3%B0',
    license: '未标注许可（来源：百度百科）',
  },
  'li-dazhao': {
    name: '李大钊',
    webp: `${dir}/li-dazhao.webp`,
    jpg: `${dir}/li-dazhao.jpg`,
    credit: '摄影者未载明（百度百科词条图）',
    source: 'https://baike.baidu.com/item/%E6%9D%8E%E5%A4%A7%E9%92%8A',
    license: '未标注许可（来源：百度百科）',
  },
  gramsci: {
    name: '安东尼奥·葛兰西',
    webp: `${dir}/gramsci.webp`,
    jpg: `${dir}/gramsci.jpg`,
    credit: '摄影者未载明（百度百科词条图）',
    source: 'https://baike.baidu.com/item/%E5%AE%89%E4%B8%9C%E5%B0%BC%E5%A5%A5%C2%B7%E8%91%9B%E5%85%B0%E8%A5%BF',
    license: '未标注许可（来源：百度百科）',
    note: '原图为青年时期的现代数字修复／上色版本，非原始档案照；已裁切至 480×480。',
  },
  castro: {
    name: '菲德尔·卡斯特罗',
    webp: `${dir}/castro.webp`,
    jpg: `${dir}/castro.jpg`,
    credit: '摄影者未载明（百度百科词条图）',
    source: 'https://baike.baidu.com/item/%E8%8F%B2%E5%BE%B7%E5%B0%94%C2%B7%E5%8D%A1%E6%96%AF%E7%89%B9%E7%BD%97',
    license: '未标注许可（来源：百度百科）',
  },
  sankara: {
    name: '托马斯·桑卡拉',
    webp: `${dir}/sankara.webp`,
    jpg: `${dir}/sankara.jpg`,
    credit: '摄影者未载明（百度百科词条图）',
    source: 'https://baike.baidu.com/item/%E6%89%98%E9%A9%AC%E6%96%AF%C2%B7%E6%A1%91%E5%8D%A1%E6%8B%89',
    license: '未标注许可（来源：百度百科）',
    note: '原图 1024×719，按人像精确取景裁切至 480×480。',
  },
}

export function getPortrait(key: string): PortraitMeta | undefined {
  return PORTRAITS[key]
}
