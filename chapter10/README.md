# Chapter 10: The Standard Library

A comprehensive reference implementation demonstrating Go's standard library through a complete JSON/XML validation and conversion CLI tool.

## Overview

This chapter covers the essential packages in Go's standard library and shows how they work together to build production-quality applications. The tool demonstrates real-world usage patterns for each package.

## Packages Covered

### 1. `fmt` - Formatted I/O
- `Print`, `Println`, `Printf` - Console output
- `Sprint`, `Sprintf` - String formatting
- Format verbs: `%v`, `%+v`, `%#v`, `%T`, `%d`, `%s`, `%q`, etc.
- `Fprintf` - Writing to any `io.Writer`
- `Errorf` - Creating formatted errors

### 2. `flag` - Command-line Arguments
- Defining flags: `flag.String()`, `flag.Bool()`, `flag.Int()`
- Parsing arguments: `flag.Parse()`
- Accessing positional arguments: `flag.Args()`
- Building help messages

### 3. `io` - Core I/O Interfaces
- `io.Reader` and `io.Writer` interfaces
- `io.ReadCloser`, `io.WriteCloser` composition
- Utility functions: `io.Copy()`, `io.ReadAll()`, `io.LimitReader()`
- Interface-based design for flexibility

### 4. `os` - Operating System Interface
- Standard streams: `os.Stdin`, `os.Stdout`, `os.Stderr`
- File operations: `os.Open()`, `os.Create()`, `os.Stat()`
- Environment variables: `os.Getenv()`, `os.Setenv()`, `os.LookupEnv()`
- Exit codes: `os.Exit()`
- Temporary files: `os.CreateTemp()`
- Atomic file operations with rename

### 5. `encoding/json` and `encoding/xml` - Structured Data
- Marshaling: converting Go structs to JSON/XML
- Unmarshaling: parsing JSON/XML into Go structs
- Struct tags for field mapping
- `json:"name"`, `xml:"name,attr"`, `omitempty`, etc.
- Custom marshaling/unmarshaling methods
- `MarshalIndent()` for formatted output

### 6. `bufio` - Buffered I/O
- `bufio.Reader` for efficient reading
- `bufio.Writer` for batched writes
- `Peek()` - looking ahead without consuming
- `ReadString()`, `ReadLine()` - line-by-line reading
- `Flush()` - ensuring buffered data is written
- Performance benefits of buffering

### 7. `regexp` - Regular Expressions
- Pattern compilation: `regexp.MustCompile()`
- Matching: `MatchString()`
- Capturing groups: `FindStringSubmatch()`
- Common patterns for validation
- Performance considerations

### 8. `strings` - String Manipulation
- Case conversion: `ToLower()`, `ToUpper()`
- Trimming: `TrimSpace()`, `Trim()`
- Searching: `Contains()`, `HasPrefix()`, `HasSuffix()`
- Splitting and joining: `Split()`, `Join()`
- Replacing: `ReplaceAll()`
- `strings.Builder` for efficient concatenation

### 9. `strconv` - String Conversions
- `Atoi()`, `Itoa()` - string to/from int
- `ParseInt()`, `ParseFloat()`, `ParseBool()` - parsing with error handling
- `FormatInt()`, `FormatFloat()` - formatting with control
- Base conversion (binary, hex, decimal)
- Precision control for floats
- `Quote()` and `Unquote()` for string literals

## Features

The CLI tool demonstrates practical usage of all these packages:

1. **Format Detection**: Automatically detects JSON or XML input
2. **Validation**: Validates configuration structure and field values using regex
3. **Conversion**: Converts between JSON and XML formats
4. **Flexible I/O**: Reads from files or stdin, writes to files or stdout
5. **Atomic Writes**: Uses temporary files to prevent corruption
6. **Error Handling**: Comprehensive error messages and exit codes
7. **Environment Support**: Configuration via environment variables
8. **Buffered Performance**: Uses buffered I/O for efficiency

## Usage

### Run Demonstrations

```bash
# See all standard library features in action
go run main.go -demo
```

This runs 10 comprehensive demonstrations:
1. fmt Package - Formatted I/O
2. flag Package - Command-line Arguments  
3. io Package - Core I/O Interfaces
4. os Package - Operating System Interface
5. encoding Package - JSON and XML
6. bufio Package - Buffered I/O
7. regexp Package - Pattern Matching
8. strings Package - String Manipulation
9. strconv Package - String Conversions
10. Complete Workflow - End-to-end example

### Practical Usage

