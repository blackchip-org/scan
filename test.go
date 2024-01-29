package scan

import "testing"

type TestCase struct {
	src  string
	toks []Token
}

func NewTestCase(src string, val string, line int, col int, type_ string) TestCase {
	return TestCase{src: src, toks: []Token{
		{Value: val, Type: type_, Pos: Pos{Line: line, Column: col}},
	}}
}

func (t TestCase) And(val string, line int, col int, type_ string) TestCase {
	t.toks = append(t.toks, Token{Value: val, Type: type_, Pos: Pos{Line: line, Column: col}})
	return t
}

func RunTests(t *testing.T, rules RuleSet, tests []TestCase) {
	for _, test := range tests {
		t.Run(test.src, func(t *testing.T) {
			scan := NewFromString("", test.src)
			r := NewRunner(scan, rules)
			toks := r.All()
			for i, tok := range toks {
				if i >= len(test.toks) {
					break
				}
				testTok := test.toks[i]
				if tok.Value != testTok.Value || tok.Type != testTok.Type {
					t.Fatalf("\n have: %v \n want: %v", tok, testTok)
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
