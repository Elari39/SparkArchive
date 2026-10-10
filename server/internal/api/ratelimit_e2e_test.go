package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

// TestRateLimitReturns429 端到端验证：同一 IP 连续超量请求会拿到 429。
func TestRateLimitReturns429(t *testing.T) {
	t.Parallel()

	h := newTestServer(t) // 已装配完整中间件链

	got429 := false
	for i := 0; i < 400; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		req.RemoteAddr = "10.0.0.9:12345" // 固定 IP，避免随机源分散计数
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code == http.StatusTooManyRequests {
			got429 = true
			break
		}
	}
	if !got429 {
		t.Fatal("连续 400 次同 IP 请求未触发 429 限流")
	}
}

// TestRateLimitSkipsStaticAssets 固定住限流的覆盖面：只限 /api/，静态资源不受影响。
//
// 回归背景：限流最初套在整个 mux 上，于是正常浏览一个页面就会抽干突发额度
// （首页约 13 个请求，其中多数是肖像图），图片随机变成 429；
// smoke.ps1 的「人物肖像资产」整段因此全线失败。
func TestRateLimitSkipsStaticAssets(t *testing.T) {
	t.Parallel()

	static := fstest.MapFS{
		"index.html":               &fstest.MapFile{Data: []byte("<!doctype html><div id=\"app\"></div>")},
		"assets/portraits/mao.jpg": &fstest.MapFile{Data: []byte("jpeg-bytes")},
	}
	h := newTestServerWithStatic(t, static)

	const ip = "10.0.0.11:12345"

	// 先把该 IP 在 /api/ 上的令牌桶抽干，作为测试前提
	drained := false
	for i := 0; i < 400; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
		req.RemoteAddr = ip
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code == http.StatusTooManyRequests {
			drained = true
			break
		}
	}
	if !drained {
		t.Fatal("未能抽干 /api/ 的令牌桶，测试前提不成立")
	}

	// 同一 IP 访问静态资源与 SPA 外壳，仍应正常返回
	for _, tc := range []struct {
		path string
		want int
	}{
		{"/assets/portraits/mao.jpg", http.StatusOK},
		{"/", http.StatusOK},
		{"/people/mao", http.StatusOK}, // 客户端路由回退
	} {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		req.RemoteAddr = ip
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code == http.StatusTooManyRequests {
			t.Errorf("%s 被限流：静态资源不应计入 /api/ 的令牌桶", tc.path)
			continue
		}
		if rec.Code != tc.want {
			t.Errorf("%s 期望 %d，实际 %d", tc.path, tc.want, rec.Code)
		}
	}
}
