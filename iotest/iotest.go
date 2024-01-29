package iotest

import (
	"io"
)

type Reader struct {
	src   io.Reader
	Limit int
	Err   error
}

func NewReader(src io.Reader) *Reader {
	return &Reader{src: src}
}

func (r *Reader) Read(p []byte) (int, error) {
	if r.Limit == 0 {
		return 0, r.Err
	}
	n, err := r.src.Read(p)
	if err != nil {
		return n, err
	}
	if n > r.Limit {
		for i := r.Limit; i < len(p); i++ {
			p[i] = 0
		}
		n = r.Limit
		err = r.Err
	}
	r.Limit -= n
	return n, err
}
