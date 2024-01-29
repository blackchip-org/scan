package scan

import "fmt"

func ExampleIsAny() {
	fmt.Printf("%c: %v\n", 'a', Any('a'))
	fmt.Printf("%c: %v\n", '4', Any('4'))
	// Output:
	// a: true
	// 4: true
}

func ExampleIsCurrency() {
	fmt.Printf("%c: %v\n", '$', Currency('$'))
	fmt.Printf("%c: %v\n", '€', Currency('€'))
	fmt.Printf("%c: %v\n", '!', Currency('!'))
	// Output:
	// $: true
	// €: true
	// !: false
}

func ExampleIsDigit() {
	fmt.Printf("%c: %v\n", '1', Digit('1'))
	fmt.Printf("%c: %v\n", '६', Digit('६'))
	fmt.Printf("%c: %v\n", 'V', Digit('V'))
	// Output:
	// 1: true
	// ६: true
	// V: false
}

func ExampleIsDigit01() {
	fmt.Printf("%c: %v\n", '1', Digit01('1'))
	fmt.Printf("%c: %v\n", '2', Digit01('2'))
	// Output:
	// 1: true
	// 2: false
}

func ExampleIsDigit07() {
	fmt.Printf("%c: %v\n", '7', Digit07('7'))
	fmt.Printf("%c: %v\n", '8', Digit07('8'))
	// Output:
	// 7: true
	// 8: false
}

func ExampleIsDigit09() {
	fmt.Printf("%c: %v\n", '9', Digit09('9'))
	fmt.Printf("%c: %v\n", 'a', Digit09('a'))
	// Output:
	// 9: true
	// a: false
}

func ExampleIsDigit0F() {
	fmt.Printf("%c: %v\n", '9', Digit0F('9'))
	fmt.Printf("%c: %v\n", 'a', Digit0F('a'))
	fmt.Printf("%c: %v\n", 'A', Digit0F('A'))
	fmt.Printf("%c: %v\n", 'g', Digit0F('g'))
	// Output:
	// 9: true
	// a: true
	// A: true
	// g: false
}

func ExampleIsLetter() {
	fmt.Printf("%c: %v\n", 'á', Letter('á'))
	fmt.Printf("%c: %v\n", '%', Letter('%'))
	// Output:
	// á: true
	// %: false
}

func ExampleIsLetterAZ() {
	fmt.Printf("%c: %v\n", 'f', LetterAZ('f'))
	fmt.Printf("%c: %v\n", 'F', LetterAZ('F'))
	fmt.Printf("%c: %v\n", '4', LetterAZ('4'))
	// Output:
	// f: true
	// F: true
	// 4: false
}

func ExampleIsNone() {
	fmt.Printf("%c: %v\n", 'f', None('f'))
	fmt.Printf("%c: %v\n", 'F', None('F'))
	fmt.Printf("%c: %v\n", '4', None('4'))
	// Output:
	// f: false
	// F: false
	// 4: false
}

func ExampleNot() {
	isA := Rune('a')
	isNotA := Not(isA)
	fmt.Printf("%c: %v\n", 'a', isNotA('a'))
	fmt.Printf("%c: %v\n", 'b', isNotA('b'))
	// Output:
	// a: false
	// b: true
}

func ExampleOr() {
	isLowerAZ := Range('a', 'z')
	isUpperAZ := Range('A', 'Z')
	isLetterAZ := Or(isLowerAZ, isUpperAZ)
	fmt.Printf("%c: %v\n", 'f', isLetterAZ('f'))
	fmt.Printf("%c: %v\n", 'F', isLetterAZ('F'))
	fmt.Printf("%c: %v\n", '4', isLetterAZ('4'))
	// Output:
	// f: true
	// F: true
	// 4: false
}

func ExampleRange() {
	isDigit09 := Range('0', '9')
	fmt.Printf("%c: %v\n", '3', isDigit09('3'))
	fmt.Printf("%c: %v\n", '6', isDigit09('6'))
	fmt.Printf("%c: %v\n", 'a', isDigit09('a'))
	// Output:
	// 3: true
	// 6: true
	// a: false
}

func ExampleRune() {
	isAB := Rune('a', 'b')
	fmt.Printf("%c: %v\n", 'a', isAB('a'))
	fmt.Printf("%c: %v\n", 'b', isAB('b'))
	fmt.Printf("%c: %v\n", 'c', isAB('c'))
	// Output:
	// a: true
	// b: true
	// c: false
}

func ExampleWhitespace() {
	fmt.Printf("%c: %v\n", ' ', Whitespace(' '))
	fmt.Printf("%c: %v\n", 'n', Whitespace('n'))
	// Output:
	//  : true
	// n: false
}
