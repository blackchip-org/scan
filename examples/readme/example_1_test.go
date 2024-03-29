package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

func Example_example1() {
	s := scan.NewScannerFromString("", "1010234abc!")
	for s.HasMore() {
		if s.This == '0' || s.This == '1' {
			s.Keep()
			continue
		}
		break
	}
	tok := s.Emit()
	fmt.Println(tok.Val)

	// Output:
	// 1010
}
