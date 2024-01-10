package scan

type NumRule struct {
	NoLeading Class
	Digit     Class
	DigitSep  Class
	DecSep    Class
	Exp       Class
}

func (r NumRule) Scan(s *Scanner) bool {
	if !s.Is(r.Digit) || s.Is(r.NoLeading) {
		return false
	}
	s.Keep()

	// for {
	// 	switch {
	// 	case s.Is(r.Digit):
	// 		s.Keep()
	// 	case s.Is(r.DigitSep):
	// 		if s.Peek(1)
	// 	}
	// }
}
