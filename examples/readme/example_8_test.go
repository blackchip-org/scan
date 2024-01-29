package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

func scanString(s *scan.Scanner) scan.Token {
	s.Type = scan.StrType
	s.Skip() // beginning quote

	for s.HasMore() {
		if s.This == '"' {
			break
		} else if s.This == '\\' && s.Next == '"' {
			s.Keep() // backquote
			s.Keep() // quote
		} else {
			s.Keep()
		}
	}
	if s.This != '"' {
		return s.Illegal("unterminated string")
	}
	s.Skip() // ending quote
	return s.Emit()
}

func scanToken8(s *scan.Scanner) scan.Token {
	switch {
	case s.Is(scan.Digit09):
		return scanInt(s)
	case s.Is(scan.LetterUnder):
		return scanIdent(s)
	case s.This == '"':
		return scanString(s)
	}
	return s.Illegal("invalid character")
}

func Example_example8() {
	var toks []scan.Token
	s := scan.NewFromString("", `   1234 abcd "hello \"world\"" 5678efgh!`)
	scan.While(s, scan.Whitespace, s.Discard)
	for s.HasMore() {
		toks = append(toks, scanToken8(s))
		scan.While(s, scan.Whitespace, s.Discard)
	}
	fmt.Println(scan.FormatTokenTable(toks))

	// Output:
	//
	//  Pos  Type     Value
	//  1:4  int      1234
	//  1:9  ident    abcd
	// 1:14  str      hello \"world\"
	// 1:32  int      5678
	// 1:36  ident    efgh
	// 1:40  illegal  ! (error: invalid character)
}
