package sgo

import (
	"testing"

	"github.com/blackchip-org/scan"
)

func TestInt(t *testing.T) {
	rules := scan.Rules(Ident, Int)
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
	scan.RunTests(t, rules, tests)
}

func TestRawString(t *testing.T) {
	rules := scan.Rules(RawString)
	tests := []scan.Test{
		scan.NewTest("`abc`", "abc", 1, 1, StringType),
		scan.NewTest("`\\n\n\\n`", "\\n\n\\n", 1, 1, StringType),
	}
	scan.RunTests(t, rules, tests)
}

func TestRune(t *testing.T) {
	rules := scan.Rules(Rune)
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
	scan.RunTests(t, rules, tests)
}

func TestString(t *testing.T) {
	rules := scan.Rules(String)
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
	scan.RunTests(t, rules, tests)
}
