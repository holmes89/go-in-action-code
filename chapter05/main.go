package main

import (
	"fmt"
)

// Chapter 5: Working With Types
// This file demonstrates Go's type system through a recipe application.
// Topics: named types, structs, embedding, methods, interfaces, pointers, type checking

func main() {
	fmt.Println("=== CHAPTER 5: WORKING WITH TYPES ===\n")

	// Run each demonstration
	demonstrateNamedTypes()
	demonstrateStructs()
	demonstrateTypeEmbedding()
	demonstrateCollectionTypes()
	demonstratePointers()
	demonstrateMethods()
	demonstrateInterfaces()
	demonstrateTypeChecking()
	demonstrateCompleteRecipe()
}

// ==============================================================================
// NAMED TYPES
// ==============================================================================
// Named types create new types from existing primitives, adding meaning and safety.
// They make code self-documenting and prevent accidental misuse.

// MeasurementSystem represents a system of measurement (metric or imperial)
type MeasurementSystem string

// Magnitude represents the numeric value of a measurement
type Magnitude float32

// Unit represents the unit of measurement (grams, cups, etc.)
type Unit string

// Define measurement systems as constants
const (
	Metric   MeasurementSystem = "metric"
	Imperial MeasurementSystem = "imperial"
)

// Define common units as constants (enum-like pattern in Go)
const (
	// Weight units
	Gram     Unit = "gram"
	KiloGram Unit = "kilogram"
	Ounce    Unit = "ounce"
	Pound    Unit = "pound"

	// Volume units
	Liter  Unit = "liter"
	Cup    Unit = "cup"
	Pint   Unit = "pint"
	Quart  Unit = "quart"
	Gallon Unit = "gallon"
)

func demonstrateNamedTypes() {
	fmt.Println("--- NAMED TYPES ---")

	// Named types give meaning to primitive values
	var weight Magnitude = 200
	var unit Unit = Gram
	var system MeasurementSystem = Metric

	fmt.Printf("Measurement: %.2f %s (%s system)\n", weight, unit, system)

	// Using constants for type-safe enum-like values
	fmt.Println("\nCommon units:")
	fmt.Printf("  Weight: %s, %s, %s, %s\n", Gram, KiloGram, Ounce, Pound)
	fmt.Printf("  Volume: %s, %s, %s, %s\n", Liter, Cup, Pint, Quart)

	// Named types prevent accidental mixing
	var metric MeasurementSystem = Metric
	var imperial MeasurementSystem = Imperial
	fmt.Printf("\nSystems: %s and %s\n", metric, imperial)

	// Type safety: can't accidentally use wrong type
	// var badUnit Unit = 123 // ERROR: cannot use 123 (int) as Unit value

	fmt.Println()
}

// ==============================================================================
// STRUCTS - GROUPING RELATED DATA
// ==============================================================================
// Structs bundle related values together into a single entity.

// Measurement combines magnitude and unit into a single type
type Measurement struct {
	Magnitude Magnitude
	Unit      Unit
}

func demonstrateStructs() {
	fmt.Println("--- STRUCTS ---")

	// Create a struct with field names (recommended)
	m1 := Measurement{
		Magnitude: 200,
		Unit:      Gram,
	}
	fmt.Printf("Measurement 1: %.2f %s\n", m1.Magnitude, m1.Unit)

	// Access fields with dot notation
	fmt.Printf("Just the magnitude: %.2f\n", m1.Magnitude)
	fmt.Printf("Just the unit: %s\n", m1.Unit)

	// Modify fields
	m1.Magnitude = 250
	fmt.Printf("After modification: %.2f %s\n", m1.Magnitude, m1.Unit)

	// Omitted fields get zero values
	m2 := Measurement{
		Unit: Cup,
		// Magnitude omitted - defaults to 0
	}
	fmt.Printf("Measurement 2: %.2f %s\n", m2.Magnitude, m2.Unit)

	// Empty struct (all zero values)
	var m3 Measurement
	fmt.Printf("Empty measurement: %.2f %s\n", m3.Magnitude, m3.Unit)

	// Structs can be compared if all fields are comparable
	m4 := Measurement{Magnitude: 200, Unit: Gram}
	m5 := Measurement{Magnitude: 200, Unit: Gram}
	fmt.Printf("\nAre m4 and m5 equal? %v\n", m4 == m5) // true

	// IMPORTANT: Uppercase = exported (public), lowercase = unexported (private)
	// Our Measurement fields are uppercase, so they're accessible outside the package

	fmt.Println()
}

