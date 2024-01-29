package iotest

import (
	"errors"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	s := "abc123def456"

	testErr := errors.New("test error")
	tests := []struct {
		want string
		err  error
	}{
		{"abc", nil},
		{"123", nil},
		{"de", testErr},
	}

	buf := make([]byte, 3)
	sr := strings.NewReader(s)
	r := NewReader(sr)
	r.Limit = 8
	r.Err = testErr

	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			n, err := r.Read(buf)
			data := string(buf[:n])
			if data != test.want {
				t.Errorf("\n have: %v \n want: %v", string(buf), test.want)
			}
			if err != test.err {
				t.Errorf("\n have: %v \n want: %v", err, test.err)
			}
		})
	}
}
