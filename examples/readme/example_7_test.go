package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

type IntRule struct {
	isDigit scan.Class
}

func NewIntRule(isDigit scan.Class) IntRule {
	return IntRule{isDigit: isDigit}
}

func (r IntRule) Eval(s *scan.Scanner) bool {
	if !r.isDigit(s.This) {
		return false
	}
	s.Type = scan.IntType
	scan.While(s, r.isDigit, s.Keep)
	return true
}

type WordRule struct {
	isLetter scan.Class
}

func NewWordRule(isLetter scan.Class) WordRule {
	return WordRule{isLetter: isLetter}
}

func (r WordRule) Eval(s *scan.Scanner) bool {
	if !r.isLetter(s.This) {
		return false
	}
	s.Type = scan.WordType
	scan.While(s, r.isLetter, s.Keep)
	return true
}

type SpaceRule struct {
	isSpace scan.Class
}

func NewSpaceRule(isSpace scan.Class) SpaceRule {
	return SpaceRule{isSpace: isSpace}
}

func (r SpaceRule) Eval(s *scan.Scanner) bool {
	if !r.isSpace(s.This) {
		return false
	}
	s.Type = scan.SpaceType
	scan.While(s, r.isSpace, s.Keep)
	return true
}

func UnexpectedUntil(c scan.Class) func(*scan.Scanner) {
	return func(s *scan.Scanner) {
		scan.Until(s, c, s.Keep)
		s.Illegal("unexpected %s", scan.Quote(s.Lit.String()))
	}
}

func Example_example7() {
	rules := scan.NewRuleSet(
		NewSpaceRule(scan.IsSpace),
		NewWordRule(scan.IsLetter),
		NewIntRule(scan.IsDigit),
	).WithNoMatchFunc(UnexpectedUntil(scan.IsSpace))

	s := scan.NewScannerFromString("example7", "abc 123 !@#  \tdef456")
	runner := scan.NewRunner(s, rules)
	toks := runner.All()
	fmt.Println(scan.FormatTokenTable(toks))

	// Output:
	//
	// Pos            Type     Value         Literal
	// example7:1:1   word     "abc"         "abc"
	// example7:1:4   space    " "           " "
	// example7:1:5   int      "123"         "123"
	// example7:1:8   space    " "           " "
	// example7:1:9   illegal  "!@#"         "!@#"
	// example7:1:12: error: unexpected "!@#"
	// example7:1:12  space    "  {!ch:\t}"  "  {!ch:\t}"
	// example7:1:15  word     "def"         "def"
	// example7:1:18  int      "456"         "456"
}