// ==============================================================================
// TYPE EMBEDDING - EXTENDING TYPES
// ==============================================================================
// Go uses composition (embedding) instead of inheritance to extend types.

// Specialized measurement types that embed Measurement
type MetricWeight struct{ Measurement }
type MetricVolume struct{ Measurement }
type ImperialWeight struct{ Measurement }
type ImperialVolume struct{ Measurement }

func demonstrateTypeEmbedding() {
	fmt.Println("--- TYPE EMBEDDING ---")

	// Create a MetricWeight (embeds Measurement)
	mw := MetricWeight{
		Measurement: Measurement{
			Magnitude: 200,
			Unit:      Gram,
		},
	}
	fmt.Printf("MetricWeight: %.2f %s\n", mw.Magnitude, mw.Unit)

	// Embedded fields are "promoted" - can access directly
	fmt.Printf("Direct access to embedded fields: %.2f %s\n", mw.Magnitude, mw.Unit)

	// Create an ImperialVolume
	iv := ImperialVolume{
		Measurement: Measurement{
			Magnitude: 2,
			Unit:      Quart,
		},
	}
	fmt.Printf("ImperialVolume: %.2f %s\n", iv.Magnitude, iv.Unit)

	// Each type is distinct even though they embed the same struct
	fmt.Printf("\nType of mw: %T\n", mw)
	fmt.Printf("Type of iv: %T\n", iv)

	// Can also set just the embedded fields
	mv := MetricVolume{
		Measurement: Measurement{
			Unit: Liter,
			// Magnitude omitted - defaults to 0
		},
	}
	fmt.Printf("MetricVolume with zero magnitude: %.2f %s\n", mv.Magnitude, mv.Unit)

	fmt.Println()
}

// ==============================================================================
// COLLECTION TYPES - SLICES AND MAPS
// ==============================================================================
// Named collection types add meaning to slices and maps.

// Ingredient represents an ingredient name
type Ingredient string

// IngredientMeasurement pairs an ingredient with its measurement
type IngredientMeasurement struct {
	Ingredient  Ingredient
	Measurement Convertible // Can be any type that implements Convertible
}

// IngredientList is a slice of ingredient measurements
type IngredientList []IngredientMeasurement

// Step represents one step in a recipe
type Step struct {
	Description string
	Ingredients IngredientList
}

// Recipe represents a complete recipe
type Recipe struct {
	Title string
	Yield uint
	Steps []Step
}

// RecipeBox is a collection of recipes
type RecipeBox []Recipe

// GroceryList maps ingredients to their total measurements
type GroceryList map[Ingredient]Measurement

