package format

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestError(t *testing.T) {
	assert.Error(t, Error("", errors.New("")))
	assert.NoError(t, Error("", nil))
}
func TestFormat(t *testing.T) {
	a := map[string]int{"a": 1, "b": 2}

	log.Println(Struct(a))
	log.Println(Struct(struct {
		A string
		B int
	}{A: "1", B: 2}))
	log.Println(Struct(nil))
}
