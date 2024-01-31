package scan_test

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

// hex integer using predefined class
func Example_example3a() {
	s := scan.NewFromString("", "1010234abc!")
	scan.While(s, scan.Digit0F, s.Keep)
	tok := s.Emit()
	fmt.Println(tok.Value)

	// Output:
	// 1010234abc
}
