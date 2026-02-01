package headerblock

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

	lst := d.GetAsFirst(ctx, block)
	assert.Equal(t, lst, "text default text bold")
	lst = d.GetAsFirst(ctx, nil)
	assert.Equal(t, lst, "")
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
				lst, err := domain.FromUnifiedToHeaderBlock(block)
				if assert.NoError(t, err) {
					if assert.Equal(t, len(lst.Data.TextData.Text), 3) {
						assert.Equal(t, lst.Data.TextData.Text[0].Style, "test")
						log.Println("apply style", format.Struct(lst))
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
				lst, err := domain.FromUnifiedToHeaderBlock(block)
				if assert.NoError(t, err) {
					if assert.Equal(t, len(lst.Data.TextData.Text), 2) {
						assert.Equal(t, lst.Data.TextData.Text[0].String, "text default test")
						log.Println("insert text", format.Struct(lst))
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
				lst, err := domain.FromUnifiedToHeaderBlock(block)
				if assert.NoError(t, err) {
					if assert.Equal(t, len(lst.Data.TextData.Text), 2) {
						assert.Equal(t, lst.Data.TextData.Text[0].String, "default")
						log.Println("delete range", format.Struct(lst))
					}
				}
			}
		}
	})

	t.Run("change_level", func(t *testing.T) {
		block := testBlock()

		lstOld, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)

		data, err := d.Op(ctx, block, "change_level", map[string]any{
			"new_level": 10,
		})
		if assert.NoError(t, err) {
			s, err := structpb.NewStruct(data)
			if assert.NoError(t, err) {
				block.Data = s
				lst, err := domain.FromUnifiedToHeaderBlock(block)
				if assert.NoError(t, err) {
					assert.Equal(t, lstOld.Data.TextData, lst.Data.TextData)
					assert.Equal(t, uint(10), lst.Data.Level)
					log.Println("todo change value", format.Struct(block))
				}
			}
		}
		data, err = d.Op(ctx, block, "change_level", map[string]any{
			"new_level": -10,
		})
		if assert.NoError(t, err) {
			s, err := structpb.NewStruct(data)
			if assert.NoError(t, err) {
				block.Data = s
				lst, err := domain.FromUnifiedToHeaderBlock(block)
				if assert.NoError(t, err) {
					assert.Equal(t, lstOld.Data.TextData, lst.Data.TextData)
					assert.Equal(t, uint(0), lst.Data.Level)
					log.Println("todo change value", format.Struct(block))
				}
			}
		}
	})
}

func TestOpNil(t *testing.T) {
	t.Parallel()
	t.Run("apply style", func(t *testing.T) {
		block := testBlockNil()

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
		block := testBlockNil()

		data, err := d.Op(ctx, block, "insert_text", map[string]any{
			"pos":      12,
			"new_text": " test",
		})
		if assert.NoError(t, err) {
			assert.Nil(t, data)
		}
	})
	t.Run("delete_range", func(t *testing.T) {
		block := testBlockNil()

		data, err := d.Op(ctx, block, "delete_range", map[string]any{
			"start":    0,
			"end":      5,
			"new_text": " test",
		})
		if assert.NoError(t, err) {
			assert.Nil(t, data)
		}
	})

	t.Run("change_level", func(t *testing.T) {
		block := testBlockNil()

		data, err := d.Op(ctx, block, "change_level", map[string]any{
			"new_level": 10,
		})
		if assert.NoError(t, err) {
			assert.Nil(t, data)
		}
	})
}

