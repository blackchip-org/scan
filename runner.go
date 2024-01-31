package scan

type Runner struct {
	scan      *Scanner
	Rules     RuleSet
	lookahead Token
}

func NewRunner(scan *Scanner, rules RuleSet) *Runner {
	r := &Runner{
		scan:  scan,
		Rules: rules,
	}
	r.lookahead = rules.Next(scan)
	return r
}

func (r *Runner) HasMore() bool {
	return r.lookahead.IsValid()
}

func (r *Runner) Next() Token {
	this := r.lookahead
	r.lookahead = Token{}
	if r.scan.HasMore() {
		r.lookahead = r.Rules.Next(r.scan)
	}
	return this
}

func (r *Runner) Peek() Token {
	return r.lookahead
}

func (r *Runner) All() []Token {
	var toks []Token
	for r.HasMore() {
		toks = append(toks, r.Next())
	}
	return toks
}
