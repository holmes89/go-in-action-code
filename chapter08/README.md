# Chapter 8: Testing and Tooling in Go

A comprehensive reference implementation demonstrating Go's testing ecosystem, including unit testing, HTTP testing, benchmarking, and fuzz testing. This demonstrates Test-Driven Development (TDD) with a Calculator example and various testing patterns.

## Quick Start

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestCalculate

# Run benchmarks
go test -bench=.

# Run benchmarks with memory stats
go test -bench=. -benchmem

# Run fuzz tests (for 30 seconds)
go test -fuzz=FuzzContains -fuzztime=30s

# Check code coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## What's Covered

This implementation demonstrates all major concepts from Chapter 8:

### 1. **Test Functions (testing.T)**
- Test function format: `func TestXxx(t *testing.T)`
- Arrange-Act-Assert pattern
- Fatal vs non-fatal failures
- Black box vs white box testing

### 2. **Test-Driven Development (TDD)**
- Writing tests before implementation
- Calculator example with TDD workflow
- Sentinel errors and error testing
- Memory feature testing

### 3. **Table-Driven Tests**
- Organizing multiple test cases
- Reducing code duplication
- Subtests with t.Run()
- Parameterized testing

### 4. **HTTP Testing (httptest)**
- Testing handlers with ResponseRecorder
- Testing servers with httptest.Server
- Testing routing and middleware
- Request/response validation

### 5. **Benchmark Testing (testing.B)**
- Performance measurement
- Comparing implementations
- Memory allocation tracking
- Interpreting benchmark results

### 6. **Fuzz Testing (testing.F)**
- Random input generation
- Edge case discovery
- Oracle testing pattern
- Corpus management

### 7. **Code Coverage**
- Measuring test coverage
- Generating coverage reports
- HTML coverage visualization
- Finding untested code

## Project Structure

```
chapter08/
├── calculator.go              # TDD example - Calculator implementation
├── calculator_test.go         # Unit tests, table-driven tests, examples
├── handlers.go                # HTTP handlers
├── handlers_test.go           # HTTP testing with httptest
├── xstrings.go                # String utilities for benchmarking/fuzzing
├── xstrings_test.go           # Benchmark and fuzz tests
├── go.mod                     # Module definition
└── README.md                  # This file
```

## Test Types Reference

### Unit Tests

```go
func TestFeature(t *testing.T) {
    // Arrange
    input := "test data"
    
    // Act
    result := FunctionToTest(input)
    
    // Assert
    if result != expected {
        t.Errorf("want: %v, got: %v", expected, result)
    }
}
```

### Table-Driven Tests

```go
func TestFeature(t *testing.T) {
    testCases := []struct {
        name string
        input string
        want string
    }{
        {name: "case1", input: "a", want: "A"},
        {name: "case2", input: "b", want: "B"},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            got := Function(tc.input)
            if got != tc.want {
                t.Errorf("want: %v, got: %v", tc.want, got)
            }
        })
    }
}
```

### Benchmark Tests

```go
func BenchmarkFeature(b *testing.B) {
    for i := 0; i < b.N; i++ {
        FunctionToTest(input)
    }
}
```

### Fuzz Tests

```go
func FuzzFeature(f *testing.F) {
    // Add corpus
    f.Add("seed1")
    f.Add("seed2")
    
    // Fuzz test
    f.Fuzz(func(t *testing.T, input string) {
        result := FunctionToTest(input)
        // Verify invariants
    })
}
```

## Testing Patterns

### 1. Arrange-Act-Assert (AAA)

```go
func TestCalculate(t *testing.T) {
    // Arrange - setup test conditions
    calc := NewCalculator()
    
    // Act - execute the code
    result, err := calc.Calculate("add", 2, 3)
    
    // Assert - verify results
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result != 5 {
        t.Errorf("want: 5, got: %d", result)
    }
}
```

### 2. Error Testing

```go
func TestErrorCondition(t *testing.T) {
    calc := NewCalculator()
    
    _, err := calc.Calculate("div", 1, 0)
    
    if !errors.Is(err, ErrDivByZero) {
        t.Errorf("expected ErrDivByZero, got: %v", err)
    }
}
```

### 3. Subtest Pattern

```go
func TestFeature(t *testing.T) {
    t.Run("subtest_name", func(t *testing.T) {
        // Test code
    })
}
```

## HTTP Testing Patterns

### Testing Individual Handlers

```go
func TestHandler(t *testing.T) {
    handler := http.HandlerFunc(YourHandler)
    rr := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/path", nil)
    
    handler.ServeHTTP(rr, req)
    
    if rr.Code != http.StatusOK {
        t.Errorf("wrong status code")
    }
}
```

