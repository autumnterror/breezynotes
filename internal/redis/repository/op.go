package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/redis/go-redis/v9"
)

func (s *Client) GetSessionNoteList(ctx context.Context, id string) ([]*brzrpc.NotePart, error) {
	us, err := s.getSession(ctx, id)
	if err != nil {
		return nil, err
	}
	if us == nil {
		return nil, nil
	}
	return us.NoteParts, nil
}

func (s *Client) SetSessionNoteList(ctx context.Context, id string, noteParts []*brzrpc.NotePart) error {
	us, err := s.getSession(ctx, id)
	if err != nil {
		return err
	}
	if us == nil {
		us = &userSession{Id: id}
	}
	us.NoteParts = noteParts
	return s.saveSession(ctx, us)
}

func (s *Client) GetSessionNoteTrash(ctx context.Context, id string) ([]*brzrpc.NotePart, error) {
	us, err := s.getSession(ctx, id)
	if err != nil {
		return nil, err
	}
	if us == nil {
		return nil, nil
	}
	return us.NoteTrash, nil
}

func (s *Client) SetSessionNoteTrash(ctx context.Context, id string, noteTrash []*brzrpc.NotePart) error {
	us, err := s.getSession(ctx, id)
	if err != nil {
		return err
	}
	if us == nil {
		us = &userSession{Id: id}
	}
	us.NoteTrash = noteTrash
	return s.saveSession(ctx, us)
}

func (s *Client) GetSessionNotes(ctx context.Context, id string) ([]*brzrpc.NoteWithBlocks, error) {
	us, err := s.getSession(ctx, id)
	if err != nil {
		return nil, err
	}
	if us == nil {
		return nil, nil
	}
	return us.Notes, nil
}

func (s *Client) SetSessionNotes(ctx context.Context, id string, notes []*brzrpc.NoteWithBlocks) error {
	us, err := s.getSession(ctx, id)
	if err != nil {
		return err
	}
	if us == nil {
		us = &userSession{Id: id}
	}
	us.Notes = notes

	if err := s.saveSession(ctx, us); err != nil {
		return err
	}

	sessionKey := userSessionKey(id)

	pipe := s.Rdb.Pipeline()
	for _, n := range notes {
		if n == nil || n.Id == "" {
			continue
		}
		pipe.SAdd(ctx, noteSessionsKey(n.Id), sessionKey)

		pipe.Expire(ctx, noteSessionsKey(n.Id), ExpSession*2)
	}
	_, _ = pipe.Exec(ctx)

	return nil
}

func (s *Client) GetSessionTags(ctx context.Context, id string) ([]*brzrpc.Tag, error) {
	us, err := s.getSession(ctx, id)
	if err != nil {
		return nil, err
	}
	if us == nil {
		return nil, nil
	}
	return us.Tags, nil
}

func (s *Client) SetSessionTags(ctx context.Context, id string, tags []*brzrpc.Tag) error {
	us, err := s.getSession(ctx, id)
	if err != nil {
		return err
	}
	if us == nil {
		us = &userSession{Id: id}
	}
	us.Tags = tags
	return s.saveSession(ctx, us)
}

func (s *Client) CleanNoteById(ctx context.Context, noteID string) error {
	const op = "redis.CleanNoteById"
	if noteID == "" {
		return format.Error(op, fmt.Errorf("noteID is empty"))
	}

	idxKey := noteSessionsKey(noteID)

	sessionKeys, err := s.Rdb.SMembers(ctx, idxKey).Result()
	if err != nil {
		return format.Error(op, fmt.Errorf("redis SMEMBERS %s: %w", idxKey, err))
	}
	if len(sessionKeys) == 0 {
		_ = s.Rdb.Del(ctx, idxKey).Err()
		return nil
	}

	getPipe := s.Rdb.Pipeline()
	getCmds := make([]*redis.StringCmd, 0, len(sessionKeys))
	for _, sk := range sessionKeys {
		getCmds = append(getCmds, getPipe.Get(ctx, sk))
	}
	_, _ = getPipe.Exec(ctx)

	setPipe := s.Rdb.Pipeline()

	for i, sk := range sessionKeys {
		data, getErr := getCmds[i].Bytes()

		if getErr == redis.Nil {
			setPipe.SRem(ctx, idxKey, sk)
			continue
		}
		if getErr != nil {
			continue
		}

		var rs userSessionRedis
		if err := json.Unmarshal(data, &rs); err != nil {
			setPipe.Del(ctx, sk)
			setPipe.SRem(ctx, idxKey, sk)
			continue
		}

		rs.Notes = nil

		newData, err := json.Marshal(&rs)
		if err != nil {
			continue
		}
		setPipe.Set(ctx, sk, newData, ExpSession)

		setPipe.SRem(ctx, idxKey, sk)
	}

	_, _ = setPipe.Exec(ctx)
	_ = s.Rdb.Del(ctx, idxKey).Err()

	return nil
}
