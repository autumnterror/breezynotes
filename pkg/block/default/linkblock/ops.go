package linkblock

import (
	"encoding/json"
	"github.com/autumnterror/breezynotes/pkg/domain"
)

// changeText handles updating the link's display text.
func changeText(b *domain.LinkBlock, raw []byte) (map[string]any, error) {
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

// changeText handles updating the link's display text.
func changeUrl(b *domain.LinkBlock, raw []byte) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}

	var req struct {
		NewUrl string `json:"new_url"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}

	b.Data.Url = req.NewUrl
	return b.Data.ToMap(), nil
}
