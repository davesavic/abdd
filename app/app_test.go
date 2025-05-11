package app_test

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/davesavic/abdd/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAbddArgsValidate(t *testing.T) {
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
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")
	testFolder := filepath.Join(tempDir, "tests")
	testFile := filepath.Join(testFolder, "test1.yaml")

	configContent := `
global:
  config:
    base_url: https://jsonplaceholder.typicode.com
    headers:
      Content-Type: application/json
    timeout: 30`

	err := os.WriteFile(configFile, []byte(configContent), 0o644)
	assert.NoError(t, err)

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
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")
	invalidConfigFile := filepath.Join(tempDir, "invalid.yaml")

	validConfig := `
global:
  config:
    base_url: https://jsonplaceholder.typicode.com
    headers:
      Content-Type: application/json
    timeout: 30`

	err := os.WriteFile(configFile, []byte(validConfig), 0o644)
	assert.NoError(t, err)

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

	tests := []struct {
		name        string
		setup       func() []string
		wantTests   int
		wantErr     bool
		validate    func(t *testing.T, tests []app.Test)
		validateErr func(t *testing.T, err error)
	}{
		{
			name: "Single folder",
			setup: func() []string {
				folder := filepath.Join(tempDir, "tests1")
				assert.NoError(t, os.Mkdir(folder, 0o755))

				content := `tests:
- name: Test1
  description: First test
  request:
    method: GET
    url: /api/resource1
  expect:
    status: 200`

				assert.NoError(t, os.WriteFile(filepath.Join(folder, "test1.yaml"), []byte(content), 0o644))
				return []string{folder}
			},
			wantTests: 1,
			validate: func(t *testing.T, tests []app.Test) {
				assert.Equal(t, "Test1", tests[0].Name)
				assert.Equal(t, "GET", tests[0].Request.Method)
			},
		},
		{
			name: "Multiple folders",
			setup: func() []string {
				folder1 := filepath.Join(tempDir, "multi1")
				folder2 := filepath.Join(tempDir, "multi2")
				require.NoError(t, os.Mkdir(folder1, 0o755))
				require.NoError(t, os.Mkdir(folder2, 0o755))

				content1 := `tests:
- name: Test1
  description: First test
  request:
    method: GET
    url: /api/resource1
  expect:
    status: 200`

				content2 := `tests:
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

				require.NoError(t, os.WriteFile(filepath.Join(folder1, "test1.yaml"), []byte(content1), 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(folder2, "test2.yaml"), []byte(content2), 0o644))
				return []string{folder1, folder2}
			},
			wantTests: 3,
			validate: func(t *testing.T, tests []app.Test) {
				names := make(map[string]bool)
				for _, test := range tests {
					names[test.Name] = true
				}
				assert.True(t, names["Test1"])
				assert.True(t, names["Test2"])
				assert.True(t, names["Test3"])
			},
		},
		{
			name: "Invalid test file",
			setup: func() []string {
				folder := filepath.Join(tempDir, "invalid")
				require.NoError(t, os.Mkdir(folder, 0o755))

				validContent := `tests:
- name: Test2
  description: Second test
  request:
    method: POST
    url: /api/resource2
  expect:
    status: 201`

				invalidContent := `invalid: yaml: :`

				require.NoError(t, os.WriteFile(filepath.Join(folder, "valid.yaml"), []byte(validContent), 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(folder, "invalid.yaml"), []byte(invalidContent), 0o644))
				return []string{folder}
			},
			wantTests: 0,
			wantErr:   true,
		},
		{
			name: "Non-existent folder",
			setup: func() []string {
				return []string{filepath.Join(tempDir, "nonexistent")}
			},
			wantTests: 0,
			wantErr:   false,
		},
		{
			name: "yaml and yml files",
			setup: func() []string {
				folder1 := filepath.Join(tempDir, "yml1")
				folder2 := filepath.Join(tempDir, "yml2")
				require.NoError(t, os.Mkdir(folder1, 0o755))
				require.NoError(t, os.Mkdir(folder2, 0o755))

				content1 := `tests:
- name: Test1
  description: First test
  request:
    method: GET
    url: /api/resource1
  expect:
    status: 200`

				content2 := `tests:
- name: Test4
  description: Fourth test
  request:
    method: GET
    url: /api/resource4
  expect:
    status: 200`

				require.NoError(t, os.WriteFile(filepath.Join(folder1, "test1.yaml"), []byte(content1), 0o644))
				require.NoError(t, os.WriteFile(filepath.Join(folder2, "test4.yaml"), []byte(content2), 0o644))
				return []string{folder1, folder2}
			},
			wantTests: 2,
		},
		{
			name: "ordered by dependency",
			setup: func() []string {
				folder := filepath.Join(tempDir, "ordered")
				require.NoError(t, os.Mkdir(folder, 0o755))

				content := `tests:
  - name: Test1
    description: First test
    depends:
      - Test2
    request:
      method: GET
      url: /api/resource1
    expect:
      status: 200
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
    depends:
      - Test1
    request:
      method: DELETE
      url: /api/resource3
    expect:
      status: 204`

				assert.NoError(t, os.WriteFile(filepath.Join(folder, "ordered.yaml"), []byte(content), 0o644))
				return []string{folder}
			},
			wantTests: 3,
			validate: func(t *testing.T, tests []app.Test) {
				assert.Equal(t, "Test2", tests[0].Name)
				assert.Equal(t, "Test1", tests[1].Name)
				assert.Equal(t, "Test3", tests[2].Name)
			},
		},
		{
			name: "circular dependency",
			setup: func() []string {
				folder := filepath.Join(tempDir, "circular")
				require.NoError(t, os.Mkdir(folder, 0o755))

				content := `tests:
  - name: Test1
    description: First test
    depends:
      - Test2
    request:
      method: GET
      url: /api/resource1
    expect:
      status: 200
  - name: Test2
    description: Second test
    depends:
      - Test1
    request:
      method: POST
      url: /api/resource2
      body: '{"key": "value"}'
    expect:
      status: 201
  - name: Test3
    description: Third test
    depends:
      - Test1
    request:
      method: DELETE
      url: /api/resource3
    expect:
      status: 204`

				assert.NoError(t, os.WriteFile(filepath.Join(folder, "circular.yaml"), []byte(content), 0o644))
				return []string{folder}
			},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "circular dependency detected")
			},
		},
		{
			name: "missing dependency",
			setup: func() []string {
				folder := filepath.Join(tempDir, "missing")
				require.NoError(t, os.Mkdir(folder, 0o755))

				content := `tests:
  - name: Test1
    description: First test
    depends:
      - Test2
    request:
      method: GET
      url: /api/resource1
    expect:
      status: 200
  - name: Test2
    description: Second test
    depends:
      - Test100
    request:
      method: POST
      url: /api/resource2
      body: '{"key": "value"}'
    expect:
      status: 201
  - name: Test3
    description: Third test
    depends:
      - Test1
    request:
      method: DELETE
      url: /api/resource3
    expect:
      status: 204`

				assert.NoError(t, os.WriteFile(filepath.Join(folder, "missing.yaml"), []byte(content), 0o644))
				return []string{folder}
			},
			wantErr: true,
			validateErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "test 'Test2' depends on non-existent test 'Test100'")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			folders := tt.setup()

			a := &app.Abdd{}
			err := a.LoadTests(folders, "")

			if tt.wantErr {
				assert.Error(t, err)
				if tt.validateErr != nil {
					tt.validateErr(t, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, a.Tests, tt.wantTests)

				if tt.validate != nil && tt.wantTests > 0 {
					tt.validate(t, a.Tests)
				}
			}
		})
	}
}

func TestRun(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"id": 1, "name": "John Doe"}`)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}

	test := &app.Test{
		Name: "Test1",
		Request: &app.TestRequest{
			Method:  "GET",
			URL:     "/api/users",
			Headers: map[string]string{"Accept": "application/json"},
		},
		Expect: app.TestExpect{
			Status: toPointer(200),
		},
		Extract: []app.TestExtract{
			{Path: "id", As: "userId"},
			{Path: "name", As: "userName"},
		},
	}

	a := app.Abdd{
		Global: app.Global{
			Config: app.Config{
				BaseURL: "https://example.com",
			},
		},
		Tests: []app.Test{*test},
		Store: map[string]any{},
		Client: &http.Client{
			Transport: &app.TestMockRoundTripper{Response: resp},
		},
	}

	err := a.Run()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(a.Tests))
	assert.Equal(t, 200, *a.LastResponse.Code)
	assert.Equal(t, "application/json", a.LastResponse.Headers["Content-Type"])
	assert.Equal(t, "{\"id\": 1, \"name\": \"John Doe\"}", *a.LastResponse.Body)
	assert.Equal(t, "Test1", a.Tests[0].Name)
	assert.Equal(t, "GET", a.Tests[0].Request.Method)
	assert.Equal(t, "/api/users", a.Tests[0].Request.URL)
	assert.Equal(t, "application/json", a.Tests[0].Request.Headers["Accept"])
	assert.Equal(t, 2, len(a.Store))
}
