package main

import (
	"fmt"
	"math"
	"unicode"
	"unicode/utf8"
)

// Chapter 3: Primitive Types And Operators
// This file demonstrates all the primitive types in Go with practical examples.
// Each section is runnable and commented for easy reference.

func main() {
	fmt.Println("=== CHAPTER 3: PRIMITIVE TYPES AND OPERATORS ===\n")

	// Run each demonstration
	demonstrateIntegers()
	demonstrateIntegerLiterals()
	demonstrateUnsignedWraparound()
	demonstrateFloats()
	demonstrateSpecialFloats()
	demonstrateComplexNumbers()
	demonstrateMathOperators()
	demonstrateBooleans()
	demonstrateStructs()
	demonstratePointers()
	demonstrateStrings()
	demonstrateRunes()
}

// ==============================================================================
// INTEGER TYPES
// ==============================================================================
// Go has signed (int) and unsigned (uint) integers in various sizes.
// int is the most common and should be your default choice.

func demonstrateIntegers() {
	fmt.Println("--- INTEGER TYPES ---")

	// Declare a variable with a type of int (default value is 0).
	var myInt int
	fmt.Printf("myInt: %d (type: %T)\n", myInt, myInt)

	// Declare variable with type int64
	var largeInt int64
	fmt.Printf("largeInt: %d (type: %T)\n", largeInt, largeInt)

	// Use short declaration to create variable with type int and a value of 3.
	i := 3
	fmt.Printf("i: %d (type: %T)\n", i, i)

	// Use a type conversion for other integer types.
	u := uint64(4)
	fmt.Printf("u: %d (type: %T)\n", u, u)

	// Common integer types:
	// int8, int16, int32, int64 (signed)
	// uint8, uint16, uint32, uint64 (unsigned)
	// byte (alias for uint8)
	// rune (alias for int32, used for Unicode)

	fmt.Println()
}

// ==============================================================================
// INTEGER LITERALS - DIFFERENT NOTATIONS
// ==============================================================================
// Integer literals can be written in decimal, hex, octal, or binary.

func demonstrateIntegerLiterals() {
	fmt.Println("--- INTEGER LITERALS ---")

	// These are different ways of writing the same value (1000).
	decInt := 1000              // Base 10 Notation (no prefix)
	hexInt := 0x3E8             // Hexadecimal Notation (Prefix: "0x")
	octInt := 01750             // Octal Notation (Prefix: "0")
	octIntAlt := 0o1750         // Alternate Octal Notation (Prefix: "0o")
	binInt := 0b1111101000      // Binary Notation (Prefix: "0b")

	// Optional underscores can be used to separate digits for readability.
	withSep := 1_000
	hexWithSep := 0x3_E_8

	// Notation is purely cosmetic, and does not affect value:
	fmt.Println("All represent 1000:", decInt, hexInt, octInt, octIntAlt, binInt, withSep, hexWithSep)

	// Negative integer literals:
	negativeInt := -10
	negativeHexInt := -0xa
	fmt.Printf("Negative values: %d, %d\n", negativeInt, negativeHexInt)

	fmt.Println()
}

// ==============================================================================
// UNSIGNED INTEGER WRAPAROUND
// ==============================================================================
// WARNING: Unsigned integers wrap around when they go below 0!

func demonstrateUnsignedWraparound() {
	fmt.Println("--- UNSIGNED INTEGER WRAPAROUND ---")

	var u uint64 // 0
	fmt.Println("Initial value:", u)

	u = u - 1
	fmt.Println("After subtracting 1:", u) // Wraps to max value!

	// Similar behavior with uint8
	var u8 uint8 = 0
	u8 = u8 - 1
	fmt.Println("uint8 wraparound:", u8) // 255

	// BEST PRACTICE: Avoid using unsigned integers to mean "only positive"
	// Instead, use signed integers and check for negative values:
	exampleValue := -5
	if err := acceptOnlyPositive(exampleValue); err != nil {
		fmt.Println("Error caught:", err)
	}

	fmt.Println()
}

