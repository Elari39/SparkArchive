package main

import (
	"fmt"
	"os"
)

// main 是进程入口：先分发子命令（`archive ingest` 走构建期入库），
// 其余一律走默认 serve 流程。
//
// 注意：子命令分发**不能**放在 init 里用 os.Exit —— 那会绕过所有 defer
// （ingest 路径下 `defer st.Close()` 将永不执行，数据库句柄与 WAL 文件无法正常收尾）。
func main() {
	if err := dispatch(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "启动失败: %v\n", err)
		os.Exit(1)
	}
}

func dispatch(args []string) error {
	if len(args) > 0 && args[0] == "ingest" {
		if err := runIngest(args[1:]); err != nil {
			return fmt.Errorf("入库失败: %w", err)
		}
		return nil
	}
	return run()
}
