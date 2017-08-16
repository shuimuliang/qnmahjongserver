package pay

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"qnmahjong/cache"
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/notice"
	"qnmahjong/pf"
	"qnmahjong/util"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	e           *echo.Echo
	iosIpay     *IpayHelper
	androidIpay *IpayHelper
	mu          sync.RWMutex
)

// Start pay server
func Start() {
	// Echo instance
	e = echo.New()

	// Set static prefix
	e.Static("/", "admin/pay")

	// Pre-compile templates
	t := &Template{
		templates: template.Must(template.ParseGlob("admin/pay/*.html")),
	}
	e.Renderer = t

	// ios Ipay instance
	iosIpay, err := NewIpayHelperWithPem(APP_ID_IOS,
		"./admin/pay/ipay_private_key_ios.pem", "./admin/pay/ipay_public_key_ios.pem")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrStartPay)
	}

	// android Ipay instance
	androidIpay, err := NewIpayHelperWithPem(APP_ID_ANDROID,
		"./admin/pay/ipay_private_key_android.pem", "./admin/pay/ipay_public_key_android.pem")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrStartPay)
	}

	// Debug mode
	// e.Debug = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ipay_order
	e.GET("/ipay_order", func(c echo.Context) (err error) {
		mu.RLock()
		mu.RUnlock()

		playerIDStr := c.QueryParam("player_id")
		channelStr := c.QueryParam("channel")
		waresIDStr := c.QueryParam("wares_id")

		playerID, err := strconv.Atoi(playerIDStr)
		if err != nil {
			return
		}

		channel, err := strconv.Atoi(channelStr)
		if err != nil {
			return
		}

		waresID, err := strconv.Atoi(waresIDStr)
		if err != nil {
			return
		}

		send := pf.OrderApplySend{
			PlayerID: int32(playerID),
			Channel:  int32(channel),
			WaresID:  int32(waresID),
		}
		util.LogSend(int32(pf.OrderApply), int32(playerID), 0, send, "OrderApply")

		order := cache.GetOrderWaresID(send.Channel, send.PlayerID, def.PayTypeIPay, waresIDStr)
		if order == nil {
			util.LogError(err, "order", order, send.PlayerID, def.ErrCreateDBOrder)
			return
		}

		err = order.Insert(db.Pool)
		if err != nil {
			util.LogError(err, "order", order, send.PlayerID, def.ErrInsertOrder)
			return
		}

		var ipay *IpayHelper
		var appid string
		switch int32(channel) {
		case def.ChannelIOSHB, def.ChannelIOSHN:
			ipay = iosIpay
			appid = APP_ID_IOS
		case def.ChannelAndroidHB, def.ChannelAndroidHN:
			ipay = androidIpay
			appid = APP_ID_ANDROID
		}

		cpOrderId := order.OrderID
		price := float32(order.Price) / 100
		transId, err := ipay.CreateIpayOrder(waresID, order.WaresName, cpOrderId, price, playerIDStr, "", "")
		if err != nil {
			util.LogError(err, "order", order, send.PlayerID, def.ErrCreateIpayOrder)
			return
		}

		order.TransID = transId
		util.LogSend(int32(pf.OrderApply), send.PlayerID, 0, order, "OrderApply")

		err = order.Update(db.Pool)
		if err != nil {
			util.LogError(err, "order", order, send.PlayerID, def.ErrUpdateOrder)
			return
		}

		payUrl, err := ipay.GetNewHtml5RedirectUrl(transId, appid, "", "")
		if err != nil {
			util.LogError(err, "order", order, send.PlayerID, def.ErrGetHTML5RedirectURL)
			return
		}

		return c.Render(http.StatusOK, "ipay_order", payUrl)
	})

	// ipay_notice
	e.POST("/ipay_notice/ios", func(c echo.Context) error {
		mu.RLock()
		mu.RUnlock()

		msg, _ := ioutil.ReadAll(c.Request().Body)
		transdata, err := iosIpay.ParseNotifyInfo(msg)
		if err != nil {
			return c.String(http.StatusOK, "FAILED")
		}

		util.LogRecv(int32(pf.OrderApply), 0, 0, *transdata, "OrderApply")
		order, err := dao.OrderByOrderID(db.Pool, transdata.CpOrderId)
		if err != nil {
			return c.String(http.StatusOK, "FAILED")
		}

		order.Status = int32(transdata.Result)
		order.ReviseTime = time.Now()
		err = order.Update(db.Pool)
		if err != nil {
			util.LogError(err, "order", order, order.PlayerID, def.ErrUpdateOrder)
			return c.String(http.StatusOK, "FAILED")
		}

		result, err := iosIpay.QueryResult(transdata.CpOrderId)
		if err != nil {
			util.LogError(err, "transdata", transdata, order.PlayerID, def.ErrQueryIpayResult)
			return c.String(http.StatusOK, "FAILED")
		}

		if result.Result != int(def.IpayStatusSuccess) {
			return c.String(http.StatusOK, "FAILED")
		}

		notice.LogicHandleOrder(order.OrderID)
		return c.String(http.StatusOK, "SUCCESS")
	})

	// ipay_notice
	e.POST("/ipay_notice/android", func(c echo.Context) error {
		mu.RLock()
		mu.RUnlock()

		msg, _ := ioutil.ReadAll(c.Request().Body)
		transdata, err := androidIpay.ParseNotifyInfo(msg)
		if err != nil {
			return c.String(http.StatusOK, "FAILED")
		}

		util.LogRecv(int32(pf.OrderApply), 0, 0, *transdata, "OrderApply")
		order, err := dao.OrderByOrderID(db.Pool, transdata.CpOrderId)
		if err != nil {
			return c.String(http.StatusOK, "FAILED")
		}

		order.Status = int32(transdata.Result)
		order.ReviseTime = time.Now()
		err = order.Update(db.Pool)
		if err != nil {
			util.LogError(err, "order", order, order.PlayerID, def.ErrUpdateOrder)
			return c.String(http.StatusOK, "FAILED")
		}

		result, err := androidIpay.QueryResult(transdata.CpOrderId)
		if err != nil {
			util.LogError(err, "transdata", transdata, order.PlayerID, def.ErrQueryIpayResult)
			return c.String(http.StatusOK, "FAILED")
		}

		if result.Result != int(def.IpayStatusSuccess) {
			return c.String(http.StatusOK, "FAILED")
		}

		notice.LogicHandleOrder(order.OrderID)
		return c.String(http.StatusOK, "SUCCESS")
	})

	// Start server
	err = e.Start(":5004")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrStartPay)
		fmt.Println("start pay failed")
		os.Exit(-1)
	}
}

// Shutdown tool server
func Shutdown() {
	if e != nil {
		// shut down gracefully, but wait no longer than 5 seconds before halting
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := e.Server.Shutdown(c)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error(def.ErrShutdownPay)
		}
		fmt.Println("pay shut down")
	}
}
