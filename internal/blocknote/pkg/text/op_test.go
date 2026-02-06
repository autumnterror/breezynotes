package text

import (
	"fmt"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/autumnterror/utils_go/pkg/utils/format"
	"github.com/stretchr/testify/assert"
	"testing"
)

//----------------------InsertDeleteText-----------------

func TestInsertText(t *testing.T) {
	tests := []struct {
		name           string
		initialString  []Part
		pos            int
		newString      string
		expectedResult []Part
		description    string
	}{
		{
			name:           "insert on empty block",
			initialString:  []Part{},
			pos:            0,
			newString:      "Hello",
			expectedResult: []Part{{Style: "default", String: "Hello"}},
			description:    "Should create new segment with default style",
		},
		{
			name: "insert at beginning",
			initialString: []Part{
				{Style: "A", String: "World"},
			},
			pos:       0,
			newString: "Hello ",
			expectedResult: []Part{
				{Style: "A", String: "Hello World"},
			},
			description: "Should insert text at beginning and preserve existing style",
		},
		{
			name: "insert at end",
			initialString: []Part{
				{Style: "A", String: "Hello"},
			},
			pos:       5,
			newString: " World",
			expectedResult: []Part{
				{Style: "A", String: "Hello World"},
			},
			description: "Should append text at end with A style",
		},
		{
			name: "insert in middle of segment",
			initialString: []Part{
				{Style: "A", String: "HelloWorld"},
			},
			pos:       5,
			newString: " ",
			expectedResult: []Part{
				{Style: "A", String: "Hello World"},
			},
			description: "Should split segment and insert new text with A style",
		},
		{
			name: "insert at boundary between segments",
			initialString: []Part{
				{Style: "A", String: "Hello"},
				{Style: "B", String: "World"},
			},
			pos:       5,
			newString: " ",
			expectedResult: []Part{
				{Style: "A", String: "Hello "},
				{Style: "B", String: "World"},
			},
			description: "Should insert between segments without splitting them",
		},
		{
			name: "insert with position out of bounds",
			initialString: []Part{
				{Style: "A", String: "Hello"},
			},
			pos:       100,
			newString: " World",
			expectedResult: []Part{
				{Style: "A", String: "Hello World"},
			},
			description: "Should clamp position to end and append text",
		},
		{
			name: "insert with negative position",
			initialString: []Part{
				{Style: "A", String: "World"},
			},
			pos:       -10,
			newString: "Hello ",
			expectedResult: []Part{
				{Style: "A", String: "Hello World"},
			},
			description: "Should clamp position to start and insert at beginning",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log.Printf("=== TEST: %s ===", test.name)
			log.Printf("Description: %s", test.description)
			log.Printf("Initial text: %+v", test.initialString)
			log.Printf("Inserting '%s' at position %d", test.newString, test.pos)

			tb := &Data{Text: test.initialString}
			err := tb.InsertText(test.pos, test.newString)

			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, tb.Text, "Result doesn't match expected")

			log.Printf("Result: %+v", tb.Text)
			log.Printf("=== TEST PASSED ===\n")
		})
	}
}

