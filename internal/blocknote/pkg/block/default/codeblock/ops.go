package codeblock

import (
	"encoding/json"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2/domainblocks"
	"github.com/autumnterror/breezynotes/utils/lang"
)

func changeText(b *domainblocks.CodeBlock, raw []byte) (map[string]any, error) {
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
	b.Data.Lang = lang.AnalyzeLanguage(b.Data.Text)

	nb, err := b.ToUnified()
	if err != nil {
		return nil, err
	}

	return nb.GetData().AsMap(), nil
}

func analyseLang(b *domainblocks.CodeBlock) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}

	b.Data.Lang = lang.AnalyzeLanguage(b.Data.Text)

	nb, err := b.ToUnified()
	if err != nil {
		return nil, err
	}

	return nb.GetData().AsMap(), nil
}
