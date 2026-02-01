package text

import (
	"github.com/autumnterror/utils_go/pkg/log"
	"testing"
)

func TestMap(t *testing.T) {
	tst := Data{Text: []Part{
		{
			Style:  "test",
			String: "test",
		},
	}}

	log.Green(tst.ToMap())
}
