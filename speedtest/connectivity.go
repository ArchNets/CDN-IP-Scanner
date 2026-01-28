package speedtest

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// ConnectivityTest checks if the connection works by testing http://www.google.com/generate_204
// Returns true if connection is successful (status 204 or 200), false otherwise
func ConnectivityTest(proxies map[string]string, timeout time.Duration) (bool, error) {
	// set proxy to request
	var proxy *url.URL
	// set proxies to req header
	for _, v := range proxies {
		proxy, _ = url.Parse(v)
	}

	// Create request to Google's generate_204 endpoint (lightweight connectivity check)
	req, err := http.NewRequest("GET", "http://www.google.com/generate_204", nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}

	// Set up client with shorter timeout for connectivity check
	client := &http.Client{
		Timeout:   timeout * time.Second,
		Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("connection failed: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// Ignore close errors
		}
	}(resp.Body)

	// Check if we got a successful response (204 No Content or 200 OK)
	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
