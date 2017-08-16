package login

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"qnmahjong/def"
	"qnmahjong/msg"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	e *echo.Echo
)

// Start login server
func Start() {
	// Echo instance
	e = echo.New()

	// Debug mode
	// e.Debug = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login Http Routes => Handler
	e.POST("/", func(c echo.Context) error {
		m, _ := ioutil.ReadAll(c.Request().Body)
		recv, err := msg.LoginHandle(m, c)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error(def.ErrHandleLogin)
			return nil
		}
		return c.HTML(http.StatusOK, base64.StdEncoding.EncodeToString(recv))
	})

	// Start server
	err := e.Start(":5001")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrStartLogin)
		fmt.Println("start login failed")
		os.Exit(-1)
	}
}

// Shutdown login server
func Shutdown() {
	if e != nil {
		// shut down gracefully, but wait no longer than 5 seconds before halting
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := e.Server.Shutdown(c)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error(def.ErrShutdownLogin)
		}
		fmt.Println("login shut down")
	}
}
