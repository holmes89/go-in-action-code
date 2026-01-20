package main

import "fmt"

// A Go Struct
// This example shows the basic structure of a Go type.

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func main() {
	andy := Person{
		FirstName: "Andy",
		LastName:  "Walker",
		Age:       43,
	}
	fmt.Printf("%+v\n", andy)
}
