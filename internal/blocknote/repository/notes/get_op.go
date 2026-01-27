package notes

import (
	"context"
	"errors"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Get return note by id can return mongo.ErrNotFound
func (a *API) Get(ctx context.Context, id string) (*domain.Note, error) {
	const op = "notes.Get"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	res := a.noteAPI.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, format.Error(op, domain.ErrNotFound)
		}
		return nil, format.Error(op, res.Err())
	}

	var n domain.Note
	err := res.Decode(&n)
	if err != nil {
		return nil, format.Error(op, err)
	}

	return &n, nil
}

// GetNoteListByUser use func a.blockAPI.GetAsFirst
func (a *API) GetNoteListByUser(ctx context.Context, id string) (*domain.NoteParts, error) {
	const op = "notes.GetNoteListByUser"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	cur, err := a.noteAPI.Find(ctx, bson.M{"author": id}, options.Find().SetSort(bson.M{"updated_at": -1}))

	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	nts := &domain.NoteParts{
		Ntps: []*domain.NotePart{},
	}

	for cur.Next(ctx) {
		var n domain.Note
		if err = cur.Decode(&n); err != nil {
			return nil, format.Error(op, err)
		}
		fb := ""
		if len(n.Blocks) > 0 {
			nfb, err := a.blockAPI.GetAsFirst(ctx, n.Blocks[0])
			if err == nil {
				fb = nfb
			} else {
				log.Warn(op, "get as first", err)
			}
		}
		nts.Ntps = append(nts.Ntps, &domain.NotePart{
			Id:         n.Id,
			Title:      n.Title,
			Tag:        n.Tag,
			FirstBlock: fb,
			UpdatedAt:  n.UpdatedAt,
		})

	}

	return nts, nil
}

// GetNoteListByTag use func a.blockAPI.GetAsFirst
func (a *API) GetNoteListByTag(ctx context.Context, idTag, idUser string) (*domain.NoteParts, error) {
	const op = "notes.GetNoteListByUser"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	cur, err := a.noteAPI.Find(ctx, bson.M{"tag._id": idTag, "author": idUser}, options.Find().SetSort(bson.M{"updated_at": -1}))
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	nts := &domain.NoteParts{
		Ntps: []*domain.NotePart{},
	}

	idx := 0
	for cur.Next(ctx) {
		var n domain.Note
		if err = cur.Decode(&n); err != nil {
			return nts, format.Error(op, err)
		}
		fb := ""
		if len(n.Blocks) > 0 {
			nfb, err := a.blockAPI.GetAsFirst(ctx, n.Blocks[0])
			if err == nil {
				fb = nfb
			}
		}
		nts.Ntps = append(nts.Ntps, &domain.NotePart{
			Id:         n.Id,
			Title:      n.Title,
			Tag:        n.Tag,
			FirstBlock: fb,
			UpdatedAt:  n.UpdatedAt,
		})
		idx++
	}

	return nts, nil
}

// GetAllByUser return note by id author
func (a *API) getAllByUser(ctx context.Context, id string) (*domain.Notes, error) {
	const op = "notes.Get"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	cur, err := a.noteAPI.Find(ctx, bson.M{"author": id}, options.Find().SetSort(bson.M{"updated_at": -1}))
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	nts := &domain.Notes{
		Nts: []*domain.Note{},
	}

	for cur.Next(ctx) {
		var n domain.Note
		if err = cur.Decode(&n); err != nil {
			return nts, format.Error(op, err)
		}
		nts.Nts = append(nts.Nts, &n)
	}

	return nts, nil
}
