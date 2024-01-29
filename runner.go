package scan

type Runner struct {
	scan      *Scanner
	rules     RuleSet
	lookahead Token
}

func NewRunner(scan *Scanner, rules RuleSet) *Runner {
	r := &Runner{
		scan:  scan,
		rules: rules,
	}
	r.lookahead = rules.Next(scan)
	return r
}

func (r *Runner) More() bool {
	return len(r.lookahead.Type) > 0
}

func (r *Runner) Next() Token {
	this := r.lookahead
	r.lookahead = r.rules.Next(r.scan)
	return this
}

func (r *Runner) Peek() Token {
	return r.lookahead
}

func (r *Runner) All() []Token {
	var toks []Token
	for r.More() {
		toks = append(toks, r.Next())
	}
	return toks
}
