package domainblocks

import (
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFile(t *testing.T) {
	op := "test file block"
	t.Run(op, func(t *testing.T) {
		file := FileBlock{
			Id:        "test",
			Type:      "text",
			NoteId:    "test",
			CreatedAt: 0,
			UpdatedAt: 0,
			IsUsed:    false,
			Data: &FileData{
				Src: "images/example.png",
			},
		}

		unif, err := file.ToUnified()
		if assert.NoError(t, err) {
			log.Println(op+"file.ToUnified()", format.Struct(unif))
		}
		nFile, err := FromUnifiedToFileBlock(unif)
		if assert.NoError(t, err) {
			log.Println(op+" FromUnifiedToFileBlock(unif)", format.Struct(nFile))
		}

		assert.Equal(t, file, *nFile)
	})
}
