package fileblock

import (
	"encoding/json"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2/domainblocks"
)

// changeSrc handles the logic for updating the file source URL.
func changeSrc(b *domainblocks.FileBlock, raw []byte) (map[string]any, error) {
	if b.Data == nil {
		b.Data = &domainblocks.FileData{}
	}

	var req struct {
		NewSrc string `json:"new_src"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}

	b.Data.Src = req.NewSrc

	nb, err := b.ToUnified()
	if err != nil {
		return nil, err
	}

	return nb.GetData().AsMap(), nil
}
