package main

import (
	"fmt"
)

// Chapter 4: Collection Types
// This file demonstrates arrays, slices, and maps with practical examples.
// Each section is runnable and commented for easy reference.

func main() {
	fmt.Println("=== CHAPTER 4: COLLECTION TYPES ===\n")

	// Run each demonstration
	demonstrateArrays()
	demonstrateArrayDeclarations()
	demonstrateArrayIteration()
	demonstrateArrayValues()
	demonstrateMultidimensionalArrays()
	demonstrateSlices()
	demonstrateSliceDeclarations()
	demonstrateSliceAppend()
	demonstrateSliceExpressions()
	demonstrateSliceCopy()
	demonstrateSlicePanics()
	demonstrateMaps()
	demonstrateMapOperations()
	demonstrateMapIteration()
}

// ==============================================================================
// ARRAYS - FIXED-LENGTH COLLECTIONS
// ==============================================================================
// Arrays are fixed-length collections containing contiguous blocks of elements.
// The length is part of the type, so [5]int is different from [4]int.

func demonstrateArrays() {
	fmt.Println("--- ARRAYS BASICS ---")

	// Declare an array with var (all values default to zero)
	var intArray [5]int
	fmt.Printf("Declared array: %v\n", intArray) // [0 0 0 0 0]

	// Set values using index notation
	intArray[0] = 10
	intArray[1] = 20
	intArray[2] = 30
	intArray[3] = 40
	intArray[4] = 50
	fmt.Printf("After assignment: %v\n", intArray) // [10 20 30 40 50]

	// Access values by index
	firstValue := intArray[0]
	lastValue := intArray[4]
	fmt.Printf("First: %d, Last: %d\n", firstValue, lastValue)

	// Get array length with len()
	fmt.Printf("Array length: %d\n", len(intArray))

	// IMPORTANT: Array length is part of the type
	// These are DIFFERENT types and cannot be assigned to each other:
	var array4 [4]int
	var array5 [5]int
	// array5 = array4 // ERROR: cannot use array4 (type [4]int) as type [5]int

	fmt.Printf("array4 type: %T\n", array4) // [4]int
	fmt.Printf("array5 type: %T\n", array5) // [5]int

	fmt.Println()
}

// ==============================================================================
// ARRAY DECLARATIONS - DIFFERENT FORMS
// ==============================================================================
// Arrays can be declared in several ways with different initialization options.

func demonstrateArrayDeclarations() {
	fmt.Println("--- ARRAY DECLARATIONS ---")

	// 1. Array literal with all values
	fullArray := [5]int{10, 20, 30, 40, 50}
	fmt.Printf("Full array: %v\n", fullArray)

	// 2. Array literal with partial values (rest are zero)
	partialArray := [5]int{10, 20, 30}
	fmt.Printf("Partial array: %v\n", partialArray) // [10 20 30 0 0]

	// 3. Array literal with sparse values (specify indices)
	sparseArray := [5]int{1: 20, 3: 40}
	fmt.Printf("Sparse array: %v\n", sparseArray) // [0 20 0 40 0]

	// 4. Array literal starting from specific index
	lastThree := [5]int{2: 30, 40, 50}
	fmt.Printf("Last three: %v\n", lastThree) // [0 0 30 40 50]

	// 5. Automatic length with ... (compiler counts elements)
	autoArray := [...]int{5, 4, 3, 2, 1}
	fmt.Printf("Auto-sized array: %v (len: %d)\n", autoArray, len(autoArray))

	// 6. Automatic length with sparse values (size based on highest index)
	autoSparse := [...]int{0: 5, 3: 2}
	fmt.Printf("Auto sparse: %v (len: %d)\n", autoSparse, len(autoSparse)) // [5 0 0 2]

	// Arrays of different types
	stringArray := [3]string{"hello", "world", "!"}
	floatArray := [3]float64{1.1, 2.2, 3.3}
	boolArray := [3]bool{true, false, true}

	fmt.Printf("Strings: %v\n", stringArray)
	fmt.Printf("Floats: %v\n", floatArray)
	fmt.Printf("Bools: %v\n", boolArray)

	fmt.Println()
}

