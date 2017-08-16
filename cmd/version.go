package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	branch = "v0.0.1"
	commit = "not set"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version server for qnmahjong",
	Long:  `Version server for qnmahjong.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("version called")
		fmt.Printf("branch: %s\n", branch)
		fmt.Printf("commit: %s\n", commit)
		os.Exit(-1)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// SetBranchCommit get branch and commit from git
func SetBranchCommit(b, c string) {
	branch = b
	commit = c
}
