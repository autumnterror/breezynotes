package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Move(s []string, old, new int) []string {
	if old < 0 || old >= len(s) {
		return nil
	}
	if new < 0 {
		new = 0
	} else if new > len(s) {
		new = len(s)
	}

	if old == new {
		return s
	}
	elem := s[old]

	sWO := append(s[:old], s[old+1:]...)

	s = append(sWO[:new], append([]string{elem}, sWO[new:]...)...)
	return s
}

type TestMove struct {
	inS   []string
	needS []string
	old   int
	new   int
	err   bool
}

func TestTest(t *testing.T) {
	// log.Blue(Move([]string{"0", "1"}, 1, 0))
	tts := []TestMove{
		{
			inS:   []string{"0", "1", "2", "3", "4"},
			needS: []string{"0", "2", "1", "3", "4"},
			old:   1,
			new:   2,
		},
		{
			inS:   []string{"0", "1"},
			needS: []string{"1", "0"},
			old:   0,
			new:   1,
		},
		{
			inS:   []string{"0", "1", "2", "3", "4"},
			needS: []string{"1", "0", "2", "3", "4"},
			old:   0,
			new:   1,
		},
		{
			inS:   []string{"0", "1", "2", "3", "4"},
			needS: []string{"0", "1", "2", "3", "4"},
			old:   0,
			new:   0,
		},
		{
			inS:   []string{"0", "1", "2", "3", "4"},
			needS: []string{"0", "1", "2", "3", "4"},
			old:   -1,
			new:   -1,
			err:   true,
		},
	}

	for _, tt := range tts {
		dst := make([]string, len(tt.inS))
		copy(dst, tt.inS)

		res := Move(tt.inS, tt.old, tt.new)
		if tt.err {
			assert.Equal(t, []string(nil), res)
		} else {
			assert.Equal(t, tt.needS, res)
		}
		res = Move(dst, tt.old, tt.new)
		if tt.err {
			assert.Equal(t, []string(nil), res)
		} else {
			assert.Equal(t, tt.needS, res)
		}
	}
}
