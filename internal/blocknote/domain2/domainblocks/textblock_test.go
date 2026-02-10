package domainblocks

import (
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/text"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestText(t *testing.T) {
	op := "test text block"
	t.Run(op, func(t *testing.T) {
		txt := TextBlock{
			Id:        "test",
			Type:      "text",
			NoteId:    "test",
			CreatedAt: 0,
			UpdatedAt: 0,
			IsUsed:    false,
			Data: &text.Data{Text: []text.Part{
				{
					Style:  "default",
					String: "test def",
				},
				{
					Style:  "bald",
					String: "test bald",
				},
			},
			},
		}

		unif, err := txt.ToUnified()
		if assert.NoError(t, err) {
			log.Println(op+"txt.ToUnified()", format.Struct(unif))
		}
		ntxt, err := FromUnifiedToTextBlock(unif)
		if assert.NoError(t, err) {
			log.Println(op+" FromUnifiedToTextBlock(unif)", format.Struct(ntxt))
		}

		assert.Equal(t, txt, *ntxt)
		if assert.NoError(t, ntxt.Data.ApplyStyle(0, 2, "new")) {
			if assert.Equal(t, len(ntxt.Data.Text), 3) {
				assert.Equal(t, ntxt.Data.Text[0].Style, "new")
			}
			log.Println(op+" .ApplyStyle(0, 2, \"new\")", format.Struct(ntxt))
		}
	})
}
