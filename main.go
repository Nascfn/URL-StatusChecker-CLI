// Declare the main package for executable programs.
package main

// Import the standard formatting package.
import (
	// Import fmt for formatted I/O.
	"fmt"
	// Import net/http for HTTP client requests.
	"net/http"
	// Import net/url for URL validation and parsing.
	"net/url"
	// Import os for CLI arguments and exit codes.
	"os"
	// Import strings for input trimming and prefix checks.
	"strings"
	// Import time for request timeout duration.
	"time"
)

// Set a fixed timeout for outbound HTTP requests.
const requestTimeout = 10 * time.Second

// main is the program entry point.
func main() {
	// Ensure exactly one CLI argument (the URL) is provided.
	if len(os.Args) != 2 {
		// Print usage instructions when arguments are invalid.
		printUsage()
		// Exit with code 1 for invalid usage.
		os.Exit(1)
	}

	// Remove surrounding whitespace from the input URL argument.
	input := strings.TrimSpace(os.Args[1])
	// Reject empty input after trimming.
	if input == "" {
		// Print an explicit empty-URL error message.
		fmt.Fprintln(os.Stderr, "Error: URL cannot be empty")
		// Print usage instructions after the error.
		printUsage()
		// Exit with code 1 for invalid input.
		os.Exit(1)
	}

	// Normalize and validate the provided URL.
	normalizedURL, err := normalizeURL(input)
	// Handle invalid URL parsing/validation errors.
	if err != nil {
		// Print the URL validation error to stderr.
		fmt.Fprintf(os.Stderr, "Invalid URL: %v\n", err)
		// Exit with code 1 for invalid input.
		os.Exit(1)
	}

	// Perform the HTTP status check against the normalized URL.
	statusCode, err := checkURLStatus(normalizedURL)
	// Handle network/DNS/timeout errors.
	if err != nil {
		// Print DOWN status with the connection error detail.
		fmt.Printf("DOWN: %s (%v)\n", normalizedURL, err)
		// Exit with code 2 for unreachable/down URL.
		os.Exit(2)
	}

	// Treat 2xx and 3xx responses as LIVE.
	if statusCode >= 200 && statusCode < 400 {
		// Print LIVE status with the returned HTTP code.
		fmt.Printf("LIVE: %s (HTTP %d)\n", normalizedURL, statusCode)
		// Exit with code 0 for success.
		os.Exit(0)
	}

	// Print DOWN status for non-2xx/3xx HTTP responses.
	fmt.Printf("DOWN: %s (HTTP %d)\n", normalizedURL, statusCode)
	// Exit with code 2 for down/unhealthy URL status.
	os.Exit(2)
}

// printUsage shows expected command usage examples.
func printUsage() {
	// Print the primary usage format.
	fmt.Fprintln(os.Stderr, "Usage: go run . <url>")
	// Print a concrete usage example.
	fmt.Fprintln(os.Stderr, "Example: go run . https://example.com")
}

// normalizeURL adds a default scheme and validates structure.
func normalizeURL(raw string) (string, error) {
	// Add https scheme if no http/https prefix exists.
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		// Prefix with https by default.
		raw = "https://" + raw
	}

	// Parse and validate the URL as a request URI.
	parsed, err := url.ParseRequestURI(raw)
	// Return parsing error if URL is invalid.
	if err != nil {
		// Propagate parse error to caller.
		return "", err
	}

	// Ensure the parsed URL includes a host.
	if parsed.Host == "" {
		// Return a custom host-missing validation error.
		return "", fmt.Errorf("missing host")
	}

	// Return normalized URL string on success.
	return parsed.String(), nil
}

// checkURLStatus performs an HTTP GET and returns the status code.
func checkURLStatus(target string) (int, error) {
	// Create an HTTP client with a timeout.
	client := &http.Client{Timeout: requestTimeout}

	// Send GET request to the target URL.
	resp, err := client.Get(target)
	// Return error when request fails.
	if err != nil {
		// Return zero status and the underlying request error.
		return 0, err
	}
	// Ensure the response body is closed.
	defer resp.Body.Close()

	// Return the HTTP status code.
	return resp.StatusCode, nil
}