package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

func scanInt(s *scan.Scanner) scan.Token {
	s.Type = scan.IntType
	scan.While(s, scan.Digit09, s.Keep)
	return s.Emit()
}

func scanIdent(s *scan.Scanner) scan.Token {
	s.Type = scan.IdentType
	s.Keep()
	scan.While(s, scan.LetterDigitUnder, s.Keep)
	return s.Emit()
}

func scanToken7(s *scan.Scanner) scan.Token {
	switch {
	case s.Is(scan.Digit09):
		return scanInt(s)
	case s.Is(scan.LetterUnder):
		return scanIdent(s)
	}
	return s.Illegal("invalid character")
}

func Example_example7() {
	var toks []scan.Token
	s := scan.NewFromString("", "   1234 abcd 5678efgh!")
	scan.While(s, scan.Whitespace, s.Discard)
	for s.HasMore() {
		toks = append(toks, scanToken7(s))
		scan.While(s, scan.Whitespace, s.Discard)
	}
	fmt.Println(scan.FormatTokenTable(toks))

	// Output:
	//
	//  Pos  Type     Value
	//  1:4  int      1234
	//  1:9  ident    abcd
	// 1:14  int      5678
	// 1:18  ident    efgh
	// 1:22  illegal  ! (error: invalid character)
}