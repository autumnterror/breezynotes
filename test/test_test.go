package test

import (
	"fmt"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"slices"
	"strings"
	"testing"
	"time"

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

func TestErrors(t *testing.T) {
	err := format.Error("op", status.Error(codes.FailedPrecondition, "test"))

	s, ok := status.FromError(err)
	if assert.True(t, ok) {
		log.Println(s.Code())
	}
}

type Test struct {
	A int
	B int
}

func TestArray(t *testing.T) {
	ta := []Test{{1, 2}, {3, 4}, {5, 6}}

	idx := slices.IndexFunc(ta, func(tt Test) bool {
		return tt.A == 3
	})
	assert.Equal(t, idx, 1)

	ta = append(ta[:idx], ta[idx+1:]...)
	log.Println(ta)
	assert.Equal(t, 2, len(ta))
	idx = slices.IndexFunc(ta, func(tt Test) bool {
		return tt.B == 2
	})
	assert.Equal(t, idx, 0)

	idx = slices.IndexFunc(ta, func(tt Test) bool {
		return tt.B == 33
	})
	assert.Equal(t, idx, -1)

}

func TestTestTest(t *testing.T) {
	pr := "ЫыьЬ HhYy"
	log.Blue("", strings.ToLower(pr))
}

func TestBCrypt(t *testing.T) {
	password := []byte("some-password")
	for cost := 8; cost <= 14; cost++ {
		start := time.Now()
		_, err := bcrypt.GenerateFromPassword(password, cost)
		if err != nil {
			panic(err)
		}
		elapsed := time.Since(start)
		fmt.Printf("cost=%d -> %v\n", cost, elapsed)
	}
}
