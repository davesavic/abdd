/*
Copyright © 2025 Dave Savic
*/
package cmd

import (
	"github.com/davesavic/abdd/app"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		folders, err := cmd.Flags().GetStringSlice("folders")
		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			return
		}

		a, err := app.New(app.AbddArgs{
			ConfigFile: cmd.Flag("config").Value.String(),
			Folders:    folders,
		})
		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			return
		}

		err = a.Run()
		if err != nil {
			cmd.PrintErrf("Error: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringSliceP("folders", "f", []string{}, "Folders to run tests from")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
