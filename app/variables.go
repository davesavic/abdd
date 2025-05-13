package app

import (
	"fmt"
	"regexp"
)

func (a *Abdd) ReplaceVariables(t *Test) error {
	if t == nil {
		return fmt.Errorf("test cannot be nil")
	}
	if t.Request == nil && t.Command == nil {
		return fmt.Errorf("test must have either a request or a command")
	}

	if t.Command != nil {
		if t.Command.Directory != "" {
			t.Command.Directory = a.replaceVariablesInText(t.Command.Directory)
		}
		t.Command.Command = a.replaceVariablesInText(t.Command.Command)
	}

	if t.Request != nil {
		headers := map[string]string{}
		for key, value := range t.Request.Headers {
			headers[key] = a.replaceVariablesInText(value)
		}

		var body *string
		if t.Request.Body != nil {
			bodyValue := a.replaceVariablesInText(*t.Request.Body)
			body = &bodyValue
		}

		t.Request = &TestRequest{
			Method:  t.Request.Method,
			URL:     a.replaceVariablesInText(t.Request.URL),
			Body:    body,
			Headers: headers,
		}
	}

	if t.Expect.Headers != nil {
		headers := map[string]string{}
		for key, value := range t.Expect.Headers {
			headers[key] = a.replaceVariablesInText(value)
		}
		t.Expect.Headers = headers
	}

	if t.Expect.Json != nil {
		json := map[string]any{}
		for key, value := range t.Expect.Json {
			json[key] = a.replaceVariablesInText(fmt.Sprintf("%v", value))
		}
		t.Expect.Json = json
	}

	return nil
}

func (a *Abdd) replaceVariablesInText(text string) string {
	r := regexp.MustCompile(`\${([^}]+)}`)
	return r.ReplaceAllStringFunc(text, func(match string) string {
		// Extract key name without ${ and }
		key := match[2 : len(match)-1]
		if val, ok := a.Store[key]; ok {
			return fmt.Sprintf("%v", val)
		}
		return match
	})
}
