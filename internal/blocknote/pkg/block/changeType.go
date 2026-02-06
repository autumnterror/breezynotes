package block

import (
	"github.com/autumnterror/breezynotes/internal/blocknote/domain/domainblocks"
	"github.com/autumnterror/breezynotes/internal/blocknote/pkg/text"
	"github.com/autumnterror/breezynotes/utils/lang"
)

func ChangeTypeUnif(textData *text.Data, plainText, newType string, levelList uint, valueOrdered int) (map[string]any, error) {
	if valueOrdered == 0 {
		valueOrdered = 1
	}

	var newData map[string]any
	switch newType {
	case domainblocks.TextBlockType:
		newData = textData.ToMap()
	case domainblocks.ListBlockToDoType:
		nd := domainblocks.ListData{
			TextData: textData,
			Level:    levelList,
			Type:     domainblocks.ListBlockToDoType,
			Value:    0,
		}
		newData = nd.ToMap()
	case domainblocks.ListBlockUnorderedType:
		nd := domainblocks.ListData{
			TextData: textData,
			Level:    levelList,
			Type:     domainblocks.ListBlockUnorderedType,
			Value:    0,
		}
		newData = nd.ToMap()
	case domainblocks.ListBlockOrderedType:
		nd := domainblocks.ListData{
			TextData: textData,
			Level:    levelList,
			Type:     domainblocks.ListBlockOrderedType,
			Value:    valueOrdered,
		}
		newData = nd.ToMap()
	case domainblocks.CodeBlockType:
		nd := domainblocks.CodeData{
			Text: plainText,
			Lang: lang.AnalyzeLanguage(plainText),
		}

		newData = nd.ToMap()
	case domainblocks.HeaderBlockType1:
		nd := domainblocks.HeaderData{
			TextData: textData,
			Level:    1,
		}
		newData = nd.ToMap()
	case domainblocks.HeaderBlockType2:
		nd := domainblocks.HeaderData{
			TextData: textData,
			Level:    2,
		}
		newData = nd.ToMap()
	case domainblocks.HeaderBlockType3:
		nd := domainblocks.HeaderData{
			TextData: textData,
			Level:    3,
		}
		newData = nd.ToMap()
	case domainblocks.FileBlockType:
		nd := domainblocks.FileData{
			Src: "",
		}
		newData = nd.ToMap()
	case domainblocks.LinkBlockType:
		nd := domainblocks.LinkData{
			Text: plainText,
		}
		newData = nd.ToMap()
	case domainblocks.ImgBlockType:
		nd := domainblocks.ImgData{
			Src: "",
			Alt: plainText,
		}
		newData = nd.ToMap()
	case domainblocks.QuoteBlockType:
		nd := domainblocks.QuoteData{
			Text: plainText,
		}
		newData = nd.ToMap()
	default:
		return nil, domainblocks.ErrUnsupportedType
	}
	return newData, nil
}
