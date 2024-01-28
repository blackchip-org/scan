package scan

import "testing"

func TestRepeat(t *testing.T) {
	s := NewFromString("", "12345")
	tests := []struct {
		val string
		n   int
	}{
		{"123", 3},
		{"45", 3},
	}

	for _, test := range tests {
		t.Run(test.val, func(t *testing.T) {
			Repeat(s.Keep, test.n)
			tok := s.Emit()
			if tok.Value != test.val {
				t.Errorf("\n have: %v \n want: %v", tok.Value, test.val)
			}
		})
	}
}

func TestWhile(t *testing.T) {
	s := NewFromString("", "   123abc")
	tests := []struct {
		val   string
		class Class
		fn    func()
	}{
		{"", IsWhitespace, s.Discard},
		{"123", IsDigit, s.Keep},
		{"abc", IsAny, s.Keep},
	}

	for _, test := range tests {
		t.Run(test.val, func(t *testing.T) {
			While(s, test.class, test.fn)
			tok := s.Emit()
			if tok.Value != test.val {
				t.Errorf("\n have: %v \n want: %v", tok.Value, test.val)
			}
		})
	}
}
