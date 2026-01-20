package main

import (
	"cmp"
	"errors"
	"fmt"
	"maps"
	"slices"

	"golang.org/x/exp/constraints"
)

// ============================================================================
// GENERICS IN GO
// ============================================================================
// This file demonstrates Go's generics system, introduced in Go 1.18.
// Generics allow us to write reusable, type-safe code that works with
// different types while maintaining compile-time safety.
//
// Key Concepts Covered:
// 1. Type Parameters - placeholders for types in generic definitions
// 2. Type Constraints - restrictions on what types can be used
// 3. Generic Interfaces - interfaces that use type parameters
// 4. Generic Structs - data structures that work with any type
// 5. Multiple Type Parameters - relating multiple types together
// 6. Standard Library Support - constraints, slices, maps, cmp packages
//
// Think of generics as a way to write "templates" for code that can
// operate on many different types without sacrificing type safety.
// ============================================================================

// ============================================================================
// 1. BASIC GENERIC STRUCTURES
// ============================================================================
// We start with the fundamental building blocks of generic programming:
// type parameters and simple generic types.

// node is a generic linked list node that can hold any type.
// The [T any] syntax declares a type parameter T that can be any type.
// This is the foundation of generic data structures in Go.
type node[T any] struct {
	value T          // The value can be of any type T
	next  *node[T]   // The next pointer references another node of the same type
}

// ============================================================================
// 2. GENERIC INTERFACES
// ============================================================================
// Generic interfaces define contracts that work with any type.
// They allow us to write reusable abstractions for collections,
// algorithms, and other patterns.

// Collection defines a generic interface for any collection type.
// This provides a common set of operations that all collections should support,
// regardless of what type of elements they contain.
type Collection[T any] interface {
	Add(item T)           // Add an item to the collection
	Remove(item T) bool   // Remove an item, return true if found
	Contains(item T) bool // Check if item exists in collection
	Size() int            // Return number of items
	Clear()               // Remove all items
}

// List extends Collection to provide indexed access to elements.
// This demonstrates interface embedding with generics - List "is a" Collection
// with additional capabilities specific to ordered sequences.
type List[T any] interface {
	Collection[T]                      // Embed the Collection interface
	Get(index int) (T, error)          // Get item at index
	Set(index int, item T) error       // Set item at index
	Insert(index int, item T) error    // Insert item at index
	RemoveAt(index int) (T, error)     // Remove item at index
}

// ============================================================================
// 3. TYPE CONSTRAINTS
// ============================================================================
// Type constraints restrict which types can be used with a generic.
// The most common built-in constraint is 'comparable', which ensures
// types support == and != operators.

// LinkedList is a generic linked list that only accepts comparable types.
// We use 'comparable' instead of 'any' because we need to compare elements
// in methods like Contains and Remove.
//
// Design Decision: Why comparable?
// - Enables equality checks (==, !=)
// - Required for Contains and Remove operations
// - More restrictive than 'any' but provides necessary functionality
type LinkedList[T comparable] struct {
	head *node[T] // Pointer to first node
	size int      // Number of elements in list
}

// NewLinkedList creates a new empty linked list for any comparable type.
// This is a generic constructor function - it returns a LinkedList[T]
// where T is determined by the caller.
//
// Usage examples:
//   intList := NewLinkedList[int]()
//   stringList := NewLinkedList[string]()
func NewLinkedList[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{}
}

// ============================================================================
// IMPLEMENTING COLLECTION INTERFACE
// ============================================================================
// These methods fulfill the Collection[T] interface contract.
// Each method maintains the invariants of the linked list structure.

// Add appends an item to the end of the list.
// This is an O(n) operation as we must traverse to the end.
//
// Implementation notes:
// - Creates a new node wrapping the value
// - Handles empty list case (set as head)
// - Otherwise traverses to end and appends
func (l *LinkedList[T]) Add(item T) {
	n := &node[T]{value: item}

	if l.head == nil {
		// Empty list: new node becomes head
		l.head = n
	} else {
		// Non-empty: traverse to end and append
		curr := l.head
		for curr.next != nil {
			curr = curr.next
		}
		curr.next = n
	}
	l.size++
}