// acceptOnlyPositive demonstrates proper handling of negative values
// instead of using unsigned integers.
func acceptOnlyPositive(i int) error {
	if i < 0 {
		return fmt.Errorf("value must be positive, got: %d", i)
	}
	fmt.Println("Valid positive value:", i)
	return nil
}

// ==============================================================================
// FLOATING-POINT TYPES
// ==============================================================================
// float32 and float64 represent decimal numbers.
// Use float64 by default (better precision, good hardware support).

func demonstrateFloats() {
	fmt.Println("--- FLOATING-POINT TYPES ---")

	// Declare with default value (0.0)
	var doubleFloat float64
	fmt.Printf("doubleFloat: %f (type: %T)\n", doubleFloat, doubleFloat)

	var singleFloat float32
	fmt.Printf("singleFloat: %f (type: %T)\n", singleFloat, singleFloat)

	// Short declaration (defaults to float64)
	f := 12.1
	fmt.Printf("f: %f (type: %T)\n", f, f)

	// Explicit float32
	g := float32(12.1)
	fmt.Printf("g: %f (type: %T)\n", g, g)

	// Different ways to write the same float value
	floatVal := 12.         // fractional part optional
	floatVal = 12e0        // scientific notation
	floatVal = .12e+2      // integer part optional
	floatVal = 1_2.        // underscores for readability
	floatVal = float64(12) // convert from int
	fmt.Printf("All represent 12.0: %f\n", floatVal)

	// Negative floats
	negativeFloat := -12.5
	fmt.Printf("Negative float: %f\n", negativeFloat)

	fmt.Println()
}

// ==============================================================================
// SPECIAL FLOATING-POINT VALUES
// ==============================================================================
// Floats can have special values: +Inf, -Inf, and NaN (Not a Number)

func demonstrateSpecialFloats() {
	fmt.Println("--- SPECIAL FLOATING-POINT VALUES ---")

	// Positive and negative infinity
	posInf := math.Inf(1)
	negInf := math.Inf(-1)
	fmt.Println("Positive infinity:", posInf)
	fmt.Println("Negative infinity:", negInf)

	// Not a Number (NaN)
	nan := math.NaN()
	fmt.Println("NaN:", nan)

	// Division by zero creates infinity (not panic like integers!)
	var x float64 = 1.0
	var y float64 = 0.0
	result := x / y
	fmt.Println("1.0 / 0.0 =", result)

	// Invalid operations create NaN
	invalid := math.Sqrt(-1)
	fmt.Println("sqrt(-1) =", invalid)

	// Checking for special values
	fmt.Println("Is result Inf?", math.IsInf(result, 0))
	fmt.Println("Is invalid NaN?", math.IsNaN(invalid))

	// Large number overflow
	f := 2.0
	huge := math.Pow(f, 10_000)
	fmt.Println("2^10000 =", huge)

	fmt.Println()
}

// ==============================================================================
// COMPLEX NUMBER TYPES
// ==============================================================================
// complex64 and complex128 represent complex numbers with real and imaginary parts.

func demonstrateComplexNumbers() {
	fmt.Println("--- COMPLEX NUMBER TYPES ---")

	// Create using complex() function
	cmplx1 := complex(1.1, 2.2)                   // complex128
	cmplx2 := complex(float32(1.1), float32(2.2)) // complex64

	fmt.Printf("complex128: %v (type: %T)\n", cmplx1, cmplx1)
	fmt.Printf("complex64: %v (type: %T)\n", cmplx2, cmplx2)

	// Create using complex expressions
	cmplx3 := 1.1 + 2.2i  // equivalent to complex(1.1, 2.2)
	cmplx4 := 1 + 2i      // equivalent to complex(1, 2)
	cmplx5 := -3i         // equivalent to complex(0, -3)
	cmplx6 := -3.3 - 4.4i // equivalent to complex(-3.3, -4.4)

	fmt.Println("Using expressions:")
	fmt.Println("  1.1 + 2.2i =", cmplx3)
	fmt.Println("  1 + 2i =", cmplx4)
	fmt.Println("  -3i =", cmplx5)
	fmt.Println("  -3.3 - 4.4i =", cmplx6)

	// Extract real and imaginary parts
	fmt.Printf("Real part of cmplx1: %f\n", real(cmplx1))
	fmt.Printf("Imaginary part of cmplx1: %f\n", imag(cmplx1))

	fmt.Println()
}

