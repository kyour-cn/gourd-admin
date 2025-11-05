package dto

import "app/internal/modules/common/auth"

type UserResetPasswordReq struct {
	UserPassword       string `json:"user_password" validate:"required,min=6,max=32"`
	NewPassword        string `json:"new_password" validate:"required,min=6,max=32"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required,min=6,max=32"`
	Claims             *auth.UserClaims
}

type UserUpdateNameReq struct {
	Nickname string `json:"nickname" validate:"required,min=2,max=20"`
	Claims   *auth.UserClaims
}
