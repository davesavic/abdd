package app

import "fmt"

func (a *Abdd) PrintStartTest(t *Test) {
	fmt.Printf("\n%s %s\n", infoText("▶"), t.Name)
	fmt.Printf("  %s: %s\n", infoText("Description"), t.Description)
}

func (a *Abdd) PrintGenerateFakeDataStep(t *Test) {
	fmt.Printf("  %s Generated fake data\n", infoText("•"))
	fmt.Printf("  %s: %s\n", infoText("Fake"), t.Fake)
}

func (a *Abdd) PrintReplaceVariablesStep(t *Test) {
	fmt.Printf("  %s Replaced variables\n", infoText("•"))
	fmt.Printf("  %s: %+v\n", infoText("Request"), t.Request)
}

func (a *Abdd) PrintExecuteCommandStep(t *Test) {
	fmt.Printf("  %s Executed command\n", infoText("•"))
	fmt.Printf("  %s: %+v\n", infoText("Command"), t.Command)
}

func (a *Abdd) PrintMakeRequestStep(t *Test) {
	fmt.Printf("  %s Made request\n", infoText("•"))
	fmt.Printf("  %s: [%s]%s %+v %+v\n", infoText("Request"), t.Request.Method, t.Request.URL, *t.Request.Body, t.Request.Headers)
	fmt.Printf("  %s: %+v\n", infoText("Store"), a.Store)
}

func (a *Abdd) PrintValidateResponseStep(t *Test) {
	fmt.Printf("  %s Validated response\n", infoText("•"))
	fmt.Printf("  %s: %+v\n", infoText("Response"), a.LastResponse)
}

func (a *Abdd) PrintExtractDataStep(t *Test) {
	fmt.Printf("  %s Extracted data\n", infoText("•"))
	fmt.Printf("  %s: %+v\n", infoText("Extracted"), a.Store)
}

func (a *Abdd) PrintFailureDetails(t *Test) {
	fmt.Println(failureText("\n❯ Test Failure Details:"))

	fmt.Printf("  %s: %s\n", infoText("Test"), t.Name)
	fmt.Printf("  %s: %s\n", infoText("Description"), t.Description)

	if t.Request != nil {
		fmt.Printf("\n  %s:\n", infoText("Request"))
		fmt.Printf("    %s: %s\n", infoText("Method"), t.Request.Method)
		fmt.Printf("    %s: %s\n", infoText("URL"), a.Global.Config.BaseURL+t.Request.URL)

		if t.Request.Headers != nil {
			fmt.Printf("    %s:\n", infoText("Headers"))
			for k, v := range t.Request.Headers {
				fmt.Printf("      %s: %s\n", k, v)
			}
		}

		if t.Request.Body != nil {
			fmt.Printf("    %s: %s\n", infoText("Body"), *t.Request.Body)
		}
	}

	if t.Command != nil {
		fmt.Printf("    %s:\n", infoText("Command"))
		fmt.Printf("    %s: %s\n", infoText("Command"), t.Command.Command)
		if t.Command.Directory != "" {
			fmt.Printf("    %s: %s\n", infoText("Directory"), t.Command.Directory)
		}
	}

	if a.LastResponse != nil {
		fmt.Printf("  %s:\n", infoText("Response"))
		if a.LastResponse.Code != nil {
			fmt.Printf("    %s: %d\n", infoText("Status"), *a.LastResponse.Code)
		}

		if a.LastResponse.Headers != nil {
			fmt.Printf("    %s:\n", infoText("Headers"))
			for k, v := range a.LastResponse.Headers {
				fmt.Printf("      %s: %s\n", k, v)
			}
		}

		if a.LastResponse.Body != nil {
			fmt.Printf("    %s: %s\n", infoText("Body"), *a.LastResponse.Body)
		}
	}

	if len(a.Store) > 0 {
		fmt.Printf("\n  %s:\n", infoText("Store"))
		for k, v := range a.Store {
			fmt.Printf("    %s: %v\n", k, v)
		}
	}

	fmt.Printf("\n  %s:\n", infoText("Expected"))
	if t.Expect.Status != nil {
		fmt.Printf("    %s: %d\n", infoText("Status"), *t.Expect.Status)
	}

	if t.Expect.Headers != nil {
		fmt.Printf("    %s:\n", infoText("Headers"))
		for k, v := range t.Expect.Headers {
			fmt.Printf("      %s: %s\n", k, v)
		}
	}

	if t.Expect.Json != nil {
		fmt.Printf("    %s:\n", infoText("JSON"))
		for k, v := range t.Expect.Json {
			fmt.Printf("      %s: %v\n", k, v)
		}
	}

	fmt.Println()
}
