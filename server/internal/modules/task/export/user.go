package export

import (
	"context"
	"encoding/json"
	"fmt"

	"app/internal/http/admin/dto"
	"app/internal/orm/model"
	"app/internal/orm/query"
	"app/internal/util/excel"
)

func UserExport(ctx context.Context, task *model.Task) error {
	params := &dto.UserExportReq{}
	err := json.Unmarshal([]byte(task.Content), params)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("uploads/export/export_user_%d.xlsx", task.ID)
	e := excel.NewExcel()
	defer func(e *excel.Excel) {
		_ = e.Close()
	}(e)

	err = e.SetCols([]excel.Column{
		{Name: "ID", Width: 10},
		{Name: "昵称", Width: 20},
		{Name: "用户名", Width: 20},
		{Name: "状态", Width: 10},
		{Name: "注册时间", Width: 20},
		{Name: "最后登录时间", Width: 20},
	}, 1)
	if err != nil {
		return err
	}

	// 从数据库查询用户数据
	users, err := query.User.WithContext(ctx).
		Find()
	if err != nil {
		return err
	}

	// 写入用户数据
	for _, user := range users {
		loginTime := "未登录"
		if !user.LoginTime.IsZero() {
			loginTime = user.LoginTime.Format("2006-01-02 15:04:05")
		}

		err = e.WriteLine(e.CurrentRow, []any{
			user.ID,
			user.Nickname,
			user.Username,
			user.Status,
			user.CreatedAt.Format("2006-01-02 15:04:05"),
			loginTime,
		})
		if err != nil {
			return err
		}
	}

	err = e.Save("web/" + fileName)
	if err != nil {
		return err
	}

	result, err := json.Marshal(map[string]string{
		"file": fileName,
	})

	// 更新任务状态为完成
	_, err = query.Task.WithContext(ctx).
		Where(query.Task.ID.Eq(task.ID)).
		Updates(&model.Task{
			Status: 2,
			Result: string(result),
		})
	if err != nil {
		return err
	}

	// 用户导出任务
	return nil
}
