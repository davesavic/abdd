package app_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/davesavic/abdd/app"
	"github.com/stretchr/testify/assert"
)

func TestAbddArgsValidate(t *testing.T) {
	// Create temporary test directories and files
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")
	testFolder := filepath.Join(tempDir, "tests")

	err := os.WriteFile(configFile, []byte("global:\n  base_url: http://example.com\n  timeout: 30"), 0o644)
	assert.NoError(t, err)

	err = os.Mkdir(testFolder, 0o755)
	assert.NoError(t, err)

	notADir := filepath.Join(tempDir, "file.txt")
	err = os.WriteFile(notADir, []byte("not a directory"), 0o644)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		args    app.AbddArgs
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Empty config file",
			args:    app.AbddArgs{ConfigFile: "", Folders: []string{testFolder}},
			wantErr: true,
			errMsg:  "config file not provided",
		},
		{
			name:    "Config file doesn't exist",
			args:    app.AbddArgs{ConfigFile: filepath.Join(tempDir, "nonexistent.yaml"), Folders: []string{testFolder}},
			wantErr: true,
			errMsg:  "config file " + filepath.Join(tempDir, "nonexistent.yaml") + " does not exist",
		},
		{
			name:    "Config file is a directory",
			args:    app.AbddArgs{ConfigFile: testFolder, Folders: []string{testFolder}},
			wantErr: true,
			errMsg:  "config file " + testFolder + " is a directory",
		},
		{
			name:    "No folders provided",
			args:    app.AbddArgs{ConfigFile: configFile, Folders: []string{}},
			wantErr: true,
			errMsg:  "no folders provided",
		},
		{
			name:    "Folder doesn't exist",
			args:    app.AbddArgs{ConfigFile: configFile, Folders: []string{filepath.Join(tempDir, "nonexistent")}},
			wantErr: true,
			errMsg:  "folder " + filepath.Join(tempDir, "nonexistent") + " does not exist",
		},
		{
			name:    "Folder is not a directory",
			args:    app.AbddArgs{ConfigFile: configFile, Folders: []string{notADir}},
			wantErr: true,
			errMsg:  notADir + " is not a directory",
		},
		{
			name:    "Valid arguments",
			args:    app.AbddArgs{ConfigFile: configFile, Folders: []string{testFolder}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNew(t *testing.T) {
	// Create temporary test directories and files
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")
	testFolder := filepath.Join(tempDir, "tests")
	testFile := filepath.Join(testFolder, "test1.yaml")

	// Create config file
	configContent := `
global:
  config:
    base_url: https://jsonplaceholder.typicode.com
    headers:
      Content-Type: application/json
    timeout: 30`

	err := os.WriteFile(configFile, []byte(configContent), 0o644)
	assert.NoError(t, err)

	// Create test folder and test file
	err = os.Mkdir(testFolder, 0o755)
	assert.NoError(t, err)

	testContent := `tests:
- name: Test1
  description: First test
  request:
    method: GET
    url: /api/resource
  expect:
    status: 200
`
	err = os.WriteFile(testFile, []byte(testContent), 0o644)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		args    app.AbddArgs
		wantErr bool
	}{
		{
			name:    "Invalid arguments",
			args:    app.AbddArgs{ConfigFile: "", Folders: []string{}},
			wantErr: true,
		},
		{
			name:    "Valid arguments",
			args:    app.AbddArgs{ConfigFile: configFile, Folders: []string{testFolder}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := app.New(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, a)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, a)
				assert.Equal(t, "https://jsonplaceholder.typicode.com", a.Global.Config.BaseURL)
				assert.Equal(t, 30, a.Global.Config.Timeout)
				assert.Len(t, a.Tests, 1)
				assert.Equal(t, "Test1", a.Tests[0].Name)
			}
		})
	}
}

