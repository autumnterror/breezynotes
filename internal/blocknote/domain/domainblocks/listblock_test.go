package domainblocks

import (
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/text"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestList(t *testing.T) {
	op := "test list block"
	t.Run(op, func(t *testing.T) {
		test := ListBlock{
			Id:        "test",
			Type:      "list",
			NoteId:    "test",
			CreatedAt: 0,
			UpdatedAt: 0,
			IsUsed:    false,
			Data: &ListData{
				TextData: &text.Data{
					Text: []text.Part{
						{
							Style:  "default",
							String: "test",
						},
					}},
				Level: 2,
				Type:  ListBlockToDoType,
				Value: 1,
			},
		}
		unif, err := test.ToUnified()
		if assert.NoError(t, err) {
			log.Println(op+"ListBlock.ToUnified()", format.Struct(unif))
		}

		newTest, err := FromUnifiedToListBlock(unif)
		if assert.NoError(t, err) {
			log.Println(op+"ListBlock.FromUnifiedToListBlock()", format.Struct(newTest))
		}
		assert.Equal(t, test, *newTest)

		if assert.NoError(t, newTest.Data.TextData.ApplyStyle(0, 2, "new")) {
			if assert.Equal(t, len(newTest.Data.TextData.Text), 2) {
				assert.Equal(t, newTest.Data.TextData.Text[0].Style, "new")
			}
			log.Println(op+" .ApplyStyle(0, 2, \"new\")", format.Struct(newTest))
		}
	})

}
