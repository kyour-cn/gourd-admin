package service

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gen"
	"gourd/internal/config"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories"
	"strconv"
	"strings"
	"time"
)

func LoginUser(username string, password string) (*model.User, error) {

	// TODO: 登录频率限制锁

	ru := repositories.User{}
	uq := query.User

	// 查询用户
	userModel, err := ru.Query().
		Where(
			uq.Username.Eq(username),
			uq.Password.Eq(password),
		).
		Select(
			uq.ID,
			uq.Realname,
			uq.Username,
			uq.Mobile,
			uq.Avatar,
			uq.RegisterTime,
			uq.Status,
			uq.RoleID,
		).
		First()
	if err != nil {
		return nil, err
	}

	return userModel, nil
}

// GenerateToken 生成token
func GenerateToken(user *model.User) (string, error) {

	conf, err := config.GetJwtConfig()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"iss": "gourd_admin",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Second * time.Duration(conf.Expire)).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(conf.Secret))
}

type menuTree struct {
	ParentId  int32      `json:"parentId"`
	Id        int32      `json:"id"`
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	RuleId    int32      `json:"ruleId"`
	Component string     `json:"component"`
	Sort      int32      `json:"sort"`
	Meta      *menuMeta  `json:"meta"`
	AppId     int32      `json:"appId"`
	Children  []menuTree `json:"children"`
}
type menuTreeArr []menuTree
type menuMeta struct {
	Title string `json:"title"`
	Type  string `json:"type"`
	Icon  string `json:"icon"`
}

func GetMenu(userInfo *model.User) (any, error) {

	rr := repositories.Role{}
	rm := repositories.Menu{}

	// 查询角色
	roleInfo, err := rr.Query().
		Where(
			query.Role.ID.Eq(userInfo.RoleID),
		).
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
		conditions = append(conditions, qm.RuleID.In(ids...))
	}

	menu, err := rm.Query().
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

			mate := &menuMeta{}
			if menu.Meta != "" {
				_ = json.Unmarshal([]byte(menu.Meta), mate)
			}

			menuTree := menuTree{
				ParentId:  menu.Pid,
				Id:        menu.ID,
				Name:      menu.Name,
				Path:      menu.Path,
				RuleId:    menu.RuleID,
				Component: menu.Component,
				Sort:      menu.Sort,
				Meta:      mate,
				AppId:     menu.AppID,
				Children:  children,
			}
			menuTreeArr = append(menuTreeArr, menuTree)
		}
	}
	return menuTreeArr
}
