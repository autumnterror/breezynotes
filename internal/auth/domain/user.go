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

func UserToRpc(u *User) *brzrpc.User {
	if u == nil {
		return nil
	}
	return &brzrpc.User{
		Id:       u.Id,
		Login:    u.Login,
		Email:    u.Email,
		About:    u.About,
		Photo:    u.Photo,
		Password: u.Password,
	}
}
