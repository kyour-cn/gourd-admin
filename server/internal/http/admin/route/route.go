package route

import (
	"github.com/go-chi/chi/v5"
	"gourd/internal/http/admin/ctl"
	"gourd/internal/http/middleware"
)

// RegisterRoute 注册路由组
func RegisterRoute(r chi.Router) {

	r.Use(middleware.CorsMiddleware)

	// 注册登录相关路由
	authCtl := ctl.AuthCtl{}
	r.HandleFunc("/auth/captcha", authCtl.Captcha)
	r.HandleFunc("/auth/login", authCtl.Login)
	r.With(middleware.AuthJwtMiddleware).HandleFunc("/auth/menu", authCtl.GetMenu)

	// 注册app相关路由
	r.Mount("/app", chi.NewRouter().
		Group(func(r chi.Router) {
			appCtl := ctl.AppCtl{}
			r.Use(middleware.AuthJwtMiddleware)
			r.HandleFunc("/list", appCtl.List)
			r.HandleFunc("/add", appCtl.Add)
			r.HandleFunc("/edit", appCtl.Edit)
		}))

	// 测试相关路由
	testsCtl := ctl.TestsCtl{}
	r.HandleFunc("/tests/test", testsCtl.Test)
}
