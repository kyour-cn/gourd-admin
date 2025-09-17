package system

import (
	"app/internal/http/admin/dto"
	"app/internal/http/admin/services"
	"app/internal/http/common/controller"
	"net/http"
)

// Log 用户控制器
type Log struct {
	controller.Base //继承基础控制器
}

// TypeList 日志类型列表
func (c *Log) TypeList(w http.ResponseWriter, r *http.Request) {
	req := &dto.LogTypeListReq{}
	if err := c.QueryReqUnmarshal(r, req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常："+err.Error(), "")
		return
	}

	service := services.NewLogTypeService(r.Context())
	res, err := service.GetTypeList(req)
	if err != nil {
		_ = c.Fail(w, 500, "获取列表失败", err.Error())
		return
	}
	_ = c.Success(w, "", res)
}

// List 日志列表
func (c *Log) List(w http.ResponseWriter, r *http.Request) {
	req := &dto.LogListReq{}
	if err := c.QueryReqUnmarshal(r, req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常："+err.Error(), "")
		return
	}

	service := services.NewLogTypeService(r.Context())
	res, err := service.GetList(req)
	if err != nil {
		_ = c.Fail(w, 500, "获取列表失败", err.Error())
		return
	}
	_ = c.Success(w, "", res)
}

// LogStat 日志统计
func (c *Log) LogStat(w http.ResponseWriter, r *http.Request) {
	req := &dto.LogStatReq{}
	if err := c.QueryReqUnmarshal(r, req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常："+err.Error(), "")
		return
	}

	service := services.NewLogTypeService(r.Context())
	res, err := service.GetLogStat(req)
	if err != nil {
		_ = c.Fail(w, 500, "获取列表失败", err.Error())
		return
	}
	_ = c.Success(w, "", res)
}
