/*
Copyright © 2025 Dave Savic
*/
package cmd

import (
	"os"

	"github.com/davesavic/abdd/app"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a new set of tests",
	Long:  `Initialise a new set of tests for the current project. This command will create a new directory structure and necessary files to get started with testing.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Generate an abdd.yaml file in the current directory using the abdd struct from app.Abdd
		// This will be the default configuration for the project.
		cfgDefaults := app.Abdd{
			Global: app.Global{
				Config: app.Config{
					BaseURL: "https://jsonplaceholder.typicode.com",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Timeout:     30,
					StopOnError: true,
					Verbose:     false,
				},
			},
		}

		yamlData, err := yaml.Marshal(cfgDefaults)
		if err != nil {
			cmd.PrintErrf("Error generating default configuration: %v\n", err)
			return
		}

		err = os.WriteFile("abdd.yaml", yamlData, 0o644)
		if err != nil {
			cmd.PrintErrf("Error writing abdd.yaml file: %v\n", err)
			return
		}

		// Create a tests directory if it doesn't exist
		if _, err = os.Stat("tests"); os.IsNotExist(err) {
			err = os.Mkdir("tests", 0o755)
			if err != nil {
				cmd.PrintErrf("Error creating tests directory: %v\n", err)
				return
			}
		}

		// decalre a struct with a yaml tag for the tests
		type SampleTestFile struct {
			Tests []app.Test `yaml:"tests"`
		}
		sampleTestFile := SampleTestFile{
			Tests: []app.Test{
				{
					Name:        "Create post",
					Description: "This sample test creates a new post",
					Fake: map[string]string{
						"title": "{sentence:5}",
						"body":  "{paragraph:1,3,50,\n}",
					},
					Request: &app.TestRequest{
						Method:  "POST",
						URL:     "/posts",
						Body:    toPointer(`{"title": "${title}", "body": "${body}", "userId": 1}`),
						Headers: map[string]string{},
					},
					Command: &app.TestCommand{
						Command: "echo 'Creating post with title: ${title}'",
					},
					Expect: app.TestExpect{
						Headers: map[string]string{
							"Content-Type": "application/json; charset=utf-8",
						},
						Status: toPointer(201),
					},
					Extract: []app.TestExtract{
						{
							Path: "id",
							As:   "postId",
						},
					},
				},
				{
					Name:        "Create comment",
					Description: "This sample test creates a new comment for the post created in the previous test",
					Depends:     []string{"Create post"},
					Fake: map[string]string{
						"comment_email": "{email}",
						"comment_name":  "{name}",
						"comment_body":  "{paragraph:1,3,40,\n}",
					},
					Request: &app.TestRequest{
						Method: "POST",
						URL:    "/comments",
						Body:   toPointer(`{"postId": "${postId}", "name": "${comment_name}", "email": "${comment_email}", "body": "${comment_body}"}`),
					},
					Expect: app.TestExpect{
						Status: toPointer(201),
					},
					Extract: []app.TestExtract{
						{
							Path: "id",
							As:   "commentId",
						},
					},
				},
				{
					Name:        "Get comments for post",
					Description: "This sample test retrieves comments for the post created in the previous test",
					Depends:     []string{"Create post", "Create comment"},
					Request: &app.TestRequest{
						Method: "GET",
						URL:    "/posts/1/comments",
					},
					Expect: app.TestExpect{
						Status: toPointer(200),
						Json: map[string]any{
							"#":        5,
							"0.postId": 1,
						},
					},
				},
			},
		}

		sampleTestFileData, err := yaml.Marshal(sampleTestFile)
		if err != nil {
			cmd.PrintErrf("Error generating sample tests: %v\n", err)
			return
		}

		err = os.WriteFile("tests/posts.yaml", sampleTestFileData, 0o644)
		if err != nil {
			cmd.PrintErrf("Error writing sample test file: %v\n", err)
			return
		}

		cmd.Printf("Project initialised successfully!\n")
	},
}

func toPointer[T any](v T) *T {
	return &v
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
