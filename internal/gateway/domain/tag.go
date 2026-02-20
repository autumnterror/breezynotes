package domain

import brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"

type Tag struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Color    string `json:"color"`
	Emoji    string `json:"emoji"`
	UserId   string `json:"user_id"`
	IsPinned bool   `json:"is_pinned"`
}
type Tags struct {
	Tgs []Tag `json:"tags"`
}

func ToTag(t *brzrpc.Tag) *Tag {
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

func ToTags(t *brzrpc.Tags) *Tags {
	if t == nil {
		return nil
	}

	var tgs []Tag

	for _, tg := range t.GetItems() {
		tgs = append(tgs, *ToTag(tg))
	}

	return &Tags{
		Tgs: tgs,
	}
}

type CreateTagRequest struct {
	Title string `json:"title"`
	Color string `json:"color"`
	Emoji string `json:"emoji"`
}

type UpdateTagTitleRequest struct {
	IdTag string `json:"id_tag"`
	Title string `json:"title"`
}
type UpdateTagEmojiRequest struct {
	IdTag string `json:"id_tag"`
	Emoji string `json:"emoji"`
}
type UpdateTagColorRequest struct {
	IdTag string `json:"id_tag"`
	Color string `json:"color"`
}
type UpdatePinnedEmojiRequest struct {
	IdTag string `json:"id_tag"`
}