func demonstrateCollectionTypes() {
	fmt.Println("--- COLLECTION TYPES ---")

	// Create an ingredient list
	ingredients := IngredientList{
		{Ingredient: "Flour", Measurement: MetricWeight{Measurement{200, Gram}}},
		{Ingredient: "Sugar", Measurement: MetricWeight{Measurement{50, Gram}}},
		{Ingredient: "Milk", Measurement: MetricVolume{Measurement{300, Liter}}},
	}

	fmt.Println("Ingredient List:")
	for i, ing := range ingredients {
		fmt.Printf("  %d. %s: %.2f %s\n", i+1, ing.Ingredient, 
			ing.Measurement.ToMetric().Magnitude, ing.Measurement.ToMetric().Unit)
	}

	// Create a recipe
	recipe := Recipe{
		Title: "Simple Pancakes",
		Yield: 4,
		Steps: []Step{
			{
				Description: "Mix dry ingredients",
				Ingredients: IngredientList{
					{Ingredient: "Flour", Measurement: MetricWeight{Measurement{200, Gram}}},
					{Ingredient: "Sugar", Measurement: MetricWeight{Measurement{50, Gram}}},
				},
			},
			{
				Description: "Add wet ingredients",
				Ingredients: IngredientList{
					{Ingredient: "Milk", Measurement: MetricVolume{Measurement{300, Liter}}},
					{Ingredient: "Egg", Measurement: MetricWeight{Measurement{50, Gram}}},
				},
			},
		},
	}

	fmt.Printf("\nRecipe: %s (serves %d)\n", recipe.Title, recipe.Yield)
	for i, step := range recipe.Steps {
		fmt.Printf("  Step %d: %s\n", i+1, step.Description)
	}

	// Create a RecipeBox (slice of recipes)
	box := RecipeBox{recipe}
	fmt.Printf("\nRecipe box contains %d recipe(s)\n", len(box))

	// Create a GroceryList (map)
	groceries := make(GroceryList)
	groceries["Flour"] = Measurement{200, Gram}
	groceries["Sugar"] = Measurement{50, Gram}
	groceries["Milk"] = Measurement{300, Liter}

	fmt.Println("\nGrocery List:")
	for ingredient, measurement := range groceries {
		fmt.Printf("  %s: %.2f %s\n", ingredient, measurement.Magnitude, measurement.Unit)
	}

	fmt.Println()
}

// ==============================================================================
// POINTERS - SHARING MEMORY
// ==============================================================================
// Pointers allow multiple variables to refer to the same memory location.

func demonstratePointers() {
	fmt.Println("--- POINTERS ---")

	// Copy by value (default in Go)
	m1 := Measurement{200, Gram}
	m2 := m1 // Creates a COPY
	m2.Magnitude = 300

	fmt.Println("Copy by value:")
	fmt.Printf("  m1: %.2f %s\n", m1.Magnitude, m1.Unit) // Still 200
	fmt.Printf("  m2: %.2f %s\n", m2.Magnitude, m2.Unit) // Now 300

	// Using pointers to share memory
	m3 := Measurement{200, Gram}
	m4 := &m3 // m4 is a pointer to m3
	m4.Magnitude = 300 // Modifies m3!

	fmt.Println("\nUsing pointers:")
	fmt.Printf("  m3: %.2f %s\n", m3.Magnitude, m3.Unit) // Now 300
	fmt.Printf("  *m4: %.2f %s\n", m4.Magnitude, m4.Unit) // Also 300

	// Go automatically dereferences pointers for struct field access
	// m4.Magnitude is the same as (*m4).Magnitude

	// Creating pointers with new
	m5 := new(Measurement)
	m5.Magnitude = 100
	m5.Unit = Ounce
	fmt.Printf("\nCreated with new: %.2f %s\n", m5.Magnitude, m5.Unit)

	// Pointers with struct literals (common pattern)
	m6 := &Measurement{
		Magnitude: 500,
		Unit:      Gram,
	}
	fmt.Printf("Pointer from literal: %.2f %s\n", m6.Magnitude, m6.Unit)

	// Zero value of pointer is nil
	var m7 *Measurement
	fmt.Printf("\nNil pointer: %v\n", m7)
	if m7 == nil {
		fmt.Println("  m7 is nil (doesn't point to anything)")
	}

	// Pointers are useful for large structs (avoid copying)
	// and when you need to modify the original value

	fmt.Println()
}

// ==============================================================================
// METHODS - ADDING BEHAVIOR TO TYPES
// ==============================================================================
// Methods attach behavior to types. The receiver can be a value or pointer.

// String implements the Stringer interface for Unit
// VALUE RECEIVER - doesn't modify the receiver
func (u Unit) String() string {
	switch u {
	case Gram:
		return "g"
	case KiloGram:
		return "kg"
	case Ounce:
		return "oz"
	case Pound:
		return "lb"
	case Liter:
		return "L"
	case Cup:
		return "cup"
	case Pint:
		return "pt"
	case Quart:
		return "qt"
	case Gallon:
		return "gal"
	}
	return string(u)
}

// AddStep adds a step to a recipe
// POINTER RECEIVER - modifies the receiver
func (r *Recipe) AddStep(step Step) {
	r.Steps = append(r.Steps, step)
}

