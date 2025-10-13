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
