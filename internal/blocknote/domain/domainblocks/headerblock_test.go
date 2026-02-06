package domainblocks

import (
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/text"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHeader(t *testing.T) {
	op := "test header block"
	t.Run(op, func(t *testing.T) {
		test := HeaderBlock{
			Id:        "test",
			Type:      "header",
			NoteId:    "test",
			CreatedAt: 0,
			UpdatedAt: 0,
			IsUsed:    false,
			Data: &HeaderData{
				TextData: &text.Data{
					Text: []text.Part{
						{
							Style:  "default",
							String: "test",
						},
					}},
				Level: 2,
			},
		}
		header, err := test.ToUnified()
		if assert.NoError(t, err) {
			log.Println(op+"HeaderBlock.ToUnified()", format.Struct(header))
		}

		newHeader, err := FromUnifiedToHeaderBlock(header)
		if assert.NoError(t, err) {
			log.Println(op+"HeaderBlock.FromUnifiedToHeaderBlock()", format.Struct(newHeader))
		}
		assert.Equal(t, test, *newHeader)

		if assert.NoError(t, newHeader.Data.TextData.ApplyStyle(0, 2, "new")) {
			if assert.Equal(t, len(newHeader.Data.TextData.Text), 2) {
				assert.Equal(t, newHeader.Data.TextData.Text[0].Style, "new")
			}
			log.Println(op+" .ApplyStyle(0, 2, \"new\")", format.Struct(newHeader))
		}
	})

}
