package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

type IntRule2 struct {
	digit    scan.Class
	digitSep scan.Class
}

func NewIntRule2(digit scan.Class) IntRule2 {
	return IntRule2{digit: digit}
}

func (r IntRule2) WithDigitSep(digitSep scan.Class) IntRule2 {
	r.digitSep = digitSep
	return r
}

func (r IntRule2) Eval(s *scan.Scanner) bool {
	if !s.Is(r.digit) {
		return false
	}
	s.Keep()
	for s.HasMore() {
		if s.Is(r.digit) {
			s.Keep()
		} else if s.Is(r.digitSep) {
			s.Skip()
		} else {
			break
		}
	}
	return true
}

func Example_example8() {
	rules := scan.NewRuleSet(
		NewSpaceRule(scan.Whitespace),
		NewWordRule(scan.Letter),
		NewIntRule2(scan.Digit).
			WithDigitSep(scan.Rune(',')),
	).WithNoMatchFunc(UnexpectedUntil(scan.Whitespace))

	s := scan.NewScannerFromString("example8", "abc 1,234 !@#  \tdef45,678")
	runner := scan.NewRunner(s, rules)
	toks := runner.All()
	fmt.Println(scan.FormatTokenTable(toks))

	// Output:
	//
	// Pos            Type     Value         Literal
	// example8:1:1   word     "abc"         "abc"
	// example8:1:4   space    " "           " "
	// example8:1:5   1234     "1234"        "1,234"
	// example8:1:10  space    " "           " "
	// example8:1:11  illegal  "!@#"         "!@#"
	// example8:1:14: error: unexpected "!@#"
	// example8:1:14  space    "  {!ch:\t}"  "  {!ch:\t}"
	// example8:1:17  word     "def"         "def"
	// example8:1:20  45678    "45678"       "45,678"
}
