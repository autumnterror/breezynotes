package textblock

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain/domainblocks"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/text"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
	"testing"
)

func TestGetAsFirst(t *testing.T) {
	t.Parallel()
	block := testBlock()

	txt := d.GetAsFirst(ctx, block)
	assert.Equal(t, txt, "text default text bold")
	txt = d.GetAsFirst(ctx, nil)
	assert.Equal(t, txt, "")
}

func TestOp(t *testing.T) {
	t.Parallel()
	t.Run("apply style", func(t *testing.T) {
		block := testBlock()

		data, err := d.Op(ctx, block, "apply_style", map[string]any{
			"start": 0,
			"end":   2,
			"style": "test",
		})
		if assert.NoError(t, err) {
			s, err := structpb.NewStruct(data)
			if assert.NoError(t, err) {
				block.Data = s
				txt, err := domainblocks.FromUnifiedToTextBlock(block)
				if assert.NoError(t, err) {
					if assert.Equal(t, len(txt.Data.Text), 3) {
						assert.Equal(t, txt.Data.Text[0].Style, "test")
						log.Println("apply style", format.Struct(txt))
					}
				}
			}
		}
	})
	t.Run("insert_text", func(t *testing.T) {
		block := testBlock()

		data, err := d.Op(ctx, block, "insert_text", map[string]any{
			"pos":      12,
			"new_text": " test",
		})
		if assert.NoError(t, err) {
			s, err := structpb.NewStruct(data)
			if assert.NoError(t, err) {
				block.Data = s
				txt, err := domainblocks.FromUnifiedToTextBlock(block)
				if assert.NoError(t, err) {
					if assert.Equal(t, len(txt.Data.Text), 2) {
						assert.Equal(t, txt.Data.Text[0].String, "text default test")
						log.Println("insert text", format.Struct(txt))
					}
				}
			}
		}
	})
	t.Run("delete_range", func(t *testing.T) {
		block := testBlock()

		data, err := d.Op(ctx, block, "delete_range", map[string]any{
			"start":    0,
			"end":      5,
			"new_text": " test",
		})
		if assert.NoError(t, err) {
			s, err := structpb.NewStruct(data)
			if assert.NoError(t, err) {
				block.Data = s
				txt, err := domainblocks.FromUnifiedToTextBlock(block)
				if assert.NoError(t, err) {
					if assert.Equal(t, len(txt.Data.Text), 2) {
						assert.Equal(t, txt.Data.Text[0].String, "default")
						log.Println("delete range", format.Struct(txt))
					}
				}
			}
		}
	})
}

func TestOpBad(t *testing.T) {
	t.Parallel()
	t.Run("apply style", func(t *testing.T) {
		block := testNilBlock()
		data, err := d.Op(ctx, block, "apply_style", map[string]any{
			"start": 0,
			"end":   2,
			"style": "test",
		})
		if assert.NoError(t, err) {
			assert.Nil(t, data)
		}
	})
	t.Run("insert_text", func(t *testing.T) {
		block := testNilBlock()

		data, err := d.Op(ctx, block, "insert_text", map[string]any{
			"pos":      12,
			"new_text": " test",
		})
		if assert.NoError(t, err) {
			assert.Nil(t, data)
		}
	})
	t.Run("delete_range", func(t *testing.T) {
		block := testNilBlock()

		data, err := d.Op(ctx, block, "delete_range", map[string]any{
			"start":    0,
			"end":      5,
			"new_text": " test",
		})
		if assert.NoError(t, err) {
			assert.Nil(t, data)
		}
	})
}

func TestChangeType(t *testing.T) {
	t.Parallel()
	t.Run(domainblocks.TextBlockType, func(t *testing.T) {
		block := testBlock()

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.TextBlockType))
	})
	t.Run(domainblocks.ListBlockToDoType, func(t *testing.T) {
		block := testBlock()

		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockToDoType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on todo", format.Struct(block))
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, domainblocks.ListBlockToDoType, lb.Data.Type)
	})
	t.Run(domainblocks.ListBlockUnorderedType, func(t *testing.T) {
		block := testBlock()

		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockUnorderedType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered", format.Struct(block))
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, domainblocks.ListBlockUnorderedType, lb.Data.Type)
	})
	t.Run(domainblocks.ListBlockOrderedType, func(t *testing.T) {
		block := testBlock()

		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockOrderedType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on ordered", format.Struct(block))
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, domainblocks.ListBlockOrderedType, lb.Data.Type)
	})
	t.Run(domainblocks.CodeBlockType, func(t *testing.T) {
		block := testBlock()

		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.CodeBlockType))
		lb, err := domainblocks.FromUnifiedToCodeBlock(block)
		assert.NoError(t, err)
		log.Println("after change on code", format.Struct(block))
		assert.Equal(t, txt.Data.PlainText(), lb.Data.Text)
	})
	t.Run(domainblocks.HeaderBlockType1, func(t *testing.T) {
		block := testBlock()

		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType1))
		lb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_1", format.Struct(block))
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, uint(1), lb.Data.Level)
	})
	t.Run(domainblocks.HeaderBlockType2, func(t *testing.T) {
		block := testBlock()

		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType2))
		lb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_2", format.Struct(block))
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, uint(2), lb.Data.Level)
	})
	t.Run(domainblocks.HeaderBlockType3, func(t *testing.T) {
		block := testBlock()

		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType3))
		lb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_3", format.Struct(block))
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, uint(3), lb.Data.Level)
	})
	t.Run(domainblocks.FileBlockType, func(t *testing.T) {
		block := testBlock()

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.FileBlockType))
		lb, err := domainblocks.FromUnifiedToFileBlock(block)
		assert.NoError(t, err)
		log.Println("after change on file", format.Struct(block))
		assert.Equal(t, "", lb.Data.Src)
	})
	t.Run(domainblocks.ImgBlockType, func(t *testing.T) {
		block := testBlock()

		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ImgBlockType))
		lb, err := domainblocks.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		log.Println("after change on img", format.Struct(block))
		assert.Equal(t, txt.Data.PlainText(), lb.Data.Alt)
	})
	t.Run(domainblocks.LinkBlockType, func(t *testing.T) {
		block := testBlock()

		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.LinkBlockType))
		lb, err := domainblocks.FromUnifiedToLinkBlock(block)
		assert.NoError(t, err)
		log.Println("after change on link", format.Struct(block))
		assert.Equal(t, txt.Data.PlainText(), lb.Data.Text)
	})
	t.Run(domainblocks.QuoteBlockType, func(t *testing.T) {
		block := testBlock()

		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.QuoteBlockType))
		lb, err := domainblocks.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		log.Println("after change on quote", format.Struct(block))
		assert.Equal(t, txt.Data.PlainText(), lb.Data.Text)

	})
}

