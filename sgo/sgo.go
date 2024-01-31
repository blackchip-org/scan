package sgo

import "github.com/blackchip-org/scan"

const (
	BinType    = "bin"
	FloatType  = "float"
	HexType    = "hex"
	IdentType  = "ident"
	IntType    = "int"
	OctType    = "oct"
	RuneType   = "rune"
	StringType = "string"
)

var (
	EscapeRules = []scan.Rule{
		scan.NewCharEncRule(
			scan.AlertEnc,
			scan.BackspaceEnc,
			scan.FormFeedEnc,
			scan.LineFeedEnc,
			scan.CarriageReturnEnc,
			scan.HorizontalTabEnc,
			scan.VerticalTabEnc,
		),
		scan.Hex2Enc,
		scan.Hex4Enc,
		scan.Hex8Enc,
		scan.OctEnc,
	}
	Ident = scan.Ident.WithKeywords(Keywords...)
	Int   = scan.Rules(
		scan.Hex0x.
			WithIntType(HexType).
			WithDigitSep(scan.Rune('_')).
			WithLeadingDigitSepAllowed(true),
		scan.Oct0o.
			WithIntType(OctType).
			WithDigitSep(scan.Rune('_')).
			WithLeadingDigitSepAllowed(true),
		scan.Bin0b.
			WithIntType(BinType).
			WithDigitSep(scan.Rune('_')).
			WithLeadingDigitSepAllowed(true),
		scan.RealExp.
			WithIntType(IntType).
			WithRealType(FloatType).
			WithDigitSep(scan.Rune('_')),
	)
	Keywords = []string{
		"break", "case", "chan", "const", "continue",
		"default", "defer", "else", "fallthrough", "for",
		"func", "go", "goto", "if", "import",
		"interface", "map", "package", "range", "return",
		"select", "struct", "switch", "type", "var",
	}
	RawString = scan.NewStrRule('`', '`').
			WithType(StringType).
			WithMultiline(true)
	Rune = scan.NewStrRule('\'', '\'').
		WithType(RuneType).
		WithMaxLen(1).
		WithEscape('\\').
		WithEscapeRules(EscapeRules...)
	String = scan.NewStrRule('"', '"').
		WithType(StringType).
		WithEscape('\\').
		WithEscapeRules(EscapeRules...)
)
