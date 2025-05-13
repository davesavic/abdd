package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/goccy/go-yaml"
)

var (
	ErrUnexpectedStatusCode        = errors.New("unexpected status code")
	ErrHeaderNotFound              = errors.New("header not found")
	ErrHeaderNotEqual              = errors.New("header not equal")
	ErrJsonPathNotFound            = errors.New("json path not found")
	ErrJsonPathNotEqual            = errors.New("json path not equal")
	ErrExtractionPathEmpty         = errors.New("extraction path is empty")
	ErrExtractionVariableNameEmpty = errors.New("extraction variable name is empty")
	ErrExtractionPathNotFound      = errors.New("extraction path not found")
)

var (
	successText = color.New(color.FgGreen, color.Bold).SprintFunc()
	failureText = color.New(color.FgRed, color.Bold).SprintFunc()
	headerText  = color.New(color.FgCyan).SprintFunc()
	infoText    = color.New(color.FgYellow).SprintFunc()
)

type Config struct {
	BaseURL     string            `yaml:"base_url"`
	Headers     map[string]string `yaml:"headers"`
	Timeout     int               `yaml:"timeout"`
	StopOnError bool              `yaml:"stop_on_error"`
	Verbose     bool              `yaml:"verbose"`
}

type Global struct {
	Config Config `yaml:"config"`
}

type Abdd struct {
	Global Global `yaml:"global"`
	Tests  []Test `yaml:"-"`

	Store        map[string]any `yaml:"-"`
	LastResponse *LastResponse  `yaml:"-"`
	Client       *http.Client   `yaml:"-"`
}

type LastResponse struct {
	Body    *string
	Code    *int
	Headers map[string]string
}

