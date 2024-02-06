package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

type IntRule2 struct {
	isDigit    scan.Class
	isDigitSep scan.Class
}

func NewIntRule2(isDigit scan.Class) IntRule2 {
	return IntRule2{
		isDigit:    isDigit,
		isDigitSep: scan.IsNone,
	}
}

func (r IntRule2) WithDigitSep(isDigitSep scan.Class) IntRule2 {
	r.isDigitSep = isDigitSep
	return r
}

func (r IntRule2) Eval(s *scan.Scanner) bool {
	if !r.isDigit(s.This) {
		return false
	}
	s.Keep()
	for s.HasMore() {
		if r.isDigit(s.This) {
			s.Keep()
		} else if r.isDigitSep(s.This) {
			s.Skip()
		} else {
			break
		}
	}
	return true
}

func Example_example8() {
	rules := scan.NewRuleSet(
		NewSpaceRule(scan.IsSpace),
		NewWordRule(scan.IsLetter),
		NewIntRule2(scan.IsDigit).
			WithDigitSep(scan.Rune(',')),
	).WithNoMatchFunc(UnexpectedUntil(scan.IsSpace))

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
