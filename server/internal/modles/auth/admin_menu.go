package auth

import (
	"app/internal/orm/model"
	"app/internal/orm/query"
	"encoding/json"
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
		Preload(query.Menu.MenuApi).
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
		Preload(query.Menu.MenuApi).
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
				ApiList:   menu.MenuApi,
			}
			menuTreeArr = append(menuTreeArr, menuTree)
		}
	}
	return menuTreeArr
}
