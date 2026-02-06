package notes

import (
	"context"

	"github.com/autumnterror/breezynotes/internal/blocknote/domain"

	"github.com/autumnterror/utils_go/pkg/utils/format"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// CleanTrash rm all in a.trash() by uid and field author
func (a *API) CleanTrash(ctx context.Context, uid string) error {
	const op = "notes.CleanTrash"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	_, err := a.trashAPI.DeleteMany(ctx, bson.M{"author": uid})
	if err != nil {
		return format.Error(op, err)
	}

	return nil
}

// GetNotesFromTrash by author id
func (a *API) GetNotesFromTrash(ctx context.Context, uid string) (*domain.NoteParts, error) {
	const op = "notes.CleanTrash"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	cur, err := a.trashAPI.Find(ctx, bson.M{"author": uid})
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	pts := &domain.NoteParts{
		Ntps: []*domain.NotePart{},
	}

	for cur.Next(ctx) {
		var n domain.Note
		if err = cur.Decode(&n); err != nil {
			return pts, format.Error(op, err)
		}
		fb := ""
		if len(n.Blocks) > 0 {
			nfb, err := a.blockAPI.GetAsFirst(ctx, n.Blocks[0])
			if err == nil {
				fb = nfb
			}
		}
		np := domain.NotePart{
			Id:         n.Id,
			Title:      n.Title,
			Tag:        n.Tag,
			FirstBlock: fb,
			UpdatedAt:  n.UpdatedAt,
		}
		pts.Ntps = append(pts.Ntps, &np)
	}

	return pts, nil
}

// GetNotesFullFromTrash by author id
func (a *API) GetNotesFullFromTrash(ctx context.Context, uid string) (*domain.Notes, error) {
	const op = "notes.CleanTrash"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	cur, err := a.trashAPI.Find(ctx, bson.M{"author": uid})
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer cur.Close(ctx)

	pts := &domain.Notes{
		Nts: make([]*domain.Note, 0),
	}

	for cur.Next(ctx) {
		var n domain.Note
		if err = cur.Decode(&n); err != nil {
			return pts, format.Error(op, err)
		}

		pts.Nts = append(pts.Nts, &n)
	}

	return pts, nil
}

// ToTrash get note from a.Notes() Insert on a.trash() and remove from a.Notes()
func (a *API) ToTrash(ctx context.Context, id string) error {
	const op = "notes.ToTrash"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	n, err := a.Get(ctx, id)
	if err != nil {
		return format.Error(op, err)
	}

	if _, err := a.trashAPI.InsertOne(ctx, n); err != nil {
		return format.Error(op, err)
	}

	if err := a.delete(ctx, n.Id); err != nil {
		return format.Error(op, err)
	}
	return nil
}

// ToTrashAll get notes from a.Notes() Insert on a.trash() and remove from a.Notes()
func (a *API) ToTrashAll(ctx context.Context, idUser string) error {
	const op = "notes.ToTrashAll"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	n, err := a.getAllByUser(ctx, idUser)
	if err != nil {
		return format.Error(op, err)
	}
	if n == nil {
		return nil
	}
	if len(n.Nts) == 0 {
		return nil
	}

	if _, err := a.trashAPI.InsertMany(ctx, n.Nts); err != nil {
		return format.Error(op, err)
	}

	var ids []string
	for _, n := range n.Nts {
		ids = append(ids, n.Id)
	}

	if err := a.deleteMany(ctx, ids); err != nil {
		return format.Error(op, err)
	}
	return nil
}

// FromTrash remove note in a.trash() and Insert in a.Notes()
func (a *API) FromTrash(ctx context.Context, id string) error {
	const op = "notes.FromTrash"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	res := a.trashAPI.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return format.Error(op, res.Err())
	}

	var n domain.Note
	if err := res.Decode(&n); err != nil {
		return format.Error(op, err)
	}

	if err := a.insert(ctx, &n); err != nil {
		return format.Error(op, err)
	}

	_, err := a.trashAPI.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return format.Error(op, err)
	}

	return nil
}

// FindOnTrash return note by id from trash
func (a *API) FindOnTrash(ctx context.Context, id string) (*domain.Note, error) {
	const op = "notes.FindOnTrash"

	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	res := a.trashAPI.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return nil, format.Error(op, res.Err())
	}

	var n domain.Note
	err := res.Decode(&n)
	if err != nil {
		return nil, format.Error(op, err)
	}

	return &n, nil
}
