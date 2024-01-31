# scan

An all-purpose scanner for your lexical needs.

## Introduction

A scanner is used to extract the fundamental and basic elements of a structured
text document. These basic elements are called tokens and represent the atomic
units used when parsing a document. As an example, we will use the following
GeoJSON file that is featured on the front page of https://geojson.org:

```json
{
  "type": "Feature",
  "geometry": {
    "type": "Point",
    "coordinates": [125.6, 10.1]
  },
  "properties": {
    "name": "Dinagat Islands"
  }
}
```

Some things in this document don't matter when parsing. Whitespace is one of
those things. This document is indented by two spaces but it could have been
four, eight, or even zero spaces without changing its meaning.

Some things really do matter. The brackets that are used to define an object.
The strings that are enclosed by double quotes. The colon that separates a
property name from its value. These are the tokens of the document. Parsing
a document is much easier when working with a stream of tokens instead of
a individual characters.

Let's look at the third line of the example. It starts with a string token with
the value of `geometry`. It is important to capture both bits of information--
its value of `geometry` and its type of `string`. A parser for this document is
interested if the line starts with a string token or not but is not concerned
about its value. Users of the parse result will be interested in the value.

When there are errors parsing the document, good error messages are useful to
help identify where problems are located. Recording the line and column number
with each token is important for this reporting.

The examples directory contains a JSON scanner using the scanner provided by
this package. Run it with the following:

A table should be printed to the terminal with the stream of tokens found
in the file. Its output starts with the following:

```
Pos                 Type  Value              Literal
example.json:1:1    {     "{"                "{"
example.json:2:5    str   "type"             '"type"'
example.json:2:11   :     ":"                ":"
example.json:2:13   str   "Feature"          '"Feature"'
example.json:2:22   ,     ","                ","
example.json:3:5    str   "geometry"         '"geometry"'
example.json:3:15   :     ":"                ":"
example.json:3:17   {     "{"                "{"
example.json:4:9    str   "type"             '"type"'
```

Building a scanner manually isn't difficult but it can be a bit tedious at
times. This scanner provides low level functions for iterating over a stream
and for building tokens. It also has higher level functions that are useful for
chaining together common scanning patterns. Scanners for files that follow
standard patterns can be constructed easily using the definitions found in this
package. The example JSON scanner is about 30 lines in length. Complex or
non-standard cases can stick with the lower level functions which still provide
a lot of value in getting a scanner up and running quickly.

## Creating a Scanner

A scanner is initialized with a filename followed by an input source which can
be either a `Reader`, a `string` or a `byte` slice. Use one of the
corresponding `New` functions to create a scanner or use an `Init` function to
initialize a zero value scanner or to reset an existing scanner.

Use either:

```go
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close()
    s := scan.NewScanner(filename, f)
```

or:

```go
    var s scan.Scanner
    s.InitString("", "hello world")
```

## `Keep` and `Emit`

When using the lower level functions, a stream is processed one rune at a time.
A buffer is used to accumulate these runes until a complete token has been
seen. The basic scanning loop is to decide which of the following functions
should be called for each rune:

- `Keep`: add the rune to the current token buffer and advance the scanner to
the next rune; or
- `Emit`: return the current token and reset the token buffer

The current rune is accessed with the `This` field and the `HasMore` method
returns `true` as long as there is more data in the stream. Scanning a single
binary number from an input stream can be done with the following:

```go
	s := scan.NewScannerFromString("", "1010234abc!")
	for s.HasMore() {
		if s.This == '0' || s.This == '1' {
			s.Keep()
			continue
		}
		break
	}
	tok := s.Emit()
	fmt.Println(tok.Val)

	// Output:
	// 1010
```
[Example 1](doc/example_1_test.go)

This scans the input stream while it contains zeros on ones. As soon as it
finds something else, it breaks out of the loop, emits the token, and prints
the value. The contents for `235abc!` are not scanned.

## Classes

Now we will update the example to scan for an unsigned integer in base 10. The
only change to the code is to check for digits 0 through 9 instead of just 0
and 1. This check can be placed into a separate function that takes
a `rune` as input and returns a `bool` indicating if the rune is one of those
digits. This type of function is a `Class` function and it checks to see
if a rune is member of that class. Example:

```go
    digit09 := func(r rune) bool {
        return r >= '0' && r <= '9'
    }

    s := scan.NewScannerFromString("", "1010234abc!")
    for s.HasMore() && s.Is(digit09) {
        s.Keep()
    }
    tok := s.Emit()
    fmt.Println(tok.Val)

    // Output:
    // 1010234

```
[Example 2](doc/example_2_test.go)