// Scale multiplies ingredient amounts by a factor
// POINTER RECEIVER - modifies the recipe in place
func (r *Recipe) Scale(factor uint) {
	for i := range r.Steps {
		for j := range r.Steps[i].Ingredients {
			r.Steps[i].Ingredients[j].Measurement = scaleMeasurement(
				r.Steps[i].Ingredients[j].Measurement,
				float32(factor),
			)
		}
	}
	r.Yield *= factor
}

// Helper function for scaling measurements
func scaleMeasurement(m Convertible, factor float32) Convertible {
	switch v := m.(type) {
	case MetricWeight:
		return MetricWeight{Measurement{v.Magnitude * Magnitude(factor), v.Unit}}
	case MetricVolume:
		return MetricVolume{Measurement{v.Magnitude * Magnitude(factor), v.Unit}}
	case ImperialWeight:
		return ImperialWeight{Measurement{v.Magnitude * Magnitude(factor), v.Unit}}
	case ImperialVolume:
		return ImperialVolume{Measurement{v.Magnitude * Magnitude(factor), v.Unit}}
	default:
		return m
	}
}

func demonstrateMethods() {
	fmt.Println("--- METHODS ---")

	// Value receiver example: String() method
	unit := Gram
	fmt.Printf("Unit: %s (full name: %s)\n", unit, unit) // Uses String() method
	fmt.Printf("Abbreviated: %s\n", unit.String())

	// Pointer receiver example: AddStep() and Scale()
	recipe := Recipe{
		Title: "Test Recipe",
		Yield: 2,
	}

	fmt.Printf("\nOriginal recipe: %s (serves %d)\n", recipe.Title, recipe.Yield)

	// Add steps using pointer receiver method
	recipe.AddStep(Step{
		Description: "Mix ingredients",
		Ingredients: IngredientList{
			{Ingredient: "Flour", Measurement: MetricWeight{Measurement{100, Gram}}},
		},
	})
	recipe.AddStep(Step{
		Description: "Bake",
		Ingredients: IngredientList{},
	})

	fmt.Printf("After adding steps: %d steps\n", len(recipe.Steps))

	// Scale the recipe
	recipe.Scale(3)
	fmt.Printf("After scaling by 3: serves %d\n", recipe.Yield)
	fmt.Printf("  Flour amount: %.2f %s\n",
		recipe.Steps[0].Ingredients[0].Measurement.ToMetric().Magnitude,
		recipe.Steps[0].Ingredients[0].Measurement.ToMetric().Unit)

	// VALUE vs POINTER RECEIVER RULES:
	// - Use value receiver when method doesn't modify the receiver
	// - Use pointer receiver when method needs to modify the receiver
	// - Use pointer receiver for large structs (avoid copying)
	// - Be consistent: if one method uses pointer receiver, all should

	fmt.Println()
}

// ==============================================================================
// INTERFACES - DEFINING BEHAVIOR
// ==============================================================================
// Interfaces define a set of methods that types must implement.
// Go uses "duck typing" - if it walks like a duck, it's a duck!

// Convertible interface defines types that can convert between systems
type Convertible interface {
	ToImperial() Measurement
	ToMetric() Measurement
}

// Implement Convertible for MetricWeight
func (m MetricWeight) ToImperial() Measurement {
	// 1 gram ≈ 0.035274 ounces
	return Measurement{
		Magnitude: m.Magnitude * 0.035274,
		Unit:      Ounce,
	}
}

func (m MetricWeight) ToMetric() Measurement {
	return m.Measurement
}

// Implement Convertible for ImperialWeight
func (i ImperialWeight) ToMetric() Measurement {
	// 1 ounce ≈ 28.3495 grams
	return Measurement{
		Magnitude: i.Magnitude * 28.3495,
		Unit:      Gram,
	}
}

func (i ImperialWeight) ToImperial() Measurement {
	return i.Measurement
}

// Implement Convertible for MetricVolume
func (m MetricVolume) ToImperial() Measurement {
	// 1 liter ≈ 2.11338 pints
	return Measurement{
		Magnitude: m.Magnitude * 2.11338,
		Unit:      Pint,
	}
}

