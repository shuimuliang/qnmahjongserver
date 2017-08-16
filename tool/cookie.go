package tool

import (
	"qnmahjong/cache"
	"net/http"
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

func checkCookies(c echo.Context) (email string, ok bool) {
	email, err := readCookie("toolEmail", c)
	if err != nil {
		return "", false
	}

	password, err := readCookie("toolPassword", c)
	if err != nil {
		return "", false
	}

	ok = cache.CheckAccount(email, password)
	return email, ok
}

func saveCookies(email, password string, c echo.Context) {
	writeCookie("toolEmail", email, c)
	writeCookie("toolPassword", password, c)
}
