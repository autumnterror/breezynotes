package text

import (
	"sort"
	"sync"
)

func MergeSameStyles(segments []Part) []Part {
	if len(segments) == 0 {
		return segments
	}
	out := make([]Part, 0, len(segments))
	cur := segments[0]
	for i := 1; i < len(segments); i++ {
		if segments[i].Style == cur.Style {
			cur.String += segments[i].String
		} else {
			if cur.String != "" {
				out = append(out, cur)
			}
			cur = segments[i]
		}
	}
	if cur.String != "" {
		out = append(out, cur)
	}
	return out
}

func MergeSameStylesParallel(segments []Part, workers int) []Part {
	if len(segments) == 0 {
		return segments
	}

	if len(segments) < 1000 || workers <= 1 {
		return MergeSameStyles(segments)
	}

	batchSize := len(segments) / workers
	results := make([][]Part, workers)
	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()

			start := worker * batchSize
			end := start + batchSize
			if worker == workers-1 {
				end = len(segments)
			}

			results[worker] = MergeSameStyles(segments[start:end])
		}(w)
	}
	wg.Wait()

	// Объединяем результаты
	var totalLen int
	for _, r := range results {
		totalLen += len(r)
	}

	out := make([]Part, 0, totalLen)
	for _, r := range results {
		out = append(out, r...)
	}

	// Финальное слияние границ между батчами
	return MergeSameStyles(out)
}

// buildPrefixLens Вычисляет префиксные суммы длин сегментов.
//
//	segments := []Part{
//	   {Part: "Hello", Style: "A"},
//	   {Part: "World", Style: "B"},
//	}
//
// pref, total := buildPrefixLens(segments)
//
// //pref = [5, 10], total = 10
func buildPrefixLens(segments []Part) ([]int, int) {
	n := len(segments)
	pref := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		l := len([]rune(segments[i].String))
		sum += l
		pref[i] = sum
	}
	return pref, sum
}

// findSegmentByPos Находит индекс сегмента и смещение в нем для заданной позиции.
//
//	segments := []Part{
//	   {Part: "Hello", Style: "A"},
//	   {Part: "World", Style: "B"},
//	}
//
// pref := []int{5, 10}
// segIdx, offset := findSegmentByPos(segments, pref, 7)
//
// // segIdx = 1, offset = 2 (позиция 7 - это 2-й символ в "World")
func findSegmentByPos(segments []Part, pref []int, pos int) (segIdx int, offset int) {
	n := len(segments)
	if n == 0 {
		return 0, 0
	}
	total := pref[n-1]
	if pos <= 0 {
		return 0, 0
	}
	if pos >= total {
		return n, 0
	}
	// ищем первый i с pref[i] > pos  => эквивалентно pref[i] >= pos+1
	i := sort.SearchInts(pref, pos+1)
	segIdx = i
	var prevSum int
	if segIdx > 0 {
		prevSum = pref[segIdx-1]
	}
	offset = pos - prevSum
	return segIdx, offset
}
