package scan

import (
	"bufio"
	"errors"
	"io"
)

type MarkReader struct {
	src    *bufio.Reader
	record bool
	undo   []rune
}

func NewMarkReader(r io.Reader) *MarkReader {
	return &MarkReader{src: bufio.NewReader(r)}
}

func (r *MarkReader) Read() (rune, error) {
	if !r.record {
		if len(r.undo) > 0 {
			var ch rune
			ch, r.undo = r.undo[0], r.undo[1:]
			return ch, nil
		}
		ch, _, err := r.src.ReadRune()
		return ch, err
	}
	ch, _, err := r.src.ReadRune()
	if err != nil {
		return 0, err
	}
	r.undo = append(r.undo, ch)
	return ch, nil
}

func (r *MarkReader) ReadAll() (string, error) {
	var all []rune
	var err error
	for {
		var ch rune
		ch, err = r.Read()
		if err != nil {
			break
		}
		all = append(all, ch)
	}
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return string(all), err
}

func (r *MarkReader) Mark() {
	if r.record {
		panic("stream already marked")
	}
	r.record = true
}

func (r *MarkReader) Reset() {
	r.record = false
}

func (r *MarkReader) Unmark() {
	r.record = true
	r.undo = nil
}
