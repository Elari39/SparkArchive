package api

import (
	"testing"
	"time"
)

// TestRateLimiterBurstThenRefill 直接验证令牌桶语义：
// 突发容量内全部放行，持续调用被拒，冷却后恢复。
func TestRateLimiterBurstThenRefill(t *testing.T) {
	t.Parallel()

	now := time.Unix(0, 0)
	rl := newRateLimiter(20, 60)
	rl.now = func() time.Time { return now }

	// 前 60 次（burst 容量）应全部放行
	for i := 0; i < 60; i++ {
		if !rl.allow("1.2.3.4") {
			t.Fatalf("突发内第 %d 次被拒，期望放行", i+1)
		}
	}
	// 第 61 次应被拒
	if rl.allow("1.2.3.4") {
		t.Fatal("超过 burst 后仍放行，期望拒绝")
	}

	// 时间前进 1 秒 → 补充 20 个令牌 → 应能再放行 20 次
	now = now.Add(time.Second)
	for i := 0; i < 20; i++ {
		if !rl.allow("1.2.3.4") {
			t.Fatalf("补充后第 %d 次被拒，期望放行", i+1)
		}
	}
	if rl.allow("1.2.3.4") {
		t.Fatal("补充额度用尽后仍放行，期望拒绝")
	}

	// 不同 IP 各自独立计数
	if !rl.allow("5.6.7.8") {
		t.Fatal("不同 IP 应各有独立桶，期望放行")
	}
}
