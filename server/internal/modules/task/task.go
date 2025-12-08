package task

import (
	"context"
	"log/slog"
	"time"

	"app/internal/orm/model"
	"app/internal/orm/query"
)

func Init(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Debug("task stop.")
				return
			default:
				Run(ctx)
				// 间隔1秒执行
				time.Sleep(time.Second)
			}
		}
	}()
}

// Run 运行任务
func Run(ctx context.Context) {
	defer func() {
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
	_, err := q.WithContext(ctx).
		Where(q.ID.Eq(task.ID)).
		Update(q.Status, 1)
	if err != nil {
		return err
	}

	switch task.Type {
	case "export":
		// 导出任务
		err := ExportTask(ctx, task)
		if err != nil {
			// 导出任务失败，更新任务状态为失败
			_, _ = q.WithContext(ctx).
				Where(q.ID.Eq(task.ID)).
				Updates(&model.Task{
					Status: -1,
					Result: err.Error(),
				})
			return err
		}
	default:
		// 未知类型，改为失败
		_, _ = q.WithContext(ctx).
			Where(q.ID.Eq(task.ID)).
			Updates(&model.Task{
				Status: -1,
				Result: "未知类型",
			})
	}

	return nil
}
