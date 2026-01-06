package controller

import (
	"app/internal/http/common/services"
	"net/http"
)

// Site 用户控制器
type Site struct {
	Base //继承基础控制器
}

func (c Site) Config(w http.ResponseWriter, r *http.Request) {
	serv := services.NewConfigService(r.Context())

	data, err := serv.GetSite()
	if err != nil {
		_ = c.Fail(w, 1, "获取配置失败", err.Error())
		return
	}
	_ = c.Success(w, "", data)
}
