package sale

import (
	"context"
	"fmt"
	"html/template"
	"qnmahjong/def"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	e *echo.Echo
)

// Start sale server
func Start() {
	// Echo instance
	e = echo.New()

	// Set static prefix
	e.Static("/", "admin")

	// Pre-compile templates
	t := &Template{
		templates: template.Must(template.ParseGlob("admin/sale/*.html")),
	}
	e.Renderer = t

	// Debug mode
	// e.Debug = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 登录
	e.GET("/", func(c echo.Context) error {
		return handleGetLogin(c)
	})

	e.GET("/login", func(c echo.Context) error {
		return handleGetLogin(c)
	})

	e.POST("/login", func(c echo.Context) error {
		return handlePostLogin(c)
	})

	// 注册
	e.GET("/register", func(c echo.Context) error {
		return handleGetRegister(c)
	})

	e.POST("/register", func(c echo.Context) error {
		return handlePostRegister(c)
	})

	// 主页
	e.GET("/index", func(c echo.Context) error {
		return handleGetIndex(c)
	})

	// 今日售出玉
	e.GET("/user-cards", func(c echo.Context) error {
		return handleGetUserCards(c)
	})

	e.POST("/user-cards", func(c echo.Context) error {
		return handlePostUserCards(c)
	})

	// 代理提成
	e.GET("/user-balance", func(c echo.Context) error {
		return handleGetUserBalance(c)
	})

	// 我的代理
	e.GET("/user-agents", func(c echo.Context) error {
		return handleGetUserAgents(c)
	})

	// 下级代理授权
	e.GET("/agent-auth", func(c echo.Context) error {
		return handleGetAgentAuth(c)
	})

	e.POST("/agent-auth", func(c echo.Context) error {
		return handlePostAgentAuth(c)
	})

	// 充值查询
	e.GET("/query-order", func(c echo.Context) error {
		return handleGetQueryOrder(c)
	})

	e.POST("/query-order", func(c echo.Context) error {
		return handlePostQueryOrder(c)
	})

	// 返现明细查询
	e.GET("/query-balance", func(c echo.Context) error {
		return handleGetQueryBalance(c)
	})

	// 我的代理
	e.GET("/my-agents", func(c echo.Context) error {
		return handleGetMyAgents(c)
	})

	// 我的玩家
	e.GET("/my-players", func(c echo.Context) error {
		return handleGetMyPlayers(c)
	})

	// 推广活动
	e.GET("/activities", func(c echo.Context) error {
		return handleGetActivities(c)
	})

	// 提现审核
	e.GET("/verify-balance", func(c echo.Context) error {
		return handleGetVerifyBalance(c)
	})

	// 提现审核详情
	e.GET("/verify-balance-detail", func(c echo.Context) error {
		return handleGetVerifyBalanceDetail(c)
	})

	e.POST("/verify-balance-detail", func(c echo.Context) error {
		return handlePostVerifyBalanceDetail(c)
	})

	// 个人信息
	e.GET("/edit-info", func(c echo.Context) error {
		return handleGetEditInfo(c)
	})

	// 个人资料
	e.GET("/edit-profile", func(c echo.Context) error {
		return handleGetEditProfile(c)
	})

	e.POST("/edit-profile", func(c echo.Context) error {
		return handlePostEditProfile(c)
	})

	// 修改密码
	e.GET("/edit-pwd", func(c echo.Context) error {
		return handleGetEditPWD(c)
	})

	e.POST("/edit-pwd", func(c echo.Context) error {
		return handlePostEditPWD(c)
	})

	// Start server
	err := e.Start(":5005")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrStartSale)
		fmt.Println("start sale failed")
		os.Exit(-1)
	}
}

// Shutdown sale server
func Shutdown() {
	if e != nil {
		// shut down gracefully, but wait no longer than 5 seconds before halting
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := e.Server.Shutdown(c)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error(def.ErrShutdownSale)
		}
		fmt.Println("sale shut down")
	}
}
