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
		scan.NewTest("0o600", "600", 1, 1, OctType),
		scan.NewTest("0O600", "600", 1, 1, OctType),
		scan.NewTest("0xBadFace", "BadFace", 1, 1, HexType),
		scan.NewTest("0xBad_Face", "BadFace", 1, 1, HexType),
		scan.NewTest("0x_67_7a_2f_cc_40_c6", "677a2fcc40c6", 1, 1, HexType),
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
