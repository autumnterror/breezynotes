package domain

import brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"

type Note struct {
	Id        string   `bson:"_id"`
	Title     string   `bson:"title"`
	CreatedAt int64    `bson:"created_at"`
	UpdatedAt int64    `bson:"updated_at"`
	Tag       *Tag     `bson:"tag"`
	Author    string   `bson:"author"`
	Editors   []string `bson:"editors"`
	Readers   []string `bson:"readers"`
	Blocks    []string `bson:"blocks"`
}

type Notes struct {
	Nts []*Note
}

func ToNoteDb(n *brzrpc.Note) *Note {
	if n == nil {
		return nil
	}
	nn := brzrpc.Note{Blocks: n.Blocks, Editors: n.Editors, Readers: n.Readers}
	if n.Blocks == nil {
		nn.Blocks = []string{}
	}
	if n.Editors == nil {
		nn.Editors = []string{}
	}
	if n.Readers == nil {
		nn.Readers = []string{}
	}
	return &Note{
		Id:        n.Id,
		Title:     n.Title,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Tag:       ToTagDb(n.Tag),
		Author:    n.Author,
		Editors:   nn.Editors,
		Readers:   nn.Readers,
		Blocks:    nn.Blocks,
	}
}

func FromNoteDb(n *Note) *brzrpc.Note {
	if n == nil {
		return nil
	}
	nn := brzrpc.Note{Blocks: n.Blocks, Editors: n.Editors, Readers: n.Readers}
	if n.Blocks == nil {
		nn.Blocks = []string{}
	}
	if n.Editors == nil {
		nn.Editors = []string{}
	}
	if n.Readers == nil {
		nn.Readers = []string{}
	}
	return &brzrpc.Note{
		Id:        n.Id,
		Title:     n.Title,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		Tag:       FromTagDb(n.Tag),
		Author:    n.Author,
		Editors:   nn.Editors,
		Readers:   nn.Readers,
		Blocks:    nn.Blocks,
	}
}

func ToNotesDb(n *brzrpc.Notes) *Notes {
	if n == nil {
		return nil
	}

	var nts []*Note
	for _, nt := range n.GetItems() {
		nts = append(nts, ToNoteDb(nt))
	}

	return &Notes{
		Nts: nts,
	}
}

func FromNotesDb(n *Notes) *brzrpc.Notes {
	if n == nil {
		return nil
	}

	var nts []*brzrpc.Note
	for _, nt := range n.Nts {
		nts = append(nts, FromNoteDb(nt))
	}

	return &brzrpc.Notes{
		Items: nts,
	}
}

type NotePart struct {
	Id         string
	Title      string
	Tag        *Tag
	FirstBlock string
	UpdatedAt  int64
}
type NoteParts struct {
	Ntps []*NotePart
}

func FromNotePartDb(n *NotePart) *brzrpc.NotePart {
	return &brzrpc.NotePart{
		Id:         n.Id,
		Title:      n.Title,
		Tag:        FromTagDb(n.Tag),
		FirstBlock: n.FirstBlock,
		UpdatedAt:  n.UpdatedAt,
	}
}

func ToNotePartDb(n *brzrpc.NotePart) *NotePart {
	return &NotePart{
		Id:         n.GetId(),
		Title:      n.GetTitle(),
		Tag:        ToTagDb(n.Tag),
		FirstBlock: n.GetFirstBlock(),
		UpdatedAt:  n.GetUpdatedAt(),
	}
}

func ToNotePartsDb(n *brzrpc.NoteParts) *NoteParts {
	if n == nil {
		return nil
	}

	var nts []*NotePart
	for _, nt := range n.GetItems() {
		nts = append(nts, ToNotePartDb(nt))
	}

	return &NoteParts{
		Ntps: nts,
	}
}

func FromNotePartsDb(n *NoteParts) *brzrpc.NoteParts {
	if n == nil {
		return nil
	}

	var nts []*brzrpc.NotePart
	for _, nt := range n.Ntps {
		nts = append(nts, FromNotePartDb(nt))
	}

	return &brzrpc.NoteParts{
		Items: nts,
	}
}
