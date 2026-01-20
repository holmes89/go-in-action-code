# Chapter 4: Collection Types

This directory contains a comprehensive reference implementation covering arrays, slices, and maps in Go.

## Main File

**main.go** - Complete demonstration of all collection types with runnable examples

## What's Inside

The main.go file covers:

1. **Arrays** - Fixed-length collections and their characteristics
2. **Array Declarations** - Different initialization forms
3. **Array Iteration** - Range and traditional loops
4. **Arrays as Values** - Copy semantics and pointers
5. **Multidimensional Arrays** - 2D and 3D arrays
6. **Slices** - Dynamic ordered collections
7. **Slice Declarations** - Literals, make, and nil slices
8. **Append** - Growing slices dynamically
9. **Slice Expressions** - Creating slices from slices/arrays
10. **Copy** - Duplicating slice data
11. **Avoiding Panics** - Safe slice access patterns
12. **Maps** - Key-value collections
13. **Map Operations** - Common patterns and techniques
14. **Map Iteration** - Looping through maps safely

## Running the Program

Run the complete demonstration:
```bash
go run main.go
```

Each section is clearly marked and demonstrates the concepts from Chapter 4.

## Quick Reference

### Arrays

```go
// Declaration forms
var arr [5]int                    // Zero values
arr := [5]int{1, 2, 3, 4, 5}     // Full initialization
arr := [5]int{1, 2, 3}            // Partial (rest are zeros)
arr := [5]int{1: 10, 3: 30}      // Sparse initialization
arr := [...]int{1, 2, 3}          // Automatic length

// Key characteristics
// - Fixed size (part of type)
// - [5]int != [4]int
// - Assignment copies all elements
// - Pass *[N]Type to functions to avoid copies
```

### Slices

```go
// Declaration forms
slice := []int{1, 2, 3}           // Literal
slice := make([]int, 5)           // Length 5, capacity 5
slice := make([]int, 5, 10)       // Length 5, capacity 10
slice := make([]int, 0, 5)        // Length 0, ready to grow
var slice []int                   // nil slice (not usable yet)

// Operations
slice = append(slice, 4)          // Add one element
slice = append(slice, 5, 6, 7)    // Add multiple
slice = append(slice, other...)   // Append another slice

// Slice expressions
sub := slice[1:4]                 // Elements at indices 1, 2, 3
sub := slice[:3]                  // First 3 elements
sub := slice[2:]                  // From index 2 to end
sub := slice[:]                   // Entire slice
sub := slice[1:4:5]               // With capacity limit

// Copy (creates independent copy)
dest := make([]int, len(src))
copy(dest, src)

// Check length before accessing
if i < len(slice) {
    value := slice[i]
}

// Use range for safe iteration
for i, v := range slice {
    // i is index, v is value (copy)
}
```

### Maps

```go
// Declaration forms
m := map[string]int{}             // Empty map
m := make(map[string]int)         // Empty map (same)
m := map[string]int{              // With initial values
    "key1": 10,
    "key2": 20,
}
var m map[string]int              // nil map (not usable yet)

// Operations
m["key"] = value                  // Set value
value := m["key"]                 // Get value (zero if missing)
value, found := m["key"]          // Get with existence check
delete(m, "key")                  // Remove key
length := len(m)                  // Number of keys

// Idiomatic existence check
if value, found := m["key"]; found {
    // key exists
}

// Iteration (order is random!)
for key, value := range m {
    // Process key and value
}

for key := range m {
    // Just keys
}

// Using zero values
counts := make(map[string]int)
counts["word"]++                  // Zero value allows this
```

## Key Concepts

### Arrays
- **Fixed size**: Length is part of the type
- **Value semantics**: Assignment copies all elements
- **Use cases**: Fixed-size data, but slices are usually better
- **Performance**: Pass pointers to avoid copying large arrays

### Slices
- **Dynamic**: Can grow and shrink
- **Structure**: Pointer + length + capacity
- **Shared backing**: Slice expressions share the backing array
- **Growth**: append doubles capacity when small, grows 25% when large
- **Safety**: Check `len()` before accessing indices

### Slice Expressions
- `slice[low:high]` - Elements from low to high-1
- `slice[:high]` - From start to high-1
- `slice[low:]` - From low to end
- `slice[:]` - Entire slice (new slice header)
- `slice[low:high:max]` - With capacity limit

### Maps
- **Unordered**: Iteration order is random (by design)
- **Fast lookup**: O(1) average case
- **Zero values**: Missing keys return zero value
- **Nil maps**: Must initialize before use
- **Key types**: Must be comparable (support ==)

## Common Patterns

### Growing Slices Efficiently
```go
// Preallocate if you know the size
slice := make([]int, 0, expectedSize)
for i := 0; i < expectedSize; i++ {
    slice = append(slice, value)
}
```

### Independent Slice Copy
```go
// Copy to avoid shared backing array
original := []int{1, 2, 3, 4, 5}
copy := make([]int, len(original))
copy(copy, original)
```

### Safe Slice Access
```go
// Use range for safe iteration
for i, v := range slice {
    // Always safe
}

// Or check length
if i < len(slice) {
    value := slice[i]
}
```

### Map Membership Testing
```go
// Two-value retrieval
if value, found := map[key]; found {
    // key exists
}

// Simple membership (bool values)
members := map[string]bool{
    "Alice": true,
    "Bob":   true,
}
if members["Alice"] {
    // Alice is a member
}
```

### Counter Pattern
```go
// Zero value allows incrementing
counts := make(map[string]int)
for _, word := range words {
    counts[word]++ // Works even if key doesn't exist
}
```

## Common Pitfalls

❌ **Array type confusion**
```go
var a [4]int
var b [5]int
a = b // ERROR: different types
```

❌ **Forgetting to assign append**
```go
slice := []int{1, 2, 3}
append(slice, 4) // Wrong! Result is lost
```

✅ **Always assign append result**
```go
slice = append(slice, 4) // Correct
```

❌ **Modifying via range values**
```go
for _, s := range strings {
    s = s + "!" // Modifies copy, not original
}
```

✅ **Modify via index**
```go
for i := range strings {
    strings[i] = strings[i] + "!"
}
```

❌ **Adding to nil map**
```go
var m map[string]int
m["key"] = 1 // PANIC!
```

✅ **Initialize first**
```go
m := make(map[string]int)
m["key"] = 1 // OK
```

❌ **Using nil check for empty slice**
```go
if slice == nil { // Doesn't catch all empty slices
    // ...
}
```

✅ **Use len check**
```go
if len(slice) == 0 { // Correct
    // ...
}
```

## Best Practices

✅ **DO:**
- Use slices for ordered collections (not arrays)
- Preallocate slices with `make()` if you know the size
- Always assign `append()` results back
- Use `range` loops for safe iteration
- Check `len()` before accessing indices
- Use `copy()` when you need independent slices
- Initialize maps with `make()` or literals
- Use two-value retrieval to check map existence
- Use `len()` to check if collections are empty

❌ **DON'T:**
- Use arrays unless you need fixed size
- Forget to assign `append()` results
- Rely on map iteration order
- Add to nil maps
- Use `== nil` to check if slice is empty
- Access indices without checking length
- Pass large arrays by value to functions

## Testing Your Understanding

Try modifying the code to:
1. Create a function that reverses a slice in place
2. Implement a function to remove duplicates from a slice
3. Build a word frequency counter using maps
4. Create a 2D slice (slice of slices) for a grid
5. Implement a set type using maps
6. Write a function that safely gets the nth element with a default value
