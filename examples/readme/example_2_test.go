package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

func Example_example2() {
	isDigit09 := func(r rune) bool {
		return r >= '0' && r <= '9'
	}

	s := scan.NewScannerFromString("", "1010234abc!")
	for s.HasMore() && isDigit09(s.This) {
		s.Keep()
	}
	tok := s.Emit()
	fmt.Println(tok.Val)

	// Output:
	// 1010234
}
