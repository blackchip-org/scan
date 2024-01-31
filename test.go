package scan

import (
	"strings"
	"testing"
)

type Test struct {
	src  string
	toks []Token
	errs [][]string
}

func NewTest(src string, val string, line int, col int, type_ string) Test {
	return Test{
		src: src,
		toks: []Token{{
			Value: val,
			Type:  type_,
			Pos: Pos{
				Line:   line,
				Column: col,
			}},
		},
		errs: [][]string{nil},
	}
}

func (t Test) WithError(message string) Test {
	last := len(t.toks) - 1
	t.errs[last] = append(t.errs[last], message)
	return t
}

func (t Test) And(val string, line int, col int, type_ string) Test {
	t.toks = append(t.toks, Token{
		Value: val,
		Type:  type_,
		Pos: Pos{
			Line:   line,
			Column: col,
		},
	})
	t.errs = append(t.errs, nil)
	return t
}

func RunTests(t *testing.T, rules RuleSet, tests []Test) {
	for _, test := range tests {
		t.Run(test.src, func(t *testing.T) {
			scan := NewFromString("", test.src)
			r := NewRunner(scan, rules)
			toks := r.All()
			t.Log("\n" + FormatTokenTable(toks))
			for i, tok := range toks {
				if i >= len(test.toks) {
					break
				}
				testTok := test.toks[i]
				if tok.Value != testTok.Value || tok.Type != testTok.Type || tok.Pos != testTok.Pos {
					t.Fatalf("\n have: %v \n want: %v", tok, testTok)
				}
				if len(test.errs) > 0 {
					have := tok.Errs.Error()
					want := strings.Join(test.errs[i], "\n")
					if have != want {
						t.Fatalf("\n have errors: %v\n want errors: %v", have, want)
					}
				}
			}
			if len(toks) < len(test.toks) {
				idx := len(toks)
				t.Errorf("not enough tokens, missing: %v", test.toks[idx:])
			}
			if len(toks) > len(test.toks) {
				idx := len(test.toks)
				t.Errorf("too many tokens, extra: %v", toks[idx:])
			}
		})
	}
}
