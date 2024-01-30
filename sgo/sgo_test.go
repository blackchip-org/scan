package sgo

import (
	"testing"

	"github.com/blackchip-org/scan"
)

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
			WithError("foo"),
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
			WithError(`1:4: error: not valid: "D800"`),
		scan.NewTest(`"\U00110000"`, "00110000", 1, 1, scan.IllegalType).
			WithError(`1:4: error: not valid: "00110000"`),
	}
	scan.RunTests(t, rules, tests)
}
