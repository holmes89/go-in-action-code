# Chapter 3: Primitive Types And Operators

This directory contains a comprehensive reference implementation covering all primitive types in Go.

## Main File

**main.go** - Complete demonstration of all primitive types and operators with runnable examples

## What's Inside

The main.go file covers:

1. **Integer Types** - Signed and unsigned integers, type selection
2. **Integer Literals** - Decimal, hex, octal, binary notation
3. **Unsigned Wraparound** - Common pitfall and how to avoid it
4. **Floating-Point Types** - float32 and float64
5. **Special Float Values** - Infinity and NaN handling
6. **Complex Numbers** - complex64 and complex128
7. **Math Operators** - Arithmetic and bitwise operations
8. **Boolean Type** - true/false and conditional logic
9. **Struct Types** - Grouping related data
10. **Pointer Types** - Memory addresses and references
11. **String Types** - UTF-8 text handling
12. **Runes** - Unicode character manipulation

## Running the Program

Run the complete demonstration:
```bash
go run main.go
```

Each section is clearly marked and demonstrates the concepts from Chapter 3.

## Key Concepts

### Type Selection Guidelines

**Integers:**
- Use `int` for most cases (default, most compatible)
- Use `int64` for timestamps, file sizes, large ranges
- Avoid unsigned integers for "positive only" validation
- Use sized types (int8, int16, etc.) only when necessary

**Floats:**
- Use `float64` by default (better precision)
- Only use `float32` for proven performance needs

**Strings:**
- Remember: `len()` returns bytes, not characters
- Use `range` loops for Unicode iteration
- Convert to `[]rune` for character indexing

### Common Pitfalls

1. **Unsigned Wraparound**: `var u uint = 0; u = u - 1` → huge number!
2. **String Length**: `len("地鼠")` returns 6 (bytes), not 2 (characters)
3. **Float Special Values**: Check for Inf and NaN with math package
4. **Integer Division**: `7 / 3` = 2, not 2.33 (use floats for decimals)

## Type Quick Reference

```go
// Integers
var i int              // Most common, architecture-dependent
var b byte             // Alias for uint8
var r rune             // Alias for int32, for Unicode

// Floats
var f float64          // Default for decimals

// Complex
c := 1 + 2i           // complex128

// Boolean
var flag bool          // true or false

// String
s := "Hello, 世界"     // UTF-8 text

// Struct
type Person struct {
    name string
    age  int
}

// Pointer
ptr := &value          // Get address
val := *ptr           // Dereference
```

## Mathematical Operations

```go
// Arithmetic
a + b, a - b, a * b, a / b, a % b

// Bitwise (integers only)
x & y   // AND
x | y   // OR
x ^ y   // XOR
x &^ y  // bit clear
x << 2  // left shift
x >> 2  // right shift
```

## String and Rune Operations

```go
// String length (bytes)
len("hello")           // 5

// Unicode length (characters)
len([]rune("地鼠"))    // 2

// Iterate characters
for i, r := range "Hello, 世界" {
    // i is byte position
    // r is rune (character)
}

// Character inspection
unicode.IsLetter(r)
unicode.IsUpper(r)
unicode.ToLower(r)
utf8.RuneLen(r)
```

## Best Practices

✅ **DO:**
- Use `int` as your default integer type
- Use `float64` as your default floating-point type
- Check for negative values explicitly (don't rely on unsigned)
- Use `range` loops for Unicode strings
- Use boolean helper functions for complex conditions
- Use named struct types for reusability

❌ **DON'T:**
- Use unsigned integers to mean "positive only"
- Use `len()` to count Unicode characters
- Mix arithmetic operations between different types
- Forget to check for Inf/NaN in float calculations
- Index strings directly for Unicode text

## Testing Your Understanding

Try modifying the code to:
1. Add a function that converts Celsius to Fahrenheit (floats)
2. Create a struct representing a Rectangle with methods
3. Write a function that reverses a Unicode string properly
4. Demonstrate bitwise operations for flag management
5. Handle special float values in calculations
