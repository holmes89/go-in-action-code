# Chapter 7: Error Handling in Go

A comprehensive reference implementation demonstrating Go's approach to error handling, which treats errors as ordinary values rather than exceptions. This file shows how to generate, handle, wrap, inspect, and log errors using a Key-Value store example.

## Quick Start

```bash
# Run the complete demonstration
go run main.go
```

## What's Covered

This implementation demonstrates all major concepts from Chapter 7:

### 1. **Errors as Values**
- The error interface (`Error() string`)
- Returning errors from functions
- Checking errors with `if err != nil`
- Error as last return value convention

### 2. **Generating Errors**
- `errors.New()` for simple errors
- `fmt.Errorf()` for formatted errors
- When to use each approach

### 3. **Sentinel Errors**
- Package-level error variables
- Naming convention (`ErrConflict`, `ErrNotExist`)
- Using sentinel errors for common conditions

### 4. **Error Wrapping**
- Adding context with `fmt.Errorf()` and `%w`
- Preserving original errors
- Building error chains

### 5. **Error Inspection**
- `errors.Is()` for checking specific errors
- Works with wrapped errors
- Comparing against sentinel errors

### 6. **Custom Error Types**
- Implementing the error interface
- Storing structured data in errors
- `errors.As()` for extracting custom types

### 7. **Logging Errors**
- Using `fmt` package for simple output
- Using `log` package for timestamped logs
- Customizing log flags
- Annotating errors with context

### 8. **Panic and Recover**
- When to use panic vs errors
- The `defer`/`recover` pattern
- Converting panics to errors
- Best practices for panic usage

### 9. **Error Handling Patterns**
- Early return pattern
- Named return for cleanup
- Error accumulation
- Error translation

## Key Concepts

### The Error Interface

```go
type error interface {
    Error() string
}
```

Any type implementing this interface is an error. That's it!

### Error Handling Flow

```go
// Pattern 1: Immediate check
err := store.Put("key", "value")
if err != nil {
    // Handle error
}

// Pattern 2: Inline check
if err := store.Put("key", "value"); err != nil {
    // Handle error
}

// Pattern 3: Multiple returns
val, err := store.Get("key")
if err != nil {
    // Handle error
}
// Use val
```

### Error Hierarchy

```
error (interface)
├── Simple errors (errors.New, fmt.Errorf)
├── Sentinel errors (package-level variables)
├── Wrapped errors (fmt.Errorf with %w)
└── Custom error types (struct implementing error)
```

## Code Organization

```
main.go
├── Error Type Definition (interface explanation)
├── Store Interface (methods that return errors)
├── Sentinel Errors (ErrConflict, ErrNotExist, ErrUninitialized)
├── KVStore Implementation
│   ├── NewKVStore (constructor)
│   ├── Put, Get, Delete (basic operations)
│   ├── ForcePut (errors.Is example)
│   ├── BulkGet (custom error type)
│   └── Keys, SafeKeys (panic/recover)
├── Custom Error Type (bulkGetError)
├── Error Handling Utilities
│   ├── HandleBulkGetError (errors.As example)
│   ├── wrappingExample (error wrapping)
│   └── Pattern functions (various patterns)
└── Demonstration Functions
```

## Sentinel Errors Reference

| Error | When It Occurs | How to Check |
|-------|---------------|--------------|
| `ErrConflict` | Key already exists | `errors.Is(err, ErrConflict)` |
| `ErrNotExist` | Key doesn't exist | `errors.Is(err, ErrNotExist)` |
| `ErrUninitialized` | Store not initialized | `errors.Is(err, ErrUninitialized)` |

### Standard Library Sentinel Errors

```go
io.EOF              // End of file/stream
sql.ErrNoRows       // Query returned no rows
os.ErrNotExist      // File/directory doesn't exist
os.ErrPermission    // Permission denied
```

## Error Wrapping

### Without Wrapping (❌ Loses Context)

```go
if err != nil {
    return fmt.Errorf("operation failed: %v", err)
}
// Original error is formatted into string, can't use errors.Is/As
```

### With Wrapping (✅ Preserves Original)

```go
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
// Original error is preserved, errors.Is/As work
```

### Why Use %w?

- Preserves error chain for `errors.Is()`
- Allows `errors.As()` to extract original types
- Adds context at each layer
- Enables error inspection up the call stack

## errors.Is() vs errors.As()

### errors.Is() - Check Error Identity

```go
err := store.Get("missing")
if errors.Is(err, ErrNotExist) {
    // Handle missing key specifically
}
```

Use when:
- Checking for sentinel errors
- Making decisions based on error type
- You only need to know IF an error occurred

