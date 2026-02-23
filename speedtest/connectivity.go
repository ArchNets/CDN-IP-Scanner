package speedtest

import (
	"crypto/tls"
	"io"
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

	req, err := http.NewRequest("GET", testUrl, nil)
	if err != nil {
		return false
	}

	client := &http.Client{
		Timeout: timeout * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Accept any 2xx or 3xx as a successful connection test
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return true
	}

	return false
}
