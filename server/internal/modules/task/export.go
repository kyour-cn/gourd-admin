package task

import (
	"context"
	"errors"

	"app/internal/modules/task/export"
	"app/internal/orm/model"
)

func ExportTask(ctx context.Context, task *model.Task) error {
	switch task.Label {
	case "user":
		// 用户导出任务
		err := export.UserExport(ctx, task)
		if err != nil {
			return err
		}
	default:
		// 未知类型，改为失败
		return errors.New("未知导出类型")
	}

	// 导出任务
	return nil
}
