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
	preTokenFunc  func(*Scanner)
	postTokenFunc func(*Scanner, Token) Token
	noMatchFunc   func(*Scanner)
}

func NewRuleSet(rules ...Rule) RuleSet {
	return RuleSet{
		rules:       rules,
		noMatchFunc: UnexpectedRune(),
	}
}

func (r RuleSet) WithPreTokenFunc(fn func(*Scanner)) RuleSet {
	r.preTokenFunc = fn
	return r
}

func (r RuleSet) WithPostTokenFunc(fn func(*Scanner, Token) Token) RuleSet {
	r.postTokenFunc = fn
	return r
}

func (r RuleSet) WithNoMatchFunc(fn func(*Scanner)) RuleSet {
	r.noMatchFunc = fn
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
	start := s.Pos

	for {
		if r.preTokenFunc != nil {
			r.preTokenFunc(s)
		}

		match := false
		for _, rule := range r.rules {
			if match = rule.Eval(s); match {
				tok = s.Emit()
				break
			}
		}
		if !match {
			r.noMatchFunc(s)
			return s.Emit()
		}
		if r.postTokenFunc != nil {
			tok = r.postTokenFunc(s, tok)
		}
		if tok.IsValid() {
			break
		}
		if s.Pos == start {
			panic("scanner did not advance")
		}
		start = s.Pos
	}
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
	s.Lit.WriteRune(s.This)
	s.Val.WriteRune(to)
	s.Skip()
	return true
}

type ClassRule struct {
	isClass Class
	type_   string
}

func NewClassRule(isClass Class) ClassRule {
	return ClassRule{isClass: isClass}
}

func (r ClassRule) WithType(t string) ClassRule {
	r.type_ = t
	return r
}

func (r ClassRule) Eval(s *Scanner) bool {
	if !r.isClass(s.This) {
		return false
	}
	s.Keep()
	if r.type_ != "" {
		s.Type = r.type_
	}
	return true
}

type CommentRule struct {
	begin LiteralRule
	end   LiteralRule
	keep  *bool
}

func NewCommentRule(begin LiteralRule, end LiteralRule) CommentRule {
	var keep bool
	return CommentRule{
		begin: begin.WithSkip(true),
		end:   end.WithSkip(true),
		keep:  &keep,
	}
}

func (r CommentRule) WithKeep(keep *bool) CommentRule {
	r.keep = keep
	return r
}

func (r CommentRule) Eval(s *Scanner) bool {
	if !r.begin.Eval(s) {
		return false
	}
	var action func()
	if *r.keep {
		s.Type = CommentType
		action = s.Keep
	} else {
		action = s.Skip
	}

	UntilRule(s, r.end, action)
	return true
}

type HexEncRule struct {
	flag   rune
	digits int
	asByte bool
}

func NewHexEncRule(flag rune, digits int) HexEncRule {
	if digits != 2 && digits != 4 && digits != 8 {
		panic(fmt.Sprintf("invalid digits '%v' for hex encoding map", digits))
	}
	return HexEncRule{flag: flag, digits: digits}
}

func (r HexEncRule) AsByte(b bool) HexEncRule {
	if b && r.digits != 2 {
		panic("digits must be 2 for encoding bytes")
	}
	r.asByte = b
	return r
}

func (r HexEncRule) Eval(s *Scanner) bool {
	if s.This != r.flag {
		return false
	}
	s.Skip()
	digits := make([]rune, r.digits)
	for i := 0; i < r.digits; i++ {
		ch := s.Peek(i)
		if !IsDigit0F(ch) {
			s.Illegal("invalid encoding: %v", Quote(string(digits[:i])))
			Repeat(s.Keep, i)
			return true
		}
		digits[i] = s.Peek(i)
	}
	val, err := strconv.ParseUint(string(digits), 16, 32)
	if err != nil {
		s.Illegal("invalid encoding: %v", Quote(string(digits)))
		Repeat(s.Keep, r.digits)
		return true
	}
	if r.asByte {
		s.Val.WriteByte(uint8(val))
	} else {
		ch := rune(val)
		if !utf8.ValidRune(ch) {
			s.Illegal("invalid encoding: %v", Quote(string(digits)))
			Repeat(s.Keep, r.digits)
			return true
		}
		s.Val.WriteRune(ch)
	}
	Repeat(s.Skip, r.digits)
	return true
}

var (
	Hex2EncRule = NewHexEncRule('x', 2)
	Hex4EncRule = NewHexEncRule('u', 4)
	Hex8EncRule = NewHexEncRule('U', 8)
)

type IdentRule struct {
	isHead   Class
	isTail   Class
	keywords map[string]struct{}
}

