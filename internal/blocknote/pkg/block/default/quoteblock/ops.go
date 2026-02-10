package quoteblock

import (
	"encoding/json"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2/domainblocks"
)

// changeText handles the logic for updating the quote's text.
func changeText(b *domainblocks.QuoteBlock, raw []byte) (map[string]any, error) {
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

	return b.Data.ToMap(), nil
}
