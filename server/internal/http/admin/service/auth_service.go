package service

import (
	"app/internal/config"
	"app/internal/orm/model"
	"app/internal/orm/query"
	"app/internal/util/redisutil"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gen"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// LoginUser 登录用户
func LoginUser(ctx context.Context, username string, password string) (*model.User, error) {

	rdb, err := redisutil.GetRedis(ctx)
	if err != nil {
		return nil, errors.New("redis连接失败")
	}

	// 登录频率限制锁 10秒
	key := "login_lock:" + username
	val, _ := rdb.Get(ctx, key).Result()
	failures, _ := strconv.Atoi(val)
	// 10秒登录失败次数超过3次，禁止登录
	if failures > 3 {
		return nil, errors.New("登录失败次数过多，请稍后再试")
	}

	uq := query.User

	// 查询用户
	userModel, err := uq.
		Where(
			uq.Username.Eq(username),
			uq.Password.Eq(password),
		).
		Select(
			uq.ID,
			uq.Nickname,
			uq.Username,
			uq.Mobile,
			uq.Avatar,
			uq.CreateTime,
			uq.Status,
			uq.RoleID,
		).
		First()
	if err != nil {

		// 登录失败次数+1
		rdb.Incr(ctx, key)
		rdb.Expire(ctx, key, 10*time.Second)

		return nil, errors.New("用户名或密码错误")
	}

	rdb.Del(ctx, key)

	return userModel, nil
}

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

type MenuMate struct {
	Title            string `json:"title"`
	Icon             string `json:"icon"`
	Active           string `json:"active"`
	Color            string `json:"color"`
	Type             string `json:"type"`
	Fullpage         bool   `json:"fullpage"`
	Tag              string `json:"tag"`
	Hidden           bool   `json:"hidden"`
	HiddenBreadcrumb bool   `json:"hiddenBreadcrumb"`
}
type menuTree struct {
	ParentId  int32           `json:"parentId"`
	Id        int32           `json:"id"`
	Name      string          `json:"name"`
	Path      string          `json:"path"`
	Component string          `json:"component"`
	Sort      int32           `json:"sort"`
	Meta      *MenuMate       `json:"meta"`
	AppId     int32           `json:"appId"`
	Children  []menuTree      `json:"children"`
	ApiList   []model.MenuAPI `json:"apiList"`
}
type menuTreeArr []menuTree

// GetMenu 获取用户菜单
func GetMenu(userInfo *model.User) (any, error) {

	// 查询角色
	roleInfo, err := query.Role.
		Where(query.Role.ID.Eq(userInfo.RoleID)).
		First()
	if err != nil {
		return nil, err
	}

	qm := query.Menu
	conditions := []gen.Condition{
		qm.Status.Eq(1),
		qm.AppID.Eq(roleInfo.AppID),
	}

	// 判断是否管理员
	if roleInfo.IsAdmin == 0 {
		var ids []int32
		arr := strings.Split(roleInfo.Rules, ",")
		for _, v := range arr {
			num, _ := strconv.Atoi(v)
			ids = append(ids, int32(num))
		}
		conditions = append(conditions, qm.ID.In(ids...))
	}

	menu, err := query.Menu.
		Preload(query.Menu.ApiList).
		Where(conditions...).
		Find()
	if err != nil {
		return nil, err
	}
	// 构建菜单树
	menuTreeArr := recursionMenu(menu, 0)
	return menuTreeArr, nil

}

// GetMenuFormApp 获取指定应用的菜单
func GetMenuFormApp(appId int32) (any, error) {

	qm := query.Menu
	conditions := []gen.Condition{
		qm.Status.Eq(1),
		qm.AppID.Eq(appId),
	}

	menu, err := query.Menu.
		Preload(query.Menu.ApiList).
		Where(conditions...).
		Find()
	if err != nil {
		return nil, err
	}
	// 构建菜单树
	menuTreeArr := recursionMenu(menu, 0)
	return menuTreeArr, nil
}

func recursionMenu(menus []*model.Menu, parentId int32) menuTreeArr {
	var menuTreeArr menuTreeArr
	for _, menu := range menus {
		if menu.Pid == parentId {
			children := recursionMenu(menus, menu.ID)

			mate := &MenuMate{}
			if menu.Meta != "" {
				_ = json.Unmarshal([]byte(menu.Meta), mate)
			}

			menuTree := menuTree{
				ParentId:  menu.Pid,
				Id:        menu.ID,
				Name:      menu.Name,
				Path:      menu.Path,
				Component: menu.Component,
				Sort:      menu.Sort,
				Meta:      mate,
				AppId:     menu.AppID,
				Children:  children,
				ApiList:   menu.ApiList,
			}
			menuTreeArr = append(menuTreeArr, menuTree)
		}
	}
	return menuTreeArr
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
