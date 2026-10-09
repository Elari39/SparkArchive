package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
