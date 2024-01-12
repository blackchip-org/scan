package scan

import (
	"fmt"
	"strings"
	"testing"
)

func TestMarkReader(t *testing.T) {
	tests := []struct {
		in      string
		initial int
		mark    int
		reset   bool
		out     string
	}{
		{"123456789", 0, 0, true, "123456789"},
		{"123456789", 0, 3, true, "123456789"},
		{"123456789", 0, 3, false, "456789"},
		{"123456789", 3, 3, true, "456789"},
		{"123456789", 3, 30, true, "456789"},
		{"123456789", 3, 3, false, "789"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			r := NewMarkReader(strings.NewReader(test.in))
			for i := 0; i < test.initial; i++ {
				r.Read()
			}
			r.Mark()
			for i := 0; i < test.mark; i++ {
				r.Read()
			}
			if test.reset {
				r.Reset()
			}
			have, err := r.ReadAll()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(have) != test.out {
				t.Errorf("\n have: %v \n want: %v", string(have), test.out)
			}
		})
	}
}
