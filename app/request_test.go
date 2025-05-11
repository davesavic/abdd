package app_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/davesavic/abdd/app"
	"github.com/stretchr/testify/assert"
)

type mockRoundTripper struct {
	response *http.Response
	err      error
}

func (m *mockRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestMakeRequest(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(*app.Abdd, *app.Test)
		expects func(app.Abdd, *app.Test, error)
	}{
		{
			name: "Successful request",
			setup: func(a *app.Abdd, test *app.Test) {
				mockBody := "response body"
				a.Global = app.Global{
					Config: app.Config{
						BaseURL: "https://example.com",
					},
				}

				resp := &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(mockBody)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}

				a.Client = &http.Client{
					Transport: &mockRoundTripper{response: resp},
				}

				test.Request = &app.TestRequest{
					Method:  "GET",
					URL:     "/api/users",
					Headers: map[string]string{"Accept": "application/json"},
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, a.LastResponse)
				assert.Equal(t, 200, *a.LastResponse.Code)
				assert.Equal(t, "response body", *a.LastResponse.Body)
				assert.Equal(t, "application/json", a.LastResponse.Headers["Content-Type"])
			},
		},
		{
			name: "Request with body",
			setup: func(a *app.Abdd, test *app.Test) {
				mockBody := "response data"

				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					body, _ := io.ReadAll(r.Body)
					assert.Equal(t, "request payload", string(body))
					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(mockBody))
				}))

				a.Global = app.Global{
					Config: app.Config{
						BaseURL: server.URL,
					},
				}
				a.Client = server.Client()

				requestBody := "request payload"
				test.Request = &app.TestRequest{
					Method:  "POST",
					URL:     "/api/data",
					Body:    &requestBody,
					Headers: map[string]string{"Content-Type": "text/plain"},
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, a.LastResponse)
				assert.Equal(t, 200, *a.LastResponse.Code)
				assert.Equal(t, "response data", *a.LastResponse.Body)
			},
		},
		{
			name: "Client returns error",
			setup: func(a *app.Abdd, test *app.Test) {
				a.Global = app.Global{
					Config: app.Config{
						BaseURL: "https://example.com",
					},
				}

				a.Client = &http.Client{
					Transport: &mockRoundTripper{err: http.ErrHandlerTimeout},
				}

				test.Request = &app.TestRequest{
					Method: "GET",
					URL:    "/api/users",
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to make request")
			},
		},
		{
			name: "Invalid request URL",
			setup: func(a *app.Abdd, test *app.Test) {
				a.Global = app.Global{
					Config: app.Config{
						BaseURL: ":",
					},
				}

				test.Request = &app.TestRequest{
					Method: "GET",
					URL:    "/api/users",
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to create request")
			},
		},
		{
			name: "204 No Content response",
			setup: func(a *app.Abdd, test *app.Test) {
				a.Store = make(map[string]any)
				a.Global = app.Global{
					Config: app.Config{
						BaseURL: "https://example.com",
					},
				}

				resp := &http.Response{
					StatusCode: 204,
					Body:       io.NopCloser(strings.NewReader("")),
					Header:     http.Header{},
				}

				a.Client = &http.Client{
					Transport: &mockRoundTripper{response: resp},
				}

				test.Request = &app.TestRequest{
					Method:  "DELETE",
					URL:     "/api/users/123",
					Headers: map[string]string{"Accept": "application/json"},
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, a.LastResponse)
				assert.Equal(t, 204, *a.LastResponse.Code)
				assert.Equal(t, "", *a.LastResponse.Body)
			},
		},
		{
			name: "400 Bad Request response",
			setup: func(a *app.Abdd, test *app.Test) {
				a.Store = make(map[string]any)
				a.Global = app.Global{
					Config: app.Config{
						BaseURL: "https://example.com",
					},
				}

				errorBody := `{"error": "Invalid request parameters"}`
				resp := &http.Response{
					StatusCode: 400,
					Body:       io.NopCloser(strings.NewReader(errorBody)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}

				a.Client = &http.Client{
					Transport: &mockRoundTripper{response: resp},
				}

				test.Request = &app.TestRequest{
					Method:  "POST",
					URL:     "/api/users",
					Headers: map[string]string{"Content-Type": "application/json"},
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err) // Function should still succeed even with 4xx status
				assert.NotNil(t, a.LastResponse)
				assert.Equal(t, 400, *a.LastResponse.Code)
				assert.Equal(t, `{"error": "Invalid request parameters"}`, *a.LastResponse.Body)
			},
		},
		{
			name: "500 Server Error response",
			setup: func(a *app.Abdd, test *app.Test) {
				a.Store = make(map[string]any)
				a.Global = app.Global{
					Config: app.Config{
						BaseURL: "https://example.com",
					},
				}

				errorBody := `{"error": "Internal server error"}`
				resp := &http.Response{
					StatusCode: 500,
					Body:       io.NopCloser(strings.NewReader(errorBody)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}

				a.Client = &http.Client{
					Transport: &mockRoundTripper{response: resp},
				}

				test.Request = &app.TestRequest{
					Method: "GET",
					URL:    "/api/status",
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err) // Function should still succeed even with 5xx status
				assert.NotNil(t, a.LastResponse)
				assert.Equal(t, 500, *a.LastResponse.Code)
				assert.Equal(t, `{"error": "Internal server error"}`, *a.LastResponse.Body)
			},
		},
		{
			name: "Response with multiple header values",
			setup: func(a *app.Abdd, test *app.Test) {
				a.Store = make(map[string]any)
				a.Global = app.Global{
					Config: app.Config{
						BaseURL: "https://example.com",
					},
				}

				resp := &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader("success")),
					Header: http.Header{
						"Set-Cookie":    []string{"session=123", "user=john"},
						"Cache-Control": []string{"no-cache", "no-store"},
					},
				}

				a.Client = &http.Client{
					Transport: &mockRoundTripper{response: resp},
				}

				test.Request = &app.TestRequest{
					Method: "GET",
					URL:    "/api/auth",
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, a.LastResponse)
				assert.Equal(t, 200, *a.LastResponse.Code)
				assert.Equal(t, "session=123, user=john", a.LastResponse.Headers["Set-Cookie"])
				assert.Equal(t, "no-cache, no-store", a.LastResponse.Headers["Cache-Control"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := &app.Abdd{
				Store: make(map[string]any),
			}
			test := &app.Test{}
			tc.setup(a, test)
			a.Tests = []app.Test{*test}

			err := a.MakeRequest(test)

			tc.expects(*a, test, err)
		})
	}
}
