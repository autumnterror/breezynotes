package text

import (
	"errors"
	"github.com/autumnterror/utils_go/pkg/utils/alg"
	"runtime"
)

func (tb *Data) PlainText() string {
	if tb == nil {
		return ""
	}
	total := make([]rune, 0, 128)
	for _, s := range tb.Text {
		total = append(total, []rune(s.String)...)
	}
	return string(total)
}

func (tb *Data) ApplyStyle(start, end int, style string) error {
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
	var newData []Part

	// добавить сегменты до segStart
	if segStart > 0 {
		newData = append(newData, alg.SliceCopy(tb.Text, 0, segStart)...)
	}

	// обработать стартовый сегмент (если вставляем в середину)
	if segStart < len(tb.Text) && offStart > 0 {
		seg := tb.Text[segStart]
		r := []rune(seg.String)
		left := string(r[:offStart])
		mid := string(r[offStart:])

		// left keeps old style
		if left != "" {
			newData = append(newData, Part{Style: seg.Style, String: left})
		}

		// Создаем временный tail для обработки и помещаем все что после начала обрабатываемого фрагмента
		var tail []Part
		tail = append(tail, Part{Style: seg.Style, String: mid})
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
				r := []rune(seg.String)
				left := string(r[:offStartProc])
				newData = append(newData, Part{Style: seg.Style, String: left})
			}

			//если работа в одном сегменте
			if segStartProc == segEndProc {
				r := []rune(seg.String)
				startIdx := offStartProc
				endIdx := offEndProc
				//выделяем в отдельный блок
				if endIdx > startIdx {
					midText := string(r[startIdx:endIdx])
					newData = append(newData, Part{Style: style, String: midText})
				}
				//оставшееся в сегменте сохраняем
				right := string(r[endIdx:])
				if right != "" {
					newData = append(newData, Part{Style: seg.Style, String: right})
				}
				//оставшееся в массиве сохраняем
				if segEndProc+1 < len(tail) {
					newData = append(newData, alg.SliceCopy(tail, segStartProc+1, len(tail))...)
				}
			} else {
				//если работа в нескольких сегментах
				r := []rune(seg.String)
				if offStartProc < len(r) {
					//сохраняем часть текущего сегмента
					mid := string(r[offStartProc:])
					if mid != "" {
						newData = append(newData, Part{Style: style, String: mid})
					}
				}
				for k := segStartProc + 1; k < segEndProc && k < len(tail); k++ {
					//идем по сегментам, сохраняя данные
					segK := tail[k]
					newData = append(newData, Part{Style: style, String: segK.String})
				}
				//сохраняем оставшиеся в последнем рабочем сегменте
				if segEndProc < len(tail) {
					rEnd := []rune(tail[segEndProc].String)
					if offEndProc > 0 {
						left := string(rEnd[:offEndProc])
						newData = append(newData, Part{Style: style, String: left})
					}
					//после обработки в тот же стиль
					if offEndProc < len(rEnd) {
						right := string(rEnd[offEndProc:])
						newData = append(newData, Part{Style: tail[segEndProc].Style, String: right})
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
			r := []rune(seg.String)
			segLen := len(r)
			segStartPos := cursor
			segEndPos := cursor + segLen

			overlapStart := alg.MaxInt(start, segStartPos)
			overlapEnd := alg.MinInt(end, segEndPos)

			if overlapStart > segStartPos {
				leftCnt := overlapStart - segStartPos
				if leftCnt > 0 {
					newData = append(newData, Part{Style: seg.Style, String: string(r[:leftCnt])})
				}
			}

			if overlapEnd > overlapStart {
				midL := overlapStart - segStartPos
				midR := overlapEnd - segStartPos
				newData = append(newData, Part{Style: style, String: string(r[midL:midR])})
			}

			if overlapEnd < segEndPos {
				rightL := overlapEnd - segStartPos
				newData = append(newData, Part{Style: seg.Style, String: string(r[rightL:])})
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
func (tb *Data) InsertText(pos int, newText string) error {
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
		tb.Text = []Part{{Style: "default", String: newText}}
		return nil
	}

	segIdx, offset := findSegmentByPos(tb.Text, pref, pos)

	if segIdx == len(tb.Text) {
		if segIdx > 0 {
			seg := tb.Text[segIdx-1]
			tb.Text = append(tb.Text, Part{Style: seg.Style, String: newText})
			tb.Text = MergeSameStyles(tb.Text)
			return nil
		}
		tb.Text = append(tb.Text, Part{Style: "default", String: newText})
		return nil
	}

	seg := tb.Text[segIdx]
	r := []rune(seg.String)
	left := string(r[:offset])
	right := string(r[offset:])

	newSegs := make([]Part, 0, len(tb.Text)+2)

	if segIdx > 0 {
		newSegs = append(newSegs, tb.Text[:segIdx]...)
	}

	if left != "" {
		newSegs = append(newSegs, Part{Style: seg.Style, String: left})
	}
	if offset == 0 && segIdx > 0 {
		segBack := tb.Text[segIdx-1]
		newSegs = append(newSegs, Part{Style: segBack.Style, String: newText})
	} else {
		newSegs = append(newSegs, Part{Style: seg.Style, String: newText})
	}

	if right != "" {
		newSegs = append(newSegs, Part{Style: seg.Style, String: right})
	}

	if segIdx+1 < len(tb.Text) {
		newSegs = append(newSegs, tb.Text[segIdx+1:]...)
	}

	tb.Text = MergeSameStyles(newSegs)
	return nil
}

// DeleteRange удаляет диапазон [start, end)
func (tb *Data) DeleteRange(start, end int) error {
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

	newSegs := make([]Part, 0, len(tb.Text))

	// Добавляем сегменты до segStart
	if segStart > 0 {
		newSegs = append(newSegs, tb.Text[:segStart]...)
	}

	// Обрабатываем начальный сегмент (если нужно)
	if segStart < len(tb.Text) && offStart > 0 {
		seg := tb.Text[segStart]
		r := []rune(seg.String)
		left := string(r[:offStart])
		if left != "" {
			newSegs = append(newSegs, Part{Style: seg.Style, String: left})
		}
	}

	// Обрабатываем конечный сегмент (если нужно)
	if segEnd < len(tb.Text) && offEnd > 0 {
		seg := tb.Text[segEnd]
		r := []rune(seg.String)
		right := string(r[offEnd:])
		if right != "" {
			newSegs = append(newSegs, Part{Style: seg.Style, String: right})
		}
	}

	// Добавляем сегменты после segEnd
	if segEnd+1 < len(tb.Text) {
		newSegs = append(newSegs, tb.Text[segEnd+1:]...)
	}

	tb.Text = MergeSameStyles(newSegs)
	return nil
}