func TestChangeTypeFromList(t *testing.T) {
	t.Parallel()
	t.Run(domain.TextBlockType, func(t *testing.T) {
		block := testBlock()
		list, _ := domain.FromUnifiedToHeaderBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.TextBlockType))
		tb, err := domain.FromUnifiedToTextBlock(block)
		assert.NoError(t, err)
		log.Println("after change on text", format.Struct(block))
		assert.Equal(t, list.Data.TextData, tb.Data)
	})

	t.Run(domain.ListBlockToDoType, func(t *testing.T) {
		block := testBlock()
		list, _ := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockToDoType))
		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockToDoType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered", format.Struct(block))
		assert.Equal(t, list.Data.TextData, lb.Data.TextData)
		assert.Equal(t, domain.ListBlockToDoType, lb.Data.Type)
	})

	t.Run(domain.ListBlockUnorderedType, func(t *testing.T) {
		block := testBlock()
		list, _ := domain.FromUnifiedToHeaderBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockUnorderedType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered", format.Struct(block))
		assert.Equal(t, list.Data.TextData, lb.Data.TextData)
		assert.Equal(t, domain.ListBlockUnorderedType, lb.Data.Type)
	})

	t.Run(domain.ListBlockOrderedType, func(t *testing.T) {
		block := testBlock()
		list, _ := domain.FromUnifiedToHeaderBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockOrderedType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on ordered", format.Struct(block))
		assert.Equal(t, list.Data.TextData, lb.Data.TextData)
		assert.Equal(t, domain.ListBlockOrderedType, lb.Data.Type)
	})

	t.Run(domain.CodeBlockType, func(t *testing.T) {
		block := testBlock()
		list, _ := domain.FromUnifiedToHeaderBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.CodeBlockType))
		cb, err := domain.FromUnifiedToCodeBlock(block)
		assert.NoError(t, err)
		log.Println("after change on code", format.Struct(block))
		assert.Equal(t, list.Data.TextData.PlainText(), cb.Data.Text)
	})

	t.Run(domain.HeaderBlockType1, func(t *testing.T) {
		block := testBlock()
		list, _ := domain.FromUnifiedToHeaderBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType1))
		hb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_1", format.Struct(block))
		assert.Equal(t, list.Data.TextData, hb.Data.TextData)
		assert.Equal(t, uint(1), hb.Data.Level)
	})

	t.Run(domain.HeaderBlockType2, func(t *testing.T) {
		block := testBlock()
		list, _ := domain.FromUnifiedToHeaderBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType2))
		hb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_2", format.Struct(block))
		assert.Equal(t, list.Data.TextData, hb.Data.TextData)
		assert.Equal(t, uint(2), hb.Data.Level)
	})

	t.Run(domain.HeaderBlockType3, func(t *testing.T) {
		block := testBlock()
		list, _ := domain.FromUnifiedToHeaderBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType3))
		hb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_3", format.Struct(block))
		assert.Equal(t, list.Data.TextData, hb.Data.TextData)
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
		list, _ := domain.FromUnifiedToHeaderBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ImgBlockType))
		ib, err := domain.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		log.Println("after change on img", format.Struct(block))
		assert.Equal(t, list.Data.TextData.PlainText(), ib.Data.Alt)
	})

	t.Run(domain.LinkBlockType, func(t *testing.T) {
		block := testBlock()
		list, _ := domain.FromUnifiedToHeaderBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.LinkBlockType))
		linkB, err := domain.FromUnifiedToLinkBlock(block)
		assert.NoError(t, err)
		log.Println("after change on link", format.Struct(block))
		assert.Equal(t, list.Data.TextData.PlainText(), linkB.Data.Text)
	})

	t.Run(domain.QuoteBlockType, func(t *testing.T) {
		block := testBlock()
		list, _ := domain.FromUnifiedToHeaderBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.QuoteBlockType))
		qb, err := domain.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		log.Println("after change on quote", format.Struct(block))
		assert.Equal(t, list.Data.TextData.PlainText(), qb.Data.Text)
	})
}

func TestChangeTypeFromListNil(t *testing.T) {
	t.Parallel()
	t.Run(domain.TextBlockType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domain.TextBlockType))
		tb, err := domain.FromUnifiedToTextBlock(block)
		assert.NoError(t, err)
		assert.Nil(t, tb.Data)
	})

	t.Run(domain.ListBlockToDoType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockToDoType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, domain.ListBlockToDoType, lb.Data.Type)
		assert.Nil(t, lb.Data.TextData)
	})
	t.Run(domain.ListBlockUnorderedType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockUnorderedType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, domain.ListBlockUnorderedType, lb.Data.Type)
		assert.Nil(t, lb.Data.TextData)
	})
	t.Run(domain.ListBlockOrderedType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockOrderedType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, domain.ListBlockOrderedType, lb.Data.Type)
		assert.Nil(t, lb.Data.TextData)
	})

	t.Run(domain.CodeBlockType, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domain.CodeBlockType))
		cb, err := domain.FromUnifiedToCodeBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", cb.Data.Text)
	})

	t.Run(domain.HeaderBlockType1, func(t *testing.T) {
		block := testBlockNil()
		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType1))
		hb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), hb.Data.Level)
		assert.Nil(t, hb.Data.TextData)
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

var (
	d   = Driver{}
	ctx = context.Background()
)

func testBlock() *brzrpc.Block {
	lst := domain.HeaderBlock{
		Id:        "test",
		Type:      "list",
		NoteId:    "test",
		CreatedAt: 0,
		UpdatedAt: 0,
		IsUsed:    false,
		Data: &domain.HeaderData{
			TextData: &text.Data{
				Text: []text.Part{
					{
						Style:  "default",
						String: "text default",
					},
					{
						Style:  "bold",
						String: " text bold",
					},
				},
			},
			Level: 3,
		},
	}

	lstUnif, err := lst.ToUnified()
	if err != nil {
		return nil
	}

	return lstUnif
}
func testBlockNil() *brzrpc.Block {
	lst := domain.HeaderBlock{
		Id:        "test",
		Type:      "list",
		NoteId:    "test",
		CreatedAt: 0,
		UpdatedAt: 0,
		IsUsed:    false,
		Data:      nil,
	}

	lstUnif, err := lst.ToUnified()
	if err != nil {
		return nil
	}

	return lstUnif
}
