package router

import (
	"app/internal/http/admin/controller/system"
	"app/internal/http/middleware"
	"github.com/go-chi/chi/v5"
)

// Router 注册路由组
func Router(r chi.Router) {

	// 跨域中间件
	r.Use(middleware.CorsMiddleware)

	// 登录相关路由
	r.Route("/auth", func(r chi.Router) {
		c := system.Auth{}
		r.HandleFunc("/captcha", c.Captcha)
		r.HandleFunc("/login", c.Login)
		r.With(middleware.AuthJwtMiddleware).
			HandleFunc("/menu", c.GetMenu)
	})

	// 系统相关路由
	r.Route("/system", func(r chi.Router) {

		// app相关路由
		r.Route("/app", func(r chi.Router) {
			c := system.App{}
			r.Use(middleware.AuthJwtMiddleware)
			r.Get("/list", c.List)
			r.Post("/add", c.Add)
			r.Post("/edit", c.Edit)
			r.Post("/delete", c.Delete)
		})

		// menu相关路由
		r.Route("/menu", func(r chi.Router) {
			c := system.Menu{}
			r.Use(middleware.AuthJwtMiddleware)
			r.Get("/list", c.List)
			r.Post("/add", c.Add)
			r.Post("/edit", c.Edit)
			r.Post("/delete", c.Delete)
		})

		// role相关路由
		r.Route("/role", func(r chi.Router) {
			c := system.Role{}
			r.Use(middleware.AuthJwtMiddleware)
			r.Get("/list", c.List)
			r.Post("/add", c.Add)
			r.Post("/edit", c.Edit)
			r.Post("/delete", c.Delete)
		})

		// user相关路由
		r.Route("/user", func(r chi.Router) {
			c := system.User{}
			r.Use(middleware.AuthJwtMiddleware)
			r.HandleFunc("/list", c.List)
			r.HandleFunc("/add", c.Add)
			r.HandleFunc("/edit", c.Edit)
			r.HandleFunc("/delete", c.Delete)
		})

		// log相关路由
		r.Route("/log", func(r chi.Router) {
			c := system.Log{}
			r.Use(middleware.AuthJwtMiddleware)
			r.HandleFunc("/typeList", c.TypeList)
			r.HandleFunc("/list", c.List)
			r.HandleFunc("/logStat", c.LogStat)
		})
	})

}
