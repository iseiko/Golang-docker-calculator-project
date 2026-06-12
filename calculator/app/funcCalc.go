package main

import "math"

// Searchable defines the interface that all mathematical operations must implement
// It ensures that each operation type has a ShowAnswers method that takes two floats and returns a result
type Searchable interface {
	ShowAnswers(n1, n2 float64) float64
}

// Sum represents addition operation
// It embeds the Operation struct to inherit its properties
type Sum struct {
	Operation
}

// ShowAnswers performs addition on two numbers and returns the result
// n1: first number
// n2: second number
// Returns: n1 + n2
func (s Sum) ShowAnswers(n1, n2 float64) float64 {
	return n1 + n2
}

// Sub represents subtraction operation
// It embeds the Operation struct to inherit its properties
type Sub struct {
	Operation
}

// ShowAnswers performs subtraction on two numbers and returns the result
// n1: first number
// n2: second number
// Returns: n1 - n2
func (s Sub) ShowAnswers(n1, n2 float64) float64 {
	return n1 - n2
}

// Div represents division operation
// It embeds the Operation struct to inherit its properties
type Div struct {
	Operation
}

// ShowAnswers performs division on two numbers and returns the result
// IMPORTANT: Handles division by zero by returning 0 to prevent panics
// n1: dividend (first number)
// n2: divisor (second number)
// Returns: n1 / n2 if n2 != 0, otherwise returns 0
func (s Div) ShowAnswers(n1, n2 float64) float64 {
	if n2 != 0 {
		return n1 / n2
	}
	// Return 0 instead of dividing by zero (could also return NaN or error)
	return 0
}

// Mul represents multiplication operation
// It embeds the Operation struct to inherit its properties
type Mul struct {
	Operation
}

// ShowAnswers performs multiplication on two numbers and returns the result
// n1: first number
// n2: second number
// Returns: n1 * n2
func (s Mul) ShowAnswers(n1, n2 float64) float64 {
	return n1 * n2
}

// Pow represents exponentiation (power) operation
// It embeds the Operation struct to inherit its properties
type Pow struct {
	Operation
}

// ShowAnswers performs exponentiation on two numbers and returns the result
// Uses the math.Pow function from Go's standard library
// n1: base (first number)
// n2: exponent (second number)
// Returns: n1 ^ n2 (n1 raised to the power of n2)
func (s Pow) ShowAnswers(n1, n2 float64) float64 {
	return math.Pow(n1, n2)
}