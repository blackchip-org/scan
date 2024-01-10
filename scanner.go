package scan

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

const EndOfText = rune(-1)

// Pos represents a position within an input stream.
type Pos struct {
	Name   string `json:",omitempty"`
	Line   int
	Column int
}

// NewPos returns a new position with the given name and the line number
// and column number set to one.
func NewPos(name string) Pos {
	return Pos{Name: name, Line: 1, Column: 1}
}

func (p Pos) String() string {
	if p.Name != "" {
		return fmt.Sprintf("%v:%v:%v", p.Name, p.Line, p.Column)
	}
	return fmt.Sprintf("%v:%v", p.Line, p.Column)
}

// General token tags
const (
	BinTag     = "bin"
	CommentTag = "comment"
	ErrTag     = "err"
	RealTag    = "real"
	ExpTag     = "exp"
	HexTag     = "hex"
	IdentTag   = "ident"
	IntTag     = "int"
	NumTag     = "num"
	OctTag     = "oct"
	StrTag     = "str"
)

type Error struct {
	Pos     Pos
	Message string
	Cause   error
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) Unwrap() error {
	return e.Cause
}

type Token struct {
	Value   string
	Literal string `json:",omitempty"`
	Tags    []string
	Pos     Pos
	Err     *Error
}

type Scanner struct {
	This    rune
	Value   strings.Builder
	Literal strings.Builder
	Tags    []string
	src     *PeekReader
	err     *Error
	thisPos Pos
	tokPos  Pos
}

func New(name string, src io.Reader) *Scanner {
	s := &Scanner{src: NewPeekReader(src)}
	s.next()
	s.thisPos = NewPos(name)
	s.tokPos = NewPos(name)
	return s
}

func NewFromString(name string, src string) *Scanner {
	return New(name, strings.NewReader(src))
}

func NewFromBytes(name string, src []byte) *Scanner {
	return New(name, bytes.NewReader(src))
}

// HasMore return true if there is more data to consume in the stream.
func (s *Scanner) HasMore() bool {
	return s.This != EndOfText
}

func (s *Scanner) Is(c Class) bool {
	if c == nil {
		return false
	}
	return c(s.This)
}

// Keep advances the stream to the next rune and adds the current rune to
// the token value and literal.
func (s *Scanner) Keep() {
	if s.This != EndOfText {
		s.Value.WriteRune(s.This)
		s.Literal.WriteRune(s.This)
		s.next()
	}
}

// Skip advances the string to the next rune without adding the current rune to
// the token value. This current rune is written to the literal value.
func (s *Scanner) Skip() {
	s.Literal.WriteRune(s.This)
	s.next()
}

// Discard advances the string to the next rune without adding the current
// rune to the token.
func (s *Scanner) Discard() {
	s.next()
}

// Peek looks ahead to the nth rune in the input stream and returns that
// value. If an error or end of stream is encountered, EndOfText is returned.
// When called with n == 0, the current rune is returned.
func (s *Scanner) Peek(n int) rune {
	if s.This == EndOfText {
		return EndOfText
	}
	if n == 0 {
		return s.This
	}
	ahead, _ := s.src.PeekAt(n)
	return ahead
}

// PeekTo looks ahead and returns the string containing the nth runes in
// the stream. The string returned includes the current rune. If an error or
// end of stream is encountered, a blank string is returned.
func (s *Scanner) PeekTo(n int) string {
	if s.This == EndOfText {
		return ""
	}
	ahead, _ := s.src.PeekTo(n)
	return string(s.This) + ahead
}

// Emit returns the token that has been built and resets the builder for the
// next token.
func (s *Scanner) Emit() Token {
	var t Token

	t.Value = s.Value.String()
	t.Literal = s.Literal.String()
	t.Tags = slices.Clone(s.Tags)
	t.Pos = s.tokPos
	t.Err = s.err

	// Only include the literal if it differs from the value
	if t.Value == t.Literal {
		t.Literal = ""
	}
	// If there are no tags yet, add the value as a tag
	if len(t.Tags) == 0 && t.Value != "" {
		t.Tags = append(t.Tags, t.Value)
	}
	if s.err != nil {
		t.Tags = append(t.Tags, ErrTag)
	}

	s.Value.Reset()
	s.Literal.Reset()
	s.Tags = nil
	s.tokPos = s.thisPos

	return t
}

func (s *Scanner) next() {
	if s.This == EndOfText {
		return
	}
	if s.This == '\n' {
		s.thisPos.Line++
		s.thisPos.Column = 1
	} else {
		s.thisPos.Column++
	}

	var err error
	s.This, err = s.src.Read()

	if err != nil {
		// Mark the stream as done when seeing an EOF but don't retain that
		// as an actual error
		if !errors.Is(err, io.EOF) {
			s.err = &Error{
				Pos:     s.thisPos,
				Message: fmt.Sprintf("error reading stream: %v", err),
				Cause:   err,
			}
		}
		s.This = EndOfText
	}
}

// Repeat calls function fn in a loop for n times.
func Repeat(fn func(), n int) {
	for i := 0; i < n; i++ {
		fn()
	}
}
