package domain

import brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"

type Tag struct {
	Id       string `bson:"_id"`
	Title    string `bson:"title"`
	Color    string `bson:"color"`
	Emoji    string `bson:"emoji"`
	UserId   string `bson:"user_id"`
	IsPinned bool   `bson:"is_pinned"`
}

type Tags struct {
	Tgs []*Tag
}

func ToTagDb(t *brzrpc.Tag) *Tag {
	if t == nil {
		return nil
	}
	return &Tag{
		Id:       t.Id,
		Title:    t.Title,
		Color:    t.Color,
		Emoji:    t.Emoji,
		UserId:   t.UserId,
		IsPinned: t.IsPinned,
	}
}

func FromTagDb(t *Tag) *brzrpc.Tag {
	if t == nil {
		return nil
	}
	return &brzrpc.Tag{
		Id:       t.Id,
		Title:    t.Title,
		Color:    t.Color,
		Emoji:    t.Emoji,
		UserId:   t.UserId,
		IsPinned: t.IsPinned,
	}
}

func ToTagsDb(t *brzrpc.Tags) *Tags {
	if t == nil {
		return nil
	}

	var tgs []*Tag

	for _, tg := range t.GetItems() {
		tgs = append(tgs, ToTagDb(tg))
	}

	return &Tags{
		Tgs: tgs,
	}
}

func FromTagsDb(t *Tags) *brzrpc.Tags {
	if t == nil {
		return nil
	}
	var tgs []*brzrpc.Tag

	for _, tg := range t.Tgs {
		tgs = append(tgs, FromTagDb(tg))
	}

	return &brzrpc.Tags{
		Items: tgs,
	}
}
