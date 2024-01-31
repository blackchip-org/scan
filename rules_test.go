package scan

import "testing"

func TestBin(t *testing.T) {
	rules := Rules(Bin)
	tests := []Test{
		NewTest("012", "01", 1, 1, BinType).
			And("2", 1, 3, IllegalType).
			WithError(`1:3: error: unexpected "2"`),
	}
	RunTests(t, rules, tests)
}

func TestBin0b(t *testing.T) {
	rules := Rules(Bin0b, Int)
	tests := []Test{
		NewTest("0b0101", "0101", 1, 1, BinType),
		NewTest("0B0101", "0101", 1, 1, BinType),
		NewTest("0123", "0123", 1, 1, IntType),
	}
	RunTests(t, rules, tests)
}

func TestHex(t *testing.T) {
	rules := Rules(Hex)
	tests := []Test{
		NewTest("0123456789abcdefg", "0123456789abcdef", 1, 1, HexType).
			And("g", 1, 17, IllegalType).
			WithError(`1:17: error: unexpected "g"`),
	}
	RunTests(t, rules, tests)
}

func TestHex0x(t *testing.T) {
	rules := Rules(Hex0x, Int)
	tests := []Test{
		NewTest("0x12bc", "12bc", 1, 1, HexType),
		NewTest("0X12bc", "12bc", 1, 1, HexType),
		NewTest("12bc", "12", 1, 1, IntType).
			And("b", 1, 3, IllegalType).
			WithError(`1:3: error: unexpected "b"`).
			And("c", 1, 4, IllegalType).
			WithError(`1:4: error: unexpected "c"`),
	}
	RunTests(t, rules, tests)
}

func TestIdent(t *testing.T) {
	rules := Rules(Ident.WithKeywords("true", "false"))
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
	rules := Rules(Int)
	tests := []Test{
		NewTest("12345", "12345", 1, 1, IntType),
		NewTest("-1234", "-", 1, 1, IllegalType).
			WithError(`1:1: error: unexpected "-"`).
			And("1234", 1, 2, IntType),
	}
	RunTests(t, rules, tests)
}

func TestLiterals(t *testing.T) {
	rules := Rules(Literals("=", "===", "+", "+=", "/"))
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
	rules := Rules(Oct)
	tests := []Test{
		NewTest("012345678", "01234567", 1, 1, OctType).
			And("8", 1, 9, IllegalType).
			WithError(`1:9: error: unexpected "8"`),
	}
	RunTests(t, rules, tests)
}

func TestOct0o(t *testing.T) {
	rules := Rules(Oct0o, Int)
	tests := []Test{
		NewTest("0o755", "755", 1, 1, OctType),
		NewTest("0O755", "755", 1, 1, OctType),
		NewTest("7558", "7558", 1, 1, IntType),
	}
	RunTests(t, rules, tests)
}

func TestReal(t *testing.T) {
	rules := Rules(Real)
	tests := []Test{
		NewTest("12.345", "12.345", 1, 1, RealType),
		NewTest("-12.34", "-", 1, 1, IllegalType).
			WithError(`1:1: error: unexpected "-"`).
			And("12.34", 1, 2, RealType),
	}
	RunTests(t, rules, tests)
}

func TestSignedInt(t *testing.T) {
	rules := Rules(SignedInt)
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
	rules := Rules(SignedReal)
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
	rules := Rules(SignedRealExp)
	tests := []Test{
		NewTest("123", "123", 1, 1, IntType),
		NewTest("123e10", "123e10", 1, 1, RealType),
		NewTest("1.23", "1.23", 1, 1, RealType),
		NewTest("-123.45E-10", "-123.45E-10", 1, 1, RealType),
	}
	RunTests(t, rules, tests)
}

func TestSignedRealWithDigitSep(t *testing.T) {
	rules := Rules(SignedRealExp.WithDigitSep(Rune('_')))
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
	rules := Rules(
		StrDoubleQuote,
		StrSingleQuote,
	)
	tests := []Test{
		NewTest(`"a"`, `a`, 1, 1, StrType),
		NewTest("'a'", "a", 1, 1, StrType),
		NewTest(`'a\'b'`, `a'b`, 1, 1, StrType),
	}
	RunTests(t, rules, tests)
}