func (m MetricVolume) ToMetric() Measurement {
	return m.Measurement
}

// Implement Convertible for ImperialVolume
func (i ImperialVolume) ToMetric() Measurement {
	// 1 pint ≈ 0.473176 liters
	return Measurement{
		Magnitude: i.Magnitude * 0.473176,
		Unit:      Liter,
	}
}

func (i ImperialVolume) ToImperial() Measurement {
	return i.Measurement
}

// Function that accepts any Convertible type
func printConversion(c Convertible, name string) {
	metric := c.ToMetric()
	imperial := c.ToImperial()
	fmt.Printf("%s:\n", name)
	fmt.Printf("  Metric:   %.2f %s\n", metric.Magnitude, metric.Unit)
	fmt.Printf("  Imperial: %.2f %s\n", imperial.Magnitude, imperial.Unit)
}

func demonstrateInterfaces() {
	fmt.Println("--- INTERFACES ---")

	// Create different measurement types
	mw := MetricWeight{Measurement{200, Gram}}
	iw := ImperialWeight{Measurement{8, Ounce}}
	mv := MetricVolume{Measurement{1, Liter}}
	iv := ImperialVolume{Measurement{2, Pint}}

	// All these types implement Convertible
	// We can pass them to functions that expect Convertible
	printConversion(mw, "Metric Weight")
	fmt.Println()
	printConversion(iw, "Imperial Weight")
	fmt.Println()
	printConversion(mv, "Metric Volume")
	fmt.Println()
	printConversion(iv, "Imperial Volume")

	// Store different types in a slice of Convertible
	fmt.Println("\nMixed measurements (all Convertible):")
	measurements := []Convertible{mw, iw, mv, iv}
	for i, m := range measurements {
		metric := m.ToMetric()
		fmt.Printf("  %d. %.2f %s\n", i+1, metric.Magnitude, metric.Unit)
	}

	// INTERFACE BENEFITS:
	// - Polymorphism: different types, same interface
	// - Flexibility: easy to add new types that satisfy the interface
	// - Testing: easy to mock interfaces
	// - Decoupling: code depends on behavior, not concrete types

	fmt.Println()
}

// ==============================================================================
// TYPE CHECKING - TYPE ASSERTIONS AND TYPE SWITCHES
// ==============================================================================
// Determine the concrete type behind an interface at runtime.

// MeasurementSystemFromConvertible returns the measurement system
func MeasurementSystemFromConvertible(m Convertible) MeasurementSystem {
	switch m.(type) {
	case MetricWeight, MetricVolume:
		return Metric
	case ImperialWeight, ImperialVolume:
		return Imperial
	default:
		return ""
	}
}

// Type assertion example
func isMetricWeight(m Convertible) bool {
	_, ok := m.(MetricWeight)
	return ok
}

// Type switch with actions for each type
func describeMeasurement(m Convertible) string {
	switch v := m.(type) {
	case MetricWeight:
		return fmt.Sprintf("Metric weight: %.2f %s", v.Magnitude, v.Unit)
	case MetricVolume:
		return fmt.Sprintf("Metric volume: %.2f %s", v.Magnitude, v.Unit)
	case ImperialWeight:
		return fmt.Sprintf("Imperial weight: %.2f %s", v.Magnitude, v.Unit)
	case ImperialVolume:
		return fmt.Sprintf("Imperial volume: %.2f %s", v.Magnitude, v.Unit)
	default:
		return "Unknown measurement type"
	}
}

