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
be either a `io.Reader`, a `string` or a `byte` slice. Use one of the
corresponding new functions to create a scanner or use an init function to
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

## Basic Run Loop

When using the lower level functions, a stream is processed one rune at a time.
A buffer is used to accumulate these runes until a complete token has been
seen. The basic scanning loop is to decide which of the following functions
should be called for each rune:

#### `s.Keep()`

Adds the current rune to the token buffer and advances the scanner to
the next rune.

#### `s.Emit()`

Returns the current token and resets the token buffer.

### Example

The current rune is accessed with the `s.This` field and the `s.HasMore()`
method returns `true` as long as there is more data in the stream. Scanning a
single binary number from an input stream can be done with the following:

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
and 1. This check can be placed into a separate function that takes a `rune` as
input and returns a `bool` indicating if the rune is one of those digits. This
type of function is a `scan.Class` function and it checks to see if a rune is
member of that class. Example:

```go
    digit09 := func(r rune) bool {
        return r >= '0' && r <= '9'
    }

    s := scan.NewScannerFromString("", "1010234abc!")
    for s.HasMore() && digit09(s.This) {
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

#### `scan.Rune()`

Creates a `scan.Class` function that returns true if it matches any of the
specified runes. An example of a rune class that checks for common
whitespace:

```go
    whitespace := scan.Rune(' ', '\n', '\t', '\f', '\r')
```

#### `scan.Range()`

Creates a `scan.Class` function that returns true if it matches a range of
runes as defined by their code point values. An example of a rune class that
checks for digits between 0 and 9:

```go
    digit09 := scan.Range('0', '9')
```

#### `scan.Or()`

Creates a `scan.Class` function that returns true if any of the specified
classes returns true. An example of a rune class that checks for hexadecimal
digits:

```go
    digit0F := scan.Or(
        scan.Range('0', '9'),
        scan.Range('A', 'F'),
        scan.Range('a', 'f'),
    )
```

#### `scan.Not()`

Creates a `scan.Class` function that negates another class. For example, to
create a class that accepts all non-whitespace runes:

```go
    notWhitespace := scan.Not(whitespace)
```

### Predefined Classes

Common rune classes are provided by this package and can be used instead. These
classes return `true` under the following conditions:

- `scan.Any`: always, as long as there are more runes available in the stream
- `scan.Currency`: currency symbol as defined by Unicode
- `scan.Digit`: digit as defined by Unicode
- `scan.Digit01`: binary digits
- `scan.Digit07`: octal digits
- `scan.Digit09`: decimal digits
- `scan.Digit0F`: hexadecimal digits
- `scan.Letter`: letter as defined by Unicode
- `scan.LetterAZ`: letters `A` through `Z` and `a` through `z`
- `scan.LetterUnder`: letter or an underscore
- `scan.LetterDigitUnder`: letter, digit, or an underscore
- `scan.None`: never, always returns `false`
- `scan.Printable`: printable as defined by Unicode
- `scan.Rune8`: a code point that can fit in an 8-bit number
- `scan.Rune16`: a code point that can fit in a 16-bit number
- `scan.Whitespace`: whitespace as defined by Unicode.

### `s.Is()`

The current rune found in `s.This` can be tested with a class using the `s.Is`
method of the scanner:

```go
    s.Is(scan.Whitespace)
```

This is equivalent to:

```go
    scan.Whitespace(s.This)
```

The advantage to using `s.Is()` is that the function will return `false` when
passed a `nil` class function.

### Updated Example

The previous example can now be updated using the predefined class for
base-10 numbers and the `s.Is()` method:

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

A common loop while scanning is to iterate over the stream while there is more
data and then break out of the loop once a condition is met. The following
functions can be used in simple loop operations:

#### `scan.While()`

Consume input from the stream while characters belong to a certain class. The
function takes a `*scan.Scanner`, a `scan.Class` used to check membership, and
a function that should be called for each matching character. For example, to
add characters to the current token while digits are seen:

```go
    scan.While(s, scan.Digit09, s.Keep)
