package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"app/internal/config"
	auth2 "app/internal/modules/common/auth"
)

// AuthJwtMiddleware JWT鉴权中间件
// 验证JWT token并获取用户信息
func AuthJwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从请求头中获取 JWT token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		headerToken := authHeader[len("Bearer "):]

		// 获取JWT配置
		conf, err := config.GetJwtConfig()
		if err != nil {
			http.Error(w, "Internal Server Error: No JWT config found", http.StatusInternalServerError)
			return
		}

		// 解析 JWT token到Claims
		claims := auth2.UserClaims{}
		token, err := jwt.ParseWithClaims(headerToken, &claims, func(token *jwt.Token) (any, error) {
			return []byte(conf.Secret), nil
		})
		if err != nil {
			http.Error(w, "Invalid JWT token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// 验证 JWT token 有效性
		if !token.Valid {
			http.Error(w, "Invalid JWT token: Token is expired or invalid", http.StatusUnauthorized)
			return
		}

		// 验证 JWT token 的接口权限
		if !auth2.CheckPath(claims, r) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// 将 JWT token 信息存入 context 中
		ctx := context.WithValue(r.Context(), "jwt", claims)

		// 继续处理实际请求
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
