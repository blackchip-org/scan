package scan

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/blackchip-org/scan/peek"
)

const (
	EndOfText = rune(-1)
)

const (
	ErrorType = "error"
)

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
	Type    string
	Pos     Pos
	Err     *Error
}

type Scanner struct {
	This    rune
	Next    rune
	Value   strings.Builder
	Literal strings.Builder
	Type    string
	src     *peek.Reader
	err     *Error
	nextErr *Error
	thisPos Pos
	tokPos  Pos
}

func New(name string, src io.Reader) *Scanner {
	s := &Scanner{src: peek.NewReader(src)}
	s.next()
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

func (s *Scanner) NextIs(c Class) bool {
	if c == nil {
		return false
	}
	return c(s.Next)
}

func (s *Scanner) Peek(i int) rune {
	if i == 0 {
		return s.This
	}
	if i == 1 {
		return s.Next
	}
	if i > 1 {
		ch, err := s.src.Peek(i - 2 + 1)
		if err != nil {
			return EndOfText
		}
		return ch
	}
	chs := []rune(s.Literal.String())
	li := len(chs) + i
	if li < 0 {
		return EndOfText
	}
	return chs[li]
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

func (s *Scanner) Undo() {
	chs := []rune(s.Literal.String())
	slices.Reverse(chs)
	for _, ch := range chs {
		if s.Next != EndOfText {
			s.src.Unread(s.Next)
		}
		if s.This != EndOfText {
			s.Next = s.This
		}
		s.This = ch
	}
	s.Value.Reset()
	s.Literal.Reset()
	s.Type = ""
	s.thisPos = s.tokPos
}

// Emit returns the token that has been built and resets the builder for the
// next token.
func (s *Scanner) Emit() Token {
	var t Token

	t.Value = s.Value.String()
	t.Literal = s.Literal.String()
	t.Type = s.Type
	t.Pos = s.tokPos
	t.Err = s.err

	// Only include the literal if it differs from the value
	if t.Value == t.Literal {
		t.Literal = ""
	}
	// If there are no type yet, set it to the value
	if t.Type == "" {
		t.Type = t.Value
	}
	// But if there was an error, set the type to error
	if s.err != nil {
		t.Type = ErrorType
	}

	s.Value.Reset()
	s.Literal.Reset()
	s.Type = ""
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
	s.This, s.err = s.Next, s.nextErr
	s.Next, err = s.src.Read()

	if err != nil {
		// Mark the stream as done when seeing an EOF but don't retain that
		// as an actual error
		if !errors.Is(err, io.EOF) {
			s.nextErr = &Error{
				Pos:     s.thisPos,
				Message: fmt.Sprintf("error reading stream: %v", err),
				Cause:   err,
			}
		}
		s.Next = EndOfText
	}
}