```

#### `scan.Until()`

This function is similar to `scan.While()` but consumes characters from the
stream until a character belongs to the provided class. For example, to add any
non-whitespace characters to the curren token:

```go
    scan.Until(s, scan.Whitespace, s.Keep)
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
    scan.While(s, scan.Digit09, s.Keep)
    tok := s.Emit()
    fmt.Println(tok.Val)

    // Output:
    // 1010234
```
[Example 4](doc/example_4_test.go)

## `s.Discard()`

For runes like whitespace, it useful to simply discard them without creating
tokens. Use the `s.Discard()` method to reset the current token buffer, discard
the current rune, and advance the scanner to the next rune. The example below
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

## Types

The next example is going to scan the entire stream and create tokens for each
integer and word that is found. A word is going to be any sequence of letters
and letters only. Integers and words must be separated by whitespace. For this
scanner, each token type will be placed in a separate function and each token
will be marked with a type. A type can be any string value but there are some
constants predefined in the scan package. The functions return a boolean value
to indicate if there was a match for that token.

The function for the integer is as follows:

```go
    func scanInt(s *scan.Scanner) bool {
        if !s.Is(scan.Digit09) {
            return false
        }
        s.Type = scan.IntType
        scan.While(s, scan.Digit09, s.Keep)
        return true
    }
```
[Example 6](doc/example_6_test.go)

The function for the word is similar and is as follows:

```go
    func scanWord(s *scan.Scanner) bool {
        if !s.Is(scan.Letter) {
            return false
        }
        s.Type = scan.WordType
        scan.While(s, scan.Letter, s.Keep)
        return true
    }
```
[Example 6](doc/example_6_test.go)

Since whitespace is now significant, it will be returned as a token instead
of being discarded:

```go
    func scanWhitespace(s *scan.Scanner) bool {
        if !s.Is(scan.Whitespace) {
            return false
        }
        s.Type = scan.SpaceType
        scan.While(s, scan.Whitespace, s.Keep)
        return true
    }
