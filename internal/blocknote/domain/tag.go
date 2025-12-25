package domain

import brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"

type Tag struct {
	Id     string `bson:"_id"`
	Title  string `bson:"title"`
	Color  string `bson:"color"`
	Emoji  string `bson:"emoji"`
	UserId string `bson:"user_id"`
}

func ToTagDb(t *brzrpc.Tag) *Tag {
	if t == nil {
		return nil
	}
	return &Tag{
		Id:     t.Id,
		Title:  t.Title,
		Color:  t.Color,
		Emoji:  t.Emoji,
		UserId: t.UserId,
	}
}

func FromTagDb(t *Tag) *brzrpc.Tag {
	if t == nil {
		return nil
	}
	return &brzrpc.Tag{
		Id:     t.Id,
		Title:  t.Title,
		Color:  t.Color,
		Emoji:  t.Emoji,
		UserId: t.UserId,
	}
}
