package app_test

import (
	"testing"

	"github.com/davesavic/abdd/app"
	"github.com/stretchr/testify/assert"
)

func TestGenerateFakeData(t *testing.T) {
	testCases := []struct {
		name        string
		test        app.Test
		initStore   map[string]any
		expects     func(app.Abdd)
		expectedErr error
	}{
		{
			name: "Valid fake data generation",
			test: app.Test{
				Fake: map[string]string{
					"customerEmail": "{email}",
				},
			},
			expects: func(a app.Abdd) {
				assert.NotEmpty(t, a.Store["customerEmail"])
				assert.Contains(t, a.Store["customerEmail"].(string), "@")
			},
		},
		{
			name: "Invalid fake data tag returns same string",
			test: app.Test{
				Fake: map[string]string{
					"customerEmail": "{invalidTag}",
				},
			},
			expects: func(a app.Abdd) {
				assert.Equal(t, "{invalidTag}", a.Store["customerEmail"])
			},
		},
		{
			name: "Empty fake data map",
			test: app.Test{
				Fake: map[string]string{},
			},
			expects: func(a app.Abdd) {
				assert.Empty(t, a.Store)
			},
		},
		{
			name: "Mixed valid tag and invalid tag",
			test: app.Test{
				Fake: map[string]string{
					"customerEmail": "{email}",
					"customerName":  "{invalidTag}",
				},
			},
			expects: func(a app.Abdd) {
				assert.NotEmpty(t, a.Store["customerEmail"])
				assert.Contains(t, a.Store["customerEmail"].(string), "@")
				assert.Equal(t, "{invalidTag}", a.Store["customerName"])
			},
		},
		{
			name: "Valid tag with static string",
			test: app.Test{
				Fake: map[string]string{
					"customerEmail": "{email} - static",
				},
			},
			expects: func(a app.Abdd) {
				assert.NotEmpty(t, a.Store["customerEmail"])
				assert.Contains(t, a.Store["customerEmail"].(string), "@")
				assert.Contains(t, a.Store["customerEmail"].(string), " - static")
			},
		},
		{
			name: "Overwrite existing store value",
			test: app.Test{
				Fake: map[string]string{
					"customerEmail": "{email}",
				},
			},
			initStore: map[string]any{
				"customerEmail": "oldValue",
			},
			expects: func(a app.Abdd) {
				assert.NotEmpty(t, a.Store["customerEmail"])
				assert.Contains(t, a.Store["customerEmail"].(string), "@")
				assert.NotEqual(t, "oldValue", a.Store["customerEmail"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := make(map[string]any)
			if tc.initStore != nil {
				store = tc.initStore
			}

			a := app.Abdd{
				Tests: []app.Test{tc.test},
				Store: store,
			}

			err := a.GenerateFakeData(&tc.test)
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				tc.expects(a)
			}
		})
	}
}
