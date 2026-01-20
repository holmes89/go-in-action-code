package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

// ============================================================================
// ERROR HANDLING IN GO
// ============================================================================
// This file demonstrates Go's approach to error handling, which treats errors
// as ordinary values rather than using exceptions. This makes error handling
// explicit, predictable, and part of the normal control flow.
//
// Key Concepts Covered:
// 1. Errors as Values - the error type and basic handling
// 2. Error Generation - errors.New() and fmt.Errorf()
// 3. Sentinel Errors - predefined package-level error values
// 4. Error Wrapping - adding context with %w
// 5. Error Inspection - errors.Is() and errors.As()
// 6. Custom Error Types - implementing the error interface
// 7. Logging Errors - fmt and log packages
// 8. Panic and Recover - exceptional situations
//
// Throughout this file, we build a Key-Value store to demonstrate
// real-world error handling patterns.
// ============================================================================

// ============================================================================
// 1. THE ERROR TYPE
// ============================================================================
// In Go, errors are just values of the 'error' type. The error type is
// actually an interface with a single method:
//
//   type error interface {
//       Error() string
//   }
//
// This simple interface means:
// - Any type can be an error (just implement Error() string)
// - Errors are returned explicitly from functions
// - Error handling is part of normal control flow
// - No hidden exceptions or try/catch blocks

// ============================================================================
// 2. STORE INTERFACE
// ============================================================================
// Our Store interface defines operations for a Key-Value database.
// Notice how error is the last return value in functions that can fail.
// This is Go convention.

type Store interface {
	Put(key string, val string) error              // Only returns error
	Get(key string) (string, error)                // Returns value and error
	Delete(key string) error                       // Only returns error
	ForcePut(key string, val string) error         // Overwrites existing values
	BulkGet(keys ...string) ([]string, error)      // Gets multiple values
	Keys() []string                                // Returns all keys (can panic if uninitialized)
}

// ============================================================================
// 3. SENTINEL ERRORS
// ============================================================================
// Sentinel errors are predefined package-level error values that represent
// specific conditions. They allow callers to check for specific errors
// programmatically using errors.Is().
//
// Convention:
// - Defined as package-level variables
// - Named with "Err" prefix
// - Grouped in a var block for visibility
// - Messages are lowercase (not sentences)

var (
	// ErrConflict indicates a key already exists when it shouldn't
	ErrConflict = errors.New("key already exists")

	// ErrNotExist indicates a key doesn't exist when it should
	ErrNotExist = errors.New("no value exists for key")

	// ErrUninitialized indicates the store wasn't properly initialized
	ErrUninitialized = errors.New("store not properly initialized")
)

// Why sentinel errors?
// - Consistent error values across package
// - Enable programmatic error handling with errors.Is()
// - Callers can make decisions based on error type
// - Common in standard library (io.EOF, sql.ErrNoRows, etc.)

// ============================================================================
// 4. KVSTORE IMPLEMENTATION
// ============================================================================
// KVStore is a simple in-memory key-value store that demonstrates
// various error handling patterns.

type KVStore struct {
	m map[string]string // The underlying storage
}

// NewKVStore creates a properly initialized KVStore.
// This is the safe way to create a store - it ensures the map is initialized.
//
// Design note: Constructors like this prevent initialization errors and
// make it clear to users how to properly create the type.
func NewKVStore() *KVStore {
	return &KVStore{
		m: make(map[string]string),
	}
}

// ============================================================================
// 5. BASIC ERROR GENERATION AND HANDLING
// ============================================================================

// Put adds a key-value pair to the store.
// Returns ErrConflict if the key already exists.
//
// Error handling pattern demonstrated:
// - Check for error condition first
// - Return sentinel error for known cases
// - Return nil on success
func (store *KVStore) Put(key string, val string) error {
	// Check if key exists using the "comma ok" idiom
	if _, ok := store.m[key]; ok {
		// Return sentinel error for this specific condition
		return ErrConflict
	}

	// Success path - store the value
	store.m[key] = val
	return nil // nil means no error
}

// Get retrieves a value for a given key.
// Returns the value and nil error if found.
// Returns empty string and ErrNotExist if not found.
//
// Important pattern: When returning multiple values with an error,
// return meaningful "zero values" when there's an error.
func (store *KVStore) Get(key string) (string, error) {
	v, ok := store.m[key]
	if !ok {
		// Key doesn't exist - return zero value and error
		return "", ErrNotExist
	}

	// Success - return value and nil error
	return v, nil
}

