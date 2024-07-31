package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"gourd/internal/config"
	"net/http"
)

func AuthJwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		header := authHeader[len("Bearer "):]

		conf, err := config.GetJwtConfig()
		if err != nil {
			http.Error(w, "Internal Server Error: No JWT config found", http.StatusInternalServerError)
			return
		}

		jwtData := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(header, jwtData, func(token *jwt.Token) (interface{}, error) {
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

		// TODO: 验证 JWT token 权限

		ctx := context.WithValue(r.Context(), "jwt", jwtData)

		// 继续处理实际请求
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
