package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gen"

	"app/internal/config"
	"app/internal/http/common/dto"
	"app/internal/orm/model"
	"app/internal/orm/query"
	"app/internal/util/cache"
)

func NewAuthService(ctx context.Context) *AuthService {
	return &AuthService{
		ctx: ctx,
	}
}

type AuthService struct {
	ctx context.Context
}

type MenuMate struct {
	Title            string `json:"title"`
	Icon             string `json:"icon"`
	Active           string `json:"active"`
	Color            string `json:"color"`
	Type             string `json:"type"`
	FullPage         bool   `json:"fullPage"`
	Tag              string `json:"tag"`
	Hidden           bool   `json:"hidden"`
	HiddenBreadcrumb bool   `json:"hiddenBreadcrumb"`
}
type menuTree struct {
	ParentId  int64           `json:"pid"`
	Id        int64           `json:"id"`
	Name      string          `json:"name"`
	Title     string          `json:"title"`
	Path      string          `json:"path"`
	Component string          `json:"component"`
	Sort      int64           `json:"sort"`
	Meta      *MenuMate       `json:"meta"`
	AppId     int64           `json:"appId"`
	Children  []menuTree      `json:"children"`
	ApiList   []model.MenuAPI `json:"apiList"`
}
type MenuTreeArr []menuTree

func (s *AuthService) GetMenu(userInfo *model.User, appId int64) (MenuTreeArr, error) {
	if len(userInfo.UserRole) == 0 {
		return nil, errors.New("用户角色不存在")
	}

	isAdmin := false

	// 筛选指定appid的菜单列表
	mIds := make([]int64, 0)
	for _, v := range userInfo.UserRole {
		if v.Role.AppID != appId {
			continue
		}
		if v.Role.IsAdmin == 1 {
			isAdmin = true
			break
		}
		for _, v := range strings.Split(v.Role.Rules, ",") {
			i, _ := strconv.Atoi(v)
			mIds = append(mIds, int64(i))
		}
	}

	qm := query.Menu
	conds := []gen.Condition{
		qm.AppID.Eq(appId),
		qm.Type.Eq("menu"),
	}

	// 判断是否管理员
	if !isAdmin {
		if len(mIds) == 0 {
			return nil, errors.New("暂无权限")
		}
		conds = append(conds, qm.ID.In(mIds...))
	}

	// 获取菜单
	menu, err := qm.WithContext(s.ctx).
		Preload(qm.MenuApi).
		Where(conds...).
		Find()
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	return s.RecursionMenu(menu, 0), nil
}

// RecursionMenu 递归菜单
func (s *AuthService) RecursionMenu(menus []*model.Menu, parentId int64) MenuTreeArr {
	var arr MenuTreeArr
	for _, menu := range menus {
		if menu.Pid == parentId {
			children := s.RecursionMenu(menus, menu.ID)

			mate := &MenuMate{}
			if menu.Meta != "" {
				_ = json.Unmarshal([]byte(menu.Meta), mate)
			}

			menuTree := menuTree{
				ParentId:  menu.Pid,
				Id:        menu.ID,
				Name:      menu.Name,
				Title:     menu.Title,
				Path:      menu.Path,
				Component: menu.Component,
				Sort:      menu.Sort,
				Meta:      mate,
				AppId:     menu.AppID,
				Children:  children,
				ApiList:   menu.MenuApi,
			}
			arr = append(arr, menuTree)
		}
	}

	// 按Sort字段从小到大排序
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].Sort < arr[j].Sort
	})

	return arr
}

// GetMenuFormApp 获取指定应用的菜单
func (s *AuthService) GetMenuFormApp(appId int64) (any, error) {

	qm := query.Menu
	conds := []gen.Condition{
		qm.AppID.Eq(appId),
	}

	menu, err := query.Menu.
		Preload(query.Menu.MenuApi).
		Where(conds...).
		Find()
	if err != nil {
		return nil, err
	}
	// 构建菜单树
	return s.RecursionMenu(menu, 0), nil
}

func (s *AuthService) GetPermission(userInfo *model.User, appId int64) ([]string, error) {
	if len(userInfo.UserRole) == 0 {
		return nil, errors.New("用户角色不存在")
	}

	isAdmin := false

	// 筛选指定appid的菜单列表
	mIds := make([]int64, 0)
	for _, v := range userInfo.UserRole {
		if v.Role.AppID != appId {
			continue
		}
		if v.Role.IsAdmin == 1 {
			isAdmin = true
			break
		}
		for _, v := range strings.Split(v.Role.Rules, ",") {
			i, _ := strconv.Atoi(v)
			mIds = append(mIds, int64(i))
		}
	}

	qma := query.MenuAPI
	conds := []gen.Condition{
		qma.AppID.Eq(appId),
	}

	// 判断是否管理员
	if !isAdmin {
		if len(mIds) == 0 {
			return []string{}, nil
		}
		conds = append(conds, qma.MenuID.In(mIds...))
	}

	// 查询菜单的全部接口权限
	list, err := query.MenuAPI.
		Where(conds...).
		Select(query.MenuAPI.ID, query.MenuAPI.Tag).
		Find()
	if err != nil {
		return nil, errors.New("获取权限列表失败")
	}

	var arr []string
	for _, v := range list {
		arr = append(arr, v.Tag)
	}

	return arr, nil
}

// LoginUser 登录用户
func (s *AuthService) LoginUser(username string, password string) (*model.User, error) {
	// 登录频率限制锁 10秒
	key := "login_lock:" + username
	val, ok := cache.Get(key)
	if !ok {
		val = 0
	}
	failures, _ := val.(int)
	// 10秒登录失败次数超过3次，禁止登录
	if failures > 3 {
		return nil, errors.New("登录失败次数过多，请稍后再试")
	}

	uq := query.User

	// 查询用户
	userModel, err := uq.
		Preload(uq.UserRole, uq.UserRole.Role, uq.UserRole.Role.App).
		Where(
			uq.Username.Eq(username),
			uq.Password.Eq(password),
		).
		Select(uq.ID, uq.Nickname, uq.Username, uq.Avatar, uq.CreatedAt, uq.Status).
		First()
	if err != nil {
		// 登录失败次数+1
		cache.Set(key, failures+1, 10*time.Second)
		return nil, errors.New("用户名或密码错误")
	}
	cache.Delete(key)

	return userModel, nil
}

// GenerateToken 生成token
func (s *AuthService) GenerateToken(claims dto.UserClaims) (string, error) {
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

// CheckPath 检查Token接口权限
func (s *AuthService) CheckPath(claims dto.UserClaims, r *http.Request) bool {
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
	ruleSet := make(map[int64]bool)
	for _, v := range userInfo.UserRole {
		// 管理员角色拥有所有权限
		if v.Role.IsAdmin == 1 {
			return true
		}
		for _, ruleIDStr := range strings.Split(v.Role.Rules, ",") {
			ruleID, _ := strconv.Atoi(ruleIDStr)
			ruleSet[int64(ruleID)] = true
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
