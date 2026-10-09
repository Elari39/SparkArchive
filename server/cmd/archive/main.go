// Command archive 启动星火档案馆的 HTTP 服务。
//
// 数据库在构建期由 `archive ingest` 生成并烘焙进镜像，运行时以只读方式打开。
package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sparkarchive/server/internal/api"
	"github.com/sparkarchive/server/internal/content"
	"github.com/sparkarchive/server/internal/store"
)

//go:embed all:webdist
var webDist embed.FS

func run() error {
	var (
		addr        = flag.String("addr", envOr("ADDR", ":12026"), "监听地址")
		dbPath      = flag.String("db", envOr("DB_PATH", "/data/archive.db"), "数据库路径")
		sitePath    = flag.String("site", envOr("SITE_PATH", ""), "site.yaml 路径（可选，覆盖数据库内站点信息）")
		logLevel    = flag.String("log-level", envOr("LOG_LEVEL", "info"), "日志级别: debug|info|warn|error")
		shutdownTTL = flag.Duration("shutdown-timeout", 10*time.Second, "优雅关闭超时")
	)
	flag.Parse()

	log := newLogger(*logLevel)

	st, err := store.OpenReadOnly(*dbPath)
	if err != nil {
		return err
	}
	defer st.Close()

	// 站点信息优先从 site.yaml 读取（便于只改文案不发版），否则用内置默认值
	site, err := loadSite(*sitePath, log)
	if err != nil {
		return err
	}

	static, err := fs.Sub(webDist, "webdist")
	if err != nil {
		return fmt.Errorf("准备静态资源: %w", err)
	}
	if _, err := fs.Stat(static, "index.html"); err != nil {
		log.Warn("未找到内嵌前端资源，仅提供 API", "err", err)
		static = nil
	}

	srv := &http.Server{
		Addr:              *addr,
		Handler:           api.New(st, site, log, static).Handler(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// 信号驱动优雅关闭
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Info("星火档案馆已启动", "addr", *addr, "db", *dbPath)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("监听 %s: %w", *addr, err)
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Info("收到退出信号，正在优雅关闭")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), *shutdownTTL)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("优雅关闭失败: %w", err)
	}
	log.Info("已关闭")
	return nil
}

// loadSite 读取站点配置：显式路径 > 默认相对路径 > 内置兜底。
func loadSite(path string, log *slog.Logger) (content.Site, error) {
	for _, p := range []string{path, "content/site.yaml", "/app/content/site.yaml"} {
		if p == "" {
			continue
		}
		s, err := content.LoadSite(p)
		if err == nil {
			return s, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return content.Site{}, fmt.Errorf("解析站点配置 %s: %w", p, err)
		}
	}
	log.Warn("未找到 site.yaml，使用内置站点信息")
	return content.Site{
		Name:     "星火档案馆",
		NameEn:   "SPARK ARCHIVE",
		Tagline:  "全世界无产者，联合起来！",
		Subtitle: "无产阶级革命理论与实践的文献档案",
	}, nil
}

func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func newLogger(level string) *slog.Logger {
	var lv slog.Level
	if err := lv.UnmarshalText([]byte(level)); err != nil {
		lv = slog.LevelInfo
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lv}))
}
