package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Hit 一条检索结果。
type Hit struct {
	Kind    string `json:"kind"`
	Ref     string `json:"ref"`
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
	Person  string `json:"person"`
}

// minTrigramLen 是 FTS5 trigram 分词器能匹配的最短查询长度。
// 少于 3 个字符的查询（如"列宁"）无法用 trigram 命中，需降级为 LIKE。
const minTrigramLen = 3

// FTS snippet() 中摘要的最大 token 数（左起，超出以省略号收尾）。
const snippetTokens = 24

// LIKE 降级路径中摘要截取的字符数（CJK 场景下 120 字足够展示一行余）。
const snippetChars = 120

// fallbackLimit 是 limit 未由上游指定时的防御性默认页大小。
const fallbackLimit = 50

// Search 执行全文检索。
//
// 中文处理要点：FTS5 的 trigram 分词器可以匹配 CJK 子串，但要求查询至少 3 个字符；
// 因此短查询降级为 LIKE 子串匹配，保证"列宁""游击"这类词仍能命中。
//
// limit/offset 的取值边界由调用方（api 层）负责，此处只做防御性兜底。
func (s *Store) Search(ctx context.Context, q, kind string, limit, offset int) ([]Hit, int, error) {
	q = strings.TrimSpace(q)
	if q == "" {
		return []Hit{}, 0, nil
	}
	if limit <= 0 {
		limit = fallbackLimit
	}
	if offset < 0 {
		offset = 0
	}

	var (
		hits  []Hit
		total int
		err   error
	)
	if utf8.RuneCountInString(q) >= minTrigramLen {
		hits, total, err = s.searchFTS(ctx, q, kind, limit, offset)
	} else {
		hits, total, err = s.searchLike(ctx, q, kind, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	return hits, total, nil
}

// ftsQuery 把用户输入转成安全的 FTS5 短语查询。
// 用双引号包裹整个查询，内部的双引号翻倍转义，避免语法错误与注入。
func ftsQuery(q string) string {
	return `"` + strings.ReplaceAll(q, `"`, `""`) + `"`
}

// searchQuery 两条检索路径（FTS 与 LIKE 降级）的公共骨架：
// 先 COUNT 总数，再按 listSQL 取当页，统一走 scanHits。
//
// 注：FTS5 的 snippet() 无法与窗口函数同查询使用
// （"unable to use function snippet in the requested context"），
// 故 total 仍需单独 COUNT。数据只读且命中集小，代价可接受。
func (s *Store) searchQuery(ctx context.Context, where string, args []any, listSQL string, limit, offset int) ([]Hit, int, error) {
	var total int
	countSQL := "SELECT COUNT(*) FROM search_fts WHERE " + where
	if err := s.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("统计检索结果: %w", err)
	}

	listArgs := append(append([]any{}, args...), limit, offset)
	rows, err := s.db.QueryContext(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("执行检索: %w", err)
	}
	defer rows.Close()
	hits, err := scanHits(rows)
	if err != nil {
		return nil, 0, err
	}
	return hits, total, nil
}

// searchFTS 长查询走 FTS5 trigram MATCH，按 bm25 排序并用 snippet 生成摘要。
func (s *Store) searchFTS(ctx context.Context, q, kind string, limit, offset int) ([]Hit, int, error) {
	where := "search_fts MATCH ?"
	args := []any{ftsQuery(q)}
	if kind != "" {
		where += " AND kind = ?"
		args = append(args, kind)
	}
	listSQL := fmt.Sprintf(`SELECT kind, ref, title, snippet(search_fts, 1, '', '', '…', %d), person
	            FROM search_fts WHERE `, snippetTokens) + where + ` ORDER BY rank LIMIT ? OFFSET ?`
	return s.searchQuery(ctx, where, args, listSQL, limit, offset)
}

// searchLike 是短查询的降级路径：对标题与正文做子串匹配。
func (s *Store) searchLike(ctx context.Context, q, kind string, limit, offset int) ([]Hit, int, error) {
	like := "%" + escapeLike(q) + "%"
	where := "(title LIKE ? ESCAPE '\\' OR body LIKE ? ESCAPE '\\')"
	args := []any{like, like}
	if kind != "" {
		where += " AND kind = ?"
		args = append(args, kind)
	}
	listSQL := fmt.Sprintf(`SELECT kind, ref, title, substr(body, 1, %d), person
	            FROM search_fts WHERE `, snippetChars) + where + ` LIMIT ? OFFSET ?`
	return s.searchQuery(ctx, where, args, listSQL, limit, offset)
}

// escapeLike 转义 LIKE 的通配符，使用户输入按字面量匹配。
func escapeLike(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return r.Replace(s)
}

func scanHits(rows *sql.Rows) ([]Hit, error) {
	hits := []Hit{} // 非 nil，保证 JSON 序列化为 [] 而不是 null
	for rows.Next() {
		var h Hit
		if err := rows.Scan(&h.Kind, &h.Ref, &h.Title, &h.Snippet, &h.Person); err != nil {
			return nil, fmt.Errorf("扫描检索结果: %w", err)
		}
		h.Snippet = strings.TrimSpace(strings.ReplaceAll(h.Snippet, "\n", " "))
		hits = append(hits, h)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历检索结果: %w", err)
	}
	return hits, nil
}
