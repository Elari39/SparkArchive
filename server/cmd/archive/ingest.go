package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/sparkarchive/server/internal/content"
	"github.com/sparkarchive/server/internal/store"
)

// runIngest 实现 `archive ingest` 子命令：把 content/ 编译为 SQLite 数据库。
// 构建镜像时调用，运行时不再需要 content/ 目录。
func runIngest(args []string) error {
	fs := flag.NewFlagSet("ingest", flag.ExitOnError)
	contentDir := fs.String("content", "content", "内容目录")
	dbPath := fs.String("db", "/data/archive.db", "输出数据库路径")
	if err := fs.Parse(args); err != nil {
		return err
	}

	archive, err := content.Load(*contentDir)
	if err != nil {
		return err
	}

	// 每次都从零重建，保证构建可重复
	if err := os.Remove(*dbPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("清理旧数据库: %w", err)
	}

	st, err := store.Open(*dbPath)
	if err != nil {
		return err
	}
	defer st.Close()

	if err := st.Migrate(); err != nil {
		return err
	}
	if err := st.Ingest(context.Background(), archive); err != nil {
		return err
	}

	counts, err := st.Counts(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("入库完成 %s: 人物=%d 著作=%d 事件=%d 术语=%d\n",
		*dbPath, counts.People, counts.Works, counts.Events, counts.Terms)
	return nil
}
