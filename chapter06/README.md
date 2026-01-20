# Chapter 6: Generics

A comprehensive reference implementation demonstrating Go's generics system, introduced in Go 1.18. This file shows how to write reusable, type-safe code using type parameters, constraints, and generic data structures.

## Quick Start

```bash
# Run the complete demonstration
go run main.go

# Run with module support (required for golang.org/x/exp/constraints)
go get golang.org/x/exp/constraints
go run main.go
```

## What's Covered

This implementation demonstrates all major concepts from Chapter 6:

### 1. **Type Parameters Basics**
- Generic node structures
- Type parameter syntax `[T any]`
- Using type parameters in struct fields
- Generic constructors

### 2. **Generic Interfaces**
- Collection[T] - basic collection operations
- List[T] - indexed access with generics
- Map[K, V] - multiple type parameters
- Interface embedding with generics

### 3. **Type Constraints**
- Built-in constraints (`any`, `comparable`)
- Custom constraints (`Numeric`, `Ordered`)
- The tilde operator (`~int` for underlying types)
- Interface constraints (`fmt.Stringer`)

### 4. **Generic Data Structures**
- LinkedList[T comparable] - generic linked list
- HashMap[K comparable, V any] - generic map
- StringableList[T fmt.Stringer] - interface-constrained list

### 5. **Generic Functions**
- Functions with type parameters
- Generic algorithms (ReverseList, MergeMaps)
- Type inference vs explicit type arguments
- Working with generic interfaces

### 6. **Multiple Type Parameters**
- Map[K, V] with independent key/value types
- Constraints on each parameter
- Real-world use cases (dictionaries, pairs)

### 7. **Standard Library Support**
- constraints package (Ordered, Integer, Float)
- slices package (Sort, Contains, Equal)
- maps package (Equal, Clone, Copy)
- cmp package (Compare, Ordered)

## Key Concepts

### Type Parameter Syntax

```go
// Single type parameter
type List[T any] interface { ... }

// With constraint
type LinkedList[T comparable] struct { ... }

// Multiple type parameters
type Map[K comparable, V any] interface { ... }

// Generic function
func ReverseList[T comparable](list List[T]) { ... }
```

### Constraints Hierarchy

```
any (interface{})
├── comparable (supports == and !=)
├── constraints.Ordered (supports <, >, <=, >=)
│   ├── Integer types
│   ├── Float types
│   └── string
└── Custom interfaces (fmt.Stringer, etc.)
```

### When to Use Each Constraint

| Constraint | Use When | Example |
|------------|----------|---------|
| `any` | No special requirements | `type Box[T any]` |
| `comparable` | Need equality checks | `Contains()`, `Remove()` |
| `constraints.Ordered` | Need comparisons | `Sort()`, `Max()`, `Min()` |
| `constraints.Integer` | Integer-specific ops | Bit operations |
| Custom interface | Need specific methods | `fmt.Stringer` for printing |

### The Tilde (`~`) Operator

```go
// Without tilde - only exact types
type Numeric interface {
    int | float64  // Only built-in int and float64
}

// With tilde - underlying types too
type NumericFlexible interface {
    ~int | ~float64  // int, float64, AND user-defined types like MyInt
}

type MyInt int  // Satisfies NumericFlexible, not Numeric
```

## Code Organization

```
main.go
├── Generic Structures (node, LinkedList)
├── Generic Interfaces (Collection, List, Map)
├── LinkedList Implementation (Collection + List methods)
├── HashMap Implementation (Map methods)
├── Generic Functions (ReverseList, MergeMaps, etc.)
├── Custom Constraints (Numeric, Ordered)
├── Interface Constraints (StringableList)
├── Standard Library Demonstrations
└── Main Function (runs all examples)
```

## Demonstration Functions

Each demonstration is self-contained and shows specific aspects:

1. **demonstrateBasicGenerics()** - LinkedList with different types
2. **demonstrateListOperations()** - Get, Set, Insert, Remove, Reverse
3. **demonstrateCustomConstraints()** - Sum, Max, Min with constraints
4. **demonstrateMultipleTypeParameters()** - HashMap and MergeMaps
5. **demonstrateInterfaceConstraints()** - StringableList with fmt.Stringer
6. **demonstrateConstraintsPackage()** - constraints.Ordered
7. **demonstrateSlicesPackage()** - Generic slice operations
8. **demonstrateMapsPackage()** - Generic map operations
9. **demonstrateCmpPackage()** - Generic comparisons

## Design Patterns

### Generic Constructor Pattern

```go
func NewLinkedList[T comparable]() *LinkedList[T] {
    return &LinkedList[T]{}
}

// Usage - type parameter specified explicitly
list := NewLinkedList[int]()
```

### Zero Value Pattern

```go
func Get(index int) (T, error) {
    var zero T  // Zero value of unknown type T
    if index < 0 {
        return zero, errors.New("invalid index")
    }
    // ... actual implementation
    return value, nil
}
```

### Type Assertion with Generics

```go
func PrintCollection[T any](c Collection[T]) {
    // Check if collection also implements List
    if list, ok := c.(List[T]); ok {
        // Can now use List methods
        val, _ := list.Get(0)
    }
}
```

### Multiple Type Parameter Pattern

