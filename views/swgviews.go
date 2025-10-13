package views

type SWGMessage struct {
	Message string `json:"message" example:"some info"`
}

type SWGError struct {
	Error string `json:"error" example:"error"`
}
