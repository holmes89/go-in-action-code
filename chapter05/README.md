# Chapter 5: Working With Types

This directory contains a comprehensive reference implementation covering Go's type system through a recipe application.

## Main File

**main.go** - Complete demonstration of types, structs, methods, and interfaces with a practical recipe builder

## What's Inside

The main.go file covers:

1. **Named Types** - Creating meaningful types from primitives
2. **Structs** - Grouping related data together
3. **Type Embedding** - Extending types through composition
4. **Collection Types** - Named slices and maps
5. **Pointers** - Sharing memory and avoiding copies
6. **Methods** - Value vs pointer receivers
7. **Interfaces** - Defining and implementing behavior
8. **Type Checking** - Assertions and type switches
9. **Complete Recipe Example** - All concepts working together

## Running the Program

Run the complete demonstration:
```bash
go run main.go
```

Each section demonstrates concepts from Chapter 5 using a recipe application as the domain model.

## Domain Model Overview

The recipe application demonstrates:

```
MeasurementSystem (string)
├── Metric
└── Imperial

Magnitude (float32)    Unit (string)
└── Measurement struct ─┬── MetricWeight
                        ├── MetricVolume
                        ├── ImperialWeight
                        └── ImperialVolume

Ingredient (string)
└── IngredientMeasurement ─┬── IngredientList (slice)
                           └── Step.Ingredients

Recipe
├── Title (string)
├── Yield (uint)
└── Steps ([]Step)

RecipeBox ([]Recipe)
GroceryList (map[Ingredient]Measurement)
```

## Quick Reference

### Named Types

```go
// Create meaningful types from primitives
type Magnitude float32
type Unit string
type MeasurementSystem string

// Enum-like constants
const (
    Gram     Unit = "gram"
    KiloGram Unit = "kilogram"
    Ounce    Unit = "ounce"
)

// Benefits:
// - Self-documenting code
// - Prevents accidental misuse
// - Enables method attachment
```

### Structs

```go
// Group related data
type Measurement struct {
    Magnitude Magnitude
    Unit      Unit
}

// Create and use
m := Measurement{
    Magnitude: 200,
    Unit:      Gram,
}
fmt.Println(m.Magnitude) // Access with dot notation

// Zero values for omitted fields
m2 := Measurement{Unit: Cup} // Magnitude = 0
```

### Type Embedding

```go
// Extend types through composition
type MetricWeight struct{ Measurement }

// Create
mw := MetricWeight{
    Measurement: Measurement{200, Gram},
}

// Promoted fields (direct access)
fmt.Println(mw.Magnitude) // Not mw.Measurement.Magnitude
```

### Pointers

```go
// Value semantics (copy)
m1 := Measurement{200, Gram}
m2 := m1      // Copies m1
m2.Magnitude = 300 // Doesn't affect m1

// Pointer semantics (shared)
m3 := &Measurement{200, Gram}
m4 := m3      // Both point to same memory
m4.Magnitude = 300 // Affects m3!

// When to use pointers:
// - Large structs (avoid copying)
// - Need to modify original
// - Optional values (nil check)
```

### Methods

```go
// Value receiver (doesn't modify)
func (u Unit) String() string {
    return abbreviation(u)
}

// Pointer receiver (modifies original)
func (r *Recipe) Scale(factor uint) {
    r.Yield *= factor
    // Modify ingredient amounts...
}

// Rules:
// - Value: small types, no modification
// - Pointer: large types, modification needed
// - Be consistent per type
```

### Interfaces

```go
// Define behavior
type Convertible interface {
    ToMetric() Measurement
    ToImperial() Measurement
}

// Implement implicitly (no "implements" keyword)
func (m MetricWeight) ToMetric() Measurement {
    return m.Measurement
}

func (m MetricWeight) ToImperial() Measurement {
    return Measurement{m.Magnitude * 0.035274, Ounce}
}

// Use polymorphically
func printConversion(c Convertible) {
    metric := c.ToMetric()
    imperial := c.ToImperial()
    // Works for any type that implements Convertible
}
```

### Type Checking

```go
// Type assertion
if weight, ok := m.(MetricWeight); ok {
    // m is a MetricWeight
    fmt.Println(weight.Magnitude)
}

// Type switch
switch v := m.(type) {
case MetricWeight:
    fmt.Println("Metric weight:", v.Magnitude)
case ImperialWeight:
    fmt.Println("Imperial weight:", v.Magnitude)
default:
    fmt.Println("Unknown type")
}
```

