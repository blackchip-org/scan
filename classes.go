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
	// IsAny returns true as long as there is a rune available in the input
	// stream.
	IsAny Class = func(r rune) bool {
		return r != EndOfText
	}

	// IsCurrency returns true when given a rune that is a currency symbol as
	// defined by Unicode.
	IsCurrency Class = func(r rune) bool {
		return unicode.Is(unicode.Sc, r)
	}

	// IsDigit returns true when given a digit as defined by Unicode.
	IsDigit Class = unicode.IsDigit

	// IsDigit01 returns true when given a valid binary digit.
	IsDigit01 Class = Rune('0', '1')

	// IsDigit07 returns true when given a valid octal digit.
	IsDigit07 Class = Range('0', '7')

	// IsDigit09 returns true when given a valid decimal digit.
	IsDigit09 Class = Range('0', '9')

	// IsDigit0F returns true when given a valid hexadecimal digit.
	IsDigit0F Class = Or(
		IsDigit09,
		Range('a', 'f'),
		Range('A', 'F'),
	)

	// IsLetter returns true when given rune is a letter as defined by Unicode.
	IsLetter Class = unicode.IsLetter

	// IsLetterAZ returns true when given letters from the Latin alphabet.
	IsLetterAZ Class = Or(
		Range('a', 'z'),
		Range('A', 'Z'),
	)

	// IsLetterUnder returns true when given letters as defined by Unicode
	// or an underscore.
	IsLetterUnder = Or(
		IsLetter,
		Rune('_'),
	)

	// IsLetterDigitUnder returns true when given letters as digits as defined
	// by Unicode or an underscore.
	IsLetterDigitUnder = Or(
		IsLetterUnder,
		IsDigit,
	)

	// IsNone always returns false.
	IsNone Class = func(r rune) bool {
		return false
	}

	// IsPrintable returns true when given a rune that is printable as defined
	// by Unicode.
	IsPrintable Class = unicode.IsPrint

	// IsRune8 returns true when given a rune that can be represented by
	// an 8-bit number.
	IsRune8 Class = Range(0, 0xff)

	// IsRune16 returns true when given a rune that can be represented by
	// a 16-bit number.
	IsRune16 Class = Range(0, 0xffff)

	// IsSign returns true when the rune is a positive (+) or negative (-)
	// numeric symbol.
	IsSign Class = Rune('+', '-')

	// IsWhitespace returns true when the given rune is whitespace as
	// defined by Unicode.
	IsWhitespace Class = unicode.IsSpace
)