// TestDeleteRange tests various cases of DeleteRange function
func TestDeleteRange(t *testing.T) {
	tests := []struct {
		name           string
		initialString  []Part
		start          int
		end            int
		expectedResult []Part
		description    string
	}{
		{
			name:           "del from empty block",
			initialString:  []Part{},
			start:          0,
			end:            5,
			expectedResult: []Part{},
			description:    "Should do nothing on empty block",
		},
		{
			name: "delete entire single segment",
			initialString: []Part{
				{Style: "A", String: "Hello"},
			},
			start:          0,
			end:            5,
			expectedResult: []Part{},
			description:    "Should remove the entire segment",
		},
		{
			name: "delete beginning of segment",
			initialString: []Part{
				{Style: "A", String: "HelloWorld"},
			},
			start: 0,
			end:   5,
			expectedResult: []Part{
				{Style: "A", String: "World"},
			},
			description: "Should remove beginning of segment",
		},
		{
			name: "delete end of segment",
			initialString: []Part{
				{Style: "A", String: "HelloWorld"},
			},
			start: 5,
			end:   10,
			expectedResult: []Part{
				{Style: "A", String: "Hello"},
			},
			description: "Should remove end of segment",
		},
		{
			name: "delete middle of segment",
			initialString: []Part{
				{Style: "A", String: "HelloWorld"},
			},
			start: 3,
			end:   7,
			expectedResult: []Part{
				{Style: "A", String: "Helrld"},
			},
			description: "Should remove middle part of segment",
		},
		{
			name: "delete across multiple segments",
			initialString: []Part{
				{Style: "A", String: "Hello"},
				{Style: "B", String: " "},
				{Style: "C", String: "World"},
			},
			start: 3,
			end:   7,
			expectedResult: []Part{
				{Style: "A", String: "Hel"},
				{Style: "C", String: "orld"},
			},
			description: "Should remove parts across multiple segments",
		},
		{
			name: "delete with range out of bounds",
			initialString: []Part{
				{Style: "A", String: "Hello"},
			},
			start:          -5,
			end:            100,
			expectedResult: []Part{},
			description:    "Should clamp range and remove everything",
		},
		{
			name: "delete with invalid range",
			initialString: []Part{
				{Style: "A", String: "Hello"},
			},
			start: 7,
			end:   3,
			expectedResult: []Part{
				{Style: "A", String: "Hello"},
			},
			description: "Should do nothing with invalid range (start > end)",
		},
		{
			name: "delete empty range",
			initialString: []Part{
				{Style: "A", String: "Hello"},
			},
			start: 3,
			end:   3,
			expectedResult: []Part{
				{Style: "A", String: "Hello"},
			},
			description: "Should do nothing with empty range",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log.Printf("=== TEST: %s ===", test.name)
			log.Printf("Description: %s", test.description)
			log.Printf("Initial text: %+v", test.initialString)
			log.Printf("Deleting range [%d, %d)", test.start, test.end)
			tb := &Data{Text: test.initialString}
			err := tb.DeleteRange(test.start, test.end)

			assert.NoError(t, err)
			assert.Equal(t, test.expectedResult, tb.Text, "Result doesn't match expected")

			log.Printf("Result: %+v", tb.Text)
			log.Printf("=== TEST PASSED ===\n")
		})
	}
}

// TestIntegration tests combined operations
func TestIntegration(t *testing.T) {
	t.Run("insert and delete operations", func(t *testing.T) {
		log.Printf("=== INTEGRATION TEST ===")
		tb := &Data{Text: []Part{}}

		// insert into empty block
		err := tb.InsertText(0, "Hello World")
		assert.NoError(t, err)
		log.Printf("After initial insert: %+v", tb.Text)

		// insert in the middle
		err = tb.InsertText(5, " Beautiful")
		assert.NoError(t, err)
		log.Printf("After middle insert: %+v", tb.Text)

		// delete part of the text
		err = tb.DeleteRange(0, 6)
		assert.NoError(t, err)
		log.Printf("After deletion: %+v", tb.Text)

		// Apply style to part of the text
		err = tb.ApplyStyle(0, 5, "bold")
		assert.NoError(t, err)
		log.Printf("After applying style: %+v", tb.Text)

		expected := []Part{
			{Style: "bold", String: "Beaut"},
			{Style: "default", String: "iful World"},
		}
		assert.Equal(t, expected, tb.Text, "Integration test failed")

		log.Printf("Final result: %+v", tb.Text)
		log.Printf("=== INTEGRATION TEST PASSED ===\n")
	})
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("Unicode characters", func(t *testing.T) {
		log.Printf("=== UNICODE TEST ===")

		tb := &Data{Text: []Part{
			{Style: "A", String: "Привет"},
			{Style: "B", String: "Мир"},
		}}

		// insert in the middle of Unicode text
		err := tb.InsertText(3, " Крутой")
		assert.NoError(t, err)
		log.Printf("After Unicode insert: %+v", tb.Text)

		// delete part of Unicode text
		err = tb.DeleteRange(0, 7)
		assert.NoError(t, err)
		log.Printf("After Unicode delete: %+v", tb.Text)

		log.Printf("=== UNICODE TEST PASSED ===\n")
	})

	t.Run("Empty operations", func(t *testing.T) {
		log.Printf("=== EMPTY OPERATIONS TEST ===")

		tb := &Data{Text: []Part{
			{Style: "A", String: "Hello"},
		}}

		// insert empty text
		err := tb.InsertText(2, "")
		assert.NoError(t, err)
		log.Printf("After empty insert: %+v", tb.Text)

		// delete empty range
		err = tb.DeleteRange(2, 2)
		assert.NoError(t, err)
		log.Printf("After empty delete: %+v", tb.Text)

		// Should be unchanged
		expected := []Part{
			{Style: "A", String: "Hello"},
		}
		assert.Equal(t, expected, tb.Text, "Empty operations test failed")

		log.Printf("=== EMPTY OPERATIONS TEST PASSED ===\n")
	})
}

