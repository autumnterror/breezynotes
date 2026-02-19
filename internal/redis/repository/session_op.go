package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

type userSession struct {
	Id        string                   `json:"id"`
	NoteParts []*brzrpc.NotePart       `json:"note_parts"`
	Notes     []*brzrpc.NoteWithBlocks `json:"notes"`
	Tags      []*brzrpc.Tag            `json:"tags"`
	NoteTrash []*brzrpc.NotePart       `json:"note_trash"`
}

type userSessionRedis struct {
	Id        string            `json:"id"`
	NoteParts []json.RawMessage `json:"note_parts"`
	Notes     []json.RawMessage `json:"notes"`
	Tags      []json.RawMessage `json:"tags"`
	NoteTrash []json.RawMessage `json:"note_trash"`
}

func userSessionKey(id string) string {
	return userSessionKeyPrefix + id
}

const (
	ExpSession            = time.Hour
	userSessionKeyPrefix  = "user_session:"
	noteSessionsKeyPrefix = "note_sessions:"
)

var (
	ErrNotFound = fmt.Errorf("session not found")
)

func noteSessionsKey(noteID string) string {
	return noteSessionsKeyPrefix + noteID
}

func (s *Client) saveSession(ctx context.Context, us *userSession) error {
	const op = "redis.saveSession"
	if us == nil {
		return format.Error(op, fmt.Errorf("userSession is nil"))
	}
	if us.Id == "" {
		return format.Error(op, fmt.Errorf("userSession.Id is empty"))
	}

	pj := protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}

	rs := userSessionRedis{
		Id:        us.Id,
		NoteParts: make([]json.RawMessage, 0, len(us.NoteParts)),
		NoteTrash: make([]json.RawMessage, 0, len(us.NoteTrash)),
		Notes:     make([]json.RawMessage, 0, len(us.Notes)),
		Tags:      make([]json.RawMessage, 0, len(us.Tags)),
	}

	for _, np := range us.NoteParts {
		if np == nil {
			rs.NoteParts = append(rs.NoteParts, nil)
			continue
		}
		b, err := pj.Marshal(np)
		if err != nil {
			return format.Error(op, fmt.Errorf("marshal NotePart: %w", err))
		}
		rs.NoteParts = append(rs.NoteParts, b)
	}

	for _, np := range us.NoteTrash {
		if np == nil {
			rs.NoteTrash = append(rs.NoteTrash, nil)
			continue
		}
		b, err := pj.Marshal(np)
		if err != nil {
			return format.Error(op, fmt.Errorf("marshal NoteTrash: %w", err))
		}
		rs.NoteTrash = append(rs.NoteTrash, b)
	}

	for _, n := range us.Notes {
		if n == nil {
			rs.Notes = append(rs.Notes, nil)
			continue
		}
		b, err := pj.Marshal(n)
		if err != nil {
			return format.Error(op, fmt.Errorf("marshal NoteWithBlocks: %w", err))
		}
		rs.Notes = append(rs.Notes, b)
	}

	for _, t := range us.Tags {
		if t == nil {
			rs.Tags = append(rs.Tags, nil)
			continue
		}
		b, err := pj.Marshal(t)
		if err != nil {
			return format.Error(op, fmt.Errorf("marshal Tag: %w", err))
		}
		rs.Tags = append(rs.Tags, b)
	}

	data, err := json.Marshal(&rs)
	if err != nil {
		return format.Error(op, fmt.Errorf("marshal userSessionRedis: %w", err))
	}

	key := userSessionKey(us.Id)
	if err := s.Rdb.Set(ctx, key, data, ExpSession).Err(); err != nil {
		return format.Error(op, fmt.Errorf("redis SET %s: %w", key, err))
	}

	return nil
}

