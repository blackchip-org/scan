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
			s := NewScannerFromString("", src)
			Repeat(s.Keep, test.n)
			have := s.Emit()
			if have.Val != test.value {
				t.Errorf("\n have: [%v] \n want: [%v]", have.Val, test.value)
			}
		})
	}
}

func TestUndo(t *testing.T) {
	s := NewScannerFromString("", "123abcd")
	While(s, Digit, s.Discard)
	s.Keep()
	s.Keep()
	s.Undo()
	While(s, Any, s.Keep)
	tok := s.Emit()
	want := "abcd"
	if tok.Val != want {
		t.Errorf("\n have: %v \n want: %v", tok.Val, want)
	}
}

func TestPeek(t *testing.T) {
	s := NewScannerFromString("", "xyz0123")
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
	s := NewScannerFromString("", "")
	s.Keep()
	want := Token{
		Pos:  Pos{Line: 1, Col: 1},
		Type: EndOfTextType,
	}
	tok := s.Emit()
	if !tok.Equal(want) {
		t.Errorf("\n have: %v \n want: %v", tok, want)
	}
}

func TestEnd(t *testing.T) {
	s := NewScannerFromString("", "1")
	s.Keep()
	s.Emit()
	Repeat(s.Keep, 10)

	want := Token{
		Pos:  Pos{Line: 1, Col: 2},
		Type: EndOfTextType,
	}
	tok := s.Emit()
	if !tok.Equal(want) {
		t.Errorf("\n have: %v \n want: %v", tok, want)
	}
}

func TestReInit(t *testing.T) {
	s := NewScannerFromString("foo", "abc")
	s.Discard()
	s.Keep()
	have := s.Emit()
	want := "b"
	if have.Val != want {
		t.Errorf("\n have: %v \n want: %v", have, want)
	}
	wantPos := Pos{Name: "foo", Line: 1, Col: 2}
	if have.Pos != wantPos {
		t.Errorf("\n have: %v \n want: %v", have.Pos, wantPos)
	}

	s.InitFromString("bar", "123")
	s.Discard()
	s.Keep()
	have = s.Emit()
	want = "2"
	if have.Val != want {
		t.Errorf("\n have: %v \n want: %v", have, want)
	}
	wantPos = Pos{Name: "bar", Line: 1, Col: 2}
	if have.Pos != wantPos {
		t.Errorf("\n have: %v \n want: %v", have.Pos, wantPos)
	}
}

func TestInit(t *testing.T) {
	var s Scanner
	s.InitFromString("foo", "abc")
	s.Discard()
	s.Keep()
	have := s.Emit()
	want := "b"
	if have.Val != want {
		t.Errorf("\n have: %v \n want: %v", have, want)
	}
}
