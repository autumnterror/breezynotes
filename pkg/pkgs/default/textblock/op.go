package textblock

import (
	"errors"
	"github.com/autumnterror/breezynotes/pkg/utils/alg"
	"runtime"
)

func (tb *TextBlock) PlainText() string {
	total := make([]rune, 0, 128)
	for _, s := range tb.Text {
		total = append(total, []rune(s.Text)...)
	}
	return string(total)
}

func (tb *TextBlock) ApplyStyle(start, end int, style string) error {
	if start >= end {
		return errors.New("invalid range: start >= end")
	}

	//получаем длины
	prefSumText, total := buildPrefixLens(tb.Text)
	//валидация
	if total == 0 {
		return nil
	}
	if start < 0 {
		start = 0
	}
	if end > total {
		end = total
	}
	if start >= end {
		return errors.New("start >= end")
	}

	// найти сегменты
	segStart, offStart := findSegmentByPos(tb.Text, prefSumText, start)
	// соберём новый слайс
	var newData []TextData

	// добавить сегменты до segStart
	if segStart > 0 {
		newData = append(newData, alg.SliceCopy(tb.Text, 0, segStart)...)
	}

	// обработать стартовый сегмент (если вставляем в середину)
	if segStart < len(tb.Text) && offStart > 0 {
		seg := tb.Text[segStart]
		r := []rune(seg.Text)
		left := string(r[:offStart])
		mid := string(r[offStart:])

		// left keeps old style
		if left != "" {
			newData = append(newData, TextData{Style: seg.Style, Text: left})
		}

		// Создаем временный tail для обработки и помещаем все что после начала обрабатываемого фрагмента
		var tail []TextData
		tail = append(tail, TextData{Style: seg.Style, Text: mid})
		if segStart+1 < len(tb.Text) {
			tail = append(tail, alg.SliceCopy(tb.Text, segStart+1, len(tb.Text))...)
		}

		//пересчитываем длины без учета необрабатываемого куска
		tailPref, _ := buildPrefixLens(tail)
		var sumBefore int
		if segStart > 0 {
			sumBefore = prefSumText[segStart-1]
		}
		startProc := start - (sumBefore + offStart)
		endProc := end - (sumBefore + offStart)

		// ищем позиции в tail
		segStartProc, offStartProc := findSegmentByPos(tail, tailPref, startProc)
		segEndProc, offEndProc := findSegmentByPos(tail, tailPref, endProc)

		if segStartProc < len(tail) {
			seg := tail[segStartProc]
			//сохраняем нерабочую часть
			if offStartProc > 0 {
				r := []rune(seg.Text)
				left := string(r[:offStartProc])
				newData = append(newData, TextData{Style: seg.Style, Text: left})
			}

			//если работа в одном сегменте
			if segStartProc == segEndProc {
				r := []rune(seg.Text)
				startIdx := offStartProc
				endIdx := offEndProc
				//выделяем в отдельный блок
				if endIdx > startIdx {
					midText := string(r[startIdx:endIdx])
					newData = append(newData, TextData{Style: style, Text: midText})
				}
				//оставшееся в сегменте сохраняем
				right := string(r[endIdx:])
				if right != "" {
					newData = append(newData, TextData{Style: seg.Style, Text: right})
				}
				//оставшееся в массиве сохраняем
				if segEndProc+1 < len(tail) {
					newData = append(newData, alg.SliceCopy(tail, segStartProc+1, len(tail))...)
				}
			} else {
				//если работа в нескольких сегментах
				r := []rune(seg.Text)
				if offStartProc < len(r) {
					//сохраняем часть текущего сегмента
					mid := string(r[offStartProc:])
					if mid != "" {
						newData = append(newData, TextData{Style: style, Text: mid})
					}
				}
				for k := segStartProc + 1; k < segEndProc && k < len(tail); k++ {
					//идем по сегментам, сохраняя данные
					segK := tail[k]
					newData = append(newData, TextData{Style: style, Text: segK.Text})
				}
				//сохраняем оставшиеся в последнем рабочем сегменте
				if segEndProc < len(tail) {
					rEnd := []rune(tail[segEndProc].Text)
					if offEndProc > 0 {
						left := string(rEnd[:offEndProc])
						newData = append(newData, TextData{Style: style, Text: left})
					}
					//после обработки в тот же стиль
					if offEndProc < len(rEnd) {
						right := string(rEnd[offEndProc:])
						newData = append(newData, TextData{Style: tail[segEndProc].Style, Text: right})
					}
					//сохраняем оставшиеся в массиве
					if segEndProc+1 < len(tail) {
						newData = append(newData, alg.SliceCopy(tail, segEndProc+1, len(tail))...)
					}
				}
			}
		}
		// нормализация
		tb.Text = MergeSameStylesParallel(newData, runtime.NumCPU())
		return nil
	}

	if segStart < len(tb.Text) {
		cursor := 0
		if segStart > 0 {
			cursor = prefSumText[segStart-1]
		}
		i := segStart
		for i < len(tb.Text) && cursor < end {
			seg := tb.Text[i]
			r := []rune(seg.Text)
			segLen := len(r)
			segStartPos := cursor
			segEndPos := cursor + segLen

			overlapStart := alg.MaxInt(start, segStartPos)
			overlapEnd := alg.MinInt(end, segEndPos)

			if overlapStart > segStartPos {
				leftCnt := overlapStart - segStartPos
				if leftCnt > 0 {
					newData = append(newData, TextData{Style: seg.Style, Text: string(r[:leftCnt])})
				}
			}

			if overlapEnd > overlapStart {
				midL := overlapStart - segStartPos
				midR := overlapEnd - segStartPos
				newData = append(newData, TextData{Style: style, Text: string(r[midL:midR])})
			}

			if overlapEnd < segEndPos {
				rightL := overlapEnd - segStartPos
				newData = append(newData, TextData{Style: seg.Style, Text: string(r[rightL:])})
			}
			cursor += segLen
			i++
		}

		if i < len(tb.Text) {
			newData = append(newData, alg.SliceCopy(tb.Text, i, len(tb.Text))...)
		}
	}

	// normalize
	tb.Text = MergeSameStyles(newData)
	return nil
}

