package views

import brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"

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
type SWGNoteWithBlocks struct {
	Title     string      `json:"title,omitempty"`
	CreatedAt int64       `json:"created_at,omitempty"`
	UpdatedAt int64       `json:"updated_at,omitempty"`
	Tag       *brzrpc.Tag `json:"tag,omitempty"`
	Id        string      `json:"id,omitempty"`
	Author    string      `json:"author,omitempty"`
	Editors   []string    `json:"editors,omitempty"`
	Readers   []string    `json:"readers,omitempty"`
	Blocks    []SWGBlock  `json:"blocks,omitempty"`
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
