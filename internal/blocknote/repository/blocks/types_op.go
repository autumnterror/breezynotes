package blocks

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/block"
	"github.com/autumnterror/utils_go/pkg/utils/format"
)

// GetAsFirst universal func that get block, calls func GetAsFirst by field _type. NEED TO REGISTER TYPE BEFORE USE
func (a *API) GetAsFirst(ctx context.Context, id string) (string, error) {
	const op = "blocks.GetAsFirst"
	ctx, done := context.WithTimeout(ctx, domain.WaitTime)
	defer done()

	b, err := a.Get(ctx, id)
	if err != nil {
		return "", format.Error(op, err)
	}
	if block.Registry[b.Type] == nil {
		return "", domain.ErrTypeNotDefined
	}

	return block.Registry[b.Type].GetAsFirst(ctx, domain.FromBlockDb(b)), nil
}
