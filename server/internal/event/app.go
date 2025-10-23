package event

import (
	"app/internal/initialize"
	"context"
	"log/slog"

	"github.com/go-gourd/gourd/event"
)

// AppEvent 事件注册
func AppEvent(_ context.Context) {

	// Boot事件(应用) -初始化应用时执行
	event.Listen("app.boot", func(ctx context.Context) {
		slog.Debug("boot event.")

		// 初始化一些全局配置、工具等
		initialize.InitCommon(ctx)

		// 初始化数据库
		err := initialize.InitDatabase()
		if err != nil {
			panic(err)
		}

		// 初始化命令行
		initialize.InitCmd()
	})

	// Init事件(应用) -初始化完成执行
	event.Listen("app.init", func(context.Context) {
		slog.Debug("init event.")
	})

	// Start事件(应用) -启动后执行
	event.Listen("app.start", func(context.Context) {
		slog.Debug("start event.")

		// 初始化定时任务
		initialize.InitCron()

		// 初始化Http服务
		initialize.InitHttpServer()
	})

	// Stop事件(应用) -停止时执行
	event.Listen("app.stop", func(context.Context) {
		slog.Debug("stop event.")
	})

}
