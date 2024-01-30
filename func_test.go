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
		{"", Whitespace, s.Discard},
		{"123", Digit, s.Keep},
		{"abc", Any, s.Keep},
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

func TestQuote(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{`abcd`, `"abcd"`},
		{`ab"c"d`, `'ab"c"d'`},
		{`ab'c'd`, `"ab'c'd"`},
		{"ab`c`d", "\"ab`c`d\""},
		{"a\"b\"'c'd", "`a\"b\"'c'd`"},
		{"a\"b\"'c'`d`", "{!quote:a\"b\"'c'`d`}"},
		{"ab\x00cd", `"ab{!ch:00}cd"`},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			have := Quote(test.in)
			if have != test.want {
				t.Errorf("\n have: %v \n want: %v", have, test.want)
			}
		})
	}
}
