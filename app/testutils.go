package app

import "net/http"

// TestMockRoundTripper is a mock implementation of http.RoundTripper for testing purposes only
type TestMockRoundTripper struct {
	Response *http.Response
	Err      error
}

func (m *TestMockRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return m.Response, m.Err
}
