package sgo

import (
	"testing"

	"github.com/blackchip-org/scan"
)

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
			WithError("1:4: error: character not valid: 'D800'"),
	}
	scan.RunTests(t, rules, tests)
}
