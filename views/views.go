package views

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
	Title   string   `json:"title,omitempty"`
	TagId   string   `json:"tag,omitempty"`
	Editors []string `json:"editors,omitempty"`
	Readers []string `json:"readers,omitempty"`
	Blocks  []string `json:"blocks,omitempty"`
	Status  int      `json:"status"`
}
type TagReq struct {
	Title string `json:"title"`
	Color string `json:"color"`
	Emoji string `json:"emoji"`
}
