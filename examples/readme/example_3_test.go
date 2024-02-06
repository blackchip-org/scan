package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

func Example_example3() {
	s := scan.NewScannerFromString("", "1010234abc!")
	for s.HasMore() && scan.IsDigit09(s.This) {
		s.Keep()
	}
	tok := s.Emit()
	fmt.Println(tok.Val)

	// Output:
	// 1010234
}
