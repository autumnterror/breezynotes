package domain

import brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"

type User struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	About    string `json:"about"`
	Photo    string `json:"photo"`
	Password string `json:"password"`
}

func UserFromRpc(u *brzrpc.User) *User {
	if u == nil {
		return nil
	}
	return &User{
		Id:       u.GetId(),
		Login:    u.GetLogin(),
		Email:    u.GetEmail(),
		About:    u.GetAbout(),
		Photo:    u.GetPhoto(),
		Password: u.GetPassword(),
	}
}

type UpdateAboutRequest struct {
	NewAbout string `json:"new_about"`
}
type UpdateEmailRequest struct {
	NewEmail string `json:"new_email"`
}
type UpdatePhotoRequest struct {
	NewPhoto string `json:"new_photo"`
}
type ChangePasswordRequest struct {
	OldPassword  string `json:"old_password"`
	NewPassword  string `json:"new_password"`
	NewPassword2 string `json:"new_password_2"`
}
type UpdatePasswordRequest struct {
	OldPassword  string `json:"old_password"`
	NewPassword  string `json:"new_password"`
	NewPassword2 string `json:"new_password_2"`
}
type AuthRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserRegister struct {
	Login string `json:"login"`
	Email string `json:"email"`
	Pw1   string `json:"pw1"`
	Pw2   string `json:"pw2"`
}
