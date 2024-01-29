package peek

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/blackchip-org/scan/iotest"
)

func TestPeek(t *testing.T) {
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
			r := NewReader(strings.NewReader(text))
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
			r := NewReader(strings.NewReader(text))
			peek, _ := r.Peek(test.n)
			if peek != test.val {
				t.Errorf("\n have: [%c] \n want: [%c]", peek, test.val)
			}
		})
	}
}

func TestReaderEOF(t *testing.T) {
	r := NewReader(strings.NewReader("123"))
	v, err := r.PeekTo(4)
	if err != nil {
		t.Fatal(err)
	}
	if v != "123" {
		t.Fatalf("\n have: %v \n want: %v", v, "123")
	}
	for i := 0; i < 3; i++ {
		_, err := r.Read()
		if err != nil {
			t.Fatal(err)
		}
	}
	_, err = r.Read()
	if err != io.EOF {
		t.Errorf("\n have: %v \n want: EOF", err)
	}
}

func TestUnreadAll(t *testing.T) {
	text := "1234567"
	tests := []struct {
		n      int    // read these runes first
		unread string // then unread these runes
		val    string // remaining reads should contain this
	}{
		{0, "abc", "abc1234567"},
		{4, "abc", "abc567"},
		{7, "abc", "abc"},
	}

	for _, test := range tests {
		t.Run(test.val, func(t *testing.T) {
			r := NewReader(strings.NewReader(text))
			for i := 0; i < test.n; i++ {
				_, err := r.Read()
				if err != nil {
					t.Fatal(err)
				}
			}
			r.UnreadAll([]rune(test.unread))
			var have []rune
			for {
				ch, err := r.Read()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					t.Fatal(err)
				}
				have = append(have, ch)
			}
			if string(have) != test.val {
				t.Errorf("\n have: %v \n want: %v", string(have), test.val)
			}
		})
	}
}

func TestReadError(t *testing.T) {
	text := "123456"
	sr := iotest.NewReader(strings.NewReader(text))
	testErr := errors.New("test error")
	sr.Limit = 3
	sr.Err = testErr

	r := NewReader(sr)
	r.Read()              // 1
	r.Read()              // 2
	v, err := r.PeekTo(3) // 3 only
	if v != "3" {
		t.Fatalf("have: %v \n want: 3", v)
	}
	if err != testErr {
		t.Fatal(err)
	}

	ch, err := r.Read()
	if ch != '3' {
		t.Fatalf("have: %b \n want: 3", ch)
	}
	if err != nil {
		t.Fatal(err)
	}

	_, err = r.Read()
	if err != testErr {
		t.Fatal(err)
	}

	r.Unread('X')
	ch, err = r.Read()
	if ch != 'X' {
		t.Fatalf("have: %b \n want: X", ch)
	}
	if err != nil {
		t.Fatal(err)
	}

	_, err = r.Read()
	if err != testErr {
		t.Fatal(err)
	}
}
