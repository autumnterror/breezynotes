package views

type SWGMessage struct {
	Message string `json:"message" example:"some info"`
}

type SWGError struct {
	Error string `json:"error" example:"error"`
}

type SWGBlock struct {
	Id        string         `json:"id,omitempty"`
	Type      string         `json:"type,omitempty"`
	NoteId    string         `json:"note_id,omitempty"`
	Order     int32          `json:"order,omitempty"`
	CreatedAt int64          `json:"created_at,omitempty"`
	UpdatedAt int64          `json:"updated_at,omitempty"`
	IsUsed    bool           `json:"is_used,omitempty"`
	Data      map[string]any `json:"data,omitempty"`
}
type SWGBlocks struct {
	Items []*SWGBlock `json:"items,omitempty"`
}
type SWGCreateBlockRequest struct {
	Type string         `json:"type,omitempty"`
	Data map[string]any `json:"data,omitempty"`
}
type SWGOpBlockRequest struct {
	Id   string         `json:"id,omitempty"`
	Op   string         `json:"op,omitempty"`
	Data map[string]any `json:"data,omitempty"`
}
