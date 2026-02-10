package notes

import (
	"context"
	"errors"

	"github.com/autumnterror/breezynotes/internal/blocknote/domain2"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/alg"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Get return note by id can return mongo.ErrNotFound
func (a *API) Get(ctx context.Context, idNote, idUser string) (*domain2.Note, error) {
	const op = "notes.Get"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()

	res := a.noteAPI.FindOne(ctx, bson.M{"_id": idNote})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, format.Error(op, domain2.ErrNotFound)
		}
		return nil, format.Error(op, res.Err())
	}

	var n domain2.Note
	err := res.Decode(&n)
	if err != nil {
		return nil, format.Error(op, err)
	}

	res = a.noteTagsAPI.FindOne(ctx, bson.M{"note_id": idNote, "tag.user_id": idUser})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return &n, nil
		}
		return nil, format.Error(op, res.Err())
	}
	var nt domain2.NoteTags
	err = res.Decode(&n)
	if err != nil {
		return nil, format.Error(op, err)
	}

	n.Tag = nt.Tag

	return &n, nil
}

// GetNoteListByUser use func a.blockAPI.GetAsFirst
func (a *API) GetNoteListByUser(ctx context.Context, id string) (*domain2.NoteParts, error) {
	const op = "notes.GetNoteListByUser"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()

	cur, err := a.noteTagsAPI.Find(ctx, bson.M{"tag.user_id": id})
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	noteTag := make(map[string]*domain2.Tag)

	for cur.Next(ctx) {
		nt := domain2.NoteTags{}
		if err = cur.Decode(&nt); err != nil {
			return nil, format.Error(op, err)
		}
		noteTag[nt.NoteId] = nt.Tag
	}

	cur, err = a.noteAPI.Find(ctx,
		bson.M{
			"$or": []bson.M{
				{"author": id},
				{"editors": id},
				{"readers": id},
			},
		},
		options.Find().SetSort(bson.M{"updated_at": -1}),
	)

	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	nts := &domain2.NoteParts{
		Ntps: []*domain2.NotePart{},
	}

	for cur.Next(ctx) {
		var n domain2.Note
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
		var role string
		switch {
		case n.Author == id:
			role = "author"
		case alg.IsIn(id, n.Editors):
			role = "editor"
		case alg.IsIn(id, n.Readers):
			role = "reader"
		}
		nts.Ntps = append(nts.Ntps, &domain2.NotePart{
			Id:         n.Id,
			Title:      n.Title,
			Tag:        noteTag[n.Id],
			FirstBlock: fb,
			UpdatedAt:  n.UpdatedAt,
			Role:       role,
			IsBlog:     n.IsBlog,
			IsPublic:   n.IsPublic,
		})

	}

	return nts, nil
}

// GetNoteListByTag use func a.blockAPI.GetAsFirst
func (a *API) GetNoteListByTag(ctx context.Context, idTag, idUser string) (*domain2.NoteParts, error) {
	const op = "notes.GetNoteListByTag"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()

	nts := &domain2.NoteParts{
		Ntps: []*domain2.NotePart{},
	}

	tag, err := a.tagAPI.Get(ctx, idTag)
	if err != nil {
		if errors.Is(err, domain2.ErrNotFound) {
			return nts, nil
		}
		return nil, format.Error(op, err)
	}

	cur, err := a.noteTagsAPI.Find(ctx, bson.M{"tag._id": idTag})
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	var ntIds []string

	for cur.Next(ctx) {
		nt := domain2.NoteTags{}
		if err = cur.Decode(&nt); err != nil {
			return nil, format.Error(op, err)
		}
		ntIds = append(ntIds, nt.NoteId)
	}

	if len(ntIds) == 0 {
		return nts, nil
	}

	cur, err = a.noteAPI.Find(
		ctx,
		bson.D{{"_id", bson.D{{"$in", ntIds}}}},
		options.Find().SetSort(bson.M{"updated_at": -1}))
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var n domain2.Note
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
		var role string
		switch {
		case n.Author == idUser:
			role = "author"
		case alg.IsIn(idUser, n.Editors):
			role = "editor"
		case alg.IsIn(idUser, n.Readers):
			role = "reader"
		}
		nts.Ntps = append(nts.Ntps, &domain2.NotePart{
			Id:         n.Id,
			Title:      n.Title,
			Tag:        tag,
			FirstBlock: fb,
			UpdatedAt:  n.UpdatedAt,
			Role:       role,
			IsBlog:     n.IsBlog,
			IsPublic:   n.IsPublic,
		})
	}

	return nts, nil
}

// GetAllByUser return note by id author
func (a *API) getAllByUser(ctx context.Context, id string) (*domain2.Notes, error) {
	const op = "notes.getAllByUser"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()

	cur, err := a.noteAPI.Find(ctx, bson.M{"author": id}, options.Find().SetSort(bson.M{"updated_at": -1}))
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	nts := &domain2.Notes{
		Nts: []*domain2.Note{},
	}

	for cur.Next(ctx) {
		var n domain2.Note
		if err = cur.Decode(&n); err != nil {
			return nts, format.Error(op, err)
		}
		nts.Nts = append(nts.Nts, &n)
	}

	return nts, nil
}

func (a *API) Search(ctx context.Context, id, prompt string) (<-chan *domain2.NotePart, error) {
	const op = "notes.Search"

	ctx, done := context.WithTimeout(ctx, domain2.WaitTime)
	defer done()

	notesChan := make(chan *domain2.NotePart)

	go func() {
		defer close(notesChan)

		cur, err := a.noteAPI.Find(
			ctx,
			bson.M{"author": id, "title": bson.M{"$regex": prompt, "$options": "i"}},
			options.Find().SetSort(bson.M{"updated_at": -1}),
		)
		if err == nil {
			for cur.Next(ctx) {
				var n domain2.Note
				if err = cur.Decode(&n); err != nil {
					continue
				}
				select {
				case <-ctx.Done():
					cur.Close(ctx)
					return
				default:
					notesChan <- &domain2.NotePart{
						Id:    n.Id,
						Title: n.Title,
						//Tag:        n.Tag,
						FirstBlock: "",
						UpdatedAt:  n.UpdatedAt,
						IsBlog:     n.IsBlog,
						IsPublic:   n.IsPublic,
					}
				}
			}
			cur.Close(ctx)
		}
	}()

	return notesChan, nil
}