// ==============================================================================
// MATHEMATICAL OPERATORS
// ==============================================================================
// Go supports standard arithmetic and bitwise operations.

func demonstrateMathOperators() {
	fmt.Println("--- MATHEMATICAL OPERATORS ---")

	// Arithmetic operators
	a := 7
	b := 3

	fmt.Println("Arithmetic Operations:")
	fmt.Printf("  %d + %d = %d\n", a, b, a+b) // 10
	fmt.Printf("  %d - %d = %d\n", a, b, a-b) // 4
	fmt.Printf("  %d * %d = %d\n", a, b, a*b) // 21
	fmt.Printf("  %d / %d = %d\n", a, b, a/b) // 2 (integer division)
	fmt.Printf("  %d %% %d = %d\n", a, b, a%b) // 1 (remainder/modulo)

	// Float division gives decimal result
	fmt.Printf("  %.2f / %.2f = %.2f\n", 7.0, 3.0, 7.0/3.0) // 2.33

	// Bitwise operators (integers only)
	fmt.Println("\nBitwise Operations:")
	x := 12  // 1100 in binary
	y := 10  // 1010 in binary
	fmt.Printf("  %d & %d = %d (AND)\n", x, y, x&y)       // 8 (1000)
	fmt.Printf("  %d | %d = %d (OR)\n", x, y, x|y)        // 14 (1110)
	fmt.Printf("  %d ^ %d = %d (XOR)\n", x, y, x^y)       // 6 (0110)
	fmt.Printf("  %d &^ %d = %d (bit clear)\n", x, y, x&^y) // 4 (0100)
	fmt.Printf("  %d << 2 = %d (left shift)\n", x, x<<2)  // 48
	fmt.Printf("  %d >> 2 = %d (right shift)\n", x, x>>2) // 3

	fmt.Println()
}

// ==============================================================================
// BOOLEAN TYPE
// ==============================================================================
// The bool type represents true/false values.

func demonstrateBooleans() {
	fmt.Println("--- BOOLEAN TYPE ---")

	// Default value is false
	var boolVar bool
	fmt.Println("Default bool value:", boolVar)

	// Boolean expressions
	x := 0
	y := 1
	fmt.Printf("%d < %d: %v\n", x, y, x < y)
	fmt.Printf("%d == %d: %v\n", x, y, x == y)
	fmt.Printf("%d != %d: %v\n", x, y, x != y)

	// Boolean operations
	truth := true
	falsehood := false
	fmt.Println("\nBoolean Operations:")
	fmt.Printf("  true && false: %v\n", truth && falsehood)
	fmt.Printf("  true || false: %v\n", truth || falsehood)
	fmt.Printf("  !true: %v\n", !truth)

	// Using booleans in conditionals
	if isEven(4) {
		fmt.Println("4 is even")
	}

	// Boolean helper functions make code more readable
	numbers := []int{1, 2, 3, 4, 5, 6}
	fmt.Print("Even numbers: ")
	for _, n := range numbers {
		if isEven(n) {
			fmt.Print(n, " ")
		}
	}
	fmt.Println()

	fmt.Println()
}

// isEven is a boolean helper function that makes conditionals more readable.
func isEven(i int) bool {
	return i%2 == 0
}

// ==============================================================================
// STRUCT TYPES
// ==============================================================================
// Structs group related values together into a single unit.

