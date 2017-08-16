package cmd

import (
	"fmt"
	"qnmahjong/cron"
	"qnmahjong/db"
	"qnmahjong/log"
	"qnmahjong/login"
	"qnmahjong/notice"
	"qnmahjong/util"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login server for qnmahjong",
	Long:  `Login server for qnmahjong.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("login called")
		fmt.Printf("branch: %s\n", branch)
		fmt.Printf("commit: %s\n", commit)
		defer util.Stack()
		log.Config("login")
		cron.Start("login")
		db.Start()
		notice.StartLoginNotice()
		login.Start()
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	loginCmd.Flags().StringP("host", "", "0.0.0.0", "The host to listen")
	loginCmd.Flags().Int32P("port", "", 5001, "The port to listen")
}
