package sale

import (
	"qnmahjong/cache"
	"qnmahjong/def"
	"qnmahjong/notice"
	"qnmahjong/util"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// 登录
func handleGetLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func handlePostLogin(c echo.Context) error {
	agIDStr := c.FormValue("ag_id")
	password := c.FormValue("password")
	agID, _ := strconv.Atoi(agIDStr)
	sha1Password := util.Sha1Password(password)
	ok := cache.CheckAgAgent(int32(agID), sha1Password)
	if !ok {
		return c.Render(http.StatusOK, "failed", "账号或密码错误，请重新登录！")
	}

	saveCookies(agIDStr, sha1Password, c)
	return c.Redirect(http.StatusMovedPermanently, "/index")
}

// 注册
func handleGetRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "register", nil)
}

func handlePostRegister(c echo.Context) error {
	agUpperIDStr := c.FormValue("ag_upper_id")
	agIDStr := c.FormValue("ag_id")
	telephone := c.FormValue("telephone")
	password := c.FormValue("password")
	passwordEnsure := c.FormValue("password_ensure")
	agUpperID, _ := strconv.Atoi(agUpperIDStr)
	agID, _ := strconv.Atoi(agIDStr)
	if password != passwordEnsure {
		return c.Render(http.StatusOK, "failed", "注册失败")
	}

	ok := cache.CreateAgAgent(int32(agUpperID), int32(agID), password, telephone)
	if !ok {
		return c.Render(http.StatusOK, "failed", "注册失败")
	}

	notice.SaleInitAgAccount()
	return c.Render(http.StatusOK, "init-login", "注册成功，请点此处登录你的账号！")
}

// 主页
func handleGetIndex(c echo.Context) error {
	agID, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rd := cache.GetSaleIndexRD(agID)
	return c.Render(http.StatusOK, "index", rd)
}

// 今日售出玉
func handleGetUserCards(c echo.Context) error {
	return c.Render(http.StatusOK, "user-cards", nil)
}

func handlePostUserCards(c echo.Context) error {
	startDate := c.FormValue("start-date")
	endDate := c.FormValue("end-date")
	agID, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rd := cache.GetUserCardsRD(agID, startDate, endDate)
	return c.Render(http.StatusOK, "user-cards-deatil", rd)
}

// 代理提成
func handleGetUserBalance(c echo.Context) error {
	agID, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rd := cache.GetUserBalanceRD(agID)
	return c.Render(http.StatusOK, "user-balance", rd)
}

// 我的代理
func handleGetUserAgents(c echo.Context) error {
	return c.Render(http.StatusOK, "user-agents", nil)
}

// 下级代理授权
func handleGetAgentAuth(c echo.Context) error {
	agID, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rd := cache.GetAgentAuthRD(agID)
	return c.Render(http.StatusOK, "agent-auth", rd)
}

func handlePostAgentAuth(c echo.Context) error {
	agIDStr := c.FormValue("ag_id")
	agLevelStr := c.FormValue("ag_level")
	agID, _ := strconv.Atoi(agIDStr)
	agLevel, _ := strconv.Atoi(agLevelStr)
	agUpperID, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	ok = cache.CreateAgAuth(int32(agUpperID), int32(agID), int32(agLevel))
	if !ok {
		return c.Render(http.StatusOK, "failed", "授权失败")
	}

	notice.SaleInitAgAuth()
	return c.Render(http.StatusOK, "success", "授权成功")
}

// 充值查询
func handleGetQueryOrder(c echo.Context) error {
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	return c.Render(http.StatusOK, "query-order", nil)
}

func handlePostQueryOrder(c echo.Context) error {
	playerIDStr := c.FormValue("player-id")
	startDate := c.FormValue("start-date")
	endDate := c.FormValue("end-date")
	playerID, _ := strconv.Atoi(playerIDStr)
	agID, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rd := cache.GetQueryOrderRD(agID, int32(playerID), startDate, endDate)
	return c.Render(http.StatusOK, "query-order-detail", rd)
}

