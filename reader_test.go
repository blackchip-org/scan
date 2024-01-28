package scan

import (
	"io"
	"strconv"
	"strings"
	"testing"
)

func TestPeekTo(t *testing.T) {
	text := "1234567"
	tests := []struct {
		n   int
		val string
	}{
		{0, ""},
		{1, "1"},
		{3, "123"},
		{8, "1234567"},
	}

	for _, test := range tests {
		t.Run(strconv.Itoa(test.n), func(t *testing.T) {
			r := NewPeekReader(strings.NewReader(text))
			peek, _ := r.PeekTo(test.n)
			if peek != test.val {
				t.Errorf("\n have: [%v] \n want: [%v]", peek, test.val)
			}
			var valBuf []rune
			for {
				ch, err := r.Read()
				if err != nil {
					break
				}
				valBuf = append(valBuf, ch)
			}
			val := string(valBuf)
			if val != text {
				t.Errorf("\n full \n have: [%v] \n want: [%v]", val, text)
			}
		})
	}
}

func TestPeekAt(t *testing.T) {
	text := "1234567"
	tests := []struct {
		n   int
		val rune
	}{
		{1, '1'},
		{3, '3'},
		{8, EndOfText},
	}

	for _, test := range tests {
		t.Run(strconv.Itoa(test.n), func(t *testing.T) {
			r := NewPeekReader(strings.NewReader(text))
			peek, _ := r.PeekAt(test.n)
			if peek != test.val {
				t.Errorf("\n have: [%c] \n want: [%c]", peek, test.val)
			}
		})
	}
}

func TestReaderEOF(t *testing.T) {
	r := NewPeekReader(strings.NewReader("123"))
	v, err := r.PeekTo(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "123" {
		t.Fatalf("\n have: %v \n want: %v", v, "123")
	}
	for i := 0; i < 3; i++ {
		_, err := r.Read()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	_, err = r.Read()
	if err != io.EOF {
		t.Errorf("\n have: %v \n want: EOF", err)
	}
}
