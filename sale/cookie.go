package sale

import (
	"qnmahjong/cache"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

func writeCookie(name, value string, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = time.Now().AddDate(0, 0, 7)
	c.SetCookie(cookie)
}

func readCookie(name string, c echo.Context) (value string, err error) {
	cookie, err := c.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func checkCookies(c echo.Context) (id int32, ok bool) {
	agIDStr, err := readCookie("saleAgID", c)
	if err != nil {
		return 0, false
	}

	agID, _ := strconv.Atoi(agIDStr)
	password, err := readCookie("salePassword", c)
	if err != nil {
		return 0, false
	}

	ok = cache.CheckAgAgent(int32(agID), password)
	return int32(agID), ok
}

func saveCookies(agID, password string, c echo.Context) {
	writeCookie("saleAgID", agID, c)
	writeCookie("salePassword", password, c)
}
