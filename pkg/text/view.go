package text

import "fmt"

type Data struct {
	Text []Part `json:"text" bson:"text"`
}

func (tb *Data) ToMap() map[string]any {
	if tb == nil {
		return nil
	}
	if tb.Text != nil && len(tb.Text) > 0 {
		textArr := make([]any, 0, len(tb.Text))
		for _, t := range tb.Text {
			textArr = append(textArr, map[string]any{
				"style":  t.Style,
				"string": t.String,
			})
		}

		return map[string]any{
			"text": textArr,
		}
	}
	return nil
}

type Part struct {
	Style  string `json:"style" bson:"style"`
	String string `json:"string" bson:"string"`
}

func NewDataFromMap(m map[string]any) (*Data, error) {
	rawText, ok := m["text"]
	if !ok {
		return nil, nil
	}

	list, ok := rawText.([]any)
	if !ok {
		return nil, fmt.Errorf(`field "text" has unexpected type %T, want []any`, rawText)
	}

	texts := make([]Part, 0, len(list))
	for i, v := range list {
		obj, ok := v.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("text[%d] has unexpected type %T, want map[string]any", i, v)
		}

		style, styleOk := obj["style"].(string)
		str, strOk := obj["string"].(string)

		if !styleOk || !strOk {
			return nil, fmt.Errorf("text[%d] has missing or invalid 'style' or 'string' fields", i)
		}

		texts = append(texts, Part{
			Style:  style,
			String: str,
		})
	}

	return &Data{Text: texts}, nil
}
