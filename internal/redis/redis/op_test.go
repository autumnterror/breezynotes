package redis

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/redis/config"
	brzrpc "github.com/autumnterror/breezynotes/pkg/protos/proto/gen"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
	"testing"
)

func TestCrudSession(t *testing.T) {
	c := New(config.Test())
	idTest := "test"
	s, err := structpb.NewStruct(map[string]any{"text": []any{map[string]any{"style": "default", "text": "test"}, map[string]any{"style": "italic", "text": "ghost"}}})
	assert.NoError(t, err)
	assert.NoError(t, c.saveSession(context.TODO(), &userSession{
		Id: idTest,
		NoteParts: []*brzrpc.NotePart{
			{
				Id:    "1",
				Title: "1",
				Tag: &brzrpc.Tag{
					Id:     "1",
					Title:  "1",
					Color:  "1",
					Emoji:  "1",
					UserId: "1",
				},
				FirstBlock: "1",
				UpdatedAt:  1,
			},
			{
				Id:    "1",
				Title: "1",
				Tag: &brzrpc.Tag{
					Id:     "1",
					Title:  "1",
					Color:  "1",
					Emoji:  "1",
					UserId: "1",
				},
				FirstBlock: "1",
				UpdatedAt:  1,
			},
		},
		Notes: []*brzrpc.NoteWithBlocks{
			{
				Id:        "1",
				Title:     "1",
				CreatedAt: 1,
				UpdatedAt: 1,
				Tag: &brzrpc.Tag{
					Id:     "1",
					Title:  "1",
					Color:  "1",
					Emoji:  "1",
					UserId: "1",
				},
				Author:  "1",
				Editors: []string{},
				Readers: []string{},
				Blocks: []*brzrpc.Block{
					{
						Id:        "1",
						Type:      "1",
						NoteId:    "1",
						Order:     1,
						CreatedAt: 1,
						UpdatedAt: 1,
						IsUsed:    true,
						Data:      s,
					},
					{
						Id:        "1",
						Type:      "1",
						NoteId:    "1",
						Order:     1,
						CreatedAt: 1,
						UpdatedAt: 1,
						IsUsed:    true,
						Data:      s,
					},
				},
			},
			{
				Id:        "1",
				Title:     "1",
				CreatedAt: 1,
				UpdatedAt: 1,
				Tag:       nil,
				Author:    "1",
				Editors:   []string{},
				Readers:   []string{},
				Blocks: []*brzrpc.Block{
					{
						Id:        "1",
						Type:      "1",
						NoteId:    "1",
						Order:     1,
						CreatedAt: 1,
						UpdatedAt: 1,
						IsUsed:    true,
						Data:      s,
					},
					{
						Id:        "1",
						Type:      "1",
						NoteId:    "1",
						Order:     1,
						CreatedAt: 1,
						UpdatedAt: 1,
						IsUsed:    true,
						Data:      s,
					}},
			},
		},
		Tags: []*brzrpc.Tag{
			{
				Id:     "1",
				Title:  "1",
				Color:  "1",
				Emoji:  "1",
				UserId: "1",
			},
			{
				Id:     "1",
				Title:  "1",
				Color:  "1",
				Emoji:  "1",
				UserId: "1",
			},
		},
		NoteTrash: []*brzrpc.NotePart{
			{
				Id:    "1",
				Title: "1",
				Tag: &brzrpc.Tag{
					Id:     "1",
					Title:  "1",
					Color:  "1",
					Emoji:  "1",
					UserId: "1",
				},
				FirstBlock: "1",
				UpdatedAt:  1,
			},
			{
				Id:    "1",
				Title: "1",
				Tag: &brzrpc.Tag{
					Id:     "1",
					Title:  "1",
					Color:  "1",
					Emoji:  "1",
					UserId: "1",
				},
				FirstBlock: "1",
				UpdatedAt:  1,
			},
		},
	}))

	if s, err := c.getSession(context.Background(), idTest); assert.NoError(t, err) {
		assert.Equal(t, 2, len(s.Notes))
		assert.Equal(t, 2, len(s.NoteParts))
		assert.Equal(t, 2, len(s.Tags))
		assert.Equal(t, 2, len(s.NoteTrash))
	}

	assert.NoError(t, c.deleteSession(context.TODO(), idTest))
	if _, err := c.getSession(context.Background(), idTest); assert.Error(t, err) {
	}
}
func TestCrudFields(t *testing.T) {
	c := New(config.Test())
	idTest := "TestCrudFields"
	ctx := context.TODO()
	s, err := structpb.NewStruct(map[string]any{"text": []any{map[string]any{"style": "default", "text": "test"}, map[string]any{"style": "italic", "text": "ghost"}}})
	assert.NoError(t, err)
	assert.Error(t, c.CheckSession(ctx, idTest))
	assert.NoError(t, c.CreateSession(ctx, idTest))

	assert.NoError(t, c.CheckSession(ctx, idTest))

	assert.NoError(t, c.SetSessionTags(ctx, idTest, []*brzrpc.Tag{
		{
			Id:     "test",
			Title:  "test",
			Color:  "test",
			Emoji:  "test",
			UserId: "test",
		},
		{
			Id:     "test2",
			Title:  "test2",
			Color:  "test2",
			Emoji:  "test2",
			UserId: "test2",
		},
	}))

	if tgs, err := c.GetSessionTags(ctx, idTest); assert.NoError(t, err) {
		if assert.Equal(t, 2, len(tgs)) {
			assert.Equal(t, "test", tgs[0].Id)
			assert.Equal(t, "test2", tgs[1].Id)
		}
	}

	assert.NoError(t, c.SetSessionNoteList(ctx, idTest, []*brzrpc.NotePart{
		{
			Id:    "test",
			Title: "test",
			Tag: &brzrpc.Tag{
				Id:     "test",
				Title:  "test",
				Color:  "test",
				Emoji:  "test",
				UserId: "test",
			},
			FirstBlock: "test",
			UpdatedAt:  1,
		},
		{
			Id:    "test2",
			Title: "test2",
			Tag: &brzrpc.Tag{
				Id:     "test2",
				Title:  "test2",
				Color:  "test2",
				Emoji:  "test2",
				UserId: "test2",
			},
			FirstBlock: "test2",
			UpdatedAt:  1,
		},
	}))

	if nps, err := c.GetSessionNoteList(ctx, idTest); assert.NoError(t, err) {
		if assert.Equal(t, 2, len(nps)) {
			assert.Equal(t, "test", nps[0].Id)
			assert.Equal(t, "test2", nps[1].Id)
		}
	}

	assert.NoError(t, c.SetSessionNoteTrash(ctx, idTest, []*brzrpc.NotePart{
		{
			Id:    "test",
			Title: "test",
			Tag: &brzrpc.Tag{
				Id:     "test",
				Title:  "test",
				Color:  "test",
				Emoji:  "test",
				UserId: "test",
			},
			FirstBlock: "test",
			UpdatedAt:  1,
		},
		{
			Id:    "test2",
			Title: "test2",
			Tag: &brzrpc.Tag{
				Id:     "test2",
				Title:  "test2",
				Color:  "test2",
				Emoji:  "test2",
				UserId: "test2",
			},
			FirstBlock: "test2",
			UpdatedAt:  1,
		},
	}))

	if nps, err := c.GetSessionNoteTrash(ctx, idTest); assert.NoError(t, err) {
		if assert.Equal(t, 2, len(nps)) {
			assert.Equal(t, "test", nps[0].Id)
			assert.Equal(t, "test2", nps[1].Id)
		}
	}

	assert.NoError(t, c.SetSessionNotes(ctx, idTest, []*brzrpc.NoteWithBlocks{
		{
			Id:        "test",
			Title:     "test",
			CreatedAt: 1,
			UpdatedAt: 1,
			Tag: &brzrpc.Tag{
				Id:     "test",
				Title:  "test",
				Color:  "test",
				Emoji:  "test",
				UserId: "test",
			},
			Author:  "test",
			Editors: []string{},
			Readers: []string{},
			Blocks: []*brzrpc.Block{
				{
					Id:        "test",
					Type:      "test",
					NoteId:    "test",
					Order:     1,
					CreatedAt: 1,
					UpdatedAt: 1,
					IsUsed:    true,
					Data:      s,
				},
				{
					Id:        "test",
					Type:      "test",
					NoteId:    "test",
					Order:     1,
					CreatedAt: 1,
					UpdatedAt: 1,
					IsUsed:    true,
					Data:      s,
				},
			},
		},
		{
			Id:        "test2",
			Title:     "test2",
			CreatedAt: 1,
			UpdatedAt: 1,
			Tag:       nil,
			Author:    "test2",
			Editors:   []string{},
			Readers:   []string{},
			Blocks: []*brzrpc.Block{
				{
					Id:        "test2",
					Type:      "test2",
					NoteId:    "test2",
					Order:     1,
					CreatedAt: 1,
					UpdatedAt: 1,
					IsUsed:    true,
					Data:      s,
				},
				{
					Id:        "test2",
					Type:      "test2",
					NoteId:    "test2",
					Order:     1,
					CreatedAt: 1,
					UpdatedAt: 1,
					IsUsed:    true,
					Data:      s,
				}},
		},
	}))

	if nps, err := c.GetSessionNotes(ctx, idTest); assert.NoError(t, err) {
		if assert.Equal(t, 2, len(nps)) {
			assert.Equal(t, "test", nps[0].Id)
			assert.Equal(t, "test2", nps[1].Id)
		}
	}
	assert.NoError(t, c.SetSessionNotes(ctx, idTest, nil))
	assert.NoError(t, c.SetSessionNoteTrash(ctx, idTest, nil))
	assert.NoError(t, c.SetSessionTags(ctx, idTest, nil))
	assert.NoError(t, c.SetSessionNoteList(ctx, idTest, nil))

	assert.NoError(t, c.deleteSession(context.TODO(), idTest))
	if _, err := c.getSession(context.Background(), idTest); assert.Error(t, err) {

	}
}
