package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

func Example_example9() {
	rules := scan.NewRuleSet(
		scan.KeepSpaceRule,
		scan.NewWhileRule(scan.IsLetter, scan.WordType),
		scan.IntRule.WithDigitSep(scan.Rune(',')),
	).WithNoMatchFunc(scan.UnexpectedUntil(scan.IsSpace))

	s := scan.NewScannerFromString("example9", "abc 1,234 !@#  \tdef45,678")
	runner := scan.NewRunner(s, rules)
	toks := runner.All()
	fmt.Println(scan.FormatTokenTable(toks))

	// Output:
	//
	// Pos            Type     Value         Literal
	// example9:1:1   word     "abc"         "abc"
	// example9:1:4   space    " "           " "
	// example9:1:5   int      "1234"        "1,234"
	// example9:1:10  space    " "           " "
	// example9:1:11  illegal  "!@#"         "!@#"
	// example9:1:14: error: unexpected "!@#"
	// example9:1:14  space    "  {!ch:\t}"  "  {!ch:\t}"
	// example9:1:17  word     "def"         "def"
	// example9:1:20  int      "45678"       "45,678"
}
