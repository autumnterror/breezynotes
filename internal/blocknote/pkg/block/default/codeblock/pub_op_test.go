package codeblock

import (
	"context"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2/domainblocks"
	"testing"

	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"github.com/autumnterror/breezynotes/utils/lang"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
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
		data, err := d.Op(ctx, block, "change_text", map[string]any{
			"new_text": "<html></html>",
		})
		if assert.NoError(t, err) {
			s, err := structpb.NewStruct(data)
			if assert.NoError(t, err) {
				block.Data = s
				cb, err := domainblocks.FromUnifiedToCodeBlock(block)
				log.Blue(format.Struct(cb))
				if assert.NoError(t, err) {
					assert.Equal(t, "<html></html>", cb.Data.Text)
					assert.Equal(t, "HTML", cb.Data.Lang)
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
				cb, err := domainblocks.FromUnifiedToCodeBlock(block)
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
	t.Run(domainblocks.TextBlockType, func(t *testing.T) {
		block := testBlock()
		cb, _ := domainblocks.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.TextBlockType))
		tb, err := domainblocks.FromUnifiedToTextBlock(block)
		assert.NoError(t, err)
		log.Println("after change on text", format.Struct(tb))
		assert.Equal(t, cb.Data.Text, tb.Data.PlainText())
	})

	t.Run(domainblocks.ListBlockUnorderedType, func(t *testing.T) {
		block := testBlock()
		cb, _ := domainblocks.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockUnorderedType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered code", format.Struct(lb))
		assert.Equal(t, cb.Data.Text, lb.Data.TextData.PlainText())
		assert.Equal(t, domainblocks.ListBlockUnorderedType, lb.Data.Type)
	})
	t.Run(domainblocks.ListBlockToDoType, func(t *testing.T) {
		block := testBlock()
		code, _ := domainblocks.FromUnifiedToCodeBlock(block)
		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockToDoType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on unordered", format.Struct(block))
		assert.Equal(t, code.Data.Text, lb.Data.TextData.PlainText())
		assert.Equal(t, domainblocks.ListBlockToDoType, lb.Data.Type)
	})
	t.Run(domainblocks.ListBlockOrderedType, func(t *testing.T) {
		block := testBlock()
		code, _ := domainblocks.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ListBlockOrderedType))
		lb, err := domainblocks.FromUnifiedToListBlock(block)
		assert.NoError(t, err)
		log.Println("after change on ordered", format.Struct(block))
		assert.Equal(t, code.Data.Text, lb.Data.TextData.PlainText())
		assert.Equal(t, domainblocks.ListBlockOrderedType, lb.Data.Type)
	})

	t.Run(domainblocks.HeaderBlockType1, func(t *testing.T) {
		block := testBlock()
		cb, _ := domainblocks.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType1))
		hb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_1", format.Struct(hb))
		assert.Equal(t, cb.Data.Text, hb.Data.TextData.PlainText())
		assert.Equal(t, uint(1), hb.Data.Level)
	})
	t.Run(domainblocks.HeaderBlockType2, func(t *testing.T) {
		block := testBlock()
		code, _ := domainblocks.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType2))
		hb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_2", format.Struct(block))
		assert.Equal(t, code.Data.Text, hb.Data.TextData.PlainText())
		assert.Equal(t, uint(2), hb.Data.Level)
	})

	t.Run(domainblocks.HeaderBlockType3, func(t *testing.T) {
		block := testBlock()
		code, _ := domainblocks.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.HeaderBlockType3))
		hb, err := domainblocks.FromUnifiedToHeaderBlock(block)
		assert.NoError(t, err)
		log.Println("after change on header_3", format.Struct(block))
		assert.Equal(t, code.Data.Text, hb.Data.TextData.PlainText())
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
		code, _ := domainblocks.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.ImgBlockType))
		ib, err := domainblocks.FromUnifiedToImgBlock(block)
		assert.NoError(t, err)
		log.Println("after change on img", format.Struct(block))
		assert.Equal(t, code.Data.Text, ib.Data.Alt)
	})

	t.Run(domainblocks.LinkBlockType, func(t *testing.T) {
		block := testBlock()
		code, _ := domainblocks.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.LinkBlockType))
		linkB, err := domainblocks.FromUnifiedToLinkBlock(block)
		assert.NoError(t, err)
		log.Println("after change on link", format.Struct(block))
		assert.Equal(t, code.Data.Text, linkB.Data.Text)
	})

	t.Run(domainblocks.QuoteBlockType, func(t *testing.T) {
		block := testBlock()
		cb, _ := domainblocks.FromUnifiedToCodeBlock(block)

		assert.NoError(t, d.ChangeType(ctx, block, domainblocks.QuoteBlockType))
		qb, err := domainblocks.FromUnifiedToQuoteBlock(block)
		assert.NoError(t, err)
		log.Println("after change on quote", format.Struct(qb))
		assert.Equal(t, cb.Data.Text, qb.Data.Text)
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
		cb, err := domainblocks.FromUnifiedToCodeBlock(block)
		assert.NoError(t, err)
		assert.Equal(t, "", cb.Data.Text)
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

var (
	d   = Driver{}
	ctx = context.Background()
)

var textTest = "\tpackage main\n\t\timport \"fmt\"\n\t\tfunc main() { fmt.Println(\"Hello\") }"

func testBlock() *brzrpc.Block {
	cb := domainblocks.CodeBlock{
		Id:     "test_code",
		Type:   "code",
		NoteId: "test_note",
		Data: &domainblocks.CodeData{
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
	cb := domainblocks.CodeBlock{
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
	// lexer := lexers.Analyse(code)
	// if lexer != nil {
	// 	lang := lexer.Config().Name
	// 	fmt.Printf("Для фрагмента:\n%s\n... определён язык: %s\n", code, lang)
	// 	fmt.Println("---")
	// } else {
	// 	fmt.Printf("Для фрагмента:\n%s\n... не удалось определить язык.\n", code)
	// 	fmt.Println("---")
	// }
	log.Green(code, lang.AnalyzeLanguage(code))
}
