package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

// discarding whitespace
func Example_example5() {
	s := scan.NewFromString("", " \t 1010234abc!")
	scan.While(s, scan.Whitespace, s.Discard)
	scan.While(s, scan.Digit0F, s.Keep)
	tok := s.Emit()
	fmt.Println(tok.Value)

	// Output:
	// 1010234abc
}
