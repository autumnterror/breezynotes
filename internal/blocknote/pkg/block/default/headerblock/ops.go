package headerblock

import (
	"encoding/json"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain/domainblocks"
)

func applyStyleOp(b *domainblocks.HeaderBlock, raw []byte) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}
	if b.Data.TextData == nil {
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
	if err := b.Data.TextData.ApplyStyle(req.Start, req.End, req.Style); err != nil {
		return nil, err
	}
	nb, err := b.ToUnified()
	if err != nil {
		return nil, err
	}

	return nb.GetData().AsMap(), nil
}
func insertTextOp(b *domainblocks.HeaderBlock, raw []byte) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}
	if b.Data.TextData == nil {
		return nil, nil
	}

	var req struct {
		Pos     int    `json:"pos"`
		NewText string `json:"new_text"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}
	if err := b.Data.TextData.InsertText(req.Pos, req.NewText); err != nil {
		return nil, err
	}
	nb, err := b.ToUnified()
	if err != nil {
		return nil, err
	}

	return nb.GetData().AsMap(), nil
}

func deleteRangeOp(b *domainblocks.HeaderBlock, raw []byte) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}
	if b.Data.TextData == nil {
		return nil, nil
	}

	var req struct {
		Start int `json:"start"`
		End   int `json:"end"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}
	if err := b.Data.TextData.DeleteRange(req.Start, req.End); err != nil {
		return nil, err
	}
	nb, err := b.ToUnified()
	if err != nil {
		return nil, err
	}

	return nb.GetData().AsMap(), nil
}

func changeLevel(b *domainblocks.HeaderBlock, raw []byte) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}

	var req struct {
		NewLevel int `json:"new_level"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}
	if req.NewLevel < 0 {
		req.NewLevel = 0
	}
	b.Data.Level = uint(req.NewLevel)

	nb, err := b.ToUnified()
	if err != nil {
		return nil, err
	}

	return nb.GetData().AsMap(), nil
}
