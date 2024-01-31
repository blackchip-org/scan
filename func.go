package scan

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Repeat calls function fn in a loop for n times.
func Repeat(fn func(), n int) {
	for i := 0; i < n; i++ {
		fn()
	}
}

// While calls function fn while this rune is a member of class c.
func While(s *Scanner, c Class, fn func()) {
	for c(s.This) {
		fn()
	}
}

// Until calls function fn until this rune is a member of class c.
func Until(s *Scanner, c Class, fn func()) {
	for s.HasMore() && !c(s.This) {
		fn()
	}
}

func Escape(ch rune) string {
	switch {
	case ch == '\n':
		return "{!ch:\\n}"
	case ch <= 0xff:
		return fmt.Sprintf("{!ch:%02x}", ch)
	case ch <= 0xffff:
		return fmt.Sprintf("{!ch:%04x}", ch)
	default:
		return fmt.Sprintf("{!ch:%08x}", ch)
	}
}

func Quote(s string) string {
	var qs strings.Builder
	for _, ch := range s {
		switch {
		case !unicode.IsPrint(ch) || ch == 0xfffd:
			qs.WriteString(Escape(ch))
		default:
			qs.WriteRune(ch)
		}
	}
	switch {
	case !strings.Contains(s, `"`):
		return `"` + qs.String() + `"`
	case !strings.Contains(s, `'`):
		return `'` + qs.String() + `'`
	case !strings.Contains(s, "`"):
		return "`" + qs.String() + "`"
	default:
		return "{!quote:" + s + "}"
	}
}

func QuoteRune(r rune) string {
	return Quote(string(r))
}

func FormatTokenTable(ts []Token) string {
	var out strings.Builder

	posHeading := "Pos"
	typeHeading := "Type"
	valHeading := "Value"
	litHeading := "Literal"

	posLen := utf8.RuneCountInString(posHeading)
	typeLen := utf8.RuneCountInString(typeHeading)
	valLen := utf8.RuneCountInString(valHeading)
	litLen := utf8.RuneCountInString(litHeading)
	for _, t := range ts {
		posLen = max(posLen, utf8.RuneCountInString(t.Pos.String()))
		typeLen = max(typeLen, utf8.RuneCountInString(t.Type))
		valLen = max(valLen, utf8.RuneCountInString(Quote(t.Value)))
		litLen = max(litLen, utf8.RuneCountInString(Quote(t.Literal)))
	}
	line := fmt.Sprintf("%*s  %-*s  %-*s  %-*s",
		posLen, posHeading,
		typeLen, typeHeading,
		valLen, valHeading,
		litLen, litHeading,
	)
	out.WriteString(strings.TrimRight(line, " "))
	out.WriteRune('\n')
	for _, t := range ts {
		line := fmt.Sprintf("%*s  %-*s  %-*s  %-*s",
			posLen, t.Pos.String(),
			typeLen, t.Type,
			valLen, Quote(t.Value),
			litLen, Quote(t.Literal),
		)
		out.WriteString(strings.TrimRight(line, " "))
		out.WriteRune('\n')
		if len(t.Errs) > 0 {
			fmt.Fprintln(&out, t.Errs)
		}
	}
	return out.String()
}
