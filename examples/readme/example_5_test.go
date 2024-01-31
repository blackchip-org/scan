package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

// discarding whitespace
func Example_example5() {
	s := scan.NewScannerFromString("", " \t 1010234abc!")
	scan.While(s, scan.Whitespace, s.Discard)
	scan.While(s, scan.Digit09, s.Keep)
	tok := s.Emit()
	fmt.Println(tok.Val)

	// Output:
	// 1010234
}
