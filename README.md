# scan

An all-purpose scanner for your lexical needs.

[![Go Reference](https://pkg.go.dev/badge/github.com/blackchip-org/scan.svg)](https://pkg.go.dev/github.com/blackchip-org/scan)

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
individual characters.

Let's look at the third line of the example. It starts with a string token with
the value of `geometry`. It is important to capture both bits of information--
its value of `geometry` and its type of `string`. A parser for this document is
interested if the next token is a string or not but is not concerned about its
value. Users of the parse result will be interested in the value.

When there are errors scanning and parsing the document, good error messages
are useful to help identify where problems are located. Recording the line and
column number with each token is important for this reporting.

The examples directory contains a JSON scanner using the scanner provided by
this package. Run it with the following:

    go run cmd/scan-json/scan-json.go example.json

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
package. Complex or non-standard cases can stick with the lower level functions
which still provide a lot of value in getting a scanner up and running quickly.

## Creating a Scanner

A scanner is initialized with a filename followed by an input source which can
be either an `io.Reader`, a `string` or a `byte` slice. If the filename is not
important or doesn't apply to the input, it can set to a blank value. Use one
of the corresponding new functions to create a scanner or use an init function
to either initialize a zero value scanner or to reset an existing scanner.

An example of using new to read from a file:

```go
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close()
    s := scan.NewScanner(filename, f)
```

An example of using init to read from a string:

```go
    var s scan.Scanner
    s.InitString("", "hello world")
```

## Basic Run Loop

When using the lower level functions, a stream is processed one rune at a time.
A buffer is used to accumulate these runes until a complete token has been
seen. The basic scanning loop is to look at the upcoming rune in the input
and to decide to either keep the rune for the current token or to emit the
token. The following fields and methods are used for this:

#### `s.This`

The rune at the current position in the input stream.

#### `s.Next`

The next rune found in the input stream. When the scanner advances, this value
will become `s.This`

#### `s.HasMore()`

Returns true if there is more runes in the input stream. Returns false when
the end of file or an error condition is encountered.

#### `s.Keep()`

Adds the current rune to the token buffer and advances the scanner to
the next rune.

#### `s.Emit()`

Returns the current token and resets the token buffer.

### Example

Scanning a single binary number from an input stream can be done by looping
over the stream while it has data and adding to the token buffer as long as the
current rune is a zero or a one:

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
[Example 1](examples/readme/example_1_test.go)

Once the scanner encounters a non-binary digit, it breaks out of the loop and
emits the token accumulated so far. The remaining contents of the stream,
"234abc!", are not read.

## Classes

Now we will update the example to scan for an unsigned integer in base 10. The
only change necessary is to now check for digits 0 through 9 instead of just 0
and 1. This check can be placed into a separate function that takes a `rune` as
input and returns a `bool` indicating if the rune is one of those digits. This
type of function is a `scan.Class` function and it is used to determine if a
rune is member of that class. Example:

```go
    isDigit09 := func(r rune) bool {
        return r >= '0' && r <= '9'
    }

    s := scan.NewScannerFromString("", "1010234abc!")
    for s.HasMore() && isDigit09(s.This) {
        s.Keep()
    }
    tok := s.Emit()
    fmt.Println(tok.Val)

    // Output:
    // 1010234
```
[Example 2](examples/readme/example_2_test.go)

### Class Functions

The scanner package provides some functions to easily construct various
rune classes:

#### `scan.Rune()`

Creates a `scan.Class` function that returns true if it matches any of the
specified runes. An example of a rune class that checks for common
whitespace:

```go
    isSpace := scan.Rune(' ', '\n', '\t', '\f', '\r')
```

#### `scan.Range()`

Creates a `scan.Class` function that returns true if it matches a range of
runes as defined by their code point values. An example of a rune class that
checks for digits between 0 and 9:

```go
    isDigit09 := scan.Range('0', '9')
```

#### `scan.Or()`

Creates a `scan.Class` function that returns true if any of the specified
classes returns true. An example of a rune class that checks for hexadecimal
digits:

```go
    isDigit0F := scan.Or(
        scan.Range('0', '9'),
        scan.Range('A', 'F'),
        scan.Range('a', 'f'),
    )
```

#### `scan.Not()`

Creates a `scan.Class` function that negates another class. For example, to
create a class that accepts all non-whitespace runes:

```go
    isNotSpace := scan.Not(isSpace)
```

### Predefined Classes

Common rune classes are provided by this package and can be used instead. These
classes return `true` under the following conditions:

- `scan.IsAny`: always, as long as there are more runes available in the stream
- `scan.IsCurrency`: currency symbol as defined by Unicode
- `scan.IsDigit`: digit as defined by Unicode
- `scan.IsDigit01`: binary digits
- `scan.IsDigit07`: octal digits
- `scan.IsDigit09`: decimal digits
- `scan.IsDigit0F`: hexadecimal digits
- `scan.IsLetter`: letter as defined by Unicode
- `scan.IsLetterAZ`: letters `A` through `Z` and `a` through `z`
- `scan.IsLetterUnder`: letter or an underscore
- `scan.IsLetterDigitUnder`: letter, digit, or an underscore
- `scan.IsNone`: never, always returns `false`
- `scan.IsPrintable`: printable as defined by Unicode
- `scan.IsRune8`: a code point that can fit in an 8-bit number
- `scan.IsRune16`: a code point that can fit in a 16-bit number
- `scan.IsWhitespace`: whitespace as defined by Unicode.

### Updated Example

The previous example can now be updated using the predefined class for
base 10 numbers and the `s.Is()` method:

```go
    s := scan.NewScannerFromString("", "1010234abc!")
    for s.HasMore() && scan.IsDigit09(s.This) {
        s.Keep()
    }
    tok := s.Emit()
    fmt.Println(tok.Val)

    // Output:
    // 1010234

```
[Example 3](examples/readme/example_3_test.go)

## Loop Functions

A common loop while scanning is to iterate over the stream while there is more
data and then break out of the loop once a condition is met. The following
functions can be used in simple loop operations:

#### `scan.While()`

Consume input from the stream while characters belong to a certain class. The
function takes a `*scan.Scanner`, a `scan.Class` used to check membership, and
a function that should be called for each matching character. For example, to
add characters to the current token while digits are seen:

```go
    scan.While(s, scan.IsDigit09, s.Keep)
```

#### `scan.Until()`

This function is similar to `scan.While()` but consumes characters from the
stream until a character belongs to the provided class. For example, to add any
non-whitespace characters to the curren token:

```go
    scan.Until(s, scan.IsSpace, s.Keep)
```

#### `scan.Repeat()`

Sometimes it is necessary to call a function like `scan.Keep()` a fixed number
of times to advance the scanner. Use the `scan.Repeat()` function for this:

```go
    scan.Repeat(s.Keep, 3)
```

### Updated Example

The example can now be updated to remove the loop with a call to
`scan.While()`:

```go
    s := scan.NewScannerFromString("", "1010234abc!")
    scan.While(s, scan.IsDigit09, s.Keep)
    tok := s.Emit()
    fmt.Println(tok.Val)

    // Output:
    // 1010234
```
[Example 4](examples/readme/example_4_test.go)

## Discarding

For runes like whitespace, it useful to simply discard them without creating
tokens.

#### `s.Discard()`

Resets the current token buffer, discards the current rune, and advances the
scanner to the next rune. The example below now discards any whitespace found
at the beginning of the stream:

```go
    s := scan.NewScannerFromString("", " \t 1010234abc!")
    scan.While(s, scan.IsSpace, s.Discard)
    scan.While(s, scan.IsDigit09, s.Keep)
    tok := s.Emit()
    fmt.Println(tok.Val)

    // Output:
    // 1010234
```
[Example 5](examples/readme/example_5_test.go)

## Types

The next example is going to scan the entire stream and create tokens for each
integer and word that is found. A word is going to be any sequence of letters
found in the stream. Integers and words must be separated by whitespace. We
will now place each token type will into a separate function and each token
will now be identified with a type. A type can be any string value but there
are some constants already defined in the scan package. The functions return a
boolean value to indicate if there was a match for that token.

The function for the integer is as follows:

```go
    func scanInt(s *scan.Scanner) bool {
        if !scan.IsDigit09(s.This) {
            return false
        }
        s.Type = scan.IntType
        scan.While(s, scan.IsDigit09, s.Keep)
        return true
    }
```
[Example 6](examples/readme/example_6_test.go)

The function for the word is similar and is as follows:

```go
    func scanWord(s *scan.Scanner) bool {
        if !scan.IsLetter(s.This) {
            return false
        }
        s.Type = scan.WordType
        scan.While(s, scan.IsLetter, s.Keep)
        return true
    }
```
[Example 6](examples/readme/example_6_test.go)

Since whitespace is now significant, it will be returned as a token instead
of being discarded:

```go
    func scanSpace(s *scan.Scanner) bool {
        if !scan.IsSpace(s.This) {
            return false
        }
        s.Type = scan.SpaceType
        scan.While(s, scan.IsSpace, s.Keep)
        return true
    }
```
[Example 6](examples/readme/example_6_test.go)

Now create a main loop that checks each function for a match and then collects
all of the tokens seen along the way:

```go
    scanFuncs := []func(*scan.Scanner) bool{
        scanSpace,
        scanWord,
        scanInt,
    }

    var toks []scan.Token
    s := scan.NewScannerFromString("example6", "abc 123 !@#  \tdef456")
    for s.HasMore() {
        match := false
        for _, fn := range scanFuncs {
            if match = fn(s); match {
                toks = append(toks, s.Emit())
                break
            }
        }
        if !match {
            scan.Until(s, scan.IsSpace, s.Keep)
            s.Illegal("unexpected %v", scan.Quote(s.Val.String()))
            toks = append(toks, s.Emit())
        }
    }
    fmt.Println(scan.FormatTokenTable(toks))

    // Output:
    //
    // Pos            Type     Value         Literal
    // example6:1:1   word     "abc"         "abc"
    // example6:1:4   space    " "           " "
    // example6:1:5   int      "123"         "123"
    // example6:1:8   space    " "           " "
    // example6:1:9   illegal  "!@#"         "!@#"
    // example6:1:12: error: unexpected "!@#"
    // example6:1:12  space    "  {!ch:\t}"  "  {!ch:\t}"
    // example6:1:15  word     "def"         "def"
    // example6:1:18  int      "456"         "456"

```
[Example 6](examples/readme/example_6_test.go)

For each token that is generated, we can check the `tok.Type` field to see if
it is an int, word, space, or an illegal token. The `tok.Pos` field contains
the name of the file (or item) being scanned, the line number, and the column
number of where the token starts.

Note that the role of the scanner is to simply emit tokens of the various types
seen in the input stream. The "def456" string in the stream is illegal by this
example's definition because a space must separate the word and the integer.
That check should be handled by the parser that is processing the tokens, not
by the scanner itself.

There are a few new functions or methods introduced in this example:

#### `s.Illegal()`

This function has the same arguments as `fmt.Printf()` and attaches an error
to the token which can be found in `tok.Errs`. It also sets the token type
field to `scan.IllegalType`.

#### `scan.Quote()`

This function surrounds a value with double quotes to make it clear where the
value starts and ends. If the value contains double quotes, then single quotes
will be used instead. If the value contains both double and single quotes, then
backticks are used. In the rare event that all quotes are in use, the value
is prefixed with `{!quote:` and ends with `}`.

Non-printable characters found in values are escaped with a prefix of
`{!ch:` and ends with `}`. For characters that are normally escaped, they
will be displayed like `{!ch:\n}` otherwise they are shown with a
hexadecimal number like `{!ch:a0}`

#### `scan.FormatTokenTable()`

This function takes a slice of tokens and returns a string with the tokens
formatted in a table to be used for display. For each token, it shows the
position of where the token starts, the type of token, and the value
and literal. In this example, the value and literal are always the same but
that does not always have to be the case. If there are errors attached to a
token, those are listed on subsequent lines.

## Rules

Scanning for integers, words, and spaces are common things to scan for in a
text document. In the next example, the scanning functions are modified to
allow configuration of what is considered a digit, a letter, and a space.
Instead of defining a function, a struct will now be used that contains the
configuration and the scanning function will now be a method called `Eval`.

These structs will be called *rules* and they satisfy the `scan.Rule`
interface. Changing the integer function to a rule looks like this:

```go
    type IntRule struct {
        isDigit scan.Class
    }

    func NewIntRule(isDigit scan.Class) IntRule {
        return IntRule{isDigit: isDigit}
    }

    func (r IntRule) Eval(s *scan.Scanner) bool {
        if !r.isDigit(s.This) {
            return false
        }
        s.Type = scan.IntType
        scan.While(s, r.isDigit, s.Keep)
        return true
    }
```
[Example 7](examples/readme/example_7_test.go)

With this change, this one rule can be used for scanning decimal, binary,
octal, and hexadecimal numbers simply by calling new with the corresponding
digit class .The rules for word and whitespace looks similar and aren't
show here. The full source code for the example can be found in the
Example 7 links.

Now a function is needed to handle the case where a token does not match. This
also has to be configurable to allow different definitions of whitespace. This
will be a function that returns a function that scans the unexpected value
until the specified class is found:

```go
func UnexpectedUntil(c scan.Class) func(*scan.Scanner) {
    return func(s *scan.Scanner) {
        scan.Until(s, c, s.Keep)
        s.Illegal("unexpected %s", scan.Quote(s.Lit.String()))
    }
}
```

These new rules and the no match function can now be bundled up with a
`scan.RuleSet`. Then a `scan.Runner` can be created with the scanner and rule
set to iterate through the tokens using `runner.HasMore` and `runner.Next`. Or
all the tokens can be obtained by calling `runner.All()`:

```go
    rules := scan.NewRuleSet(
        NewSpaceRule(scan.Whitespace),
        NewWordRule(scan.Letter),
        NewIntRule(scan.Digit),
    ).WithNoMatchFunc(UnexpectedUntil(scan.Whitespace))

    s := scan.NewScannerFromString("example7", "abc 123 !@#  \tdef456")
    runner := scan.NewRunner(s, rules)
    toks := runner.All()
    fmt.Println(scan.FormatTokenTable(toks))

    // Output:
    //
    // Pos            Type     Value         Literal
    // example7:1:1   word     "abc"         "abc"
    // example7:1:4   space    " "           " "
    // example7:1:5   int      "123"         "123"
    // example7:1:8   space    " "           " "
    // example7:1:9   illegal  "!@#"         "!@#"
    // example7:1:12: error: unexpected "!@#"
    // example7:1:12  space    "  {!ch:\t}"  "  {!ch:\t}"
    // example7:1:15  word     "def"         "def"
    // example7:1:18  int      "456"         "456"
```
[Example 7](examples/readme/example_7_test.go)

## Preprocessing

Sometimes it is convenient to perform some light preprocessing of the token
value before passing it off to the parser. A good example of this is allowing
digit separators in an integer value such as "1,000". The parser is most likely
going to pass this responsibility over to `strconv.ParseInt()` and the comma
gets in the way. Another use case is to convert runes to another rune. For
example, converting all whitespace runes to a single space value or changing
newlines to semicolons. The following can be used for these cases:

#### `s.Skip()`

Adds the current rune to the token's literal but not to the token's value. The
scanner then advances to the next rune.

#### `s.Value`

This is a string builder containing the token value seen so far. If newlines
need to be converted to semicolons, this can be done in the following way:

```go
    if s.This == '\n' {
        s.Skip()
        s.Value.WriteRune(';')
    }
```

Now the integer rule can be updated to include a digit separator. The digit
class is mandatory but the digit separator is optional so only the digit class
is found in the constructor. The digit separator can be configured with a
chaining function:

```go
type IntRule2 struct {
    isDigit    scan.Class
    isDigitSep scan.Class
}

func NewIntRule2(isDigit scan.Class) IntRule2 {
    return IntRule2{
        isDigit:    isDigit,
        isDigitSep: scan.IsNone,
    }
}

func (r IntRule2) WithDigitSep(isDigitSep scan.Class) IntRule2 {
    r.isDigitSep = isDigitSep
    return r
}

func (r IntRule2) Eval(s *scan.Scanner) bool {
    if !r.isDigit(s.This) {
        return false
    }
    s.Keep()
    for s.HasMore() {
        if r.isDigit(s.This) {
            s.Keep()
        } else if r.isDigitSep(s.This) {
            s.Skip()
        } else {
            break
        }
    }
    return true
}
```
[Example 8](examples/readme/example_8_test.go)

Now update the rule set with the new rule and modify the input to include some
digit separators:

```go
    rules := scan.NewRuleSet(
        NewSpaceRule(scan.IsSpace),
        NewWordRule(scan.IsLetter),
        NewIntRule2(scan.IsDigit).
            WithDigitSep(scan.Rune(',')),
    ).WithNoMatchFunc(UnexpectedUntil(scan.IsSpace))

    s := scan.NewScannerFromString("example8", "abc 1,234 !@#  \tdef45,678")
    runner := scan.NewRunner(s, rules)
    toks := runner.All()
    fmt.Println(scan.FormatTokenTable(toks))

    // Output:
    //
    // Pos            Type     Value         Literal
    // example8:1:1   word     "abc"         "abc"
    // example8:1:4   space    " "           " "
    // example8:1:5   1234     "1234"        "1,234"
    // example8:1:10  space    " "           " "
    // example8:1:11  illegal  "!@#"         "!@#"
    // example8:1:14: error: unexpected "!@#"
    // example8:1:14  space    "  {!ch:\t}"  "  {!ch:\t}"
    // example8:1:17  word     "def"         "def"
    // example8:1:20  45678    "45678"       "45,678"
```
[Example 8](examples/readme/example_8_test.go)

## Predefined Rules

The scan package already has rules defined for reading the token types used
so far. The example can now be updated to use those rules instead:

```go
    rules := scan.NewRuleSet(
        scan.KeepSpaceRule,
        scan.NewWhileRule(scan.IsLetter, scan.WordType),
        scan.IntRule.WithDigitSep(scan.Rune(',')),
    ).WithNoMatchFunc(scan.UnexpectedUntil(scan.IsSpace))

    s := scan.NewScannerFromString("example9", "abc 1,234 !@#  \tdef45,678")
    runner := scan.NewRunner(s, rules)
    toks := runner.All()
    fmt.Println(scan.FormatTokenTable(toks))

    // Output:
    //
    // Pos            Type     Value         Literal
    // example9:1:1   word     "abc"         "abc"
    // example9:1:4   space    " "           " "
    // example9:1:5   int      "1234"        "1,234"
    // example9:1:10  space    " "           " "
    // example9:1:11  illegal  "!@#"         "!@#"
    // example9:1:14: error: unexpected "!@#"
    // example9:1:14  space    "  {!ch:\t}"  "  {!ch:\t}"
    // example9:1:17  word     "def"         "def"
    // example9:1:20  int      "45678"       "45,678"
```
[Example 9](examples/readme/example_9_test.go)

The predefined rules are as follows:

* `scan.CharEncRule`: For converting character encodings such as `\n` to their
actual values like `U+000A`
* `scan.ClassRule`: For runes that are a member of a class.
* `scan.CommentRule`: Comments which can be configured, at runtime, to emit
tokens or discard the input.
* `scan.HexEncRule`: For converting hex escape sequences such as `\x7f` to
their actual values.
* `scan.IdentRule`: Identifiers used in programming languages that usually
need two characters classes: one for the first character, and one for the
rest. Can also be configured with a list of keywords.
* `scan.LiteralRule`: For matching exact strings in the input. Useful for
matching against operators, punctuation, and keywords when not using an
`scan.IdentRule`.
* `scan.NumRule`: Integer and real number values. This is quite configurable
for most standard needs but it doesn't cover every possibility that can
be found in the wild.
* `scan.OctEncRule`: For converting octal escape sequences such as `\377` to
their actual values.
* `scan.StrRule`: For scanning strings which are bracketed by a start and
ending rune. Can be configured to include escape sequences and multiline
behavior.
* `scan.SpaceRule`: For discarding whitespace but can be configured to emit
tokens if needed.
* `scan.WhileRule`: Scanning a token while a character class is true.

Some rules have predefined configurations as well. For example `scan.Hex0xRule`
that uses a `scan.NumRule` configured for hexadecimal digits that are
prefixed with `0x`. See the API documentation for more information.

## Lookahead and lookbehind

Most scanning operations can be done by looking at the current character,
which is stored in `s.This` and the next character in the stream that is
stored in `s.Next`. In some cases, it can be easier to look further ahead
in the stream or to even speculate and advance the scanner and bail out
if the stream doesn't contain what is expected. Use the following for these
cases:

#### `s.Peek()`

Looks ahead in the stream. It takes an integer value which is the number of
runes ahead of the current rule it should return. Using `s.Peek(0)` is
the same as `s.This` and `s.Peek(1)` is the same as `s.Next`. The lookahead
value can be as large as needed. If a peek operation would go past the end
of the stream, the value of `scan.EndOfText` is returned. A negative value
can be used to look back at the literal collected so far. For example,
`s.Peek(-1)` would return the rune that was scanned before the current one.
If a negative value goes past the length of the current literal, a value
of `scan.EndOfText` is returned.

#### `s.Undo()`

When speculating with the scanner, calling this method will place the token
literal seen so far back onto the input stream and then clear the token
buffer.

## Full Examples

There are two full examples provided with this package. The first is a
[JSON scanner](scanjson/scanjson.go). The scanner for JSON is quite simple but
JSON does have some quirks:

* Numbers can start with a negative sign, but not a positive sign.
* Numbers cannot have leading zeros. "303" is okay, but "0303" is not.
* Real numbers cannot have empty parts. "0.3" and "3.0" is valid, but ".3" and "3." is not.

There is also a [Go scanner](scango/scango.go). This one is a bit more
complicated but is mostly constructed from predefined rules. The exceptions
are for a rule to handle imaginary numbers and a post token processor for
automatic semicolon insertion.

In a separate repository, there is an example of a scanner, parser, and a
formatter for [DMS angles](https://github.com/blackchip-org/dms).

## Status

This package is still a work in progress and is subject to change. If you
find this useful, please drop a line at zc at blackchip dot org.

## License

MIT
