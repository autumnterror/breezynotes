package imgblock

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2/domainblocks"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
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
	assert.Equal(t, "A beautiful landscape", d.GetAsFirst(ctx, block))

	blockNil := testBlockNil()
	assert.Equal(t, "", d.GetAsFirst(ctx, blockNil))
}

func TestOp(t *testing.T) {
	t.Parallel()
	t.Run("change_src", func(t *testing.T) {
		block := testBlock()
		newSrc := "https://example.com/new_image.png"
		data, err := d.Op(ctx, block, "change_src", map[string]any{"new_src": newSrc})
		if assert.NoError(t, err) {
			s, _ := structpb.NewStruct(data)
			block.Data = s
			ib, _ := domainblocks.FromUnifiedToImgBlock(block)
			assert.Equal(t, newSrc, ib.Data.Src)
			assert.Equal(t, testAlt, ib.Data.Alt)
		}
	})

	t.Run("change_alt", func(t *testing.T) {
		block := testBlock()
		newAlt := "A stunning mountain view"
		data, err := d.Op(ctx, block, "change_alt", map[string]any{"new_alt": newAlt})
		if assert.NoError(t, err) {
			s, _ := structpb.NewStruct(data)
			block.Data = s
			ib, _ := domainblocks.FromUnifiedToImgBlock(block)
			assert.Equal(t, newAlt, ib.Data.Alt)
			assert.Equal(t, "https://example.com/image.jpg", ib.Data.Src)
		}
	})
}

func TestOpNilData(t *testing.T) {
	t.Parallel()
	t.Run("change_src", func(t *testing.T) {
		block := testBlockNil()
		newSrc := "https://example.com/new_image.png"
		data, err := d.Op(ctx, block, "change_src", map[string]any{"new_src": newSrc})
		if assert.NoError(t, err) {
			assert.Nil(t, data)
		}
	})

	t.Run("change_alt", func(t *testing.T) {
		block := testBlockNil()
		newAlt := "A stunning mountain view"
		data, err := d.Op(ctx, block, "change_alt", map[string]any{"new_alt": newAlt})
		if assert.NoError(t, err) {
			assert.Nil(t, data)
		}
	})
}

