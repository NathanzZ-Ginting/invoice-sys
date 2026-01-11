package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// ProxyRequest forwards a request to a service
func ProxyRequest(w http.ResponseWriter, r *http.Request, serviceURL string) {
	// Build target URL
	targetURL := serviceURL + r.URL.Path
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	log.Printf("🔀 Proxying %s %s -> %s", r.Method, r.URL.Path, targetURL)

	// Create new request
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		log.Printf("❌ Error creating proxy request: %v", err)
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Send request to service
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Printf("❌ Error forwarding request: %v", err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	io.Copy(w, resp.Body)
}

// GetServiceURL returns the URL for a service
func GetServiceURL(serviceName string) string {
	host := os.Getenv(serviceName + "_HOST")
	port := os.Getenv(serviceName + "_PORT")
	
	if host == "" {
		host = "localhost"
	}
	
	return fmt.Sprintf("http://%s:%s", host, port)
}
