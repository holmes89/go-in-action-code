package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// ============================================================================
// CHAPTER 8: TESTING AND TOOLING
// ============================================================================
// This file demonstrates Go's comprehensive testing ecosystem.
// While normally you'd put tests in _test.go files, this comprehensive
// reference shows all the concepts in one place for learning purposes.
//
// Topics Covered:
// 1. Test Functions (testing.T) - Basic unit testing
// 2. Test-Driven Development - Writing tests first
// 3. Table-Driven Tests - Testing multiple cases efficiently
// 4. HTTP Testing (httptest) - Testing web handlers
// 5. Benchmark Testing (testing.B) - Performance measurement
// 6. Fuzz Testing (testing.F) - Finding edge cases
// 7. Code Coverage - Measuring test completeness
// 8. Error Testing - Verifying error conditions
//
// IMPORTANT: This is a REFERENCE implementation.
// In real projects, tests go in separate *_test.go files!
//
// See main_test.go for the actual runnable tests.

// ============================================================================
// PART 1: CALCULATOR - TEST-DRIVEN DEVELOPMENT EXAMPLE
// ============================================================================
// This demonstrates TDD workflow: Write tests BEFORE implementation.
// TDD Benefits:
// - Clear requirements upfront
// - Better API design (you use it first!)
// - Complete test coverage
// - Catches bugs early

// Calculator sentinel errors allow users to check for specific failures
var (
	ErrDivByZero        = errors.New("divide by zero")
	ErrUnknownOperation = errors.New("unknown operation")
)

// Calculator provides basic math operations with memory storage.
// Designed using TDD to meet these requirements:
// - Support add, sub, mul, div operations
// - Store last successful calculation in memory
// - Retrieve stored memory value
// - Return errors for invalid operations or divide by zero
type Calculator struct {
	mem    int   // Memory storage
	result int   // Last calculation result
	err    error // Last error (prevents storing invalid results)
}

// NewCalculator returns a calculator ready to use.
func NewCalculator() *Calculator {
	return &Calculator{}
}

// Calculate performs a mathematical operation.
// Valid operations: "add", "sub", "mul", "div"
// Returns the result and an error if the operation fails.
func (c *Calculator) Calculate(operation string, a, b int) (int, error) {
	// Check divide by zero BEFORE the switch to prevent panic
	if operation == "div" && b == 0 {
		c.err = ErrDivByZero
		c.result = 0
		return 0, ErrDivByZero
	}

	// IMPORTANT: Reset error on successful calculation
	// This allows Store() to work after a previous error
	// (Bug discovered through code coverage testing!)
	c.err = nil

	switch operation {
	case "add":
		c.result = a + b
	case "sub":
		c.result = a - b
	case "mul":
		c.result = a * b
	case "div":
		c.result = a / b
	default:
		c.err = ErrUnknownOperation
		return 0, ErrUnknownOperation
	}

	return c.result, nil
}

// Store saves the current result to memory if valid.
// Does nothing if the last calculation resulted in an error.
func (c *Calculator) Store() {
	if c.err != nil {
		return // Don't store invalid results
	}
	c.mem = c.result
}

// Recall returns the value stored in memory.
// For a new calculator, this returns 0.
func (c *Calculator) Recall() int {
	return c.mem
}

// ============================================================================
// DEMONSTRATION 1: Basic Calculator Usage
// ============================================================================

func demonstrateCalculator() {
	fmt.Println("\n=== DEMONSTRATION 1: Calculator (TDD Example) ===")

	c := NewCalculator()

	// Basic operations
	result, _ := c.Calculate("add", 5, 3)
	fmt.Printf("5 + 3 = %d\n", result)

	result, _ = c.Calculate("mul", 4, 7)
	fmt.Printf("4 * 7 = %d\n", result)

	// Using memory
	c.Calculate("div", 100, 5) // 20
	c.Store()
	fmt.Printf("Stored in memory: %d\n", c.Recall())

	// Error handling
	_, err := c.Calculate("div", 10, 0)
	if errors.Is(err, ErrDivByZero) {
		fmt.Printf("Error caught: %v\n", err)
	}

	// Memory not updated after error
	fmt.Printf("Memory still: %d (unchanged)\n", c.Recall())

	fmt.Println("\nTesting Approach:")
	fmt.Println("1. Write tests first (define requirements)")
	fmt.Println("2. Implement to pass tests")
	fmt.Println("3. Use coverage to find gaps")
	fmt.Println("4. Refactor with confidence")
}