// 返现明细查询
func handleGetQueryBalance(c echo.Context) error {
	agID, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rds := cache.GetQueryBalanceRD(agID)
	return c.Render(http.StatusOK, "query-balance", rds)
}

// 我的代理
func handleGetMyAgents(c echo.Context) error {
	agID, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rds := cache.GetMyAgentsRD(agID)
	return c.Render(http.StatusOK, "my-agents", rds)
}

// 我的玩家
func handleGetMyPlayers(c echo.Context) error {
	agID, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rds := cache.GetMyPlayersRD(agID)
	return c.Render(http.StatusOK, "my-players", rds)
}

// 推广活动
func handleGetActivities(c echo.Context) error {
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	return c.Render(http.StatusOK, "activities", nil)
}

// 提现审核
func handleGetVerifyBalance(c echo.Context) error {
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rds := cache.GetVerifyBalanceRD()
	return c.Render(http.StatusOK, "verify-balance", rds)
}

// 提现审核详情
func handleGetVerifyBalanceDetail(c echo.Context) error {
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	start := c.QueryParam("start")
	rds := cache.GetVerifyBalanceDetailRD(start)
	return c.Render(http.StatusOK, "verify-balance-detail", rds)
}

func handlePostVerifyBalanceDetail(c echo.Context) error {
	indexIDStr := c.FormValue("index-id")
	agIDStr := c.FormValue("ag-id")
	indexID, _ := strconv.Atoi(indexIDStr)
	agID, _ := strconv.Atoi(agIDStr)
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	ok = cache.UpdateAgBillStatus(int32(agID), int32(indexID), def.AgBillStatusYidakuan)
	if !ok {
		return c.Render(http.StatusOK, "failed", "打款失败")
	}

	return c.Render(http.StatusOK, "success", "打款成功")
}

// 个人信息
func handleGetEditInfo(c echo.Context) error {
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	return c.Render(http.StatusOK, "edit-info", nil)
}

// 个人资料
func handleGetEditProfile(c echo.Context) error {
	agID, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	rd := cache.GetEditProfileRD(agID)
	return c.Render(http.StatusOK, "edit-profile", rd)
}

func handlePostEditProfile(c echo.Context) error {
	telephone := c.FormValue("telephone")
	realname := c.FormValue("realname")
	weixin := c.FormValue("weixin")
	alipay := c.FormValue("alipay")
	email := c.FormValue("email")
	agID, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	ok = cache.EditAgAgentProfile(agID, telephone, realname, weixin, alipay, email)
	if !ok {
		return c.Render(http.StatusOK, "failed", "修改个人资料失败")
	}

	return c.Render(http.StatusOK, "success", "修改个人资料成功")
}

// 修改密码
func handleGetEditPWD(c echo.Context) error {
	_, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	return c.Render(http.StatusOK, "edit-pwd", nil)
}

func handlePostEditPWD(c echo.Context) error {
	password := c.FormValue("password")
	newPassword := c.FormValue("new-password")
	newPasswordEnsure := c.FormValue("new-password_ensure")
	sha1Password := util.Sha1Password(password)
	agID, ok := checkCookies(c)
	if !ok {
		return c.Render(http.StatusOK, "relogin", nil)
	}

	ok = cache.CheckAgAgent(agID, sha1Password)
	if !ok {
		return c.Render(http.StatusOK, "failed", "原始密码不正确")
	}

	if newPassword != newPasswordEnsure {
		return c.Render(http.StatusOK, "failed", "两次输入密码不一致")
	}

	ok = cache.EditAgAgentPWD(agID, util.Sha1Password(newPassword))
	if !ok {
		return c.Render(http.StatusOK, "failed", "修改密码失败")
	}

	return c.Render(http.StatusOK, "init-login", "修改密码成功，请点此处重新登录！")
}