### Class Functions

The scanner package provides some functions to easily construct various
rune classes:

#### `Rune`

Creates a `Class` function that returns true if it matches any of the
specified runes. An example of a rune class that checks for common
whitespace runes:

```go
    whitespace := scan.Rune(' ', '\n', '\t', '\f', '\r')
```

#### `Range`

Creates a `Class` function that returns true if it matches a range of runes
as defined by their code point values. An example of a rune class that checks
for digits between 0 and 9:

```go
    digit09 := scan.Range('0', '9')
```

#### `Or`

Creates a `Class` function that returns true if any of the specified classes
returns true. An example of a rune class that checks for hexadecimal digits:

```go
    digit0F := scan.Or(
        scan.Range('0', '9'),
        scan.Range('A', 'F'),
        scan.Range('a', 'f'),
    )
```

#### `Not`

Creates a `Class` function that negates another class. For example, to create
a class that accepts all non-whitespace runes:

```go
    notWhitespace := scan.Not(whitespace)
```

### Predefined Classes

Common rune classes are provided by this package and can be used instead. These
classes return `true` under the following conditions:

- `Any`: always, as long as there are more runes available in the stream
- `Currency`: currency symbol as defined by Unicode
- `Digit`: digit as defined by Unicode
- `Digit01`: binary digits
- `Digit07`: octal digits
- `Digit09`: decimal digits
- `Digit0F`: hexadecimal digits
- `Letter`: letter as defined by Unicode
- `LetterAZ`: letters `A` through `Z` and `a` through `z`
- `LetterUnder`: letter or an underscore
- `LetterDigitUnder`: letter, digit, or an underscore
- `None`: never, always returns `false`
- `Printable`: printable as defined by Unicode
- `Rune8`: a code point that can fit in an 8-bit number
- `Rune16`: a code point that can fit in a 16-bit number
- `Whitespace`: whitespace as defined by Unicode.

### `Is`

The current rune found in `This` can be tested with a class using the `Is`
method of the scanner:

```go
    s.Is(scan.Whitespace)
```

This is equivalent to:

```go
    scan.Whitespace(s.This)
```

The advantage to using `Is` is that the function will return `false` when
passed a `nil` class function.

### Updated Example

The previous example can now be updated using the predefined class for
base-10 numbers and the `Is` method:

```go
    s := scan.NewScannerFromString("1010234abc!")
    for s.HasMore() && s.Is(scan.Digit09) {
        s.Keep()
    }
    tok := s.Emit()
    fmt.Println(tok.Val)

    // Output:
    // 1010234
```
[Example 3](doc/example_3_test.go)

## Loop Functions

A common loop while scanning is to iterate over the stream while there is more data and then break out of the loop once a condition is met. The following
functions can be used in simple loop operations:

#### `While`

Consume input from the stream while characters belong to a certain class.
The function takes a `Scanner`, a `Class` used to check membership, and
a function that should be called for each matching character. For example,
to add characters to the current token while digits are seen:

```go
    scan.While(s, scan.Digit09, s.Keep)
```

#### `Until`

This function is similar to `While` but consumes characters from the stream
until a character belongs to the provided class. For example, to add any
non-whitespace characters to the curren token:

```go
    scan.Until(s, scan.Whitespace, s.Keep)
```

#### `Repeat`

Sometimes it is necessary to call a function like `Keep` a fixed number of
times to advance the scanner. Use the `Repeat` function for this:

```go
    scan.Repeat(s.Keep, 3)
```

### Updated Example

The example can now be updated to remove the loop with a call to `While`:

```go
    s := scan.NewScannerFromString("", "1010234abc!")
    scan.While(s, scan.Digit09, s.Keep)
    tok := s.Emit()
    fmt.Println(tok.Val)

    // Output:
    // 1010234
```
[Example 4](doc/example_4_test.go)

## `Discard`

For runes like whitespace, it useful to simply discard them without creating
tokens. Use the `Discard` method to reset the current token buffer, discard the
current rune, and advance the scanner to the next rune. The example below
now discards any whitespace found at the beginning of the stream:

```go
    s := scan.NewScannerFromString("", " \t 1010234abc!")
    scan.While(s, scan.Whitespace, s.Discard)
    scan.While(s, scan.Digit09, s.Keep)
    tok := s.Emit()
    fmt.Println(tok.Val)

    // Output:
    // 1010234
```
[Example 5](doc/example_5_test.go)
```

