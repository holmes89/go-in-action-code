# Chapter 2: Diving Into Go - Word Count Program

This directory contains a comprehensive reference implementation of the word count program from Chapter 2.

## Main File

**main.go** - Complete progression from basic space counting to file reading, all in one annotated reference file

## What's Inside

The main.go file demonstrates the complete evolution of the word count program:

1. **Iteration 1** (commented): Basic space counting approach
2. **Iteration 2** (commented): Using strings.Fields for better word splitting
3. **Iteration 3** (active): Reading from files with command-line arguments
4. **Key Concepts**: Comprehensive comments explaining all Go fundamentals
5. **Exercise Ideas**: Suggestions for extending the program

## Running the Program

Run with the sample file:
```bash
go run main.go sample.txt
```

Try with your own files:
```bash
go run main.go yourfile.txt
```

## Building an Executable

Build the program:
```bash
go build -o wordcount main.go
```

Run the executable:
```bash
./wordcount sample.txt
```

## Testing Different Scenarios

The program handles:
- Multiple spaces between words
- Empty lines
- Various punctuation
- Different file types

## Key Concepts Covered

- Package declarations and imports
- Variables and type inference
- For loops and iteration
- Functions and multiple return values
- Error handling patterns
- Slices and arrays
- Command-line arguments
- File I/O with os.ReadFile
- String manipulation with strings.Fields

## Go Tools Demonstrated

Format your code:
```bash
gofmt -w main.go
```

View documentation:
```bash
go doc strings.Fields
go doc os.ReadFile
go doc fmt.Println
```

Check for errors:
```bash
go vet main.go
```

## Exercise Ideas

Try extending the program with:
1. A verbose (-v) flag showing each word
2. Line counting in addition to words
3. Byte counting
4. Support for multiple files
5. A help (-h) flag
6. Unique word counting with maps
