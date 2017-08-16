package cmd

import (
	"fmt"
	"qnmahjong/cache"
	"qnmahjong/cron"
	"qnmahjong/db"
	"qnmahjong/log"
	"qnmahjong/logic"
	"qnmahjong/notice"
	"qnmahjong/redis"
	"qnmahjong/util"

	"github.com/spf13/cobra"
)

// logicCmd represents the logic command
var logicCmd = &cobra.Command{
	Use:   "logic",
	Short: "Logic server for qnmahjong",
	Long:  `Logic server for qnmahjong.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("logic called")
		fmt.Printf("branch: %s\n", branch)
		fmt.Printf("commit: %s\n", commit)
		defer util.Stack()
		log.Config("logic")
		cron.Start("logic")
		db.Start()
		redis.Start()
		cache.InitCost()
		cache.InitGame()
		cache.InitShop()
		cache.InitPlayer()
		notice.StartLogicNotice()
		logic.Start()
	},
}

func init() {
	RootCmd.AddCommand(logicCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logicCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logicCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	logicCmd.Flags().StringP("host", "", "0.0.0.0", "The host to listen")
	logicCmd.Flags().Int32P("port", "", 5002, "The port to listen")
}
