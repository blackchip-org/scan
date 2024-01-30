package scan

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

type Rule interface {
	Eval(*Scanner) bool
}

type RuleSet struct {
	rules         []Rule
	discards      Class
	preTokenFunc  func(*Scanner)
	postTokenFunc func(*Scanner, Token) Token
}

func Rules(rules ...Rule) RuleSet {
	return RuleSet{
		rules:    rules,
		discards: Whitespace,
	}
}

func (r RuleSet) WithDiscards(c Class) RuleSet {
	r.discards = c
	return r
}

func (r RuleSet) WithPreTokenFunc(fn func(*Scanner)) RuleSet {
	r.preTokenFunc = fn
	return r
}

func (r RuleSet) WithPostTokenFunc(fn func(*Scanner, Token) Token) RuleSet {
	r.postTokenFunc = fn
	return r
}

func (r RuleSet) Eval(s *Scanner) bool {
	for _, r := range r.rules {
		if r.Eval(s) {
			return true
		}
	}
	return false
}

func (r RuleSet) Next(s *Scanner) Token {
	var tok Token
	for s.HasMore() {
		While(s, r.discards, s.Discard)
		if r.preTokenFunc != nil {
			r.preTokenFunc(s)
		}

		for _, rule := range r.rules {
			if rule.Eval(s) {
				tok = s.Emit()
				break
			}
		}
		if tok.Value == "" {
			s.Illegal("unexpected %s", QuoteRune(s.This))
			tok = s.Emit()
		}
		if r.postTokenFunc != nil {
			tok = r.postTokenFunc(s, tok)
		}
		if tok.Value != "" && tok.Type != "" {
			break
		}
	}
	While(s, r.discards, s.Discard)
	return tok
}

type CharEnc struct {
	From rune
	To   rune
}

func NewCharEnc(from rune, to rune) CharEnc {
	return CharEnc{from, to}
}

var (
	AlertEnc          CharEnc = NewCharEnc('a', '\a')
	BackspaceEnc      CharEnc = NewCharEnc('b', '\b')
	BellEnc           CharEnc = AlertEnc
	CarriageReturnEnc CharEnc = NewCharEnc('r', '\r')
	FormFeedEnc       CharEnc = NewCharEnc('f', '\f')
	HorizontalTabEnc  CharEnc = NewCharEnc('t', '\t')
	LineFeedEnc       CharEnc = NewCharEnc('n', '\n')
	NewlineEnc        CharEnc = LineFeedEnc
	TabEnc            CharEnc = HorizontalTabEnc
	VerticalTabEnc    CharEnc = NewCharEnc('v', '\v')
)

type CharEncRule struct {
	charmap map[rune]rune
}

func NewCharEncRule(maps ...CharEnc) CharEncRule {
	r := CharEncRule{charmap: make(map[rune]rune)}
	for _, cm := range maps {
		r.charmap[cm.From] = cm.To
	}
	return r
}

func (r CharEncRule) Eval(s *Scanner) bool {
	to, ok := r.charmap[s.This]
	if !ok {
		return false
	}
	s.Literal.WriteRune(s.This)
	s.Value.WriteRune(to)
	s.Skip()
	return true
}

type HexEncRule struct {
	flag   rune
	digits int
}

func NewHexEncRule(flag rune, digits int) HexEncRule {
	if digits != 2 && digits != 4 && digits != 8 {
		panic(fmt.Sprintf("invalid digits '%v' for hex encoding map", digits))
	}
	return HexEncRule{flag: flag, digits: digits}
}

func (r HexEncRule) Eval(s *Scanner) bool {
	if s.This != r.flag {
		return false
	}
	s.Skip()
	digits := make([]rune, r.digits)
	for i := 0; i < r.digits; i++ {
		digits[i] = s.Peek(i)
	}
	val, err := strconv.ParseUint(string(digits), 16, 32)
	if err != nil {
		s.Illegal("invalid encoding: '%v'", string(digits))
		return true
	}
	ch := rune(val)
	if !utf8.ValidRune(ch) {
		s.Illegal("character not valid: '%v'", string(digits))
		return true
	}
	Repeat(s.Skip, r.digits)
	s.Value.WriteRune(ch)
	return true
}