// Remove deletes the first occurrence of item from the list.
// Returns true if the item was found and removed, false otherwise.
//
// Why this needs comparable:
// We use == to compare curr.value with item. This requires T to be comparable.
func (l *LinkedList[T]) Remove(item T) bool {
	var prev *node[T]
	curr := l.head

	for curr != nil {
		if curr.value == item {
			// Found the item to remove
			if prev == nil {
				// Removing head node
				l.head = curr.next
			} else {
				// Removing middle or end node
				prev.next = curr.next
			}
			l.size--
			return true
		}
		prev = curr
		curr = curr.next
	}
	return false
}

// Contains checks if an item exists in the list.
// Returns true if found, false otherwise.
// This is an O(n) operation requiring a full traversal in worst case.
func (l *LinkedList[T]) Contains(item T) bool {
	curr := l.head
	for curr != nil {
		if curr.value == item {
			return true
		}
		curr = curr.next
	}
	return false
}

// Size returns the number of elements in the list.
// This is O(1) as we maintain a size counter.
func (l *LinkedList[T]) Size() int {
	return l.size
}

// Clear removes all elements from the list.
// This is O(1) - we just reset the head pointer and size.
// The garbage collector will clean up the orphaned nodes.
func (l *LinkedList[T]) Clear() {
	l.head = nil
	l.size = 0
}

// ============================================================================
// IMPLEMENTING LIST INTERFACE
// ============================================================================
// These methods provide indexed access to elements, making LinkedList
// behave more like an array (despite being a linked structure).

// Get retrieves the value at the specified index.
// Returns an error if the index is out of bounds.
//
// Important pattern: Zero values in generics
// We declare 'var zero T' to get the zero value of type T.
// Since we don't know what T is, we can't write 'return 0' or 'return nil'.
// The zero value is: 0 for numbers, "" for strings, nil for pointers, etc.
func (l *LinkedList[T]) Get(index int) (T, error) {
	var zero T // Zero value of type T

	if index < 0 || index >= l.size {
		return zero, errors.New("index out of bounds")
	}

	// Traverse to the index
	curr := l.head
	for i := 0; i < index; i++ {
		curr = curr.next
	}
	return curr.value, nil
}

// Set updates the value at the specified index.
// Returns an error if the index is out of bounds.
func (l *LinkedList[T]) Set(index int, item T) error {
	if index < 0 || index >= l.size {
		return errors.New("index out of bounds")
	}

	// Traverse to the index
	curr := l.head
	for i := 0; i < index; i++ {
		curr = curr.next
	}
	curr.value = item
	return nil
}

// Insert inserts a value at the specified index, shifting subsequent elements.
// Allows inserting at index == size (equivalent to append).
func (l *LinkedList[T]) Insert(index int, item T) error {
	if index < 0 || index > l.size {
		return errors.New("index out of bounds")
	}

	n := &node[T]{value: item}

	if index == 0 {
		// Insert at head
		n.next = l.head
		l.head = n
	} else {
		// Insert in middle or at end
		curr := l.head
		for i := 0; i < index-1; i++ {
			curr = curr.next
		}
		n.next = curr.next
		curr.next = n
	}
	l.size++
	return nil
}

// RemoveAt removes and returns the value at the specified index.
// Returns an error if the index is out of bounds.
func (l *LinkedList[T]) RemoveAt(index int) (T, error) {
	var zero T

	if index < 0 || index >= l.size {
		return zero, errors.New("index out of bounds")
	}

	var prev *node[T]
	curr := l.head

	// Traverse to the index
	for i := 0; i < index; i++ {
		prev = curr
		curr = curr.next
	}

	if prev == nil {
		// Removing head
		l.head = curr.next
	} else {
		// Removing middle or end
		prev.next = curr.next
	}
	l.size--
	return curr.value, nil
}

// ============================================================================
// 4. GENERIC FUNCTIONS
// ============================================================================
// Functions can also be generic, operating on any type that satisfies
// their constraints.

// PrintCollection demonstrates a generic function that works with
// any Collection[T]. This shows how we can write utilities that operate
// on generic interfaces.
//
// Type parameter T must be 'any' (not comparable) because the Collection
// interface uses 'any', making this function more flexible.
func PrintCollection[T any](c Collection[T]) {
	fmt.Printf("Collection has %d items:\n", c.Size())

	// Type assertion to check if this collection is also a List
	// This demonstrates runtime type checking with generics
	if list, ok := c.(List[T]); ok {
		for i := 0; i < c.Size(); i++ {
			val, _ := list.Get(i)
			fmt.Printf("  [%d] %v\n", i, val)
		}
	}
}

