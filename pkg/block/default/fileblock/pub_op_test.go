package fileblock

import (
	"context"
	"testing"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/domain"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	d   = Driver{}
	ctx = context.Background()
)

func TestGetAsFirst(t *testing.T) {
	t.Parallel()
	block := testBlock()
	assert.Equal(t, "file", d.GetAsFirst(ctx, block))

	blockNil := testBlockNil()
	assert.Equal(t, "file", d.GetAsFirst(ctx, blockNil))

	assert.Equal(t, "file", d.GetAsFirst(ctx, nil))
}

func TestOp(t *testing.T) {
	t.Parallel()
	t.Run("change_src", func(t *testing.T) {
		block := testBlock()
		newSrc := "path/to/another/document.docx"

		data, err := d.Op(ctx, block, "change_src", map[string]any{
			"new_src": newSrc,
		})
		if assert.NoError(t, err) {
			s, err := structpb.NewStruct(data)
			if assert.NoError(t, err) {
				block.Data = s
				fb, err := domain.FromUnifiedToFileBlock(block)
				if assert.NoError(t, err) {
					assert.Equal(t, newSrc, fb.Data.Src)
				}
			}
		}
	})
}

func TestOpNil(t *testing.T) {
	t.Parallel()
	block := testBlockNil()
	newSrc := "path/to/a/new/file.zip"

	// Operation on a block with nil Data should initialize it
	data, err := d.Op(ctx, block, "change_src", map[string]any{
		"new_src": newSrc,
	})
	if assert.NoError(t, err) {
		s, err := structpb.NewStruct(data)
		if assert.NoError(t, err) {
			block.Data = s
			fb, err := domain.FromUnifiedToFileBlock(block)
			if assert.NoError(t, err) {
				assert.NotNil(t, fb.Data)
				assert.Equal(t, newSrc, fb.Data.Src)
			}
		}
	}
}

func TestChangeType(t *testing.T) {
	t.Parallel()
	t.Run(domain.TextBlockType, func(t *testing.T) {
		block := testBlock()
		// fb, _ := domain.FromUnifiedToFileBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.TextBlockType))
		tb, err := domain.FromUnifiedToTextBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "from file", tb.Data.PlainText())
	})

	t.Run(domain.QuoteBlockType, func(t *testing.T) {
		block := testBlock()
		// fb, _ := domain.FromUnifiedToFileBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.QuoteBlockType))
		qb, err := domain.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "from file", qb.Data.Text)
	})

	t.Run(domain.ImgBlockType, func(t *testing.T) {
		block := testBlock()
		// fb, _ := domain.FromUnifiedToFileBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ImgBlockType))
		ib, err := domain.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", ib.Data.Src)
		assert.Equal(t, "from file", ib.Data.Alt)
	})
}

func testBlock() *brzrpc.Block {
	fb := domain.FileBlock{
		Id:     "test_file",
		Type:   domain.FileBlockType,
		NoteId: "test_note",
		Data: &domain.FileData{
			Src: "path/to/my/file.pdf",
		},
	}
	unif, _ := fb.ToUnified()
	return unif
}

func testBlockNil() *brzrpc.Block {
	fb := domain.FileBlock{
		Id:     "test_file_nil",
		Type:   domain.FileBlockType,
		NoteId: "test_note",
		Data:   nil,
	}
	unif, _ := fb.ToUnified()
	return unif
}
