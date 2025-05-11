package app_test

import (
	"fmt"
	"testing"

	"github.com/davesavic/abdd/app"
)

func TestExtractData(t *testing.T) {
	testCases := []struct {
		name          string
		test          app.Test
		lastResponse  *app.LastResponse
		expectedStore map[string]any
		expectedErr   error
	}{
		{
			name: "Valid extraction",
			test: app.Test{
				Extract: []app.TestExtract{
					{
						Path: "key",
						As:   "extractedKey",
					},
				},
			},
			lastResponse: &app.LastResponse{
				Body: toPointer(`{"key": "value"}`),
			},
			expectedStore: map[string]any{
				"extractedKey": "value",
			},
		},
		{
			name: "Extraction path does not exist",
			test: app.Test{
				Extract: []app.TestExtract{
					{
						Path: "nonexistentKey",
						As:   "extractedKey",
					},
				},
			},
			lastResponse: &app.LastResponse{
				Body: toPointer(`{"key": "value"}`),
			},
			expectedErr: fmt.Errorf("%w: expected %s to be present", app.ErrExtractionPathNotFound, "nonexistentKey"),
		},
		{
			name: "Empty extraction path",
			test: app.Test{
				Extract: []app.TestExtract{
					{
						Path: "",
						As:   "extractedKey",
					},
				},
			},
			lastResponse: &app.LastResponse{
				Body: toPointer(`{"key": "value"}`),
			},
			expectedErr: fmt.Errorf("%w: extraction path cannot be empty", app.ErrExtractionPathNotFound),
		},
		{
			name: "Empty variable name",
			test: app.Test{
				Extract: []app.TestExtract{
					{
						Path: "key",
						As:   "",
					},
				},
			},
			lastResponse: &app.LastResponse{
				Body: toPointer(`{"key": "value"}`),
			},
			expectedErr: fmt.Errorf("%w: extraction variable name cannot be empty", app.ErrExtractionVariableNameEmpty),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := app.Abdd{
				LastResponse: tc.lastResponse,
				Store:        make(map[string]any),
			}

			err := a.ExtractData(&tc.test)
			if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}

			if tc.expectedStore != nil {
				for key, expectedValue := range tc.expectedStore {
					actualValue, exists := a.Store[key]
					if !exists {
						t.Errorf("expected key %s to be present in store", key)
					}
					if actualValue != expectedValue {
						t.Errorf("expected key %s to be %v, got %v", key, expectedValue, actualValue)
					}
				}
			}
		})
	}
}
