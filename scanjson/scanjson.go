package scanjson

import "github.com/blackchip-org/scan"

var (
	EscapeRules = []scan.Rule{
		scan.NewCharEncRule(
			scan.BackspaceEnc,
			scan.FormFeedEnc,
			scan.LineFeedEnc,
			scan.CarriageReturnEnc,
			scan.HorizontalTabEnc,
		),
		scan.Hex4Enc,
	}
	Punct = []string{"{", "}", "[", "]", ":", ",", "true", "false", "null"}
)

var (
	String = scan.StrDoubleQuote.WithEscapeRules(EscapeRules...)
	Number = scan.RealExp.
		WithSign(scan.Rune('-')).
		WithLeadingZeroAllowed(false).
		WithEmptyPartsAllowed(false)
	Symbols    = scan.Literal(Punct...)
	Whitespace = scan.Rune(' ', '\n', '\r', '\t')
)

type Context struct {
	RuleSet scan.RuleSet
}

func NewContext() *Context {
	c := &Context{}
	c.RuleSet = scan.NewRuleSet(
		String,
		Number,
		Symbols,
	).WithDiscards(Whitespace)
	return c
}