var (
	Hex2Enc = NewHexEncRule('x', 2)
	Hex4Enc = NewHexEncRule('u', 4)
	Hex8Enc = NewHexEncRule('U', 8)
)

type IdentRule struct {
	head     Class
	tail     Class
	keywords map[string]struct{}
}

func NewIdentRule(head Class, tail Class) IdentRule {
	return IdentRule{
		head:     head,
		tail:     tail,
		keywords: make(map[string]struct{}),
	}
}

func (r IdentRule) WithKeywords(ks ...string) IdentRule {
	keywords := make(map[string]struct{})
	for _, k := range ks {
		keywords[k] = struct{}{}
	}
	r.keywords = keywords
	return r
}

func (r IdentRule) Eval(s *Scanner) bool {
	if !s.Is(r.head) {
		return false
	}
	s.Keep()
	While(s, r.tail, s.Keep)
	if _, ok := r.keywords[s.Value.String()]; !ok {
		s.Type = IdentType
	}
	return true
}

var (
	Ident IdentRule = NewIdentRule(LetterUnder, LetterDigitUnder)
)

type trieNode struct {
	children map[rune]*trieNode
	leaf     bool
}

type LiteralRule struct {
	lits *trieNode
	skip bool
}

func Literals(lits ...string) LiteralRule {
	rule := LiteralRule{lits: &trieNode{children: make(map[rune]*trieNode)}}
	for _, lit := range lits {
		prev := rule.lits
		var node *trieNode
		var ok bool
		for _, ch := range lit {
			node, ok = prev.children[ch]
			if !ok {
				node = &trieNode{children: make(map[rune]*trieNode)}
				prev.children[ch] = node
			}
			prev = node
		}
		node.leaf = true
	}
	return rule
}

func (r LiteralRule) WithSkip(b bool) LiteralRule {
	r.skip = b
	return r
}

func (r LiteralRule) Eval(s *Scanner) bool {
	i := 0
	good := 0
	prev := r.lits
	fn := s.Keep
	if r.skip {
		fn = s.Skip
	}

	for {
		ch := s.Peek(i)
		node, ok := prev.children[ch]
		if !ok {
			Repeat(fn, good)
			return good > 0
		}
		if node.leaf {
			good = i + 1
		}
		prev = node
		i++
	}
}

type NumRule struct {
	type_    string
	digit    Class
	prefix   Rule
	sign     Class
	digitSep Class
	decSep   Class
	exp      Class
	expSign  Class
}

func NewNumRule(digit Class) NumRule {
	return NumRule{
		digit:    digit,
		sign:     None,
		prefix:   TrueRule,
		digitSep: None,
		decSep:   None,
		exp:      None,
		expSign:  None,
	}
}

func (r NumRule) WithType(t string) NumRule {
	r.type_ = t
	return r
}

func (r NumRule) WithSign(c Class) NumRule {
	r.sign = c
	return r
}

func (r NumRule) WithPrefix(rule Rule) NumRule {
	r.prefix = rule
	return r
}

func (r NumRule) WithDigitSep(c Class) NumRule {
	r.digitSep = c
	return r
}

func (r NumRule) WithDecSep(c Class) NumRule {
	r.decSep = c
	return r
}

func (r NumRule) WithExp(c Class) NumRule {
	r.exp = c
	return r
}

func (r NumRule) WithExpSign(c Class) NumRule {
	r.expSign = c
	return r
}

