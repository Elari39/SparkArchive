package store

import (
	"context"
	"database/sql"
	"fmt"
)

// queryRows 执行查询并把每一行扫进 []T。
//
// 它收敛了原先散落在 query.go 各处的「切片初始化 + rows.Next + Scan + rows.Err」
// 样板；返回的切片在无结果时是非 nil 空切片，保证 JSON 序列化为 [] 而非 null。
//
// scan 接收各列的扫描目标；务必按 SELECT 的列顺序传入。
func queryRows[T any](ctx context.Context, db *sql.DB, query string, scan func(*sql.Rows) (T, error), args ...any) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("执行查询: %w", err)
	}
	defer rows.Close()

	out := []T{}
	for rows.Next() {
		v, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历结果: %w", err)
	}
	return out, nil
}