func (s *Client) CreateSession(ctx context.Context, id string) error {
	const op = "redis.CheckSession"
	key := userSessionKey(id)

	rs := userSessionRedis{
		Id:        id,
		NoteParts: nil,
		Notes:     nil,
		Tags:      nil,
	}
	data, err := json.Marshal(&rs)
	if err != nil {
		return format.Error(op, fmt.Errorf("marshal userSessionRedis: %w", err))
	}

	if err := s.Rdb.Set(ctx, key, data, ExpSession).Err(); err != nil {
		return format.Error(op, fmt.Errorf("redis SET %s: %w", key, err))
	}

	return nil
}

func (s *Client) CheckSession(ctx context.Context, id string) error {
	const op = "redis.CheckSession"
	key := userSessionKey(id)
	_, err := s.Rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return format.Error(op, ErrNotFound)
		}
		return format.Error(op, fmt.Errorf("redis GET %s: %w", key, err))
	}
	return nil
}

// GetSession достаёт всю сессию из Redis по Id.
func (s *Client) getSession(ctx context.Context, id string) (*userSession, error) {
	const op = "redis.getSession"
	if id == "" {
		return nil, format.Error(op, fmt.Errorf("id is empty"))
	}
	key := userSessionKey(id)

	data, err := s.Rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, format.Error(op, fmt.Errorf("session not found"))
		}
		return nil, format.Error(op, fmt.Errorf("redis GET %s: %w", key, err))
	}

	var rs userSessionRedis
	if err := json.Unmarshal(data, &rs); err != nil {
		return nil, format.Error(op, fmt.Errorf("unmarshal userSessionRedis: %w", err))
	}

	pj := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	us := &userSession{
		Id:        rs.Id,
		NoteParts: make([]*brzrpc.NotePart, 0, len(rs.NoteParts)),
		NoteTrash: make([]*brzrpc.NotePart, 0, len(rs.NoteTrash)),
		Notes:     make([]*brzrpc.NoteWithBlocks, 0, len(rs.Notes)),
		Tags:      make([]*brzrpc.Tag, 0, len(rs.Tags)),
	}

	for _, raw := range rs.NoteParts {
		if raw == nil {
			us.NoteParts = append(us.NoteParts, nil)
			continue
		}
		np := &brzrpc.NotePart{}
		if err := pj.Unmarshal(raw, np); err != nil {
			return nil, format.Error(op, fmt.Errorf("unmarshal NotePart: %w", err))
		}
		us.NoteParts = append(us.NoteParts, np)
	}

	for _, raw := range rs.NoteTrash {
		if raw == nil {
			us.NoteTrash = append(us.NoteTrash, nil)
			continue
		}
		np := &brzrpc.NotePart{}
		if err := pj.Unmarshal(raw, np); err != nil {
			return nil, format.Error(op, fmt.Errorf("unmarshal NotePart: %w", err))
		}
		us.NoteTrash = append(us.NoteTrash, np)
	}

	for _, raw := range rs.Notes {
		if raw == nil {
			us.Notes = append(us.Notes, nil)
			continue
		}
		n := &brzrpc.NoteWithBlocks{}
		if err := pj.Unmarshal(raw, n); err != nil {
			return nil, format.Error(op, fmt.Errorf("unmarshal NoteWithBlocks: %w", err))
		}
		us.Notes = append(us.Notes, n)
	}

	// Tags
	for _, raw := range rs.Tags {
		if raw == nil {
			us.Tags = append(us.Tags, nil)
			continue
		}
		t := &brzrpc.Tag{}
		if err := pj.Unmarshal(raw, t); err != nil {
			return nil, format.Error(op, fmt.Errorf("unmarshal Tag: %w", err))
		}
		us.Tags = append(us.Tags, t)
	}

	return us, nil
}

// deleteSession удаляет сессию по Id.
func (s *Client) deleteSession(ctx context.Context, id string) error {
	const op = "redis.deleteSession"
	if id == "" {
		return format.Error(op, fmt.Errorf("id is empty"))
	}
	key := userSessionKey(id)
	if err := s.Rdb.Del(ctx, key).Err(); err != nil {
		return format.Error(op, fmt.Errorf("redis DEL %s: %w", key, err))
	}
	return nil
}
