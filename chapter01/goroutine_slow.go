package main

import (
	"fmt"
	"time"
)

// Doing Something Slow
// This example demonstrates running a slow function concurrently
// while the main function continues with other work.

func DoSomethingSlow() {
	fmt.Println("SLOW: maybe I'm doing network stuff?")
	time.Sleep(1 * time.Second)
	fmt.Println("SLOW: Okay, finally finished!")
}

func main() {
	fmt.Println("Starting the main task...")

	// Create a goroutine to do something slow.
	go DoSomethingSlow()

	fmt.Println("Resuming the main task...")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Finished the main task!")
	time.Sleep(1 * time.Second)
}
