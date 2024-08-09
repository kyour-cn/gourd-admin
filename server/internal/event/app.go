package event

import (
	"context"
	"github.com/go-gourd/gourd/event"
	"gourd/internal/cmd"
	"gourd/internal/cron"
	"gourd/internal/http/router"
	"gourd/internal/init"
	"gourd/internal/util"
	"gourd/internal/util/redisutil"
	"log/slog"
)

// Register 事件注册
func Register(ctx context.Context) {

	// Boot事件(应用) -初始化应用时执行
	event.Listen("app.boot", func(ctx context.Context) {
		slog.Debug("boot event.")

		err := util.InitLog()
		if err != nil {
			panic(err)
		}

		err = init.DatabaseInit()
		if err != nil {
			panic(err)
		}

		_, err = redisutil.InitRedis(ctx)
		if err != nil {
			panic(err)
		}

		// 注册命令行
		cmd.Register()
	})

	// Init事件(应用) -初始化完成执行
	event.Listen("app.init", func(context.Context) {
		slog.Debug("init event.")

		// 注册定时任务
		cron.Register()

		// 注册路由
		router.Register()
	})

	// Start事件(应用) -启动后执行
	event.Listen("app.start", func(context.Context) {
		slog.Debug("start event.")

		// 启动Http服务
		router.StartServer()
	})

	// Stop事件(应用) -终止时执行
	event.Listen("app.stop", func(context.Context) {
		slog.Debug("stop event.")
	})

	// 注册更多自定义事件监听

}
