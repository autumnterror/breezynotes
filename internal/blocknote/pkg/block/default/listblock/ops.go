package listblock

import (
	"encoding/json"
	"github.com/autumnterror/breezynotes/internal/blocknote/domain2/domainblocks"
	"github.com/autumnterror/utils_go/pkg/log"
)

func applyStyleOp(b *domainblocks.ListBlock, raw []byte) (map[string]any, error) {
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
func insertTextOp(b *domainblocks.ListBlock, raw []byte) (map[string]any, error) {
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

func deleteRangeOp(b *domainblocks.ListBlock, raw []byte) (map[string]any, error) {
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

func changeType(b *domainblocks.ListBlock, raw []byte) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}

	var req struct {
		NewType      string `json:"new_type"`
		OrderedValue int    `json:"ordered_value"`
	}
	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}
	switch req.NewType {
	case domainblocks.ListBlockToDoType:
		b.Data.Value = 0
	case domainblocks.ListBlockUnorderedType:
		b.Data.Value = 0
	case domainblocks.ListBlockOrderedType:
		b.Data.Value = req.OrderedValue
	default:
		return nil, domainblocks.ErrUnsupportedType
	}
	b.Data.Type = req.NewType

	nb, err := b.ToUnified()
	if err != nil {
		return nil, err
	}

	return nb.GetData().AsMap(), nil
}

func changeValue(b *domainblocks.ListBlock, raw []byte) (map[string]any, error) {
	if b.Data == nil {
		return nil, nil
	}
	var req struct {
		NewValue int `json:"new_value"`
	}

	if err := json.Unmarshal(raw, &req); err != nil {
		return nil, err
	}
	switch b.Data.Type {
	case domainblocks.ListBlockToDoType:
		if req.NewValue >= 1 {
			b.Data.Value = 1
		} else {
			b.Data.Value = 0
		}
	case domainblocks.ListBlockUnorderedType:
		b.Data.Value = 0
	case domainblocks.ListBlockOrderedType:
		if req.NewValue < 1 {
			b.Data.Value = 1
		} else {
			b.Data.Value = req.NewValue
		}
	default:
		return nil, domainblocks.ErrUnsupportedType
	}
	log.Green(b.Data.Value)
	nb, err := b.ToUnified()
	if err != nil {
		return nil, err
	}

	return nb.GetData().AsMap(), nil
}
func changeLevel(b *domainblocks.ListBlock, raw []byte) (map[string]any, error) {
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