// ============================================================================
// PART 2: HTTP HANDLERS - WEB SERVICE TESTING
// ============================================================================
// Go's httptest package makes testing HTTP handlers easy.
// Two approaches:
// 1. httptest.ResponseRecorder - Unit test handlers directly
// 2. httptest.Server - Integration test with routing

// HelloGet handles GET requests and returns "hello".
// This is the simplest possible HTTP handler.
func HelloGet(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hello"))
}

// HelloPost handles POST requests and echoes back "hello " + body.
// Demonstrates reading request data and error handling.
func HelloPost(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hello := append([]byte("hello "), body...)
	w.WriteHeader(http.StatusCreated)
	w.Write(hello)
}

// HelloHandler sets up routing for /hello endpoint.
// Routes to different handlers based on HTTP method.
func HelloHandler(mux *http.ServeMux) {
	mux.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			HelloGet(w, req)
		case http.MethodPost:
			HelloPost(w, req)
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("invalid hello method"))
		}
	})
}

// ============================================================================
// DEMONSTRATION 2: HTTP Handler Testing
// ============================================================================

func demonstrateHTTPTesting() {
	fmt.Println("\n=== DEMONSTRATION 2: HTTP Testing ===")

	// Approach 1: Testing handler directly with ResponseRecorder
	fmt.Println("\n1. Unit Testing Handler (ResponseRecorder):")

	handler := http.HandlerFunc(HelloGet)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "", nil)

	handler.ServeHTTP(rr, req)

	fmt.Printf("   Status: %d\n", rr.Code)
	fmt.Printf("   Body: %s\n", rr.Body.String())
	fmt.Println("   ✓ Fast, tests handler logic only")

	// Approach 2: Testing with actual server
	fmt.Println("\n2. Integration Testing (httptest.Server):")

	mux := http.NewServeMux()
	HelloHandler(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Make real HTTP request
	fullURL, _ := url.JoinPath(server.URL, "/hello")
	resp, _ := server.Client().Get(fullURL)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Printf("   Status: %d\n", resp.StatusCode)
	fmt.Printf("   Body: %s\n", string(body))
	fmt.Println("   ✓ Realistic, tests routing and middleware")

	fmt.Println("\nTesting Strategies:")
	fmt.Println("• ResponseRecorder: Fast unit tests for handler logic")
	fmt.Println("• httptest.Server: Integration tests for routing")
	fmt.Println("• Table-driven tests: Test multiple endpoints efficiently")
}

// ============================================================================
// PART 3: STRING UTILITIES - BENCHMARKING AND FUZZING
// ============================================================================
// These functions demonstrate performance testing and edge case discovery.

var alphabetArray = []byte("abcdefghijklmnopqrstuvwxyz")

// RandomString generates a random string using a byte slice.
// Naive implementation for benchmarking comparison.
func RandomString(n int) string {
	b := make([]byte, 0, n)
	for range n {
		idx := rand.Int31() % 26
		b = append(b, alphabetArray[idx])
	}
	return string(b)
}

// RandomStringBuilder generates a random string using strings.Builder.
// Optimized implementation using standard library.
func RandomStringBuilder(n int) string {
	builder := strings.Builder{}
	for range n {
		idx := rand.Int31() % 26
		builder.WriteByte(alphabetArray[idx])
	}
	return builder.String()
}

// Contains checks if a string contains a given rune.
// This correctly handles Unicode by iterating over runes.
// Note: Initial implementation had a bug (using []byte instead of runes)
// that was discovered through fuzz testing!
func Contains(text string, char rune) bool {
	// Ranging over a string yields runes, not bytes
	// This correctly handles multi-byte UTF-8 characters
	for _, c := range text {
		if c == char {
			return true
		}
	}
	return false
}

// ============================================================================
// DEMONSTRATION 3: Benchmarking and Fuzzing
// ============================================================================

