package domainblocks

import (
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuote(t *testing.T) {
	op := "test quote block"
	t.Run(op, func(t *testing.T) {
		test := QuoteBlock{
			Id:        "quote-test-id",
			Type:      "quote",
			NoteId:    "note-id",
			CreatedAt: 112233,
			UpdatedAt: 445566,
			IsUsed:    false,
			Data: &QuoteData{
				Text: "To be, or not to be, that is the question.",
			},
		}

		unifiedBlock, err := test.ToUnified()
		if assert.NoError(t, err) {
			log.Println(op+": QuoteBlock.ToUnified()", format.Struct(unifiedBlock))
		}

		newQuoteBlock, err := FromUnifiedToQuoteBlock(unifiedBlock)
		if assert.NoError(t, err) {
			log.Println(op+": FromUnifiedToQuoteBlock()", format.Struct(newQuoteBlock))
		}

		assert.Equal(t, test, *newQuoteBlock)
	})
}
