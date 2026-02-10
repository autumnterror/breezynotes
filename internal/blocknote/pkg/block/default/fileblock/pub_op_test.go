package fileblock

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2/domainblocks"
	"testing"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
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
				fb, err := domainblocks.FromUnifiedToFileBlock(block)
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
			fb, err := domainblocks.FromUnifiedToFileBlock(block)
			if assert.NoError(t, err) {
				assert.NotNil(t, fb.Data)
				assert.Equal(t, newSrc, fb.Data.Src)
			}
		}
	}
}

func TestChangeType(t *testing.T) {
	t.Parallel()
	t.Run(domainblocks.TextBlockType, func(t *testing.T) {
		block := testBlock()
		// fb, _ := domain.FromUnifiedToFileBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.TextBlockType))
		tb, err := domainblocks.FromUnifiedToTextBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "from file", tb.Data.PlainText())
	})

	t.Run(domainblocks.QuoteBlockType, func(t *testing.T) {
		block := testBlock()
		// fb, _ := domain.FromUnifiedToFileBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.QuoteBlockType))
		qb, err := domainblocks.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "from file", qb.Data.Text)
	})

	t.Run(domainblocks.ImgBlockType, func(t *testing.T) {
		block := testBlock()
		// fb, _ := domain.FromUnifiedToFileBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ImgBlockType))
		ib, err := domainblocks.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", ib.Data.Src)
		assert.Equal(t, "from file", ib.Data.Alt)
	})
}

func testBlock() *brzrpc.Block {
	fb := domainblocks.FileBlock{
		Id:     "test_file",
		Type:   domainblocks.FileBlockType,
		NoteId: "test_note",
		Data: &domainblocks.FileData{
			Src: "path/to/my/file.pdf",
		},
	}
	unif, _ := fb.ToUnified()
	return unif
}

func testBlockNil() *brzrpc.Block {
	fb := domainblocks.FileBlock{
		Id:     "test_file_nil",
		Type:   domainblocks.FileBlockType,
		NoteId: "test_note",
		Data:   nil,
	}
	unif, _ := fb.ToUnified()
	return unif
}
