package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

// hex integer with class constructors
func Example_example3() {
	digit0F := scan.Or(
		scan.Range('0', '9'),
		scan.Range('A', 'F'),
		scan.Range('a', 'f'),
	)

	s := scan.NewFromString("", "1010234abc!")
	scan.While(s, digit0F, s.Keep)
	tok := s.Emit()
	fmt.Println(tok.Value)

	// Output:
	// 1010234abc
}
