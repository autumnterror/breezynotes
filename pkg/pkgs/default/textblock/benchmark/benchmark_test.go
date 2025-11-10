package benchmark

import (
	"bufio"
	"github.com/autumnterror/breezynotes/pkg/pkgs/default/textblock"
	"log"
	"os"
	"testing"
	"time"
)

func TestApplyStylePerformance(t *testing.T) {
	t.Parallel()

	// читаем файл ./text.txt
	file, err := os.Open("./text.txt")
	if err != nil {
		t.Fatalf("failed to open text.txt: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes) // читаем посимвольно
	var text []textblock.TextData
	for scanner.Scan() {
		ch := scanner.Text()
		text = append(text, textblock.TextData{Style: "default", Text: ch})
	}

	tb := &textblock.TextBlock{Text: text}
	log.Printf("Initial segments of text: %d\n", len(tb.Text))

	startTime := time.Now()

	// операция 1: выделяем средний кусок и делаем bold
	//midStart := len(tb.Text) / 2
	//midEnd := midStart + 50
	err = tb.ApplyStyle(100000, 300000, "bold")
	if err != nil {
		t.Fatalf("ApplyStyle failed: %v", err)
	}
	log.Printf("After bold middle: %d segments. Time: %d ms\n", len(tb.Text), time.Since(startTime).Milliseconds())

	startTime2 := time.Now()
	// операция 2: вставляем в конец
	endStart := len(tb.Text)
	endEnd := len(tb.Text)
	err = tb.ApplyStyle(endStart-1, endEnd, "italic")
	if err != nil {
		t.Fatalf("ApplyStyle failed: %v", err)
	}
	log.Printf("After insert at end: %d segments. Time: %d ms\n", len(tb.Text), time.Since(startTime2).Milliseconds())
	startTime3 := time.Now()

	// операция 3: перекрываем первый кусок
	err = tb.ApplyStyle(0, 100000, "underline")
	if err != nil {
		t.Fatalf("ApplyStyle failed: %v", err)
	}
	log.Printf("After underline first 100 chars: %d segments. Time: %d ms\n", len(tb.Text), time.Since(startTime3).Milliseconds())
	duration := time.Since(startTime)
	log.Printf("Total time: %s\n", duration)
}
