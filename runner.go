package scan

type Runner struct {
	scan  *Scanner
	Rules RuleSet
	This  Token
	Next  Token
}

func NewRunner(scan *Scanner, rules RuleSet) *Runner {
	r := &Runner{
		scan:  scan,
		Rules: rules,
	}
	r.This = rules.Next(scan)
	r.Next = rules.Next(scan)
	return r
}

func (r *Runner) HasMore() bool {
	return !r.This.IsEndOfText()
}

func (r *Runner) Scan() Token {
	if r.This.IsEndOfText() {
		return r.This
	}
	r.This = r.Next
	r.Next = r.Rules.Next(r.scan)
	return r.This
}

func (r *Runner) All() []Token {
	var toks []Token
	for r.HasMore() {
		toks = append(toks, r.This)
		r.Scan()
	}
	return toks
}
