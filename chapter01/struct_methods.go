package main

import "fmt"

// A Go Struct And Method
// This example demonstrates how to add methods to a struct type.

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func (p Person) Name() string {
	return p.FirstName + " " + p.LastName
}

func (p Person) Describe() string {
	return fmt.Sprintf("%s is %d years old", p.Name(), p.Age)
}

func main() {
	andy := Person{
		FirstName: "Andy",
		LastName:  "Walker",
		Age:       43,
	}
	fmt.Println(andy.Describe())
}