//----------------------ApplyStyle------------------------

func TestApplyStyle(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		initial  []Part
		start    int
		end      int
		style    string
		expected []Part
	}{
		{
			name: "simple middle range",
			initial: []Part{
				{Style: "default", String: "hello world"},
			},
			start: 6,
			end:   11,
			style: "bold",
			expected: []Part{
				{Style: "default", String: "hello "},
				{Style: "bold", String: "world"},
			},
		},
		{
			name: "apply at beginning",
			initial: []Part{
				{Style: "default", String: "hello world"},
			},
			start: 0,
			end:   5,
			style: "italic",
			expected: []Part{
				{Style: "italic", String: "hello"},
				{Style: "default", String: " world"},
			},
		},
		{
			name: "apply at end",
			initial: []Part{
				{Style: "default", String: "hello world"},
			},
			start: 6,
			end:   11,
			style: "italic",
			expected: []Part{
				{Style: "default", String: "hello "},
				{Style: "italic", String: "world"},
			},
		},
		{
			name: "nested style inside existing style",
			initial: []Part{
				{Style: "bold", String: "hello world"},
			},
			start: 6,
			end:   11,
			style: "italic",
			expected: []Part{
				{Style: "bold", String: "hello "},
				{Style: "italic", String: "world"},
			},
		},
		{
			name: "range beyond length should trim",
			initial: []Part{
				{Style: "default", String: "short"},
			},
			start: 2,
			end:   100,
			style: "bold",
			expected: []Part{
				{Style: "default", String: "sh"},
				{Style: "bold", String: "ort"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := &Data{Text: append([]Part{}, tt.initial...)}
			assert.NoError(t, block.ApplyStyle(tt.start, tt.end, tt.style))
			fmt.Printf("After ApplyStyle(%d,%d,%s): %+v\n", tt.start, tt.end, tt.style, block.Text)
			assert.Equal(t, tt.expected, block.Text)
		})
	}
}

