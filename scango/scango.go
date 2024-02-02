package scango

import (
	"fmt"

	"github.com/blackchip-org/scan"
)

const (
	CommentType = scan.CommentType
	FloatType   = "float"
	IdentType   = scan.IdentType
	IllegalType = scan.IllegalType
	ImagType    = "imag"
	IntType     = scan.IntType
	RuneType    = "rune"
	StringType  = "string"
)

var (
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
	Bin = scan.Bin0b.
		WithIntType(IntType).
		WithDigitSep(scan.Rune('_')).
		WithLeadingDigitSepAllowed(true).
		WithSuffix(ImagSuffix)
	GenComment = scan.NewCommentRule(scan.Literal("/*"), scan.Literal("*/"))
	HexFloat   = scan.NewNumRule(scan.Digit0F).
			WithIntType(IntType).
			WithRealType(FloatType).
			WithPrefix(scan.Literal("0x", "0X")).
			WithDecSep(scan.Rune('.')).
			WithExp(scan.Rune('p', 'P')).
			WithExpSign(scan.Rune('+', '-')).
			WithDigitSep(scan.Rune('_')).
			WithLeadingDigitSepAllowed(true).
			WithSuffix(ImagSuffix)
	Ident    = scan.Ident.WithKeywords(Keywords...)
	IntFloat = scan.RealExp.
			WithIntType(IntType).
			WithRealType(FloatType).
			WithDigitSep(scan.Rune('_')).
			WithSuffix(ImagSuffix)
	LineComment = scan.NewCommentRule(scan.Literal("//"), scan.Literal("\n"))
	Oct         = scan.Oct0o.
			WithIntType(IntType).
			WithDigitSep(scan.Rune('_')).
			WithLeadingDigitSepAllowed(true).
			WithSuffix(ImagSuffix)
	RawString = scan.NewStrRule('`', '`').
			WithType(StringType).
			WithMultiline(true)
	Rune = scan.NewStrRule('\'', '\'').
		WithType(RuneType).
		WithMaxLen(1).
		WithEscape('\\').
		WithEscapeRules(EscapeRules(RuneType)...)
	String = scan.NewStrRule('"', '"').
		WithType(StringType).
		WithEscape('\\').
		WithEscapeRules(EscapeRules(StringType)...)
	Symbols    = scan.Literal(OpsPunct...)
	Whitespace = scan.NewSpaceRule(scan.Rune(' ', '\t', '\r'))
)

func EscapeRules(type_ string) []scan.Rule {
	if type_ != StringType && type_ != RuneType {
		panic(fmt.Sprintf("invalid escape rule type: %v", type_))
	}
	return []scan.Rule{
		scan.NewCharEncRule(
			scan.AlertEnc,
			scan.BackspaceEnc,
			scan.FormFeedEnc,
			scan.LineFeedEnc,
			scan.CarriageReturnEnc,
			scan.HorizontalTabEnc,
			scan.VerticalTabEnc,
		),
		scan.Hex2Enc.AsByte(type_ == StringType),
		scan.Hex4Enc,
		scan.Hex8Enc,
		scan.OctEnc,
	}
}

var semiColonRequiredAfter = map[string]struct{}{
	IdentType:     {},
	IntType:       {},
	FloatType:     {},
	ImagType:      {},
	RuneType:      {},
	StringType:    {},
	"break":       {},
	"continue":    {},
	"fallthrough": {},
	"return":      {},
	"++":          {},
	"--":          {},
	")":           {},
	"]":           {},
	"}":           {},
}

func AutoSemiInsertion() func(*scan.Scanner, scan.Token) scan.Token {
	var last scan.Token
	return func(_ *scan.Scanner, t scan.Token) scan.Token {
		if t.Type == "\n" {
			if _, yes := semiColonRequiredAfter[last.Type]; yes {
				t.Type = ";"
				t.Val = ";"
			} else {
				t = scan.Token{}
			}
		}
		last = t
		return t
	}
}

type ImagSuffixRule struct{}

func (r ImagSuffixRule) Eval(s *scan.Scanner) bool {
	if s.Is(scan.Rune('i')) {
		s.Keep()
		s.Type = ImagType
		return true
	}
	return false
}

var ImagSuffix = ImagSuffixRule{}

type Context struct {
	KeepComments bool
	RuleSet      scan.RuleSet
}

func NewContext() *Context {
	c := &Context{}
	c.RuleSet = scan.NewRuleSet(
		Whitespace,
		GenComment.WithKeep(&c.KeepComments),
		LineComment.WithKeep(&c.KeepComments),
		Rune,
		HexFloat,
		Oct,
		Bin,
		IntFloat,
		String,
		RawString,
		Ident,
		Symbols,
	).WithPostTokenFunc(AutoSemiInsertion())
	return c
}
