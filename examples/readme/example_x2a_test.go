package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

// replace loop with scan.While
func Example_example2a() {
	digit09 := func(r rune) bool {
		return r >= '0' && r <= '9'
	}

	s := scan.NewScannerFromString("", "1010234abc!")
	scan.While(s, digit09, s.Keep)
	tok := s.Emit()
	fmt.Println(tok.Val)

	// Output:
	// 1010234
}
