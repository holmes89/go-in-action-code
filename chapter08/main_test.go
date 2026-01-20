package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// ============================================================================
// CALCULATOR TESTS - TDD EXAMPLE
// ============================================================================

// TestErrDivideByZero verifies divide by zero error.
func TestErrDivideByZero(t *testing.T) {
	c := NewCalculator()
	got, err := c.Calculate("div", 1, 0)
	if !errors.Is(err, ErrDivByZero) {
		t.Errorf("expected ErrDivByZero - got: %v", err)
	}
	if got != 0 {
		t.Errorf("error result should be 0 - got: %d", got)
	}
}

// TestErrUnknownOperation verifies invalid operation error.
func TestErrUnknownOperation(t *testing.T) {
	c := NewCalculator()
	for _, op := range []string{"", "mult", "asd", "ADD", "divide"} {
		t.Logf("test invalid operation %q", op)
		got, err := c.Calculate(op, 1, 1)
		if !errors.Is(err, ErrUnknownOperation) {
			t.Errorf("expected ErrUnknownOperation - got: %v", err)
		}
		if got != 0 {
			t.Errorf("error result should be 0 - got: %d", got)
		}
	}
}

// TestMemory verifies Store and Recall functionality.
func TestMemory(t *testing.T) {
	c := NewCalculator()

	// Initial memory should be 0
	got := c.Recall()
	if got != 0 {
		t.Errorf("initial memory value should be 0 - got: %d", got)
	}

	// Calculate and store
	t.Logf("Calculate - div, 12, 3")
	got, err := c.Calculate("div", 12, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 4 {
		t.Errorf("incorrect result - want: 4, got: %d", got)
	}

	// Verify stored value
	c.Store()
	got = c.Recall()
	if got != 4 {
		t.Errorf("recall - want: 4, got: %d", got)
	}

	// Store should NOT update after error
	_, _ = c.Calculate("div", 0, 0)
	c.Store()
	got = c.Recall()
	if got != 4 {
		t.Errorf("storage should not update value after invalid result")
	}

	// Storage should work again after error
	_, err = c.Calculate("mul", 4, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c.Store()
	got = c.Recall()
	if got != 12 {
		t.Errorf("recall after storage recovery - want: 12, got: %d", got)
	}
}

// TestCalculate tests all operations with table-driven approach.
func TestCalculate(t *testing.T) {
	type testCase struct {
		name string
		op   string
		a    int
		b    int
		want int
	}

	testCases := []testCase{
		{name: "divide", op: "div", a: 10, b: 5, want: 2},
		{name: "divide_negative", op: "div", a: -12, b: 4, want: -3},
		{name: "add_zero", op: "add", a: 1, b: 0, want: 1},
		{name: "add_negative", op: "add", a: 2, b: -4, want: -2},
		{name: "add_large", op: "add", a: 100, b: 200, want: 300},
		{name: "multiply_zero", op: "mul", a: 100, b: 0, want: 0},
		{name: "multiply_negative", op: "mul", a: 2, b: -4, want: -8},
		{name: "multiply", op: "mul", a: 5, b: 5, want: 25},
		{name: "subtract", op: "sub", a: 20, b: 7, want: 13},
		{name: "subtract_negative", op: "sub", a: 5, b: -4, want: 9},
		{name: "subtract_zero", op: "sub", a: 0, b: 10, want: -10},
	}

	c := NewCalculator()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := c.Calculate(tc.op, tc.a, tc.b)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tc.want {
				t.Errorf("want: %d, got: %d", tc.want, got)
			}
		})
	}
}

// ============================================================================
// HTTP HANDLER TESTS
// ============================================================================

// TestHandler tests HTTP handlers with ResponseRecorder.
func TestHandler(t *testing.T) {
	type testCase struct {
		name             string
		handler          http.HandlerFunc
		method           string
		body             []byte
		expectedStatus   int
		expectedResponse string
	}

	testCases := []testCase{
		{
			name:             "GET_hello",
			handler:          HelloGet,
			method:           "GET",
			body:             nil,
			expectedStatus:   http.StatusOK,
			expectedResponse: "hello",
		},
		{
			name:             "POST_hello_with_body",
			handler:          HelloPost,
			method:           "POST",
			body:             []byte("Go in Action"),
			expectedStatus:   http.StatusCreated,
			expectedResponse: "hello Go in Action",
		},
		{
			name:             "POST_hello_empty",
			handler:          HelloPost,
			method:           "POST",
			body:             nil,
			expectedStatus:   http.StatusCreated,
			expectedResponse: "hello ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(tc.handler)
			rr := httptest.NewRecorder()
			req, err := http.NewRequest(tc.method, "", bytes.NewBuffer(tc.body))
			if err != nil {
				t.Fatalf("unable to create request: %s", err)
			}

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d but received %d",
					tc.expectedStatus, rr.Code)
			}

			respBody, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatalf("unable to read response: %s", err)
			}

			if string(respBody) != tc.expectedResponse {
				t.Errorf("expected %q received %q",
					tc.expectedResponse, string(respBody))
			}
		})
	}
}

