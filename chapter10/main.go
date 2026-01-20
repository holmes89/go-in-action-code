package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// =============================================================================
// CHAPTER 10: THE STANDARD LIBRARY
// =============================================================================
//
// This file demonstrates the essential packages in Go's standard library
// by building a complete JSON/XML validation and conversion CLI tool.
//
// Covered packages:
// - fmt: Formatted I/O (printing, formatting, scanning)
// - flag: Command-line argument parsing
// - io: Core I/O interfaces (Reader, Writer, Closer)
// - os: Operating system functionality (files, stdin/stdout, env vars)
// - encoding/json & encoding/xml: Structured data serialization
// - bufio: Buffered I/O for performance
// - regexp: Regular expression pattern matching
// - strings: String manipulation utilities
// - strconv: String/type conversions
//
// The tool can:
// - Read JSON or XML from files or stdin
// - Validate the structure against a Config schema
// - Convert between JSON and XML formats
// - Output to files or stdout
// - Handle configuration via flags and environment variables

// =============================================================================
// DATA STRUCTURES
// =============================================================================

// Config represents a server configuration that can be serialized to both
// JSON and XML formats. Notice the struct tags that map Go fields to both
// JSON keys and XML elements/attributes.
type Config struct {
	XMLName   xml.Name `json:"-" xml:"configuration"`                  // Root element name for XML
	Name      string   `json:"name" xml:"name,attr"`                   // Application name (XML attribute)
	Version   string   `json:"version" xml:"version,attr"`             // Version string (XML attribute)
	Namespace string   `json:"-" xml:"xmlns,attr,omitempty"`           // XML namespace (optional)
	Debug     bool     `json:"debug" xml:"debug"`                      // Debug mode flag
	Port      int      `json:"port" xml:"port"`                        // Server port number
	Host      string   `json:"host" xml:"host,omitempty"`              // Hostname (omit if empty)
	LogLevel  string   `json:"log_level,omitempty" xml:"log_level,omitempty"` // Logging level
	Comment   string   `json:"-" xml:",comment"`                       // XML comment only
}

// ValidationError represents a validation error with context about which
// field failed and why. This provides better error messages to users.
type ValidationError struct {
	Field  string
	Reason string
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("validation error in field '%s': %s", ve.Field, ve.Reason)
}

// =============================================================================
// VALIDATION LOGIC
// =============================================================================

// validateConfig performs comprehensive validation on a Config struct using
// regular expressions and business logic. This demonstrates the regexp package
// and best practices for input validation.
func validateConfig(cfg *Config) []ValidationError {
	var errors []ValidationError

	// Validate name: must be non-empty and alphanumeric with hyphens/underscores
	// Pattern: ^[a-zA-Z0-9_-]+$
	// ^ = start of string, $ = end of string, + = one or more
	namePattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if cfg.Name == "" {
		errors = append(errors, ValidationError{
			Field:  "name",
			Reason: "name cannot be empty",
		})
	} else if !namePattern.MatchString(cfg.Name) {
		errors = append(errors, ValidationError{
			Field:  "name",
			Reason: "name must be alphanumeric with hyphens or underscores",
		})
	}

	// Validate version: should follow semantic versioning pattern
	// Pattern: ^\d+\.\d+\.\d+$ matches "1.0.0", "2.3.4", etc.
	// \d = digit, \. = literal period (escaped)
	versionPattern := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	if cfg.Version == "" {
		errors = append(errors, ValidationError{
			Field:  "version",
			Reason: "version cannot be empty",
		})
	} else if !versionPattern.MatchString(cfg.Version) {
		errors = append(errors, ValidationError{
			Field:  "version",
			Reason: "version must follow semantic versioning (e.g., 1.0.0)",
		})
	}

	// Validate port: must be in valid range (1-65535)
	if cfg.Port < 1 || cfg.Port > 65535 {
		errors = append(errors, ValidationError{
			Field:  "port",
			Reason: fmt.Sprintf("port must be between 1 and 65535, got %d", cfg.Port),
		})
	}

	// Validate host: should be a valid hostname or URL if provided
	if cfg.Host != "" {
		// Try to parse as URL to validate format
		if _, err := url.ParseRequestURI("http://" + cfg.Host); err != nil {
			errors = append(errors, ValidationError{
				Field:  "host",
				Reason: fmt.Sprintf("invalid host format: %v", err),
			})
		}
	}

	// Validate log level: should be one of the standard levels
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true,
	}
	if cfg.LogLevel != "" && !validLogLevels[strings.ToLower(cfg.LogLevel)] {
		errors = append(errors, ValidationError{
			Field:  "log_level",
			Reason: "log_level must be one of: debug, info, warn, error, fatal",
		})
	}

	return errors
}