// ReverseList reverses the elements of any List[T] in-place.
// This demonstrates a generic algorithm that works on an interface.
//
// Algorithm: Swap elements from both ends moving toward the middle
// - Get element at index i and index (n-1-i)
// - Swap them
// - Continue for first half of list
func ReverseList[T comparable](list List[T]) {
	n := list.Size()
	for i := 0; i < n/2; i++ {
		// Get values at mirror positions
		a, _ := list.Get(i)
		b, _ := list.Get(n - 1 - i)

		// Swap them
		list.Set(i, b)
		list.Set(n-1-i, a)
	}
}

// ============================================================================
// 5. CUSTOM TYPE CONSTRAINTS
// ============================================================================
// You can define custom constraints to restrict type parameters to
// specific sets of types. This is useful for numeric operations,
// custom interfaces, or any other common type characteristic.

// Numeric is a custom constraint that allows only numeric types.
// The | operator creates a type union - T can be int OR float32 OR float64.
//
// Note: This matches only exact types. MyInt (defined as 'type MyInt int')
// would NOT satisfy this constraint.
type Numeric interface {
	int | float32 | float64
}

// NumericWithUnderlying is a more flexible numeric constraint.
// The ~ (tilde) operator means "any type whose underlying type is...".
// This allows user-defined types like 'type MyInt int' to satisfy the constraint.
//
// Design choice: When to use ~
// - Use ~ when you want to accept user-defined types
// - Don't use ~ when you need only exact built-in types
type NumericWithUnderlying interface {
	~int | ~float32 | ~float64
}

// Sum adds all elements in a slice of any numeric type.
// This demonstrates using a custom constraint for a mathematical operation.
func Sum[T Numeric](values []T) T {
	var total T // Zero value (0 for numeric types)

	for _, v := range values {
		total += v // The + operator works because T is constrained to numeric types
	}
	return total
}

// Ordered is a custom constraint for types that can be ordered.
// This enables comparison operations like <, >, <=, >=.
// We use ~ to allow user-defined types built on these underlying types.
type Ordered interface {
	~int | ~float32 | ~float64 | ~string
}

// Max returns the maximum of two values of any ordered type.
// This works with int, float, string, or any user-defined type
// with one of these as its underlying type.
func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min returns the minimum of two values of any ordered type.
func Min[T Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// ============================================================================
// 6. MULTIPLE TYPE PARAMETERS
// ============================================================================
// Many data structures need to work with multiple types simultaneously.
// Maps, for example, have both key and value types.

// Map is a generic interface for map-like collections.
// It uses two type parameters:
// - K (key type) must be comparable (to check key equality)
// - V (value type) can be any type
//
// This is a common pattern: keys need comparison, values can be anything.
type Map[K comparable, V any] interface {
	Put(key K, value V)       // Insert or update key-value pair
	Get(key K) (V, bool)      // Get value for key, bool indicates if found
	Remove(key K) bool        // Remove key-value pair
	ContainsKey(key K) bool   // Check if key exists
	Size() int                // Number of key-value pairs
	Clear()                   // Remove all pairs
	Keys() []K                // Get all keys
	Values() []V              // Get all values
}

// HashMap is a concrete implementation of Map using Go's built-in map.
// It demonstrates how to implement a generic interface with multiple
// type parameters.
//
// Implementation notes:
// - Wraps Go's built-in map for type safety and additional methods
// - Could be extended with thread safety, ordering, etc.
type HashMap[K comparable, V any] struct {
	m map[K]V // The underlying Go map
}

// NewHashMap creates a new empty HashMap.
// This demonstrates a constructor for a type with multiple type parameters.
//
// Usage:
//   m := NewHashMap[string, int]()     // Map strings to ints
//   m := NewHashMap[int, *User]()      // Map ints to User pointers
func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{
		m: make(map[K]V),
	}
}

// Put inserts or updates a key-value pair.
func (h *HashMap[K, V]) Put(key K, value V) {
	h.m[key] = value
}

// Get retrieves the value for a key.
// Returns the value and true if found, zero value and false if not found.
//
// This follows Go's map access pattern: val, ok := m[key]
func (h *HashMap[K, V]) Get(key K) (V, bool) {
	val, ok := h.m[key]
	return val, ok
}

// Remove deletes a key-value pair.
// Returns true if the key existed, false otherwise.
func (h *HashMap[K, V]) Remove(key K) bool {
	if _, ok := h.m[key]; ok {
		delete(h.m, key)
		return true
	}
	return false
}

