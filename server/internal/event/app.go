package event

import (
	"context"
	"log/slog"

	"github.com/go-gourd/gourd/event"

	"app/internal/initialize"
	"app/internal/modules/task"
)

// AppEvent 事件注册
func AppEvent(_ context.Context) {

	// Boot事件(应用) -初始化应用时执行
	event.Listen("app.boot", func(ctx context.Context) {
		// 初始化一些全局配置、工具等
		err := initialize.InitCommon(ctx)
		if err != nil {
			panic(err)
		}

		slog.Debug("boot event.")
	})

	// Init事件(应用) -初始化完成执行
	event.Listen("app.init", func(context.Context) {
		slog.Debug("init event.")

		// 初始化命令行并解析参数
		initialize.InitCmd()
	})

	// Start事件(应用) -启动后执行
	event.Listen("app.start", func(ctx context.Context) {
		slog.Debug("start event.")

		// 初始化定时任务
		initialize.InitCron()

		// 初始化Http服务
		initialize.InitHttpServer()

		// 初始化异步任务
		task.Init(ctx)
	})

	// Stop事件(应用) -停止时执行
	event.Listen("app.stop", func(context.Context) {
		slog.Debug("stop event.")
	})

}
