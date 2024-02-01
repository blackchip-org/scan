package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

type IntRule struct {
	digit scan.Class
}

func NewIntRule(digit scan.Class) IntRule {
	return IntRule{digit: digit}
}

func (r IntRule) Eval(s *scan.Scanner) bool {
	if !s.Is(r.digit) {
		return false
	}
	s.Type = scan.IntType
	scan.While(s, r.digit, s.Keep)
	return true
}

type WordRule struct {
	letter scan.Class
}

func NewWordRule(letter scan.Class) WordRule {
	return WordRule{letter: letter}
}

func (r WordRule) Eval(s *scan.Scanner) bool {
	if !s.Is(r.letter) {
		return false
	}
	s.Type = scan.WordType
	scan.While(s, r.letter, s.Keep)
	return true
}

type SpaceRule struct {
	space scan.Class
}

func NewSpaceRule(space scan.Class) SpaceRule {
	return SpaceRule{space: space}
}

func (r SpaceRule) Eval(s *scan.Scanner) bool {
	if !s.Is(r.space) {
		return false
	}
	s.Type = scan.SpaceType
	scan.While(s, r.space, s.Keep)
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
		NewSpaceRule(scan.Whitespace),
		NewWordRule(scan.Letter),
		NewIntRule(scan.Digit),
	).WithNoMatchFunc(UnexpectedUntil(scan.Whitespace))

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