// ContainsKey checks if a key exists in the map.
func (h *HashMap[K, V]) ContainsKey(key K) bool {
	_, ok := h.m[key]
	return ok
}

// Size returns the number of key-value pairs.
func (h *HashMap[K, V]) Size() int {
	return len(h.m)
}

// Clear removes all key-value pairs from the map.
func (h *HashMap[K, V]) Clear() {
	h.m = make(map[K]V)
}

// Keys returns a slice of all keys in the map.
// Note: Order is not guaranteed (Go maps are unordered).
func (h *HashMap[K, V]) Keys() []K {
	keys := make([]K, 0, len(h.m))
	for k := range h.m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of all values in the map.
// Note: Order is not guaranteed.
func (h *HashMap[K, V]) Values() []V {
	values := make([]V, 0, len(h.m))
	for _, v := range h.m {
		values = append(values, v)
	}
	return values
}

// MergeMaps combines two maps into a new map.
// If keys conflict, the value from map b takes precedence.
//
// This demonstrates:
// - A function with multiple type parameters
// - Working with generic interfaces
// - Type assertions on generic types
func MergeMaps[K comparable, V any](a, b Map[K, V]) Map[K, V] {
	result := NewHashMap[K, V]()

	// Merge both maps
	for _, m := range []Map[K, V]{a, b} {
		// Type assert to access the underlying map
		// This is safe because we know HashMap implements Map
		if hm, ok := m.(*HashMap[K, V]); ok {
			for k, v := range hm.m {
				result.Put(k, v)
			}
		}
	}

	return result
}

// ============================================================================
// 7. INTERFACE CONSTRAINTS
// ============================================================================
// You can use existing interfaces as type constraints, ensuring that
// generic types implement specific methods.

// Stringer is Go's standard interface for types that can be converted to string.
// It's defined in the fmt package as:
//   type Stringer interface { String() string }
//
// We can use it as a constraint to ensure our generic type can be printed.

// StringableList is a linked list that only accepts types implementing fmt.Stringer.
// This ensures every element can be converted to a string via String() method.
type StringableList[T fmt.Stringer] struct {
	head *node[T]
	size int
}

// NewStringableList creates a new list that only accepts Stringers.
func NewStringableList[T fmt.Stringer]() *StringableList[T] {
	return &StringableList[T]{}
}

// Add appends a Stringer to the list.
func (s *StringableList[T]) Add(item T) {
	n := &node[T]{value: item}
	if s.head == nil {
		s.head = n
	} else {
		curr := s.head
		for curr.next != nil {
			curr = curr.next
		}
		curr.next = n
	}
	s.size++
}

// PrintAll prints all elements using their String() method.
// This is safe because we know T implements fmt.Stringer.
func (s *StringableList[T]) PrintAll() {
	curr := s.head
	fmt.Println("StringableList contents:")
	for curr != nil {
		// We can call String() because of the fmt.Stringer constraint
		fmt.Printf("  %s\n", curr.value.String())
		curr = curr.next
	}
}

// Person is an example type that implements fmt.Stringer.
// This allows it to be used with StringableList.
type Person struct {
	Name string
	Age  int
}

// String implements fmt.Stringer for Person.
func (p Person) String() string {
	return fmt.Sprintf("%s (age %d)", p.Name, p.Age)
}

// ============================================================================
// 8. STANDARD LIBRARY GENERICS
// ============================================================================
// Go's standard library provides several packages that leverage generics
// for common operations.

// demonstrateConstraintsPackage shows the constraints package from
// golang.org/x/exp/constraints. This package provides standard constraints
// for common type characteristics.
func demonstrateConstraintsPackage() {
	fmt.Println("\n=== Constraints Package ===")

	// constraints.Ordered allows any type that supports <, <=, >, >=
	// This includes all numeric types and strings
	fmt.Println("Maximum of ints:", Maximum(42, 17))
	fmt.Println("Maximum of floats:", Maximum(3.14, 2.71))
	fmt.Println("Maximum of strings:", Maximum("zebra", "aardvark"))

	// Other useful constraints from the package:
	// - constraints.Integer: all integer types
	// - constraints.Float: all float types
	// - constraints.Complex: all complex types
	// - constraints.Signed: signed integer types
	// - constraints.Unsigned: unsigned integer types
}

// Maximum uses constraints.Ordered to work with any ordered type.
func Maximum[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// demonstrateSlicesPackage shows the slices package, which provides
// generic operations on slices.
func demonstrateSlicesPackage() {
	fmt.Println("\n=== Slices Package ===")

	// slices.Sort sorts any slice of ordered types in place
	nums := []int{5, 2, 8, 1, 9}
	fmt.Printf("Before sort: %v\n", nums)
	slices.Sort(nums)
	fmt.Printf("After sort:  %v\n", nums)

	// Works with any ordered type
	words := []string{"go", "rust", "java", "python"}
	fmt.Printf("Before sort: %v\n", words)
	slices.Sort(words)
	fmt.Printf("After sort:  %v\n", words)

	// Other useful slices functions:
	// - slices.Contains: check if slice contains element
	// - slices.Equal: compare two slices for equality
	// - slices.Index: find index of element
	// - slices.Delete: remove elements from slice
	// - slices.Clone: create a copy of slice
	contains := slices.Contains(nums, 5)
	fmt.Printf("Contains 5: %v\n", contains)

	nums2 := []int{1, 2, 5, 8, 9}
	equal := slices.Equal(nums, nums2)
	fmt.Printf("Slices equal: %v\n", equal)
}

// demonstrateMapsPackage shows the maps package, which provides
// generic operations on maps.
func demonstrateMapsPackage() {
	fmt.Println("\n=== Maps Package ===")

	// Create two maps
	m1 := map[string]int{"a": 1, "b": 2, "c": 3}
	m2 := map[string]int{"a": 1, "b": 2, "c": 3}
	m3 := map[string]int{"a": 1, "b": 99, "c": 3}

	// maps.Equal compares two maps for equality
	fmt.Printf("m1 == m2: %v\n", maps.Equal(m1, m2))
	fmt.Printf("m1 == m3: %v\n", maps.Equal(m1, m3))

	// maps.Clone creates a shallow copy
	m4 := maps.Clone(m1)
	fmt.Printf("Cloned map: %v\n", m4)

	// Other useful maps functions:
	// - maps.DeleteFunc: delete entries matching predicate
	// - maps.EqualFunc: compare with custom comparison
	// - maps.Copy: copy from one map to another
}

// demonstrateCmpPackage shows the cmp package, which provides
// generic comparison functions.
func demonstrateCmpPackage() {
	fmt.Println("\n=== Cmp Package ===")

	// cmp.Compare returns -1, 0, or 1
	fmt.Printf("Compare(5, 10): %d\n", cmp.Compare(5, 10))     // -1 (5 < 10)
	fmt.Printf("Compare(10, 5): %d\n", cmp.Compare(10, 5))     // 1 (10 > 5)
	fmt.Printf("Compare(5, 5): %d\n", cmp.Compare(5, 5))       // 0 (5 == 5)

	// Works with any ordered type
	fmt.Printf("Compare(\"apple\", \"banana\"): %d\n",
		cmp.Compare("apple", "banana")) // -1

	// Use cmp.Compare in custom sorting or comparison logic
	nums := []int{5, 2, 9, 1, 5, 6}
	maxVal, ok := MaxOfSlice(nums)
	if ok {
		fmt.Printf("Maximum value: %d\n", maxVal)
	}
}

// MaxOfSlice finds the maximum value in a slice using cmp.Compare.
// This demonstrates how to use the cmp package in custom algorithms.
func MaxOfSlice[T cmp.Ordered](s []T) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}

	max := s[0]
	for _, v := range s[1:] {
		// Use cmp.Compare instead of direct comparison
		if cmp.Compare(v, max) > 0 {
			max = v
		}
	}
	return max, true
}

