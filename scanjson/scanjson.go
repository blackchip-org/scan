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
		scan.Hex4EncRule,
	}
)

var (
	String = scan.StrDoubleQuoteRule.WithEscapeRules(EscapeRules...)
	Number = scan.RealExpRule.
		WithSign(scan.Rune('-')).
		WithLeadingZeroAllowed(false).
		WithEmptyPartsAllowed(false)
	Literals = scan.Literal(
		"{", "}", "[", "]", ":", ",",
		"true", "false", "null",
	)
	Whitespace = scan.NewWhileRule(scan.Rune(' ', '\n', '\r', '\t'), "").WithKeep(false)
)

type Context struct {
	RuleSet scan.RuleSet
}

func NewContext() *Context {
	c := &Context{}
	c.RuleSet = scan.NewRuleSet(
		Whitespace,
		Literals,
		String,
		Number,
	)
	return c
}