func TestChangeTypeNil(t *testing.T) {
	t.Parallel()
	t.Run(domainblocks.TextBlockType, func(t *testing.T) {
		block := testNilBlock()
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.TextBlockType))
	})
	t.Run(domainblocks.ListBlockToDoType, func(t *testing.T) {
		block := testNilBlock()
		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockToDoType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, domainblocks.ListBlockToDoType, lb.Data.Type)
	})
	t.Run(domainblocks.ListBlockUnorderedType, func(t *testing.T) {
		block := testNilBlock()
		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockUnorderedType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, domainblocks.ListBlockUnorderedType, lb.Data.Type)
	})
	t.Run(domainblocks.ListBlockOrderedType, func(t *testing.T) {
		block := testNilBlock()
		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockOrderedType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, domainblocks.ListBlockOrderedType, lb.Data.Type)
	})
	t.Run(domainblocks.CodeBlockType, func(t *testing.T) {
		block := testNilBlock()
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.CodeBlockType))
		lb, err := domainblocks.FromUnifiedToCodeBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", lb.Data.Text)
	})
	t.Run(domainblocks.HeaderBlockType1, func(t *testing.T) {
		block := testNilBlock()

		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType1))
		lb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, uint(1), lb.Data.Level)
	})
	t.Run(domainblocks.HeaderBlockType2, func(t *testing.T) {
		block := testNilBlock()

		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType2))
		lb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, uint(2), lb.Data.Level)
	})
	t.Run(domainblocks.HeaderBlockType3, func(t *testing.T) {
		block := testNilBlock()

		txt, _ := domainblocks.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType3))
		lb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, uint(3), lb.Data.Level)
	})
	t.Run(domainblocks.FileBlockType, func(t *testing.T) {
		block := testNilBlock()
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.FileBlockType))
		lb, err := domainblocks.FromUnifiedToFileBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", lb.Data.Src)
	})
	t.Run(domainblocks.ImgBlockType, func(t *testing.T) {
		block := testNilBlock()
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ImgBlockType))
		lb, err := domainblocks.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", lb.Data.Alt)
	})
	t.Run(domainblocks.LinkBlockType, func(t *testing.T) {
		block := testNilBlock()
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.LinkBlockType))
		lb, err := domainblocks.FromUnifiedToLinkBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", lb.Data.Text)
	})
	t.Run(domainblocks.QuoteBlockType, func(t *testing.T) {
		block := testNilBlock()
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.QuoteBlockType))
		lb, err := domainblocks.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", lb.Data.Text)
	})
}

var (
	d   = Driver{}
	ctx = context.Background()
)

func testBlock() *brzrpc.Block {
	txt := domainblocks.TextBlock{
		Id:        "test",
		Type:      domainblocks.TextBlockType,
		NoteId:    "test",
		CreatedAt: 0,
		UpdatedAt: 0,
		IsUsed:    false,
		Data: &text.Data{Text: []text.Part{
			{
				Style:  "default",
				String: "text default",
			},
			{
				Style:  "bold",
				String: " text bold",
			},
		}},
	}

	txtUnif, err := txt.ToUnified()
	if err != nil {
		return nil
	}

	return txtUnif
}
func testNilBlock() *brzrpc.Block {
	txt := domainblocks.TextBlock{
		Id:        "test",
		Type:      domainblocks.TextBlockType,
		NoteId:    "test",
		CreatedAt: 0,
		UpdatedAt: 0,
		IsUsed:    false,
		Data:      nil,
	}

	txtUnif, err := txt.ToUnified()
	if err != nil {
		return nil
	}

	return txtUnif
}
