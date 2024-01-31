package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

// replace with predefined rules
func Example_example9() {
	s := scan.NewFromString("", `   1234 abcd "hello \"world\"" 5678efgh!`)
	rules := scan.Rules(
		scan.StrDoubleQuote,
		scan.Int,
		scan.Ident,
	).WithDiscards(scan.Whitespace)

	runner := scan.NewRunner(s, rules)
	toks := runner.All()
	fmt.Println(scan.FormatTokenTable(toks))

	// Output:
	//
	//  Pos  Type     Value
	//  1:4  int      1234
	//  1:9  ident    abcd
	// 1:14  str      hello "world"
	// 1:32  int      5678
	// 1:36  ident    efgh
	// 1:40  illegal  !
	// 1:40: error: unexpected "!"
}
