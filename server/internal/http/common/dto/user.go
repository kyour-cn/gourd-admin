package dto

import (
	"github.com/golang-jwt/jwt/v5"

	"app/internal/orm/model"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Sub  int32  `json:"sub"`
	Name string `json:"name"`
}

type UserLoginRes struct {
	Token    string      `json:"token"`
	UserInfo *model.User `json:"userInfo"`
	Expire   int64       `json:"expire"`
	Apps     any         `json:"apps"`
}

type UserResetPasswordReq struct {
	UserPassword       string `json:"user_password" validate:"required|minLen:6|maxLen:32"`
	NewPassword        string `json:"new_password" validate:"required|minLen:6|maxLen:32"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required|minLen:6|maxLen:32"`
	Claims             *UserClaims
}

type UserUpdateNameReq struct {
	Nickname string `json:"nickname" validate:"required|minLen:2|maxLen:20"`
	Claims   *UserClaims
}
