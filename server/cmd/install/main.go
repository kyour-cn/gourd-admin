package main

import (
	"app/internal/initialize"
	"app/internal/orm/query"
	"log/slog"
	"os"
	"strings"
)

func main() {
	// 初始化数据库
	err := initialize.InitDatabase()
	if err != nil {
		panic(err)
	}

	// 导入初始化数据
	err = InportSql("assets/migrations/schema.sql")
	if err != nil {
		panic(err)
	}
	// 导入初始化数据
	err = InportSql("assets/migrations/seed.sql")
	if err != nil {
		panic(err)
	}

	slog.Info("初始化数据导入完成")
}

func InportSql(file string) error {
	slog.Info("导入初始化数据", "file", file)
	// 导入初始化数据
	sqlBytes, _ := os.ReadFile(file)
	sqlStr := string(sqlBytes)

	db := query.Q.App.UnderlyingDB()

	// 按 ; 分割多条 SQL（简单方式）
	queries := strings.Split(sqlStr, ";")

	for _, q := range queries {
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}

		// 执行SQL语句
		tx := db.Exec(q)
		// 检查是否有错误
		if tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}