// InsertText вставляет текст newText в позицию pos; newText получает стиль "default".
// Сохраняет существующие стили по бокам.
func (tb *TextBlock) InsertText(pos int, newText string) error {
	if newText == "" {
		return nil
	}
	pref, total := buildPrefixLens(tb.Text)
	if pos < 0 {
		pos = 0
	}
	if pos > total {
		pos = total
	}

	if len(tb.Text) == 0 {
		tb.Text = []TextData{{Style: "default", Text: newText}}
		return nil
	}

	segIdx, offset := findSegmentByPos(tb.Text, pref, pos)

	if segIdx == len(tb.Text) {
		tb.Text = append(tb.Text, TextData{Style: "default", Text: newText})
		return nil
	}

	seg := tb.Text[segIdx]
	r := []rune(seg.Text)
	left := string(r[:offset])
	right := string(r[offset:])

	newSegs := make([]TextData, 0, len(tb.Text)+2)

	if segIdx > 0 {
		newSegs = append(newSegs, tb.Text[:segIdx]...)
	}

	if left != "" {
		newSegs = append(newSegs, TextData{Style: seg.Style, Text: left})
	}

	newSegs = append(newSegs, TextData{Style: "default", Text: newText})

	if right != "" {
		newSegs = append(newSegs, TextData{Style: seg.Style, Text: right})
	}

	if segIdx+1 < len(tb.Text) {
		newSegs = append(newSegs, tb.Text[segIdx+1:]...)
	}

	tb.Text = MergeSameStyles(newSegs)
	return nil
}

// DeleteRange удаляет диапазон [start, end)
func (tb *TextBlock) DeleteRange(start, end int) error {
	if start >= end {
		return nil
	}
	pref, total := buildPrefixLens(tb.Text)
	if total == 0 {
		return nil
	}
	if start < 0 {
		start = 0
	}
	if end > total {
		end = total
	}
	if start >= end {
		return nil
	}

	segStart, offStart := findSegmentByPos(tb.Text, pref, start)
	segEnd, offEnd := findSegmentByPos(tb.Text, pref, end)

	newSegs := make([]TextData, 0, len(tb.Text))

	// Добавляем сегменты до segStart
	if segStart > 0 {
		newSegs = append(newSegs, tb.Text[:segStart]...)
	}

	// Обрабатываем начальный сегмент (если нужно)
	if segStart < len(tb.Text) && offStart > 0 {
		seg := tb.Text[segStart]
		r := []rune(seg.Text)
		left := string(r[:offStart])
		if left != "" {
			newSegs = append(newSegs, TextData{Style: seg.Style, Text: left})
		}
	}

	// Обрабатываем конечный сегмент (если нужно)
	if segEnd < len(tb.Text) && offEnd > 0 {
		seg := tb.Text[segEnd]
		r := []rune(seg.Text)
		right := string(r[offEnd:])
		if right != "" {
			newSegs = append(newSegs, TextData{Style: seg.Style, Text: right})
		}
	}

	// Добавляем сегменты после segEnd
	if segEnd+1 < len(tb.Text) {
		newSegs = append(newSegs, tb.Text[segEnd+1:]...)
	}

	tb.Text = MergeSameStyles(newSegs)
	return nil
}
