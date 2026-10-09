import DOMPurify from 'dompurify'

/**
 * 集中式 HTML 净化。
 *
 * 站点会对两处内容使用 `v-html`：
 *  1. MarkdownBlock —— marked 渲染出的正文；
 *  2. SearchView —— 检索结果摘要的高亮。
 *
 * 二者最终都进入 `v-html`，因此统一在**渲染前**经白名单净化，
 * 杜绝上游数据一旦引入外部/用户内容时的存储型 XSS。
 */

/** 正文渲染：允许常规富文本标签与安全属性。 */
export function sanitizeHtml(html: string): string {
  return DOMPurify.sanitize(html, {
    ALLOWED_TAGS: [
      'a', 'p', 'br', 'hr', 'strong', 'em', 'del', 's', 'code', 'pre', 'blockquote',
      'ul', 'ol', 'li', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
      'table', 'thead', 'tbody', 'tr', 'th', 'td',
    ],
    ALLOWED_ATTR: ['href', 'title', 'class'],
    // 仅允许安全协议，阻断 javascript: / data: 等
    ALLOWED_URI_REGEXP: /^(?:https?|mailto):/i,
  })
}

/** 检索高亮：只保留 `<mark class="hit">`，其余一律剥离。 */
export function sanitizeHighlight(html: string): string {
  return DOMPurify.sanitize(html, {
    ALLOWED_TAGS: ['mark'],
    ALLOWED_ATTR: ['class'],
  })
}

/** 转义 HTML 特殊字符（含单引号），供拼接高亮片段前使用。 */
export function escapeHtml(text: string): string {
  return text.replace(/[&<>"']/g, (ch) => {
    switch (ch) {
      case '&':
        return '&amp;'
      case '<':
        return '&lt;'
      case '>':
        return '&gt;'
      case '"':
        return '&quot;'
      default:
        return '&#39;'
    }
  })
}