func (r NumRule) Eval(s *Scanner) bool {
	if s.Is(r.sign) {
		s.Keep()
	}
	if !r.prefix.Eval(s) {
		s.Undo()
		return false
	}

	s.Type = IntType
	seenDigit := false
	scanDigits := func() {
		for s.HasMore() {
			switch {
			case s.Is(r.digit):
				seenDigit = true
				s.Keep()
			case s.Is(r.digitSep) && s.NextIs(r.digit) && r.digit(s.Peek(-1)):
				s.Skip()
			default:
				return
			}
		}
	}

	scanDigits()
	if s.Is(r.decSep) {
		s.Type = RealType
		s.Keep()
		scanDigits()
	}

	if !seenDigit {
		s.Undo()
		return false
	}

	if s.Is(r.exp) {
		s.Type = RealType
		s.Keep()
		if s.Is(r.expSign) {
			s.Keep()
		}
		scanDigits()
	}
	if r.type_ != "" {
		s.Type = r.type_
	}
	return true
}

var (
	Bin           NumRule = NewNumRule(Digit01).WithType(BinType)
	Bin0b         NumRule = Bin.WithPrefix(Literals("0b", "0B").WithSkip(true))
	Hex           NumRule = NewNumRule(Digit0F).WithType(HexType)
	Hex0x         NumRule = Hex.WithPrefix(Literals("0x", "0X").WithSkip(true))
	Oct           NumRule = NewNumRule(Digit07).WithType(OctType)
	Oct0o         NumRule = Oct.WithPrefix(Literals("0o", "0O").WithSkip(true))
	Int           NumRule = NewNumRule(Digit09)
	Real          NumRule = Int.WithDecSep(Rune('.'))
	RealExp       NumRule = Real.WithExp(Rune('e', 'E')).WithExpSign(Sign)
	SignedInt     NumRule = Int.WithSign(Sign)
	SignedReal    NumRule = Real.WithSign(Sign)
	SignedRealExp NumRule = RealExp.WithSign(Sign)
)

type OctEncRule struct {
}

func NewOctEncRule() OctEncRule {
	return OctEncRule{}
}

func (r OctEncRule) Eval(s *Scanner) bool {
	if !s.Is(Digit07) {
		return false
	}
	digits := make([]rune, 3)
	for i := 0; i < 3; i++ {
		digits = append(digits, s.Peek(i))
	}
	val, err := strconv.ParseUint(string(digits), 8, 8)
	if err != nil {
		s.Illegal("invalid encoding: '%v'", string(digits))
		return true
	}
	Repeat(s.Skip, 3)
	s.Value.WriteRune(rune(val))
	return true
}

var OctEnc = NewOctEncRule()

type StrRule struct {
	type_       string
	begin       rune
	end         rune
	escape      rune
	escapeRules RuleSet
	multiline   bool
}

func NewStrRule(begin rune, end rune) StrRule {
	return StrRule{
		begin:       begin,
		end:         end,
		escapeRules: Rules(),
	}
}

func (r StrRule) WithType(t string) StrRule {
	r.type_ = t
	return r
}

func (r StrRule) WithEscape(rn rune) StrRule {
	r.escape = rn
	return r
}

func (r StrRule) WithEscapeRules(rules ...Rule) StrRule {
	r.escapeRules = Rules(rules...)
	return r
}

func (r StrRule) WithMultiline(b bool) StrRule {
	r.multiline = b
	return r
}

func (r StrRule) Eval(s *Scanner) bool {
	if s.This != r.begin {
		return false
	}
	s.Skip()

	s.Type = StrType
	if r.type_ != "" {
		s.Type = r.type_
	}

	for s.HasMore() {
		switch {
		case s.This == r.end:
			s.Skip()
			return true
		case s.This == '\n':
			if !r.multiline {
				s.Illegal("not terminated")
				return true
			}
			s.Keep()
		case s.This == r.escape:
			s.Skip()
			if s.This == r.end || s.This == r.escape {
				s.Keep()
			} else if !r.escapeRules.Eval(s) {
				s.Illegal("invalid escape sequence: '%v%v'", r.escape, s.This)
			}
		default:
			s.Keep()
		}
	}
	s.Illegal("not terminated")
	return true
}

var (
	StrDoubleQuote = NewStrRule('"', '"').WithEscape('\\')
	StrSingleQuote = NewStrRule('\'', '\'').WithEscape('\\')
)

type trueRule struct{}

func (r trueRule) Eval(s *Scanner) bool {
	return true
}

var TrueRule = trueRule{}
