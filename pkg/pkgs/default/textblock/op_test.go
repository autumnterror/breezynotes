package textblock

import (
	"fmt"
	"github.com/autumnterror/breezynotes/pkg/log"
	"github.com/autumnterror/breezynotes/pkg/utils/format"
	"github.com/stretchr/testify/assert"
	"testing"
)

//----------------------InsertDeleteText-----------------

func TestInsertText(t *testing.T) {
	tests := []struct {
		name           string
		initialText    []TextData
		pos            int
		newText        string
		expectedResult []TextData
		description    string
	}{
		{
			name:           "insert on empty block",
			initialText:    []TextData{},
			pos:            0,
			newText:        "Hello",
			expectedResult: []TextData{{Style: "default", Text: "Hello"}},
			description:    "Should create new segment with default style",
		},
		{
			name: "Insert at beginning",
			initialText: []TextData{
				{Style: "A", Text: "World"},
			},
			pos:     0,
			newText: "Hello ",
			expectedResult: []TextData{
				{Style: "default", Text: "Hello "},
				{Style: "A", Text: "World"},
			},
			description: "Should insert text at beginning and preserve existing style",
		},
		{
			name: "Insert at end",
			initialText: []TextData{
				{Style: "A", Text: "Hello"},
			},
			pos:     5,
			newText: " World",
			expectedResult: []TextData{
				{Style: "A", Text: "Hello"},
				{Style: "default", Text: " World"},
			},
			description: "Should append text at end with default style",
		},
		{
			name: "Insert in middle of segment",
			initialText: []TextData{
				{Style: "A", Text: "HelloWorld"},
			},
			pos:     5,
			newText: " ",
			expectedResult: []TextData{
				{Style: "A", Text: "Hello"},
				{Style: "default", Text: " "},
				{Style: "A", Text: "World"},
			},
			description: "Should split segment and insert new text with default style",
		},
		{
			name: "Insert at boundary between segments",
			initialText: []TextData{
				{Style: "A", Text: "Hello"},
				{Style: "B", Text: "World"},
			},
			pos:     5,
			newText: " ",
			expectedResult: []TextData{
				{Style: "A", Text: "Hello"},
				{Style: "default", Text: " "},
				{Style: "B", Text: "World"},
			},
			description: "Should insert between segments without splitting them",
		},
		{
			name: "Insert with position out of bounds",
			initialText: []TextData{
				{Style: "A", Text: "Hello"},
			},
			pos:     100,
			newText: " World",
			expectedResult: []TextData{
				{Style: "A", Text: "Hello"},
				{Style: "default", Text: " World"},
			},
			description: "Should clamp position to end and append text",
		},
		{
			name: "Insert with negative position",
			initialText: []TextData{
				{Style: "A", Text: "World"},
			},
			pos:     -10,
			newText: "Hello ",
			expectedResult: []TextData{
				{Style: "default", Text: "Hello "},
				{Style: "A", Text: "World"},
			},
			description: "Should clamp position to start and insert at beginning",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log.Printf("=== TEST: %s ===", test.name)
			log.Printf("Description: %s", test.description)
			log.Printf("Initial text: %+v", test.initialText)
			log.Printf("Inserting '%s' at position %d", test.newText, test.pos)

			tb := &TextBlock{Text: test.initialText}
			err := tb.InsertText(test.pos, test.newText)

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
		initialText    []TextData
		start          int
		end            int
		expectedResult []TextData
		description    string
	}{
		{
			name:           "del from empty block",
			initialText:    []TextData{},
			start:          0,
			end:            5,
			expectedResult: []TextData{},
			description:    "Should do nothing on empty block",
		},
		{
			name: "Delete entire single segment",
			initialText: []TextData{
				{Style: "A", Text: "Hello"},
			},
			start:          0,
			end:            5,
			expectedResult: []TextData{},
			description:    "Should remove the entire segment",
		},
		{
			name: "Delete beginning of segment",
			initialText: []TextData{
				{Style: "A", Text: "HelloWorld"},
			},
			start: 0,
			end:   5,
			expectedResult: []TextData{
				{Style: "A", Text: "World"},
			},
			description: "Should remove beginning of segment",
		},
		{
			name: "Delete end of segment",
			initialText: []TextData{
				{Style: "A", Text: "HelloWorld"},
			},
			start: 5,
			end:   10,
			expectedResult: []TextData{
				{Style: "A", Text: "Hello"},
			},
			description: "Should remove end of segment",
		},
		{
			name: "Delete middle of segment",
			initialText: []TextData{
				{Style: "A", Text: "HelloWorld"},
			},
			start: 3,
			end:   7,
			expectedResult: []TextData{
				{Style: "A", Text: "Helrld"},
			},
			description: "Should remove middle part of segment",
		},
		{
			name: "Delete across multiple segments",
			initialText: []TextData{
				{Style: "A", Text: "Hello"},
				{Style: "B", Text: " "},
				{Style: "C", Text: "World"},
			},
			start: 3,
			end:   7,
			expectedResult: []TextData{
				{Style: "A", Text: "Hel"},
				{Style: "C", Text: "orld"},
			},
			description: "Should remove parts across multiple segments",
		},
		{
			name: "Delete with range out of bounds",
			initialText: []TextData{
				{Style: "A", Text: "Hello"},
			},
			start:          -5,
			end:            100,
			expectedResult: []TextData{},
			description:    "Should clamp range and remove everything",
		},
		{
			name: "Delete with invalid range",
			initialText: []TextData{
				{Style: "A", Text: "Hello"},
			},
			start: 7,
			end:   3,
			expectedResult: []TextData{
				{Style: "A", Text: "Hello"},
			},
			description: "Should do nothing with invalid range (start > end)",
		},
		{
			name: "Delete empty range",
			initialText: []TextData{
				{Style: "A", Text: "Hello"},
			},
			start: 3,
			end:   3,
			expectedResult: []TextData{
				{Style: "A", Text: "Hello"},
			},
			description: "Should do nothing with empty range",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log.Printf("=== TEST: %s ===", test.name)
			log.Printf("Description: %s", test.description)
			log.Printf("Initial text: %+v", test.initialText)
			log.Printf("Deleting range [%d, %d)", test.start, test.end)

			tb := &TextBlock{Text: test.initialText}
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
	t.Run("Insert and delete operations", func(t *testing.T) {
		log.Printf("=== INTEGRATION TEST ===")

		tb := &TextBlock{Text: []TextData{}}

		// Insert into empty block
		err := tb.InsertText(0, "Hello World")
		assert.NoError(t, err)
		log.Printf("After initial insert: %+v", tb.Text)

		// Insert in the middle
		err = tb.InsertText(5, " Beautiful")
		assert.NoError(t, err)
		log.Printf("After middle insert: %+v", tb.Text)

		// Delete part of the text
		err = tb.DeleteRange(0, 6)
		assert.NoError(t, err)
		log.Printf("After deletion: %+v", tb.Text)

		// Apply style to part of the text
		err = tb.ApplyStyle(0, 5, "bold")
		assert.NoError(t, err)
		log.Printf("After applying style: %+v", tb.Text)

		expected := []TextData{
			{Style: "bold", Text: "Beaut"},
			{Style: "default", Text: "iful World"},
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

		tb := &TextBlock{Text: []TextData{
			{Style: "A", Text: "Привет"},
			{Style: "B", Text: "Мир"},
		}}

		// Insert in the middle of Unicode text
		err := tb.InsertText(3, " Крутой")
		assert.NoError(t, err)
		log.Printf("After Unicode insert: %+v", tb.Text)

		// Delete part of Unicode text
		err = tb.DeleteRange(0, 7)
		assert.NoError(t, err)
		log.Printf("After Unicode delete: %+v", tb.Text)

		log.Printf("=== UNICODE TEST PASSED ===\n")
	})

	t.Run("Empty operations", func(t *testing.T) {
		log.Printf("=== EMPTY OPERATIONS TEST ===")

		tb := &TextBlock{Text: []TextData{
			{Style: "A", Text: "Hello"},
		}}

		// Insert empty text
		err := tb.InsertText(2, "")
		assert.NoError(t, err)
		log.Printf("After empty insert: %+v", tb.Text)

		// Delete empty range
		err = tb.DeleteRange(2, 2)
		assert.NoError(t, err)
		log.Printf("After empty delete: %+v", tb.Text)

		// Should be unchanged
		expected := []TextData{
			{Style: "A", Text: "Hello"},
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
		initial  []TextData
		start    int
		end      int
		style    string
		expected []TextData
	}{
		{
			name: "simple middle range",
			initial: []TextData{
				{Style: "default", Text: "hello world"},
			},
			start: 6,
			end:   11,
			style: "bold",
			expected: []TextData{
				{Style: "default", Text: "hello "},
				{Style: "bold", Text: "world"},
			},
		},
		{
			name: "apply at beginning",
			initial: []TextData{
				{Style: "default", Text: "hello world"},
			},
			start: 0,
			end:   5,
			style: "italic",
			expected: []TextData{
				{Style: "italic", Text: "hello"},
				{Style: "default", Text: " world"},
			},
		},
		{
			name: "apply at end",
			initial: []TextData{
				{Style: "default", Text: "hello world"},
			},
			start: 6,
			end:   11,
			style: "italic",
			expected: []TextData{
				{Style: "default", Text: "hello "},
				{Style: "italic", Text: "world"},
			},
		},
		{
			name: "nested style inside existing style",
			initial: []TextData{
				{Style: "bold", Text: "hello world"},
			},
			start: 6,
			end:   11,
			style: "italic",
			expected: []TextData{
				{Style: "bold", Text: "hello "},
				{Style: "italic", Text: "world"},
			},
		},
		{
			name: "range beyond length should trim",
			initial: []TextData{
				{Style: "default", Text: "short"},
			},
			start: 2,
			end:   100,
			style: "bold",
			expected: []TextData{
				{Style: "default", Text: "sh"},
				{Style: "bold", Text: "ort"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := &TextBlock{Text: append([]TextData{}, tt.initial...)}
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
		initial []TextData
		start   int
		end     int
		style   string
		wantErr bool
	}{
		{
			name:    "start >= end",
			initial: []TextData{{Style: "default", Text: "hello"}},
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
			initial: []TextData{{Style: "default", Text: "abc"}},
			start:   1,
			end:     10,
			style:   "italic",
			wantErr: false, // обрежет до длины текста
		},
		{
			name:    "negative start",
			initial: []TextData{{Style: "default", Text: "abc"}},
			start:   -5,
			end:     2,
			style:   "italic",
			wantErr: false, // скорректирует start = 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &TextBlock{Text: tt.initial}
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
	block := &TextBlock{
		Text: []TextData{
			{Style: "default", Text: "hello world"},
		},
	}

	// Ожидания после каждой операции
	expectedSteps := [][]TextData{
		{
			{Style: "default", Text: "hello "},
			{Style: "bold", Text: "world"},
		},
		{
			{Style: "default", Text: "he"},
			{Style: "italic", Text: "llo wo"},
			{Style: "bold", Text: "rld"},
		},
		{
			{Style: "underline", Text: "he"},
			{Style: "italic", Text: "llo wo"},
			{Style: "bold", Text: "rld"},
		},
		{
			{Style: "underline", Text: "he"},
			{Style: "italic", Text: "llo"},
			{Style: "bold", Text: " world"},
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
	block := &TextBlock{
		Text: []TextData{
			{Style: "default", Text: "hello world!"},
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
		[]TextData{
			{Style: "bold", Text: "hello"},
			{Style: "default", Text: " world!"},
		},
		[]TextData{
			{Style: "bold", Text: "hello"},
			{Style: "default", Text: " "},
			{Style: "italic", Text: "world"},
			{Style: "default", Text: "!"},
		},
		[]TextData{
			{Style: "bold", Text: "hel"},
			{Style: "underline", Text: "lo wo"},
			{Style: "italic", Text: "rld"},
			{Style: "default", Text: "!"},
		},
		[]TextData{
			{Style: "default", Text: "hello world!"},
		},
		[]TextData{
			{Style: "bold", Text: "h"},
			{Style: "default", Text: "ello world!"},
		},
		[]TextData{
			{Style: "bold", Text: "h"},
			{Style: "default", Text: "ello world"},
			{Style: "italic", Text: "!"},
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
