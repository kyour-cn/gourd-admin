package router

import (
	"app/internal/http/admin/controller/common"
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
		c := common.Auth{}
		r.HandleFunc("/captcha", c.Captcha)
		r.HandleFunc("/login", c.Login)
		r.With(middleware.AuthJwtMiddleware).
			HandleFunc("/menu", c.GetMenu)
	})

	// 上传相关路由
	r.Route("/upload", func(r chi.Router) {
		c := common.Upload{}
		r.Use(middleware.AuthJwtMiddleware)
		r.Post("/image", c.Image)   // 上传图片
		r.Post("/file", c.File)     // 上传文件
		r.Post("/delete", c.Delete) // 上传文件
	})

	// 用户
	r.Route("/user", func(r chi.Router) {
		c := common.User{}
		r.Use(middleware.AuthJwtMiddleware)
		r.HandleFunc("/info", c.Info)        // 用户信息
		r.Post("/password", c.ResetPassword) // 修改密码
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
			r.Get("/list", c.List)
			r.Post("/add", c.Add)
			r.Post("/edit", c.Edit)
			r.Post("/delete", c.Delete)
		})

		// log相关路由
		r.Route("/log", func(r chi.Router) {
			c := system.Log{}
			r.Use(middleware.AuthJwtMiddleware)
			r.Get("/list", c.List)
			r.Get("/typeList", c.TypeList)
			r.Get("/logStat", c.LogStat)
		})
	})

}
