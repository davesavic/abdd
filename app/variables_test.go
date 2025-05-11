package app_test

import (
	"testing"

	"github.com/davesavic/abdd/app"
	"github.com/stretchr/testify/assert"
)

func TestReplaceVariables(t *testing.T) {
	a := app.Abdd{
		Store: map[string]any{
			"customerEmail": "testing@gmail.com",
		},
	}

	testCases := []struct {
		name    string
		setup   func(*app.Abdd, *app.Test)
		expects func(app.Abdd, *app.Test, error)
	}{
		{
			name: "Replace variables in request headers",
			setup: func(a *app.Abdd, test *app.Test) {
				a.Store["customerEmail"] = "jennyfromtheblock@gmail.com"

				test.Request = &app.TestRequest{
					Headers: map[string]string{
						"Customer-Email": "${customerEmail}",
					},
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err)
				assert.Equal(t, a.Store["customerEmail"], test.Request.Headers["Customer-Email"])
			},
		},
		{
			name: "Replace variables in request URL",
			setup: func(a *app.Abdd, test *app.Test) {
				a.Store["customerEmail"] = "hello@gmail.com"
				test.Request = &app.TestRequest{
					URL: "/users?search=${customerEmail}",
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "/users?search=hello@gmail.com", test.Request.URL)
			},
		},
		{
			name: "Replace variables in request body",
			setup: func(a *app.Abdd, test *app.Test) {
				a.Store["customerEmail"] = "nosugar@gmail.com"
				body := "{\"email\": \"${customerEmail}\"}"
				test.Request = &app.TestRequest{
					Body: &body,
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "{\"email\": \"nosugar@gmail.com\"}", *test.Request.Body)
			},
		},
		{
			name: "Replace variables in command",
			setup: func(a *app.Abdd, test *app.Test) {
				a.Store["customerEmail"] = "testing@gmail.com"
				test.Command = &app.TestCommand{
					Command: "echo ${customerEmail}",
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "echo testing@gmail.com", test.Command.Command)
			},
		},
		{
			name: "Replace variables in expect headers",
			setup: func(a *app.Abdd, test *app.Test) {
				a.Store["token"] = "1234567890"
				test.Expect.Headers = map[string]string{
					"Authorization": "Bearer ${token}",
				}
				test.Request = &app.TestRequest{}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "Bearer 1234567890", test.Expect.Headers["Authorization"])
			},
		},
		{
			name: "Replace variables in expect JSON",
			setup: func(a *app.Abdd, test *app.Test) {
				a.Store["customerEmail"] = "chimmy@gmail.com"
				test.Expect.Json = map[string]any{
					"extractedKey": "${customerEmail}",
				}
				test.Request = &app.TestRequest{}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "chimmy@gmail.com", test.Expect.Json["extractedKey"])
			},
		},
		{
			name:  "Test with no request or command",
			setup: func(a *app.Abdd, test *app.Test) {},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.Error(t, err)
				assert.Equal(t, "test must have either a request or a command", err.Error())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initTest := app.Test{}
			tc.setup(&a, &initTest)

			err := a.ReplaceVariables(&initTest)

			tc.expects(a, &initTest, err)
		})
	}
}
