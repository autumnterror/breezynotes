package validate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPw(t *testing.T) {
	type test struct {
		pw   string
		res  bool
		name string
	}
	tests := []test{
		{
			pw:   "Abc1!",
			res:  true,
			name: "true",
		},
		{
			pw:   "Password123!",
			res:  true,
			name: "true",
		},
		{
			pw:   "abc1!",
			res:  false,
			name: "no upper",
		},
		{
			pw:   "Abcde!",
			res:  false,
			name: "no digit",
		},
		{
			pw:   "Abc1",
			res:  false,
			name: "short",
		},
		{
			pw:   "Abc12",
			res:  false,
			name: "no symbol",
		},
		{
			pw:   "12345",
			res:  false,
			name: "no letter",
		},
		{
			pw:   "Ab!456789012345678901",
			res:  false,
			name: "too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.True(t, Password(tt.pw) == tt.res)
		})
	}
}
