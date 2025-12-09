package export

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"app/internal/http/admin/dto"
	"app/internal/orm/model"
	"app/internal/orm/query"
	"app/internal/util/excel"

	"gorm.io/gen"
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

	q := query.User
	var conds []gen.Condition

	// 关键词搜索
	if params.Keyword != "" {
		conds = append(conds, q.Where(
			q.Where(q.Username.Like("%"+params.Keyword+"%")).
				Or(q.Nickname.Like("%"+params.Keyword+"%")),
		))
	}

	// 分页查询用户数据
	pageSize := 1000
	page := 1

	// 分页处理数据，直到没有更多数据
	for {
		offset := (page - 1) * pageSize
		users, err := q.WithContext(ctx).
			Where(conds...).
			Limit(pageSize).
			Offset(offset).
			Find()
		if err != nil {
			return err
		}

		// 如果没有数据，停止分页
		if len(users) == 0 {
			break
		}

		slog.Debug("export user page", "page", page, "count", len(users))

		// 写入当前页的用户数据
		for _, user := range users {
			loginTime := "未登录"
			if !user.LoginTime.IsZero() {
				loginTime = user.LoginTime.Format("2006-01-02 15:04:05")
			}
			_ = loginTime
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
			user = nil
		}

		// 如果返回的数据少于pageSize，说明是最后一页
		if len(users) < pageSize {
			break
		}
		page++
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
