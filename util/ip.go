package util

import (
	"net"
	"net/http"

	"github.com/labstack/echo"
)

func RealIP(request *http.Request) string {
	ra := request.RemoteAddr
	if ip := request.Header.Get(echo.HeaderXForwardedFor); ip != "" {
		ra = ip
	} else if ip := request.Header.Get(echo.HeaderXRealIP); ip != "" {
		ra = ip
	} else {
		ra, _, _ = net.SplitHostPort(ra)
	}
	return ra

}
