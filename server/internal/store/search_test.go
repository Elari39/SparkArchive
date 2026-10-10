package store

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/sparkarchive/server/internal/content"
)

// newTestStore 用仓库内真实内容构建一个内存态数据库，用于检索测试。
func newTestStore(t *testing.T) *Store {
	t.Helper()
	ctx := t.Context()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	st, err := Open(dbPath)
	if err != nil {
		t.Fatalf("打开测试数据库: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	if err := st.Migrate(); err != nil {
		t.Fatalf("建表: %v", err)
	}
	a, err := content.Load(filepath.Join("..", "..", "..", "content"))
	if err != nil {
		t.Fatalf("加载内容: %v", err)
	}
	if err := st.Ingest(ctx, a); err != nil {
		t.Fatalf("入库: %v", err)
	}
	return st
}

// realPeopleCount 返回 content/ 下实际收录的人物数。
// 用于替代写死的数字：本包关心的是「入库后数量与内容源一致」这一不变量，
// 而不是某个特定数字，故新增人物时这里无需改动。
func realPeopleCount(t *testing.T) int {
	t.Helper()
	a, err := content.Load(filepath.Join("..", "..", "..", "content"))
	if err != nil {
		t.Fatalf("加载内容: %v", err)
	}
	return len(a.People)
}

// TestSearchChinese 是本项目最关键的一组测试：
// SQLite FTS5 默认的 unicode61 分词器不切分中文，若配置不当中文检索会完全失效。
func TestSearchChinese(t *testing.T) {
	t.Parallel()
	st := newTestStore(t)
	ctx := t.Context()

	tests := []struct {
		name     string
		query    string
		wantHits bool // 是否应当有结果
		minHits  int
	}{
		{name: "长词-无产阶级", query: "无产阶级", wantHits: true, minHits: 3},
		{name: "长词-剩余价值", query: "剩余价值", wantHits: true, minHits: 1},
		{name: "长词-十月革命", query: "十月革命", wantHits: true, minHits: 1},
		{name: "长词-帝国主义", query: "帝国主义", wantHits: true, minHits: 1},
		// 2 字查询无法用 trigram 匹配，必须走 LIKE 降级路径
		{name: "短词-列宁", query: "列宁", wantHits: true, minHits: 1},
		{name: "短词-游击", query: "游击", wantHits: true, minHits: 1},
		{name: "短词-群众", query: "群众", wantHits: true, minHits: 1},
		{name: "短词-古巴", query: "古巴", wantHits: true, minHits: 1},
		{name: "不存在的词", query: "紫色独角兽不存在的词", wantHits: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hits, total, err := st.Search(ctx, tc.query, "", 50, 0)
			if err != nil {
				t.Fatalf("检索 %q 出错: %v", tc.query, err)
			}
			if tc.wantHits && total < tc.minHits {
				t.Errorf("检索 %q 命中 %d 条，期望至少 %d 条", tc.query, total, tc.minHits)
			}
			if !tc.wantHits && total != 0 {
				t.Errorf("检索 %q 期望 0 条，实际 %d 条", tc.query, total)
			}
			if len(hits) > 0 && hits[0].Title == "" {
				t.Errorf("检索 %q 的结果缺少标题", tc.query)
			}
		})
	}
}

// TestSearchCrossKind 确认同一关键词能跨类型命中。
func TestSearchCrossKind(t *testing.T) {
	t.Parallel()
	st := newTestStore(t)

	hits, _, err := st.Search(t.Context(), "无产阶级", "", 100, 0)
	if err != nil {
		t.Fatalf("检索出错: %v", err)
	}
	kinds := map[string]bool{}
	for _, h := range hits {
		kinds[h.Kind] = true
	}
	if len(kinds) < 2 {
		t.Errorf("期望跨多种类型命中，实际只有 %v", kinds)
	}
}

// TestSearchKindFilter 确认类型过滤生效。
func TestSearchKindFilter(t *testing.T) {
	t.Parallel()
	st := newTestStore(t)

	for _, kind := range []string{"term", "work", "person"} {
		hits, _, err := st.Search(t.Context(), "无产阶级", kind, 50, 0)
		if err != nil {
			t.Fatalf("检索 %s 出错: %v", kind, err)
		}
		for _, h := range hits {
			if h.Kind != kind {
				t.Errorf("类型过滤 %s 失效，出现 %s", kind, h.Kind)
			}
		}
	}
}

