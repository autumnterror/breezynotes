package codeblock

import (
	"encoding/json"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/autumnterror/breezynotes/pkg/domain"
)

func changeText(b *domain.CodeBlock, raw []byte) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}

	var req struct {
		NewText string `json:"new_text"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}

	b.Data.Text = req.NewText

	nb, err := b.ToUnified()
	if err != nil {
		return nil, err
	}

	return nb.GetData().AsMap(), nil
}

func analyseLang(b *domain.CodeBlock) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}

	lexer := lexers.Analyse(b.Data.Text)
	if lexer == nil {
		b.Data.Lang = ""
	} else {
		b.Data.Lang = lexer.Config().Name
	}

	nb, err := b.ToUnified()
	if err != nil {
		return nil, err
	}

	return nb.GetData().AsMap(), nil
}