## Key Concepts

### Type System Benefits

1. **Type Safety**: Catch errors at compile time
2. **Self-Documentation**: Code is clearer and more expressive
3. **Composition**: Build complex types from simple ones
4. **Flexibility**: Interfaces enable polymorphism
5. **Maintainability**: Changes are easier and safer

### Design Principles

**Composition over Inheritance**
```go
// Not inheritance - composition/embedding
type MetricWeight struct{ Measurement }
```

**Small Interfaces**
```go
// Good: focused, single-responsibility
type Convertible interface {
    ToMetric() Measurement
    ToImperial() Measurement
}
```

**Accept Interfaces, Return Structs**
```go
// Function accepts interface (flexible)
func Convert(m Convertible) Measurement {
    return m.ToMetric() // Concrete return type
}
```

**Make Zero Values Useful**
```go
// Zero value is usable
var m Measurement  // {0, ""}
// Fields have sensible defaults
```

### Value vs Pointer Receivers

| Criteria | Value Receiver | Pointer Receiver |
|----------|---------------|------------------|
| Modifies receiver? | No | Yes |
| Large struct? | Avoid | Prefer |
| Small struct? | Prefer | Avoid |
| Consistency | Match other methods | Match other methods |

## Common Patterns

### Enum Pattern
```go
type Status string

const (
    Pending  Status = "pending"
    Complete Status = "complete"
    Failed   Status = "failed"
)
```

### Builder Pattern
```go
recipe := Recipe{Title: "Pancakes"}
recipe.AddStep(Step{...})
recipe.AddStep(Step{...})
recipe.Scale(2)
```

### Type Assertion for Optional Behavior
```go
if stringer, ok := value.(fmt.Stringer); ok {
    fmt.Println(stringer.String())
}
```

### Collection Types
```go
type GroceryList map[Ingredient]Measurement

func (g GroceryList) Add(ing Ingredient, m Measurement) {
    // Custom behavior on collection
}
```

## Best Practices

✅ **DO:**
- Use named types to add meaning
- Make zero values useful when possible
- Use small, focused interfaces
- Be consistent with pointer/value receivers
- Compose types through embedding
- Accept interfaces, return concrete types
- Use type assertions safely (check ok)
- Document exported types and methods

❌ **DON'T:**
- Mix pointer and value receivers without reason
- Create large interfaces (keep them small)
- Use inheritance patterns (use composition)
- Ignore zero values (design around them)
- Access nil pointers (always check)
- Copy large structs (use pointers)
- Force everything into interfaces (concrete types are fine)

## Common Pitfalls

❌ **Forgetting to use pointer receiver**
```go
func (r Recipe) Scale(factor uint) {
    r.Yield *= factor // Modifies copy, not original!
}
```

✅ **Use pointer receiver for modification**
```go
func (r *Recipe) Scale(factor uint) {
    r.Yield *= factor // Modifies original
}
```

❌ **Not checking type assertion**
```go
weight := m.(MetricWeight) // Panics if wrong type!
```

✅ **Always check ok value**
```go
if weight, ok := m.(MetricWeight); ok {
    // Safe to use weight
}
```

❌ **Large interfaces**
```go
type Processor interface {
    Process()
    Validate()
    Transform()
    Save()
    Load()
    // Too many methods!
}
```

✅ **Small, focused interfaces**
```go
type Validator interface {
    Validate() error
}

type Processor interface {
    Process() error
}
```

## Testing Your Understanding

Try extending the code to:
1. Add a `Temperature` measurement type with Celsius/Fahrenheit conversion
2. Implement a `Cookable` interface with `Cook(duration time.Duration)` method
3. Create a `Cookbook` type with methods to search and filter recipes
4. Add validation methods to ensure measurements are positive
5. Implement unit conversion for volume (cups to liters, etc.)
6. Create a `String()` method for Recipe to format it nicely
7. Add nutritional information using composition

## Interface Design Tips

**The smaller the interface, the better:**
- io.Reader: `Read([]byte) (int, error)`
- io.Writer: `Write([]byte) (int, error)`
- fmt.Stringer: `String() string`

**Common single-method interfaces:**
- Easy to implement
- Easy to understand
- Easy to compose
- Maximum flexibility

**When to use interfaces:**
- Testing (mocking dependencies)
- Multiple implementations
- Decoupling packages
- Polymorphic behavior

**When NOT to use interfaces:**
- Only one implementation
- No need for abstraction
- Premature generalization
- Clear concrete type is better