### Testing with Server

```go
func TestServer(t *testing.T) {
    mux := http.NewServeMux()
    RegisterHandlers(mux)
    
    server := httptest.NewServer(mux)
    defer server.Close()
    
    resp, err := server.Client().Get(server.URL + "/path")
    // Test response
}
```

## Benchmark Testing

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=.

# Run specific benchmark
go test -bench=BenchmarkRandomString

# Include memory statistics
go test -bench=. -benchmem

# Run for specific duration
go test -bench=. -benchtime=10s

# Run with specific number of iterations
go test -bench=. -benchtime=100x
```

### Interpreting Results

```
BenchmarkRandomString-8    30609    119318 ns/op    24576 B/op    100 allocs/op
```

- `BenchmarkRandomString-8`: Function name + GOMAXPROCS
- `30609`: Number of iterations
- `119318 ns/op`: Nanoseconds per operation
- `24576 B/op`: Bytes allocated per operation
- `100 allocs/op`: Number of allocations per operation

### Comparison Tips

```bash
# Save baseline
go test -bench=. > old.txt

# Make changes, then compare
go test -bench=. > new.txt
benchstat old.txt new.txt  # Requires golang.org/x/perf/cmd/benchstat
```

## Fuzz Testing

### Running Fuzz Tests

```bash
# Run for 30 seconds
go test -fuzz=FuzzContains -fuzztime=30s

# Run with specific corpus
go test -fuzz=FuzzContains -fuzztime=1000x

# Re-run with saved failing input
go test -run=FuzzContains/HASH
```

### Fuzz Test Output

```
fuzz: elapsed: 3s, execs: 1148614 (382751/sec), new interesting: 0 (total: 5)
```

- `elapsed`: Time running
- `execs`: Total executions
- `execs/sec`: Operations per second
- `new interesting`: New coverage-expanding inputs
- `total`: Total interesting inputs in corpus

### Fuzzing Strategies

1. **Oracle Testing**: Compare with known-good implementation
2. **Invariant Testing**: Test properties that should always hold
3. **Differential Testing**: Compare two implementations
4. **Round-trip Testing**: Encode/decode should equal original

## Code Coverage

### Generating Coverage

```bash
# Show coverage percentage
go test -cover

# Generate coverage profile
go test -coverprofile=coverage.out

# View in browser
go tool cover -html=coverage.out

# Coverage by function
go tool cover -func=coverage.out
```

### Coverage Tips

- Aim for 70-80% coverage (100% is often impractical)
- Focus on critical paths
- Don't test trivial code just for coverage
- Use coverage to find missing tests, not as a goal

## Common Testing Patterns

### Testing Setup/Teardown

```go
func TestFeature(t *testing.T) {
    // Setup
    setup := prepareTest()
    defer cleanup(setup)  // Teardown
    
    // Test code
}
```

### Helper Functions

```go
// Mark as test helper
func testHelper(t *testing.T, input string) result {
    t.Helper()  // Failures report caller's location
    
    // Helper code
    return result
}
```

### Parallel Tests

```go
func TestFeature(t *testing.T) {
    t.Parallel()  // Run in parallel with other parallel tests
    
    // Test code
}
```

### Skipping Tests

```go
func TestFeature(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping in short mode")
    }
    
    // Long-running test
}
```

## Calculator Example (TDD)

The calculator demonstrates complete TDD workflow:

1. **Define requirements**
   - Basic operations (add, sub, mul, div)
   - Memory storage
   - Error handling

2. **Write tests first**
   - Error condition tests
   - Operation tests (table-driven)
   - Memory feature tests

3. **Implement to pass tests**
   - Add operations
   - Handle errors
   - Implement memory

4. **Refine with coverage**
   - Find missing tests
   - Discover bugs (Store not recovering after error)
   - Achieve 100% coverage

## HTTP Testing Example

Demonstrates two approaches:

### Unit Testing Handlers (ResponseRecorder)
- Fast
- Tests handler logic in isolation
- No routing/middleware

### Integration Testing (httptest.Server)
- More realistic
- Tests full HTTP stack
- Tests routing and middleware

## Benchmark Example

Compares two string generation implementations:
- `RandomString`: Using byte slice
- `RandomStringBuilder`: Using strings.Builder

Results show which is faster and allocates less memory.

## Fuzz Example

Tests `Contains` function with random inputs:
- Finds edge cases (Unicode, empty strings)
- Compares with standard library (oracle testing)
- Discovered bug in initial implementation

## Best Practices

### Testing
✅ Write tests for exported functions
✅ Test error conditions
✅ Use table-driven tests for multiple cases
✅ Use descriptive test names
✅ Log context for failures
✅ Test one thing per test function
✅ Use black box testing when possible

❌ Don't test unexported functions (test through public API)
❌ Don't test trivial getters/setters
❌ Don't make tests depend on each other
❌ Don't use `t.Fail()` without logging why

### Benchmarking
✅ Run multiple times for stable results
✅ Use `-benchmem` to track allocations
✅ Test realistic workloads
✅ Save baselines for comparison

❌ Don't include setup in benchmark loop
❌ Don't test trivial operations
❌ Don't optimize without benchmarking first

### Fuzzing
✅ Provide diverse seed corpus
✅ Test invariants, not specific outputs
✅ Run periodically (not in every test run)
✅ Keep saved failing inputs as regression tests

❌ Don't run in CI/CD without time limit
❌ Don't test exact outputs (too strict)
❌ Don't ignore saved failures

## Common Mistakes

### 1. Not Checking Errors

```go
// ❌ Bad
result, _ := Function()

