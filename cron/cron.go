package cron

import (
	"fmt"
	"qnmahjong/cache"
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/log"
	"qnmahjong/notice"

	"github.com/Sirupsen/logrus"
	crontab "github.com/robfig/cron"
)

var (
	c *crontab.Cron
)

// Start crontab server
func Start(server string) {
	c = crontab.New()
	err := c.AddFunc("@daily", func() {
		log.Config(server)
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error(def.ErrAddCronFunc)
		return
	}

	c.Start()
}

// StartForSale ...
func StartForSale(server string) {
	c = crontab.New()
	err := c.AddFunc("@daily", func() {
		log.Config(server)
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error(def.ErrAddCronFunc)
		return
	}

	// 统计上周结算
	err = c.AddFunc("0 0 0 * * 1", func() {
		cache.CashSettlement()
		notice.SaleInitAgBill()
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error(def.ErrAddCronFunc)
		return
	}

	// 不足打款金额结余
	err = c.AddFunc("0 0 0 * * 2", func() {
		dao.UpdateAgBillStatus(db.Pool)
		notice.SaleInitAgBill()
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error(def.ErrAddCronFunc)
		return
	}

	// 代理等级提升
	// err = c.AddFunc("0 0 0 1 * *", func() {
	// })

	c.Start()
}

// Shutdown crontab server
func Shutdown() {
	if c != nil {
		c.Stop()
		fmt.Println("cron shut down")
	}
}
