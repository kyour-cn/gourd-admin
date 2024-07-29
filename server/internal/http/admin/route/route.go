package route

import (
	"github.com/go-chi/chi/v5"
	"gourd/internal/http/admin/ctl"
	"gourd/internal/http/admin/middleware"
)

// RegisterRoute 注册路由组
func RegisterRoute(r chi.Router) {

	r.Use(middleware.CorsMiddleware)

	authCtl := ctl.AuthCtl{}
	r.HandleFunc("/auth/captcha", authCtl.Captcha)
	r.HandleFunc("/auth/login", authCtl.Login)

	testsCtl := ctl.TestsCtl{}
	r.HandleFunc("/tests/test", testsCtl.Test)

}
