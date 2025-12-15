package task

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-gourd/gourd/event"

	"app/internal/orm/model"
	"app/internal/orm/query"
)

func Init(ctx context.Context) {
	go func() {
		// 监听任务运行事件
		event.Listen("task.run", func(_ context.Context) {
			go Run(ctx)
		})
		for {
			select {
			case <-ctx.Done():
				slog.Debug("task stop.")
				return
			default:
				Run(ctx)
				// 间隔执行
				time.Sleep(time.Second * 10)
			}
		}
	}()
}

var runLock = make(chan struct{}, 1)

// Run 运行任务
func Run(ctx context.Context) {
	// 加锁，防止并发执行
	select {
	case runLock <- struct{}{}:
	default:
		return
	}

	defer func() {
		// 解锁
		<-runLock

		// 处理panic
		if err := recover(); err != nil {
			slog.Error("task run error.", "err", err)
		}
	}()

	// 查询task待处理列表
	tasks, err := query.Task.WithContext(ctx).
		Where(query.Task.Status.Eq(0)).
		Find()
	if err != nil {
		slog.Error("query task error.", "err", err)
		return
	}
	// 遍历任务列表
	for _, v := range tasks {
		// 处理任务
		err := handleTask(ctx, v)
		if err != nil {
			slog.Error("handle task error.", "err", err)
		}
	}

}

func handleTask(ctx context.Context, task *model.Task) error {
	q := query.Task
	// 更新任务状态为处理中
	res, err := q.WithContext(ctx).
		Where(q.ID.Eq(task.ID)).
		Update(q.Status, 1)
	if err != nil && res.RowsAffected == 0 {
		return err
	}

	switch task.Type {
	case "export":
		// 导出任务
		err := ExportTask(ctx, task)
		if err != nil {
			result := err.Error()
			// 导出任务失败，更新任务状态为失败
			_, err = q.WithContext(ctx).
				Where(q.ID.Eq(task.ID)).
				Updates(&model.Task{
					Status: -1,
					Result: &result,
				})
			return err
		}
	default:
		result := "未知类型"
		// 未知类型，改为失败
		_, err = q.WithContext(ctx).
			Where(q.ID.Eq(task.ID)).
			Updates(&model.Task{
				Status: -1,
				Result: &result,
			})
		if err != nil {
			return err
		}
	}

	return nil
}