```go
// Keys must be comparable, values can be anything
type Map[K comparable, V any] interface {
    Put(key K, value V)
    Get(key K) (V, bool)
}
```

## Common Pitfalls

### 1. Forgetting Type Parameters in Method Receivers

```go
// ❌ Wrong - missing type parameter
func (l *LinkedList) Add(item T) { ... }

// ✅ Correct
func (l *LinkedList[T]) Add(item T) { ... }
```

### 2. Constraint Mismatch

```go
// ❌ Won't compile - needs comparable
list := NewLinkedList[[]int]()  // Slices aren't comparable

// ✅ Correct
list := NewLinkedList[int]()
```

### 3. Type Assertion with Wrong Type Parameters

```go
var c Collection[int] = NewLinkedList[int]()
list, ok := c.(List[string])  // ❌ Type mismatch - ok will be false
list, ok := c.(List[int])     // ✅ Correct
```

### 4. Using Operators Without Constraints

```go
// ❌ Won't compile
func Add[T any](a, b T) T {
    return a + b  // Error: + not defined for all types
}

// ✅ Correct - constrain to numeric types
func Add[T Numeric](a, b T) T {
    return a + b  // Now safe
}
```

### 5. Unnecessary Type Parameters

```go
// ❌ Overusing generics
func PrintString[T string](s T) { fmt.Println(s) }

// ✅ Just use concrete types when appropriate
func PrintString(s string) { fmt.Println(s) }
```

## Performance Considerations

### When Generics May Be Slower

- **Interface dispatch overhead** - Generic code may use interface-based dispatch
- **Reduced optimization** - Compiler optimizations may be less aggressive
- **Allocation patterns** - Generic code may allocate more frequently

### When Generics Are Just as Fast

- **Simple operations** - Basic operations compile efficiently
- **Proper constraints** - Specific constraints enable better optimization
- **Non-critical paths** - Most application code benefits from generics

### Best Practices

1. **Profile before optimizing** - Measure actual performance impact
2. **Use concrete types in hot paths** - Consider concrete implementations for tight loops
3. **Leverage constraints** - More specific constraints = better optimization
4. **Avoid unnecessary boxing** - Be aware of interface conversions

## Standard Library Integration

### Using constraints Package

```go
import "golang.org/x/exp/constraints"

func Maximum[T constraints.Ordered](a, b T) T {
    if a > b { return a }
    return b
}
```

### Using slices Package

```go
import "slices"

nums := []int{5, 2, 8, 1}
slices.Sort(nums)               // Generic sort
contains := slices.Contains(nums, 5)
equal := slices.Equal(nums, other)
```

### Using maps Package

```go
import "maps"

m1 := map[string]int{"a": 1}
m2 := maps.Clone(m1)            // Generic clone
equal := maps.Equal(m1, m2)     // Generic comparison
```

### Using cmp Package

```go
import "cmp"

result := cmp.Compare(5, 10)    // Returns -1, 0, or 1
// Use in custom comparison logic
```

## Real-World Use Cases

### 1. **Data Structures**
- Generic collections (List, Set, Map, Queue, Stack)
- Type-safe wrappers around built-in types
- Custom containers with specific behaviors

### 2. **Algorithms**
- Sorting, searching, filtering
- Mathematical operations
- Graph algorithms, tree traversals

### 3. **Utility Functions**
- Max, Min, Clamp for any ordered type
- Map, Filter, Reduce for slices
- Generic validation functions

### 4. **API Design**
- Type-safe builder patterns
- Generic options/configuration
- Reusable middleware

### 5. **Testing Utilities**
- Generic assertion helpers
- Mock generators
- Test data builders

## Further Exploration

Try extending the examples:

1. **Implement Set[T comparable]** - Generic set data structure
2. **Add OrderedMap[K, V]** - Map that maintains insertion order
3. **Create Stack[T] and Queue[T]** - Generic FIFO/LIFO structures
4. **Build BinaryTree[T Ordered]** - Generic binary search tree
5. **Write Filter, Map, Reduce** - Generic slice operations
6. **Implement PriorityQueue[T]** - Generic heap-based queue

## Key Takeaways

✅ **Generics eliminate code duplication** - Write once, use with any type

✅ **Type safety at compile time** - Catch errors before runtime

✅ **Constraints enable powerful abstractions** - Express requirements precisely

✅ **Standard library integration** - Use generic utilities for common tasks

✅ **Backward compatible** - Works alongside existing non-generic code

✅ **Performance is usually acceptable** - Profile before optimizing

⚠️ **Use judiciously** - Don't over-engineer simple code

⚠️ **Start with concrete types** - Refactor to generics when patterns emerge

⚠️ **Document constraints clearly** - Make requirements explicit

## Testing

Run the code to see all demonstrations:

```bash
go run main.go
```

Expected output includes:
- LinkedList operations with different types
- Custom constraint demonstrations (Sum, Max, Min)
- HashMap with multiple type parameters
- Interface constraint examples
- Standard library generic functions

## Dependencies

```bash
# Required for constraints package
go get golang.org/x/exp/constraints
```

## Summary

Generics in Go provide a powerful way to write reusable, type-safe code without sacrificing clarity or performance. This reference demonstrates the key patterns and best practices for effective generic programming in Go, from simple type parameters to complex multi-parameter abstractions with custom constraints.