// Delete removes a key-value pair from the store.
// Returns ErrNotExist if the key doesn't exist.
func (store *KVStore) Delete(key string) error {
	_, ok := store.m[key]
	if !ok {
		return ErrNotExist
	}

	delete(store.m, key)
	return nil
}

// ============================================================================
// 6. ERROR INSPECTION WITH errors.Is()
// ============================================================================
// errors.Is() checks if an error matches a specific sentinel error,
// even if the error has been wrapped with additional context.

// ForcePut adds a key-value pair, overwriting if it already exists.
// This demonstrates:
// - Using errors.Is() to check for specific errors
// - Error wrapping with fmt.Errorf and %w
// - Chaining operations with error handling
func (store *KVStore) ForcePut(key string, val string) error {
	// Try normal put first
	err := store.Put(key, val)
	if err == nil {
		// Success - no error
		return nil
	}

	// Check if the error is specifically a conflict
	// errors.Is() works even if the error has been wrapped
	if errors.Is(err, ErrConflict) {
		// Handle the conflict by deleting and re-adding

		// Delete the existing key
		if err := store.Delete(key); err != nil {
			// Wrap the error with context using %w
			// This preserves the original error for errors.Is() and errors.As()
			return fmt.Errorf("unable to delete during force put: %w", err)
		}

		// Put the new value
		if err := store.Put(key, val); err != nil {
			return fmt.Errorf("unable to put during force put: %w", err)
		}

		return nil
	}

	// Unknown error - wrap it with context
	return fmt.Errorf("unable to force put: %w", err)
}

// ============================================================================
// 7. CUSTOM ERROR TYPES
// ============================================================================
// Sometimes you need more than just a message - you need to carry
// additional data in your errors. You can do this by creating a custom
// type that implements the error interface.

// bulkGetError captures which keys were missing during a bulk get operation.
// By making this a custom type, we can:
// - Store structured data (the missing keys)
// - Format the error message dynamically
// - Allow callers to extract this data with errors.As()
type bulkGetError []string

// addKey adds a missing key to the error.
// This is a helper method for building up the error during processing.
func (e bulkGetError) addKey(key string) bulkGetError {
	return append(e, key)
}

// Error implements the error interface.
// This is the only method required to make any type an error.
//
// Design note: Error messages should be:
// - Lowercase (not starting with capital unless proper noun)
// - Not end with punctuation
// - Descriptive enough to understand what went wrong
func (e bulkGetError) Error() string {
	if len(e) == 0 {
		return ""
	}
	// Format the error message with the missing keys
	return fmt.Sprintf("missing keys: %s", strings.Join(e, ", "))
}

// ============================================================================
// 8. USING CUSTOM ERROR TYPES
// ============================================================================

// BulkGet retrieves multiple values at once.
// This demonstrates:
// - Building up a custom error during processing
// - Continuing after errors to collect all issues
// - Returning partial results with an error
func (store *KVStore) BulkGet(keys ...string) ([]string, error) {
	var res []string
	var errb bulkGetError // Our custom error type

	// Process each key
	for _, key := range keys {
		v, err := store.Get(key)

		// Check if this specific error is "not exist"
		if errors.Is(err, ErrNotExist) {
			// Add to our custom error and continue processing
			errb = errb.addKey(key)
			continue
		}

		// Unexpected error - return immediately
		if err != nil {
			return nil, fmt.Errorf("unexpected error during bulk get: %w", err)
		}

		// Success - add value to results
		res = append(res, v)
	}

	// Check if we accumulated any missing keys
	// Note: Empty slices are NOT nil in Go, but we can still check
	if len(errb) > 0 {
		// Return partial results AND the error
		// This is valid in Go - you can return meaningful data with errors
		return res, errb
	}

	return res, nil
}

// ============================================================================
// 9. ERROR INSPECTION WITH errors.As()
// ============================================================================
// errors.As() extracts the underlying error type from an error value.
// This allows you to access the structured data in custom errors.

// HandleBulkGetError demonstrates using errors.As() to extract and handle
// a custom error type. This shows how callers can programmatically respond
// to specific error types.
func HandleBulkGetError(store *KVStore, keys ...string) []string {
	vals, err := store.BulkGet(keys...)
	if err != nil {
		// Try to extract our custom error type
		var berr bulkGetError

		// errors.As() checks if err contains a bulkGetError
		// and if so, assigns it to berr
		// Note: We pass a POINTER to berr (&berr) so it can be assigned
		if errors.As(err, &berr) {
			fmt.Println("Handling missing keys from bulk get:")

			// Now we have access to the slice of missing keys
			for _, key := range berr {
				fmt.Printf("  - Key '%s' was missing, adding placeholder\n", key)
				// Handle by adding placeholder values
				_ = store.Put(key, "placeholder")
			}

			// Try again after fixing
			vals, err = store.BulkGet(keys...)
			if err != nil {
				fmt.Printf("Still got error after handling: %v\n", err)
			}
		} else {
			// Not our custom error - handle generically
			fmt.Printf("Unexpected error: %v\n", err)
		}
	}

	return vals
}

