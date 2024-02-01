package scanjson

import (
	_ "embed"
	"testing"

	"github.com/blackchip-org/scan"
)

func TestKeywords(t *testing.T) {
	ctx := NewContext()
	tests := []scan.Test{
		scan.NewTest("true", "true", 1, 1, "true"),
		scan.NewTest("false", "false", 1, 1, "false"),
		scan.NewTest("null", "null", 1, 1, "null"),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}

func TestPunct(t *testing.T) {
	ctx := NewContext()
	tests := []scan.Test{
		scan.NewTest("{", "{", 1, 1, "{"),
		scan.NewTest("}", "}", 1, 1, "}"),
		scan.NewTest("[", "[", 1, 1, "["),
		scan.NewTest("]", "]", 1, 1, "]"),
		scan.NewTest(":", ":", 1, 1, ":"),
		scan.NewTest(",", ",", 1, 1, ","),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}

func TestStrings(t *testing.T) {
	ctx := NewContext()
	tests := []scan.Test{
		scan.NewTest(`"foo"`, "foo", 1, 1, scan.StrType),
		scan.NewTest(`"\n"`, "\n", 1, 1, scan.StrType),
		scan.NewTest(`"\\"`, "\\", 1, 1, scan.StrType),
		scan.NewTest(`"\u12e4"`, "ዤ", 1, 1, scan.StrType),
		scan.NewTest(`"foo`, "foo", 1, 1, scan.IllegalType).
			WithError("1:5: error: not terminated"),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}

func TestNumbers(t *testing.T) {
	ctx := NewContext()
	tests := []scan.Test{
		scan.NewTest("0", "0", 1, 1, scan.IntType),
		scan.NewTest("42", "42", 1, 1, scan.IntType),
		scan.NewTest("72.40", "72.40", 1, 1, scan.RealType),
		scan.NewTest("-72.40", "-72.40", 1, 1, scan.RealType),
		scan.NewTest("-72.40e+10", "-72.40e+10", 1, 1, scan.RealType),
		scan.NewTest("042", "0", 1, 1, scan.IntType).
			And("42", 1, 2, scan.IntType),
		scan.NewTest("42.", "42", 1, 1, scan.IntType).
			And(".", 1, 3, scan.IllegalType).
			WithError(`1:3: error: unexpected "."`),
		scan.NewTest("+72.40", "+", 1, 1, scan.IllegalType).
			WithError(`1:1: error: unexpected "+"`).
			And("72.40", 1, 2, scan.RealType),
		scan.NewTest("-72.40e", "-72.40", 1, 1, scan.RealType).
			And("e", 1, 7, scan.IllegalType).
			WithError(`1:7: error: unexpected "e"`),
		scan.NewTest("-e", "-", 1, 1, scan.IllegalType).
			WithError(`1:1: error: unexpected "-"`).
			And("e", 1, 2, scan.IllegalType).
			WithError(`1:2: error: unexpected "e"`),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}

//go:embed example.json
var example string

//go:embed example.tokens.json
var exampleTokens string

// func TestFile(t *testing.T) {
// 	s := scan.NewScannerFromString("example.json", example)
// 	ctx := NewContext()
// 	toks := scan.NewRunner(s, ctx.RuleSet).All()
// 	data, err := json.MarshalIndent(toks, "", "    ")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	have := string(data)
// 	if have != exampleTokens {
// 		t.Errorf("have:\n%v\nwant:\n%v", have, exampleTokens)
// 	}
// }