// ==============================================================================
// ARRAY ITERATION
// ==============================================================================
// Arrays can be iterated using range loops or traditional for loops.

func demonstrateArrayIteration() {
	fmt.Println("--- ARRAY ITERATION ---")

	array := [...]int{1, 2, 3, 4, 5}

	// Method 1: range with index only
	fmt.Println("Using range with index:")
	for i := range array {
		fmt.Printf("  array[%d] = %d\n", i, array[i])
	}

	// Method 2: range with index and value (value is a COPY)
	fmt.Println("\nUsing range with index and value:")
	for i, value := range array {
		fmt.Printf("  index %d: %d\n", i, value)
	}

	// Method 3: range with value only (ignore index with _)
	fmt.Println("\nUsing range with value only:")
	for _, value := range array {
		fmt.Printf("  %d ", value)
	}
	fmt.Println()

	// Method 4: traditional for loop with len
	fmt.Println("\nUsing traditional for loop:")
	for i := 0; i < len(array); i++ {
		fmt.Printf("  array[%d] = %d\n", i, array[i])
	}

	// Traditional for loop allows custom iteration patterns
	fmt.Println("\nReverse iteration:")
	for i := len(array) - 1; i >= 0; i-- {
		fmt.Printf("  %d ", array[i])
	}
	fmt.Println()

	fmt.Println("\nSkip every other element:")
	for i := 0; i < len(array); i += 2 {
		fmt.Printf("  array[%d] = %d\n", i, array[i])
	}

	// IMPORTANT: range provides COPIES of values
	fmt.Println("\nAttempting to modify via range (doesn't work):")
	strings := [3]string{"hello", "world", "go"}
	for _, s := range strings {
		s = s + "!" // This modifies the COPY, not the original
	}
	fmt.Printf("Array unchanged: %v\n", strings)

	// To modify, use index:
	fmt.Println("\nModifying via index (works):")
	for i := range strings {
		strings[i] = strings[i] + "!"
	}
	fmt.Printf("Array modified: %v\n", strings)

	fmt.Println()
}

// ==============================================================================
// ARRAYS AS VALUES
// ==============================================================================
// Arrays are VALUES in Go - assignment creates a complete copy.

func demonstrateArrayValues() {
	fmt.Println("--- ARRAYS AS VALUES ---")

	// Arrays are copied on assignment
	colors := [3]string{"Red", "Green", "Blue"}
	colorsCopy := colors // Creates a COMPLETE COPY

	// Modifying the copy doesn't affect the original
	colorsCopy[0] = "Yellow"

	fmt.Printf("Original: %v\n", colors)     // [Red Green Blue]
	fmt.Printf("Copy:     %v\n", colorsCopy) // [Yellow Green Blue]

	// Passing arrays to functions also creates copies
	fmt.Println("\nPassing arrays to functions:")
	numbers := [3]int{1, 2, 3}
	fmt.Printf("Before function: %v\n", numbers)
	modifyArray(numbers) // Function gets a COPY
	fmt.Printf("After function:  %v\n", numbers) // Unchanged!

	// Use pointers to modify original arrays
	fmt.Println("\nUsing pointer to modify original:")
	modifyArrayPointer(&numbers)
	fmt.Printf("After pointer function: %v\n", numbers) // Modified!

	// Arrays of pointers share the pointed-to data
	fmt.Println("\nArrays of pointers:")
	val1, val2 := 10, 20
	ptrArray := [2]*int{&val1, &val2}
	ptrArrayCopy := ptrArray // Copies the pointers, not the values

	*ptrArrayCopy[0] = 100 // Modifies the value that both arrays point to
	fmt.Printf("val1 = %d (modified via copy!)\n", val1)

	fmt.Println()
}

// modifyArray receives a COPY of the array
func modifyArray(arr [3]int) {
	arr[0] = 999 // Modifies the copy only
	fmt.Printf("  Inside function: %v\n", arr)
}

