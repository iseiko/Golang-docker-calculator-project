package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// openServer initializes and starts the HTTP server for the calculator
// It creates an HTTP server configured with:
//   - Address: localhost on port 8080 (127.0.0.1:8080)
//   - Handler: FunctionHandler struct that routes requests
//   - Read timeout: 10 seconds (max time to read an incoming request)
//   - Write timeout: 10 seconds (max time to send a response)
//
// The server uses a custom handler instead of the default http.ServeMux
// This approach gives more control over request routing and validation
func openServer() {
	// Create a new HTTP server with custom configuration
	servidor := &http.Server{
		// Addr specifies the address and port the server listens on
		// 127.0.0.1 is localhost (only accessible from the machine currently running the code)
		// 8080 is the port number
		Addr: "0.0.0.0:8080",

		// Handler is the custom HTTP handler (defined below as FunctionHandler)
		// It implements the http.Handler interface with ServeHTTP method
		Handler: FunctionHandler{},

		// ReadTimeout sets the maximum duration for reading an entire request
		// If a request takes longer than 10 seconds to arrive, it will be cancelled
		// This prevents slow client attacks from tying up server resources
		ReadTimeout: 10 * time.Second,

		// WriteTimeout sets the maximum duration for writing a response
		// If the server takes longer than 10 seconds to send a response, the connection is closed
		// This prevents slow client connections from blocking the server
		WriteTimeout: 10 * time.Second,
	}

	// Start the server and listen for incoming connections
	// ListenAndServe is a blocking call - it runs until the server stops or an error occurs
	// log.Fatal prints the error and terminates the program if ListenAndServe returns an error
	// Common errors: port already in use, permission denied, etc.
	log.Fatal(servidor.ListenAndServe())
}

// FunctionHandler is a custom HTTP handler that implements the http.Handler interface
// It acts like a "class" in object-oriented programming (Go is not OO, but uses structs similarly)
// The FunctionHandler doesn't need any fields - it just needs the ServeHTTP method
type FunctionHandler struct{}

// ServeHTTP implements the http.Handler interface
// It's called for every HTTP request the server receives
// This method:
// 1. Validates that the request path and method are recognized
// 2. Routes valid requests to the appropriate handler function
// 3. Returns a 404 error for unrecognized requests
//
// Parameters:
//   - res: http.ResponseWriter - used to write the response back to the client
//   - req: *http.Request - contains information about the incoming request
func (f FunctionHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Validate the incoming request
	// Check if the requested HTTP method and path exist in our routing tables
	if !ValidServer(req) {
		// If validation fails, send a 404 Not Found error in JSON format
		// TableMessage[404] contains the error details: code 404 and "Not found" message
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusNotFound)
		res.Write(MessageToJson(TableMessage[404]))
		return
	}

	// If validation passes, route the request to the appropriate handler
	// Method[req.Method] gets the handler map for the HTTP method (e.g., "GET")
	// [req.URL.Path] gets the specific handler function for that path (e.g., "/result")
	// (res, req) calls the handler function with the response writer and request
	Method[req.Method][req.URL.Path](res, req)
}

// ValidServer checks if an incoming request can be handled by the server
// It validates that both the HTTP method and the requested path are supported
//
// Parameters:
//   - req: *http.Request - the incoming HTTP request
//
// Returns: true if the request can be handled, false if it's an unsupported method/path
func ValidServer(req *http.Request) bool {
	// Check if the request's HTTP method exists in the Method map
	// If Method[req.Method] returns nil, the method is not supported
	if Method[req.Method] == nil {
		return false
	}

	// Check if the requested path exists in the method's handler map
	// If Method[req.Method][req.URL.Path] is not nil, the path is supported
	// This returns true if both method and path are found, false otherwise
	return Method[req.Method][req.URL.Path] != nil
}

// HTTPMessage represents a standardized HTTP error response
// All error responses follow this format when sent as JSON
type HTTPMessage struct {
	// Code is the HTTP status code (e.g., 404, 500)
	Code int `json:"code"`

	// Error is a human-readable error message
	Error string `json:"error"`
}

// TableMessage is a map that stores standardized error responses for different HTTP status codes
// Each status code maps to an HTTPMessage with the appropriate code and description
var TableMessage = map[int]HTTPMessage{
	// 404 Not Found - requested resource or endpoint doesn't exist
	404: {404, "Not found"},

	// 500 Internal Server Error - something went wrong on the server side
	500: {500, "Server error"},
}

// MessageToJson converts an HTTPMessage struct to JSON format
// It's used to serialize error messages before sending them to clients
//
// Parameters:
//   - m: HTTPMessage - the error message to convert
//
// Returns: JSON representation as a byte slice
func MessageToJson(m HTTPMessage) []byte {
	// Convert (marshal) the HTTPMessage struct to JSON
	// The struct tags like `json:"code"` tell the marshaler which JSON field names to use
	jsonData, err := json.Marshal(m)
	if err != nil {
		log.Printf("ERROR: Failed to marshal HTTP message to JSON: %v", err)
		return []byte(`{"code":500,"error":"marshaling failed"}`)
	}
	return jsonData
}
