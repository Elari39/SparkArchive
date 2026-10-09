// Package api 提供只读的 JSON HTTP 接口。
package api

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/sparkarchive/server/internal/content"
	"github.com/sparkarchive/server/internal/store"
)

// 检索分页参数边界（原先散落在 store 层的限制统一收到 handler）。
const (
	defaultLimit   = 50
	maxSearchLimit = 200
	staticAssetTTL = "public, max-age=31536000, immutable"
	staticNoCache  = "no-cache"
)

// clampLimit 把 limit 约束到 (0, maxSearchLimit]，非法值回落到 defaultLimit。
func clampLimit(n int) int {
	if n <= 0 || n > maxSearchLimit {
		return defaultLimit
	}
	return n
}

// Server 持有依赖并实现 http.Handler。
type Server struct {
	store  *store.Store
	site   content.Site
	log    *slog.Logger
	static fs.FS // 已构建的前端资源，可为 nil
}

// New 构造 API 服务。static 为 nil 时不提供静态资源（开发模式）。
func New(st *store.Store, site content.Site, log *slog.Logger, static fs.FS) *Server {
	return &Server{store: st, site: site, log: log, static: static}
}

// Handler 返回装配好的路由。
// 使用 Go 1.22+ 的方法感知 ServeMux 模式与 r.PathValue 取路径参数。
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/meta", s.handleMeta)
	mux.HandleFunc("GET /api/people", s.handlePeople)
	mux.HandleFunc("GET /api/people/{slug}", s.handlePerson)
	mux.HandleFunc("GET /api/works", s.handleWorks)
	mux.HandleFunc("GET /api/works/{slug}", s.handleWork)
	mux.HandleFunc("GET /api/events", s.handleEvents)
	mux.HandleFunc("GET /api/terms", s.handleTerms)
	mux.HandleFunc("GET /api/terms/{slug}", s.handleTerm)
	mux.HandleFunc("GET /api/relations", s.handleRelations)
	mux.HandleFunc("GET /api/search", s.handleSearch)

	// 其余全部交给前端 SPA（含客户端路由回退）
	mux.HandleFunc("GET /", s.handleStatic)

	// 顺序：限流最外层 → 日志 → recovery 最内层。
	// 这样超限请求也会被日志记录；recovery 借 statusRecorder 判断头是否已写。
	return s.withRateLimit(s.withLogging(s.withRecovery(mux)))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// getJSON 统一 handler 骨架：取数据 → 失败走 fail → 成功序列化。
// 返回 any 便于各 handler 直接返回各自的 DTO 或 map。
func (s *Server) getJSON(w http.ResponseWriter, r *http.Request, load func(context.Context) (any, error)) {
	v, err := load(r.Context())
	if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (s *Server) handleMeta(w http.ResponseWriter, r *http.Request) {
	s.getJSON(w, r, func(ctx context.Context) (any, error) {
		counts, err := s.store.Counts(ctx)
		if err != nil {
			return nil, err
		}
		return map[string]any{"site": s.site, "counts": counts}, nil
	})
}

func (s *Server) handlePeople(w http.ResponseWriter, r *http.Request) {
	s.getJSON(w, r, func(ctx context.Context) (any, error) {
		return s.store.Persons(ctx)
	})
}

func (s *Server) handlePerson(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	s.getJSON(w, r, func(ctx context.Context) (any, error) {
		return s.store.Person(ctx, slug)
	})
}

func (s *Server) handleWorks(w http.ResponseWriter, r *http.Request) {
	person, q := r.URL.Query().Get("person"), r.URL.Query().Get("q")
	s.getJSON(w, r, func(ctx context.Context) (any, error) {
		return s.store.Works(ctx, person, q)
	})
}

func (s *Server) handleWork(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	s.getJSON(w, r, func(ctx context.Context) (any, error) {
		return s.store.Work(ctx, slug)
	})
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	person, category := r.URL.Query().Get("person"), r.URL.Query().Get("category")
	s.getJSON(w, r, func(ctx context.Context) (any, error) {
		return s.store.Events(ctx, person, category)
	})
}

func (s *Server) handleTerms(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	s.getJSON(w, r, func(ctx context.Context) (any, error) {
		return s.store.Terms(ctx, q)
	})
}

func (s *Server) handleTerm(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	s.getJSON(w, r, func(ctx context.Context) (any, error) {
		return s.store.Term(ctx, slug)
	})
}

func (s *Server) handleRelations(w http.ResponseWriter, r *http.Request) {
	s.getJSON(w, r, func(ctx context.Context) (any, error) {
		nodes, edges, err := s.store.Graph(ctx)
		if err != nil {
			return nil, err
		}
		return map[string]any{"nodes": nodes, "edges": edges}, nil
	})
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	kind := r.URL.Query().Get("type")
	limit := clampLimit(atoiDefault(r.URL.Query().Get("limit"), defaultLimit))
	offset := max(atoiDefault(r.URL.Query().Get("offset"), 0), 0)

	s.getJSON(w, r, func(ctx context.Context) (any, error) {
		hits, total, err := s.store.Search(ctx, q, kind, limit, offset)
		if err != nil {
			return nil, err
		}
		return map[string]any{"total": total, "items": hits}, nil
	})
}

// handleStatic 提供内嵌的前端资源，并为客户端路由做 index.html 回退。
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	if s.static == nil {
		http.NotFound(w, r)
		return
	}
	name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
	if name == "" || name == "." {
		name = "index.html"
	}
	data, err := fs.ReadFile(s.static, name)
	if err != nil {
		// 未命中的路径交给 SPA 路由处理
		data, err = fs.ReadFile(s.static, "index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		name = "index.html"
	}
	if ct := mime.TypeByExtension(path.Ext(name)); ct != "" {
		w.Header().Set("Content-Type", ct)
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	// 带内容哈希的资源可长期缓存；index.html 不缓存，保证发版即时生效
	if strings.HasPrefix(name, "assets/") {
		w.Header().Set("Cache-Control", staticAssetTTL)
	} else {
		w.Header().Set("Cache-Control", staticNoCache)
	}
	if _, err := w.Write(data); err != nil {
		s.log.WarnContext(r.Context(), "写入响应失败", "path", r.URL.Path, "err", err)
	}
}

// fail 统一错误响应：未找到映射为 404，其余为 500。
func (s *Server) fail(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, store.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "未找到"})
		return
	}
	s.log.ErrorContext(r.Context(), "请求处理失败", "path", r.URL.Path, "err", err)
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "服务器内部错误"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// 此处无法再改动响应头，编码失败只能记录
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Warn("序列化 JSON 响应失败", "err", err)
	}
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
