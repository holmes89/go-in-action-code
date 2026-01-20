package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Chapter 2: Diving Into Go - Word Count Program
// This file demonstrates the progression of building a word counting program
// from basic iteration to reading files.

// ==============================================================================
// ITERATION 1: Basic Space Counting (Commented Out)
// ==============================================================================
// This was our first approach - counting spaces to estimate words.
// Problem: Doesn't handle multiple spaces or spaces at edges well.
/*
func countSpaces() {
	// Use hard-coded text to start.
	text := "let's count some words!"

	var numSpaces int

	// Look at the text one byte at a time.
	for i := 0; i < len(text); i++ {
		if text[i] == ' ' {
			numSpaces++
		}
	}

	fmt.Println("Found", numSpaces+1, "words")
}
*/

// ==============================================================================
// ITERATION 2: Using strings.Fields (Commented Out)
// ==============================================================================
// Better approach: strings.Fields properly handles whitespace splitting.
// This is more robust than counting spaces.
/*
func countWithFields() {
	text := "let's count some words!"
	words := strings.Fields(text)
	fmt.Println("Found", len(words), "words")
}
*/

// ==============================================================================
// ITERATION 3: Reading from a File (Current Implementation)
// ==============================================================================
// Final version that reads from a file specified as a command-line argument.

func main() {
	// Check if a filename was provided as a command-line argument.
	// os.Args[0] is the program name, os.Args[1] is the first argument.
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		fmt.Println("\nExample: go run main.go sample.txt")
		os.Exit(1)
	}

	// Get the filename from the command-line arguments.
	filename := os.Args[1]

	// Read the entire file into memory.
	// os.ReadFile returns two values: the file contents and an error.
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		// If there was an error reading the file, log it and exit.
		log.Fatalf("Error reading file '%s': %v", filename, err)
	}

	// Convert the byte slice to a string.
	// Files are read as []byte, but strings.Fields expects a string.
	text := string(fileContents)

	// Split the text into words using strings.Fields.
	// This handles all whitespace characters and multiple spaces correctly.
	words := strings.Fields(text)

	// Print the result.
	fmt.Printf("File: %s\n", filename)
	fmt.Printf("Found %d words\n", len(words))

	// BONUS: Show some statistics
	if len(words) > 0 {
		fmt.Printf("First word: %s\n", words[0])
		fmt.Printf("Last word: %s\n", words[len(words)-1])
	}
}

// ==============================================================================
// KEY CONCEPTS DEMONSTRATED IN THIS CHAPTER
// ==============================================================================
//
// 1. PACKAGE DECLARATION
//    - "package main" creates an executable program
//    - Must have a main() function as the entry point
//
// 2. IMPORTS
//    - Import standard library packages with import "packagename"
//    - Group multiple imports with parentheses: import ( ... )
//    - fmt: formatted I/O
//    - strings: string manipulation
//    - os: operating system interaction
//    - log: logging and error handling
//
// 3. VARIABLES
//    - var name type: explicit type declaration
//    - name := value: short declaration with type inference
//    - Zero values: 0 for numbers, "" for strings, nil for pointers
//
// 4. FOR LOOPS
//    - Three-part: for init; condition; post { }
//    - Range: for i, v := range collection { }
//    - While-style: for condition { }
//
// 5. FUNCTIONS
//    - Multiple return values: func name() (type1, type2)
//    - Error handling pattern: check if err != nil
//    - Package.Function syntax for imported packages
//
// 6. ERROR HANDLING
//    - Check errors immediately after they're returned
//    - Handle or propagate errors appropriately
//    - log.Fatal prints error and exits with status 1
//
// 7. SLICES
//    - []type is a slice (dynamic array)
//    - len() returns the number of elements
//    - Indexed with [i], starting at 0
//    - strings.Fields returns []string
//
// 8. COMMAND-LINE ARGUMENTS
//    - os.Args is a []string containing arguments
//    - os.Args[0] is the program name
//    - os.Args[1:] are the actual arguments
//
// ==============================================================================
// USAGE EXAMPLES
// ==============================================================================
//
// Build and run:
//   $ go build -o wordcount main.go
//   $ ./wordcount sample.txt
//
// Or run directly:
//   $ go run main.go sample.txt
//
// Format code:
//   $ gofmt -w main.go
//
// View documentation:
//   $ go doc strings.Fields
//   $ go doc os.ReadFile
//
// ==============================================================================
// EXERCISE IDEAS
// ==============================================================================
//
// 1. Add a -v (verbose) flag that shows each word found
// 2. Count lines in addition to words (hint: strings.Split with "\n")
// 3. Count bytes and compare to len(fileContents)
// 4. Handle multiple files (iterate over os.Args[1:])
// 5. Add a -h flag that prints usage information
// 6. Count unique words using a map[string]int
//
// ==============================================================================