func TestApplyStyleInvalidRanges(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		initial []Part
		start   int
		end     int
		style   string
		wantErr bool
	}{
		{
			name:    "start >= end",
			initial: []Part{{Style: "default", String: "hello"}},
			start:   3,
			end:     2,
			style:   "bold",
			wantErr: true,
		},
		{
			name:    "empty text",
			initial: nil,
			start:   0,
			end:     1,
			style:   "bold",
			wantErr: false, // просто возвращаем nil
		},
		{
			name:    "end > total length",
			initial: []Part{{Style: "default", String: "abc"}},
			start:   1,
			end:     10,
			style:   "italic",
			wantErr: false, // обрежет до длины текста
		},
		{
			name:    "negative start",
			initial: []Part{{Style: "default", String: "abc"}},
			start:   -5,
			end:     2,
			style:   "italic",
			wantErr: false, // скорректирует start = 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &Data{Text: tt.initial}
			err := tb.ApplyStyle(tt.start, tt.end, tt.style)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestApplyStyleMultipleOps(t *testing.T) {
	t.Parallel()
	block := &Data{
		Text: []Part{
			{Style: "default", String: "hello world"},
		},
	}

	// Ожидания после каждой операции
	expectedSteps := [][]Part{
		{
			{Style: "default", String: "hello "},
			{Style: "bold", String: "world"},
		},
		{
			{Style: "default", String: "he"},
			{Style: "italic", String: "llo wo"},
			{Style: "bold", String: "rld"},
		},
		{
			{Style: "underline", String: "he"},
			{Style: "italic", String: "llo wo"},
			{Style: "bold", String: "rld"},
		},
		{
			{Style: "underline", String: "he"},
			{Style: "italic", String: "llo"},
			{Style: "bold", String: " world"},
		},
	}

	ops := []struct {
		start int
		end   int
		style string
	}{
		{6, 11, "bold"},     // выделяем "world"
		{2, 8, "italic"},    // выделяем "llo wo"
		{0, 2, "underline"}, // выделяем "he"
		{5, 8, "bold"},      // еще раз кусок внутри italic ("wo")
	}

	for i, op := range ops {
		assert.NoError(t, block.ApplyStyle(op.start, op.end, op.style))
		fmt.Printf("Step %d: ApplyStyle(%d,%d,%s) => %+v\n",
			i+1, op.start, op.end, op.style, block.Text)
		log.Println(block.Text)
		assert.Equal(t, expectedSteps[i], block.Text, fmt.Sprintf("step %d failed", i+1))
	}
}

func TestApplyStyleMultipleComplex(t *testing.T) {
	t.Parallel()
	block := &Data{
		Text: []Part{
			{Style: "default", String: "hello world!"},
		},
	}

	ops := []struct {
		start int
		end   int
		style string
		desc  string
	}{
		{0, 5, "bold", "apply bold to 'hello'"},
		{6, 11, "italic", "apply italic to 'world'"},
		{3, 8, "underline", "apply underline across segments 'lo wo'"},
		{0, 12, "default", "reset entire text to default"},
		{0, 1, "bold", "apply bold to first char"},
		{11, 12, "italic", "apply italic to last char"},
	}

	expectedSteps := []interface{}{
		[]Part{
			{Style: "bold", String: "hello"},
			{Style: "default", String: " world!"},
		},
		[]Part{
			{Style: "bold", String: "hello"},
			{Style: "default", String: " "},
			{Style: "italic", String: "world"},
			{Style: "default", String: "!"},
		},
		[]Part{
			{Style: "bold", String: "hel"},
			{Style: "underline", String: "lo wo"},
			{Style: "italic", String: "rld"},
			{Style: "default", String: "!"},
		},
		[]Part{
			{Style: "default", String: "hello world!"},
		},
		[]Part{
			{Style: "bold", String: "h"},
			{Style: "default", String: "ello world!"},
		},
		[]Part{
			{Style: "bold", String: "h"},
			{Style: "default", String: "ello world"},
			{Style: "italic", String: "!"},
		},
	}
	fmt.Printf("start: \n```%s```\n", format.Struct(block.Text))
	for i, op := range ops {
		err := block.ApplyStyle(op.start, op.end, op.style)
		assert.NoError(t, err, fmt.Sprintf("operation %d failed: %s", i, op.desc))
		fmt.Printf("ApplyStyle(%d,%d,%s) => \n```%s```\n",
			op.start, op.end, op.style, format.Struct(block.Text))
		assert.Equal(t, expectedSteps[i], block.Text, fmt.Sprintf("step %d mismatch", i+1))
	}
}
