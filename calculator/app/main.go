package main

// main is the entry point of the calculator server application
// It initializes logging and starts the HTTP server
// The main function is kept minimal to separate concerns - it only calls the startup function
func main() {
	// Call the openServer function to initialize and start the HTTP server
	// This function starts listening for incoming HTTP requests on localhost:8080
	openServer()
}
