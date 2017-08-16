package tool

import (
	"encoding/json"
	"io/ioutil"
	"qnmahjong/cache"
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/notice"
	"qnmahjong/util"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

// 版本检查
func handlePostCheckVersion(c echo.Context) error {
	channel, err := strconv.Atoi(c.QueryParam("channel"))
	if err != nil {
		return c.JSON(http.StatusOK, nil)
	}

	version, err := strconv.Atoi(c.QueryParam("version"))
	if err != nil {
		return c.JSON(http.StatusOK, nil)
	}

	deviceID := c.QueryParam("deviceID")
	versionCheck := cache.GetVersionCheck(int32(channel), int32(version))
	log.WithFields(log.Fields{
		"channel":      channel,
		"version":      version,
		"deviceID":     deviceID,
		"versionCheck": versionCheck,
	}).Info("VersionCheck")
	return c.JSON(http.StatusOK, versionCheck)
}

// 刷新配置
func handleGetRefreshConfig(c echo.Context) error {
	// logic
	notice.LogicInitCost()
	notice.LogicInitGame()
	notice.LogicInitShop()
	// pay
	notice.PayInitShop()
	return c.Render(http.StatusOK, "success", "刷新成功")
}

// 河南麻将下载页
func handleGetHnmj(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "http://oow62av5c.bkt.clouddn.com/hnmj.html")
}

// 登录
func handleGetLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func handlePostLogin(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	sha1Password := util.Sha1Password(password)
	ok := cache.CheckAccount(email, sha1Password)
	if !ok {
		return c.Render(http.StatusOK, "failed", "账号或密码错误，请重新登录！")
	}

	saveCookies(email, sha1Password, c)
	return c.Redirect(http.StatusMovedPermanently, "/index")
}

// 主页
func handleGetIndex(c echo.Context) error {
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	return c.Render(http.StatusOK, "index", nil)
}

// 游戏配置
func handleGetGamesManage(c echo.Context) error {
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rd, _ := dao.SelectAllGames(db.Pool)
	return c.Render(http.StatusOK, "games-manage", rd)
}

func handlePostGamesManage(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	var games struct {
		Data []*dao.Game `json:"data" form:"data" query:"data"`
	}

	err := json.Unmarshal(body, &games)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"data":  string(body),
		}).Error("Unmarshal tool games")
		return c.String(http.StatusOK, "failed")
	}

	handleGamesManage(games.Data)
	return c.String(http.StatusOK, "success")
}

// 房卡配置
func handleGetCostsManage(c echo.Context) error {
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rd, _ := dao.SelectAllCosts(db.Pool)
	return c.Render(http.StatusOK, "costs-manage", rd)
}

func handlePostCostsManage(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	var costs struct {
		Data []*dao.Cost `json:"data" form:"data" query:"data"`
	}

	err := json.Unmarshal(body, &costs)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"data":  string(body),
		}).Error("Unmarshal tool costs")
		return c.String(http.StatusOK, "failed")
	}

	handleCostsManage(costs.Data)
	return c.String(http.StatusOK, "success")
}

// 模块配置
func handleGetModulesManage(c echo.Context) error {
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rd, _ := dao.SelectAllModules(db.Pool)
	return c.Render(http.StatusOK, "modules-manage", rd)
}

func handlePostModulesManage(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	var modules struct {
		Data []*dao.Module `json:"data" form:"data" query:"data"`
	}

	err := json.Unmarshal(body, &modules)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"data":  string(body),
		}).Error("Unmarshal tool modules")
		return c.String(http.StatusOK, "failed")
	}

	handleModulesManage(modules.Data)
	return c.String(http.StatusOK, "success")
}

// 商品配置
func handleGetShopsManage(c echo.Context) error {
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rd, _ := dao.SelectAllShops(db.Pool)
	return c.Render(http.StatusOK, "shops-manage", rd)
}

func handlePostShopsManage(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	var shops struct {
		Data []*dao.Shop `json:"data" form:"data" query:"data"`
	}

	err := json.Unmarshal(body, &shops)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"data":  string(body),
		}).Error("Unmarshal tool shops")
		return c.String(http.StatusOK, "failed")
	}

	handleShopsManage(shops.Data)
	return c.String(http.StatusOK, "success")
}

// 账号管理
func handleGetAccountsManage(c echo.Context) error {
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rd, _ := dao.SelectAllAccounts(db.Pool)
	return c.Render(http.StatusOK, "accounts-manage", rd)
}

func handlePostAccountsManage(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	var accounts struct {
		Data []*dao.Account `json:"data" form:"data" query:"data"`
	}

	err := json.Unmarshal(body, &accounts)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"data":  string(body),
		}).Error("Unmarshal tool accounts")
		return c.String(http.StatusOK, "failed")
	}

	handleAccountsManage(accounts.Data)
	return c.String(http.StatusOK, "success")
}

// 角色管理
func handleGetRolesManage(c echo.Context) error {
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rd, _ := dao.SelectAllRoles(db.Pool)
	return c.Render(http.StatusOK, "roles-manage", rd)
}

func handlePostRolesManage(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	var roles struct {
		Data []*dao.Role `json:"data" form:"data" query:"data"`
	}

	err := json.Unmarshal(body, &roles)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"data":  string(body),
		}).Error("Unmarshal tool roles")
		return c.String(http.StatusOK, "failed")
	}

	handleRolesManage(roles.Data)
	return c.String(http.StatusOK, "success")
}

// 权限管理
func handleGetPermissionsManage(c echo.Context) error {
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rd, _ := dao.SelectAllPermissions(db.Pool)
	return c.Render(http.StatusOK, "permissions-manage", rd)
}

func handlePostPermissionsManage(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	var permissions struct {
		Data []*dao.Permission `json:"data" form:"data" query:"data"`
	}

	err := json.Unmarshal(body, &permissions)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"data":  string(body),
		}).Error("Unmarshal tool permissions")
		return c.String(http.StatusOK, "failed")
	}

	handlePermissionsManage(permissions.Data)
	return c.String(http.StatusOK, "success")
}
