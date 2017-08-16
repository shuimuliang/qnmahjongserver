package cmd

import (
	"fmt"
	"qnmahjong/cron"
	"qnmahjong/log"
	"qnmahjong/robot"
	"qnmahjong/util"

	"github.com/spf13/cobra"
)

// robotCmd represents the robot command
var robotCmd = &cobra.Command{
	Use:   "robot",
	Short: "Robot server for qnmahjong",
	Long:  `Robot server for qnmahjong.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("robot called")
		fmt.Printf("branch: %s\n", branch)
		fmt.Printf("commit: %s\n", commit)
		defer util.Stack()
		log.Config("robot")
		cron.Start("robot")
		robot.Start()
	},
}

func init() {
	RootCmd.AddCommand(robotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// robotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// robotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
