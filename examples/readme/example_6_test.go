package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

func scanInt(s *scan.Scanner) bool {
	if !s.Is(scan.Digit09) {
		return false
	}
	s.Type = scan.IntType
	scan.While(s, scan.Digit09, s.Keep)
	return true
}

func scanWord(s *scan.Scanner) bool {
	if !s.Is(scan.Letter) {
		return false
	}
	s.Type = scan.WordType
	scan.While(s, scan.Letter, s.Keep)
	return true
}

func scanSpace(s *scan.Scanner) bool {
	if !s.Is(scan.Whitespace) {
		return false
	}
	s.Type = scan.SpaceType
	scan.While(s, scan.Whitespace, s.Keep)
	return true
}

func Example_example6() {
	scanFuncs := []func(*scan.Scanner) bool{
		scanSpace,
		scanWord,
		scanInt,
	}

	var toks []scan.Token
	s := scan.NewScannerFromString("example6", "abc 123 !@#  \tdef456")
	for s.HasMore() {
		match := false
		for _, fn := range scanFuncs {
			if match = fn(s); match {
				toks = append(toks, s.Emit())
				break
			}
		}
		if !match {
			scan.Until(s, scan.Whitespace, s.Keep)
			s.Illegal("unexpected %v", scan.Quote(s.Val.String()))
			toks = append(toks, s.Emit())
		}
	}
	fmt.Println(scan.FormatTokenTable(toks))

	// Output:
	//
	// Pos            Type     Value         Literal
	// example6:1:1   word     "abc"         "abc"
	// example6:1:4   space    " "           " "
	// example6:1:5   int      "123"         "123"
	// example6:1:8   space    " "           " "
	// example6:1:9   illegal  "!@#"         "!@#"
	// example6:1:12: error: unexpected "!@#"
	// example6:1:12  space    "  {!ch:\t}"  "  {!ch:\t}"
	// example6:1:15  word     "def"         "def"
	// example6:1:18  int      "456"         "456"
}
