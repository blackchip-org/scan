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
	BinType     = "bin"
	CommentType = "comment"
	ErrorType   = "error"
	HexType     = "hex"
	IdentType   = "ident"
	IllegalType = "illegal"
	IntType     = "int"
	OctType     = "oct"
	RealType    = "real"
	SpaceType   = "space"
	StrType     = "str"
	WordType    = "word"
)

// Pos represents a position within an input stream.
type Pos struct {
	Name string `json:"name,omitempty"`
	Line int    `json:"line"`
	Col  int    `json:"col"`
}

// NewPos returns a new position with the given name and the line number
// and column number set to one.
func NewPos(name string) Pos {
	return Pos{Name: name, Line: 1, Col: 1}
}

func (p Pos) String() string {
	if p.Name != "" {
		return fmt.Sprintf("%v:%v:%v", p.Name, p.Line, p.Col)
	}
	return fmt.Sprintf("%v:%v", p.Line, p.Col)
}

type Error struct {
	Pos     Pos
	Message string
	Cause   error
}

func (e Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%v: error: %v: %v", e.Pos, e.Message, e.Cause)
	}
	return fmt.Sprintf("%v: error: %v", e.Pos, e.Message)
}

func (e Error) Unwrap() error {
	return e.Cause
}

type Errors []Error

func (es Errors) Error() string {
	var messages []string
	for _, e := range es {
		messages = append(messages, e.Error())
	}
	return strings.Join(messages, "\n")
}

type Token struct {
	Val  string `json:"val"`
	Lit  string `json:"lit,omitempty"`
	Type string `json:"type"`
	Pos  Pos    `json:"pos"`
	Errs Errors `json:"errs,omitempty"`
}

func (t Token) IsValid() bool {
	return t.Val != "" || t.Type != ""
}

func (t Token) String() string {
	var b strings.Builder

	if t.Pos.Line != 0 || t.Pos.Col != 0 {
		fmt.Fprintf(&b, "%v:%v ", t.Pos.Line, t.Pos.Col)
	}
	fmt.Fprintf(&b, "type:%v val:%v", t.Type, Quote(t.Val))
	if t.Lit != "" {
		fmt.Fprintf(&b, " lit:%v", Quote(t.Lit))
	}
	if len(t.Errs) > 0 {
		fmt.Fprintf(&b, "\n%v", t.Errs)
	}
	return b.String()
}

type Scanner struct {
	This    rune
	Next    rune
	Val     strings.Builder
	Lit     strings.Builder
	Errs    Errors
	Type    string
	src     *peek.Reader
	srcErr  *Error
	thisPos Pos
	tokPos  Pos
}

func NewScanner(name string, src io.Reader) *Scanner {
	s := &Scanner{}
	s.Init(name, src)
	return s
}

func NewScannerFromString(name string, src string) *Scanner {
	s := &Scanner{}
	s.InitFromString(name, src)
	return s
}

func NewScannerFromBytes(name string, src []byte) *Scanner {
	s := &Scanner{}
	s.InitFromBytes(name, src)
	return s
}

func (s *Scanner) Init(name string, src io.Reader) {
	s.This = 0
	s.Next = 0
	s.Val.Reset()
	s.Lit.Reset()
	s.Errs = nil
	s.Type = ""
	s.src = peek.NewReader(src)
	s.srcErr = nil
	s.next()
	s.next()
	s.thisPos = NewPos(name)
	s.tokPos = NewPos(name)
}

func (s *Scanner) InitFromString(name string, src string) {
	s.Init(name, strings.NewReader(src))
}

func (s *Scanner) InitFromBytes(name string, src []byte) {
	s.Init(name, bytes.NewReader(src))
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
	chs := []rune(s.Lit.String())
	li := len(chs) + i
	if li < 0 {
		return EndOfText
	}
	return chs[li]
}

func (s *Scanner) PeekIs(i int, c Class) bool {
	if c == nil {
		return false
	}
	return c(s.Peek(i))
}

// Keep advances the stream to the next rune and adds the current rune to
// the token value and literal.
func (s *Scanner) Keep() {
	if s.This != EndOfText {
		s.Val.WriteRune(s.This)
		s.Lit.WriteRune(s.This)
		s.next()
	}
}

// Skip advances the string to the next rune without adding the current rune to
// the token value. This current rune is written to the literal value.
func (s *Scanner) Skip() {
	s.Lit.WriteRune(s.This)
	s.next()
}

// Discard advances the string to the next rune without adding the current
// rune to the token.
func (s *Scanner) Discard() {
	s.next()
	s.Emit()
}

func (s *Scanner) Undo() {
	chs := []rune(s.Lit.String())
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
	s.Val.Reset()
	s.Lit.Reset()
	s.Type = ""
	s.thisPos = s.tokPos
}

// Emit returns the token that has been built and resets the builder for the
// next token.
func (s *Scanner) Emit() Token {
	var t Token

	t.Val = s.Val.String()
	t.Lit = s.Lit.String()
	t.Type = s.Type
	t.Pos = s.tokPos
	t.Errs = s.Errs

	if t.Lit == "" && s.srcErr != nil {
		t.Type = ErrorType
		t.Errs = append(t.Errs, *s.srcErr)
		return t
	}

	// If there are no type yet, set it to the value
	if t.Type == "" {
		t.Type = t.Val
	}

	s.Val.Reset()
	s.Lit.Reset()
	s.Type = ""
	s.Errs = nil
	s.tokPos = s.thisPos

	return t
}

func (s *Scanner) Illegal(format string, args ...any) {
	s.Type = IllegalType
	s.Errs = append(s.Errs, Error{
		Pos:     s.thisPos,
		Message: fmt.Sprintf(format, args...),
	})
}

func (s *Scanner) next() {
	if s.This == EndOfText {
		return
	}
	if s.This == '\n' {
		s.thisPos.Line++
		s.thisPos.Col = 1
	} else {
		s.thisPos.Col++
	}

	var err error
	s.This = s.Next
	s.Next, err = s.src.Read()

	if err != nil {
		// Mark the stream as done when seeing an EOF but don't retain that
		// as an actual error
		if !errors.Is(err, io.EOF) {
			s.srcErr = &Error{
				Pos:     s.thisPos,
				Message: fmt.Sprintf("error reading stream: %v", err),
				Cause:   err,
			}
		}
		s.Next = EndOfText
	}
}
