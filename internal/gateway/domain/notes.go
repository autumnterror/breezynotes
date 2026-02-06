package domain

import brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"

type ChangeTitleNoteRequest struct {
	Id    string `json:"note_id"`
	Title string `json:"title"`
}
type CreateNoteRequest struct {
	Title string `json:"title"`
}
type NoteListPaginationResponse struct {
	Items []*NotePart `json:"items"`
	Total int         `json:"total"`
}

type NoteTagId struct {
	TagId  string `json:"tag_id"`
	NoteId string `json:"note_id"`
}

type NoteId struct {
	NoteId string `json:"note_id"`
}

type ShareNoteRequest struct {
	NoteId string `json:"note_id"`
	Login  string `json:"login"`
	Role   string `json:"role"`
}
type ChangeRoleRequest struct {
	NoteId  string `json:"note_id"`
	Login   string `json:"login"`
	NewRole string `json:"new_role"`
}

type NotePart struct {
	Id         string
	Title      string
	Tag        *Tag
	FirstBlock string
	UpdatedAt  int64
	Role       string
}

func ToNotePart(n *brzrpc.NotePart) *NotePart {
	return &NotePart{
		Id:         n.GetId(),
		Title:      n.GetTitle(),
		Tag:        ToTag(n.Tag),
		FirstBlock: n.GetFirstBlock(),
		UpdatedAt:  n.GetUpdatedAt(),
		Role:       n.GetRole(),
	}
}
func ToNotePartList(n []*brzrpc.NotePart) []*NotePart {
	if n == nil {
		return []*NotePart{}
	}

	var npl []*NotePart

	for _, np := range n {
		npl = append(npl, ToNotePart(np))
	}
	return npl
}

type NoteWithBlocks struct {
	Title     string   `json:"title"`
	CreatedAt int64    `json:"created_at"`
	UpdatedAt int64    `json:"updated_at"`
	Tag       *Tag     `json:"tag"`
	Id        string   `json:"id"`
	Author    string   `json:"author"`
	Editors   []string `json:"editors"`
	Readers   []string `json:"readers"`
	Blocks    []Block  `json:"blocks"`
}

func ToNoteWithBlocksDb(n *brzrpc.NoteWithBlocks) *NoteWithBlocks {
	if n == nil {
		return nil
	}
	nn := brzrpc.NoteWithBlocks{Blocks: n.Blocks, Editors: n.Editors, Readers: n.Readers}
	if n.Blocks == nil {
		nn.Blocks = []*brzrpc.Block{}
	}
	if n.Editors == nil {
		nn.Editors = []string{}
	}
	if n.Readers == nil {
		nn.Readers = []string{}
	}
	return &NoteWithBlocks{
		Id:        n.Id,
		Title:     n.Title,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Tag:       ToTag(n.Tag),
		Author:    n.Author,
		Editors:   nn.Editors,
		Readers:   nn.Readers,
		//Blocks:    ToBlocksDb(&brzrpc.Blocks{Items: nn.Blocks}).Blks,
	}
}