// ✅ Good
result, err := Function()
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}
```

### 2. Poor Test Names

```go
// ❌ Bad
func TestFunc(t *testing.T) { }

// ✅ Good
func TestCalculateDivideByZero(t *testing.T) { }
```

### 3. Testing Implementation Details

```go
// ❌ Bad - testing unexported function
func TestInternalHelper(t *testing.T) { }

// ✅ Good - test through public API
func TestPublicFunction(t *testing.T) { }
```

### 4. Not Using Table-Driven Tests

```go
// ❌ Bad - repetitive
func TestAdd(t *testing.T) {
    if Add(1, 2) != 3 { t.Error() }
    if Add(2, 3) != 5 { t.Error() }
    if Add(3, 4) != 7 { t.Error() }
}

// ✅ Good - table-driven
func TestAdd(t *testing.T) {
    tests := []struct{ a, b, want int }{
        {1, 2, 3},
        {2, 3, 5},
        {3, 4, 7},
    }
    for _, tt := range tests {
        // Test
    }
}
```

## Testing Commands Cheat Sheet

```bash
# Basic testing
go test                          # Test current package
go test ./...                    # Test all packages
go test -v                       # Verbose output
go test -run TestName            # Run specific test
go test -short                   # Skip long-running tests

# Coverage
go test -cover                   # Show coverage
go test -coverprofile=c.out      # Generate profile
go tool cover -html=c.out        # View in browser
go tool cover -func=c.out        # Coverage by function

# Benchmarking
go test -bench=.                 # Run all benchmarks
go test -bench=. -benchmem       # Include memory stats
go test -bench=. -benchtime=10s  # Run for 10 seconds
go test -bench=. -cpuprofile=cpu.prof  # CPU profile

# Fuzzing
go test -fuzz=FuzzName           # Run fuzz test
go test -fuzz=. -fuzztime=30s    # Fuzz for 30 seconds
go test -fuzz=. -fuzztime=1000x  # 1000 iterations

# Profiling
go test -cpuprofile=cpu.prof     # CPU profile
go test -memprofile=mem.prof     # Memory profile
go tool pprof cpu.prof           # Analyze profile
```

## Further Exploration

Try extending the examples:

1. **Add more calculator operations** - square root, power, modulo
2. **Implement ordered map** - test with property-based testing
3. **Create REST API** - test with httptest
4. **Add middleware** - test authentication, logging
5. **Fuzz JSON parsing** - find edge cases
6. **Benchmark sorting algorithms** - compare implementations
7. **Test concurrent code** - use race detector

## Key Takeaways

✅ **Tests are documentation** - They show how code should be used

✅ **TDD drives design** - Writing tests first improves API design

✅ **Table-driven tests reduce duplication** - Easier to maintain

✅ **HTTP testing is easy** - httptest makes testing handlers simple

✅ **Benchmarking guides optimization** - Measure before optimizing

✅ **Fuzzing finds edge cases** - Discover bugs you didn't think of

✅ **Coverage finds gaps** - But don't obsess over 100%

✅ **Test the happy path AND errors** - Both are important

## Summary

Go's testing ecosystem provides:
- **Built-in testing framework** - No external dependencies
- **Table-driven test pattern** - Reduces duplication
- **HTTP testing tools** - httptest package
- **Performance measurement** - Benchmark testing
- **Fuzz testing** - Automatic edge case discovery
- **Code coverage tools** - Find untested code
- **Simple, explicit testing** - No magic, no surprises

This makes Go one of the best languages for writing testable, maintainable code.
