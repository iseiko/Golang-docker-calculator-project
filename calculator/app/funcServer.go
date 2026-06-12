package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Method is a map that routes HTTP methods (GET, POST, etc.) to their corresponding handler maps
var Method = map[string]map[string]func(w http.ResponseWriter, r *http.Request){
	"GET": MethodFunction,
}

// MethodFunction is a map that routes URL paths to their handler functions
var MethodFunction = map[string]func(w http.ResponseWriter, r *http.Request){
	"/result": MethodResult,
}

// MethodResult handles calculator operations from GET requests
// It expects a query parameter "op" with format: "number1 operator number2"
// Supports formats: "10+5", "10 + 5", or "10%20+%205"
func MethodResult(res http.ResponseWriter, req *http.Request) {
	// Set response header to indicate JSON content type
	res.Header().Set("Content-Type", "application/json")

	// Extract the "op" query parameter from the URL query string
	operation := req.URL.Query().Get("op")

	// Validate that operation is not empty
	if operation == "" {
		log.Printf("ERROR: Empty operation parameter received")
		PrintErro(res, "empty operation", http.StatusBadRequest)
		return
	}

	log.Printf("DEBUG: Raw operation received: %q", operation)

	// Parse the operation - try multiple formats
	n1, operator, n2, err := parseOperation(operation)

	if err != nil {
		log.Printf("ERROR: Failed to parse operation: %v", err)
		PrintErro(res, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("DEBUG: Parsed - n1=%f, operator=%s, n2=%f", n1, operator, n2)

	// Look up the operator in the OperationsTable
	searchable, ok := OperationsTable[operator]

	// If the operator is not found in the OperationsTable, it's not supported
	if !ok {
		log.Printf("ERROR: Unsupported operator: %s", operator)
		PrintErro(res, "unsupported operator: "+operator, http.StatusBadRequest)
		return
	}

	// Log the successful operation being performed
	log.Printf("INFO: Performing operation - %f %s %f", n1, operator, n2)

	// Create an Operation struct and perform the calculation
	r := Operation{
		Op:     operator,
		Result: searchable.ShowAnswers(n1, n2),
	}

	// Log the result
	log.Printf("INFO: Operation result: %f", r.Result)

	// Convert the result to JSON and write it to the response
	jsonData, err := json.Marshal(r)
	if err != nil {
		log.Printf("ERROR: Failed to marshal result to JSON: %v", err)
		PrintErro(res, "server error", http.StatusInternalServerError)
		return
	}

	// Set HTTP status code to 200 (OK) - operation was successful
	res.WriteHeader(http.StatusOK)

	// Write the JSON response to the client
	res.Write(jsonData)
}

// parseOperation extracts numbers and operator from the operation string
// Supports multiple formats: "10+5", "10 + 5", "10 + 5" (URL encoded)
//
// Parameters:
//   - operation: string containing the mathematical operation
//
// Returns:
//   - n1: first number as float64
//   - operator: the operator symbol (+, -, *, /, ^)
//   - n2: second number as float64
//   - error: error if parsing fails
func parseOperation(operation string) (float64, string, float64, error) {
	// Trim whitespace from both ends
	operation = strings.TrimSpace(operation)

	// Try format with spaces first: "10 + 5"
	if strings.Contains(operation, " ") {
		return parseWithSpaces(operation)
	}

	// Try format without spaces: "10+5"
	return parseWithoutSpaces(operation)
}

// parseWithSpaces handles operations like "10 + 5" or "10 + 5"
func parseWithSpaces(operation string) (float64, string, float64, error) {
	// Split by spaces
	parts := strings.Split(operation, " ")

	log.Printf("DEBUG: Split by spaces - parts: %v", parts)

	// Filter out empty strings (in case of multiple spaces)
	var filtered []string
	for _, part := range parts {
		if part != "" {
			filtered = append(filtered, part)
		}
	}

	// We expect exactly 3 parts: number, operator, number
	if len(filtered) != 3 {
		return 0, "", 0, &ParseError{
			message: "invalid format: expected 'number operator number' (e.g., '10 + 5')",
			input:   operation,
		}
	}

	// Parse first number
	n1, err := strconv.ParseFloat(filtered[0], 64)
	if err != nil {
		return 0, "", 0, &ParseError{
			message: "invalid first number: " + filtered[0],
			input:   operation,
		}
	}

	// Extract operator
	operator := filtered[1]

	// Validate operator
	if !isValidOperator(operator) {
		return 0, "", 0, &ParseError{
			message: "invalid operator: " + operator,
			input:   operation,
		}
	}

	// Parse second number
	n2, err := strconv.ParseFloat(filtered[2], 64)
	if err != nil {
		return 0, "", 0, &ParseError{
			message: "invalid second number: " + filtered[2],
			input:   operation,
		}
	}

	return n1, operator, n2, nil
}

// parseWithoutSpaces handles operations like "10+5" without spaces
func parseWithoutSpaces(operation string) (float64, string, float64, error) {
	// Use regex to find the operator
	// Pattern: number (optional decimal) operator number (optional decimal)
	// Operators: +, -, *, /, ^
	re := regexp.MustCompile(`^(-?\d+\.?\d*)([\+\-\*/\^])(-?\d+\.?\d*)$`)

	matches := re.FindStringSubmatch(operation)

	log.Printf("DEBUG: Regex matches: %v", matches)

	// Regex should return 4 matches: full string, n1, operator, n2
	if len(matches) != 4 {
		return 0, "", 0, &ParseError{
			message: "invalid format: expected 'number operator number' (e.g., '10+5', '10 + 5')",
			input:   operation,
		}
	}

	// Parse first number
	n1, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, "", 0, &ParseError{
			message: "invalid first number: " + matches[1],
			input:   operation,
		}
	}

	// Extract operator
	operator := matches[2]

	// Parse second number
	n2, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return 0, "", 0, &ParseError{
			message: "invalid second number: " + matches[3],
			input:   operation,
		}
	}

	return n1, operator, n2, nil
}

// isValidOperator checks if the operator is one we support
func isValidOperator(op string) bool {
	_, exists := OperationsTable[op]
	return exists
}

// ParseError is a custom error type for parsing failures
type ParseError struct {
	message string
	input   string
}

// Error implements the error interface
func (e *ParseError) Error() string {
	return e.message
}

// PrintErro sends an error response in JSON format back to the client
func PrintErro(res http.ResponseWriter, op string, statusCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	errorResponse := map[string]string{
		"result": "invalid expression",
		"op":     op,
	}

	jsonData, err := json.Marshal(errorResponse)
	if err != nil {
		log.Printf("ERROR: Failed to marshal error to JSON: %v", err)
		return
	}

	res.Write(jsonData)
}

// ChangeJson converts an Operation struct to a JSON byte slice
func ChangeJson(s Operation) []byte {
	jsonData, err := json.Marshal(s)
	if err != nil {
		log.Printf("ERROR: Failed to convert operation to JSON: %v", err)
		return []byte(`{"error":"marshaling failed"}`)
	}
	return jsonData
}