func demonstrateTypeChecking() {
	fmt.Println("--- TYPE CHECKING ---")

	// Create some measurements
	mw := MetricWeight{Measurement{200, Gram}}
	iw := ImperialWeight{Measurement{8, Ounce}}
	mv := MetricVolume{Measurement{1, Liter}}

	// Type switch to determine system
	fmt.Println("Measurement systems:")
	fmt.Printf("  mw: %s\n", MeasurementSystemFromConvertible(mw))
	fmt.Printf("  iw: %s\n", MeasurementSystemFromConvertible(iw))
	fmt.Printf("  mv: %s\n", MeasurementSystemFromConvertible(mv))

	// Type assertion to check specific type
	fmt.Println("\nType assertion checks:")
	fmt.Printf("  Is mw a MetricWeight? %v\n", isMetricWeight(mw))
	fmt.Printf("  Is iw a MetricWeight? %v\n", isMetricWeight(iw))

	// Type switch with descriptions
	fmt.Println("\nType switch descriptions:")
	fmt.Printf("  %s\n", describeMeasurement(mw))
	fmt.Printf("  %s\n", describeMeasurement(iw))
	fmt.Printf("  %s\n", describeMeasurement(mv))

	// Type assertion with value extraction (need interface value)
	fmt.Println("\nExtracting concrete values from interface:")
	var c Convertible = mw // Store in interface
	if weight, ok := c.(MetricWeight); ok {
		fmt.Printf("  MetricWeight magnitude: %.2f\n", weight.Magnitude)
	}

	// Safe type assertion (check ok before using)
	c = iw // Imperial weight
	if weight, ok := c.(MetricWeight); ok {
		fmt.Printf("  Would print if iw was MetricWeight: %.2f\n", weight.Magnitude)
	} else {
		fmt.Println("  iw is not a MetricWeight")
	}

	// TYPE CHECKING USE CASES:
	// - Determine concrete type for type-specific logic
	// - Extract additional information from concrete types
	// - Handle different types differently in generic code
	// - Implement type-specific optimizations

	fmt.Println()
}

// ==============================================================================
// COMPLETE EXAMPLE - PUTTING IT ALL TOGETHER
// ==============================================================================
// Demonstrates how all the concepts work together in practice.

// CreateGroceryList generates a shopping list from recipes
func CreateGroceryList(recipes RecipeBox) GroceryList {
	shoppingList := make(GroceryList)

	for _, recipe := range recipes {
		for _, step := range recipe.Steps {
			for _, im := range step.Ingredients {
				metric := im.Measurement.ToMetric()

				if existing, found := shoppingList[im.Ingredient]; found {
					// Add to existing amount
					shoppingList[im.Ingredient] = Measurement{
						Magnitude: existing.Magnitude + metric.Magnitude,
						Unit:      metric.Unit,
					}
				} else {
					// New ingredient
					shoppingList[im.Ingredient] = metric
				}
			}
		}
	}

	return shoppingList
}

// ConvertToImperial converts all measurements in a recipe to imperial
func (r *Recipe) ConvertToImperial() {
	for i := range r.Steps {
		for j := range r.Steps[i].Ingredients {
			imperial := r.Steps[i].Ingredients[j].Measurement.ToImperial()
			// Wrap in appropriate type based on unit
			if imperial.Unit == Ounce || imperial.Unit == Pound {
				r.Steps[i].Ingredients[j].Measurement = ImperialWeight{imperial}
			} else {
				r.Steps[i].Ingredients[j].Measurement = ImperialVolume{imperial}
			}
		}
	}
}

