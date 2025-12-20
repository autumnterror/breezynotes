package redis

import (
	"context"
	"github.com/autumnterror/breezynotes/api/proto/gen"
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

	return s.saveSession(ctx, us)
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
