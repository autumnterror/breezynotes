package log

import (
	"errors"
	"testing"
)

func TestLog(t *testing.T) {
	Error("err test", "msg", nil)
	Error("err test", "msg", errors.New("error"))

	Warn("warn test", "msg", nil)
	Warn("warn test", "msg", errors.New("error"))

	Success("suc test", "msg")

	Info("info test", "msg")

	Red(1, 2, 3, 4, 5, 6, 7, 8)
	Green(1, 2, 3, 4, 5, 6, 7, 8)
	Blue(1, 2, 3, 4, 5, 6, 7, 8)
	Yellow(1, 2, 3, 4, 5, 6, 7, 8)

	Println(1, 2, 3, 4, 5, 6, 7, 8)
}
