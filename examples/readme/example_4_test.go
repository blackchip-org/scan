package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

func Example_example4() {
	s := scan.NewScannerFromString("", "1010234abc!")
	scan.While(s, scan.Digit09, s.Keep)
	tok := s.Emit()
	fmt.Println(tok.Val)

	// Output:
	// 1010234
}