// TestSearchInjection 确认恶意输入不会破坏查询。
func TestSearchInjection(t *testing.T) {
	t.Parallel()
	st := newTestStore(t)
	ctx := t.Context()

	queries := []string{
		`'; DROP TABLE person;--`,
		`" OR 1=1 --`,
		`* AND ()`,
		`NEAR(`,
		`"`,
		`%`,
		`_`,
	}
	for _, q := range queries {
		// 不应 panic 或返回错误
		if _, _, err := st.Search(ctx, q, "", 10, 0); err != nil {
			t.Errorf("检索恶意输入 %q 出错: %v", q, err)
		}
	}
	// 表必须仍然存在且内容完好
	want := realPeopleCount(t)
	var n int
	if err := st.DB().QueryRowContext(ctx, "SELECT COUNT(*) FROM person").Scan(&n); err != nil {
		t.Fatalf("person 表已损坏: %v", err)
	}
	if n != want {
		t.Errorf("person 表行数 = %d, 期望 %d", n, want)
	}
}

// TestSearchEmpty 确认空查询返回空结果而非报错。
func TestSearchEmpty(t *testing.T) {
	t.Parallel()
	st := newTestStore(t)

	for _, q := range []string{"", "   ", "\t"} {
		hits, total, err := st.Search(t.Context(), q, "", 10, 0)
		if err != nil {
			t.Fatalf("空查询 %q 出错: %v", q, err)
		}
		if total != 0 || len(hits) != 0 {
			t.Errorf("空查询 %q 应返回 0 条", q)
		}
	}
}

// TestSearchPagination 确认分页与总数统计一致。
func TestSearchPagination(t *testing.T) {
	t.Parallel()
	st := newTestStore(t)
	ctx := t.Context()

	_, total, err := st.Search(ctx, "革命", "", 100, 0)
	if err != nil {
		t.Fatalf("检索出错: %v", err)
	}
	if total < 2 {
		t.Skipf("命中数不足以测试分页: %d", total)
	}

	page1, _, err := st.Search(ctx, "革命", "", 1, 0)
	if err != nil {
		t.Fatalf("首页检索出错: %v", err)
	}
	if len(page1) != 1 {
		t.Fatalf("首页应有 1 条，实际 %d", len(page1))
	}

	page2, _, err := st.Search(ctx, "革命", "", 1, 1)
	if err != nil {
		t.Fatalf("次页检索出错: %v", err)
	}
	if len(page2) == 1 && page2[0].Ref == page1[0].Ref && page2[0].Title == page1[0].Title {
		t.Error("分页第二页与第一页返回了同一条结果")
	}
}

// TestIngestCounts 确认入库行数与内容源一致。
func TestIngestCounts(t *testing.T) {
	t.Parallel()
	st := newTestStore(t)

	counts, err := st.Counts(t.Context())
	if err != nil {
		t.Fatalf("统计: %v", err)
	}
	if want := realPeopleCount(t); counts.People != want {
		t.Errorf("人物数 = %d, 期望 %d", counts.People, want)
	}
	if counts.Works < 10 {
		t.Errorf("著作数 = %d, 期望 >= 10", counts.Works)
	}
	if counts.Events < 15 {
		t.Errorf("事件数 = %d, 期望 >= 15", counts.Events)
	}
	if counts.Terms < 10 {
		t.Errorf("术语数 = %d, 期望 >= 10", counts.Terms)
	}
}

// TestPersonDetail 确认人物详情的四域章节完整。
func TestPersonDetail(t *testing.T) {
	t.Parallel()
	st := newTestStore(t)
	ctx := t.Context()

	for _, slug := range []string{"marx", "engels", "plekhanov", "zetkin", "katayama", "lenin", "luxemburg", "kollontai", "li-dazhao", "ho-chi-minh", "gramsci", "mao", "castro", "guevara", "sankara"} {
		p, err := st.Person(ctx, slug)
		if err != nil {
			t.Fatalf("查询人物 %s: %v", slug, err)
		}
		seen := map[string]bool{}
		for _, s := range p.Sections {
			seen[s.Kind] = true
		}
		for _, k := range []string{"politics", "economy", "culture", "thought", "controversy"} {
			if !seen[k] {
				t.Errorf("人物 %s 缺少 %s 章节", slug, k)
			}
		}
		if len(p.Timeline) == 0 {
			t.Errorf("人物 %s 缺少生平年表", slug)
		}
		if len(p.Works) == 0 {
			t.Errorf("人物 %s 没有关联著作", slug)
		}
	}
}

