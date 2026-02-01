package imgblock

import (
	"encoding/json"
	"github.com/autumnterror/breezynotes/pkg/domain"
)

// changeSrc handles updating the image source URL.
func changeSrc(b *domain.ImgBlock, raw []byte) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}

	var req struct {
		NewSrc string `json:"new_src"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}

	b.Data.Src = req.NewSrc
	return b.Data.ToMap(), nil
}

// changeAlt handles updating the image alt text.
func changeAlt(b *domain.ImgBlock, raw []byte) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}

	var req struct {
		NewAlt string `json:"new_alt"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}

	b.Data.Alt = req.NewAlt
	return b.Data.ToMap(), nil
}