// ============================================================================
// 9. DEMONSTRATION FUNCTIONS
// ============================================================================

func demonstrateBasicGenerics() {
	fmt.Println("\n=== Basic Generic LinkedList ===")

	// Create a linked list of integers
	intList := NewLinkedList[int]()
	intList.Add(10)
	intList.Add(20)
	intList.Add(30)
	fmt.Println("Integer list size:", intList.Size())
	fmt.Println("Contains 20:", intList.Contains(20))
	fmt.Println("Contains 99:", intList.Contains(99))

	// Create a linked list of strings
	strList := NewLinkedList[string]()
	strList.Add("Go")
	strList.Add("Rust")
	strList.Add("Python")
	fmt.Println("\nString list size:", strList.Size())
	fmt.Println("Contains 'Go':", strList.Contains("Go"))
	fmt.Println("Contains 'Java':", strList.Contains("Java"))

	// The same code works with both types!
	// This is the power of generics.
}

func demonstrateListOperations() {
	fmt.Println("\n=== List Operations ===")

	list := NewLinkedList[string]()
	list.Add("first")
	list.Add("second")
	list.Add("third")

	// Get operations
	fmt.Println("Initial list:")
	PrintCollection(list)

	// Insert operation
	list.Insert(1, "inserted")
	fmt.Println("\nAfter inserting at index 1:")
	PrintCollection(list)

	// Set operation
	list.Set(2, "modified")
	fmt.Println("\nAfter modifying index 2:")
	PrintCollection(list)

	// RemoveAt operation
	removed, _ := list.RemoveAt(1)
	fmt.Printf("\nAfter removing index 1 (removed: %s):\n", removed)
	PrintCollection(list)

	// Reverse operation
	ReverseList(list)
	fmt.Println("\nAfter reversing:")
	PrintCollection(list)
}

