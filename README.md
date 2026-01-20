# Go In Action - Code Examples

This directory contains all runnable code examples from the book "Go In Action (Second Edition)".

## Chapter Organization

- **[chapter01](chapter01/)** - Introducing Go: Basic goroutines, channels, interfaces, generics
- **[chapter02](chapter02/)** - Diving Into Go: Word count program, scope, and fundamentals
- **[chapter03](chapter03/)** - Primitive Types: Integers, floats, strings, runes, pointers, structs
- **[chapter04](chapter04/)** - Collection Types: Arrays, slices, and maps
- **[chapter05](chapter05/)** - Working With Types: Named types, methods, interfaces, composition
- **[chapter06](chapter06/)** - Generics: Type parameters and generic data structures
- **[chapter07](chapter07/)** - Errors: Error handling, custom errors, panic/recover
- **[chapter08](chapter08/)** - Testing: Unit tests, table-driven tests, benchmarks
- **[chapter09](chapter09/)** - Concurrency: Goroutines, channels, WaitGroups, mutexes
- **[chapter10](chapter10/)** - Standard Library: JSON, XML, flags, file I/O
- **[chapter11](chapter11/)** - Larger Projects: Modules, packages, workspaces

## Running Examples

Each chapter directory contains:
- Runnable `.go` files
- A `README.md` with specific instructions
- Any supporting files needed

Navigate to any chapter and run:
```bash
cd chapter01
go run goroutine_basic.go
```

## Requirements

- Go 1.18 or later (for generics support)
- Some examples may require additional packages (noted in chapter READMEs)

## Testing

To run all tests in a chapter:
```bash
cd chapter08
go test -v
```

## Building

To build any program:
```bash
go build -o program_name file.go
./program_name
```

Enjoy learning Go! 🎉