func TestChangeType(t *testing.T) {
	t.Parallel()
	t.Run(domainblocks.TextBlockType, func(t *testing.T) {
		block := testBlock()
		img, _ := domainblocks.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.TextBlockType))
		tb, err := domainblocks.FromUnifiedToTextBlock(block)
		assert.NoError(t, err)
		log.Println("after change on text", format.Struct(tb))
		assert.Equal(t, img.Data.Alt, tb.Data.PlainText())
	})

	t.Run(domainblocks.ListBlockUnorderedType, func(t *testing.T) {
		block := testBlock()
		img, _ := domainblocks.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockUnorderedType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered img", format.Struct(lb))
		assert.Equal(t, img.Data.Alt, lb.Data.TextData.PlainText())
		assert.Equal(t, domainblocks.ListBlockUnorderedType, lb.Data.Type)
	})
	t.Run(domainblocks.ListBlockToDoType, func(t *testing.T) {
		block := testBlock()
		img, _ := domainblocks.FromUnifiedToImgBlock(block)
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockToDoType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered", format.Struct(block))
		assert.Equal(t, img.Data.Alt, lb.Data.TextData.PlainText())
		assert.Equal(t, domainblocks.ListBlockToDoType, lb.Data.Type)
	})
	t.Run(domainblocks.ListBlockOrderedType, func(t *testing.T) {
		block := testBlock()
		img, _ := domainblocks.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockOrderedType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on ordered", format.Struct(block))
		assert.Equal(t, img.Data.Alt, lb.Data.TextData.PlainText())
		assert.Equal(t, domainblocks.ListBlockOrderedType, lb.Data.Type)
	})

	t.Run(domainblocks.HeaderBlockType1, func(t *testing.T) {
		block := testBlock()
		img, _ := domainblocks.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType1))
		hb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_1", format.Struct(hb))
		assert.Equal(t, img.Data.Alt, hb.Data.TextData.PlainText())
		assert.Equal(t, uint(1), hb.Data.Level)
	})
	t.Run(domainblocks.HeaderBlockType2, func(t *testing.T) {
		block := testBlock()
		img, _ := domainblocks.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType2))
		hb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_2", format.Struct(block))
		assert.Equal(t, img.Data.Alt, hb.Data.TextData.PlainText())
		assert.Equal(t, uint(2), hb.Data.Level)
	})

	t.Run(domainblocks.HeaderBlockType3, func(t *testing.T) {
		block := testBlock()
		img, _ := domainblocks.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType3))
		hb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_3", format.Struct(block))
		assert.Equal(t, img.Data.Alt, hb.Data.TextData.PlainText())
		assert.Equal(t, uint(3), hb.Data.Level)
	})
	t.Run(domainblocks.FileBlockType, func(t *testing.T) {
		block := testBlock()
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.FileBlockType))
		fb, err := domainblocks.FromUnifiedToFileBlock(block)
		assert.NoError(t, err)
		log.Println("after change on file", format.Struct(block))
		assert.Equal(t, "", fb.Data.Src)
	})

	t.Run(domainblocks.ImgBlockType, func(t *testing.T) {
		block := testBlock()
		img, _ := domainblocks.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ImgBlockType))
		ib, err := domainblocks.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		log.Println("after change on img", format.Struct(block))
		assert.Equal(t, img.Data.Alt, ib.Data.Alt)
	})

	t.Run(domainblocks.LinkBlockType, func(t *testing.T) {
		block := testBlock()
		img, _ := domainblocks.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.LinkBlockType))
		linkB, err := domainblocks.FromUnifiedToLinkBlock(block)
		assert.NoError(t, err)
		log.Println("after change on link", format.Struct(block))
		assert.Equal(t, img.Data.Alt, linkB.Data.Text)
	})

	t.Run(domainblocks.QuoteBlockType, func(t *testing.T) {
		block := testBlock()
		img, _ := domainblocks.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.QuoteBlockType))
		qb, err := domainblocks.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		log.Println("after change on quote", format.Struct(qb))
		assert.Equal(t, img.Data.Alt, qb.Data.Text)
	})
}

func TestChangeTypeNil(t *testing.T) {
	t.Parallel()
	t.Run(domainblocks.TextBlockType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.TextBlockType))
		tb, err := domainblocks.FromUnifiedToTextBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", tb.Data.PlainText())
	})

	t.Run(domainblocks.ListBlockToDoType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockToDoType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, domainblocks.ListBlockToDoType, lb.Data.Type)
		assert.Equal(t, "", lb.Data.TextData.PlainText())
	})

	t.Run(domainblocks.CodeBlockType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.CodeBlockType))
		code, err := domainblocks.FromUnifiedToCodeBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", code.Data.Text)
	})

	t.Run(domainblocks.HeaderBlockType1, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType1))
		hb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), hb.Data.Level)
		assert.Equal(t, "", hb.Data.TextData.PlainText())
	})

	t.Run(domainblocks.ImgBlockType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ImgBlockType))
		ib, err := domainblocks.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", ib.Data.Alt)
		assert.Equal(t, "", ib.Data.Src)
	})

	t.Run(domainblocks.LinkBlockType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.LinkBlockType))
		linkB, err := domainblocks.FromUnifiedToLinkBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", linkB.Data.Text)
	})

	t.Run(domainblocks.QuoteBlockType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.QuoteBlockType))
		qb, err := domainblocks.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", qb.Data.Text)
	})
}

var testAlt = "A beautiful landscape"

func testBlock() *brzrpc.Block {
	ib := domainblocks.ImgBlock{
		Id:   "test_img",
		Type: domainblocks.ImgBlockType,
		Data: &domainblocks.ImgData{
			Src: "https://example.com/image.jpg",
			Alt: testAlt,
		},
	}
	unif, _ := ib.ToUnified()
	return unif
}

func testBlockNil() *brzrpc.Block {
	ib := domainblocks.ImgBlock{
		Id:   "test_img_nil",
		Type: domainblocks.ImgBlockType,
		Data: nil,
	}
	unif, _ := ib.ToUnified()
	return unif
}