func NewIdentRule(isHead Class, isTail Class) IdentRule {
	return IdentRule{
		isHead:   isHead,
		isTail:   isTail,
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
	if !r.isHead(s.This) {
		return false
	}
	s.Keep()
	While(s, r.isTail, s.Keep)
	if _, ok := r.keywords[s.Val.String()]; !ok {
		s.Type = IdentType
	}
	return true
}

var (
	StandardIdentRule IdentRule = NewIdentRule(IsLetterUnder, IsLetterDigitUnder)
)

type trieNode struct {
	children map[rune]*trieNode
	leaf     bool
}

type LiteralRule struct {
	lits *trieNode
	skip bool
}

func Literal(lits ...string) LiteralRule {
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
	intType                string
	realType               string
	isDigit                Class
	prefixRule             Rule
	isSign                 Class
	isDigitSep             Class
	isDecSep               Class
	isExp                  Class
	isExpSign              Class
	suffixRules            []Rule
	leadingDigitSepAllowed bool
	leadingZeroAllowed     bool
	emptyPartsAllowed      bool
}

func NewNumRule(digit Class) NumRule {
	return NumRule{
		intType:            IntType,
		realType:           RealType,
		isDigit:            digit,
		prefixRule:         TrueRule,
		isSign:             IsNone,
		isDigitSep:         IsNone,
		isDecSep:           IsNone,
		isExp:              IsNone,
		isExpSign:          IsNone,
		leadingZeroAllowed: true,
		emptyPartsAllowed:  true,
	}
}

func (r NumRule) WithIntType(t string) NumRule {
	r.intType = t
	return r
}

func (r NumRule) WithRealType(t string) NumRule {
	r.realType = t
	return r
}

func (r NumRule) WithSign(c Class) NumRule {
	r.isSign = c
	return r
}

func (r NumRule) WithPrefix(rule Rule) NumRule {
	r.prefixRule = rule
	return r
}

func (r NumRule) WithDigitSep(c Class) NumRule {
	r.isDigitSep = c
	return r
}

func (r NumRule) WithDecSep(c Class) NumRule {
	r.isDecSep = c
	return r
}

func (r NumRule) WithExp(c Class) NumRule {
	r.isExp = c
	return r
}

func (r NumRule) WithExpSign(c Class) NumRule {
	r.isExpSign = c
	return r
}

func (r NumRule) WithSuffix(rules ...Rule) NumRule {
	r.suffixRules = rules
	return r
}

func (r NumRule) WithLeadingDigitSepAllowed(b bool) NumRule {
	r.leadingDigitSepAllowed = b
	return r
}

func (r NumRule) WithLeadingZeroAllowed(b bool) NumRule {
	r.leadingZeroAllowed = false
	return r
}

func (r NumRule) WithEmptyPartsAllowed(b bool) NumRule {
	r.emptyPartsAllowed = b
	return r
}

func (r NumRule) Eval(s *Scanner) bool {
	if r.isSign(s.This) {
		s.Keep()
	}
	if !r.prefixRule.Eval(s) {
		s.Undo()
		return false
	}
	if r.isDigitSep(s.This) && r.leadingDigitSepAllowed {
		s.Skip()
	}

	s.Type = IntType
	if r.intType != "" {
		s.Type = r.intType
	}

	if !r.leadingZeroAllowed && s.This == '0' {
		if !r.isDecSep(s.Next) && !r.isExp(s.Next) {
			s.Keep()
			return true
		}
	}

	seenDigit := false
	scanDigits := func() {
		for s.HasMore() {
			switch {
			case r.isDigit(s.This):
				seenDigit = true
				s.Keep()
			case r.isDigitSep(s.This) && r.isDigit(s.Next) && r.isDigit(s.Peek(-1)):
				s.Skip()
			default:
				return
			}
		}
	}

	scanDigits()
	if r.isDecSep(s.This) {
		if !r.emptyPartsAllowed {
			if !seenDigit {
				s.Undo()
				return false
			}
			if !r.isDigit(s.Next) {
				return true
			}
		}
		s.Type = RealType
		if r.realType != "" {
			s.Type = r.realType
		}
		s.Keep()
		scanDigits()
	}

	if !seenDigit {
		s.Undo()
		return false
	}

	if r.isExp(s.This) && (r.isDigit(s.Next) ||
		(r.isExpSign(s.Next)) && r.isDigit(s.Peek(2))) {
		s.Type = RealType
		if r.realType != "" {
			s.Type = r.realType
		}
		s.Keep()
		if r.isExpSign(s.This) {
			s.Keep()
		}
		scanDigits()
	}

	for _, r := range r.suffixRules {
		if r.Eval(s) {
			break
		}
	}
	return true
}

var (
	BinRule           NumRule = NewNumRule(IsDigit01).WithIntType(BinType)
	Bin0bRule         NumRule = BinRule.WithPrefix(Literal("0b", "0B"))
	HexRule           NumRule = NewNumRule(IsDigit0F).WithIntType(HexType)
	Hex0xRule         NumRule = HexRule.WithPrefix(Literal("0x", "0X"))
	OctRule           NumRule = NewNumRule(IsDigit07).WithIntType(OctType)
	Oct0oRule         NumRule = OctRule.WithPrefix(Literal("0o", "0O"))
	IntRule           NumRule = NewNumRule(IsDigit09)
	RealRule          NumRule = IntRule.WithDecSep(Rune('.'))
	RealExpRule       NumRule = RealRule.WithExp(Rune('e', 'E')).WithExpSign(IsSign)
	SignedIntRule     NumRule = IntRule.WithSign(IsSign)
	SignedRealRule    NumRule = RealRule.WithSign(IsSign)
	SignedRealExpRule NumRule = RealExpRule.WithSign(IsSign)
)

type OctEncRule struct {
}

func NewOctEncRule() OctEncRule {
	return OctEncRule{}
}

func (r OctEncRule) Eval(s *Scanner) bool {
	if !IsDigit07(s.This) {
		return false
	}
	digits := make([]rune, 3)
	for i := 0; i < 3; i++ {
		ch := s.Peek(i)
		if !IsDigit07(ch) {
			s.Illegal("invalid encoding: %v", Quote(string(digits[:i])))
			Repeat(s.Keep, i)
			return true
		}
		digits[i] = ch
	}
	val, err := strconv.ParseUint(string(digits), 8, 8)
	if err != nil {
		s.Illegal("invalid encoding: %v", Quote(string(digits)))
		Repeat(s.Keep, 3)
		return true
	}
	Repeat(s.Skip, 3)
	s.Val.WriteRune(rune(val))
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
	maxLen      uint
	optTerm     bool
	nesting     bool
}

func NewStrRule(begin rune, end rune) StrRule {
	return StrRule{
		begin:       begin,
		end:         end,
		escapeRules: NewRuleSet(),
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
	r.escapeRules = NewRuleSet(rules...)
	return r
}

func (r StrRule) WithMultiline(b bool) StrRule {
	r.multiline = b
	return r
}

func (r StrRule) WithMaxLen(l uint) StrRule {
	r.maxLen = l
	return r
}

func (r StrRule) WithOptionalTerminator(t bool) StrRule {
	r.optTerm = t
	return r
}

func (r StrRule) WithNesting(t bool) StrRule {
	if r.begin == r.end {
		panic("nesting not possible")
	}
	r.nesting = true
	return r
}

func (r StrRule) recover(s *Scanner) {
	Until(s, Rune(r.end), s.Skip)
	s.Skip()
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

	length := uint(0)
	level := 1
	for s.HasMore() {
		if r.nesting && s.This == r.begin {
			level++
		}
		if s.This == r.end {
			level--
			if !r.nesting || level == 0 {
				s.Skip()
				return true
			}
		}

		length++
		if r.maxLen > 0 && length > r.maxLen {
			s.Illegal("too many characters (%v)", length)
			r.recover(s)
			return true
		}

		switch {
		case s.This == '\n':
			if !r.multiline {
				if !r.optTerm {
					s.Illegal("not terminated")
				}
				return true
			}
			s.Keep()
		case s.This == r.escape:
			s.Skip()
			if s.This == r.end || s.This == r.escape {
				s.Keep()
			} else if !r.escapeRules.Eval(s) {
				s.Illegal("invalid escape sequence: '%c%c'", r.escape, s.This)
				s.Keep()
				r.recover(s)
				return true
			} else if len(s.Errs) > 0 {
				r.recover(s)
				return true
			}
		default:
			s.Keep()
		}
	}
	if !r.optTerm {
		s.Illegal("not terminated")
	}
	return true
}

var (
	StrDoubleQuoteRule = NewStrRule('"', '"').WithEscape('\\')
	StrSingleQuoteRule = NewStrRule('\'', '\'').WithEscape('\\')
)

type trueRule struct{}

func (r trueRule) Eval(s *Scanner) bool {
	return true
}

var TrueRule = trueRule{}

type WhileRule struct {
	keep    bool
	isClass Class
	type_   string
}

func NewWhileRule(space Class, type_ string) WhileRule {
	return WhileRule{isClass: space, type_: type_, keep: true}
}

func (r WhileRule) WithKeep(b bool) WhileRule {
	r.keep = b
	return r
}

func (r WhileRule) Eval(s *Scanner) bool {
	if !r.isClass(s.This) {
		return false
	}

	var action func()
	if r.keep {
		s.Type = r.type_
		action = s.Keep
	} else {
		action = s.Discard
	}

	While(s, r.isClass, action)
	return true
}

var (
	SkipSpaceRule = NewWhileRule(IsSpace, SpaceType).WithKeep(false)
	KeepSpaceRule = NewWhileRule(IsSpace, SpaceType)
	WordRule      = NewWhileRule(Not(IsSpace), WordType)
)
