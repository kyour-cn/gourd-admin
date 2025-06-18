package auth

import (
	"app/internal/config"
	"app/internal/orm/query"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GenerateToken 生成token
func GenerateToken(claims UserClaims) (string, error) {
	// 读取配置
	conf, err := config.GetJwtConfig()
	if err != nil {
		return "", err
	}

	// 设置签署时间和过期时间
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(conf.Expire)))

	// 使用HS256算法签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(conf.Secret))
}

// CheckJwtPermission 检查Token接口权限
func CheckJwtPermission(jd UserClaims, r *http.Request) bool {

	// 取出角色ID和应用ID
	if jd.RoleId == 0 || jd.AppId == 0 {
		return false
	}

	url := r.URL.Path
	apis, err := query.MenuAPI.
		Where(
			query.MenuAPI.Path.Eq(url),
			query.MenuAPI.AppID.Eq(jd.AppId),
		).
		Select(query.MenuAPI.ID).
		Find()
	if err != nil {
		// 路由未定义，不限制
		return true
	}

	// 获取用户角色
	role, err := query.Role.
		Where(
			query.Role.ID.Eq(jd.RoleId),
			query.Role.AppID.Eq(jd.AppId),
		).
		Select(
			query.Role.ID,
			query.Role.IsAdmin,
			query.Role.Rules,
		).
		First()
	if err != nil {
		return false
	}
	// 管理员角色拥有所有权限
	if role.IsAdmin == 1 {
		return true
	}

	var ruleIds []int32
	for _, rule := range apis {
		ruleIds = append(ruleIds, rule.ID)
	}

	ruleArr := strings.Split(role.Rules, ",")

	// 判断 ruleIds 和 role.Rules 是否有交集
	for _, rid := range ruleIds {
		for _, rid2 := range ruleArr {
			_id, _ := strconv.Atoi(rid2)
			if rid == int32(_id) {
				// 权限匹配成功
				return true
			}
		}
	}

	return false
}
