package quoteblock

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/domain"
	"github.com/autumnterror/utils_go/pkg/log"
	format "github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
	"testing"
)

var (
	d   = Driver{}
	ctx = context.Background()
)

func TestGetAsFirst(t *testing.T) {
	t.Parallel()
	block := testBlock()
	assert.Equal(t, testText, d.GetAsFirst(ctx, block))

	blockNil := testBlockNil()
	assert.Equal(t, "", d.GetAsFirst(ctx, blockNil))
}

func TestOp(t *testing.T) {
	t.Parallel()
	t.Run("change_text", func(t *testing.T) {
		block := testBlock()
		newText := "Be the change that you wish to see in the world."

		data, err := d.Op(ctx, block, "change_text", map[string]any{
			"new_text": newText,
		})
		if assert.NoError(t, err) {
			s, _ := structpb.NewStruct(data)
			block.Data = s
			qb, _ := domain.FromUnifiedToQuoteBlock(block)
			assert.Equal(t, newText, qb.Data.Text)
		}
	})
}

func TestOpNilData(t *testing.T) {
	t.Parallel()
	block := testBlockNil()
	newText := "A new quote begins."

	data, err := d.Op(ctx, block, "change_text", map[string]any{
		"new_text": newText,
	})
	if assert.NoError(t, err) {
		assert.Nil(t, data)
	}
}

func TestChangeType(t *testing.T) {
	t.Parallel()
	t.Run(domain.TextBlockType, func(t *testing.T) {
		block := testBlock()
		qt, _ := domain.FromUnifiedToQuoteBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.TextBlockType))
		tb, err := domain.FromUnifiedToTextBlock(block)
		assert.NoError(t, err)
		log.Println("after change on text", format.Struct(tb))
		assert.Equal(t, qt.Data.Text, tb.Data.PlainText())
	})

	t.Run(domain.ListBlockUnorderedType, func(t *testing.T) {
		block := testBlock()
		qt, _ := domain.FromUnifiedToQuoteBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockUnorderedType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered img", format.Struct(lb))
		assert.Equal(t, qt.Data.Text, lb.Data.TextData.PlainText())
		assert.Equal(t, domain.ListBlockUnorderedType, lb.Data.Type)
	})
	t.Run(domain.ListBlockToDoType, func(t *testing.T) {
		block := testBlock()
		qt, _ := domain.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockToDoType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered", format.Struct(block))
		assert.Equal(t, qt.Data.Text, lb.Data.TextData.PlainText())
		assert.Equal(t, domain.ListBlockToDoType, lb.Data.Type)
	})
	t.Run(domain.ListBlockOrderedType, func(t *testing.T) {
		block := testBlock()
		qt, _ := domain.FromUnifiedToQuoteBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockOrderedType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on ordered", format.Struct(block))
		assert.Equal(t, qt.Data.Text, lb.Data.TextData.PlainText())
		assert.Equal(t, domain.ListBlockOrderedType, lb.Data.Type)
	})

	t.Run(domain.HeaderBlockType1, func(t *testing.T) {
		block := testBlock()
		qt, _ := domain.FromUnifiedToQuoteBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType1))
		hb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_1", format.Struct(hb))
		assert.Equal(t, qt.Data.Text, hb.Data.TextData.PlainText())
		assert.Equal(t, uint(1), hb.Data.Level)
	})
	t.Run(domain.HeaderBlockType2, func(t *testing.T) {
		block := testBlock()
		qt, _ := domain.FromUnifiedToQuoteBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType2))
		hb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_2", format.Struct(block))
		assert.Equal(t, qt.Data.Text, hb.Data.TextData.PlainText())
		assert.Equal(t, uint(2), hb.Data.Level)
	})

	t.Run(domain.HeaderBlockType3, func(t *testing.T) {
		block := testBlock()
		qt, _ := domain.FromUnifiedToQuoteBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType3))
		hb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_3", format.Struct(block))
		assert.Equal(t, qt.Data.Text, hb.Data.TextData.PlainText())
		assert.Equal(t, uint(3), hb.Data.Level)
	})
	t.Run(domain.FileBlockType, func(t *testing.T) {
		block := testBlock()
		assert.NoError(t, d.ChangeType(ctx, block, domain.FileBlockType))
		fb, err := domain.FromUnifiedToFileBlock(block)
		assert.NoError(t, err)
		log.Println("after change on file", format.Struct(block))
		assert.Equal(t, "", fb.Data.Src)
	})

	t.Run(domain.ImgBlockType, func(t *testing.T) {
		block := testBlock()
		qt, _ := domain.FromUnifiedToQuoteBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ImgBlockType))
		ib, err := domain.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		log.Println("after change on img", format.Struct(block))
		assert.Equal(t, qt.Data.Text, ib.Data.Alt)
	})

	t.Run(domain.LinkBlockType, func(t *testing.T) {
		block := testBlock()
		qt, _ := domain.FromUnifiedToQuoteBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.LinkBlockType))
		linkB, err := domain.FromUnifiedToLinkBlock(block)
		assert.NoError(t, err)
		log.Println("after change on link", format.Struct(block))
		assert.Equal(t, qt.Data.Text, linkB.Data.Text)
	})

	t.Run(domain.QuoteBlockType, func(t *testing.T) {
		block := testBlock()
		qt, _ := domain.FromUnifiedToQuoteBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.QuoteBlockType))
		qb, err := domain.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		log.Println("after change on quote", format.Struct(qb))
		assert.Equal(t, qt.Data.Text, qb.Data.Text)
	})
}

func TestChangeTypeNil(t *testing.T) {
	t.Parallel()
	t.Run(domain.TextBlockType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domain.TextBlockType))
		tb, err := domain.FromUnifiedToTextBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", tb.Data.PlainText())
	})

	t.Run(domain.ListBlockToDoType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockToDoType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, domain.ListBlockToDoType, lb.Data.Type)
		assert.Equal(t, "", lb.Data.TextData.PlainText())
	})

	t.Run(domain.CodeBlockType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domain.CodeBlockType))
		code, err := domain.FromUnifiedToCodeBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", code.Data.Text)
	})

	t.Run(domain.HeaderBlockType1, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType1))
		hb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), hb.Data.Level)
		assert.Equal(t, "", hb.Data.TextData.PlainText())
	})

	t.Run(domain.ImgBlockType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domain.ImgBlockType))
		ib, err := domain.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", ib.Data.Alt)
		assert.Equal(t, "", ib.Data.Src)
	})

	t.Run(domain.LinkBlockType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domain.LinkBlockType))
		linkB, err := domain.FromUnifiedToLinkBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", linkB.Data.Text)
	})

	t.Run(domain.QuoteBlockType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domain.QuoteBlockType))
		qb, err := domain.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", qb.Data.Text)
	})
}

// --- Helpers ---

var testText = "The only true wisdom is in knowing you know nothing."

func testBlock() *brzrpc.Block {
	qb := domain.QuoteBlock{
		Id:   "test_quote",
		Type: domain.QuoteBlockType,
		Data: &domain.QuoteData{
			Text: testText,
		},
	}
	unif, _ := qb.ToUnified()
	return unif
}

func testBlockNil() *brzrpc.Block {
	qb := domain.QuoteBlock{
		Id:   "test_quote_nil",
		Type: domain.QuoteBlockType,
		Data: nil,
	}
	unif, _ := qb.ToUnified()
	return unif
}