// =============================================================================
// FORMAT DETECTION
// =============================================================================

// detectFormat examines the input to determine if it's JSON or XML.
// This demonstrates bufio.Reader's Peek functionality - looking ahead without
// consuming bytes. This is crucial for format auto-detection.
func detectFormat(r io.Reader) (string, *bufio.Reader, error) {
	br := bufio.NewReader(r)

	// Peek at first few bytes without consuming them
	// This allows us to detect format and still read the full content
	firstBytes, err := br.Peek(10)
	if err != nil && err != io.EOF {
		return "", br, fmt.Errorf("failed to peek at input: %w", err)
	}

	// Trim whitespace to find first meaningful character
	trimmed := strings.TrimSpace(string(firstBytes))
	if trimmed == "" {
		return "", br, fmt.Errorf("input is empty")
	}

	// Check for XML declaration or opening tag
	if strings.HasPrefix(trimmed, "<?xml") || strings.HasPrefix(trimmed, "<") {
		return "xml", br, nil
	}

	// Check for JSON object or array
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		return "json", br, nil
	}

	return "", br, fmt.Errorf("unknown format: input does not appear to be JSON or XML")
}

// detectFormatFromFilename checks file extension to guess format.
// This demonstrates the strings package's suffix checking.
func detectFormatFromFilename(filename string) string {
	filename = strings.ToLower(strings.TrimSpace(filename))
	
	if strings.HasSuffix(filename, ".json") {
		return "json"
	}
	if strings.HasSuffix(filename, ".xml") {
		return "xml"
	}
	
	return "unknown"
}

// =============================================================================
// PARSING AND CONVERSION
// =============================================================================

// parseJSON reads JSON data and unmarshals it into a Config struct.
// This demonstrates encoding/json package usage with error handling.
func parseJSON(data []byte) (*Config, error) {
	var cfg Config
	
	// Unmarshal the JSON into our struct
	// The json package uses struct tags to map JSON keys to Go fields
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	
	return &cfg, nil
}

// parseXML reads XML data and unmarshals it into a Config struct.
// This demonstrates encoding/xml package usage, which is similar to JSON
// but supports attributes, namespaces, and more complex structures.
func parseXML(data []byte) (*Config, error) {
	var cfg Config
	
	// Unmarshal the XML into our struct
	// The xml package uses struct tags to map elements/attributes to Go fields
	if err := xml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid XML: %w", err)
	}
	
	return &cfg, nil
}

// convertToJSON converts a Config struct to formatted JSON.
// Uses MarshalIndent for human-readable output with consistent formatting.
func convertToJSON(cfg *Config) ([]byte, error) {
	// MarshalIndent produces nicely formatted JSON with 2-space indentation
	// First param is the value to marshal
	// Second param is the prefix for each line (usually empty)
	// Third param is the indentation string (spaces or tabs)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}
	
	return data, nil
}

// convertToXML converts a Config struct to formatted XML.
// Includes XML declaration and proper indentation for readability.
func convertToXML(cfg *Config) ([]byte, error) {
	// MarshalIndent produces nicely formatted XML
	data, err := xml.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to XML: %w", err)
	}
	
	// Add XML declaration at the beginning
	// xml.MarshalIndent doesn't include this automatically
	xmlWithDeclaration := append([]byte(xml.Header), data...)
	
	return xmlWithDeclaration, nil
}

// =============================================================================
// FILE I/O OPERATIONS
// =============================================================================

