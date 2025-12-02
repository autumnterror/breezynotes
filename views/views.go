package views

import (
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
)

type ResRPC struct {
	Res interface{}
	Err error
}

type UserRegister struct {
	Login string `json:"login"`
	Email string `json:"email"`
	Pw1   string `json:"pw1"`
	Pw2   string `json:"pw2"`
}

type NoteReq struct {
	Title string `json:"title,omitempty"`
}

type NoteWithBlocks struct {
	Title     string          `json:"title,omitempty"`
	CreatedAt int64           `json:"created_at,omitempty"`
	UpdatedAt int64           `json:"updated_at,omitempty"`
	Tag       *brzrpc.Tag     `json:"tag,omitempty"`
	Id        string          `json:"id,omitempty"`
	Author    string          `json:"author,omitempty"`
	Editors   []string        `json:"editors,omitempty"`
	Readers   []string        `json:"readers,omitempty"`
	Blocks    []*brzrpc.Block `json:"blocks,omitempty"`
}
type TagReq struct {
	Title string `json:"title"`
	Color string `json:"color"`
	Emoji string `json:"emoji"`
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
	Login        string `json:"login"`
	Email        string `json:"email"`
	OldPassword  string `json:"old_password"`
	NewPassword  string `json:"new_password"`
	NewPassword2 string `json:"new_password_2"`
}