func demonstrateBenchmarkingAndFuzzing() {
	fmt.Println("\n=== DEMONSTRATION 3: Benchmarking & Fuzzing ===")

	fmt.Println("\n1. Benchmark Testing (Performance Measurement):")
	fmt.Println("   Compares RandomString vs RandomStringBuilder")
	fmt.Println("   Run with: go test -bench=. -benchmem")
	fmt.Println()
	fmt.Println("   Example output:")
	fmt.Println("   BenchmarkRandomString-8        30609   119318 ns/op")
	fmt.Println("   BenchmarkRandomStringBuilder-8 24721   120046 ns/op")
	fmt.Println()
	fmt.Println("   Interpretation:")
	fmt.Println("   • 30609/24721: Number of iterations")
	fmt.Println("   • 119318/120046 ns/op: Nanoseconds per operation")
	fmt.Println("   • Use -benchmem to see memory allocations")

	// Demo the functions
	s1 := RandomString(10)
	s2 := RandomStringBuilder(10)
	fmt.Printf("\n   RandomString(10): %s\n", s1)
	fmt.Printf("   RandomStringBuilder(10): %s\n", s2)

	fmt.Println("\n2. Fuzz Testing (Edge Case Discovery):")
	fmt.Println("   Tests Contains with random inputs")
	fmt.Println("   Run with: go test -fuzz=FuzzContains -fuzztime=30s")
	fmt.Println()
	fmt.Println("   How it works:")
	fmt.Println("   • Provide seed inputs (corpus)")
	fmt.Println("   • Fuzzer generates variations")
	fmt.Println("   • Looks for panics and failures")
	fmt.Println("   • Saves failing inputs for regression testing")

	// Demo the function
	fmt.Printf("\n   Contains(\"hello\", 'e'): %v\n", Contains("hello", 'e'))
	fmt.Printf("   Contains(\"hello\", 'z'): %v\n", Contains("hello", 'z'))
	fmt.Printf("   Contains(\"你好\", '好'): %v (Unicode!)\n", Contains("你好", '好'))

	fmt.Println("\n   Original bug: Used []byte instead of runes")
	fmt.Println("   Fuzzing found: Failed on multi-byte UTF-8 characters")
	fmt.Println("   Fix: Iterate over runes directly")
}

// ============================================================================
// TESTING CONCEPTS SUMMARY
// ============================================================================

func demonstrateTestingConcepts() {
	fmt.Println("\n=== TESTING CONCEPTS REFERENCE ===")

	fmt.Println("\n1. Test Function Types:")
	fmt.Println("   • TestXxx(t *testing.T)      - Unit tests")
	fmt.Println("   • BenchmarkXxx(b *testing.B) - Performance tests")
	fmt.Println("   • FuzzXxx(f *testing.F)      - Random input tests")
	fmt.Println("   • ExampleXxx()               - Documentation + tests")

	fmt.Println("\n2. Test File Conventions:")
	fmt.Println("   • Files: *_test.go")
	fmt.Println("   • Package: same_package (white box) or same_package_test (black box)")
	fmt.Println("   • Black box: Test as a user would (recommended)")
	fmt.Println("   • White box: Test internal implementation")

	fmt.Println("\n3. Failure Methods:")
	fmt.Println("   Non-fatal (continue test):")
	fmt.Println("   • t.Fail()         - Mark failed, no message")
	fmt.Println("   • t.Error(...)     - Log message and mark failed")
	fmt.Println("   • t.Errorf(...)    - Formatted message and mark failed")
	fmt.Println()
	fmt.Println("   Fatal (exit immediately):")
	fmt.Println("   • t.FailNow()      - Exit test, no message")
	fmt.Println("   • t.Fatal(...)     - Log message and exit")
	fmt.Println("   • t.Fatalf(...)    - Formatted message and exit")

	fmt.Println("\n4. Table-Driven Test Pattern:")
	fmt.Println("   type testCase struct {")
	fmt.Println("       name string")
	fmt.Println("       input string")
	fmt.Println("       want string")
	fmt.Println("   }")
	fmt.Println("   for _, tc := range testCases {")
	fmt.Println("       t.Run(tc.name, func(t *testing.T) {")
	fmt.Println("           got := Function(tc.input)")
	fmt.Println("           if got != tc.want { t.Error(...) }")
	fmt.Println("       })")
	fmt.Println("   }")

	fmt.Println("\n5. Code Coverage:")
	fmt.Println("   • go test -cover           - Show coverage %")
	fmt.Println("   • go test -coverprofile=c.out")
	fmt.Println("   • go tool cover -html=c.out - View in browser")
	fmt.Println("   • Aim for 70-80% (100% often impractical)")

	fmt.Println("\n6. Test Commands:")
	fmt.Println("   • go test              - Run all tests")
	fmt.Println("   • go test -v           - Verbose output")
	fmt.Println("   • go test -run TestName - Run specific test")
	fmt.Println("   • go test -short       - Skip long tests")
	fmt.Println("   • go test ./...        - Test all packages")

	fmt.Println("\n7. Best Practices:")
	fmt.Println("   ✓ Use black box testing (package_test)")
	fmt.Println("   ✓ Test error conditions")
	fmt.Println("   ✓ Use table-driven tests")
	fmt.Println("   ✓ Write descriptive test names")
	fmt.Println("   ✓ Use t.Helper() for helper functions")
	fmt.Println("   ✓ Test one thing per test function")
	fmt.Println("   ✗ Don't test trivial code")
	fmt.Println("   ✗ Don't make tests depend on each other")
	fmt.Println("   ✗ Don't use t.Fail() without logging")
}