// openValidatedInput opens a file and validates it's suitable for reading.
// This demonstrates os.Stat for file metadata and proper error handling.
// Returns an io.ReadCloser so the caller can close it when done.
func openValidatedInput(path string) (io.ReadCloser, error) {
	// Get file metadata without opening the file
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file does not exist: %s", path)
		}
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Validate it's a regular file, not a directory or special file
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("not a regular file: %s", path)
	}

	// Check if file is empty (though we'll handle this gracefully)
	if info.Size() == 0 {
		return nil, fmt.Errorf("file is empty: %s", path)
	}

	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// createAtomicOutput creates a temporary file for atomic writes.
// This pattern prevents data corruption if the program crashes mid-write.
// The caller should call finalizeAtomicOutput when done writing.
func createAtomicOutput(outputPath string) (*os.File, error) {
	// Create temporary file in the same directory as the target
	// This ensures the temp file is on the same filesystem for atomic rename
	tempFile, err := os.CreateTemp(".", ".tmp-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	return tempFile, nil
}

// finalizeAtomicOutput moves the temporary file to its final location.
// This is atomic on most filesystems, preventing partial writes.
func finalizeAtomicOutput(tempPath, finalPath string) error {
	// Rename is atomic on most POSIX systems
	if err := os.Rename(tempPath, finalPath); err != nil {
		// Clean up temp file on failure
		os.Remove(tempPath)
		return fmt.Errorf("failed to finalize output: %w", err)
	}
	return nil
}

// =============================================================================
// CORE PROCESSING LOGIC
// =============================================================================

// process is the main processing function that handles format detection,
// parsing, validation, conversion, and output. It demonstrates how to use
// io.Reader and io.Writer interfaces to work with any input/output source.
func process(input io.Reader, output io.Writer, targetFormat string, validate bool) error {
	// Use buffered reader for efficient reading and format detection
	br := bufio.NewReader(input)
	
	// Detect the input format by peeking at the first few bytes
	inputFormat, br, err := detectFormat(br)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Detected input format: %s\n", inputFormat)

	// Read all input data into memory
	// For very large files, you might want to use streaming instead
	data, err := io.ReadAll(br)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Parse based on detected format
	var cfg *Config
	switch inputFormat {
	case "json":
		cfg, err = parseJSON(data)
		if err != nil {
			return err
		}
	case "xml":
		cfg, err = parseXML(data)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported input format: %s", inputFormat)
	}

	// Validate the configuration if requested
	if validate {
		if validationErrors := validateConfig(cfg); len(validationErrors) > 0 {
			fmt.Fprintf(os.Stderr, "\nValidation errors found:\n")
			for _, ve := range validationErrors {
				fmt.Fprintf(os.Stderr, "  - %s\n", ve.Error())
			}
			return fmt.Errorf("validation failed with %d error(s)", len(validationErrors))
		}
		fmt.Fprintf(os.Stderr, "Validation: PASSED\n")
	}

	// Convert to target format
	var outputData []byte
	if targetFormat == "auto" {
		// If auto, convert to the opposite format
		if inputFormat == "json" {
			targetFormat = "xml"
		} else {
			targetFormat = "json"
		}
	}

	switch targetFormat {
	case "json":
		outputData, err = convertToJSON(cfg)
	case "xml":
		outputData, err = convertToXML(cfg)
	default:
		return fmt.Errorf("unsupported output format: %s", targetFormat)
	}

	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Output format: %s\n", targetFormat)

	// Write output using buffered writer for efficiency
	bw := bufio.NewWriter(output)
	defer bw.Flush() // Ensure all buffered data is written

	if _, err := bw.Write(outputData); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	// Add newline for better terminal output
	if _, err := bw.WriteString("\n"); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	return nil
}

// =============================================================================
// DEMONSTRATION FUNCTIONS
// =============================================================================

// demonstrateFmtPackage shows various ways to use the fmt package for
// formatted output and input. This includes Print, Println, Printf, Sprint,
// Sprintf, and formatting with various verbs.
func demonstrateFmtPackage() {
	fmt.Println("\n=== DEMONSTRATION 1: fmt Package - Formatted I/O ===")
	
	// Basic printing
	fmt.Print("Print: no newline")
	fmt.Print(" continues on same line\n")
	fmt.Println("Println: automatic newline")
	
	// Formatted printing with Printf
	app := "myapp"
	version := "1.0.0"
	port := 8080
	debug := true
	
	fmt.Printf("Printf: %s version %s running on port %d (debug=%t)\n", 
		app, version, port, debug)
	
	// Common format verbs
	fmt.Println("\nFormat verbs:")
	fmt.Printf("  %%v (default):     %v\n", Config{Name: "test", Port: 8080})
	fmt.Printf("  %%+v (with fields): %+v\n", Config{Name: "test", Port: 8080})
	fmt.Printf("  %%#v (Go syntax):   %#v\n", Config{Name: "test", Port: 8080})
	fmt.Printf("  %%T (type):         %T\n", Config{})
	fmt.Printf("  %%d (decimal):      %d\n", 42)
	fmt.Printf("  %%b (binary):       %b\n", 42)
	fmt.Printf("  %%x (hex lower):    %x\n", 42)
	fmt.Printf("  %%X (hex upper):    %X\n", 42)
	fmt.Printf("  %%f (float):        %f\n", 3.14159)
	fmt.Printf("  %%.2f (precision):   %.2f\n", 3.14159)
	fmt.Printf("  %%s (string):       %s\n", "hello")
	fmt.Printf("  %%q (quoted):       %q\n", "hello")
	
	// Sprint functions - return strings instead of printing
	message := fmt.Sprint(app, " v", version)
	formatted := fmt.Sprintf("%s - version %s", app, version)
	
	fmt.Println("\nSprint results:")
	fmt.Println("  Sprint:", message)
	fmt.Println("  Sprintf:", formatted)
	
	// Creating formatted errors
	err := fmt.Errorf("port %d is out of range (1-65535)", 99999)
	fmt.Printf("  Error: %v\n", err)
}

// demonstrateFlagPackage shows how to use the flag package for command-line
// argument parsing. Note: This runs in demo mode with preset values.
func demonstrateFlagPackage() {
	fmt.Println("\n=== DEMONSTRATION 2: flag Package - Command-line Arguments ===")
	
	fmt.Println("The flag package parses command-line arguments:")
	fmt.Println("  -help          Show help message")
	fmt.Println("  -f <file>      Input file path")
	fmt.Println("  -o <file>      Output file path")
	fmt.Println("  -format <fmt>  Target format (json, xml, auto)")
	fmt.Println("  -validate      Validate configuration")
	fmt.Println("  -banner <msg>  Display custom banner")
	fmt.Println("\nFlags are defined with flag.String(), flag.Bool(), flag.Int(), etc.")
	fmt.Println("flag.Parse() processes the command-line arguments")
	fmt.Println("flag.Args() returns non-flag arguments")
}

// demonstrateIoPackage shows the core io package interfaces and utilities.
func demonstrateIoPackage() {
	fmt.Println("\n=== DEMONSTRATION 3: io Package - Core I/O Interfaces ===")
	
	fmt.Println("The io package defines fundamental interfaces:")
	fmt.Println("  - io.Reader: anything you can read from")
	fmt.Println("  - io.Writer: anything you can write to")
	fmt.Println("  - io.Closer: anything that needs cleanup")
	fmt.Println("  - io.ReadWriter, io.ReadCloser, etc.: composed interfaces")
	
	// Demonstrate io.Copy with string reader and stdout
	jsonData := `{"name":"demo","version":"1.0.0","port":8080}`
	reader := strings.NewReader(jsonData)
	
	fmt.Println("\nExample: io.Copy transfers data from Reader to Writer")
	fmt.Print("Content: ")
	io.Copy(os.Stdout, reader)
	fmt.Println()
	
	// Demonstrate io.ReadAll
	reader2 := strings.NewReader("Read all at once")
	data, _ := io.ReadAll(reader2)
	fmt.Printf("io.ReadAll result: %s\n", data)
	
	// Demonstrate io.LimitReader
	reader3 := strings.NewReader("This is a long string that will be limited")
	limitedReader := io.LimitReader(reader3, 20) // Only read 20 bytes
	limitedData, _ := io.ReadAll(limitedReader)
	fmt.Printf("io.LimitReader (20 bytes): %s...\n", limitedData)
}

// demonstrateOsPackage shows os package functionality for files, streams,
// environment variables, and exit codes.
func demonstrateOsPackage() {
	fmt.Println("\n=== DEMONSTRATION 4: os Package - Operating System Interface ===")
	
	// Standard streams
	fmt.Println("Standard streams:")
	fmt.Println("  - os.Stdin:  standard input (io.Reader)")
	fmt.Println("  - os.Stdout: standard output (io.Writer)")
	fmt.Println("  - os.Stderr: standard error (io.Writer)")
	
	// Write to stderr to demonstrate the difference
	fmt.Fprintln(os.Stderr, "  [This message went to stderr]")
	
	// Environment variables
	fmt.Println("\nEnvironment variables:")
	
	// Set some example environment variables
	os.Setenv("APP_NAME", "myapp")
	os.Setenv("APP_DEBUG", "true")
	os.Setenv("APP_PORT", "8080")
	
	appName := os.Getenv("APP_NAME")
	fmt.Printf("  APP_NAME: %s\n", appName)
	
	// LookupEnv tells you if the variable exists
	if debug, exists := os.LookupEnv("APP_DEBUG"); exists {
		fmt.Printf("  APP_DEBUG: %s (exists)\n", debug)
	}
	
	// File operations
	fmt.Println("\nFile operations:")
	fmt.Println("  - os.Open(): open for reading")
	fmt.Println("  - os.Create(): create or truncate for writing")
	fmt.Println("  - os.OpenFile(): open with custom flags")
	fmt.Println("  - os.Stat(): get file metadata")
	fmt.Println("  - os.Remove(): delete file")
	fmt.Println("  - os.Rename(): move/rename file")
	
	// Exit codes
	fmt.Println("\nExit codes:")
	fmt.Println("  - os.Exit(0): success")
	fmt.Println("  - os.Exit(1): general error")
	fmt.Println("  - os.Exit(2): misuse of command")
	fmt.Println("  (We don't call os.Exit() in this demo)")
}

// demonstrateEncodingPackage shows JSON and XML encoding/decoding with the
// encoding/json and encoding/xml packages.
func demonstrateEncodingPackage() {
	fmt.Println("\n=== DEMONSTRATION 5: encoding Package - JSON and XML ===")
	
	// Create a sample configuration
	cfg := Config{
		Name:     "myapp",
		Version:  "1.0.0",
		Debug:    false,
		Port:     8080,
		Host:     "localhost",
		LogLevel: "info",
		Comment:  "Production configuration",
	}
	
	// Marshal to JSON
	jsonData, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println("JSON output:")
	fmt.Println(string(jsonData))
	
	// Marshal to XML
	xmlData, _ := xml.MarshalIndent(cfg, "", "  ")
	fmt.Println("\nXML output:")
	fmt.Println(xml.Header + string(xmlData))
	
	// Unmarshal JSON back to struct
	var cfgFromJSON Config
	json.Unmarshal(jsonData, &cfgFromJSON)
	fmt.Printf("\nUnmarshaled from JSON: name=%s, port=%d\n", 
		cfgFromJSON.Name, cfgFromJSON.Port)
	
	// Unmarshal XML back to struct
	var cfgFromXML Config
	xml.Unmarshal(xmlData, &cfgFromXML)
	fmt.Printf("Unmarshaled from XML: name=%s, port=%d\n", 
		cfgFromXML.Name, cfgFromXML.Port)
	
	fmt.Println("\nStruct tags control marshaling:")
	fmt.Println("  json:\"name\"       - JSON key mapping")
	fmt.Println("  xml:\"name,attr\"   - XML attribute")
	fmt.Println("  json:\"-\"          - Skip this field")
	fmt.Println("  json:\"omitempty\" - Omit if empty")
}

// demonstrateBufioPackage shows buffered I/O for performance.
func demonstrateBufioPackage() {
	fmt.Println("\n=== DEMONSTRATION 6: bufio Package - Buffered I/O ===")
	
	fmt.Println("Buffered I/O reduces system calls by batching operations")
	
	// Reading with bufio.Reader
	input := "Line 1\nLine 2\nLine 3\n"
	reader := bufio.NewReader(strings.NewReader(input))
	
	fmt.Println("\nReading line by line with bufio.Reader:")
	lineNum := 1
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Printf("  Line %d: %s", lineNum, line)
		lineNum++
	}
	
	// Peeking without consuming
	reader2 := bufio.NewReader(strings.NewReader("Peek at this"))
	peeked, _ := reader2.Peek(4)
	fmt.Printf("\nPeek (without consuming): %s\n", peeked)
	
	// Read normally - still gets the peeked data
	firstWord, _ := reader2.ReadString(' ')
	fmt.Printf("Then read normally: %s\n", firstWord)
	
	// Writing with bufio.Writer
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	
	fmt.Println("\nWriting with bufio.Writer:")
	writer.WriteString("Buffered ")
	writer.WriteString("writing ")
	writer.WriteString("is ")
	writer.WriteString("efficient")
	writer.Flush() // Must flush to ensure all data is written
	
	fmt.Printf("  Result: %s\n", buf.String())
	fmt.Println("  (Remember to call Flush() to ensure data is written!)")
}

