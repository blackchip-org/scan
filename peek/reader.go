package peek

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

const EndOfText = rune(-1)

type Reader struct {
	src   *bufio.Reader
	ahead []rune
	err   error
}

func NewReader(r io.Reader) *Reader {
	return &Reader{src: bufio.NewReader(r)}
}

func (r *Reader) Read() (rune, error) {
	if r.err != nil && len(r.ahead) == 0 {
		return 0, r.err
	}
	if len(r.ahead) > 0 {
		var ch rune
		ch, r.ahead = r.ahead[0], r.ahead[1:]
		return ch, nil
	}
	ch, _, err := r.src.ReadRune()
	r.err = err
	return ch, err
}

func (r *Reader) Unread(ch rune) {
	if errors.Is(r.err, io.EOF) {
		r.err = nil
	}
	r.UnreadAll([]rune{ch})
}

func (r *Reader) UnreadAll(chs []rune) {
	r.ahead = append(chs, r.ahead...)
}

func (r *Reader) PeekTo(n int) (string, error) {
	if n < 0 {
		return "", fmt.Errorf("invalid peek value: %v", n)
	}
	if n <= len(r.ahead) {
		return string(r.ahead[:n]), nil
	}
	diff := n - len(r.ahead)
	for i := 0; i < diff; i++ {
		ch, _, err := r.src.ReadRune()
		if err != nil {
			r.err = err
			break
		}
		r.ahead = append(r.ahead, ch)
	}
	err := r.err
	if err == io.EOF {
		err = nil
	}
	return string(r.ahead), err
}

func (r *Reader) Peek(n int) (rune, error) {
	_, err := r.PeekTo(n)
	if err != nil {
		return EndOfText, err
	}
	if n > len(r.ahead) {
		return EndOfText, r.err
	}
	return r.ahead[n-1], nil
}
