package app

import (
	"fmt"

	"github.com/tidwall/gjson"
)

func (a *Abdd) ValidateResponse(t *Test) error {
	fmt.Println("Validating response...")

	if a.LastResponse == nil {
		return fmt.Errorf("no response to validate")
	}

	if t.Expect.Status != nil && a.LastResponse.Code != nil {
		if *a.LastResponse.Code != *t.Expect.Status {
			return fmt.Errorf("w%: expected %d, got %d", ErrUnexpectedStatusCode, *t.Expect.Status, *a.LastResponse.Code)
		}
	}

	if t.Expect.Headers != nil {
		for key, expectedValue := range t.Expect.Headers {
			actualValue, exists := a.LastResponse.Headers[key]
			if !exists {
				return fmt.Errorf("w%: expected %s to be present", ErrHeaderNotFound, key)
			}
			if actualValue != expectedValue {
				return fmt.Errorf("w%: expected header %s to be %s, got %s", ErrHeaderNotEqual, key, expectedValue, actualValue)
			}
		}
	}

	if t.Expect.Json != nil && a.LastResponse.Body != nil {
		for key, expectedValue := range t.Expect.Json {
			actualValue := gjson.Get(*a.LastResponse.Body, key)
			if !actualValue.Exists() {
				return fmt.Errorf("w%: expected %s to be present", ErrJsonPathNotFound, key)
			}
			if actualValue.String() != fmt.Sprintf("%v", expectedValue) {
				return fmt.Errorf("w%: expected %s to be %v, got %v", ErrJsonPathNotEqual, key, expectedValue, actualValue.String())
			}
		}
	}

	return nil
}
