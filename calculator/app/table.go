package main

// Operation represents a mathematical operation and its result
// This struct is used to store and return the details of a calculated operation
type Operation struct {
	// Result is the outcome of the mathematical operation (e.g., 10 + 5 = 15)
	// The `json:"result"` tag tells the JSON marshaler to use "result" as the field name in JSON
	Result float64 `json:"result"`

	// Op is the operation that was performed ("+", "-", "*", "/", "^")
	// The `json:"op"` tag tells the JSON marshaler to use "op" as the field name in JSON
	Op string `json:"op"`
}

// OperationsTable is a lookup table that maps operator symbols to their corresponding operation implementations
// It implements the Strategy pattern - each operator is mapped to an object that knows how to perform it
//
// Example usage:
//   operation := OperationsTable["+"]           // Get the Sum operation
//   result := operation.ShowAnswers(10, 5)      // Call ShowAnswers to get 15
//
// Supported operations:
//   "+" -> Addition (Sum)
//   "-" -> Subtraction (Sub)
//   "*" -> Multiplication (Mul)
//   "/" -> Division (Div)
//   "^" -> Exponentiation (Pow)
var OperationsTable = map[string]Searchable{
	// Map "+" operator to Sum struct which implements addition
	"+": Sum{},

	// Map "-" operator to Sub struct which implements subtraction
	"-": Sub{},

	// Map "*" operator to Mul struct which implements multiplication
	"*": Mul{},

	// Map "/" operator to Div struct which implements division
	"/": Div{},

	// Map "^" operator to Pow struct which implements exponentiation
	"^": Pow{},
}