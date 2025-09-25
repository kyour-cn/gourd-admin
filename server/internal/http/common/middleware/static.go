package middleware

import (
	"net/http"
	"os"
	"path/filepath"
)

// StaticOrNext 静态文件处理
func StaticOrNext(dir string) func(http.Handler) http.Handler {
	root := http.Dir(dir)
	fs := http.FileServer(root)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fPath := filepath.Join(dir, r.URL.Path)

			info, err := os.Stat(fPath)
			if err == nil {
				if info.IsDir() {
					// 如果是目录，交给 FileServer (它会处理 index.html)
					fs.ServeHTTP(w, r)
					return
				}
				// 如果是文件，直接返回，避免多余的 301
				http.ServeFile(w, r, fPath)
				return
			}

			// 文件不存在，继续走 chi 路由
			next.ServeHTTP(w, r)
		})
	}
}
