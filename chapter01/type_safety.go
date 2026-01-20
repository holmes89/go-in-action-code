package main

// This example demonstrates Go's type safety.
// This code will NOT compile - it's meant to show the compile-time error.
// Uncomment the last line to see the error.

func main() {
	var number int
	var str string

	number = 1
	str = "one"
	// number = str
	// Error: cannot use str (variable of type string) as type int in assignment

	_ = number
	_ = str
}
