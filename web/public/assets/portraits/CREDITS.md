# 人物肖像图片来源与许可

本站人物肖像采用**公有领域**历史照片，均取自维基共享资源（Wikimedia Commons），
经求闻百科（qiuwenbaike.cn）镜像获取后本地打包，未作内容修改，仅作方形裁切与尺寸压缩。

照片仅用于学习与研究，本站不作商业使用。原始文件页链接见下表。

| 人物 | 本地文件 | 原始文件 | 摄影者 | 年代 | 许可 |
|---|---|---|---|---|---|
| 马克思 marx | `marx.webp` / `marx.jpg` | [File:Karl Marx 001.jpg](https://commons.wikimedia.org/wiki/File:Karl_Marx_001.jpg) | John Jabez Edwin Mayall（1813–1901） | 1875 年 8 月前 | 公有领域 |
| 恩格斯 engels | `engels.webp` / `engels.jpg` | [File:Friedrich Engels portrait.jpg](https://commons.wikimedia.org/wiki/File:Friedrich_Engels_portrait.jpg) | 作者未知 | 1879 年 | 公有领域 |
| 列宁 lenin | `lenin.webp` / `lenin.jpg` | [File:Vladimir Lenin.jpg](https://commons.wikimedia.org/wiki/File:Vladimir_Lenin.jpg) | Pavel Semyonovich Zhukov（1870–1942） | 1920 年 | 公有领域 |
| 罗莎·卢森堡 luxemburg | `luxemburg.webp` / `luxemburg.jpg` | [File:Rosa Luxemburg.jpg](https://commons.wikimedia.org/wiki/File:Rosa_Luxemburg.jpg) | 未载明 | 未载明 | 公有领域（PD-old） |
| 格奥尔基·季米特洛夫 dimitrov | `dimitrov.webp` / `dimitrov.jpg` | [File:Georgi Dimitrov.jpg](https://commons.wikimedia.org/wiki/File:Georgi_Dimitrov.jpg) | 未载明 | 未载明 | 公有领域 |
| 胡志明 ho-chi-minh | `ho-chi-minh.webp` / `ho-chi-minh.jpg` | [File:Ho Chi Minh 1946.jpg](https://commons.wikimedia.org/wiki/File:Ho_Chi_Minh_1946.jpg) | 未载明 | 1946 年 | 公有领域（PD-Vietnam：作者去世逾 50 年） |
| 毛泽东 mao | `mao.webp` / `mao.jpg` | [File:Mao Zedong 1959.jpg](https://commons.wikimedia.org/wiki/File:Mao_Zedong_1959.jpg) | 孟庆彪（新华社摄影部中央新闻组）、侯波（新华社驻中南海记者） | 1959 年 | 公有领域 |
| 切·格瓦拉 guevara | `guevara.webp` / `guevara.jpg` | [File:CheHigh.jpg](https://commons.wikimedia.org/wiki/File:CheHigh.jpg) | Alberto Korda（1928–2001） | 1960 年 3 月 5 日 | 标为公有领域（见下） |

> **克拉拉·蔡特金（zetkin）暂无照片**：镜像站未收录其公有领域照片，故前端回退为本站绘制的
> 构成主义矢量插画（见 `PortraitArt.vue` 的 `zetkin` 分支）。补图后只需把 `zetkin.jpg` /
> `zetkin.webp` 放入本目录，并在 `web/src/config/portraits.ts` 中登记即可，无需改动其他代码。

## 关于许可的说明

- **马克思、列宁、毛泽东、恩格斯**：摄影作品已因年代久远进入公有领域（作者去世逾 70 年，或首次出版超过法定保护期）。毛泽东照片摄于 1959 年、由新华社记者拍摄，中国摄影作品著作权保护期为首次发表后 50 年，现已届满。
- **罗莎·卢森堡、格奥尔基·季米特洛夫**：两张照片在来源文件页均标注为**公有领域**（卢森堡为 PD-old，即作者去世逾 70 年；季米特洛夫照片作者未载明，按拍摄年代已届满保护期）。摄影者与确切年代在来源档案中未标注，本站如实记为「未载明」，未作推测。
- **胡志明**：所用照片来自来源文件页标注的 **PD-Vietnam** 模板（作者去世逾 50 年），文件名载明摄于 1946 年。
- **切·格瓦拉**：所用为 Alberto Korda《英勇的游击队员》（Guerrillero Heroico）。维基共享资源依古巴 1994 年第 156 号法令（照片首次使用后 25 年）与美国未履行著作权手续之出版而标注为**公有领域**。需要说明的是，该照片历史上曾出现围绕摄影者道德权利的争议；本站仅作史料展示，不作任何商业用途，并在《凡例》页如实记录。

## 处理方式

- 由原始 JPG 裁切为 480×480 方形，输出 WebP（质量 82）与 JPEG（质量 80）双格式。
- 站点通过 CSS 滤镜统一为**单色**，以与构成主义红／黑／米白色板协调。
- 若某张照片缺失或加载失败，前端自动回退为本站绘制的构成主义矢量插画。
- 新肖像的原始分辨率较低者，已在 `portraits.ts` 的 `note` 字段中注明（季米特洛夫 232×299、胡志明 282×383），以便日后有条件时替换为更高分辨率版本。