func TestLoadGlobal(t *testing.T) {
	// Create temporary test directory and config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")
	invalidConfigFile := filepath.Join(tempDir, "invalid.yaml")

	// Create valid config file
	validConfig := `
global:
  config:
    base_url: https://jsonplaceholder.typicode.com
    headers:
      Content-Type: application/json
    timeout: 30`

	err := os.WriteFile(configFile, []byte(validConfig), 0o644)
	assert.NoError(t, err)

	// Create invalid config file
	invalidConfig := `invalid: yaml: :`
	err = os.WriteFile(invalidConfigFile, []byte(invalidConfig), 0o644)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "File doesn't exist",
			path:    filepath.Join(tempDir, "nonexistent.yaml"),
			wantErr: true,
		},
		{
			name:    "Invalid YAML",
			path:    invalidConfigFile,
			wantErr: true,
		},
		{
			name:    "Valid config",
			path:    configFile,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &app.Abdd{}
			err := a.LoadGlobal(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "https://jsonplaceholder.typicode.com", a.Global.Config.BaseURL)
				assert.Equal(t, 30, a.Global.Config.Timeout)
				assert.Equal(t, "application/json", a.Global.Config.Headers["Content-Type"])
			}
		})
	}
}

func TestLoadTests(t *testing.T) {
	tempDir := t.TempDir()
	testFolder1 := filepath.Join(tempDir, "tests1")
	testFolder2 := filepath.Join(tempDir, "tests2")

	err := os.Mkdir(testFolder1, 0o755)
	assert.NoError(t, err)

	err = os.Mkdir(testFolder2, 0o755)
	assert.NoError(t, err)

	// Create test files in folder1
	test1Content := `tests:
- name: Test1
  description: First test
  request:
    method: GET
    url: /api/resource1
  expect:
    status: 200`

	err = os.WriteFile(filepath.Join(testFolder1, "test1.yaml"), []byte(test1Content), 0o644)
	assert.NoError(t, err)

	// Create test files in folder2
	test2Content := `tests:
- name: Test2
  description: Second test
  request:
    method: POST
    url: /api/resource2
    body: '{"key": "value"}'
  expect:
    status: 201
- name: Test3
  description: Third test
  request:
    method: DELETE
    url: /api/resource3
  expect:
    status: 204`

	err = os.WriteFile(filepath.Join(testFolder2, "test2.yaml"), []byte(test2Content), 0o644)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		folders   []string
		wantTests int
		wantErr   bool
	}{
		{
			name:      "Single folder",
			folders:   []string{testFolder1},
			wantTests: 1,
			wantErr:   false,
		},
		{
			name:      "Multiple folders",
			folders:   []string{testFolder1, testFolder2},
			wantTests: 3,
			wantErr:   false,
		},
		{
			name:      "Invalid test file",
			folders:   []string{testFolder2},
			wantTests: 0,
			wantErr:   true,
		},
		{
			name:      "Non-existent folder",
			folders:   []string{filepath.Join(tempDir, "nonexistent")},
			wantTests: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Only create invalid file for the test that expects it
			if tt.name == "Invalid test file" {
				invalidTestContent := `invalid: yaml: :`
				err = os.WriteFile(filepath.Join(testFolder2, "invalid.yaml"), []byte(invalidTestContent), 0o644)
				assert.NoError(t, err)
				defer os.Remove(filepath.Join(testFolder2, "invalid.yaml"))
			}

			a := &app.Abdd{}
			err := a.LoadTests(tt.folders)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, a.Tests, tt.wantTests)

				if tt.wantTests > 0 && len(a.Tests) > 0 {
					if tt.folders[0] == testFolder1 && len(tt.folders) == 1 {
						assert.Equal(t, "Test1", a.Tests[0].Name)
						assert.Equal(t, "GET", a.Tests[0].Request.Method)
					} else if len(tt.folders) > 1 {
						testNames := make(map[string]bool)
						for _, test := range a.Tests {
							testNames[test.Name] = true
						}
						assert.True(t, testNames["Test1"])
						assert.True(t, testNames["Test2"])
						assert.True(t, testNames["Test3"])
					}
				}
			}
		})
	}
}
