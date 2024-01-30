package scan

import (
	"strconv"
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

func TestUndo(t *testing.T) {
	s := NewFromString("", "123abcd")
	While(s, Digit, s.Discard)
	s.Keep()
	s.Keep()
	s.Undo()
	While(s, Any, s.Keep)
	tok := s.Emit()
	want := "abcd"
	if tok.Value != want {
		t.Errorf("\n have: %v \n want: %v", tok.Value, want)
	}
}

func TestPeek(t *testing.T) {
	s := NewFromString("", "xyz0123")
	tests := []struct {
		i  int
		ch rune
	}{
		{0, '0'},
		{1, '1'},
		{2, '2'},
		{3, '3'},
		{4, EndOfText},
		{-1, 'z'},
		{-3, 'x'},
		{-4, EndOfText},
	}

	Repeat(s.Skip, 3)
	for _, test := range tests {
		t.Run(strconv.Itoa(test.i), func(t *testing.T) {
			have := s.Peek(test.i)
			if have != test.ch {
				t.Errorf("\n have: %c \n want: %c", have, test.ch)
			}
		})
	}
}

func TestEmpty(t *testing.T) {
	s := NewFromString("", "")
	s.Keep()
	tok := s.Emit()
	if tok.Value != "" || tok.Type != "" {
		t.Error(tok)
	}
}
