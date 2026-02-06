package benchmark

import (
	"bufio"
	"fmt"
	text2 "github.com/autumnterror/breezynotes/internal/blocknote/pkg/text"
	"log"
	"os"
	"testing"
	"time"
)

func TestApplyStylePerformance(t *testing.T) {
	t.Parallel()

	// читаем файл ./text.txt

	log.Println("start read file")
	file, err := os.Open("./text.txt")
	if err != nil {
		t.Fatalf("failed to open text.txt: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	count := 0
	var lenLine []int

	var text string
	for scanner.Scan() {
		count++
		ch := scanner.Text()
		text += ch
		lenLine = append(lenLine, len(ch))
	}

	tb := &text2.Data{Text: []text2.Part{
		{Style: "default", String: text},
	}}
	sum := 0
	countWithoutLow := 0
	for _, i := range lenLine {
		if i > 5 {
			countWithoutLow++
			sum += i
		}
	}

	count *= sum / len(lenLine)
	log.Printf("Initial segments of text: %d\n", count)

	//startTimeNormal := time.Now()

	//log.Println("normalization...")
	//tb.Data = textblock.MergeSameStylesParallel(tb.Data, runtime.NumCPU())
	//log.Printf("normalization success time: %d ms", time.Since(startTimeNormal).Milliseconds())

	startTime := time.Now()

	// операция 1: выделяем средний кусок и делаем bold
	//midStart := len(tb.Data) / 2
	//midEnd := midStart + 50
	log.Printf("tb.ApplyStyle(%d,%d, \"bold\"\n", count-2*count/3, count-count/3)
	err = tb.ApplyStyle(count-2*count/3, count-count/3, "bold")
	if err != nil {
		t.Fatalf("ApplyStyle failed: %v", err)
	}
	log.Printf("After bold middle Time: %d ms\n", time.Since(startTime).Milliseconds())
	for o, i := range tb.Text {
		fmt.Printf("%d: %s", o, i.Style)
	}
	fmt.Println()
	startTime2 := time.Now()

	// операция 2: вставляем в конец
	endStart := count - count/3
	endEnd := count
	log.Printf("tb.ApplyStyle(%d,%d, \"italic\"\n", endStart-1, endEnd)
	err = tb.ApplyStyle(endStart, endEnd, "italic")
	if err != nil {
		t.Fatalf("ApplyStyle failed: %v", err)
	}
	log.Printf("After insert at end Time: %d ms\n", time.Since(startTime2).Milliseconds())
	startTime3 := time.Now()
	for o, i := range tb.Text {
		fmt.Printf("%d: %s", o, i.Style)
	}
	fmt.Println()
	// операция 3: перекрываем первый кусок
	err = tb.ApplyStyle(0, count/3, "underline")
	log.Printf("tb.ApplyStyle(%d,%d, \"underline\"\n", 0, count/4)
	if err != nil {
		t.Fatalf("ApplyStyle failed: %v", err)
	}
	for o, i := range tb.Text {
		fmt.Printf("%d: %s", o, i.Style)
	}
	fmt.Println()
	log.Printf("After underline at start Time: %d ms\n", time.Since(startTime3).Milliseconds())

	startTime5 := time.Now()
	err = tb.InsertText(0, tb.PlainText())
	log.Printf("tb.InsertText(0, tb.PlainText())")
	if err != nil {
		t.Fatalf("ApplyStyle failed: %v", err)
	}
	for o, i := range tb.Text {
		fmt.Printf("%d: %s", o, i.Style)
	}
	fmt.Println()
	log.Printf("After underline at start Time: %d ms\n", time.Since(startTime5).Milliseconds())

	startTime4 := time.Now()
	err = tb.DeleteRange(0, count/3)
	log.Printf("tb.DeleteRange(0, count/3)")
	if err != nil {
		t.Fatalf("ApplyStyle failed: %v", err)
	}
	for o, i := range tb.Text {
		fmt.Printf("%d: %s", o, i.Style)
	}
	fmt.Println()
	log.Printf("After underline at start Time: %d ms\n", time.Since(startTime4).Milliseconds())

	duration := time.Since(startTime)
	log.Printf("Total time: %s\n", duration)
}
