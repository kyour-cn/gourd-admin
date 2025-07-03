package auth

import (
	"app/internal/orm/model"
	"app/internal/orm/query"
	"encoding/json"
	"errors"
	"gorm.io/gen"
	"strconv"
	"strings"
)

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
type MenuTreeArr []menuTree

// GetMenu 获取用户菜单
func GetMenu(userInfo *model.User, appId int32) (MenuTreeArr, error) {
	if len(userInfo.UserRole) == 0 {
		return nil, errors.New("用户角色不存在")
	}

	isAdmin := false

	// 筛选指定appid的菜单列表
	mIds := make([]int32, 0)
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
			mIds = append(mIds, int32(i))
		}
	}

	qm := query.Menu
	conditions := []gen.Condition{
		qm.Status.Eq(1),
		qm.AppID.Eq(appId),
	}

	// 判断是否管理员
	if !isAdmin {
		if len(mIds) == 0 {
			return nil, errors.New("暂无权限")
		}
		conditions = append(conditions, qm.ID.In(mIds...))
	}

	// 获取菜单
	menu, err := query.Menu.
		Preload(query.Menu.MenuApi).
		Where(conditions...).
		Find()
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	return recursionMenu(menu, 0), nil
}

func GetPermission(userInfo *model.User, appId int32) ([]string, error) {
	if len(userInfo.UserRole) == 0 {
		return nil, errors.New("用户角色不存在")
	}

	isAdmin := false

	// 筛选指定appid的菜单列表
	mIds := make([]int32, 0)
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
			mIds = append(mIds, int32(i))
		}
	}

	qma := query.MenuAPI
	conditions := []gen.Condition{
		qma.AppID.Eq(appId),
	}

	// 判断是否管理员
	if !isAdmin {
		if len(mIds) == 0 {
			return []string{}, nil
		}
		conditions = append(conditions, qma.MenuID.In(mIds...))
	}

	// 查询菜单的全部接口权限
	list, err := query.MenuAPI.
		Where(conditions...).
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

// GetMenuFormApp 获取指定应用的菜单
func GetMenuFormApp(appId int32) (any, error) {

	qm := query.Menu
	conditions := []gen.Condition{
		qm.Status.Eq(1),
		qm.AppID.Eq(appId),
	}

	menu, err := query.Menu.
		Preload(query.Menu.MenuApi).
		Where(conditions...).
		Find()
	if err != nil {
		return nil, err
	}
	// 构建菜单树
	return recursionMenu(menu, 0), nil
}

// 递归菜单
func recursionMenu(menus []*model.Menu, parentId int32) MenuTreeArr {
	var arr MenuTreeArr
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
				ApiList:   menu.MenuApi,
			}
			arr = append(arr, menuTree)
		}
	}
	return arr
}