### errors.As() - Extract Error Data

```go
_, err := store.BulkGet("key1", "missing")
var berr bulkGetError
if errors.As(err, &berr) {
    // Access berr slice to see which keys were missing
    for _, key := range berr {
        // Handle each missing key
    }
}
```

Use when:
- Need structured data from custom error
- Want to access error fields/methods
- Need more than just error identity

## Custom Error Type Pattern

```go
// 1. Define type with data
type MyError struct {
    Field1 string
    Field2 int
}

// 2. Implement Error() method
func (e MyError) Error() string {
    return fmt.Sprintf("error: %s (code: %d)", e.Field1, e.Field2)
}

// 3. Return from function
func DoSomething() error {
    return MyError{Field1: "failed", Field2: 42}
}

// 4. Extract in caller
var myErr MyError
if errors.As(err, &myErr) {
    // Access myErr.Field1, myErr.Field2
}
```

## Panic and Recover

### When to Panic

✅ **Do panic when:**
- Unrecoverable programmer error (nil pointer, uninitialized data)
- Impossible state that indicates a bug
- During initialization if critical resources unavailable

❌ **Don't panic when:**
- In library code (return errors instead)
- For expected errors (file not found, network error)
- When caller might want to handle the error

### Recover Pattern

```go
func SafeOperation() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    
    // Code that might panic
    return nil
}
```

### Special Cases

```go
panic(nil)  // ❌ Cannot be recovered! Avoid this.
panic("string")  // ✅ Can be recovered
panic(errors.New("error"))  // ✅ Can be recovered
```

## Logging Patterns

### Using fmt Package

```go
if err != nil {
    fmt.Println("Error:", err)              // Simple
    fmt.Printf("Error: %v\n", err)          // Formatted
    fmt.Printf("Details: %s\n", err.Error()) // Explicit
}
```

### Using log Package

```go
if err != nil {
    log.Println("Error:", err)              // With timestamp
    log.Printf("Error: %v\n", err)          // With timestamp
}

// Customize flags
log.SetFlags(log.LstdFlags | log.Lshortfile)  // Add file:line
log.SetFlags(log.Ltime)                        // Time only
log.SetFlags(0)                                // No extras
```

### Adding Context

```go
// ❌ Not helpful
if err != nil {
    log.Println(err)
}

// ✅ Better - adds context
if err != nil {
    log.Printf("failed to save user profile: %v", err)
}

// ✅ Best - wrap the error
if err != nil {
    return fmt.Errorf("failed to save user profile: %w", err)
}
```

## Error Handling Patterns

### 1. Early Return Pattern

```go
func Process() error {
    if err := step1(); err != nil {
        return err  // Return immediately
    }
    
    if err := step2(); err != nil {
        return err
    }
    
    return nil
}
```

**Use when:** Sequential operations where each depends on previous success.

### 2. Named Return Pattern

```go
func Process() (err error) {
    defer func() {
        if err != nil {
            err = fmt.Errorf("process failed: %w", err)
        }
    }()
    
    if err = step1(); err != nil {
        return  // err already set
    }
    
    return nil
}
```

**Use when:** Need to add context to any error that occurs.

### 3. Error Accumulation Pattern

```go
func ProcessAll(items []Item) error {
    var errs []error
    
    for _, item := range items {
        if err := process(item); err != nil {
            errs = append(errs, err)
            // Continue processing other items
        }
    }
    
    if len(errs) > 0 {
        return fmt.Errorf("errors: %v", errs)
    }
    return nil
}
```

**Use when:** Want to collect all errors before returning.

### 4. Error Translation Pattern

```go
func LoadConfig() error {
    err := readFile()
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            return ErrConfigNotFound  // Translate to package error
        }
        return fmt.Errorf("config load failed: %w", err)
    }
    return nil
}
```

**Use when:** Abstracting implementation details from callers.

## Common Mistakes

### 1. Ignoring Errors

```go
// ❌ Bad - silently ignoring error
store.Put("key", "value")

// ✅ Good - explicit handling
if err := store.Put("key", "value"); err != nil {
    // Handle appropriately
}

// ✅ Also acceptable if you truly don't care
_ = store.Put("key", "value")  // Explicit ignore
```

### 2. Not Adding Context

```go
// ❌ Bad - loses context
if err != nil {
    return err
}

// ✅ Good - adds context
if err != nil {
    return fmt.Errorf("failed to save user data: %w", err)
}
```

### 3. Using %v Instead of %w

```go
// ❌ Bad - can't use errors.Is/As later
return fmt.Errorf("operation failed: %v", err)

// ✅ Good - preserves error for inspection
return fmt.Errorf("operation failed: %w", err)
```