func demonstrateStructs() {
	fmt.Println("--- STRUCT TYPES ---")

	// Anonymous struct (inline definition)
	var person1 struct {
		name string
		age  int
	}
	person1.name = "Andy"
	person1.age = 42
	fmt.Printf("person1: %+v\n", person1)

	// Named struct type (preferred for reusability)
	type Person struct {
		name string
		age  int
	}

	// Create with field names (recommended)
	andy := Person{
		name: "Andy",
		age:  42,
	}

	// Create with positional values (not recommended)
	james := Person{"James", 26}

	// Zero value struct
	var empty Person
	fmt.Printf("andy: %+v\n", andy)
	fmt.Printf("james: %+v\n", james)
	fmt.Printf("empty: %+v\n", empty)

	// Access fields with dot notation
	fmt.Printf("%s is %d years old\n", andy.name, andy.age)

	// Modify fields
	andy.age = 43
	fmt.Printf("After birthday: %+v\n", andy)

	// Anonymous struct for one-off use
	point := struct {
		X, Y int
	}{10, 20}
	fmt.Printf("Point: %+v\n", point)

	fmt.Println()
}

// ==============================================================================
// POINTER TYPES
// ==============================================================================
// Pointers store memory addresses, allowing sharing without copying.

func demonstratePointers() {
	fmt.Println("--- POINTER TYPES ---")

	// Declare a pointer variable (defaults to nil)
	var intPtr *int
	fmt.Printf("Nil pointer: %v\n", intPtr)

	// Get the address of a variable with & (address operator)
	intValue := 42
	intPtr = &intValue
	fmt.Printf("intValue: %d, address: %p\n", intValue, intPtr)

	// Dereference with * to get/set the value
	fmt.Printf("Value at pointer: %d\n", *intPtr)

	*intPtr = 100
	fmt.Printf("After modifying via pointer: intValue=%d, *intPtr=%d\n", intValue, *intPtr)

	// Create pointer with new
	newIntPtr := new(int)
	*newIntPtr = 42
	fmt.Printf("Pointer from new: %d\n", *newIntPtr)

	// Pointers to structs
	type Point struct {
		X, Y int
	}

	point := Point{X: 1, Y: 2}
	ptr := &point

	// Can access fields directly (automatic dereferencing)
	ptr.X = 10
	ptr.Y = 20
	fmt.Printf("Point after modification via pointer: %+v\n", point)

	// Create pointer with struct literal
	newPoint := &Point{X: 5, Y: 10}
	fmt.Printf("Pointer to new point: %+v\n", newPoint)

	fmt.Println()
}

// ==============================================================================
// STRING TYPES
// ==============================================================================
// Strings are read-only arrays of bytes containing UTF-8 text.

func demonstrateStrings() {
	fmt.Println("--- STRING TYPES ---")

	// Interpreted strings (double quotes) support escape codes
	basicStr := "Hello, Gophers!"
	fmt.Println("Basic string:", basicStr)

	// Unicode support
	unicodeStr := "你好，地鼠!"
	fmt.Println("Unicode string:", unicodeStr)

	// Escape sequences
	strWithEscapes := "Hello\nWorld\tTabbed"
	fmt.Println("With escapes:")
	fmt.Println(strWithEscapes)

	// Raw strings (backticks) - no escape processing
	rawString := `Hello\n\u5730\u9F20`
	fmt.Println("Raw string:", rawString)

	// Raw strings can span multiple lines
	rawMultiline := `Line 1
Line 2
Line 3`
	fmt.Println("Multiline raw string:")
	fmt.Println(rawMultiline)

	// String operations
	fmt.Println("\nString Operations:")
	fmt.Printf("  \"hello\" == \"hello\": %v\n", "hello" == "hello")
	fmt.Printf("  \"abc\" < \"xyz\": %v\n", "abc" < "xyz")

	// Concatenation
	str1 := "Hello, "
	str2 := "World!"
	combined := str1 + str2
	fmt.Println("  Concatenation:", combined)

	// Length in bytes (not characters!)
	asciiStr := "hello"
	fmt.Printf("  len(\"%s\"): %d bytes\n", asciiStr, len(asciiStr))

	// Unicode strings have different byte lengths
	unicodeStr2 := "地鼠"
	fmt.Printf("  len(\"%s\"): %d bytes (2 characters, but 6 bytes!)\n", unicodeStr2, len(unicodeStr2))

	// Byte vs character visualization
	fmt.Printf("  Bytes of \"%s\": %x\n", unicodeStr2, unicodeStr2)

	fmt.Println()
}