// TestNotFound 确认查询不存在条目返回 ErrNotFound。
func TestNotFound(t *testing.T) {
	t.Parallel()
	st := newTestStore(t)
	ctx := t.Context()

	if _, err := st.Person(ctx, "nobody"); err == nil {
		t.Error("查询不存在的人物应报错")
	} else if !strings.Contains(err.Error(), "未找到") {
		t.Errorf("错误信息不符合预期: %v", err)
	}
	if _, err := st.Work(ctx, "nobody"); err == nil {
		t.Error("查询不存在的著作应报错")
	}
	if _, err := st.Term(ctx, "nobody"); err == nil {
		t.Error("查询不存在的术语应报错")
	}
}

// TestEventFilters 确认事件筛选生效。
func TestEventFilters(t *testing.T) {
	t.Parallel()
	st := newTestStore(t)
	ctx := t.Context()

	all, err := st.Events(ctx, "", "")
	if err != nil {
		t.Fatalf("查询全部事件: %v", err)
	}
	filtered, err := st.Events(ctx, "mao", "")
	if err != nil {
		t.Fatalf("按人物筛选事件: %v", err)
	}
	if len(filtered) == 0 || len(filtered) >= len(all) {
		t.Errorf("按人物筛选后数量异常: %d / %d", len(filtered), len(all))
	}

	byCat, err := st.Events(ctx, "", "理论")
	if err != nil {
		t.Fatalf("按范畴筛选事件: %v", err)
	}
	if len(byCat) == 0 {
		t.Error("按范畴筛选应至少有一条")
	}
}

// TestGraph 确认关系图数据完整。
func TestGraph(t *testing.T) {
	t.Parallel()
	st := newTestStore(t)

	nodes, edges, err := st.Graph(t.Context())
	if err != nil {
		t.Fatalf("查询关系图: %v", err)
	}
	if want := realPeopleCount(t); len(nodes) != want {
		t.Errorf("节点数 = %d, 期望 %d", len(nodes), want)
	}
	if len(edges) == 0 {
		t.Error("关系边不应为空")
	}
	known := map[string]bool{}
	for _, n := range nodes {
		known[n.Slug] = true
	}
	for _, e := range edges {
		if !known[e.From] || !known[e.To] {
			t.Errorf("关系边引用了未知节点: %s -> %s", e.From, e.To)
		}
	}
}

// TestTermsRelation 确认术语互链与出处解析正确。
func TestTermsRelation(t *testing.T) {
	t.Parallel()
	st := newTestStore(t)
	ctx := t.Context()

	tm, err := st.Term(ctx, "proletariat")
	if err != nil {
		t.Fatalf("查询术语: %v", err)
	}
	if len(tm.People) == 0 {
		t.Error("无产阶级术语应关联人物")
	}
	if len(tm.Related) == 0 {
		t.Error("无产阶级术语应关联其他术语")
	}
	if len(tm.Sources) == 0 {
		t.Error("无产阶级术语应关联出处著作")
	}
}

// TestWorkFilter 确认著作按人物与标题筛选。
func TestWorkFilter(t *testing.T) {
	t.Parallel()
	st := newTestStore(t)
	ctx := t.Context()

	all, err := st.Works(ctx, "", "")
	if err != nil {
		t.Fatalf("查询著作: %v", err)
	}
	one, err := st.Works(ctx, "lenin", "")
	if err != nil {
		t.Fatalf("按人物查询著作: %v", err)
	}
	if len(one) == 0 || len(one) >= len(all) {
		t.Errorf("按人物筛选著作数量异常: %d / %d", len(one), len(all))
	}
	byTitle, err := st.Works(ctx, "", "资本论")
	if err != nil {
		t.Fatalf("按标题查询著作: %v", err)
	}
	if len(byTitle) != 1 {
		t.Errorf("按标题'资本论'应命中 1 条，实际 %d", len(byTitle))
	}
}

// 该测试保证元信息计数接口可用。
func TestCounts(t *testing.T) {
	t.Parallel()
	st := newTestStore(t)

	c, err := st.Counts(t.Context())
	if err != nil {
		t.Fatalf("统计: %v", err)
	}
	if c.People+c.Works+c.Events+c.Terms == 0 {
		t.Error("统计结果不应全为 0")
	}
}
