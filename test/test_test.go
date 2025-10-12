package test

import (
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"testing"
)

type b struct {
	A int
	B string
}

func Foo() interface{} {
	return b{
		A: 12,
		B: "123",
	}
}

func TestTest(t *testing.T) {
	res := Foo()
	log.Println(format.Struct(res.(b)))
}