// modifyArrayPointer receives a pointer to the original array
func modifyArrayPointer(arr *[3]int) {
	arr[0] = 999 // Modifies the original
	fmt.Printf("  Inside pointer function: %v\n", *arr)
}

// ==============================================================================
// MULTIDIMENSIONAL ARRAYS
// ==============================================================================
// Arrays can contain other arrays to create multidimensional structures.

func demonstrateMultidimensionalArrays() {
	fmt.Println("--- MULTIDIMENSIONAL ARRAYS ---")

	// Declare a 2D array (array of arrays)
	var grid [3][3]int

	// Initialize with values
	grid[0][0] = 1
	grid[0][1] = 2
	grid[0][2] = 3
	grid[1][0] = 4
	grid[1][1] = 5
	grid[1][2] = 6
	grid[2][0] = 7
	grid[2][1] = 8
	grid[2][2] = 9

	fmt.Println("2D array (grid):")
	for i := range grid {
		fmt.Printf("  %v\n", grid[i])
	}

	// Initialize with array literal
	matrix := [2][3]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Printf("\nMatrix: %v\n", matrix)

	// Sparse initialization
	sparse := [3][3]int{
		1: {1: 5}, // Only set element [1][1] = 5
	}
	fmt.Printf("Sparse 2D array: %v\n", sparse)

	// Accessing elements requires multiple brackets
	fmt.Printf("\nmatrix[0][0] = %d\n", matrix[0][0])
	fmt.Printf("matrix[1][2] = %d\n", matrix[1][2])

	// You can extract sub-arrays
	row := matrix[0] // Gets [1 2 3] as a [3]int
	fmt.Printf("First row: %v (type: %T)\n", row, row)

	// 3D arrays are possible too
	cube := [2][2][2]int{
		{{1, 2}, {3, 4}},
		{{5, 6}, {7, 8}},
	}
	fmt.Printf("\n3D array: %v\n", cube)
	fmt.Printf("cube[0][1][0] = %d\n", cube[0][1][0])

	fmt.Println()
}

// ==============================================================================
// SLICES - DYNAMIC ORDERED COLLECTIONS
// ==============================================================================
// Slices are Go's primary ordered collection type - flexible and efficient.
// A slice is a view into an array with three components:
// 1. Pointer to the backing array
// 2. Length (number of elements accessible)
// 3. Capacity (total space available for growth)

