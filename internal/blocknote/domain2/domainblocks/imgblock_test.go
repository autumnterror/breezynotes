package domainblocks

import (
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImg(t *testing.T) {
	op := "test img block"
	t.Run(op, func(t *testing.T) {
		img := ImgBlock{
			Id:        "test",
			Type:      "text",
			NoteId:    "test",
			CreatedAt: 0,
			UpdatedAt: 0,
			IsUsed:    false,
			Data: &ImgData{
				Src: "images/example.png",
				Alt: "no photo",
			},
		}

		unif, err := img.ToUnified()
		if assert.NoError(t, err) {
			log.Println(op+"img.ToUnified()", format.Struct(unif))
		}
		nImg, err := FromUnifiedToImgBlock(unif)
		if assert.NoError(t, err) {
			log.Println(op+" FromUnifiedToImgBlock(unif)", format.Struct(nImg))
		}

		assert.Equal(t, img, *nImg)
	})
}
