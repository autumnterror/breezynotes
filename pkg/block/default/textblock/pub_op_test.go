package textblock

import (
	"context"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/domain"
	"github.com/autumnterror/breezynotes/pkg/text"
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
				txt, err := domain.FromUnifiedToTextBlock(block)
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
				txt, err := domain.FromUnifiedToTextBlock(block)
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
				txt, err := domain.FromUnifiedToTextBlock(block)
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
	t.Run(domain.TextBlockType, func(t *testing.T) {
		block := testBlock()

		assert.NoError(t, d.ChangeType(ctx, block, domain.TextBlockType))
	})
	t.Run(domain.ListBlockToDoType, func(t *testing.T) {
		block := testBlock()

		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockToDoType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on todo", format.Struct(block))
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, domain.ListBlockToDoType, lb.Data.Type)
	})
	t.Run(domain.ListBlockUnorderedType, func(t *testing.T) {
		block := testBlock()

		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockUnorderedType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered", format.Struct(block))
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, domain.ListBlockUnorderedType, lb.Data.Type)
	})
	t.Run(domain.ListBlockOrderedType, func(t *testing.T) {
		block := testBlock()

		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockOrderedType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on ordered", format.Struct(block))
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, domain.ListBlockOrderedType, lb.Data.Type)
	})
	t.Run(domain.CodeBlockType, func(t *testing.T) {
		block := testBlock()

		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.CodeBlockType))
		lb, err := domain.FromUnifiedToCodeBlock(block)
		assert.NoError(t, err)
		log.Println("after change on code", format.Struct(block))
		assert.Equal(t, txt.Data.PlainText(), lb.Data.Text)
	})
	t.Run(domain.HeaderBlockType1, func(t *testing.T) {
		block := testBlock()

		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType1))
		lb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_1", format.Struct(block))
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, uint(1), lb.Data.Level)
	})
	t.Run(domain.HeaderBlockType2, func(t *testing.T) {
		block := testBlock()

		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType2))
		lb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_2", format.Struct(block))
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, uint(2), lb.Data.Level)
	})
	t.Run(domain.HeaderBlockType3, func(t *testing.T) {
		block := testBlock()

		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType3))
		lb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_3", format.Struct(block))
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, uint(3), lb.Data.Level)
	})
	t.Run(domain.FileBlockType, func(t *testing.T) {
		block := testBlock()

		assert.NoError(t, d.ChangeType(ctx, block, domain.FileBlockType))
		lb, err := domain.FromUnifiedToFileBlock(block)
		assert.NoError(t, err)
		log.Println("after change on file", format.Struct(block))
		assert.Equal(t, "", lb.Data.Src)
	})
	t.Run(domain.ImgBlockType, func(t *testing.T) {
		block := testBlock()

		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ImgBlockType))
		lb, err := domain.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		log.Println("after change on img", format.Struct(block))
		assert.Equal(t, txt.Data.PlainText(), lb.Data.Alt)
	})
	t.Run(domain.LinkBlockType, func(t *testing.T) {
		block := testBlock()

		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.LinkBlockType))
		lb, err := domain.FromUnifiedToLinkBlock(block)
		assert.NoError(t, err)
		log.Println("after change on link", format.Struct(block))
		assert.Equal(t, txt.Data.PlainText(), lb.Data.Text)
	})
	t.Run(domain.QuoteBlockType, func(t *testing.T) {
		block := testBlock()

		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.QuoteBlockType))
		lb, err := domain.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		log.Println("after change on quote", format.Struct(block))
		assert.Equal(t, txt.Data.PlainText(), lb.Data.Text)

	})
}

func TestChangeTypeNil(t *testing.T) {
	t.Parallel()
	t.Run(domain.TextBlockType, func(t *testing.T) {
		block := testNilBlock()
		assert.NoError(t, d.ChangeType(ctx, block, domain.TextBlockType))
	})
	t.Run(domain.ListBlockToDoType, func(t *testing.T) {
		block := testNilBlock()
		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockToDoType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, domain.ListBlockToDoType, lb.Data.Type)
	})
	t.Run(domain.ListBlockUnorderedType, func(t *testing.T) {
		block := testNilBlock()
		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockUnorderedType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, domain.ListBlockUnorderedType, lb.Data.Type)
	})
	t.Run(domain.ListBlockOrderedType, func(t *testing.T) {
		block := testNilBlock()
		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockOrderedType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, domain.ListBlockOrderedType, lb.Data.Type)
	})
	t.Run(domain.CodeBlockType, func(t *testing.T) {
		block := testNilBlock()
		assert.NoError(t, d.ChangeType(ctx, block, domain.CodeBlockType))
		lb, err := domain.FromUnifiedToCodeBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", lb.Data.Text)
	})
	t.Run(domain.HeaderBlockType1, func(t *testing.T) {
		block := testNilBlock()

		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType1))
		lb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, uint(1), lb.Data.Level)
	})
	t.Run(domain.HeaderBlockType2, func(t *testing.T) {
		block := testNilBlock()

		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType2))
		lb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, uint(2), lb.Data.Level)
	})
	t.Run(domain.HeaderBlockType3, func(t *testing.T) {
		block := testNilBlock()

		txt, _ := domain.FromUnifiedToTextBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType3))
		lb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, txt.Data, lb.Data.TextData)
		assert.Equal(t, uint(3), lb.Data.Level)
	})
	t.Run(domain.FileBlockType, func(t *testing.T) {
		block := testNilBlock()
		assert.NoError(t, d.ChangeType(ctx, block, domain.FileBlockType))
		lb, err := domain.FromUnifiedToFileBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", lb.Data.Src)
	})
	t.Run(domain.ImgBlockType, func(t *testing.T) {
		block := testNilBlock()
		assert.NoError(t, d.ChangeType(ctx, block, domain.ImgBlockType))
		lb, err := domain.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", lb.Data.Alt)
	})
	t.Run(domain.LinkBlockType, func(t *testing.T) {
		block := testNilBlock()
		assert.NoError(t, d.ChangeType(ctx, block, domain.LinkBlockType))
		lb, err := domain.FromUnifiedToLinkBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", lb.Data.Text)
	})
	t.Run(domain.QuoteBlockType, func(t *testing.T) {
		block := testNilBlock()
		assert.NoError(t, d.ChangeType(ctx, block, domain.QuoteBlockType))
		lb, err := domain.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", lb.Data.Text)
	})
}

var (
	d   = Driver{}
	ctx = context.Background()
)

func testBlock() *brzrpc.Block {
	txt := domain.TextBlock{
		Id:        "test",
		Type:      domain.TextBlockType,
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
	txt := domain.TextBlock{
		Id:        "test",
		Type:      domain.TextBlockType,
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