// ============================================================================
// TDD WORKFLOW EXAMPLE
// ============================================================================

func demonstrateTDDWorkflow() {
	fmt.Println("\n=== TEST-DRIVEN DEVELOPMENT WORKFLOW ===")

	fmt.Println("\nStep 1: Define Requirements")
	fmt.Println("   Calculator needs:")
	fmt.Println("   • Basic operations (add, sub, mul, div)")
	fmt.Println("   • Memory storage")
	fmt.Println("   • Error handling (divide by zero, invalid operation)")

	fmt.Println("\nStep 2: Write Tests FIRST")
	fmt.Println("   • TestErrDivideByZero - verify error returned")
	fmt.Println("   • TestErrUnknownOperation - verify invalid ops fail")
	fmt.Println("   • TestCalculate - table-driven test for all ops")
	fmt.Println("   • TestMemory - verify store/recall functionality")

	fmt.Println("\nStep 3: Run Tests (they all fail!)")
	fmt.Println("   $ go test")
	fmt.Println("   --- FAIL: TestErrDivideByZero")
	fmt.Println("   --- FAIL: TestCalculate")
	fmt.Println("   This is EXPECTED in TDD!")

	fmt.Println("\nStep 4: Implement Minimum Code to Pass")
	fmt.Println("   • Add Calculate method")
	fmt.Println("   • Add Store/Recall methods")
	fmt.Println("   • Handle errors")

	fmt.Println("\nStep 5: Run Tests Again")
	fmt.Println("   $ go test")
	fmt.Println("   PASS")
	fmt.Println("   ✓ All tests passing!")

	fmt.Println("\nStep 6: Check Coverage")
	fmt.Println("   $ go test -cover")
	fmt.Println("   coverage: 93.8% of statements")
	fmt.Println("   (Oops, missing error recovery test!)")

	fmt.Println("\nStep 7: Add Missing Tests")
	fmt.Println("   • Test Store after error (shouldn't update)")
	fmt.Println("   • Test Store recovery (should work again)")

	fmt.Println("\nStep 8: Discover Bug Through Testing!")
	fmt.Println("   Test fails: Store doesn't work after error")
	fmt.Println("   Root cause: Calculate() doesn't reset c.err")
	fmt.Println("   Fix: Add c.err = nil on successful calculation")

	fmt.Println("\nStep 9: Verify Fix")
	fmt.Println("   $ go test -cover")
	fmt.Println("   PASS")
	fmt.Println("   coverage: 100.0% of statements")
	fmt.Println("   ✓ All tests passing, complete coverage!")

	fmt.Println("\nTDD Benefits Demonstrated:")
	fmt.Println("   ✓ Clear requirements from the start")
	fmt.Println("   ✓ Better API design (used it before building it)")
	fmt.Println("   ✓ Complete test coverage")
	fmt.Println("   ✓ Found bugs early (Store recovery)")
	fmt.Println("   ✓ Confidence to refactor")
}

// ============================================================================
// MAIN FUNCTION
// ============================================================================

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                               ║")
	fmt.Println("║         CHAPTER 8: TESTING AND TOOLING IN GO                  ║")
	fmt.Println("║                                                               ║")
	fmt.Println("║  A comprehensive reference for Go's testing ecosystem         ║")
	fmt.Println("║                                                               ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")

	demonstrateCalculator()
	demonstrateHTTPTesting()
	demonstrateBenchmarkingAndFuzzing()
	demonstrateTestingConcepts()
	demonstrateTDDWorkflow()

	fmt.Println("\n" + strings.Repeat("=", 65))
	fmt.Println("NEXT STEPS:")
	fmt.Println(strings.Repeat("=", 65))
	fmt.Println()
	fmt.Println("1. Run the tests:")
	fmt.Println("   $ go test -v")
	fmt.Println()
	fmt.Println("2. Check coverage:")
	fmt.Println("   $ go test -cover")
	fmt.Println("   $ go test -coverprofile=coverage.out")
	fmt.Println("   $ go tool cover -html=coverage.out")
	fmt.Println()
	fmt.Println("3. Run benchmarks:")
	fmt.Println("   $ go test -bench=. -benchmem")
	fmt.Println()
	fmt.Println("4. Run fuzz tests:")
	fmt.Println("   $ go test -fuzz=FuzzContains -fuzztime=30s")
	fmt.Println()
	fmt.Println("See main_test.go for all runnable tests!")
	fmt.Println()
}
