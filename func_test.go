package scan

import "testing"

func TestRepeat(t *testing.T) {
	s := NewScannerFromString("", "12345")
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
			if tok.Val != test.val {
				t.Errorf("\n have: %v \n want: %v", tok.Val, test.val)
			}
		})
	}
}

func TestWhile(t *testing.T) {
	s := NewScannerFromString("", "   123abc456")
	tests := []struct {
		val   string
		class Class
		fn    func()
	}{
		{"", Whitespace, s.Discard},
		{"123", Digit, s.Keep},
		{"abc", Letter, s.Keep},
		{"456", Not(Whitespace), s.Keep},
	}

	for _, test := range tests {
		t.Run(test.val, func(t *testing.T) {
			While(s, test.class, test.fn)
			tok := s.Emit()
			if tok.Val != test.val {
				t.Errorf("\n have: %v \n want: %v", tok.Val, test.val)
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
		{"\u00a0", `"{!ch:a0}"`},
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

func TestWord(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{`abcd`, `abcd`},
		{`    abcd`, `abcd`},
		{"  ab\ncd", `ab`},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			s := NewScannerFromString("", test.in)
			have := Word(s)
			if have != test.want {
				t.Errorf("\n have: %v \n want: %v", have, test.want)
			}
		})
	}
}

func TestLine(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"  abcd  \n 1234", "abcd"},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			s := NewScannerFromString("", test.in)
			have := Line(s)
			if have != test.want {
				t.Errorf("\n have: %v \n want: %v", have, test.want)
			}
		})
	}
}
