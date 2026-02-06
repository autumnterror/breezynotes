package domainblocks

import (
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLink(t *testing.T) {
	op := "test link block"
	t.Run(op, func(t *testing.T) {
		test := LinkBlock{
			Id:        "link-test-id",
			Type:      "link",
			NoteId:    "note-id",
			CreatedAt: 12345,
			UpdatedAt: 67890,
			IsUsed:    true,
			Data: &LinkData{
				Text: "https://example.com",
			},
		}

		unifiedBlock, err := test.ToUnified()
		if assert.NoError(t, err) {
			log.Println(op+": LinkBlock.ToUnified()", format.Struct(unifiedBlock))
		}

		newLinkBlock, err := FromUnifiedToLinkBlock(unifiedBlock)
		if assert.NoError(t, err) {
			log.Println(op+": FromUnifiedToLinkBlock()", format.Struct(newLinkBlock))
		}

		assert.Equal(t, test, *newLinkBlock)
	})
}
