package main

import (
	"fmt"
	"qnmahjong/cmd"
	"qnmahjong/cron"
	"qnmahjong/db"
	"qnmahjong/logic"
	"qnmahjong/login"
	"qnmahjong/pay"
	"qnmahjong/redis"
	"qnmahjong/sale"
	"qnmahjong/tool"
	"qnmahjong/util"
	"os"
	"os/signal"
	"syscall"
)

var (
	branch = "v0.0.1"
	commit = "not set"
)

func main() {
	defer util.Stack()
	cmd.SetBranchCommit(branch, commit)
	go cmd.Execute()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done
	fmt.Println("exiting...")
	login.Shutdown()
	logic.Shutdown()
	tool.Shutdown()
	pay.Shutdown()
	sale.Shutdown()
	cron.Shutdown()
	db.Shutdown()
	redis.Shutdown()
	fmt.Println("exited!")
}
