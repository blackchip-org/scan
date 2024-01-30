package scan

import "testing"

func TestBin(t *testing.T) {
	rules := Rules(Bin)
	tests := []TestCase{
		NewTestCase("012", "01", 1, 1, BinType).
			And("2", 1, 1, IllegalType),
	}
	RunTests(t, rules, tests)
}

func TestBin0b(t *testing.T) {
	rules := Rules(Bin0b, Int)
	tests := []TestCase{
		NewTestCase("0b0101", "0101", 1, 1, BinType),
		NewTestCase("0B0101", "0101", 1, 1, BinType),
		NewTestCase("0123", "0123", 1, 1, IntType),
	}
	RunTests(t, rules, tests)
}

func TestHex(t *testing.T) {
	rules := Rules(Hex)
	tests := []TestCase{
		NewTestCase("0123456789abcdefg", "0123456789abcdef", 1, 1, HexType).
			And("g", 1, 17, IllegalType),
	}
	RunTests(t, rules, tests)
}

func TestHex0x(t *testing.T) {
	rules := Rules(Hex0x, Int)
	tests := []TestCase{
		NewTestCase("0x12bc", "12bc", 1, 1, HexType),
		NewTestCase("0X12bc", "12bc", 1, 1, HexType),
		NewTestCase("12bc", "12", 1, 1, IntType).
			And("b", 1, 3, IllegalType).
			And("c", 1, 4, IllegalType),
	}
	RunTests(t, rules, tests)
}

func TestIdent(t *testing.T) {
	rules := Rules(Ident.WithKeywords("true", "false"))
	tests := []TestCase{
		NewTestCase("abc_123", "abc_123", 1, 1, IdentType),
		NewTestCase("_abc_123", "_abc_123", 1, 1, IdentType),
		NewTestCase("0abc_123", "0", 1, 1, IllegalType).
			And("abc_123", 1, 2, IdentType),
		NewTestCase("true", "true", 1, 1, "true"),
		NewTestCase("false", "false", 1, 1, "false"),
	}
	RunTests(t, rules, tests)
}

func TestInt(t *testing.T) {
	rules := Rules(Int)
	tests := []TestCase{
		NewTestCase("12345", "12345", 1, 1, IntType),
		NewTestCase("-1234", "-", 1, 1, IllegalType).
			And("1234", 1, 2, IntType),
	}
	RunTests(t, rules, tests)
}

func TestLiterals(t *testing.T) {
	rules := Rules(Literals("=", "===", "+", "+=", "/"))
	tests := []TestCase{
		NewTestCase("=", "=", 1, 1, "="),
		NewTestCase("===", "===", 1, 1, "==="),
		NewTestCase("==", "=", 1, 1, "=").
			And("=", 2, 2, "="),
		NewTestCase("+=/", "+=", 1, 1, "+=").
			And("/", 1, 2, "/"),
	}
	RunTests(t, rules, tests)
}

func TestOct(t *testing.T) {
	rules := Rules(Oct)
	tests := []TestCase{
		NewTestCase("012345678", "01234567", 1, 1, OctType).
			And("8", 1, 9, IllegalType),
	}
	RunTests(t, rules, tests)
}

func TestOct0o(t *testing.T) {
	rules := Rules(Oct0o, Int)
	tests := []TestCase{
		NewTestCase("0o755", "755", 1, 1, OctType),
		NewTestCase("0O755", "755", 1, 1, OctType),
		NewTestCase("7558", "7558", 1, 1, IntType),
	}
	RunTests(t, rules, tests)
}

func TestReal(t *testing.T) {
	rules := Rules(Real)
	tests := []TestCase{
		NewTestCase("12.345", "12.345", 1, 1, RealType),
		NewTestCase("-12.34", "-", 1, 1, IllegalType).
			And("12.34", 1, 2, RealType),
	}
	RunTests(t, rules, tests)
}

func TestSignedInt(t *testing.T) {
	rules := Rules(SignedInt)
	tests := []TestCase{
		NewTestCase("1234567890a", "1234567890", 1, 1, IntType).
			And("a", 1, 11, IllegalType),
		NewTestCase("-1234", "-1234", 1, 1, IntType),
		NewTestCase("+1234", "+1234", 1, 1, IntType),
		NewTestCase("0x1234", "0", 1, 1, IntType).
			And("x", 1, 2, IllegalType).
			And("1234", 1, 3, IntType),
	}
	RunTests(t, rules, tests)
}

func TestSignedReal(t *testing.T) {
	rules := Rules(SignedReal)
	tests := []TestCase{
		NewTestCase("1234", "1234", 1, 1, IntType),
		NewTestCase("12.345", "12.345", 1, 1, RealType),
		NewTestCase("-12.345", "-12.345", 1, 1, RealType),
		NewTestCase(".1234", ".1234", 1, 1, RealType),

		NewTestCase("-123.45e10", "-123.45", 1, 1, RealType).
			And("e", 1, 8, IllegalType).
			And("10", 1, 9, IntType),
		NewTestCase("-12.34.5", "-12.34", 1, 1, RealType).
			And(".5", 1, 6, RealType),
		NewTestCase("-12.34-5", "-12.34", 1, 1, RealType).
			And("-5", 1, 6, IntType),
		NewTestCase("1_234", "1", 1, 1, IntType).
			And("_", 1, 2, IllegalType).
			And("234", 1, 3, IntType),
	}
	RunTests(t, rules, tests)
}

func TestSignedRealExp(t *testing.T) {
	rules := Rules(SignedRealExp)
	tests := []TestCase{
		NewTestCase("123", "123", 1, 1, IntType),
		NewTestCase("123e10", "123e10", 1, 1, RealType),
		NewTestCase("1.23", "1.23", 1, 1, RealType),
		NewTestCase("-123.45E-10", "-123.45E-10", 1, 1, RealType),
	}
	RunTests(t, rules, tests)
}

func TestSignedRealWithDigitSep(t *testing.T) {
	rules := Rules(SignedRealExp.WithDigitSep(Rune('_')))
	tests := []TestCase{
		NewTestCase("1234567", "1234567", 1, 1, IntType),
		NewTestCase("1_234_567", "1234567", 1, 1, IntType),
		NewTestCase("0_600", "0600", 1, 1, IntType),
		NewTestCase("0.15e+0_2", "0.15e+02", 1, 1, RealType),
		NewTestCase("1_234__567", "1234", 1, 1, IntType).
			And("_", 1, 6, IllegalType).
			And("_", 1, 7, IllegalType).
			And("567", 1, 8, IntType),
	}
	RunTests(t, rules, tests)
}

func TestStr(t *testing.T) {
	rules := Rules(
		StrDoubleQuote,
		StrSingleQuote,
	)
	tests := []TestCase{
		NewTestCase(`"a"`, `a`, 1, 1, StrType),
		NewTestCase("'a'", "a", 1, 1, StrType),
		//NewTestCase(`'a\'b'`, `a'b`, 1, 1, StrType),
	}
	RunTests(t, rules, tests)
}
