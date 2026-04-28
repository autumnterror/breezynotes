package domain

type Id struct {
	Id string `json:"id"`
}
type Message struct {
	Message string `json:"message" example:"some info"`
}

type Error struct {
	Error string `json:"error" example:"error"`
}

type Name struct {
	Name string `json:"name"`
}
