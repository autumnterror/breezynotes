package linkblock

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2/domainblocks"
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
	assert.Equal(t, textTest, d.GetAsFirst(ctx, block))

	blockNil := testBlockNil()
	assert.Equal(t, "", d.GetAsFirst(ctx, blockNil))
}

func TestOp(t *testing.T) {
	t.Parallel()
	newText := "Visit our official website"
	t.Run("change_text", func(t *testing.T) {
		block := testBlock()
		data, err := d.Op(ctx, block, "change_text", map[string]any{"new_text": newText})
		if assert.NoError(t, err) {
			s, _ := structpb.NewStruct(data)
			block.Data = s
			lb, _ := domainblocks.FromUnifiedToLinkBlock(block)
			assert.Equal(t, newText, lb.Data.Text)
			assert.Equal(t, "uwu.com", lb.Data.Url)
		}
	})
	newUrl := "owo.com"
	t.Run("change_uri", func(t *testing.T) {
		block := testBlock()
		data, err := d.Op(ctx, block, "change_url", map[string]any{"new_url": newUrl})
		if assert.NoError(t, err) {
			s, _ := structpb.NewStruct(data)
			block.Data = s
			lb, _ := domainblocks.FromUnifiedToLinkBlock(block)
			assert.Equal(t, textTest, lb.Data.Text)
			assert.Equal(t, newUrl, lb.Data.Url)
		}
	})
}

func TestOpNilData(t *testing.T) {
	t.Parallel()
	newText := "Visit our official website"
	t.Run("change_text", func(t *testing.T) {
		block := testBlockNil()
		data, err := d.Op(ctx, block, "change_text", map[string]any{"new_text": newText})
		if assert.NoError(t, err) {
			assert.Nil(t, data)
		}
	})
	newUrl := "owo.com"
	t.Run("change_uri", func(t *testing.T) {
		block := testBlockNil()
		data, err := d.Op(ctx, block, "change_url", map[string]any{"new_url": newUrl})
		if assert.NoError(t, err) {
			assert.Nil(t, data)
		}
	})
}

func TestChangeType(t *testing.T) {
	t.Parallel()
	t.Run(domainblocks.TextBlockType, func(t *testing.T) {
		block := testBlock()
		link, _ := domainblocks.FromUnifiedToLinkBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.TextBlockType))
		tb, err := domainblocks.FromUnifiedToTextBlock(block)
		assert.NoError(t, err)
		log.Println("after change on text", format.Struct(tb))
		assert.Equal(t, link.Data.Text, tb.Data.PlainText())
	})

	t.Run(domainblocks.ListBlockUnorderedType, func(t *testing.T) {
		block := testBlock()
		link, _ := domainblocks.FromUnifiedToLinkBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockUnorderedType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered link", format.Struct(lb))
		assert.Equal(t, link.Data.Text, lb.Data.TextData.PlainText())
		assert.Equal(t, domainblocks.ListBlockUnorderedType, lb.Data.Type)
	})
	t.Run(domainblocks.ListBlockToDoType, func(t *testing.T) {
		block := testBlock()
		link, _ := domainblocks.FromUnifiedToLinkBlock(block)
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockToDoType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered", format.Struct(block))
		assert.Equal(t, link.Data.Text, lb.Data.TextData.PlainText())
		assert.Equal(t, domainblocks.ListBlockToDoType, lb.Data.Type)
	})
	t.Run(domainblocks.ListBlockOrderedType, func(t *testing.T) {
		block := testBlock()
		link, _ := domainblocks.FromUnifiedToLinkBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockOrderedType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on ordered", format.Struct(block))
		assert.Equal(t, link.Data.Text, lb.Data.TextData.PlainText())
		assert.Equal(t, domainblocks.ListBlockOrderedType, lb.Data.Type)
	})

	t.Run(domainblocks.HeaderBlockType1, func(t *testing.T) {
		block := testBlock()
		link, _ := domainblocks.FromUnifiedToLinkBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType1))
		hb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_1", format.Struct(hb))
		assert.Equal(t, link.Data.Text, hb.Data.TextData.PlainText())
		assert.Equal(t, uint(1), hb.Data.Level)
	})
	t.Run(domainblocks.HeaderBlockType2, func(t *testing.T) {
		block := testBlock()
		link, _ := domainblocks.FromUnifiedToLinkBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType2))
		hb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_2", format.Struct(block))
		assert.Equal(t, link.Data.Text, hb.Data.TextData.PlainText())
		assert.Equal(t, uint(2), hb.Data.Level)
	})

	t.Run(domainblocks.HeaderBlockType3, func(t *testing.T) {
		block := testBlock()
		link, _ := domainblocks.FromUnifiedToLinkBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType3))
		hb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_3", format.Struct(block))
		assert.Equal(t, link.Data.Text, hb.Data.TextData.PlainText())
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
		link, _ := domainblocks.FromUnifiedToLinkBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ImgBlockType))
		ib, err := domainblocks.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		log.Println("after change on link", format.Struct(block))
		assert.Equal(t, link.Data.Text, ib.Data.Alt)
	})

	t.Run(domainblocks.LinkBlockType, func(t *testing.T) {
		block := testBlock()
		link, _ := domainblocks.FromUnifiedToLinkBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.LinkBlockType))
		linkB, err := domainblocks.FromUnifiedToLinkBlock(block)
		assert.NoError(t, err)
		log.Println("after change on link", format.Struct(block))
		assert.Equal(t, link.Data.Text, linkB.Data.Text)
	})

	t.Run(domainblocks.QuoteBlockType, func(t *testing.T) {
		block := testBlock()
		link, _ := domainblocks.FromUnifiedToLinkBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.QuoteBlockType))
		qb, err := domainblocks.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		log.Println("after change on quote", format.Struct(qb))
		assert.Equal(t, link.Data.Text, qb.Data.Text)
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

var textTest = "this is uwutube"

func testBlock() *brzrpc.Block {
	lb := domainblocks.LinkBlock{
		Id:   "test_link",
		Type: domainblocks.LinkBlockType,
		Data: &domainblocks.LinkData{
			Text: textTest,
			Url:  "uwu.com",
		},
	}
	unif, _ := lb.ToUnified()
	return unif
}

func testBlockNil() *brzrpc.Block {
	lb := domainblocks.LinkBlock{
		Id:   "test_link_nil",
		Type: domainblocks.LinkBlockType,
		Data: nil,
	}
	unif, _ := lb.ToUnified()
	return unif
}
