package api

import (
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

// 简易令牌桶限流：按客户端 IP 维护一个容量为 burst、每秒补充 rate 个令牌的桶。
//
// 目标是保护 CPU 密集的降级检索路径（2 字中文查询走全表 LIKE 扫描），
// 而非精确额度控制；数据只读、访问量小，故实现从简。
type rateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rate    float64 // 每秒补充令牌数
	burst   float64 // 桶容量
	now     func() time.Time
}

type bucket struct {
	tokens float64
	last   time.Time
}

func newRateLimiter(rate, burst float64) *rateLimiter {
	return &rateLimiter{
		buckets: make(map[string]*bucket),
		rate:    rate,
		burst:   burst,
		now:     time.Now,
	}
}

// allow 判定 key 当前是否有可用令牌，并惰性清理过期的桶。
func (l *rateLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{tokens: l.burst, last: now}
		l.buckets[key] = b
	}
	// 按流逝时间补充
	elapsed := now.Sub(b.last).Seconds()
	if elapsed > 0 {
		b.tokens = min(l.burst, b.tokens+elapsed*l.rate)
		b.last = now
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// withRateLimit 对全部请求做 IP 级限流；超限返回 429。
func (s *Server) withRateLimit(next http.Handler) http.Handler {
	// 平均 20 req/s，突发 60：足够正常浏览，又能压住脚本化滥用。
	limiter := newRateLimiter(20, 60)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.allow(clientIP(r)) {
			s.log.WarnContext(r.Context(), "请求被限流", "path", r.URL.Path, "ip", clientIP(r))
			writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "请求过于频繁，请稍后重试"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// clientIP 取客户端地址；容器内通常由 Docker 端口代理直连，RemoteAddr 即可。
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// statusRecorder 记录响应状态码与「是否已写响应头」，供访问日志与 panic 兜底使用。
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(code int) {
	if r.wroteHeader {
		// 已写过响应头则忽略重复调用，避免 "superfluous WriteHeader" 与响应体错乱。
		return
	}
	r.wroteHeader = true
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Write 隐式触发 WriteHeader(200)，需同步记录状态。
func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	return r.ResponseWriter.Write(b)
}

func (s *Server) withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		s.log.InfoContext(r.Context(), "http",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rec.status),
			slog.Duration("dur", time.Since(start)),
		)
	})
}

// withRecovery 兜住 panic，避免单个请求打挂整个进程。
// 若响应头已写出（handler 中途 panic），则不再尝试写 JSON，只记录日志。
func (s *Server) withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				s.log.ErrorContext(r.Context(), "请求 panic", "path", r.URL.Path, "panic", rec)
				if sr, ok := w.(*statusRecorder); ok && sr.wroteHeader {
					// 头部已发送，无法再改状态码或响应体，交由连接层处理
					return
				}
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "服务器内部错误"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