// ============================================================================
// 10. PANIC AND RECOVER
// ============================================================================
// Panic is Go's mechanism for truly exceptional situations that should
// terminate the program. Unlike errors, panics:
// - Break normal control flow
// - Unwind the stack
// - Execute deferred functions
// - Terminate the program (unless recovered)
//
// When to panic:
// - Programmer errors (uninitialized data, impossible state)
// - Situations where recovery is meaningless
// - Never in library code (return errors instead)

// Keys returns all keys in the store.
// Panics if the store is uninitialized (map is nil).
//
// Design note: This demonstrates when to panic vs return an error.
// Since an uninitialized store is a programmer error (they didn't use
// NewKVStore()), we panic rather than returning an error.
func (store *KVStore) Keys() []string {
	// Check for uninitialized state
	if store.m == nil {
		// Panic with a descriptive message
		// panic() accepts any value, but strings and errors are most common
		panic("uninitialized kv store: use NewKVStore() to create")
	}

	// Collect all keys
	out := make([]string, 0, len(store.m))
	for k := range store.m {
		out = append(out, k)
	}
	return out
}

// SafeKeys demonstrates using recover() to handle a potential panic.
// This shows the defer/recover pattern for graceful panic handling.
//
// Important: recover() only works inside deferred functions!
func (store *KVStore) SafeKeys() (keys []string, err error) {
	// defer executes after the function returns (or panics)
	defer func() {
		// recover() returns nil if there was no panic
		// Otherwise it returns the value passed to panic()
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)

			// We can convert the panic to an error
			// Named return values let us set err here
			err = fmt.Errorf("panic during Keys(): %v", r)
		}
	}()

	// This might panic, but we'll recover
	keys = store.Keys()
	return keys, nil
}

// ============================================================================
// 11. ERROR WRAPPING PATTERNS
// ============================================================================
// Error wrapping adds context while preserving the original error.
// Use %w in fmt.Errorf() to wrap errors properly.

// wrappingExample demonstrates different levels of error wrapping.
// This shows how errors accumulate context as they bubble up.
func wrappingExample() {
	store := NewKVStore()

	// Level 1: Base error (from Put)
	err1 := store.Put("key1", "value1")
	err1 = store.Put("key1", "value1") // Causes conflict
	fmt.Printf("Base error: %v\n", err1)
	fmt.Printf("Is ErrConflict? %v\n\n", errors.Is(err1, ErrConflict))

	// Level 2: Wrapped once (from ForcePut)
	err2 := store.ForcePut("key2", "value2")
	if err2 != nil {
		fmt.Printf("Wrapped error: %v\n", err2)
	}

	// Demonstrate wrapping with additional context
	if _, err := store.Get("nonexistent"); err != nil {
		wrapped := fmt.Errorf("failed to load configuration: %w", err)
		fmt.Printf("Wrapped with context: %v\n", wrapped)
		fmt.Printf("Is ErrNotExist? %v\n\n", errors.Is(wrapped, ErrNotExist))
	}
}

// ============================================================================
// 12. LOGGING ERRORS
// ============================================================================
// Go provides two main packages for logging: fmt and log.
// - fmt: Simple printing, no timestamps
// - log: Adds timestamps and optional file/line info

