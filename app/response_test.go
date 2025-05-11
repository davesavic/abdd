package app_test

import (
	"fmt"
	"testing"

	"github.com/davesavic/abdd/app"
)

func toPointer[T any](v T) *T {
	return &v
}

func TestValidateResponse(t *testing.T) {
	testCases := []struct {
		name         string
		test         app.Test
		lastResponse *app.LastResponse
		expectedErr  error
	}{
		{
			name: "Valid response with code, headers, and json body",
			test: app.Test{
				Expect: app.TestExpect{
					Status: toPointer(200),
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Json: map[string]any{
						"key": "value",
					},
				},
			},
			lastResponse: &app.LastResponse{
				Code: toPointer(200),
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: toPointer(`{"key": "value"}`),
			},
		},
		{
			name: "Unexpected status code",
			test: app.Test{
				Expect: app.TestExpect{
					Status: toPointer(200),
				},
			},
			lastResponse: &app.LastResponse{
				Code: toPointer(404),
			},
			expectedErr: fmt.Errorf("%w: expected %d, got %d", app.ErrUnexpectedStatusCode, 200, 404),
		},
		{
			name: "Header not found",
			test: app.Test{
				Expect: app.TestExpect{
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
			lastResponse: &app.LastResponse{},
			expectedErr:  fmt.Errorf("%w: expected %s to be present", app.ErrHeaderNotFound, "Content-Type"),
		},
		{
			name: "Header mismatch",
			test: app.Test{
				Expect: app.TestExpect{
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
			lastResponse: &app.LastResponse{
				Headers: map[string]string{
					"Content-Type": "text/html",
				},
			},
			expectedErr: fmt.Errorf("%w: expected header %s to be %s, got %s", app.ErrHeaderNotEqual, "Content-Type", "application/json", "text/html"),
		},
		{
			name: "JSON path not found",
			test: app.Test{
				Expect: app.TestExpect{
					Json: map[string]any{
						"key": "value",
					},
				},
			},
			lastResponse: &app.LastResponse{
				Body: toPointer(`{"other_key": "other_value"}`),
			},
			expectedErr: fmt.Errorf("%w: expected %s to be present", app.ErrJsonPathNotFound, "key"),
		},
		{
			name: "JSON value mismatch",
			test: app.Test{
				Expect: app.TestExpect{
					Json: map[string]any{
						"key": "value",
					},
				},
			},
			lastResponse: &app.LastResponse{
				Body: toPointer(`{"key": "other_value"}`),
			},
			expectedErr: fmt.Errorf("%w: expected %s to be %v, got %v", app.ErrJsonPathNotEqual, "key", "value", "other_value"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := &app.Abdd{
				Tests:        []app.Test{tc.test},
				LastResponse: tc.lastResponse,
			}

			err := a.ValidateResponse(&a.Tests[0])
			if err != nil && err.Error() != tc.expectedErr.Error() {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
