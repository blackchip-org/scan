package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

func Example_example9() {
	s := scan.NewFromString("", "   1234 abcd 5678efgh!")
	rules := scan.Rules(
		scan.StrDoubleQuote,
		scan.Int,
		scan.Ident,
	)
	runner := scan.NewRunner(s, rules)
	toks := runner.All()
	fmt.Println(scan.FormatTokenTable(toks))

	// Output:
	//
	//  Pos  Type     Value
	//  1:4  int      1234
	//  1:9  ident    abcd
	// 1:14  int      5678
	// 1:18  ident    efgh
	// 1:22  illegal  !
	// 1:22: error: unexpected "!"
}