func demonstrateLogging() {
	fmt.Println("\n=== Logging Demonstrations ===")

	store := NewKVStore()

	// Using fmt for error output
	fmt.Println("\n1. Using fmt package:")
	_, err := store.Get("missing")
	if err != nil {
		fmt.Println("Error occurred:", err)                    // Simple
		fmt.Printf("Error with formatting: %v\n", err)         // Printf
		fmt.Printf("Error with more context: %s\n", err)       // %s also works
	}

	// Using log for error output
	fmt.Println("\n2. Using log package:")
	_, err = store.Get("missing")
	if err != nil {
		log.Println("Error occurred:", err)                    // Adds timestamp
		log.Printf("Error with formatting: %v\n", err)         // Printf with timestamp
	}

	// Customizing log output
	fmt.Println("\n3. Customizing log flags:")
	
	// Save original flags
	originalFlags := log.Flags()
	
	// Add file and line number to logs
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Log with file:line info")

	// Just timestamp, no date
	log.SetFlags(log.Ltime)
	log.Println("Log with time only")

	// No flags at all
	log.SetFlags(0)
	log.Println("Log with no extras")

	// Restore original flags
	log.SetFlags(originalFlags)

	// Annotating errors with context
	fmt.Println("\n4. Annotating errors:")
	store.Put("key1", "value1")
	if err := store.Put("key1", "value2"); err != nil {
		log.Printf("failed to save user preference: %v", err)
	}

	// Even better - wrap the error
	if err := store.Put("key1", "value3"); err != nil {
		wrapped := fmt.Errorf("failed to save user preference: %w", err)
		log.Println(wrapped)
	}
}

// ============================================================================
// 13. ERROR HANDLING PATTERNS
// ============================================================================

// Pattern 1: Early return
// Check errors immediately and return if found
func earlyReturnPattern(store *KVStore) error {
	// Check error immediately
	if err := store.Put("key1", "value1"); err != nil {
		return err // Return immediately
	}

	// Only continue if no error
	if err := store.Put("key2", "value2"); err != nil {
		return err
	}

	return nil
}

// Pattern 2: Named return for cleanup
// Use named returns to set error in defer
func namedReturnPattern(store *KVStore) (err error) {
	defer func() {
		if err != nil {
			// Add context to any error that occurs
			err = fmt.Errorf("namedReturnPattern failed: %w", err)
		}
	}()

	if err = store.Put("key1", "value1"); err != nil {
		return // err is already set
	}

	return nil
}

// Pattern 3: Error accumulation
// Collect multiple errors before returning
func errorAccumulationPattern(store *KVStore, keys []string) error {
	var errs []error

	for _, key := range keys {
		if err := store.Put(key, "value"); err != nil {
			// Collect error but continue processing
			errs = append(errs, fmt.Errorf("key %s: %w", key, err))
		}
	}

	// Return all errors together
	if len(errs) > 0 {
		return fmt.Errorf("multiple errors: %v", errs)
	}

	return nil
}

// Pattern 4: Error translation
// Convert errors from one package to errors in your package
func errorTranslationPattern(store *KVStore) error {
	_, err := store.Get("key")
	if err != nil {
		// Don't leak internal errors - translate to your package's errors
		if errors.Is(err, ErrNotExist) {
			// Could return a different error that makes sense for this layer
			return fmt.Errorf("configuration not found: %w", err)
		}
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	return nil
}

// ============================================================================
// DEMONSTRATION FUNCTIONS
// ============================================================================

func demonstrateBasicErrors() {
	fmt.Println("\n=== Basic Error Handling ===")

	store := NewKVStore()

	// Success case
	fmt.Println("\n1. Successful operations:")
	if err := store.Put("user1", "Alice"); err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Successfully added user1")
	}

	val, err := store.Get("user1")
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Retrieved: %s\n", val)
	}

	// Error case
	fmt.Println("\n2. Error cases:")
	err2 := store.Put("user1", "Bob") // Conflict
	if err2 != nil {
		fmt.Printf("Expected error: %v\n", err2)
	}

	_, err3 := store.Get("nonexistent")
	if err3 != nil {
		fmt.Printf("Expected error: %v\n", err3)
	}
}

func demonstrateSentinelErrors() {
	fmt.Println("\n=== Sentinel Errors with errors.Is() ===")

	store := NewKVStore()
	store.Put("existing", "value")

	// Test for specific sentinel errors
	fmt.Println("\n1. Checking for ErrConflict:")
	err := store.Put("existing", "newvalue")
	if errors.Is(err, ErrConflict) {
		fmt.Println("Detected conflict - key already exists")
		fmt.Println("Could prompt user for confirmation to overwrite")
	}

	fmt.Println("\n2. Checking for ErrNotExist:")
	_, err = store.Get("missing")
	if errors.Is(err, ErrNotExist) {
		fmt.Println("Detected missing key")
		fmt.Println("Could create default value or prompt user")
	}

	// Demonstrate errors.Is() with wrapped errors
	fmt.Println("\n3. errors.Is() works with wrapped errors:")
	wrappedErr := fmt.Errorf("operation failed: %w", ErrConflict)
	if errors.Is(wrappedErr, ErrConflict) {
		fmt.Println("Found ErrConflict even though it was wrapped!")
	}
}