// ==============================================================================
// RUNES AND UNICODE
// ==============================================================================
// Runes (alias for int32) represent Unicode code points.

func demonstrateRunes() {
	fmt.Println("--- RUNES AND UNICODE ---")

	// A rune is an int32 representing a Unicode code point
	var r rune = 'A'
	fmt.Printf("Rune 'A': %c (value: %d, type: %T)\n", r, r, r)

	// Unicode characters
	emoji := '😀'
	fmt.Printf("Emoji rune: %c (value: %d, hex: %X)\n", emoji, emoji, emoji)

	// Character size varies
	fmt.Println("\nCharacter sizes in bytes:")
	fmt.Printf("  'A': %d byte(s)\n", utf8.RuneLen('A'))
	fmt.Printf("  'À': %d byte(s)\n", utf8.RuneLen('À'))
	fmt.Printf("  '地': %d byte(s)\n", utf8.RuneLen('地'))
	fmt.Printf("  '😀': %d byte(s)\n", utf8.RuneLen('😀'))

	// Iterating strings with range gives you runes
	text := "Hello, 世界"
	fmt.Printf("\nIterating \"%s\" with range:\n", text)
	for i, r := range text {
		fmt.Printf("  Position %d: %c (U+%04X)\n", i, r, r)
	}

	// Convert string to []rune for indexed access
	runes := []rune(text)
	fmt.Printf("\nAs rune slice: %d characters\n", len(runes))
	fmt.Printf("  First character: %c\n", runes[0])
	fmt.Printf("  Last character: %c\n", runes[len(runes)-1])

	// Unicode package functions
	fmt.Println("\nUnicode character inspection:")
	chars := "Àa1😀"
	for _, r := range chars {
		fmt.Printf("  %c:", r)
		if unicode.IsLetter(r) {
			fmt.Print(" letter")
			if unicode.IsUpper(r) {
				fmt.Print(" uppercase")
			}
		}
		if unicode.IsDigit(r) {
			fmt.Print(" digit")
		}
		if unicode.IsSymbol(r) {
			fmt.Print(" symbol")
		}
		fmt.Println()
	}

	fmt.Println()
}

// ==============================================================================
// KEY CONCEPTS SUMMARY
// ==============================================================================
//
// INTEGERS
// - Use 'int' for most cases (signed, architecture-dependent)
// - Specific sizes: int8, int16, int32, int64
// - Unsigned: uint, uint8, uint16, uint32, uint64
// - Aliases: byte (uint8), rune (int32)
// - Avoid unsigned integers for "positive only" - use int with error checking
//
// FLOATING-POINT
// - Use float64 by default (better precision)
// - float32 only for specific performance needs
// - Can have special values: +Inf, -Inf, NaN
// - Division by zero doesn't panic (creates infinity)
//
// COMPLEX NUMBERS
// - complex64 (float32 parts), complex128 (float64 parts)
// - Created with complex(real, imag) or expressions (1+2i)
// - Extract parts with real() and imag()
//
// BOOLEANS
// - Two values: true and false
// - Default value is false
// - Use boolean functions for readable conditionals
//
// STRUCTS
// - Group related values together
// - Access fields with dot notation
// - Use named types for reusability
//
// POINTERS
// - Store memory addresses
// - & gets address, * dereferences
// - Allows sharing without copying
// - Default value is nil
//
// STRINGS
// - Read-only byte arrays with UTF-8 text
// - len() returns bytes, not characters
// - Use range loops to iterate Unicode characters
// - Interpreted strings ("") support escapes
// - Raw strings (``) take everything literally
//
// RUNES
// - Represent Unicode code points (int32 alias)
// - Use for character manipulation
// - Convert string to []rune for indexed character access
//
// ==============================================================================
