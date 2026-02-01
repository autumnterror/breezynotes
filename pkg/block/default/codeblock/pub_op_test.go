package codeblock

import (
	"context"
	"fmt"
	"github.com/alecthomas/chroma/v2/lexers"
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/pkg/domain"
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
	assert.Equal(t, txt, textTest)
	txt = d.GetAsFirst(ctx, nil)
	assert.Equal(t, txt, "")
}

func TestOp(t *testing.T) {
	t.Parallel()
	t.Run("change_text", func(t *testing.T) {
		block := testBlock()
		startBlock, _ := domain.FromUnifiedToCodeBlock(block)
		data, err := d.Op(ctx, block, "change_text", map[string]any{
			"new_text": "print('hello world')",
		})
		if assert.NoError(t, err) {
			s, err := structpb.NewStruct(data)
			if assert.NoError(t, err) {
				block.Data = s
				cb, err := domain.FromUnifiedToCodeBlock(block)
				log.Blue(format.Struct(cb))
				if assert.NoError(t, err) {
					assert.Equal(t, "print('hello world')", cb.Data.Text)
					assert.Equal(t, startBlock.Data.Lang, cb.Data.Lang)
					log.Println("change text", format.Struct(cb))
				}
			}
		}
	})

	t.Run("analyse_lang", func(t *testing.T) {
		block := testBlock()
		data, err := d.Op(ctx, block, "analyse_lang", nil)
		if assert.NoError(t, err) {
			s, err := structpb.NewStruct(data)
			if assert.NoError(t, err) {
				block.Data = s
				cb, err := domain.FromUnifiedToCodeBlock(block)
				if assert.NoError(t, err) {
					assert.Equal(t, textTest, cb.Data.Text)
					assert.Equal(t, "Go", cb.Data.Lang)
					log.Println("analyse lang", format.Struct(cb))
				}
			}
		}
	})
}

func TestOpNil(t *testing.T) {
	t.Parallel()
	t.Run("change_text", func(t *testing.T) {
		block := testBlockNil()

		data, err := d.Op(ctx, block, "change_text", map[string]any{
			"new_text": "test",
		})
		if assert.NoError(t, err) {
			assert.Nil(t, data)
		}
	})

	t.Run("analyse_lang", func(t *testing.T) {
		block := testBlockNil()

		data, err := d.Op(ctx, block, "analyse_lang", nil)
		if assert.NoError(t, err) {
			assert.Nil(t, data)
		}
	})
}

