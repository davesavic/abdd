package app_test

import (
	"os"
	"testing"

	"github.com/davesavic/abdd/app"
	"github.com/stretchr/testify/assert"
)

func TestExecuteCommand(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(*app.Abdd, *app.Test)
		expects func(app.Abdd, *app.Test, error)
	}{
		{
			name: "Valid command execution",
			setup: func(a *app.Abdd, test *app.Test) {
				test.Command = &app.TestCommand{
					Command: "echo Hello World",
					As:      "greeting",
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "Hello World", a.Store["greeting"])
			},
		},
		{
			name: "Valid command execution with directory",
			setup: func(a *app.Abdd, test *app.Test) {
				test.Command = &app.TestCommand{
					Command:   "echo Hello from temp directory",
					Directory: "/tmp",
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err)
				assert.Empty(t, a.Store)
			},
		},
		{
			name: "Valid command execution with directory and save output",
			setup: func(a *app.Abdd, test *app.Test) {
				tempDir := t.TempDir()
				testFile := tempDir + "/testfile.txt"
				err := os.WriteFile(testFile, []byte("Hellooooo there"), 0o644)
				assert.NoError(t, err)

				test.Command = &app.TestCommand{
					Command:   "cat testfile.txt",
					Directory: tempDir,
					As:        "greeting",
				}
			},
			expects: func(a app.Abdd, test *app.Test, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "Hellooooo there", a.Store["greeting"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := app.Abdd{
				Store: map[string]any{},
			}
			test := &app.Test{}
			tc.setup(&a, test)

			err := a.ExecuteCommand(test)
			tc.expects(a, test, err)
		})
	}
}
