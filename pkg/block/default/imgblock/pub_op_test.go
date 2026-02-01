package imgblock

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/domain"
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
			ib, _ := domain.FromUnifiedToImgBlock(block)
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
			ib, _ := domain.FromUnifiedToImgBlock(block)
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
	t.Run(domain.TextBlockType, func(t *testing.T) {
		block := testBlock()
		img, _ := domain.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.TextBlockType))
		tb, err := domain.FromUnifiedToTextBlock(block)
		assert.NoError(t, err)
		log.Println("after change on text", format.Struct(tb))
		assert.Equal(t, img.Data.Alt, tb.Data.PlainText())
	})

	t.Run(domain.ListBlockUnorderedType, func(t *testing.T) {
		block := testBlock()
		img, _ := domain.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockUnorderedType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered img", format.Struct(lb))
		assert.Equal(t, img.Data.Alt, lb.Data.TextData.PlainText())
		assert.Equal(t, domain.ListBlockUnorderedType, lb.Data.Type)
	})
	t.Run(domain.ListBlockToDoType, func(t *testing.T) {
		block := testBlock()
		img, _ := domain.FromUnifiedToImgBlock(block)
		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockToDoType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered", format.Struct(block))
		assert.Equal(t, img.Data.Alt, lb.Data.TextData.PlainText())
		assert.Equal(t, domain.ListBlockToDoType, lb.Data.Type)
	})
	t.Run(domain.ListBlockOrderedType, func(t *testing.T) {
		block := testBlock()
		img, _ := domain.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockOrderedType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on ordered", format.Struct(block))
		assert.Equal(t, img.Data.Alt, lb.Data.TextData.PlainText())
		assert.Equal(t, domain.ListBlockOrderedType, lb.Data.Type)
	})

	t.Run(domain.HeaderBlockType1, func(t *testing.T) {
		block := testBlock()
		img, _ := domain.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType1))
		hb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_1", format.Struct(hb))
		assert.Equal(t, img.Data.Alt, hb.Data.TextData.PlainText())
		assert.Equal(t, uint(1), hb.Data.Level)
	})
	t.Run(domain.HeaderBlockType2, func(t *testing.T) {
		block := testBlock()
		img, _ := domain.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType2))
		hb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_2", format.Struct(block))
		assert.Equal(t, img.Data.Alt, hb.Data.TextData.PlainText())
		assert.Equal(t, uint(2), hb.Data.Level)
	})

	t.Run(domain.HeaderBlockType3, func(t *testing.T) {
		block := testBlock()
		img, _ := domain.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType3))
		hb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_3", format.Struct(block))
		assert.Equal(t, img.Data.Alt, hb.Data.TextData.PlainText())
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
		img, _ := domain.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ImgBlockType))
		ib, err := domain.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		log.Println("after change on img", format.Struct(block))
		assert.Equal(t, img.Data.Alt, ib.Data.Alt)
	})

	t.Run(domain.LinkBlockType, func(t *testing.T) {
		block := testBlock()
		img, _ := domain.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.LinkBlockType))
		linkB, err := domain.FromUnifiedToLinkBlock(block)
		assert.NoError(t, err)
		log.Println("after change on link", format.Struct(block))
		assert.Equal(t, img.Data.Alt, linkB.Data.Text)
	})

	t.Run(domain.QuoteBlockType, func(t *testing.T) {
		block := testBlock()
		img, _ := domain.FromUnifiedToImgBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.QuoteBlockType))
		qb, err := domain.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		log.Println("after change on quote", format.Struct(qb))
		assert.Equal(t, img.Data.Alt, qb.Data.Text)
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

var testAlt = "A beautiful landscape"

func testBlock() *brzrpc.Block {
	ib := domain.ImgBlock{
		Id:   "test_img",
		Type: domain.ImgBlockType,
		Data: &domain.ImgData{
			Src: "https://example.com/image.jpg",
			Alt: testAlt,
		},
	}
	unif, _ := ib.ToUnified()
	return unif
}

func testBlockNil() *brzrpc.Block {
	ib := domain.ImgBlock{
		Id:   "test_img_nil",
		Type: domain.ImgBlockType,
		Data: nil,
	}
	unif, _ := ib.ToUnified()
	return unif
}
