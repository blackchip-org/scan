package scan

import (
	"testing"
)

func TestKeep(t *testing.T) {
	src := "1234567"
	tests := []struct {
		n     int
		value string
	}{
		{0, ""},
		{1, "1"},
		{3, "123"},
		{8, "1234567"},
	}

	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			s := NewFromString("", src)
			Repeat(s.Keep, test.n)
			have := s.Emit()
			if have.Value != test.value {
				t.Errorf("\n have: [%v] \n want: [%v]", have.Value, test.value)
			}
			if have.Literal != "" {
				t.Errorf("\n have: [%v] \n want: empty literal", have.Literal)
			}
		})
	}
}

// func TestPeekToScanner(t *testing.T) {
// 	src := "1234567"
// 	tests := []struct {
// 		n     int
// 		value string
// 	}{
// 		{0, "1"},
// 		{1, "12"},
// 		{3, "1234"},
// 		{8, "1234567"},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.value, func(t *testing.T) {
// 			s := NewFromString("", src)
// 			value := s.PeekTo(test.n)
// 			if value != test.value {
// 				t.Errorf("\n have: [%v] \n want: [%v]", value, test.value)
// 			}
// 			Repeat(s.Keep, 10)
// 			full := s.Emit().Value
// 			testFull := "1234567"
// 			if full != testFull {
// 				t.Errorf("\n full \n have: [%v] \n want: [%v]", full, testFull)
// 			}
// 		})
// 	}

// }
