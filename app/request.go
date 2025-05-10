package app

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (a *Abdd) MakeRequest(t *Test) error {
	fmt.Println("Making request...")

	var payload *bytes.Buffer
	if t.Request.Body != nil {
		payload = bytes.NewBufferString(*t.Request.Body)
	}

	req, err := http.NewRequest(t.Request.Method, a.Global.Config.BaseURL+t.Request.URL, payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	for key, value := range t.Request.Headers {
		req.Header.Set(key, value)
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respHeaders := map[string]string{}
	for key, values := range resp.Header {
		if len(values) > 0 {
			respHeaders[key] = strings.Join(values, ", ")
		}
	}

	respBody := string(bodyBytes)
	lr := LastResponse{
		Headers: respHeaders,
		Body:    &respBody,
		Code:    &resp.StatusCode,
	}
	a.LastResponse = &lr

	return nil
}