### 4. Swallowing Panics

```go
// ❌ Bad - panic becomes invisible
defer func() {
    recover()  // Silent! No one knows panic happened
}()

// ✅ Good - log or return panic
defer func() {
    if r := recover(); r != nil {
        log.Printf("Panic: %v", r)
        // Or convert to error and return
    }
}()
```

### 5. Panicking in Libraries

```go
// ❌ Bad - library code shouldn't panic
func ProcessData(data []byte) {
    if len(data) == 0 {
        panic("empty data")
    }
}

// ✅ Good - return error instead
func ProcessData(data []byte) error {
    if len(data) == 0 {
        return errors.New("empty data")
    }
    return nil
}
```

## Best Practices

✅ **Do:**
- Check errors immediately after they occur
- Return errors to the caller when you can't handle them
- Add context when wrapping errors with `%w`
- Use sentinel errors for common, expected conditions
- Use custom error types when you need structured data
- Use `errors.Is()` to check for specific errors
- Use `errors.As()` to extract custom error data
- Document which errors your functions can return
- Log errors at appropriate levels
- Use defer/recover only for truly exceptional situations

❌ **Don't:**
- Ignore errors (unless you have a good reason)
- Panic in library code
- Use panic for normal error handling
- Return errors without context
- Create new errors when wrapping (use `%w`)
- Handle errors far from where they occur
- Swallow panics silently
- Use `panic(nil)`

## Testing Error Handling

```go
func TestPutConflict(t *testing.T) {
    store := NewKVStore()
    
    // First put succeeds
    err := store.Put("key", "value1")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    // Second put should fail
    err = store.Put("key", "value2")
    if !errors.Is(err, ErrConflict) {
        t.Errorf("expected ErrConflict, got %v", err)
    }
}

func TestBulkGetError(t *testing.T) {
    store := NewKVStore()
    store.Put("key1", "value1")
    
    _, err := store.BulkGet("key1", "missing")
    if err == nil {
        t.Fatal("expected error")
    }
    
    var berr bulkGetError
    if !errors.As(err, &berr) {
        t.Fatal("expected bulkGetError")
    }
    
    if len(berr) != 1 || berr[0] != "missing" {
        t.Errorf("unexpected missing keys: %v", berr)
    }
}
```

## Real-World Examples

### Configuration Loading

```go
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            return DefaultConfig(), nil  // Use defaults
        }
        return nil, fmt.Errorf("read config: %w", err)
    }
    
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parse config: %w", err)
    }
    
    return &cfg, nil
}
```

### Network Requests with Retry

```go
func FetchData(url string) ([]byte, error) {
    var lastErr error
    
    for i := 0; i < 3; i++ {
        resp, err := http.Get(url)
        if err == nil {
            defer resp.Body.Close()
            return io.ReadAll(resp.Body)
        }
        
        lastErr = err
        time.Sleep(time.Second * time.Duration(i+1))
    }
    
    return nil, fmt.Errorf("fetch failed after 3 retries: %w", lastErr)
}
```

### Database Operations

```go
func GetUser(id int) (*User, error) {
    var user User
    err := db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user)
    
    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrUserNotFound  // Custom sentinel
    }
    if err != nil {
        return nil, fmt.Errorf("query user %d: %w", id, err)
    }
    
    return &user, nil
}
```

## Further Exploration

Try extending the examples:

1. **Add error metrics** - Count error types
2. **Implement retry logic** - Retry on specific errors
3. **Create error middleware** - Centralized error handling
4. **Add structured logging** - Use `slog` package (Chapter 11)
5. **Implement error codes** - Numeric error codes for APIs
6. **Add stack traces** - Enhanced debugging information

## Key Takeaways

✅ **Errors are values** - Handle them like any other data

✅ **Check immediately** - Don't delay error handling

✅ **Add context** - Use `fmt.Errorf` with `%w`

✅ **Use sentinel errors** - For known, common conditions

✅ **Custom types for data** - When you need structured info

✅ **errors.Is() for identity** - Check which error occurred

✅ **errors.As() for data** - Extract structured error data

✅ **Panic rarely** - Only for programmer errors

✅ **Recover carefully** - Use defer/recover pattern

✅ **Log appropriately** - With context and timestamps

## Summary

Go's error handling philosophy emphasizes:
- **Explicitness** over hidden exceptions
- **Immediacy** over delayed handling  
- **Simplicity** over complex hierarchies
- **Flexibility** through interfaces
- **Safety** through compile-time checking

This approach makes error handling predictable, testable, and maintainable, even if it requires more upfront code.
