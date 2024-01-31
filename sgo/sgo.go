package sgo

import "github.com/blackchip-org/scan"

const (
	BinType     = scan.BinType
	CommentType = scan.CommentType
	FloatType   = "float"
	HexType     = scan.HexType
	IdentType   = scan.IdentType
	IntType     = scan.IntType
	ImagType    = "imag"
	OctType     = scan.OctType
	RuneType    = "rune"
	StringType  = "string"
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
	Keywords = []string{
		"break", "case", "chan", "const", "continue",
		"default", "defer", "else", "fallthrough", "for",
		"func", "go", "goto", "if", "import",
		"interface", "map", "package", "range", "return",
		"select", "struct", "switch", "type", "var",
	}
	OpsPunct = []string{
		"+", "&", "+=", "&=", "&&", "==", "!=", "(", ")",
		"-", "|", "-=", "|=", "||", "<", "<=", "[", "]",
		"*", "^", "*=", "^=", "<-", ">", ">=", "{", "}",
		"/", "<<", "/=", "<<=", "++", "=", ":=", ",", ";",
		"%", ">>", "%=", ">>=", "--", "!", "...", ".", ":",
		"&^", "&^=", "~",
		"\n",
	}
)

var (
	GeneralComment = scan.NewCommentRule(scan.Literals("/*"), scan.Literals("*/"))
	Ident          = scan.Ident.WithKeywords(Keywords...)
	Int            = scan.Rules(
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
	LineComment = scan.NewCommentRule(scan.Literals("//"), scan.Literals("\n"))
	RawString   = scan.NewStrRule('`', '`').
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
	Symbols    = scan.Literals(OpsPunct...)
	Whitespace = scan.Rune(' ', '\t', '\r')
)

var semiColonRequiredAfter = map[string]struct{}{
	IdentType:  {},
	IntType:    {},
	BinType:    {},
	OctType:    {},
	HexType:    {},
	FloatType:  {},
	ImagType:   {},
	RuneType:   {},
	StringType: {},
	"++":       {},
	"--":       {},
	")":        {},
	"]":        {},
	"}":        {},
}

func AutoSemiInsertion() func(*scan.Scanner, scan.Token) scan.Token {
	var last scan.Token
	return func(_ *scan.Scanner, t scan.Token) scan.Token {
		if t.Type == "\n" {
			if _, yes := semiColonRequiredAfter[last.Type]; yes {
				t.Type = ";"
				t.Value = ";"
			} else {
				t = scan.Token{Type: scan.EmptyType}
			}
		}
		last = t
		return t
	}
}

type Context struct {
	KeepComments bool
	RuleSet      scan.RuleSet
}

func NewContext() *Context {
	c := &Context{}
	c.RuleSet = scan.Rules(
		GeneralComment.WithKeep(&c.KeepComments),
		LineComment.WithKeep(&c.KeepComments),
		Symbols,
		Rune,
		Int,
		String,
		RawString,
		Ident,
	).
		WithDiscards(Whitespace).
		WithPostTokenFunc(AutoSemiInsertion())
	return c
}