// TestServer tests HTTP handlers with httptest.Server.
func TestServer(t *testing.T) {
	type testCase struct {
		name             string
		path             string
		method           string
		body             []byte
		expectedStatus   int
		expectedResponse string
	}

	testCases := []testCase{
		{
			name:             "GET_hello",
			path:             "/hello",
			method:           "GET",
			body:             nil,
			expectedStatus:   http.StatusOK,
			expectedResponse: "hello",
		},
		{
			name:             "POST_hello",
			path:             "/hello",
			method:           "POST",
			body:             []byte("Go in Action"),
			expectedStatus:   http.StatusCreated,
			expectedResponse: "hello Go in Action",
		},
		{
			name:             "PUT_unsupported",
			path:             "/hello",
			method:           "PUT",
			body:             []byte("bad request"),
			expectedStatus:   http.StatusNotFound,
			expectedResponse: "invalid hello method",
		},
		{
			name:             "GET_wrong_path",
			path:             "/goodbye",
			method:           "GET",
			body:             nil,
			expectedStatus:   http.StatusNotFound,
			expectedResponse: "404 page not found\n",
		},
	}

	mux := http.NewServeMux()
	HelloHandler(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fullURL, err := url.JoinPath(server.URL, tc.path)
			if err != nil {
				t.Fatalf("unable to create URL: %s", err)
			}

			req, err := http.NewRequest(tc.method, fullURL, bytes.NewBuffer(tc.body))
			if err != nil {
				t.Fatalf("unable to create request: %s", err)
			}

			resp, err := server.Client().Do(req)
			if err != nil {
				t.Fatalf("request failed: %s", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("expected status %d but received %d",
					tc.expectedStatus, resp.StatusCode)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("unable to read response: %s", err)
			}

			if string(respBody) != tc.expectedResponse {
				t.Errorf("expected %q received %q",
					tc.expectedResponse, string(respBody))
			}
		})
	}
}

// ============================================================================
// STRING UTILITY TESTS
// ============================================================================

// TestContains validates the Contains function.
func TestContains(t *testing.T) {
	testCases := []struct {
		name string
		text string
		char rune
		want bool
	}{
		{name: "found", text: "hello", char: 'e', want: true},
		{name: "not_found", text: "hello", char: 'z', want: false},
		{name: "empty", text: "", char: 'a', want: false},
		{name: "unicode", text: "你好世界", char: '好', want: true},
		{name: "emoji", text: "Hello 🎉!", char: '🎉', want: true},
		{name: "first_char", text: "hello", char: 'h', want: true},
		{name: "last_char", text: "hello", char: 'o', want: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Contains(tc.text, tc.char)
			if got != tc.want {
				t.Errorf("Contains(%q, %q) = %v, want %v",
					tc.text, tc.char, got, tc.want)
			}
		})
	}
}

// TestRandomString validates both string generation implementations.
func TestRandomString(t *testing.T) {
	testCases := []struct {
		name   string
		length int
	}{
		{name: "zero", length: 0},
		{name: "one", length: 1},
		{name: "ten", length: 10},
		{name: "hundred", length: 100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result1 := RandomString(tc.length)
			result2 := RandomStringBuilder(tc.length)

			if len(result1) != tc.length {
				t.Errorf("RandomString: wrong length - want: %d, got: %d",
					tc.length, len(result1))
			}
			if len(result2) != tc.length {
				t.Errorf("RandomStringBuilder: wrong length - want: %d, got: %d",
					tc.length, len(result2))
			}

			for _, c := range result1 {
				if c < 'a' || c > 'z' {
					t.Errorf("RandomString: invalid character %c", c)
				}
			}
			for _, c := range result2 {
				if c < 'a' || c > 'z' {
					t.Errorf("RandomStringBuilder: invalid character %c", c)
				}
			}
		})
	}
}

// ============================================================================
// BENCHMARK TESTS
// ============================================================================

// BenchmarkRandomString benchmarks the byte slice implementation.
func BenchmarkRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := RandomString(i % 1000)
		if len(res) != i%1000 {
			b.Error("result has the wrong length")
		}
	}
}

// BenchmarkRandomStringBuilder benchmarks the strings.Builder implementation.
func BenchmarkRandomStringBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := RandomStringBuilder(i % 1000)
		if len(res) != i%1000 {
			b.Error("result has the wrong length")
		}
	}
}

// ============================================================================
// FUZZ TESTS
// ============================================================================

// FuzzContains fuzzes the Contains function with random inputs.
func FuzzContains(f *testing.F) {
	// Add seed corpus
	seeds := []struct {
		word string
		char rune
	}{
		{word: "hello world", char: 'e'},
		{word: "goodbye", char: 'e'},
		{word: "", char: 'o'},
		{word: "a longer example", char: 'z'},
		{word: "unicode: 你好", char: '好'},
		{word: "emoji: 🎉", char: '🎉'},
	}

	for _, s := range seeds {
		f.Add(s.word, s.char)
	}

	// Fuzz test - compare with standard library (oracle testing)
	f.Fuzz(func(t *testing.T, in string, r rune) {
		got := Contains(in, r)
		want := strings.Contains(in, string(r))

		if got != want {
			t.Errorf("Contains(%q, %q) = %v, want %v", in, r, got, want)
		}
	})
}

// ============================================================================
// EXAMPLE TESTS (appear in godoc)
// ============================================================================

// ExampleCalculator demonstrates basic calculator usage.
func ExampleCalculator() {
	c := NewCalculator()
	result, _ := c.Calculate("add", 5, 3)
	fmt.Println(result)

	result, _ = c.Calculate("mul", 4, 7)
	fmt.Println(result)

	// Output:
	// 8
	// 28
}

// ExampleCalculator_Store demonstrates the memory feature.
func ExampleCalculator_Store() {
	c := NewCalculator()
	c.Calculate("add", 10, 5)
	c.Store()
	result := c.Recall()
	fmt.Println(result)

	// Output:
	// 15
}
