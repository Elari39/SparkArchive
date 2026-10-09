// Package store 负责把解析后的内容写入 SQLite，并提供只读查询。
package store

import (
	"database/sql"
	"embed"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaFS embed.FS

// Store 持有数据库句柄。
type Store struct {
	db *sql.DB

	// 数据在运行时只读不变，故各类条目数量可安全缓存。
	countsOnce sync.Once
	counts     Counts
	countsErr  error
}

// Open 打开（或创建）位于 path 的数据库。
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("打开数据库 %s: %w", path, err)
	}
	// 单连接：构建期写入无需并发，且避免 SQLite 的锁竞争。
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("连接数据库 %s: %w", path, err)
	}
	return &Store{db: db}, nil
}

// OpenReadOnly 以只读方式打开已构建的数据库（运行时使用）。
func OpenReadOnly(path string) (*Store, error) {
	dsn := fmt.Sprintf("file:%s?mode=ro&immutable=1", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("只读打开数据库 %s: %w", path, err)
	}
	db.SetMaxOpenConns(4)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("连接数据库 %s: %w", path, err)
	}
	return &Store{db: db}, nil
}

// Close 关闭数据库。
func (s *Store) Close() error { return s.db.Close() }

// DB 暴露底层句柄供包内与测试使用。
func (s *Store) DB() *sql.DB { return s.db }

// Schema 返回建表 SQL。
func Schema() (string, error) {
	b, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return "", fmt.Errorf("读取 schema: %w", err)
	}
	return string(b), nil
}

// Migrate 建表。要求目标是空库。
func (s *Store) Migrate() error {
	sqlText, err := Schema()
	if err != nil {
		return err
	}
	if _, err := s.db.Exec(sqlText); err != nil {
		return fmt.Errorf("执行 schema: %w", err)
	}
	return nil
}
