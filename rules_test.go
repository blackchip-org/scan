package scan

import "testing"

func TestBin(t *testing.T) {
	rules := NewRuleSet(BinRule)
	tests := []Test{
		NewTest("012", "01", 1, 1, BinType).
			And("2", 1, 3, IllegalType).
			WithError(`1:3: error: unexpected "2"`),
	}
	RunTests(t, rules, tests)
}

func TestBin0b(t *testing.T) {
	rules := NewRuleSet(Bin0bRule, IntRule)
	tests := []Test{
		NewTest("0b0101", "0b0101", 1, 1, BinType),
		NewTest("0B0101", "0B0101", 1, 1, BinType),
		NewTest("0123", "0123", 1, 1, IntType),
	}
	RunTests(t, rules, tests)
}

func TestClassRule(t *testing.T) {
	rules := NewRuleSet(
		NewClassRule(Rune('-')),
		NewClassRule(Rune('+')).WithType("plus"),
	)
	tests := []Test{
		NewTest("-", "-", 1, 1, "-"),
		NewTest("+", "+", 1, 1, "plus"),
	}
	RunTests(t, rules, tests)
}

func TestEscapeRules(t *testing.T) {
	rule := NewRuleSet(StrDoubleQuoteRule.WithEscapeRules(
		NewCharEncRule(
			AlertEnc,
			BackspaceEnc,
			FormFeedEnc,
			LineFeedEnc,
			CarriageReturnEnc,
			HorizontalTabEnc,
			VerticalTabEnc,
		),
		Hex2EncRule.AsByte(true),
		Hex4EncRule,
		Hex8EncRule,
		OctEnc,
	))
	tests := []Test{
		NewTest(`"日本語"`, `日本語`, 1, 1, StrType),
		NewTest(`"\u65e5\u672c\u8a9e"`, `日本語`, 1, 1, StrType),
		NewTest(`"\U000065e5\U0000672c\U00008a9e"`, `日本語`, 1, 1, StrType),
		NewTest(`"\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e"`, `日本語`, 1, 1, StrType),
	}
	RunTests(t, rule, tests)

}

func TestHex(t *testing.T) {
	rules := NewRuleSet(HexRule)
	tests := []Test{
		NewTest("0123456789abcdefg", "0123456789abcdef", 1, 1, HexType).
			And("g", 1, 17, IllegalType).
			WithError(`1:17: error: unexpected "g"`),
	}
	RunTests(t, rules, tests)
}

func TestHex0x(t *testing.T) {
	rules := NewRuleSet(Hex0xRule, IntRule)
	tests := []Test{
		NewTest("0x12bc", "0x12bc", 1, 1, HexType),
		NewTest("0X12bc", "0X12bc", 1, 1, HexType),
		NewTest("12bc", "12", 1, 1, IntType).
			And("b", 1, 3, IllegalType).
			WithError(`1:3: error: unexpected "b"`).
			And("c", 1, 4, IllegalType).
			WithError(`1:4: error: unexpected "c"`),
	}
	RunTests(t, rules, tests)
}

func TestIdent(t *testing.T) {
	rules := NewRuleSet(StandardIdentRule.WithKeywords("true", "false"))
	tests := []Test{
		NewTest("abc_123", "abc_123", 1, 1, IdentType),
		NewTest("_abc_123", "_abc_123", 1, 1, IdentType),
		NewTest("0abc_123", "0", 1, 1, IllegalType).
			WithError(`1:1: error: unexpected "0"`).
			And("abc_123", 1, 2, IdentType),
		NewTest("true", "true", 1, 1, "true"),
		NewTest("false", "false", 1, 1, "false"),
	}
	RunTests(t, rules, tests)
}

func TestInt(t *testing.T) {
	rules := NewRuleSet(IntRule)
	tests := []Test{
		NewTest("12345", "12345", 1, 1, IntType),
		NewTest("-1234", "-", 1, 1, IllegalType).
			WithError(`1:1: error: unexpected "-"`).
			And("1234", 1, 2, IntType),
	}
	RunTests(t, rules, tests)
}

func TestLiterals(t *testing.T) {
	rules := NewRuleSet(Literal("=", "===", "+", "+=", "/"))
	tests := []Test{
		NewTest("=", "=", 1, 1, "="),
		NewTest("===", "===", 1, 1, "==="),
		NewTest("==", "=", 1, 1, "=").
			And("=", 1, 2, "="),
		NewTest("+=/", "+=", 1, 1, "+=").
			And("/", 1, 3, "/"),
	}
	RunTests(t, rules, tests)
}

