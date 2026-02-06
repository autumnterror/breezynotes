package domainblocks

import (
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCode(t *testing.T) {
	op := "test code block"
	t.Run(op, func(t *testing.T) {
		test := CodeBlock{
			Id:        "code-test-id",
			Type:      "code",
			NoteId:    "note-id",
			CreatedAt: 54321,
			UpdatedAt: 9876,
			IsUsed:    false,
			Data: &CodeData{
				Text: "fmt.Println(\"Hello, World!\")",
				Lang: "go",
			},
		}

		unifiedBlock, err := test.ToUnified()
		if assert.NoError(t, err) {
			log.Println(op+": CodeBlock.ToUnified()", format.Struct(unifiedBlock))
		}

		newCodeBlock, err := FromUnifiedToCodeBlock(unifiedBlock)
		if assert.NoError(t, err) {
			log.Println(op+": FromUnifiedToCodeBlock()", format.Struct(newCodeBlock))
		}

		assert.Equal(t, test, *newCodeBlock)
	})
}
