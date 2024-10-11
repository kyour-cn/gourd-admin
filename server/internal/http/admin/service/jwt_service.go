package service

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	jwt.RegisteredClaims
	Uid    int32  `json:"uid"`
	Uname  string `json:"uname"`
	RoleId int32  `json:"role_id"`
	AppId  int32  `json:"app_id"`
}
