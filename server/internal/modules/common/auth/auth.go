package auth

import (
	"app/internal/orm/query"
	"net/http"
	"strconv"
	"strings"
)

// CheckPath 检查Token接口权限
func CheckPath(claims UserClaims, r *http.Request) bool {
	url := r.URL.Path

	apis, err := query.MenuAPI.
		Where(query.MenuAPI.Path.Eq(url)).
		Select(query.MenuAPI.ID).
		Find()
	if err == nil && len(apis) == 0 {
		// 路由未定义，不限制
		return true
	} else if err != nil {
		return false
	}

	uq := query.User

	userInfo, err := uq.WithContext(r.Context()).
		Preload(uq.UserRole, uq.UserRole.Role).
		Where(
			uq.ID.Eq(claims.Sub),
			uq.Status.Eq(1),
		).
		First()
	if err != nil {
		// 用户状态已失效
		return false
	}

	// 权限匹配
	ruleSet := make(map[int32]bool)
	for _, v := range userInfo.UserRole {
		// 管理员角色拥有所有权限
		if v.Role.IsAdmin == 1 {
			return true
		}
		for _, ruleIDStr := range strings.Split(v.Role.Rules, ",") {
			ruleID, _ := strconv.Atoi(ruleIDStr)
			ruleSet[int32(ruleID)] = true
		}
	}

	// 判断是否有交集
	for _, api := range apis {
		if ruleSet[api.ID] {
			return true
		}
	}

	return false
}
