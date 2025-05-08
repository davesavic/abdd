/*
Copyright © 2025 Dave Savic
*/
package cmd

import (
	"fmt"

	"github.com/davesavic/abdd/app"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		a, err := app.New(app.AbddArgs{
			ConfigFile: "./_examples/abdd.yaml",
			Folders:    []string{"./_examples/tests"},
		})
		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			return
		}

		fmt.Printf("Loaded %d tests\n", len(a.Tests))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