type TestRequest struct {
	Method  string            `yaml:"method"`
	URL     string            `yaml:"url"`
	Body    *string           `yaml:"body,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

type TestCommand struct {
	Command   string `yaml:"command"`
	Directory string `yaml:"directory,omitempty"`
	As        string `yaml:"as,omitempty"`
}

type TestExpect struct {
	Headers map[string]string `yaml:"headers,omitempty"`
	Status  *int              `yaml:"status,omitempty"`
	Json    map[string]any    `yaml:"json,omitempty"`
}

type TestExtract struct {
	Path string `yaml:"path"`
	As   string `yaml:"as"`
}

type Test struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Depends     []string          `yaml:"depends,omitempty"`
	Fake        map[string]string `yaml:"fake,omitempty"`
	Request     *TestRequest      `yaml:"request,omitempty"`
	Command     *TestCommand      `yaml:"command,omitempty"`
	Expect      TestExpect        `yaml:"expect"`
	Extract     []TestExtract     `yaml:"extract,omitempty"`
}

type AbddArgs struct {
	ConfigFile string
	Folders    []string
	Verbose    bool
}

func (args *AbddArgs) Validate() error {
	if args.ConfigFile == "" {
		return fmt.Errorf("config file not provided")
	}

	c, err := os.Stat(args.ConfigFile)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("config file %s does not exist", args.ConfigFile)
	}
	if c.IsDir() {
		return fmt.Errorf("config file %s is a directory", args.ConfigFile)
	}

	if len(args.Folders) == 0 {
		return fmt.Errorf("no folders provided")
	}
	for _, folder := range args.Folders {
		i, err := os.Stat(folder)
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("folder %s does not exist", folder)
		}
		if !i.IsDir() {
			return fmt.Errorf("%s is not a directory", folder)
		}
	}
	return nil
}

func New(args AbddArgs) (*Abdd, error) {
	err := args.Validate()
	if err != nil {
		return nil, err
	}

	// Create a new Abdd instance
	a := &Abdd{
		Store:  make(map[string]any),
		Client: http.DefaultClient,
	}

	// Load the global config from the specified file
	err = a.LoadGlobal(args.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load global config: %w", err)
	}

	if a.Global.Config.Timeout != 0 {
		a.Client.Timeout = time.Duration(a.Global.Config.Timeout) * time.Second
	}

	if args.Verbose {
		a.Global.Config.Verbose = true
	}

	// Load tests from the specified folders
	err = a.LoadTests(args.Folders, args.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load tests: %w", err)
	}

	return a, nil
}

func (a *Abdd) LoadGlobal(path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var abdd Abdd
	if err := yaml.Unmarshal(f, &abdd); err != nil {
		return fmt.Errorf("failed to unmarshal config file: %w", err)
	}
	a.Global = abdd.Global
	return nil
}

// LoadTests globs the given folders for *.yaml and *.yml files and loads them into the Abdd instance.
func (a *Abdd) LoadTests(folders []string, exclude string) error {
	var allFiles []string
	for _, folder := range folders {
		df, err := filepath.Glob(filepath.Join(folder, "*.yaml"))
		if err != nil {
			return fmt.Errorf("failed to glob folder %s: %w", folder, err)
		}
		allFiles = append(allFiles, df...)

		df, err = filepath.Glob(filepath.Join(folder, "*.yml"))
		if err != nil {
			return fmt.Errorf("failed to glob folder %s: %w", folder, err)
		}

		allFiles = append(allFiles, df...)
	}

	var tests []Test
	for _, file := range allFiles {
		if file == exclude {
			continue
		}

		f, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read test file %s: %w", file, err)
		}
		var testFile struct {
			Tests []Test `yaml:"tests"`
		}
		if err := yaml.Unmarshal(f, &testFile); err != nil {
			return fmt.Errorf("failed to unmarshal test file %s: %w", file, err)
		}

		tests = append(tests, testFile.Tests...)
	}

	// Create a map of test names to tests
	testMap := make(map[string]Test)
	for _, test := range tests {
		testMap[test.Name] = test
	}

	// Check for missing dependencies
	for _, test := range tests {
		for _, dep := range test.Depends {
			if _, ok := testMap[dep]; !ok {
				return fmt.Errorf("test '%s' depends on non-existent test '%s'", test.Name, dep)
			}
		}
	}

	// Perform topological sort
	var sorted []Test
	visited := make(map[string]bool)
	tempMark := make(map[string]bool)

	var visit func(string) error
	visit = func(name string) error {
		if tempMark[name] {
			return fmt.Errorf("circular dependency detected involving test '%s'", name)
		}
		if visited[name] {
			return nil
		}

		tempMark[name] = true

		// Visit all dependencies first
		test := testMap[name]
		for _, dep := range test.Depends {
			if err := visit(dep); err != nil {
				return err
			}
		}

		tempMark[name] = false
		visited[name] = true
		sorted = append(sorted, test)

		return nil
	}

	// Visit all tests
	for _, test := range tests {
		if !visited[test.Name] {
			if err := visit(test.Name); err != nil {
				return err
			}
		}
	}

	a.Tests = sorted
	return nil
}

func (a *Abdd) Run() error {
	fmt.Println(headerText("┌─────────────────────────────────┐"))
	fmt.Println(headerText("               Tests               "))

	totalTests := len(a.Tests)
	passedTests := 0
	failedTests := 0

	for i, test := range a.Tests {
		if a.Global.Config.Verbose {
			a.PrintStartTest(&test)
		}

		var err error
		err = a.GenerateFakeData(&test)
		if err == nil {
			if a.Global.Config.Verbose {
				a.PrintGenerateFakeDataStep(&test)
			}
			err = a.ReplaceVariables(&test)
		}

		if err == nil {
			if a.Global.Config.Verbose {
				a.PrintReplaceVariablesStep(&test)
			}
			err = a.ExecuteCommand(&test)
		}

		if err == nil {
			if a.Global.Config.Verbose {
				a.PrintExecuteCommandStep(&test)
			}
			err = a.MakeRequest(&test)
		}

		if err == nil {
			if a.Global.Config.Verbose {
				a.PrintMakeRequestStep(&test)
			}
			err = a.ValidateResponse(&test)
		}

		if err == nil {
			if a.Global.Config.Verbose {
				a.PrintValidateResponseStep(&test)
			}
			err = a.ExtractData(&test)
		}

		if err == nil {
			if a.Global.Config.Verbose {
				a.PrintExtractDataStep(&test)
			}
		}

		if err == nil {
			passedTests++

			fmt.Printf("[%d/%d] %s %s\n", i+1, totalTests, successText("✓"), test.Name)
			continue
		}

		failedTests++
		fmt.Printf("[%d/%d] %s %s\n", i+1, totalTests, failureText("✗"), test.Name)
		fmt.Printf("       %s %v\n", failureText("→"), err)

		a.PrintFailureDetails(&test)

		if a.Global.Config.StopOnError {
			break
		}
	}

	fmt.Println()
	fmt.Println(headerText("└─────────────────────────────────┘"))

	fmt.Println(headerText("┌─────────────────────────────────┐"))
	fmt.Println(headerText("              Summary              "))

	totalStr := fmt.Sprintf("Total: %d", totalTests)
	passedStr := fmt.Sprintf("%s: %d", successText("Passed"), passedTests)
	failedStr := fmt.Sprintf("%s: %d", failureText("Failed"), failedTests)
	rateStr := fmt.Sprintf("Pass rate: %.1f%%", float64(passedTests)/float64(totalTests)*100)

	fmt.Println(totalStr)
	fmt.Println(passedStr)
	fmt.Println(failedStr)
	fmt.Println(rateStr)

	fmt.Println()
	fmt.Println(headerText("└─────────────────────────────────┘"))

	if failedTests > 0 && !a.Global.Config.StopOnError {
		return fmt.Errorf(failureText("%d tests failed"), failedTests)
	}

	return nil
}