```bash
# Show help
go run main.go -help

# Convert JSON to XML
go run main.go -f config.json -o config.xml -format xml

# Convert XML to JSON with validation
go run main.go -f config.xml -o config.json -format json -validate

# Auto-detect and convert (JSON→XML or XML→JSON)
go run main.go -f input.json -o output.xml

# Read from stdin, write to stdout
cat config.json | go run main.go -format xml

# Validate without converting
go run main.go -f config.xml -validate

# Custom banner
go run main.go -banner "My Config Tool" -f config.json
```

### Example Configuration Files

**config.json:**
```json
{
  "name": "myapp",
  "version": "1.0.0",
  "debug": false,
  "port": 8080,
  "host": "localhost",
  "log_level": "info"
}
```

**config.xml:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<configuration name="myapp" version="1.0.0">
  <debug>false</debug>
  <port>8080</port>
  <host>localhost</host>
  <log_level>info</log_level>
</configuration>
```

## Validation Rules

The tool validates configurations with:

- **Name**: Must be alphanumeric with hyphens/underscores
- **Version**: Must follow semantic versioning (X.Y.Z)
- **Port**: Must be in range 1-65535
- **Host**: Must be a valid hostname/URL
- **LogLevel**: Must be one of: debug, info, warn, error, fatal

## Key Patterns Demonstrated

### 1. Interface-Based Design
```go
func process(input io.Reader, output io.Writer, format string, validate bool) error
```
Works with any input/output source - files, stdin/stdout, network, memory buffers.

### 2. Struct Tags for Multiple Formats
```go
type Config struct {
    Name string `json:"name" xml:"name,attr"`
}
```
Single struct works with both JSON and XML.

### 3. Atomic File Operations
```go
tempFile, _ := os.CreateTemp(".", ".tmp-")
// Write to tempFile...
os.Rename(tempFile.Name(), finalPath)
```
Prevents corruption from partial writes.

### 4. Buffered I/O for Performance
```go
br := bufio.NewReader(input)
bw := bufio.NewWriter(output)
defer bw.Flush()
```
Reduces system calls dramatically.

### 5. Peek for Format Detection
```go
firstBytes, _ := br.Peek(10)
if strings.HasPrefix(string(firstBytes), "{") {
    // It's JSON
}
```
Detect format without consuming data.

### 6. Compiled Regex for Efficiency
```go
var namePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
// Reuse namePattern many times
```
Compile once, use many times.

### 7. Environment Variable Fallback
```go
if *banner == "" {
    *banner = os.Getenv("BANNER")
}
```
Flags override environment variables.

### 8. Proper Error Handling
```go
if err := json.Unmarshal(data, &cfg); err != nil {
    return fmt.Errorf("invalid JSON: %w", err)
}
```
Wrap errors with context using `%w`.

## Standard Library Philosophy

Go's standard library demonstrates several key principles:

1. **Simplicity**: Small, focused packages with clear responsibilities
2. **Composability**: Interfaces that work together seamlessly
3. **Performance**: Efficient implementations without sacrificing clarity
4. **Completeness**: Everything needed for production applications
5. **Consistency**: Similar patterns across different packages
6. **Documentation**: Comprehensive examples and clear godocs

## Learning Resources

- [Official Go Standard Library](https://pkg.go.dev/std)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Regex101](https://regex101.com/) - Regex testing playground

## Architecture

```
Command Line
     ↓
  Flags (flag)
     ↓
Input Source (os.Stdin or os.File)
     ↓
Format Detection (bufio + strings)
     ↓
Parsing (encoding/json or encoding/xml)
     ↓
Validation (regexp + custom logic)
     ↓
Conversion (encoding/json or encoding/xml)
     ↓
Buffered Output (bufio)
     ↓
Output Destination (os.Stdout or os.File)
     ↓
Atomic Finalization (os.Rename)
```

## Performance Considerations

1. **Buffering**: Always use `bufio` for file I/O
2. **Regex**: Compile patterns once with `regexp.MustCompile()`
3. **String Building**: Use `strings.Builder` for concatenation
4. **Streaming**: For large files, consider streaming instead of `ReadAll()`
5. **Validation**: Validate early to fail fast
6. **Resource Cleanup**: Use `defer` for closing files and flushing buffers

## Testing

The code includes a comprehensive demo mode that tests all functionality:

```bash
go run main.go -demo
```

This runs through all standard library features and shows you the output from each demonstration.

## Summary

This chapter's code demonstrates that Go's standard library provides everything needed to build robust, performant command-line tools. By understanding these core packages and their patterns, you can build production-quality applications without external dependencies.

The key takeaway: **Learn the standard library first**. It's well-designed, well-tested, and covers most common use cases. Only reach for third-party packages when you need specialized functionality not provided by the standard library.