// demonstrateRegexpPackage shows pattern matching with regular expressions.
func demonstrateRegexpPackage() {
	fmt.Println("\n=== DEMONSTRATION 7: regexp Package - Pattern Matching ===")
	
	// Compile patterns (do this once, reuse many times)
	namePattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	versionPattern := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	
	// Test various strings
	names := []string{"myapp", "my-app", "my_app", "my app", "123app"}
	versions := []string{"1.0.0", "2.3.4", "1.0", "v1.0.0"}
	emails := []string{"user@example.com", "invalid.email", "test@test"}
	
	fmt.Println("Name validation (alphanumeric, hyphens, underscores):")
	for _, name := range names {
		match := namePattern.MatchString(name)
		fmt.Printf("  %-15s %v\n", name, match)
	}
	
	fmt.Println("\nVersion validation (semantic versioning X.Y.Z):")
	for _, version := range versions {
		match := versionPattern.MatchString(version)
		fmt.Printf("  %-15s %v\n", version, match)
	}
	
	fmt.Println("\nEmail validation:")
	for _, email := range emails {
		match := emailPattern.MatchString(email)
		fmt.Printf("  %-20s %v\n", email, match)
	}
	
	// Capturing groups
	logPattern := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2}) (\d{2}):(\d{2}):(\d{2})`)
	logLine := "2026-01-20 14:30:45 Server started"
	matches := logPattern.FindStringSubmatch(logLine)
	
	if matches != nil {
		fmt.Println("\nCapturing groups from log line:")
		fmt.Printf("  Full match: %s\n", matches[0])
		fmt.Printf("  Year: %s, Month: %s, Day: %s\n", matches[1], matches[2], matches[3])
		fmt.Printf("  Hour: %s, Minute: %s, Second: %s\n", matches[4], matches[5], matches[6])
	}
}

// demonstrateStringsPackage shows string manipulation utilities.
func demonstrateStringsPackage() {
	fmt.Println("\n=== DEMONSTRATION 8: strings Package - String Manipulation ===")
	
	// Case conversion
	original := "  Hello World  "
	fmt.Printf("Original: %q\n", original)
	fmt.Printf("ToLower: %q\n", strings.ToLower(original))
	fmt.Printf("ToUpper: %q\n", strings.ToUpper(original))
	fmt.Printf("TrimSpace: %q\n", strings.TrimSpace(original))
	
	// Prefix and suffix
	filename := "config.json"
	fmt.Printf("\nFilename: %s\n", filename)
	fmt.Printf("  HasPrefix(\"config\"): %v\n", strings.HasPrefix(filename, "config"))
	fmt.Printf("  HasSuffix(\".json\"): %v\n", strings.HasSuffix(filename, ".json"))
	fmt.Printf("  HasSuffix(\".xml\"): %v\n", strings.HasSuffix(filename, ".xml"))
	
	// Contains and counting
	text := "The quick brown fox jumps over the lazy dog"
	fmt.Printf("\nText: %q\n", text)
	fmt.Printf("  Contains(\"fox\"): %v\n", strings.Contains(text, "fox"))
	fmt.Printf("  Contains(\"cat\"): %v\n", strings.Contains(text, "cat"))
	fmt.Printf("  Count(\"o\"): %d\n", strings.Count(text, "o"))
	
	// Splitting and joining
	csv := "apple,banana,cherry,date"
	fruits := strings.Split(csv, ",")
	fmt.Printf("\nSplit %q by comma:\n", csv)
	for i, fruit := range fruits {
		fmt.Printf("  [%d] %s\n", i, fruit)
	}
	
	joined := strings.Join(fruits, " | ")
	fmt.Printf("Join with \" | \": %s\n", joined)
	
	// Replacing
	sentence := "I love cats. Cats are great!"
	replaced := strings.ReplaceAll(sentence, "cats", "dogs")
	fmt.Printf("\nOriginal: %s\n", sentence)
	fmt.Printf("ReplaceAll: %s\n", replaced)
	
	// String builder for efficient concatenation
	var builder strings.Builder
	builder.WriteString("Building ")
	builder.WriteString("strings ")
	builder.WriteString("efficiently")
	fmt.Printf("\nstrings.Builder result: %s\n", builder.String())
}

// demonstrateStrconvPackage shows type conversions between strings and other types.
func demonstrateStrconvPackage() {
	fmt.Println("\n=== DEMONSTRATION 9: strconv Package - String Conversions ===")
	
	// String to int
	portStr := "8080"
	port, err := strconv.Atoi(portStr)
	fmt.Printf("Atoi(%q) = %d (error: %v)\n", portStr, port, err)
	
	invalidPort := "abc"
	_, err = strconv.Atoi(invalidPort)
	fmt.Printf("Atoi(%q) = error: %v\n", invalidPort, err)
	
	// Int to string
	num := 42
	numStr := strconv.Itoa(num)
	fmt.Printf("\nItoa(%d) = %q\n", num, numStr)
	
	// Parse with base and bit size
	hexStr := "FF"
	hexNum, _ := strconv.ParseInt(hexStr, 16, 64)
	fmt.Printf("\nParseInt(%q, base 16) = %d\n", hexStr, hexNum)
	
	binStr := "1010"
	binNum, _ := strconv.ParseInt(binStr, 2, 64)
	fmt.Printf("ParseInt(%q, base 2) = %d\n", binStr, binNum)
	
	// Parse bool
	boolTests := []string{"true", "false", "1", "0", "T", "F", "invalid"}
	fmt.Println("\nParseBool:")
	for _, test := range boolTests {
		result, err := strconv.ParseBool(test)
		if err != nil {
			fmt.Printf("  %q -> error: %v\n", test, err)
		} else {
			fmt.Printf("  %q -> %v\n", test, result)
		}
	}
	
	// Parse float
	floatStr := "3.14159"
	floatNum, _ := strconv.ParseFloat(floatStr, 64)
	fmt.Printf("\nParseFloat(%q) = %f\n", floatStr, floatNum)
	
	// Format float with precision
	formatted := strconv.FormatFloat(floatNum, 'f', 2, 64)
	fmt.Printf("FormatFloat(%.5f, 'f', 2) = %q\n", floatNum, formatted)
	
	// Quote and unquote strings
	quoted := strconv.Quote("Hello\nWorld")
	fmt.Printf("\nQuote: %s\n", quoted)
	
	unquoted, _ := strconv.Unquote(quoted)
	fmt.Printf("Unquote: %s\n", unquoted)
}

// demonstrateCompleteWorkflow shows a complete end-to-end workflow of
// reading JSON, validating it, and converting to XML.
func demonstrateCompleteWorkflow() {
	fmt.Println("\n=== DEMONSTRATION 10: Complete Workflow ===")
	
	// Sample input JSON
	inputJSON := `{
  "name": "myapp",
  "version": "1.0.0",
  "debug": true,
  "port": 8080,
  "host": "localhost",
  "log_level": "info"
}`
	
	fmt.Println("Input JSON:")
	fmt.Println(inputJSON)
	
	// Process: JSON → validate → convert to XML
	var outputBuf bytes.Buffer
	
	err := process(
		strings.NewReader(inputJSON),
		&outputBuf,
		"xml",
		true, // validate
	)
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Println("\nOutput XML:")
	fmt.Println(outputBuf.String())
}

// =============================================================================
// MAIN FUNCTION
// =============================================================================

func main() {
	// Define command-line flags
	help := flag.Bool("help", false, "Show help message")
	inputFile := flag.String("f", "", "Input file path (JSON or XML)")
	outputFile := flag.String("o", "", "Output file path")
	targetFormat := flag.String("format", "auto", "Target format: json, xml, or auto (converts to opposite)")
	validate := flag.Bool("validate", false, "Validate configuration structure")
	banner := flag.String("banner", "", "Display a custom banner message")
	demo := flag.Bool("demo", false, "Run demonstration mode showcasing all standard library features")

	flag.Parse()

	// Handle banner from environment or flag
	if *banner == "" {
		if envBanner := os.Getenv("BANNER"); envBanner != "" {
			*banner = envBanner
		}
	}

	if *banner != "" {
		fmt.Println("=" + strings.Repeat("=", len(*banner)+2) + "=")
		fmt.Printf("| %s |\n", *banner)
		fmt.Println("=" + strings.Repeat("=", len(*banner)+2) + "=")
	}

	// Handle help flag
	if *help {
		fmt.Println("myapp - JSON/XML Validation and Conversion Tool")
		fmt.Println("\nUsage:")
		fmt.Println("  myapp [options]")
		fmt.Println("\nOptions:")
		fmt.Println("  -help          Show this help message")
		fmt.Println("  -f <file>      Input file (JSON or XML). If not specified, reads from stdin")
		fmt.Println("  -o <file>      Output file. If not specified, writes to stdout")
		fmt.Println("  -format <fmt>  Target format: 'json', 'xml', or 'auto' (default: auto)")
		fmt.Println("  -validate      Validate configuration structure and field values")
		fmt.Println("  -banner <msg>  Display a custom banner (can also use BANNER env var)")
		fmt.Println("  -demo          Run demonstration mode showcasing standard library features")
		fmt.Println("\nExamples:")
		fmt.Println("  myapp -f config.json -o config.xml -format xml")
		fmt.Println("  cat config.json | myapp -format xml > config.xml")
		fmt.Println("  myapp -f config.xml -validate")
		fmt.Println("  myapp -demo")
		os.Exit(0)
	}

	// Run demonstration mode if requested
	if *demo {
		fmt.Println("╔════════════════════════════════════════════════════════════════╗")
		fmt.Println("║    GO STANDARD LIBRARY - COMPREHENSIVE DEMONSTRATION           ║")
		fmt.Println("╚════════════════════════════════════════════════════════════════╝")
		
		demonstrateFmtPackage()
		demonstrateFlagPackage()
		demonstrateIoPackage()
		demonstrateOsPackage()
		demonstrateEncodingPackage()
		demonstrateBufioPackage()
		demonstrateRegexpPackage()
		demonstrateStringsPackage()
		demonstrateStrconvPackage()
		demonstrateCompleteWorkflow()
		
		fmt.Println("\n" + strings.Repeat("=", 64))
		fmt.Println("Standard library demonstrations complete!")
		fmt.Println("\nFor practical usage, run without -demo flag:")
		fmt.Println("  myapp -f input.json -o output.xml -validate")
		return
	}

	// Normal operation mode
	var input io.Reader = os.Stdin
	var output io.Writer = os.Stdout
	var tempFile *os.File
	var tempFilePath string

	// Handle input file
	if *inputFile != "" {
		file, err := openValidatedInput(*inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	}

	// Handle output file with atomic writes
	if *outputFile != "" {
		var err error
		tempFile, err = createAtomicOutput(*outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		tempFilePath = tempFile.Name()
		defer func() {
			tempFile.Close()
			// Clean up temp file if we didn't finalize it
			if tempFilePath != "" {
				os.Remove(tempFilePath)
			}
		}()
		output = tempFile
	}

	// Process the input
	if err := process(input, output, *targetFormat, *validate); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Finalize atomic output if we used a temp file
	if tempFile != nil {
		tempFile.Close()
		if err := finalizeAtomicOutput(tempFilePath, *outputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error finalizing output: %v\n", err)
			os.Exit(1)
		}
		tempFilePath = "" // Mark as finalized
		fmt.Fprintf(os.Stderr, "Successfully wrote to: %s\n", *outputFile)
	}
}