func demonstrateCompleteRecipe() {
	fmt.Println("--- COMPLETE RECIPE EXAMPLE ---")

	// Create a recipe
	recipe := Recipe{
		Title: "Classic Pancakes",
		Yield: 4,
	}

	// Add steps
	recipe.AddStep(Step{
		Description: "Mix dry ingredients",
		Ingredients: IngredientList{
			{Ingredient: "Flour", Measurement: MetricWeight{Measurement{200, Gram}}},
			{Ingredient: "Sugar", Measurement: MetricWeight{Measurement{50, Gram}}},
			{Ingredient: "Baking Powder", Measurement: MetricWeight{Measurement{10, Gram}}},
		},
	})

	recipe.AddStep(Step{
		Description: "Add wet ingredients and mix",
		Ingredients: IngredientList{
			{Ingredient: "Milk", Measurement: MetricVolume{Measurement{0.3, Liter}}},
			{Ingredient: "Egg", Measurement: MetricWeight{Measurement{50, Gram}}},
			{Ingredient: "Butter", Measurement: MetricWeight{Measurement{30, Gram}}},
		},
	})

	recipe.AddStep(Step{
		Description: "Cook on griddle until golden",
		Ingredients: IngredientList{},
	})

	// Print original recipe
	fmt.Printf("Recipe: %s (serves %d)\n\n", recipe.Title, recipe.Yield)
	for i, step := range recipe.Steps {
		fmt.Printf("Step %d: %s\n", i+1, step.Description)
		if len(step.Ingredients) > 0 {
			fmt.Println("  Ingredients:")
			for _, ing := range step.Ingredients {
				metric := ing.Measurement.ToMetric()
				fmt.Printf("    - %s: %.2f %s\n", ing.Ingredient, metric.Magnitude, metric.Unit)
			}
		}
	}

	// Scale the recipe
	fmt.Println("\n--- Scaling Recipe (x2) ---")
	recipe.Scale(2)
	fmt.Printf("Now serves: %d\n", recipe.Yield)
	fmt.Println("Updated ingredient amounts:")
	for _, step := range recipe.Steps {
		for _, ing := range step.Ingredients {
			metric := ing.Measurement.ToMetric()
			fmt.Printf("  %s: %.2f %s\n", ing.Ingredient, metric.Magnitude, metric.Unit)
		}
	}

	// Convert to imperial
	fmt.Println("\n--- Converting to Imperial ---")
	recipe.ConvertToImperial()
	for _, step := range recipe.Steps {
		for _, ing := range step.Ingredients {
			imperial := ing.Measurement.ToImperial()
			fmt.Printf("  %s: %.2f %s\n", ing.Ingredient, imperial.Magnitude, imperial.Unit)
		}
	}

	// Create grocery list
	fmt.Println("\n--- Grocery List ---")
	box := RecipeBox{recipe}
	groceries := CreateGroceryList(box)
	for ingredient, measurement := range groceries {
		fmt.Printf("  %s: %.2f %s\n", ingredient, measurement.Magnitude, measurement.Unit)
	}

	fmt.Println()
}

// ==============================================================================
// KEY CONCEPTS SUMMARY
// ==============================================================================
//
// NAMED TYPES
// - Create new types from primitives: type MyInt int
// - Add meaning and prevent accidental misuse
// - Use constants for enum-like values
// - Makes code self-documenting
//
// STRUCTS
// - Group related data: type Person struct { Name string; Age int }
// - Access fields with dot notation: person.Name
// - Zero values for omitted fields
// - Uppercase fields = exported (public)
// - Lowercase fields = unexported (private)
//
// TYPE EMBEDDING
// - Composition over inheritance
// - Embed types without field name: type A struct { B }
// - Embedded fields are "promoted" (direct access)
// - Each type is distinct even with same embedded type
//
// COLLECTION TYPES
// - Named slices: type MyList []Item
// - Named maps: type MyMap map[string]int
// - Add methods to collection types
// - Makes collections more meaningful
//
// POINTERS
// - Store memory addresses: *Type
// - Get address: &value
// - Dereference: *pointer
// - Auto-dereference for struct fields
// - Zero value is nil
// - Avoid copying large structs
// - Enable modification of original value
//
// METHODS
// - Value receiver: func (v Type) Method()
//   * Operates on a copy
//   * Use when not modifying receiver
// - Pointer receiver: func (v *Type) Method()
//   * Operates on original
//   * Use when modifying receiver or for large structs
// - Be consistent in receiver choice
//
// INTERFACES
// - Define behavior: type I interface { Method() }
// - Implemented implicitly (duck typing)
// - Any type with required methods satisfies interface
// - Enable polymorphism and flexibility
// - Small interfaces are better (1-3 methods)
// - Code to interfaces, not concrete types
//
// TYPE CHECKING
// - Type assertion: value, ok := iface.(ConcreteType)
// - Type switch: switch v := iface.(type) { case Type1: ... }
// - Determine concrete type at runtime
// - Extract type-specific information
// - Handle multiple types differently
//
// DESIGN PRINCIPLES
// - Composition over inheritance
// - Small, focused interfaces
// - Accept interfaces, return structs
// - Make zero values useful
// - Use pointer receivers for mutation
// - Be consistent with receiver types
// - Design for clarity and safety
//
// ==============================================================================