func demonstrateForcePut() {
	fmt.Println("\n=== ForcePut with Error Handling ===")

	store := NewKVStore()

	// Add initial value
	store.Put("config", "version1")
	val, _ := store.Get("config")
	fmt.Printf("Initial value: %s\n", val)

	// Normal Put would fail
	if err := store.Put("config", "version2"); err != nil {
		fmt.Printf("Put failed as expected: %v\n", err)
	}

	// ForcePut handles the conflict
	if err := store.ForcePut("config", "version2"); err != nil {
		log.Printf("ForcePut failed: %v\n", err)
	} else {
		fmt.Println("ForcePut succeeded!")
		val, _ = store.Get("config")
		fmt.Printf("Updated value: %s\n", val)
	}
}

func demonstrateCustomErrors() {
	fmt.Println("\n=== Custom Error Types ===")

	store := NewKVStore()
	store.Put("key1", "value1")
	store.Put("key2", "value2")

	// BulkGet with some missing keys
	fmt.Println("\n1. BulkGet with missing keys:")
	vals, err := store.BulkGet("key1", "key2", "missing1", "missing2")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Got %d values despite error\n", len(vals))
	}

	// Using errors.As() to handle the specific error
	fmt.Println("\n2. Handling custom error with errors.As():")
	vals = HandleBulkGetError(store, "key1", "missing3", "missing4")
	fmt.Printf("After handling, got %d values\n", len(vals))
}

func demonstratePanicAndRecover() {
	fmt.Println("\n=== Panic and Recover ===")

	// Demonstrate panic with uninitialized store
	fmt.Println("\n1. Uninitialized store causes panic:")
	uninitStore := &KVStore{} // No map initialized!

	// This would panic and crash the program:
	// uninitStore.Keys() // DON'T DO THIS

	// But SafeKeys recovers from the panic:
	fmt.Println("\n2. SafeKeys recovers from panic:")
	keys, err := uninitStore.SafeKeys()
	if err != nil {
		fmt.Printf("Caught panic as error: %v\n", err)
	} else {
		fmt.Printf("Got keys: %v\n", keys)
	}

	// Properly initialized store works fine
	fmt.Println("\n3. Properly initialized store:")
	properStore := NewKVStore()
	properStore.Put("key1", "value1")
	properStore.Put("key2", "value2")

	keys, err = properStore.SafeKeys()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Keys: %v\n", keys)
	}
}

func demonstrateErrorWrapping() {
	fmt.Println("\n=== Error Wrapping ===")
	wrappingExample()
}

func demonstrateErrorPatterns() {
	fmt.Println("\n=== Error Handling Patterns ===")

	store := NewKVStore()
	store.Put("existing", "value")

	fmt.Println("\n1. Early return pattern:")
	if err := earlyReturnPattern(store); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Success!")
	}

	fmt.Println("\n2. Named return pattern:")
	if err := namedReturnPattern(store); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Success!")
	}

	fmt.Println("\n3. Error accumulation pattern:")
	keys := []string{"new1", "existing", "new2"}
	if err := errorAccumulationPattern(store, keys); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n4. Error translation pattern:")
	if err := errorTranslationPattern(store); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// ============================================================================
// MAIN FUNCTION
// ============================================================================

func main() {
	fmt.Println("Go Error Handling Comprehensive Reference")
	fmt.Println("==========================================")

	// 1. Basic error handling patterns
	demonstrateBasicErrors()

	// 2. Sentinel errors and errors.Is()
	demonstrateSentinelErrors()

	// 3. ForcePut demonstrating error inspection
	demonstrateForcePut()

	// 4. Custom error types and errors.As()
	demonstrateCustomErrors()

	// 5. Logging errors with fmt and log
	demonstrateLogging()

	// 6. Error wrapping with %w
	demonstrateErrorWrapping()

	// 7. Panic and recover
	demonstratePanicAndRecover()

	// 8. Various error handling patterns
	demonstrateErrorPatterns()

	fmt.Println("\n=== Summary ===")
	fmt.Println("Error handling in Go:")
	fmt.Println("✓ Errors are values, not exceptions")
	fmt.Println("✓ Check errors immediately after they occur")
	fmt.Println("✓ Use sentinel errors for known conditions")
	fmt.Println("✓ Wrap errors with context using %w")
	fmt.Println("✓ Use errors.Is() to check for specific errors")
	fmt.Println("✓ Use errors.As() to extract custom error types")
	fmt.Println("✓ Panic only for programmer errors")
	fmt.Println("✓ Recover from panics in deferred functions")
	fmt.Println("✓ Log errors with appropriate context")
}