func TestOct(t *testing.T) {
	rules := NewRuleSet(OctRule)
	tests := []Test{
		NewTest("012345678", "01234567", 1, 1, OctType).
			And("8", 1, 9, IllegalType).
			WithError(`1:9: error: unexpected "8"`),
	}
	RunTests(t, rules, tests)
}

func TestOct0o(t *testing.T) {
	rules := NewRuleSet(Oct0oRule, IntRule)
	tests := []Test{
		NewTest("0o755", "0o755", 1, 1, OctType),
		NewTest("0O755", "0O755", 1, 1, OctType),
		NewTest("7558", "7558", 1, 1, IntType),
	}
	RunTests(t, rules, tests)
}

func TestReal(t *testing.T) {
	rules := NewRuleSet(RealRule)
	tests := []Test{
		NewTest("12.345", "12.345", 1, 1, RealType),
		NewTest("-12.34", "-", 1, 1, IllegalType).
			WithError(`1:1: error: unexpected "-"`).
			And("12.34", 1, 2, RealType),
	}
	RunTests(t, rules, tests)
}

func TestSignedInt(t *testing.T) {
	rules := NewRuleSet(SignedIntRule)
	tests := []Test{
		NewTest("1234567890a", "1234567890", 1, 1, IntType).
			And("a", 1, 11, IllegalType).
			WithError(`1:11: error: unexpected "a"`),
		NewTest("-1234", "-1234", 1, 1, IntType),
		NewTest("+1234", "+1234", 1, 1, IntType),
		NewTest("0x1234", "0", 1, 1, IntType).
			And("x", 1, 2, IllegalType).
			WithError(`1:2: error: unexpected "x"`).
			And("1234", 1, 3, IntType),
	}
	RunTests(t, rules, tests)
}

func TestSignedReal(t *testing.T) {
	rules := NewRuleSet(SignedRealRule)
	tests := []Test{
		NewTest("1234", "1234", 1, 1, IntType),
		NewTest("12.345", "12.345", 1, 1, RealType),
		NewTest("-12.345", "-12.345", 1, 1, RealType),
		NewTest(".1234", ".1234", 1, 1, RealType),

		NewTest("-123.45e10", "-123.45", 1, 1, RealType).
			And("e", 1, 8, IllegalType).
			WithError(`1:8: error: unexpected "e"`).
			And("10", 1, 9, IntType),
		NewTest("-12.34.5", "-12.34", 1, 1, RealType).
			And(".5", 1, 7, RealType),
		NewTest("-12.34-5", "-12.34", 1, 1, RealType).
			And("-5", 1, 7, IntType),
		NewTest("1_234", "1", 1, 1, IntType).
			And("_", 1, 2, IllegalType).
			WithError(`1:2: error: unexpected "_"`).
			And("234", 1, 3, IntType),
	}
	RunTests(t, rules, tests)
}

func TestSignedRealExp(t *testing.T) {
	rules := NewRuleSet(SignedRealExpRule)
	tests := []Test{
		NewTest("123", "123", 1, 1, IntType),
		NewTest("123e10", "123e10", 1, 1, RealType),
		NewTest("1.23", "1.23", 1, 1, RealType),
		NewTest("-123.45E-10", "-123.45E-10", 1, 1, RealType),
	}
	RunTests(t, rules, tests)
}

func TestSignedRealWithDigitSep(t *testing.T) {
	rules := NewRuleSet(SignedRealExpRule.WithDigitSep(Rune('_')))
	tests := []Test{
		NewTest("1234567", "1234567", 1, 1, IntType),
		NewTest("1_234_567", "1234567", 1, 1, IntType),
		NewTest("0_600", "0600", 1, 1, IntType),
		NewTest("0.15e+0_2", "0.15e+02", 1, 1, RealType),
		NewTest("1_234__567", "1234", 1, 1, IntType).
			And("_", 1, 6, IllegalType).
			WithError(`1:6: error: unexpected "_"`).
			And("_", 1, 7, IllegalType).
			WithError(`1:7: error: unexpected "_"`).
			And("567", 1, 8, IntType),
	}
	RunTests(t, rules, tests)
}

func TestStr(t *testing.T) {
	rules := NewRuleSet(
		StrDoubleQuoteRule,
		StrSingleQuoteRule,
	)
	tests := []Test{
		NewTest(`"a"`, `a`, 1, 1, StrType),
		NewTest("'a'", "a", 1, 1, StrType),
		NewTest(`'a\'b'`, `a'b`, 1, 1, StrType),
	}
	RunTests(t, rules, tests)
}
