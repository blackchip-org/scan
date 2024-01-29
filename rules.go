package scan

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
			tok = s.Illegal("invalid character")
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

type NumRule struct {
	type_    string
	digit    Class
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
	Bin      NumRule = NewNumRule(Digit01).WithType(BinType)
	Hex      NumRule = NewNumRule(Digit0F).WithType(HexType)
	Int      NumRule = UInt.WithSign(Sign)
	Oct      NumRule = NewNumRule(Digit07).WithType(OctType)
	Real     NumRule = UReal.WithSign(Sign)
	RealExp  NumRule = URealExp.WithSign(Sign)
	UInt     NumRule = NewNumRule(Digit09)
	UReal    NumRule = UInt.WithDecSep(Rune('.'))
	URealExp NumRule = UReal.WithExp(Rune('e', 'E')).WithExpSign(Sign)
)