func TestChangeType(t *testing.T) {
	t.Parallel()
	t.Run(domain.TextBlockType, func(t *testing.T) {
		block := testBlock()
		cb, _ := domain.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.TextBlockType))
		tb, err := domain.FromUnifiedToTextBlock(block)
		assert.NoError(t, err)
		log.Println("after change on text", format.Struct(tb))
		assert.Equal(t, cb.Data.Text, tb.Data.PlainText())
	})

	t.Run(domain.ListBlockUnorderedType, func(t *testing.T) {
		block := testBlock()
		cb, _ := domain.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockUnorderedType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered code", format.Struct(lb))
		assert.Equal(t, cb.Data.Text, lb.Data.TextData.PlainText())
		assert.Equal(t, domain.ListBlockUnorderedType, lb.Data.Type)
	})
	t.Run(domain.ListBlockToDoType, func(t *testing.T) {
		block := testBlock()
		code, _ := domain.FromUnifiedToCodeBlock(block)
		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockToDoType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered", format.Struct(block))
		assert.Equal(t, code.Data.Text, lb.Data.TextData.PlainText())
		assert.Equal(t, domain.ListBlockToDoType, lb.Data.Type)
	})
	t.Run(domain.ListBlockOrderedType, func(t *testing.T) {
		block := testBlock()
		code, _ := domain.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ListBlockOrderedType))
		lb, err := domain.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on ordered", format.Struct(block))
		assert.Equal(t, code.Data.Text, lb.Data.TextData.PlainText())
		assert.Equal(t, domain.ListBlockOrderedType, lb.Data.Type)
	})

	t.Run(domain.HeaderBlockType1, func(t *testing.T) {
		block := testBlock()
		cb, _ := domain.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType1))
		hb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_1", format.Struct(hb))
		assert.Equal(t, cb.Data.Text, hb.Data.TextData.PlainText())
		assert.Equal(t, uint(1), hb.Data.Level)
	})
	t.Run(domain.HeaderBlockType2, func(t *testing.T) {
		block := testBlock()
		code, _ := domain.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType2))
		hb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_2", format.Struct(block))
		assert.Equal(t, code.Data.Text, hb.Data.TextData.PlainText())
		assert.Equal(t, uint(2), hb.Data.Level)
	})

	t.Run(domain.HeaderBlockType3, func(t *testing.T) {
		block := testBlock()
		code, _ := domain.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.HeaderBlockType3))
		hb, err := domain.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_3", format.Struct(block))
		assert.Equal(t, code.Data.Text, hb.Data.TextData.PlainText())
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
		code, _ := domain.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.ImgBlockType))
		ib, err := domain.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		log.Println("after change on img", format.Struct(block))
		assert.Equal(t, code.Data.Text, ib.Data.Alt)
	})

	t.Run(domain.LinkBlockType, func(t *testing.T) {
		block := testBlock()
		code, _ := domain.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.LinkBlockType))
		linkB, err := domain.FromUnifiedToLinkBlock(block)
		assert.NoError(t, err)
		log.Println("after change on link", format.Struct(block))
		assert.Equal(t, code.Data.Text, linkB.Data.Text)
	})

	t.Run(domain.QuoteBlockType, func(t *testing.T) {
		block := testBlock()
		cb, _ := domain.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domain.QuoteBlockType))
		qb, err := domain.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		log.Println("after change on quote", format.Struct(qb))
		assert.Equal(t, cb.Data.Text, qb.Data.Text)
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

var (
	d   = Driver{}
	ctx = context.Background()
)

var textTest = "\tpackage main\n\t\timport \"fmt\"\n\t\tfunc main() { fmt.Println(\"Hello\") }"

func testBlock() *brzrpc.Block {
	cb := domain.CodeBlock{
		Id:     "test_code",
		Type:   "code",
		NoteId: "test_note",
		Data: &domain.CodeData{
			Text: textTest,
			Lang: "Go",
		},
	}

	cbUnif, err := cb.ToUnified()
	if err != nil {
		return nil
	}

	return cbUnif
}

func testBlockNil() *brzrpc.Block {
	cb := domain.CodeBlock{
		Id:     "test_code_nil",
		Type:   "code",
		NoteId: "test_note",
		Data:   nil,
	}

	cbUnif, err := cb.ToUnified()
	if err != nil {
		return nil
	}

	return cbUnif
}

func TestLang(t *testing.T) {
	// Пример 1: HTML
	htmlSnippet := `
		function addNumbers(a, b) {
			return a + b;
		}
		
		// Call the function and store the result
		let result = addNumbers(5, 10);
		
		// Output the result (15)
		console.log(result);
`
	analyzeAndPrint(htmlSnippet)

	// Пример 2: Go
	goSnippet := `
		package main
		import "fmt"
		func main() { fmt.Println("Hello") }
	`
	analyzeAndPrint(goSnippet)

	// Пример 3: Python
	pythonSnippet := `
		#include <stdio.h>
		int main() {
		   // printf() displays the string inside quotation
		   printf("Hello, World!");
		   return 0;
		}
	`
	analyzeAndPrint(pythonSnippet)
}

func analyzeAndPrint(code string) {
	lexer := lexers.Analyse(code)
	if lexer != nil {
		lang := lexer.Config().Name
		fmt.Printf("Для фрагмента:\n%s\n... определён язык: %s\n", code, lang)
		fmt.Println("---")
	} else {
		fmt.Printf("Для фрагмента:\n%s\n... не удалось определить язык.\n", code)
		fmt.Println("---")
	}
}
