package api

import (
	"encoding/json"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/sparkarchive/server/internal/content"
	"github.com/sparkarchive/server/internal/store"
)

// discardLogger 返回一个丢弃全部输出的日志器，避免测试噪声。
func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// newTestServer 用真实内容构建一个可用的 API 服务（不注入静态资源）。
func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	return newTestServerWithStatic(t, nil)
}

// newTestServerWithStatic 同上，但注入一个静态资源文件系统，用于验证静态资源的服务与回退行为。
func newTestServerWithStatic(t *testing.T, static fs.FS) http.Handler {
	t.Helper()
	ctx := t.Context()

	st, err := store.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatalf("打开数据库: %v", err)
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
	return New(st, a.Site, discardLogger(), static).Handler()
}

func get(t *testing.T, h http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// decodeMap 解析为泛型 map，用于断言 JSON 字段名。
func decodeMap(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &m); err != nil {
		t.Fatalf("解析响应失败: %v (body=%s)", err, rec.Body.String())
	}
	return m
}

func TestHealth(t *testing.T) {
	t.Parallel()
	h := newTestServer(t)
	rec := get(t, h, "/api/health")
	if rec.Code != http.StatusOK {
		t.Fatalf("状态码 = %d, 期望 200", rec.Code)
	}
}

// TestMetaContract 锁定 /api/meta 的 JSON 字段名，
// 前端 zod schema 依赖这些名字，改名必须同步修改两端。
func TestMetaContract(t *testing.T) {
	t.Parallel()
	h := newTestServer(t)

	m := decodeMap(t, get(t, h, "/api/meta"))

	site, ok := m["site"].(map[string]any)
	if !ok {
		t.Fatal("响应缺少 site 对象")
	}
	for _, key := range []string{"name", "name_en", "tagline", "subtitle", "description", "footer_note", "license_note"} {
		if _, ok := site[key]; !ok {
			t.Errorf("site 缺少字段 %q（前端 SiteMeta schema 需要）", key)
		}
	}
	counts, ok := m["counts"].(map[string]any)
	if !ok {
		t.Fatal("响应缺少 counts 对象")
	}
	for _, key := range []string{"people", "works", "events", "terms"} {
		if _, ok := counts[key]; !ok {
			t.Errorf("counts 缺少字段 %q", key)
		}
	}
}

// TestPersonDetailContract 锁定人物详情的字段名，特别是 terms 用 term 而非 name。
func TestPersonDetailContract(t *testing.T) {
	t.Parallel()
	h := newTestServer(t)

	m := decodeMap(t, get(t, h, "/api/people/marx"))
	for _, key := range []string{"slug", "name", "name_en", "birth", "death", "epithet",
		"portrait_key", "summary", "thesis", "timeline", "sections", "works", "terms"} {
		if _, ok := m[key]; !ok {
			t.Errorf("人物详情缺少字段 %q（前端 PersonDetail schema 需要）", key)
		}
	}

	terms, ok := m["terms"].([]any)
	if !ok || len(terms) == 0 {
		t.Fatal("terms 应为非空数组")
	}
	first, ok := terms[0].(map[string]any)
	if !ok {
		t.Fatal("terms[0] 应为对象")
	}
	if _, ok := first["term"]; !ok {
		t.Errorf("terms[0] 缺少 term 字段，实际: %v", first)
	}

	works, ok := m["works"].([]any)
	if !ok || len(works) == 0 {
		t.Fatal("works 应为非空数组")
	}
	w0 := works[0].(map[string]any)
	for _, key := range []string{"slug", "title", "year"} {
		if _, ok := w0[key]; !ok {
			t.Errorf("works[0] 缺少字段 %q", key)
		}
	}

	sections, ok := m["sections"].([]any)
	if !ok || len(sections) == 0 {
		t.Fatal("sections 应为非空数组")
	}
	s0 := sections[0].(map[string]any)
	for _, key := range []string{"kind", "title", "body_md", "ord"} {
		if _, ok := s0[key]; !ok {
			t.Errorf("sections[0] 缺少字段 %q", key)
		}
	}
}

// TestTermDetailContract 锁定术语详情中 related/sources 的字段名。
func TestTermDetailContract(t *testing.T) {
	t.Parallel()
	h := newTestServer(t)

	m := decodeMap(t, get(t, h, "/api/terms/proletariat"))
	for _, key := range []string{"slug", "term", "aliases", "definition", "body", "people", "related", "sources"} {
		if _, ok := m[key]; !ok {
			t.Errorf("术语详情缺少字段 %q（前端 TermDetail schema 需要）", key)
		}
	}

	related := m["related"].([]any)
	if _, ok := related[0].(map[string]any)["term"]; !ok {
		t.Errorf("related[0] 缺少 term 字段，实际: %v", related[0])
	}
	sources := m["sources"].([]any)
	if _, ok := sources[0].(map[string]any)["title"]; !ok {
		t.Errorf("sources[0] 缺少 title 字段，实际: %v", sources[0])
	}
	people := m["people"].([]any)
	if _, ok := people[0].(map[string]any)["name"]; !ok {
		t.Errorf("people[0] 缺少 name 字段，实际: %v", people[0])
	}
}