```
[Example 6](doc/example_6_test.go)

Now create a main loop that checks each function for a match and collects all
tokens seen along the way:

```go
    scanFuncs := []func(*scan.Scanner) bool {
        scanWhitespace,
        scanWord,
        scanInt,
    }

    var toks []scan.Token
    s := scan.NewScannerFromString("example6", "abc 123 !@# def456")
    for s.HasMore() {
        match := false
        for _, fn := range scanFuncs {
            if match = fn(s); match {
                toks = append(toks, s.Emit())
                break
            }
        }
        if !match {
            scan.Until(s, scan.Whitespace, s.Keep)
            s.Illegal("unexpected: %v", scan.Quote(s.Val.String()))
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
    // example6:1:12: error: unexpected: "!@#"
    // example6:1:12  space    "  {!ch:\t}"  "  {!ch:\t}"
    // example6:1:15  word     "def"         "def"
    // example6:1:18  int      "456"         "456"
```
[Example 6](doc/example_6_test.go)

For each token that is generated, we can check the `tok.Type` field to see if
it is an int, word, space, or an illegal token. The `tok.Pos` field contains
the name of the file (or item) being scanned and the line and column number of
where the token starts.

Note that the role of the scanner is to simply emit tokens of the various
types seen in the input stream. The "def456" is illegal by this example's
definition because a space must separate the word and the integer. That
check should be handled by the parser that is processing the tokens, not by
the scanner itself.

There are a few new functions or methods introduced in the example:

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
        digit scan.Class
    }

    func NewIntRule(digit scan.Class) IntRule {
        return IntRule{digit: digit}
    }

    func (r IntRule) Eval(s *scan.Scanner) bool {
        if !s.Is(r.digit) {
            return false
        }
        s.Type = scan.IntType
        scan.While(s, r.digit, s.Keep)
        return true
    }
```
[Example 7](doc/example_7_test.go)

The rules for word and whitespace looks similar and can be found in the links
to the example.

Now a function is needed to handle the case where a token does not match.
This also has to be configurable to scan until whitespace is found. This
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
`scan.RuleSet`. Then, a `scan.Runner` can be created with the scanner and rule
set to iterate through the tokens using `runner.HasMore` and `runner.Next`. Or,
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
[Example 7](doc/example_7_test.go)

## Preprocessing

Sometimes it is convenient to perform some light preprocessing of the token
value before passing it off to the parser. A good example of this is allowing
digit separators in an integer value such as 1,000. The parser is most likely
going to pass this responsibility over to `strconv.ParseInt()` and the
comma gets in the way. Another use case is to convert runes to another rune.
For example, converting all whitespace runes to a space value or changing
newlines to semicolons. The following can be used for these cases:

#### s.Skip()

Adds the current character to the token's literal but not to the token's value.
The scanner then advances to the next character.

#### s.Value

This is a string builder containing the token value seen so far. If newlines
need to be converted to semicolons, this can be done like in the following:

```go
    if s.This == '\n' {
        s.Skip()
        s.Value.WriteRune(';')
    }
```

Now the integer rule can be updated to include a digit separator. The digit
class is mandatory but the digit separator is optional so only the digit
class is found in the constructor, but the digit separator can be included
with a chaining function:

```go
type IntRule2 struct {
    digit    scan.Class
    digitSep scan.Class
}

func NewIntRule2(digit scan.Class) IntRule2 {
    return IntRule2{digit: digit}
}

func (r IntRule2) WithDigitSep(digitSep scan.Class) IntRule2 {
    r.digitSep = digitSep
    return r
}

func (r IntRule2) Eval(s *scan.Scanner) bool {
    if !s.Is(r.digit) {
        return false
    }
    s.Keep()
    for s.HasMore() {
        if s.Is(r.digit) {
            s.Keep()
        } else if s.Is(r.digitSep) {
            s.Skip()
        } else {
            break
        }
    }
    return true
}
```
[Example 8](doc/example_8_test.go)

Now update the rule set with the new rule and modify the input to include some
digit separators:

```go
    rules := scan.NewRuleSet(
        NewSpaceRule(scan.Whitespace),
        NewWordRule(scan.Letter),
        NewIntRule2(scan.Digit).
            WithDigitSep(scan.Rune(',')),
    ).WithNoMatchFunc(UnexpectedUntil(scan.Whitespace))

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
[Example 8](doc/example_8_test.go)

## Predefined Rules

The scan package already has rules defined for reading the token types used
so far. The example can now be updated to use those rules instead:

```go
    rules := scan.NewRuleSet(
        scan.NewSpaceRule(scan.Whitespace).WithKeep(true),
        scan.NewWordRule(scan.Letter),
        scan.Int.WithDigitSep(scan.Rune(',')),
    ).WithNoMatchFunc(scan.UnexpectedUntil(scan.Whitespace))

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
[Example 9](doc/example_9_test.go)

The predefined rules are as follows:

* `scan.CharEncRule`: For converting character encodings such as `\n` to their
actual values like `U+000A`
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
* `scan.WordRule`: Tokens that match against a single character class.

Some rules have predefined configurations as well. For example `scan.Hex0x`
that uses a `scan.NumRule` configured for hexadecimal digits that are
prefixed with `0x`. See the API documentation for more information.

## Lookahead and lookbehind

Most scanning operations can be done by looking at the current character,
which is stored in `s.This` and the next character in the stream that is
stored in `s.Next`. In some cases, it can be easier to look further ahead
in the stream or to even speculate and advance the scanner and cancel out
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

#### `s.PeekIs()`

The same as `s.Is()` but takes a peek position as well.

#### `s.Undo()`

When speculating with the scanner, calling this method will place the token
literal seen so far back onto the input stream and then clear the token
buffer.

## Full Examples

There are two full examples provided with this package. The first is a
[JSON scanner](scanjson/scanjson.go). The scanner for JSON is quite simple but
JSON does have some quirks:

* Number can start with a negative sign, but not a positive sign
* Numbers cannot have leading zeros. "303" is okay, but "0303" is not.
* Real numbers cannot have empty parts. "0.3" is valid and "3.0" is valid, but ".3" and "3." are not.

There is also a [Go scanner](scango/scango.go). This one is a bit more
complicated but be constructed mostly from predefined rules. The exceptions
are for a rule handling imaginary numbers and a post token processor for
automatic semicolon insertion.

## Status

This package is still a work in progress and is subject to change. If you
find this useful, please drop a line at zc at blackchip dot org.

## License

MIT