func demonstrateSlices() {
	fmt.Println("--- SLICES BASICS ---")

	// Slice literal (most common way to create slices)
	slice := []int{1, 2, 3, 4, 5}
	fmt.Printf("Slice: %v\n", slice)
	fmt.Printf("  Length: %d\n", len(slice))
	fmt.Printf("  Capacity: %d\n", cap(slice))

	// Access elements just like arrays
	fmt.Printf("  First: %d, Last: %d\n", slice[0], slice[len(slice)-1])

	// Modify elements
	slice[0] = 10
	fmt.Printf("After modification: %v\n", slice)

	// Slices can be used with range loops
	fmt.Print("Elements: ")
	for _, v := range slice {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// Key difference from arrays: slices can GROW
	slice = append(slice, 6, 7, 8)
	fmt.Printf("After append: %v (len: %d, cap: %d)\n", slice, len(slice), cap(slice))

	fmt.Println()
}

// ==============================================================================
// SLICE DECLARATIONS - DIFFERENT FORMS
// ==============================================================================
// Slices can be created in several ways, each with different use cases.

func demonstrateSliceDeclarations() {
	fmt.Println("--- SLICE DECLARATIONS ---")

	// 1. Slice literal (creates backing array automatically)
	literal := []int{1, 2, 3}
	fmt.Printf("Literal: %v (len: %d, cap: %d)\n", literal, len(literal), cap(literal))

	// 2. make with length only (creates slice with zeroed values)
	withLen := make([]int, 3)
	fmt.Printf("make([]int, 3): %v (len: %d, cap: %d)\n", withLen, len(withLen), cap(withLen))

	// 3. make with length and capacity (preallocates for growth)
	withCap := make([]int, 3, 5)
	fmt.Printf("make([]int, 3, 5): %v (len: %d, cap: %d)\n", withCap, len(withCap), cap(withCap))
	// The 2 extra capacity elements are allocated but not accessible yet

	// 4. make with zero length but non-zero capacity (ready for append)
	readyToGrow := make([]int, 0, 5)
	fmt.Printf("make([]int, 0, 5): %v (len: %d, cap: %d)\n", readyToGrow, len(readyToGrow), cap(readyToGrow))

	// 5. var declaration (creates nil slice)
	var nilSlice []int
	fmt.Printf("var nilSlice: %v (len: %d, cap: %d, nil: %v)\n", 
		nilSlice, len(nilSlice), cap(nilSlice), nilSlice == nil)

	// NIL SLICES: Not usable until initialized
	// This would panic: value := nilSlice[0]

	// Empty slice vs nil slice
	emptySlice := []int{}
	fmt.Printf("\nEmpty slice:   %v (nil: %v)\n", emptySlice, emptySlice == nil)
	fmt.Printf("Nil slice:     %v (nil: %v)\n", nilSlice, nilSlice == nil)

	// Both have length 0, but only nil slice == nil
	fmt.Println("\nBoth are empty but behave differently:")
	fmt.Printf("  len(emptySlice) == 0: %v\n", len(emptySlice) == 0)
	fmt.Printf("  len(nilSlice) == 0:   %v\n", len(nilSlice) == 0)
	fmt.Printf("  emptySlice == nil:    %v\n", emptySlice == nil)
	fmt.Printf("  nilSlice == nil:      %v\n", nilSlice == nil)

	// BEST PRACTICE: Use len() to check if slice is empty, not == nil
	fmt.Println("\nChecking for empty slices:")
	fmt.Printf("  Empty check: len(emptySlice) == 0: %v\n", len(emptySlice) == 0)
	fmt.Printf("  Empty check: len(nilSlice) == 0:   %v\n", len(nilSlice) == 0)

	fmt.Println()
}

// ==============================================================================
// APPEND - GROWING SLICES
// ==============================================================================
// The append function adds elements to slices, allocating storage as needed.

func demonstrateSliceAppend() {
	fmt.Println("--- APPEND OPERATION ---")

	// Start with a slice
	slice := []int{1, 2, 3}
	fmt.Printf("Initial: %v (len: %d, cap: %d)\n", slice, len(slice), cap(slice))

	// Append a single value
	slice = append(slice, 4)
	fmt.Printf("After append(4): %v (len: %d, cap: %d)\n", slice, len(slice), cap(slice))
	// Note: capacity doubled!

	// Append multiple values
	slice = append(slice, 5, 6)
	fmt.Printf("After append(5,6): %v (len: %d, cap: %d)\n", slice, len(slice), cap(slice))

	// Append another slice (use ... to unpack)
	moreValues := []int{7, 8, 9}
	slice = append(slice, moreValues...)
	fmt.Printf("After append(slice...): %v (len: %d, cap: %d)\n", slice, len(slice), cap(slice))

	// Append works on nil slices
	var nilSlice []string
	nilSlice = append(nilSlice, "first")
	fmt.Printf("\nAppend to nil slice: %v\n", nilSlice)

	// Demonstrating capacity growth strategy
	fmt.Println("\nCapacity growth pattern:")
	demo := []int{}
	for i := 0; i < 20; i++ {
		demo = append(demo, i)
		fmt.Printf("  len: %2d, cap: %2d\n", len(demo), cap(demo))
	}
	// Notice how capacity grows: doubles when small, grows by 25% when large

	// IMPORTANT: Always assign the result of append back to the slice
	// This is necessary because append may return a new slice with new backing array

	// Demonstrating why you must assign back
	fmt.Println("\nWhy you must assign append result:")
	original := []int{1, 2, 3}
	demonstration := original
	fmt.Printf("Before append - original: %v, demonstration: %v\n", original, demonstration)
	
	demonstration = append(demonstration, 4, 5, 6, 7, 8) // New backing array!
	fmt.Printf("After append - original: %v, demonstration: %v\n", original, demonstration)
	fmt.Println("They no longer share the same backing array!")

	fmt.Println()
}

// ==============================================================================
// SLICE EXPRESSIONS - CREATING SLICES FROM SLICES/ARRAYS
// ==============================================================================
// Slice expressions create new slices that share backing arrays.

func demonstrateSliceExpressions() {
	fmt.Println("--- SLICE EXPRESSIONS ---")

	// Create a slice from another slice
	original := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("Original: %v\n", original)

	// Slice expression: [low:high] (low inclusive, high exclusive)
	sub1 := original[2:5] // Elements at indices 2, 3, 4
	fmt.Printf("original[2:5]:   %v\n", sub1)

	sub2 := original[:3] // From start up to index 3 (exclusive)
	fmt.Printf("original[:3]:    %v\n", sub2)

	sub3 := original[7:] // From index 7 to end
	fmt.Printf("original[7:]:    %v\n", sub3)

	sub4 := original[:] // Entire slice (creates new slice header)
	fmt.Printf("original[:]:     %v\n", sub4)

	// IMPORTANT: Slices created with expressions SHARE the backing array
	fmt.Println("\nDemonstrating shared backing array:")
	numbers := []int{1, 2, 3, 4, 5}
	subset := numbers[1:4] // [2 3 4]
	fmt.Printf("numbers: %v, subset: %v\n", numbers, subset)

	// Modifying subset affects original
	subset[0] = 99
	fmt.Printf("After subset[0]=99:\n")
	fmt.Printf("  numbers: %v (changed!)\n", numbers)
	fmt.Printf("  subset:  %v\n", subset)

	// Modifying original affects subset
	numbers[3] = 88
	fmt.Printf("After numbers[3]=88:\n")
	fmt.Printf("  numbers: %v\n", numbers)
	fmt.Printf("  subset:  %v (changed!)\n", subset)

	// Slice expressions with capacity limit [low:high:max]
	fmt.Println("\nSlice expressions with capacity control:")
	original2 := []int{0, 1, 2, 3, 4, 5}
	
	// Without max (inherits full remaining capacity)
	slice1 := original2[:2]
	fmt.Printf("original2[:2]:     %v (len: %d, cap: %d)\n", slice1, len(slice1), cap(slice1))
	
	// With max (limits capacity)
	slice2 := original2[:2:3]
	fmt.Printf("original2[:2:3]:   %v (len: %d, cap: %d)\n", slice2, len(slice2), cap(slice2))

	// Creating slice from array
	fmt.Println("\nCreating slice from array:")
	array := [5]int{10, 20, 30, 40, 50}
	sliceFromArray := array[1:4]
	fmt.Printf("array[1:4]: %v (type: %T)\n", sliceFromArray, sliceFromArray)
	// Now it's a slice, backed by the array

	fmt.Println()
}

// ==============================================================================
// COPY - DUPLICATING SLICES
// ==============================================================================
// The copy function creates independent copies of slice data.

func demonstrateSliceCopy() {
	fmt.Println("--- COPY OPERATION ---")

	// Create source slice
	source := []int{1, 2, 3, 4, 5}
	fmt.Printf("Source: %v\n", source)

	// Create destination with same length
	dest := make([]int, len(source))
	
	// Copy returns number of elements copied
	copied := copy(dest, source)
	fmt.Printf("Copied %d elements\n", copied)
	fmt.Printf("Destination: %v\n", dest)

	// Now they're independent - modifying one doesn't affect the other
	dest[0] = 999
	fmt.Printf("After dest[0]=999:\n")
	fmt.Printf("  source: %v (unchanged)\n", source)
	fmt.Printf("  dest:   %v\n", dest)

	// Copy behavior with different lengths
	fmt.Println("\nCopy with different lengths:")

	// Destination is shorter - copies only what fits
	short := make([]int, 3)
	n := copy(short, source)
	fmt.Printf("Short dest (len=3): %v (copied %d elements)\n", short, n)

	// Destination is longer - copies all source, leaves rest unchanged
	long := make([]int, 7)
	n = copy(long, source)
	fmt.Printf("Long dest (len=7): %v (copied %d elements)\n", long, n)

	// Source is shorter - copies all source
	longSrc := []int{10, 20, 30, 40, 50, 60, 70}
	shortDest := make([]int, 3)
	n = copy(shortDest, longSrc)
	fmt.Printf("Short dest, long src: %v (copied %d elements)\n", shortDest, n)

	// IMPORTANT: copy doesn't allocate - destination must be pre-sized
	var nilDest []int
	n = copy(nilDest, source) // Nothing copied!
	fmt.Printf("\nCopy to nil slice: %v (copied %d elements)\n", nilDest, n)

	// Practical use: creating independent copy of a sub-slice
	fmt.Println("\nCreating independent copy of sub-slice:")
	original := []int{1, 2, 3, 4, 5, 6}
	subSlice := original[2:5] // [3 4 5] - shares backing array
	
	// Create independent copy
	independent := make([]int, len(subSlice))
	copy(independent, subSlice)
	
	original[3] = 999
	fmt.Printf("After modifying original:\n")
	fmt.Printf("  original:    %v\n", original)
	fmt.Printf("  subSlice:    %v (affected!)\n", subSlice)
	fmt.Printf("  independent: %v (not affected!)\n", independent)

	fmt.Println()
}

// ==============================================================================
// AVOIDING SLICE PANICS
// ==============================================================================
// Accessing invalid indices causes runtime panics. Use len() to avoid them.

func demonstrateSlicePanics() {
	fmt.Println("--- AVOIDING SLICE PANICS ---")

	// Safe slice access
	numbers := []int{1, 2, 3}
	fmt.Printf("Slice: %v (len: %d)\n", numbers, len(numbers))

	// This is safe
	for i := 0; i < len(numbers); i++ {
		fmt.Printf("  numbers[%d] = %d\n", i, numbers[i])
	}

	// Use len() to check before accessing
	index := 5
	if index < len(numbers) {
		fmt.Printf("numbers[%d] = %d\n", index, numbers[index])
	} else {
		fmt.Printf("Index %d is out of bounds (len: %d)\n", index, len(numbers))
	}

	// Use range loops to avoid index issues
	fmt.Println("\nUsing range (always safe):")
	for i, v := range numbers {
		fmt.Printf("  [%d]: %d\n", i, v)
	}

	// Flexible function that handles any slice size
	fmt.Println("\nFlexible medal printing:")
	printMedalists("Race 1", []string{"Alice", "Bob", "Charlie", "Dave"})
	printMedalists("Race 2", []string{"Eve", "Frank"})
	printMedalists("Race 3", []string{})

	fmt.Println()
}

// printMedalists handles slices of any size using range
func printMedalists(race string, contestants []string) {
	fmt.Printf("\n%s results:\n", race)
	
	if len(contestants) == 0 {
		fmt.Println("  No contestants")
		return
	}

	medals := []string{"Gold", "Silver", "Bronze"}
	
	for i, contestant := range contestants {
		if i < len(medals) {
			fmt.Printf("  %s Medal: %s\n", medals[i], contestant)
		} else {
			fmt.Printf("  Finisher: %s\n", contestant)
		}
	}
}

// ==============================================================================
// MAPS - KEY-VALUE COLLECTIONS
// ==============================================================================
// Maps store key-value pairs with fast lookups by key.
// Maps are unordered and use hash tables internally.

func demonstrateMaps() {
	fmt.Println("--- MAPS BASICS ---")

	// Map literal (most common way to create maps)
	ages := map[string]int{
		"Alice":   30,
		"Bob":     25,
		"Charlie": 35,
	}
	fmt.Printf("Ages map: %v\n", ages)

	// make creates an empty map
	scores := make(map[string]int)
	fmt.Printf("Empty map: %v (len: %d)\n", scores, len(scores))

	// Add values to map
	scores["Alice"] = 100
	scores["Bob"] = 95
	scores["Charlie"] = 87
	fmt.Printf("After adding: %v (len: %d)\n", scores, len(scores))

	// Retrieve values by key
	aliceScore := scores["Alice"]
	fmt.Printf("Alice's score: %d\n", aliceScore)

	// Accessing non-existent key returns zero value
	missingScore := scores["Dave"] // Returns 0 (zero value for int)
	fmt.Printf("Dave's score (not in map): %d\n", missingScore)

	// Two-value retrieval to check existence
	bobScore, found := scores["Bob"]
	fmt.Printf("Bob's score: %d (found: %v)\n", bobScore, found)

	daveScore, found := scores["Dave"]
	fmt.Printf("Dave's score: %d (found: %v)\n", daveScore, found)

	// Idiomatic existence check with if
	if score, found := scores["Alice"]; found {
		fmt.Printf("Found Alice with score: %d\n", score)
	}

	// Delete key from map
	delete(scores, "Bob")
	fmt.Printf("After deleting Bob: %v\n", scores)

	// Deleting non-existent key is safe (no panic)
	delete(scores, "NotThere")
	fmt.Printf("After deleting non-existent key: %v\n", scores)

	// Nil map (not usable until initialized)
	var nilMap map[string]int
	fmt.Printf("\nNil map: %v (len: %d, nil: %v)\n", nilMap, len(nilMap), nilMap == nil)
	
	// IMPORTANT: Cannot add to nil map (would panic)
	// nilMap["key"] = value // panic: assignment to entry in nil map
	
	// Must initialize first
	if nilMap == nil {
		nilMap = make(map[string]int)
	}
	nilMap["key"] = 42
	fmt.Printf("After initialization: %v\n", nilMap)

	fmt.Println()
}

// ==============================================================================
// MAP OPERATIONS
// ==============================================================================
// Common patterns for working with maps.

func demonstrateMapOperations() {
	fmt.Println("--- MAP OPERATIONS ---")

	// Using zero values for convenience
	fmt.Println("1. Zero values for membership sets:")
	members := map[string]bool{
		"Alice": true,
		"Bob":   true,
		"Carol": true,
	}

	for _, name := range []string{"Alice", "Dave", "Carol"} {
		if members[name] {
			fmt.Printf("  %s is a member\n", name)
		} else {
			fmt.Printf("  %s is not a member\n", name)
		}
	}

	// Using zero values for counters
	fmt.Println("\n2. Zero values for counters:")
	wordCount := make(map[string]int)
	words := []string{"apple", "banana", "apple", "cherry", "banana", "apple"}
	
	for _, word := range words {
		wordCount[word]++ // Zero value (0) allows this to work
	}
	fmt.Printf("  Word counts: %v\n", wordCount)

	// Building complex maps
	fmt.Println("\n3. Maps with complex value types:")
	
	// Map with slice values
	groups := make(map[string][]string)
	groups["fruits"] = []string{"apple", "banana"}
	groups["vegetables"] = []string{"carrot", "broccoli"}
	groups["fruits"] = append(groups["fruits"], "orange")
	
	fmt.Printf("  Groups: %v\n", groups)

	// Map with struct values
	type Person struct {
		Name string
		Age  int
	}
	
	people := map[string]Person{
		"alice": {Name: "Alice", Age: 30},
		"bob":   {Name: "Bob", Age: 25},
	}
	fmt.Printf("  People: %v\n", people)

	// Map with struct keys
	type Point struct {
		X, Y int
	}
	
	locations := map[Point]string{
		{0, 0}:   "origin",
		{10, 20}: "point A",
		{30, 40}: "point B",
	}
	fmt.Printf("  Locations: %v\n", locations)
	
	// Look up by struct key
	if name, found := locations[Point{10, 20}]; found {
		fmt.Printf("  Found location: %s\n", name)
	}

	fmt.Println()
}

// ==============================================================================
// MAP ITERATION
// ==============================================================================
// Maps can be iterated with range, but order is not guaranteed.

func demonstrateMapIteration() {
	fmt.Println("--- MAP ITERATION ---")

	colors := map[string]string{
		"red":    "#FF0000",
		"green":  "#00FF00",
		"blue":   "#0000FF",
		"yellow": "#FFFF00",
	}

	// Iterate with key only
	fmt.Println("Keys only:")
	for color := range colors {
		fmt.Printf("  %s\n", color)
	}

	// Iterate with key and value
	fmt.Println("\nKeys and values:")
	for color, hex := range colors {
		fmt.Printf("  %s: %s\n", color, hex)
	}

	// Ignore key with _
	fmt.Println("\nValues only:")
	for _, hex := range colors {
		fmt.Printf("  %s\n", hex)
	}

	// IMPORTANT: Map iteration order is random
	fmt.Println("\nMultiple iterations (notice random order):")
	for i := 0; i < 3; i++ {
		fmt.Printf("Iteration %d: ", i+1)
		for color := range colors {
			fmt.Printf("%s ", color)
		}
		fmt.Println()
	}

	// Filtering maps by deleting during iteration
	fmt.Println("\nFiltering map (remove short hex codes):")
	testMap := map[string]string{
		"a": "#F00",     // Short (will be removed)
		"b": "#00FF00",  // Long
		"c": "#00F",     // Short (will be removed)
		"d": "#0000FF",  // Long
	}
	
	fmt.Printf("Before: %v\n", testMap)
	for key, value := range testMap {
		if len(value) < 7 {
			delete(testMap, key)
		}
	}
	fmt.Printf("After filtering: %v\n", testMap)

	// Counting and aggregating
	fmt.Println("\nCounting and aggregating:")
	inventory := map[string]int{
		"apples":  10,
		"oranges": 5,
		"bananas": 8,
		"grapes":  15,
	}
	
	total := 0
	for _, count := range inventory {
		total += count
	}
	fmt.Printf("Total items: %d\n", total)

	fmt.Println()
}

// ==============================================================================
// KEY CONCEPTS SUMMARY
// ==============================================================================
//
// ARRAYS
// - Fixed-length collections: [5]int is different from [4]int
// - Length is part of the type
// - Values, not references: assignment copies all elements
// - Use for fixed-size data, but slices are usually better
// - Pass pointers to functions to avoid copying
//
// SLICES
// - Dynamic, flexible ordered collections (primary collection type)
// - Three components: pointer, length, capacity
// - Created with literals: []int{1,2,3} or make: make([]int, len, cap)
// - Grow with append: slice = append(slice, values...)
// - Share backing arrays: modifications visible across slice expressions
// - Use copy() to create independent copies
// - Check len() before accessing to avoid panics
// - Use range loops for flexibility
// - nil check: use len(slice) == 0, not slice == nil
//
// SLICE EXPRESSIONS
// - Create new slices from existing: slice[low:high]
// - Share backing array with original
// - Can limit capacity: slice[low:high:max]
// - low defaults to 0, high defaults to len()
//
// APPEND
// - Always assign result back: slice = append(slice, values...)
// - May allocate new backing array (capacity growth)
// - Works on nil slices
// - Can append multiple values or another slice with ...
//
// MAPS
// - Unordered key-value collections: map[K]V
// - Create with literal or make: make(map[string]int)
// - Access: value := map[key]
// - Check existence: value, found := map[key]
// - Delete: delete(map, key)
// - Zero value returned for missing keys
// - nil maps must be initialized before use
// - Iteration order is random (by design)
// - Keys must be comparable (can use ==)
//
// BEST PRACTICES
// - Use slices for ordered collections (not arrays)
// - Use make() to preallocate if you know the size
// - Always check len() before accessing indices
// - Use range loops for safe iteration
// - Assign append results back to slice
// - Use copy() when you need independent slices
// - Initialize maps with make() or literals (not var)
// - Use two-value retrieval to check map key existence
//
// ==============================================================================
