package block

import (
	"github.com/autumnterror/breezynotes/pkg/domain"
	"github.com/autumnterror/breezynotes/pkg/text"
)

func ChangeTypeUnif(textData *text.Data, plainText, newType string, levelList uint, valueOrdered int) (map[string]any, error) {
	if valueOrdered == 0 {
		valueOrdered = 1
	}

	var newData map[string]any
	switch newType {
	case domain.TextBlockType:
		newData = textData.ToMap()
	case domain.ListBlockToDoType:
		nd := domain.ListData{
			TextData: textData,
			Level:    levelList,
			Type:     domain.ListBlockToDoType,
			Value:    0,
		}
		newData = nd.ToMap()
	case domain.ListBlockUnorderedType:
		nd := domain.ListData{
			TextData: textData,
			Level:    levelList,
			Type:     domain.ListBlockUnorderedType,
			Value:    0,
		}
		newData = nd.ToMap()
	case domain.ListBlockOrderedType:
		nd := domain.ListData{
			TextData: textData,
			Level:    levelList,
			Type:     domain.ListBlockOrderedType,
			Value:    valueOrdered,
		}
		newData = nd.ToMap()
	case domain.CodeBlockType:
		nd := domain.CodeData{
			Text: plainText,
			Lang: "undefined",
		}
		newData = nd.ToMap()
	case domain.HeaderBlockType1:
		nd := domain.HeaderData{
			TextData: textData,
			Level:    1,
		}
		newData = nd.ToMap()
	case domain.HeaderBlockType2:
		nd := domain.HeaderData{
			TextData: textData,
			Level:    2,
		}
		newData = nd.ToMap()
	case domain.HeaderBlockType3:
		nd := domain.HeaderData{
			TextData: textData,
			Level:    3,
		}
		newData = nd.ToMap()
	case domain.FileBlockType:
		nd := domain.FileData{
			Src: "",
		}
		newData = nd.ToMap()
	case domain.LinkBlockType:
		nd := domain.LinkData{
			Text: plainText,
		}
		newData = nd.ToMap()
	case domain.ImgBlockType:
		nd := domain.ImgData{
			Src: "",
			Alt: plainText,
		}
		newData = nd.ToMap()
	case domain.QuoteBlockType:
		nd := domain.QuoteData{
			Text: plainText,
		}
		newData = nd.ToMap()
	default:
		return nil, domain.ErrUnsupportedType
	}
	return newData, nil
}
