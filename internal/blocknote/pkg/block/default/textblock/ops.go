package textblock

import (
	"encoding/json"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain/domainblocks"
)

func applyStyleOp(b *domainblocks.TextBlock, raw []byte) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}
	var req struct {
		Start int    `json:"start"`
		End   int    `json:"end"`
		Style string `json:"style"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}
	if err := b.Data.ApplyStyle(req.Start, req.End, req.Style); err != nil {
		return nil, err
	}
	nb, err := b.ToUnified()
	if err != nil {
		return nil, err
	}

	return nb.GetData().AsMap(), nil
}
func insertTextOp(b *domainblocks.TextBlock, raw []byte) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}
	var req struct {
		Pos     int    `json:"pos"`
		NewText string `json:"new_text"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}
	if err := b.Data.InsertText(req.Pos, req.NewText); err != nil {
		return nil, err
	}
	nb, err := b.ToUnified()
	if err != nil {
		return nil, err
	}

	return nb.GetData().AsMap(), nil
}

func deleteRangeOp(b *domainblocks.TextBlock, raw []byte) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}
	var req struct {
		Start int `json:"start"`
		End   int `json:"end"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}
	if err := b.Data.DeleteRange(req.Start, req.End); err != nil {
		return nil, err
	}
	nb, err := b.ToUnified()
	if err != nil {
		return nil, err
	}

	return nb.GetData().AsMap(), nil
}
