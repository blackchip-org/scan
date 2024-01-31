package sgo

import (
	"testing"

	"github.com/blackchip-org/scan"
)

func TestComments(t *testing.T) {
	ctx := NewContext()
	tests := []scan.Test{
		scan.NewTest("abc // cde \n fgh", "abc", 1, 1, IdentType).
			And("fgh", 2, 2, IdentType),
		scan.NewTest("abc /* cde \n fgh */ ijk", "abc", 1, 1, IdentType).
			And("ijk", 2, 9, IdentType),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}

func TestCommentsKeep(t *testing.T) {
	ctx := NewContext()
	ctx.KeepComments = true
	tests := []scan.Test{
		scan.NewTest("abc // cde \n fgh", "abc", 1, 1, IdentType).
			And(" cde ", 1, 5, CommentType).
			And("fgh", 2, 2, IdentType),
		scan.NewTest("abc /* cde \n fgh */ ijk", "abc", 1, 1, IdentType).
			And(" cde \n fgh ", 1, 5, CommentType).
			And("ijk", 2, 9, IdentType),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}

func TestFloat(t *testing.T) {
	ctx := NewContext()
	tests := []scan.Test{
		scan.NewTest("0.", "0.", 1, 1, FloatType),
		scan.NewTest("72.40", "72.40", 1, 1, FloatType),
		scan.NewTest("072.40", "072.40", 1, 1, FloatType),
		scan.NewTest("2.71828", "2.71828", 1, 1, FloatType),
		scan.NewTest("1.e+0", "1.e+0", 1, 1, FloatType),
		scan.NewTest("6.67428e-11", "6.67428e-11", 1, 1, FloatType),
		scan.NewTest("1E6", "1E6", 1, 1, FloatType),
		scan.NewTest(".25", ".25", 1, 1, FloatType),
		scan.NewTest(".12345E+5", ".12345E+5", 1, 1, FloatType),
		scan.NewTest("1_5.", "15.", 1, 1, FloatType),
		scan.NewTest("0.15e+0_2", "0.15e+02", 1, 1, FloatType),
		scan.NewTest("0x1p-2", "0x1p-2", 1, 1, FloatType),
		scan.NewTest("0x2.p10", "0x2.p10", 1, 1, FloatType),
		scan.NewTest("0x1.Fp+0", "0x1.Fp+0", 1, 1, FloatType),
		scan.NewTest("0X.8p-0", "0X.8p-0", 1, 1, FloatType),
		scan.NewTest("0X_1FFFP-16", "0X1FFFP-16", 1, 1, FloatType),
		scan.NewTest("0x15e-2", "0x15e", 1, 1, IntType).
			And("-", 1, 6, "-").
			And("2", 1, 7, IntType),
		scan.NewTest("0x.p1", "0", 1, 1, IntType).
			And("x", 1, 2, IdentType).
			And(".", 1, 3, ".").
			And("p1", 1, 4, IdentType),
		scan.NewTest("1p-2", "1", 1, 1, IntType).
			And("p", 1, 2, IdentType).
			And("-", 1, 3, "-").
			And("2", 1, 4, IntType),
		scan.NewTest("0x1.5e-2", "0x1.5e", 1, 1, FloatType).
			And("-", 1, 7, "-").
			And("2", 1, 8, IntType),
		scan.NewTest("1_.5", "1", 1, 1, IntType).
			And("_", 1, 2, IdentType).
			And(".5", 1, 3, FloatType),
		scan.NewTest("1._5", "1.", 1, 1, FloatType).
			And("_5", 1, 3, IdentType),
		scan.NewTest("1.5_e1", "1.5", 1, 1, FloatType).
			And("_e1", 1, 4, IdentType),
		scan.NewTest("1.5e_1", "1.5", 1, 1, FloatType).
			And("e_1", 1, 4, IdentType),
		scan.NewTest("1.5e1_", "1.5e1", 1, 1, FloatType).
			And("_", 1, 6, IdentType),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}

func TestIdent(t *testing.T) {
	ctx := NewContext()
	tests := []scan.Test{
		scan.NewTest("a", "a", 1, 1, IdentType),
		scan.NewTest("_x9", "_x9", 1, 1, IdentType),
		scan.NewTest("ThisVariableIsExported", "ThisVariableIsExported", 1, 1, IdentType),
		scan.NewTest("αβ", "αβ", 1, 1, IdentType),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}

func TestInt(t *testing.T) {
	ctx := NewContext()
	tests := []scan.Test{
		scan.NewTest("42", "42", 1, 1, IntType),
		scan.NewTest("4_2", "42", 1, 1, IntType),
		scan.NewTest("0600", "0600", 1, 1, IntType),
		scan.NewTest("0_600", "0600", 1, 1, IntType),
		scan.NewTest("0o600", "0o600", 1, 1, IntType),
		scan.NewTest("0O600", "0O600", 1, 1, IntType),
		scan.NewTest("0xBadFace", "0xBadFace", 1, 1, IntType),
		scan.NewTest("0xBad_Face", "0xBadFace", 1, 1, IntType),
		scan.NewTest("0x_67_7a_2f_cc_40_c6", "0x677a2fcc40c6", 1, 1, IntType),
		scan.NewTest("170141183460469231731687303715884105727", "170141183460469231731687303715884105727", 1, 1, IntType),
		scan.NewTest("170_141183_460469_231731_687303_715884_105727", "170141183460469231731687303715884105727", 1, 1, IntType),
		scan.NewTest("_42", "_42", 1, 1, IdentType),
		scan.NewTest("4__2", "4", 1, 1, IntType).
			And("__2", 1, 2, IdentType),
		scan.NewTest("0_xBadFace", "0", 1, 1, IntType).
			And("_xBadFace", 1, 2, IdentType),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}

func TestImag(t *testing.T) {
	ctx := NewContext()
	tests := []scan.Test{
		scan.NewTest("0i", "0i", 1, 1, ImagType),
		scan.NewTest("0123i", "0123i", 1, 1, ImagType),
		scan.NewTest("0o123i", "0o123i", 1, 1, ImagType),
		scan.NewTest("0xabci", "0xabci", 1, 1, ImagType),
		scan.NewTest("0.i", "0.i", 1, 1, ImagType),
		scan.NewTest("2.71828i", "2.71828i", 1, 1, ImagType),
		scan.NewTest("1.e+0i", "1.e+0i", 1, 1, ImagType),
		scan.NewTest("6.67428e-11i", "6.67428e-11i", 1, 1, ImagType),
		scan.NewTest("1E6i", "1E6i", 1, 1, ImagType),
		scan.NewTest(".25i", ".25i", 1, 1, ImagType),
		scan.NewTest(".12345E+5i", ".12345E+5i", 1, 1, ImagType),
		scan.NewTest("0x1p-2i", "0x1p-2i", 1, 1, ImagType),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}

func TestKeywords(t *testing.T) {
	ctx := NewContext()
	tests := []scan.Test{
		scan.NewTest("break", "break", 1, 1, "break"),
		scan.NewTest("case", "case", 1, 1, "case"),
		scan.NewTest("chan", "chan", 1, 1, "chan"),
		scan.NewTest("const", "const", 1, 1, "const"),
		scan.NewTest("continue", "continue", 1, 1, "continue"),
		scan.NewTest("default", "default", 1, 1, "default"),
		scan.NewTest("defer", "defer", 1, 1, "defer"),
		scan.NewTest("else", "else", 1, 1, "else"),
		scan.NewTest("fallthrough", "fallthrough", 1, 1, "fallthrough"),
		scan.NewTest("for", "for", 1, 1, "for"),
		scan.NewTest("func", "func", 1, 1, "func"),
		scan.NewTest("go", "go", 1, 1, "go"),
		scan.NewTest("goto", "goto", 1, 1, "goto"),
		scan.NewTest("if", "if", 1, 1, "if"),
		scan.NewTest("import", "import", 1, 1, "import"),
		scan.NewTest("interface", "interface", 1, 1, "interface"),
		scan.NewTest("map", "map", 1, 1, "map"),
		scan.NewTest("package", "package", 1, 1, "package"),
		scan.NewTest("range", "range", 1, 1, "range"),
		scan.NewTest("return", "return", 1, 1, "return"),
		scan.NewTest("select", "select", 1, 1, "select"),
		scan.NewTest("struct", "struct", 1, 1, "struct"),
		scan.NewTest("switch", "switch", 1, 1, "switch"),
		scan.NewTest("type", "type", 1, 1, "type"),
		scan.NewTest("var", "var", 1, 1, "var"),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}

func TestRawString(t *testing.T) {
	ctx := NewContext()
	tests := []scan.Test{
		scan.NewTest("`abc`", "abc", 1, 1, StringType),
		scan.NewTest("`\\n\n\\n`", "\\n\n\\n", 1, 1, StringType),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}

func TestRune(t *testing.T) {
	ctx := NewContext()
	tests := []scan.Test{
		scan.NewTest(`'a'`, "a", 1, 1, RuneType),
		scan.NewTest(`'ä'`, "ä", 1, 1, RuneType),
		scan.NewTest(`'本'`, "本", 1, 1, RuneType),
		scan.NewTest(`'\t'`, "\t", 1, 1, RuneType),
		scan.NewTest(`'\000'`, "\x00", 1, 1, RuneType),
		scan.NewTest(`'\007'`, "\a", 1, 1, RuneType),
		scan.NewTest(`'\377'`, "ÿ", 1, 1, RuneType),
		scan.NewTest(`'\x07'`, "\a", 1, 1, RuneType),
		scan.NewTest(`'\xff'`, "ÿ", 1, 1, RuneType),
		scan.NewTest(`'\u12e4'`, "ዤ", 1, 1, RuneType),
		scan.NewTest(`'\U00101234'`, "\U00101234", 1, 1, RuneType),
		scan.NewTest(`'\''`, `'`, 1, 1, RuneType),
		scan.NewTest(`'aa'`, "a", 1, 1, scan.IllegalType).
			WithError("1:3: error: too many characters (2)"),
		scan.NewTest(`'\k'`, "k", 1, 1, scan.IllegalType).
			WithError(`1:3: error: invalid escape sequence: '\k'`),
		scan.NewTest(`'\xa'`, "a", 1, 1, scan.IllegalType).
			WithError(`1:4: error: invalid encoding: "a"`),
		scan.NewTest(`'\0'`, "0", 1, 1, scan.IllegalType).
			WithError(`1:3: error: invalid encoding: "0"`),
		scan.NewTest(`'\400'`, "400", 1, 1, scan.IllegalType).
			WithError(`1:3: error: invalid encoding: "400"`),
		scan.NewTest(`'\uDFFF'`, "DFFF", 1, 1, scan.IllegalType).
			WithError(`1:4: error: invalid encoding: "DFFF"`),
		scan.NewTest(`'\U00110000'`, "00110000", 1, 1, scan.IllegalType).
			WithError(`1:4: error: invalid encoding: "00110000"`),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}

func TestSemiColons(t *testing.T) {
	ctx := NewContext()
	tests := []scan.Test{
		scan.NewTest("a++\n42", "a", 1, 1, IdentType).
			And("++", 1, 2, "++").
			And(";", 1, 4, ";").
			And("42", 2, 1, IntType),
		scan.NewTest("case\n42", "case", 1, 1, "case").
			And("42", 2, 1, IntType),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}

func TestString(t *testing.T) {
	ctx := NewContext()
	tests := []scan.Test{
		scan.NewTest(`"abc"`, `abc`, 1, 1, StringType),
		scan.NewTest(`"\a\b\f\n\r\t\v\\\""`, "\a\b\f\n\r\t\v\\\"", 1, 1, StringType),
		scan.NewTest(`"a\"b"`, `a"b`, 1, 1, StringType),
		scan.NewTest(`"日本語"`, "日本語", 1, 1, StringType),
		scan.NewTest(`"\u65e5本\U00008a9e"`, "日本語", 1, 1, StringType),
		scan.NewTest(`"\xff\u00FF"`, "ÿÿ", 1, 1, StringType),
		scan.NewTest(`"\uD800"`, "D800", 1, 1, scan.IllegalType).
			WithError(`1:4: error: invalid encoding: "D800"`),
		scan.NewTest(`"\U00110000"`, "00110000", 1, 1, scan.IllegalType).
			WithError(`1:4: error: invalid encoding: "00110000"`),
	}
	scan.RunTests(t, ctx.RuleSet, tests)
}
