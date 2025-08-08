package fetchData

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Custom transport that always generates error
type errorTransport struct{}

func (t *errorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("network error simulated")
}

// Test with valid data
func TestFetchByYears_ValidData(t *testing.T) {
	// Test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id": 1, "name": "Test Meeting"}]`))
	}))
	defer server.Close()

	done := make(chan struct{})
	go func() {
		err := FetchByYears(done)
		if err != nil {
			t.Errorf("FetchByYears() error = %v, wantErr %v", err, false)
		}
	}()

	// Wait for function to finish after 5 seconds
	select {
	case <-done:
		// Test réussi
	case <-time.After(5 * time.Second):
		t.Error("FetchByYears() timeout - Function did not finish in time")
	}
}

// Test with network error - Using custom transport
func TestFetchByYears_NetworkError(t *testing.T) {
	// Save default HTTP transport
	originalTransport := http.DefaultTransport
	defer func() {
		http.DefaultTransport = originalTransport
	}()

	// Replace with transport that always generates error
	http.DefaultTransport = &errorTransport{}

	done := make(chan struct{}, 1)
	err := FetchByYears(done)

	// Restore transport immediately to avoid affecting other tests
	http.DefaultTransport = originalTransport

	if err == nil {
		t.Error("FetchByYears() should return error when network problem occurs")
	}

	if err != nil && err.Error() != "cannot Access API" {
		t.Errorf("FetchByYears() error = %v, want 'cannot Access API'", err)
	}
}

// Test with consecutive empty responses
func TestFetchByYears_ConsecutiveEmptyResponses(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`)) // Empty response
	}))
	defer server.Close()

	done := make(chan struct{})
	go func() {
		err := FetchByYears(done)
		if err != nil {
			t.Errorf("FetchByYears() error = %v, wantErr %v", err, false)
		}
	}()

	// Wait for function to finish
	select {
	case <-done:
		// Function should stop after 3 consecutive empty responses
	case <-time.After(10 * time.Second):
		t.Error("FetchByYears() timeout - function should stop after 3 empty responses")
	}
}

// Test with invalid JSON
func TestFetchByYears_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid json`)) // Malformed JSON
	}))
	defer server.Close()

	done := make(chan struct{})
	go func() {
		err := FetchByYears(done)
		if err != nil {
			t.Errorf("FetchByYears() error = %v, wantErr %v", err, false)
		}
	}()

	// Function should handle invalid JSON and continue
	select {
	case <-done:
		// Test successful
	case <-time.After(10 * time.Second):
		t.Error("FetchByYears() timeout")
	}
}

// Test the done channel
func TestFetchByYears_ChannelSignaling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	done := make(chan struct{})

	go func() {
		_ = FetchByYears(done)
	}()

	// Verify the signal is corrtcly sent to the channel
	select {
	case <-done:
		// Test successful
	case <-time.After(15 * time.Second):
		t.Error("Le canal 'done' n'a pas reçu de signal dans les temps")
	}
}

// Test with mixed data
func TestFetchByYears_MixedData(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Alternante between empty and complete response
		if requestCount%2 == 0 {
			w.Write([]byte(`[]`)) // With empty response
		} else {
			w.Write([]byte(`[{"id": 1, "name": "Test"}]`)) // With data response
		}
	}))
	defer server.Close()

	done := make(chan struct{})
	go func() {
		err := FetchByYears(done)
		if err != nil {
			t.Errorf("FetchByYears() error = %v, wantErr %v", err, false)
		}
	}()

	// Waiting the function to finish
	select {
	case <-done:
		// Test successful
	case <-time.After(20 * time.Second):
		t.Error("FetchByYears() timeout")
	}
}

// Benchmark to measure performance
func BenchmarkFetchByYears(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	for i := 0; i < b.N; i++ {
		done := make(chan struct{})
		go func() {
			FetchByYears(done)
		}()
		<-done
	}
}
