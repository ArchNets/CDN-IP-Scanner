package speedtest

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// ConnectivityTest conducts a quick GET request to testUrl to check simple connectivity.
// It returns true if the status code is 2xx or 3xx.
func ConnectivityTest(testUrl string, proxies map[string]string, timeout time.Duration) bool {
	var proxy *url.URL
	for _, v := range proxies {
		proxy, _ = url.Parse(v)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyURL(proxy),
			DisableKeepAlives: true, // Don't wait for connection pooling, close fast
		},
	}

	// v2rayN tries the ping up to 2 times with a short delay if it fails
	// and inside it measures time. We just want a successful connectivity check.
	for i := 0; i < 2; i++ {
		// Create a context with timeout for each request attempt
		ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
		req, err := http.NewRequestWithContext(ctx, "GET", testUrl, nil)
		if err != nil {
			cancel()
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			cancel()
			if i < 1 { // Only sleep if we're going to try again
				time.Sleep(100 * time.Millisecond) // v2rayN uses 100ms between attempts, not 500ms
			}
			continue
		}

		// Close body after reading
		_ = resp.Body.Close()
		cancel()

		// Accept any 2xx or 3xx as a successful connection test
		if resp.StatusCode >= 200 && resp.StatusCode < 400 {
			return true
		}

		time.Sleep(100 * time.Millisecond)
	}

	return false
}
