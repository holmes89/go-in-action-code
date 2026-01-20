package main

import (
	"fmt"
	"time"
)

// Running A Function As A Goroutine
// This example shows the basic syntax for creating a goroutine

func makeMeConcurrent() {
	// Do some concurrent work...
	fmt.Println("Running concurrently!")
}

func main() {
	go makeMeConcurrent()
	// Continue with the main goroutine...
	fmt.Println("Main goroutine continues...")
	
	// Give the goroutine time to complete before program exits
	time.Sleep(100 * time.Millisecond)
}
