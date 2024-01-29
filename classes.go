package scan

import "unicode"

// Class returns true if the given rune is a member.
type Class func(rune) bool

// Not returns a Class function that returns true when given a rune that does
// not match class c.
func Not(c Class) Class {
	return func(r0 rune) bool {
		return !c(r0)
	}
}

// Or returns a Class that returns true when given a rune that matches any
// class found in cs.
func Or(cs ...Class) Class {
	return func(r0 rune) bool {
		for _, c := range cs {
			if c(r0) {
				return true
			}
		}
		return false
	}
}

// Range returns a Class function that returns true when given a rune that is
// between from and to inclusive.
func Range(from rune, to rune) Class {
	return func(r0 rune) bool {
		return r0 >= from && r0 <= to
	}
}

// Rune returns a Class function that returns true when given any rune found
// in rs.
func Rune(rs ...rune) Class {
	switch len(rs) {
	case 0:
		return func(rune) bool { return false }
	case 1:
		return func(r0 rune) bool { return r0 == rs[0] }
	}
	return func(r0 rune) bool {
		for _, r := range rs {
			if r == r0 {
				return true
			}
		}
		return false
	}
}

var (
	// Any returns true as long as there is a rune available in the input
	// stream.
	Any Class = func(r rune) bool {
		return r != EndOfText
	}

	// Currency returns true when given a rune that is a currency symbol as
	// defined by Unicode.
	Currency Class = func(r rune) bool {
		return unicode.Is(unicode.Sc, r)
	}

	// Digit returns true when given a digit as defined by Unicode.
	Digit Class = unicode.IsDigit

	// Digit01 returns true when given a valid binary digit.
	Digit01 Class = Rune('0', '1')

	// Digit07 returns true when given a valid octal digit.
	Digit07 Class = Range('0', '7')

	// Digit09 returns true when given a valid decimal digit.
	Digit09 Class = Range('0', '9')

	// Digit0F returns true when given a valid hexadecimal digit.
	Digit0F Class = Or(
		Digit09,
		Range('a', 'f'),
		Range('A', 'F'),
	)

	// Letter returns true when given rune is a letter as defined by Unicode.
	Letter Class = unicode.IsLetter

	// LetterAZ returns true when given letters from the Latin alphabet.
	LetterAZ Class = Or(
		Range('a', 'z'),
		Range('A', 'Z'),
	)

	// LetterUnder returns true when given letters as defined by Unicode
	// or an underscore.
	LetterUnder = Or(
		Letter,
		Rune('_'),
	)

	// LetterDigitUnder returns true when given letters as digits as defined
	// by Unicode or an underscore.
	LetterDigitUnder = Or(
		LetterUnder,
		Digit,
	)

	// None always returns false.
	None Class = func(r rune) bool {
		return false
	}

	// Printable returns true when given a rune that is printable as defined
	// by Unicode.
	Printable Class = unicode.IsPrint

	// Rune8 returns true when given a rune that can be represented by
	// an 8-bit number.
	Rune8 Class = Range(0, 0xff)

	// Rune16 returns true when given a rune that can be represented by
	// a 16-bit number.
	Rune16 Class = Range(0, 0xffff)

	// Sign returns true when the rune is a positive (+) or negative (-)
	// numeric symbol.
	Sign Class = Rune('+', '-')

	// Whitespace returns true when the given rune is whitespace as
	// defined by Unicode.
	Whitespace Class = unicode.IsSpace
)
