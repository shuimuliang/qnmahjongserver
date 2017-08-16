package tool

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

// Start tool server
func Start() {
	// Echo instance
	e = echo.New()

	// Set static prefix
	e.Static("/", "admin")

	// Pre-compile templates
	t := &Template{
		templates: template.Must(template.ParseGlob("admin/tool/*.html")),
	}
	e.Renderer = t

	// Debug mode
	// e.Debug = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 版本检查
	e.POST("/check_version", func(c echo.Context) error {
		return handlePostCheckVersion(c)
	})

	// 刷新配置
	e.GET("/refresh_config", func(c echo.Context) error {
		return handleGetRefreshConfig(c)
	})

	// 河南麻将下载
	e.GET("/download_hnmj", func(c echo.Context) error {
		return handleGetHnmj(c)
	})

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

	// 主页
	e.GET("/index", func(c echo.Context) error {
		return handleGetIndex(c)
	})

	// 游戏配置
	e.GET("/games-manage", func(c echo.Context) error {
		return handleGetGamesManage(c)
	})

	e.POST("/games-manage", func(c echo.Context) error {
		return handlePostGamesManage(c)
	})

	// 房卡配置
	e.GET("/costs-manage", func(c echo.Context) error {
		return handleGetCostsManage(c)
	})

	e.POST("/costs-manage", func(c echo.Context) error {
		return handlePostCostsManage(c)
	})

	// 模块配置
	e.GET("/modules-manage", func(c echo.Context) error {
		return handleGetModulesManage(c)
	})

	e.POST("/modules-manage", func(c echo.Context) error {
		return handlePostModulesManage(c)
	})

	// 商品配置
	e.GET("/shops-manage", func(c echo.Context) error {
		return handleGetShopsManage(c)
	})

	e.POST("/shops-manage", func(c echo.Context) error {
		return handlePostShopsManage(c)
	})

	// 账号管理
	e.GET("/accounts-manage", func(c echo.Context) error {
		return handleGetAccountsManage(c)
	})

	e.POST("/accounts-manage", func(c echo.Context) error {
		return handlePostAccountsManage(c)
	})

	// 角色管理
	e.GET("/roles-manage", func(c echo.Context) error {
		return handleGetRolesManage(c)
	})

	e.POST("/roles-manage", func(c echo.Context) error {
		return handlePostRolesManage(c)
	})

	// 权限管理
	e.GET("/permissions-manage", func(c echo.Context) error {
		return handleGetPermissionsManage(c)
	})

	e.POST("/permissions-manage", func(c echo.Context) error {
		return handlePostPermissionsManage(c)
	})

	// Start server
	err := e.Start(":5003")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrStartTool)
		fmt.Println("start tool failed")
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
			}).Error(def.ErrShutdownTool)
		}
		fmt.Println("tool shut down")
	}
}