func demonstrateCustomConstraints() {
	fmt.Println("\n=== Custom Constraints ===")

	// Numeric constraint allows arithmetic operations
	intNums := []int{1, 2, 3, 4, 5}
	fmt.Printf("Sum of ints: %d\n", Sum(intNums))

	floatNums := []float64{1.1, 2.2, 3.3}
	fmt.Printf("Sum of floats: %.2f\n", Sum(floatNums))

	// Ordered constraint allows comparisons
	fmt.Printf("Max of 10 and 20: %d\n", Max(10, 20))
	fmt.Printf("Min of 10 and 20: %d\n", Min(10, 20))
	fmt.Printf("Max of strings: %s\n", Max("apple", "banana"))
}

func demonstrateMultipleTypeParameters() {
	fmt.Println("\n=== Multiple Type Parameters ===")

	// Create a map from strings to integers
	scores := NewHashMap[string, int]()
	scores.Put("Alice", 95)
	scores.Put("Bob", 87)
	scores.Put("Charlie", 92)

	fmt.Printf("Map size: %d\n", scores.Size())
	fmt.Printf("Keys: %v\n", scores.Keys())
	fmt.Printf("Values: %v\n", scores.Values())

	// Get a specific value
	if score, ok := scores.Get("Alice"); ok {
		fmt.Printf("Alice's score: %d\n", score)
	}

	// Create another map
	moreScores := NewHashMap[string, int]()
	moreScores.Put("David", 88)
	moreScores.Put("Bob", 90) // Will override Bob's score

	// Merge the maps
	merged := MergeMaps(scores, moreScores)
	fmt.Printf("\nAfter merge:\n")
	fmt.Printf("Keys: %v\n", merged.Keys())
	fmt.Printf("Values: %v\n", merged.Values())
}

func demonstrateInterfaceConstraints() {
	fmt.Println("\n=== Interface Constraints ===")

	// Create a list that only accepts types implementing fmt.Stringer
	people := NewStringableList[Person]()
	people.Add(Person{Name: "Alice", Age: 30})
	people.Add(Person{Name: "Bob", Age: 25})
	people.Add(Person{Name: "Charlie", Age: 35})

	// Print all using the String() method
	people.PrintAll()

	// This wouldn't compile because int doesn't implement fmt.Stringer:
	// numbers := NewStringableList[int]() // Compile error!
}

// ============================================================================
// MAIN FUNCTION
// ============================================================================

func main() {
	fmt.Println("Go Generics Comprehensive Reference")
	fmt.Println("====================================")

	// 1. Basic generic collections
	demonstrateBasicGenerics()

	// 2. List operations (Get, Set, Insert, Remove, Reverse)
	demonstrateListOperations()

	// 3. Custom type constraints (Numeric, Ordered)
	demonstrateCustomConstraints()

	// 4. Multiple type parameters (HashMap)
	demonstrateMultipleTypeParameters()

	// 5. Interface constraints (fmt.Stringer)
	demonstrateInterfaceConstraints()

	// 6. Standard library support
	demonstrateConstraintsPackage()
	demonstrateSlicesPackage()
	demonstrateMapsPackage()
	demonstrateCmpPackage()

	fmt.Println("\n=== Summary ===")
	fmt.Println("Generics enable:")
	fmt.Println("- Type-safe code reuse")
	fmt.Println("- Flexible data structures")
	fmt.Println("- Compile-time type checking")
	fmt.Println("- Reduced code duplication")
	fmt.Println("- Better abstraction without sacrificing safety")
}
