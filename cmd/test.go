package cmd

import (
	"fmt"
	"qnmahjong/cache"
	"qnmahjong/db"
	"os"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("test called")
		fmt.Printf("branch: %s\n", branch)
		fmt.Printf("commit: %s\n", commit)
		db.Start()
		// fmt.Println(dao.AgIDsFromAgAccount(db.Pool))
		// fmt.Println(dao.AgUpperIDsFromAgAuth(db.Pool))
		// fmt.Println(dao.StartTimesFromAgBill(db.Pool))
		// fmt.Println(dao.EmailsFromAccount(db.Pool))
		// fmt.Println(dao.MjTypesFromCost(db.Pool))
		// fmt.Println(dao.ChannelsFromGame(db.Pool))
		// fmt.Println(dao.ModulesFromModule(db.Pool))
		// fmt.Println(dao.PmsnTypesFromPermission(db.Pool))
		// fmt.Println(dao.RolesFromRole(db.Pool))
		// fmt.Println(dao.ChannelsFromShop(db.Pool))
		// fmt.Println(dao.PlayerIDsFromPlayer(db.Pool))
		cache.InitAgAccount()
		cache.InitAgBill()
		cache.InitAgAuth()
		cache.CashSettlement()
		os.Exit(-1)
	},
}

func init() {
	RootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
