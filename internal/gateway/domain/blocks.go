package domain

type Block struct {
	Id        string         `json:"id"`
	Type      string         `json:"type"`
	NoteId    string         `json:"note_id"`
	Order     int32          `json:"order"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	IsUsed    bool           `json:"is_used"`
	Data      map[string]any `json:"data"`
}

type Blocks struct {
	Items []*Block `json:"items,omitempty"`
}
type CreateBlockRequest struct {
	NoteId string         `json:"note_id"`
	Pos    int            `json:"pos"`
	Type   string         `json:"type"`
	Data   map[string]any `json:"data"`
}
type OpBlockRequest struct {
	BlockId string         `json:"block_id"`
	Op      string         `json:"op"`
	Data    map[string]any `json:"data"`
	NoteId  string         `json:"note_id"`
}
type ChangeTypeBlockRequest struct {
	BlockId string `json:"block_id"`
	NewType string `json:"new_type"`
	NoteId  string `json:"note_id"`
}
type ChangeBlockOrderRequest struct {
	OldOrder int    `json:"old_order"`
	NewOrder int    `json:"new_order"`
	NoteId   string `json:"note_id"`
}

type BlockNoteId struct {
	BlockId string `json:"block_id"`
	NoteId  string `json:"note_id"`
}