// TestEventContract 锁定事件字段与关联人物结构。
func TestEventContract(t *testing.T) {
	t.Parallel()
	h := newTestServer(t)

	rec := get(t, h, "/api/events")
	var events []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &events); err != nil {
		t.Fatalf("解析事件数组失败: %v", err)
	}
	if len(events) == 0 {
		t.Fatal("事件数组不应为空")
	}
	for _, key := range []string{"id", "year", "date_text", "title", "body", "location", "category", "people"} {
		if _, ok := events[0][key]; !ok {
			t.Errorf("事件缺少字段 %q", key)
		}
	}
	people := events[0]["people"].([]any)
	if len(people) == 0 {
		t.Fatal("事件应至少关联一个人物")
	}
	if _, ok := people[0].(map[string]any)["name"]; !ok {
		t.Errorf("事件关联人物缺少 name 字段，实际: %v", people[0])
	}
}

// TestSearchContract 锁定检索响应结构。
func TestSearchContract(t *testing.T) {
	t.Parallel()
	h := newTestServer(t)

	m := decodeMap(t, get(t, h, "/api/search?q=%E6%97%A0%E4%BA%A7%E9%98%B6%E7%BA%A7"))
	if _, ok := m["total"]; !ok {
		t.Error("检索响应缺少 total")
	}
	items, ok := m["items"].([]any)
	if !ok || len(items) == 0 {
		t.Fatal("items 应为非空数组")
	}
	for _, key := range []string{"kind", "ref", "title", "snippet", "person"} {
		if _, ok := items[0].(map[string]any)[key]; !ok {
			t.Errorf("检索结果缺少字段 %q", key)
		}
	}
}

// TestRelationsContract 锁定关系图结构。
func TestRelationsContract(t *testing.T) {
	t.Parallel()
	h := newTestServer(t)

	m := decodeMap(t, get(t, h, "/api/relations"))
	nodes, ok := m["nodes"].([]any)
	if !ok || len(nodes) != 15 {
		t.Fatalf("nodes 应为 15 个，实际 %v", m["nodes"])
	}
	n0 := nodes[0].(map[string]any)
	for _, key := range []string{"slug", "name", "epithet"} {
		if _, ok := n0[key]; !ok {
			t.Errorf("节点缺少字段 %q", key)
		}
	}
	edges := m["edges"].([]any)
	if len(edges) == 0 {
		t.Fatal("edges 不应为空")
	}
	e0 := edges[0].(map[string]any)
	for _, key := range []string{"from", "to", "kind", "label", "note"} {
		if _, ok := e0[key]; !ok {
			t.Errorf("边缺少字段 %q", key)
		}
	}
}

// TestEmptyArraysAreNotNull 确认空结果序列化为 [] 而非 null，避免前端解构报错。
func TestEmptyArraysAreNotNull(t *testing.T) {
	t.Parallel()
	h := newTestServer(t)

	// 检索无结果时 items 必须是 []
	rec := get(t, h, "/api/search?q=zzzznonexistentzzz")
	if body := rec.Body.String(); !strings.Contains(body, `"items":[]`) {
		t.Errorf("空检索结果应为 []，实际: %s", body)
	}

	// 按不存在的人物筛选著作应返回 []
	rec = get(t, h, "/api/works?person=nobody")
	if body := rec.Body.String(); body != "[]\n" {
		t.Errorf("空著作列表应为 []，实际: %q", body)
	}

	// 按不存在的范畴筛选事件应返回 []
	rec = get(t, h, "/api/events?category=nonexistent")
	if body := rec.Body.String(); body != "[]\n" {
		t.Errorf("空事件列表应为 []，实际: %q", body)
	}
}

// TestStaticAssetFallbackBehavior 锁定静态资源的两条行为边界：
//   - /assets/ 下不存在的资源必须 404，不得回退成 index.html，
//     否则请求方会拿到 200 + text/html 却按 image/* 解码（排查资源缺失时极易被误导）；
//   - 其余未知路径仍应回退到 index.html，保证客户端路由可直接访问。
func TestStaticAssetFallbackBehavior(t *testing.T) {
	t.Parallel()

	static := fstest.MapFS{
		"index.html":    &fstest.MapFile{Data: []byte("<!doctype html><div id=\"app\"></div>")},
		"assets/app.js": &fstest.MapFile{Data: []byte("console.log(1)")},
	}
	h := newTestServerWithStatic(t, static)

	rec := get(t, h, "/assets/app.js")
	if rec.Code != http.StatusOK {
		t.Errorf("存在的静态资源应 200，实际 %d", rec.Code)
	}

	rec = get(t, h, "/assets/portraits/missing.jpg")
	if rec.Code != http.StatusNotFound {
		t.Errorf("缺失的 /assets/ 资源应 404，实际 %d（body=%q）", rec.Code, rec.Body.String())
	}

	rec = get(t, h, "/people/missing-person")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `id="app"`) {
		t.Errorf("未知客户端路由应回退到 index.html，实际 %d", rec.Code)
	}
}

func TestNotFoundStatuses(t *testing.T) {
	t.Parallel()
	h := newTestServer(t)

	for _, path := range []string{"/api/people/nobody", "/api/works/nobody", "/api/terms/nobody"} {
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			if rec := get(t, h, path); rec.Code != http.StatusNotFound {
				t.Errorf("%s 状态码 = %d, 期望 404", path, rec.Code)
			}
		})
	}
}

func TestNoStaticFallsThroughTo404(t *testing.T) {
	t.Parallel()
	h := newTestServer(t)

	if rec := get(t, h, "/some/page"); rec.Code != http.StatusNotFound {
		t.Errorf("无静态资源时非 API 路径应 404，实际 %d", rec.Code)
	}
}
